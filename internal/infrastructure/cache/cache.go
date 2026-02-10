package cache

import (
	"context"
	"encoding/json"
	"strings"
	"time"

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
	return strings.Join(all, ":")
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
		return err
	}
	return json.Unmarshal(data, dest)
}
