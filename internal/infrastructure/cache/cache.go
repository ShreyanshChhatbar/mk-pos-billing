package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/elliotchance/phpserialize"
	"github.com/redis/go-redis/v9"
)

// Service wraps Redis operations with namespaced key helpers.
type Service struct {
	client    *redis.Client
	namespace string
}

func NewCacheService(client *redis.Client, namespace string) *Service {
	return &Service{client: client, namespace: namespace}
}

// BuildKey constructs a namespaced cache key: namespace:part1:part2...
func (s *Service) BuildKey(parts ...string) string {
	all := append([]string{s.namespace}, parts...)
	return strings.Join(all, "")
}

// GetTagVersion attempts to find the current version hash for a given tag.
// It first looks for a key ending in :key, then falls back to discovering it from the :entries set.
func (s *Service) GetTagVersion(ctx context.Context, tag string) (string, error) {
	// Try standard Laravel version key format: tag:TAGNAME:key
	versionKey := s.BuildKey("tag:", tag, ":key")
	version, err := s.client.Get(ctx, versionKey).Result()
	if err == nil && version != "" {
		return version, nil
	}

	// Fallback: Discover from :entries if the version key is missing
	entriesKey := s.BuildKey("tag:", tag, ":entries")
	members, err := s.client.ZRange(ctx, entriesKey, 0, 0).Result()
	if err == nil && len(members) > 0 {
		// Laravel stores entries as "hash:original_key"
		parts := strings.Split(members[0], ":")
		if len(parts) > 1 {
			return parts[0], nil
		}
	}

	return "", fmt.Errorf("tag version not found for tag: %s", tag)
}

// BuildTaggedKey constructs a key that includes the Laravel tag version/hash.
func (s *Service) BuildTaggedKey(ctx context.Context, tag string, parts ...string) string {
	version, err := s.GetTagVersion(ctx, tag)
	if err != nil {
		return s.BuildKey(parts...)
	}

	// Prepend hash followed by colon to the key parts
	allParts := append([]string{version + ":"}, parts...)
	return s.BuildKey(allParts...)
}

func (s *Service) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return s.client.Set(ctx, key, value, ttl).Err()
}

func (s *Service) Get(ctx context.Context, key string) (string, error) {
	return s.client.Get(ctx, key).Result()
}

func (s *Service) Delete(ctx context.Context, key string) error {
	return s.client.Del(ctx, key).Err()
}

// SetJSON marshals a value as JSON and sets it with TTL.
func (s *Service) SetJSON(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return s.client.Set(ctx, key, b, ttl).Err()
}

// GetJSON fetches a JSON value and unmarshals into dest.
func (s *Service) GetJSON(ctx context.Context, key string, dest interface{}) error {
	data, err := s.client.Get(ctx, key).Bytes()
	if err != nil {
		fmt.Println("Redis Error: ", err)
		return err
	}

	// First try to parse as PHP serialized data since Laravel uses this
	phpDest := make(map[interface{}]interface{})
	if err := phpserialize.Unmarshal(data, &phpDest); err == nil {
		// Convert the map[interface{}]interface{} to map[string]interface{}
		// JSON cannot encode map[interface{}]interface{} natively
		stringMap := convertToStringMap(phpDest)

		// Re-encode to JSON so we can unmarshal it into the user's strongly typed struct
		jsonData, err := json.Marshal(stringMap)

		if err == nil {
			return json.Unmarshal(jsonData, dest)
		}
	}

	// Fallback to standard JSON parsing if it's not a PHP serialized string
	return json.Unmarshal(data, dest)
}

// convertToStringMap recursively converts map[interface{}]interface{} to map[string]interface{}
// This is necessary because the JSON marshaler complains about map[interface{}] keys
func convertToStringMap(m map[interface{}]interface{}) map[string]interface{} {
	res := make(map[string]interface{})
	for k, v := range m {
		keyStr := fmt.Sprintf("%v", k)

		switch val := v.(type) {
		case map[interface{}]interface{}:
			res[keyStr] = convertToStringMap(val)
		case []interface{}:
			for i, elem := range val {
				if mElem, ok := elem.(map[interface{}]interface{}); ok {
					val[i] = convertToStringMap(mElem)
				}
			}
			res[keyStr] = val
		default:
			res[keyStr] = v
		}
	}
	return res
}
