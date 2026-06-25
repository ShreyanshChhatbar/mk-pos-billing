package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
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

func (s *Service) Add(ctx context.Context, key string, value interface{}, ttl time.Duration) (bool, error) {
	return s.client.SetNX(ctx, key, value, ttl).Result()
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
	// 1. Try PHP Unmarshal
	var raw interface{}
	var phpErr error

	// Detect PHP type and use specific unmarshaler
	if len(data) > 2 && data[1] == ':' {
		switch data[0] {
		case 'a', 'O': // Array or Object
			m, err := phpserialize.UnmarshalAssociativeArray(data)
			if err == nil {
				raw = phpValueToJSONCompatible(m)
			} else {
				phpErr = err
			}
		case 's':
			raw, phpErr = phpserialize.UnmarshalString(data)
		case 'i':
			raw, phpErr = phpserialize.UnmarshalInt(data)
		case 'd':
			raw, phpErr = phpserialize.UnmarshalFloat(data)
		case 'b':
			raw, phpErr = phpserialize.UnmarshalBool(data)
		}
	}

	if phpErr == nil && raw != nil {
		// JSON Bridge: Safe population behavior
		tempJSON, jerr := json.Marshal(raw)
		if jerr == nil {
			if jerry := json.Unmarshal(tempJSON, dest); jerry == nil {
				return nil
			}
		}
	}

	// Fallback to standard JSON parsing
	return json.Unmarshal(data, dest)
}

// phpValueToJSONCompatible converts phpserialize output into values that can be
// safely marshaled to JSON without losing array shape or panicking on non-string keys.
func phpValueToJSONCompatible(value interface{}) interface{} {
	switch v := value.(type) {
	case map[interface{}]interface{}:
		return phpMapToJSONCompatible(v)
	case []interface{}:
		out := make([]interface{}, len(v))
		for i, elem := range v {
			out[i] = phpValueToJSONCompatible(elem)
		}
		return out
	default:
		return value
	}
}

func phpMapToJSONCompatible(m map[interface{}]interface{}) interface{} {
	if len(m) == 0 {
		return map[string]interface{}{}
	}

	indexes := make([]int, 0, len(m))
	indexedValues := make(map[int]interface{}, len(m))

	for k, v := range m {
		index, ok := phpArrayIndex(k)
		if !ok {
			result := make(map[string]interface{}, len(m))
			for key, val := range m {
				result[fmt.Sprint(key)] = phpValueToJSONCompatible(val)
			}
			return result
		}

		indexes = append(indexes, index)
		indexedValues[index] = phpValueToJSONCompatible(v)
	}

	sort.Ints(indexes)
	for i, index := range indexes {
		if index != i {
			result := make(map[string]interface{}, len(m))
			for key, val := range m {
				result[fmt.Sprint(key)] = phpValueToJSONCompatible(val)
			}
			return result
		}
	}

	out := make([]interface{}, len(m))
	for index, value := range indexedValues {
		out[index] = value
	}
	return out
}

func phpArrayIndex(key interface{}) (int, bool) {
	maxInt := int(^uint(0) >> 1)

	switch v := key.(type) {
	case int:
		return v, v >= 0
	case int8:
		return int(v), v >= 0
	case int16:
		return int(v), v >= 0
	case int32:
		return int(v), v >= 0
	case int64:
		if v < 0 || v > int64(maxInt) {
			return 0, false
		}
		return int(v), true
	case uint:
		if v > uint(maxInt) {
			return 0, false
		}
		return int(v), true
	case uint8:
		return int(v), true
	case uint16:
		return int(v), true
	case uint32:
		if uint64(v) > uint64(maxInt) {
			return 0, false
		}
		return int(v), true
	case uint64:
		if v > uint64(maxInt) {
			return 0, false
		}
		return int(v), true
	default:
		return 0, false
	}
}
