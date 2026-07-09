package service

import (
	"context"
	"fmt"
	"mk-pos-billing/internal/domain/model"
	"mk-pos-billing/internal/infrastructure/cache"
	"mk-pos-billing/internal/infrastructure/config"

	"go.uber.org/zap"
)

// Define the interface
type CacheMasterService interface {
	GetTillCache(ctx context.Context, storeID int, dbFallback bool) (*model.TillCache, error)
	GetUserAuthCache(ctx context.Context, token string) (*model.UserAuthCache, error)
	GetStoreCache(ctx context.Context, storeID int, dbFallback bool) (*model.StoreCache, error)
	GetUserCache(ctx context.Context, userID int, dbFallback bool) (*model.UserCache, error)
	GetDeviceCache(ctx context.Context, deviceToken string) (*model.DeviceCache, error)
	GetProductCache(ctx context.Context, productID int) (*model.ProductCache, error)
}

// Ensure implementation
var _ CacheMasterService = (*cacheMasterService)(nil)

type cacheMasterService struct {
	cacheSvc *cache.Service
	prefixes config.CachePrefixes
	tags     config.CacheTags
	expiry   config.CacheExpiry
}

func NewCacheMasterService(cacheSvc *cache.Service, prefixes config.CachePrefixes, tags config.CacheTags, expiry config.CacheExpiry) CacheMasterService {
	return &cacheMasterService{
		cacheSvc: cacheSvc,
		prefixes: prefixes,
		tags:     tags,
		expiry:   expiry,
	}
}

// GetTillCache mimics the Laravel getTillCache function
func (s *cacheMasterService) GetTillCache(ctx context.Context, storeID int, dbFallback bool) (*model.TillCache, error) {
	// The Laravel code uses cache tags, which go-redis handles differently (often via sets),
	// but mapping to a simple string key approach is standard in basic Redis structures.
	// Use BuildTaggedKey to handle Laravel's cache tagging (prepending the version hash)
	cacheKey := s.cacheSvc.BuildTaggedKey(ctx, s.tags.Till, s.prefixes.Till, fmt.Sprint(storeID))

	var tillData model.TillCache
	err := s.cacheSvc.GetPHPSerialized(ctx, cacheKey, &tillData)
	if err != nil {
		// Cache miss
		if dbFallback {
			// Placeholder: SetTillCache will be implemented here
			// err = s.SetTillCache(ctx, storeID)
			// return s.GetTillCache(ctx, storeID, false)
			return nil, fmt.Errorf("cache miss: dbFallback not fully implemented")
		}
		return nil, err
	}

	return &tillData, nil
}

// GetUserAuthCache mimics the Laravel getUserAuthCache function
func (s *cacheMasterService) GetUserAuthCache(ctx context.Context, token string) (*model.UserAuthCache, error) {
	userCacheKey := s.cacheSvc.BuildKey(s.prefixes.PosAuthUser, token)

	var findUserCache struct {
		ID int `json:"id"`
	}
	err := s.cacheSvc.GetPHPSerialized(ctx, userCacheKey, &findUserCache)
	if err != nil {
		var userID int
		if rawErr := s.cacheSvc.GetPHPSerialized(ctx, userCacheKey, &userID); rawErr != nil || userID == 0 {
			return nil, nil
		}
		findUserCache.ID = userID
	}

	userID := findUserCache.ID

	userDetailsKey := s.cacheSvc.BuildKey(s.prefixes.PosAuthToken, fmt.Sprint(userID))
	var userDetails model.UserAuthCache
	// Note: We ignore the error here because the Laravel logic checks for keys inside array
	// (it does not return empty automatically if user details are missing).
	_ = s.cacheSvc.GetPHPSerialized(ctx, userDetailsKey, &userDetails)

	// Ensure userID is set from the index lookup since it's the anchor
	userDetails.UserID = userID

	return &userDetails, nil
}

// GetStoreCache mimics the Laravel getStoreCache function
func (s *cacheMasterService) GetStoreCache(ctx context.Context, storeID int, dbFallback bool) (*model.StoreCache, error) {
	if storeID == 0 {
		return nil, nil // Equivalant to Laravel return null
	}

	// Use BuildTaggedKey to handle Laravel's cache tagging (prepending the version hash)
	cacheKey := s.cacheSvc.BuildTaggedKey(ctx, s.tags.Store, s.prefixes.Store, fmt.Sprint(storeID))

	var storeData model.StoreCache
	err := s.cacheSvc.GetPHPSerialized(ctx, cacheKey, &storeData)
	zap.L().Info("SPOSCheckPermissionsMiddleware: storeCache", zap.Any("cacheKey", cacheKey), zap.Any("err", err))

	if err != nil {
		if dbFallback {
			// Placeholder: SetStoreCache implementation
			return nil, fmt.Errorf("cache miss: dbFallback not fully implemented")
		}
		// Equivalant to Laravel return false if no dbFallback
		return nil, nil
	}

	return &storeData, nil
}

// GetUserCache mimics the Laravel getUserCache function
func (s *cacheMasterService) GetUserCache(ctx context.Context, userID int, dbFallback bool) (*model.UserCache, error) {
	cacheKey := s.cacheSvc.BuildKey(s.prefixes.User, fmt.Sprint(userID))

	var userData model.UserCache
	err := s.cacheSvc.GetPHPSerialized(ctx, cacheKey, &userData)

	if err != nil || dbFallback {
		// Either cache miss or forced db fallback
		if dbFallback {
			// Placeholder: SetUserCache
			return nil, fmt.Errorf("cache miss/rebuild: dbFallback not fully implemented")
		}
		return nil, err
	}

	return &userData, nil
}

// GetDeviceCache mimics the Laravel getDeviceCache function
func (s *cacheMasterService) GetDeviceCache(ctx context.Context, deviceToken string) (*model.DeviceCache, error) {
	// First, fetch the array/list of allowed device tokens
	deviceRememberKey := s.cacheSvc.BuildKey(s.prefixes.DeviceRemember, "")
	isKnownDevice := false
	var deviceTokenMap map[string]string
	if err := s.cacheSvc.GetPHPSerialized(ctx, deviceRememberKey, &deviceTokenMap); err != nil {
		return nil, nil
	}
	for _, dt := range deviceTokenMap {
		if fmt.Sprint(dt) == deviceToken {
			isKnownDevice = true
			break
		}
	}

	if !isKnownDevice {
		return nil, nil
	}

	deviceDetailKey := s.cacheSvc.BuildKey(s.prefixes.Device, deviceToken)
	var details model.DeviceCache

	if err := s.cacheSvc.GetPHPSerialized(ctx, deviceDetailKey, &details); err != nil {
		return nil, nil
	}

	return &details, nil
}

// GetProductCache mimics the Laravel getProductCache function
func (s *cacheMasterService) GetProductCache(ctx context.Context, productID int) (*model.ProductCache, error) {
	// Use BuildTaggedKey to handle Laravel's cache tagging (prepending the version hash)
	cacheKey := s.cacheSvc.BuildTaggedKey(ctx, s.tags.Product, s.prefixes.Product, fmt.Sprint(productID))

	var productData model.ProductCache
	err := s.cacheSvc.GetPHPSerialized(ctx, cacheKey, &productData)

	if err != nil {
		// Cache miss
		// Placeholder: SetProductCache
		return nil, fmt.Errorf("cache miss: dbFallback not fully implemented")
	}

	return &productData, nil
}
