package cache

import (
	"context"
	"fmt"
	"mk-pos-billing/internal/infrastructure/cache"
	"mk-pos-billing/internal/infrastructure/config"
)

type TillTransaction struct {
	ID     int    `php:"id"`
	Status string `php:"status"`
	Date   string `php:"date"`
}

type Till struct {
	ID              int              `php:"id"`
	TillNumber      string           `php:"till_number"`
	IsActive        bool             `php:"is_active"`
	StoreID         int              `php:"store_id"`
	OpenedDatetime  string           `php:"opened_datetime"`
	TillTransaction *TillTransaction `php:"till_transaction"`
}

type TillCache interface {
	Get(ctx context.Context, storeID int) (*Till, error)
}

type tillCache struct {
	cacheSvc *cache.Service
	prefixes *config.CachePrefixes
	tags     *config.CacheTags
}

func NewTillCache(cacheSvc *cache.Service, prefixes *config.CachePrefixes, tags *config.CacheTags) TillCache {
	return &tillCache{
		cacheSvc: cacheSvc,
		prefixes: prefixes,
		tags:     tags,
	}
}

func (s *tillCache) Get(ctx context.Context, storeID int) (*Till, error) {
	// The Laravel code uses cache tags, which go-redis handles differently (often via sets),
	// but mapping to a simple string key approach is standard in basic Redis structures.
	// Use BuildTaggedKey to handle Laravel's cache tagging (prepending the version hash)
	cacheKey := s.cacheSvc.BuildTaggedKey(ctx, s.tags.Till, s.prefixes.Till, fmt.Sprint(storeID))

	var tillData Till
	err := s.cacheSvc.GetPHPSerialized(ctx, cacheKey, &tillData)
	if err != nil {
		return nil, err
	}

	return &tillData, nil
}
