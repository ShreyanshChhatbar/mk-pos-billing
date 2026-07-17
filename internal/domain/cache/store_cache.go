package cache

import (
	"context"
	"fmt"
	"mk-pos-billing/internal/infrastructure/cache"
	"mk-pos-billing/internal/infrastructure/config"
)

type StoreCache interface {
	Get(ctx context.Context, storeID int) (*Store, error)
}

type storeCache struct {
	cacheSvc *cache.Service
	prefixes *config.CachePrefixes
	tags     *config.CacheTags
}

func NewStoreCache(cacheSvc *cache.Service, prefixes *config.CachePrefixes, tags *config.CacheTags) StoreCache {
	return &storeCache{
		cacheSvc: cacheSvc,
		prefixes: prefixes,
		tags:     tags,
	}
}

func (s *storeCache) Get(ctx context.Context, storeID int) (*Store, error) {
	if storeID == 0 {
		return nil, nil // Equivalant to Laravel return null
	}

	// Use BuildTaggedKey to handle Laravel's cache tagging (prepending the version hash)
	cacheKey := s.cacheSvc.BuildTaggedKey(ctx, s.tags.Store, s.prefixes.Store, fmt.Sprint(storeID))

	var storeData Store
	err := s.cacheSvc.GetPHPSerialized(ctx, cacheKey, &storeData)

	if err != nil {
		return nil, err
	}

	return &storeData, nil
}
