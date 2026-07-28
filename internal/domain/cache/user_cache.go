package cache

import (
	"context"
	"fmt"
	"mk-pos-billing/internal/infrastructure/cache"
	"mk-pos-billing/internal/infrastructure/config"
)

type UserCache interface {
	Get(ctx context.Context, userID int) (*User, error)
}

type userCache struct {
	cacheSvc *cache.Service
	prefixes *config.CachePrefixes
}

func NewUserCache(cacheSvc *cache.Service, prefixes *config.CachePrefixes) UserCache {
	return &userCache{
		cacheSvc: cacheSvc,
		prefixes: prefixes,
	}
}

func (s *userCache) Get(ctx context.Context, userID int) (*User, error) {
	cacheKey := s.cacheSvc.BuildKey(s.prefixes.User, fmt.Sprint(userID))

	var userData User
	err := s.cacheSvc.GetPHPSerialized(ctx, cacheKey, &userData)

	if err != nil {
		return nil, err
	}

	return &userData, nil
}
