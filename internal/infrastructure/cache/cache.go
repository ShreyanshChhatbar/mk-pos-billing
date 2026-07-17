package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/justkeval/go-phpserialize"
	"github.com/redis/go-redis/v9"
)

var ErrKeyNotFound = errors.New("key not found")

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
	res, err := s.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", ErrKeyNotFound
		}
		return "", err
	}
	return res, nil
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
		if errors.Is(err, redis.Nil) {
			return ErrKeyNotFound
		}
		return err
	}
	return json.Unmarshal(data, dest)
}

func (s *Service) SetPHPSerialized(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	serialized, err := phpserialize.Marshal(value)
	if err != nil {
		return err
	}
	return s.client.Set(ctx, key, serialized, ttl).Err()
}

func (s *Service) GetPHPSerialized(ctx context.Context, key string, dest interface{}) error {
	data, err := s.client.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return ErrKeyNotFound
		}
		return err
	}
	return phpserialize.Unmarshal(data, dest)
}
