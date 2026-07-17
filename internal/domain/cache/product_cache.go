package cache

import (
	"context"
	"fmt"
	"mk-pos-billing/internal/infrastructure/cache"
	"mk-pos-billing/internal/infrastructure/config"
)

type ProductCache interface {
	Get(ctx context.Context, productID int) (*Product, error)
}

type productCache struct {
	cacheSvc *cache.Service
	prefixes *config.CachePrefixes
	tags     *config.CacheTags
}

func NewProductCache(cacheSvc *cache.Service, prefixes *config.CachePrefixes, tags *config.CacheTags) ProductCache {
	return &productCache{
		cacheSvc: cacheSvc,
		prefixes: prefixes,
		tags:     tags,
	}
}

func (s *productCache) Get(ctx context.Context, productID int) (*Product, error) {
	// Use BuildTaggedKey to handle Laravel's cache tagging (prepending the version hash)
	cacheKey := s.cacheSvc.BuildTaggedKey(ctx, s.tags.Product, s.prefixes.Product, fmt.Sprint(productID))

	var productData Product
	err := s.cacheSvc.GetPHPSerialized(ctx, cacheKey, &productData)

	if err != nil {
		return nil, err
	}

	return &productData, nil
}
