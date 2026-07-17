package cache

import (
	"context"
	"fmt"
	"mk-pos-billing/internal/infrastructure/cache"
	"mk-pos-billing/internal/infrastructure/config"
)

type Device struct {
	DeviceID      int     `php:"device_id"`
	StoreID       *int    `php:"store_id"`
	IsActive      bool    `php:"is_active"`
	IsExpired     bool    `php:"is_expired"`
	DeviceType    string  `php:"device_type"`
	UserAuthToken *string `php:"user_auth_token"`
}

type DeviceCache interface {
	Get(ctx context.Context, deviceToken string) (*Device, error)
}

type deviceCache struct {
	cacheSvc *cache.Service
	prefixes config.CachePrefixes
}

func NewDeviceCache(cacheSvc *cache.Service, prefixes config.CachePrefixes) DeviceCache {
	return &deviceCache{
		cacheSvc: cacheSvc,
		prefixes: prefixes,
	}
}

func (s *deviceCache) Get(ctx context.Context, deviceToken string) (*Device, error) {
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
	var details Device

	if err := s.cacheSvc.GetPHPSerialized(ctx, deviceDetailKey, &details); err != nil {
		return nil, err
	}

	return &details, nil
}
