package cache

import (
	"context"
	"fmt"
	"mk-pos-billing/internal/infrastructure/cache"
	"mk-pos-billing/internal/infrastructure/config"
)

type UserAuth struct {
	UserID      int                            `php:"user_id"`
	AuthToken   string                         `php:"auth_token"`
	DeviceToken string                         `php:"device_token"`
	StoreID     string                         `php:"store_id"`
	Permissions map[string]map[string][]string `php:"permissions"`
}

type UserAuthCache interface {
	Get(ctx context.Context, token string) (*UserAuth, error)
}

type userAuthCache struct {
	cacheSvc *cache.Service
	prefixes *config.CachePrefixes
}

func NewUserAuthCache(cacheSvc *cache.Service, prefixes *config.CachePrefixes) UserAuthCache {
	return &userAuthCache{
		cacheSvc: cacheSvc,
		prefixes: prefixes,
	}
}

func (s *userAuthCache) Get(ctx context.Context, token string) (*UserAuth, error) {
	userCacheKey := s.cacheSvc.BuildKey(s.prefixes.PosAuthUser, token)

	var findUserCache struct {
		ID int `php:"id"`
	}
	err := s.cacheSvc.GetPHPSerialized(ctx, userCacheKey, &findUserCache)
	if err != nil {
		var userID int
		if err := s.cacheSvc.GetPHPSerialized(ctx, userCacheKey, &userID); err != nil || userID == 0 {
			return nil, err
		}
		findUserCache.ID = userID
	}

	userID := findUserCache.ID

	userDetailsKey := s.cacheSvc.BuildKey(s.prefixes.PosAuthToken, fmt.Sprint(userID))
	var userDetails UserAuth
	// Note: We ignore the error here because the Laravel logic checks for keys inside array
	// (it does not return empty automatically if user details are missing).
	_ = s.cacheSvc.GetPHPSerialized(ctx, userDetailsKey, &userDetails)

	// Ensure userID is set from the index lookup since it's the anchor
	userDetails.UserID = userID

	return &userDetails, nil
}
