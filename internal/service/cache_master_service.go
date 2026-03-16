package service

import (
	"context"
	"fmt"
	"mk-pos-billing/internal/infrastructure/cache"
	"mk-pos-billing/internal/infrastructure/config"
)

// Define the interface
type CacheMasterService interface {
	GetTillCache(ctx context.Context, storeID int, dbFallback bool) (map[string]interface{}, error)
	GetUserAuthCache(ctx context.Context, token string) (map[string]interface{}, error)
	GetStoreCache(ctx context.Context, storeID int, dbFallback bool) (map[string]interface{}, error)
	GetUserCache(ctx context.Context, userID int, dbFallback bool) (map[string]interface{}, error)
	GetDeviceCache(ctx context.Context, deviceToken string) (map[string]interface{}, error)
	GetProductCache(ctx context.Context, productID int) (map[string]interface{}, error)
}

// Ensure implementation
var _ CacheMasterService = (*cacheMasterService)(nil)

type cacheMasterService struct {
	cacheSvc *cache.Service
	prefixes config.CachePrefixes
	tags	 config.CacheTags
	expiry	 config.CacheExpiry
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
func (s *cacheMasterService) GetTillCache(ctx context.Context, storeID int, dbFallback bool) (map[string]interface{}, error) {
	// The Laravel code uses cache tags, which go-redis handles differently (often via sets), 
	// but mapping to a simple string key approach is standard in basic Redis structures.
	// Use BuildTaggedKey to handle Laravel's cache tagging (prepending the version hash)
	cacheKey := s.cacheSvc.BuildTaggedKey(ctx, s.tags.Till, s.prefixes.Till, fmt.Sprint(storeID))	
	
	var tillData map[string]interface{}
	err := s.cacheSvc.GetJSON(ctx, cacheKey, &tillData)
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

	return tillData, nil
}

// GetUserAuthCache mimics the Laravel getUserAuthCache function
func (s *cacheMasterService) GetUserAuthCache(ctx context.Context, token string) (map[string]interface{}, error) {
	userCacheKey := s.cacheSvc.BuildKey(s.prefixes.PosAuthUser, token)
	
	var findUserCache struct {
		ID int `json:"id"`
	}
	err := s.cacheSvc.GetJSON(ctx, userCacheKey, &findUserCache)
	if err != nil {
		// If empty, return empty response per Laravel logic
		return map[string]interface{}{}, nil
	}
	
	userID := findUserCache.ID
	
	userDetailsKey := s.cacheSvc.BuildKey(s.prefixes.PosAuthToken, fmt.Sprint(userID))
	var userDetails map[string]interface{}
	// Note: We ignore the error here because the Laravel logic checks for keys inside array 
	// (it does not return empty automatically if user details are missing).
	_ = s.cacheSvc.GetJSON(ctx, userDetailsKey, &userDetails)
	
	// Reconstruct the response map
	result := map[string]interface{}{
		"user_id": userID,
	}

	if authToken, ok := userDetails["auth_token"].(string); ok {
		result["auth_token"] = authToken
	} else {
		result["auth_token"] = ""
	}

	if deviceToken, ok := userDetails["device_token"].(string); ok {
		result["device_token"] = deviceToken
	} else {
		result["device_token"] = ""
	}

	if storeId, ok := userDetails["store_id"].(float64); ok {
		result["store_id"] = int(storeId)
	} else {
		result["store_id"] = ""
	}

	if permissions, ok := userDetails["permissions"]; ok && permissions != nil {
		result["permissions"] = permissions
	} else {
		result["permissions"] = []string{}
	}

	return result, nil
}

// GetStoreCache mimics the Laravel getStoreCache function
func (s *cacheMasterService) GetStoreCache(ctx context.Context, storeID int, dbFallback bool) (map[string]interface{}, error) {
	if storeID == 0 {
		return nil, nil // Equivalant to Laravel return null
	}
	
	// Use BuildTaggedKey to handle Laravel's cache tagging (prepending the version hash)
	cacheKey := s.cacheSvc.BuildTaggedKey(ctx, s.tags.Store, s.prefixes.Store, fmt.Sprint(storeID))

	var storeData map[string]interface{}
	err := s.cacheSvc.GetJSON(ctx, cacheKey, &storeData)
	
	if err != nil {
		if dbFallback {
			// Placeholder: SetStoreCache implementation
			return nil, fmt.Errorf("cache miss: dbFallback not fully implemented")
		}
		// Equivalant to Laravel return false if no dbFallback
		return nil, nil 
	}
	
	// Array column transformation for product categories discounts
	if discounts, ok := storeData["product_categories_discounts"].([]interface{}); ok && len(discounts) > 0 {
		mappedDiscounts := make(map[string]interface{})
		for _, v := range discounts {
			if discountItem, ok := v.(map[string]interface{}); ok {
				if catID, exists := discountItem["category_id"]; exists {
					mappedDiscounts[fmt.Sprint(catID)] = discountItem
				}
			}
		}
		storeData["product_categories_discounts"] = mappedDiscounts
	}

	return storeData, nil
}

// GetUserCache mimics the Laravel getUserCache function
func (s *cacheMasterService) GetUserCache(ctx context.Context, userID int, dbFallback bool) (map[string]interface{}, error) {
	cacheKey := s.cacheSvc.BuildKey(s.prefixes.User, fmt.Sprint(userID))
	
	var userData map[string]interface{}
	err := s.cacheSvc.GetJSON(ctx, cacheKey, &userData)
	
	if err != nil || dbFallback {
		// Either cache miss or forced db fallback
		if dbFallback {
			// Placeholder: SetUserCache
			return nil, fmt.Errorf("cache miss/rebuild: dbFallback not fully implemented")
		}
		return nil, err
	}
	
	return userData, nil
}

// GetDeviceCache mimics the Laravel getDeviceCache function
func (s *cacheMasterService) GetDeviceCache(ctx context.Context, deviceToken string) (map[string]interface{}, error) {
	// First, fetch the array/list of allowed device tokens
	deviceRememberKey := s.cacheSvc.BuildKey(s.prefixes.DeviceRemember, "")
	var deviceTokens []string
	
	err := s.cacheSvc.GetJSON(ctx, deviceRememberKey, &deviceTokens)
	if err != nil {
		return nil, nil // No known devices
	}
	
	isKnownDevice := false
	for _, dt := range deviceTokens {
		if dt == deviceToken {
			isKnownDevice = true
			break
		}
	}
	
	if !isKnownDevice {
		return nil, nil
	}
	
	deviceDetailKey := s.cacheSvc.BuildKey(s.prefixes.Device, deviceToken)
	var details map[string]interface{}
	
	err = s.cacheSvc.GetJSON(ctx, deviceDetailKey, &details)
	if err != nil {
		return nil, nil
	}
	
	return details, nil
}

// GetProductCache mimics the Laravel getProductCache function
func (s *cacheMasterService) GetProductCache(ctx context.Context, productID int) (map[string]interface{}, error) {
	// Use BuildTaggedKey to handle Laravel's cache tagging (prepending the version hash)
	cacheKey := s.cacheSvc.BuildTaggedKey(ctx, s.tags.Product, s.prefixes.Product, fmt.Sprint(productID))
	
	var productData map[string]interface{}
	err := s.cacheSvc.GetJSON(ctx, cacheKey, &productData)
	
	if err != nil {
		// Cache miss
		// Placeholder: SetProductCache
		return nil, fmt.Errorf("cache miss: dbFallback not fully implemented")
	}
	
	return productData, nil
}
