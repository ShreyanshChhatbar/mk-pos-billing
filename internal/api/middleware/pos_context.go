package middleware

import (
	"github.com/gin-gonic/gin"
)

const (
	posUserIDKey            = "pos.user_id"
	posStoreIDKey           = "pos.store_id"
	posDeviceMasterIDKey    = "pos.device_master_id"
	posTillIDKey            = "pos.till_id"
	posTillTransactionIDKey = "pos.till_transaction_id"
	posPermissionsKey       = "pos.permissions"
)

func SetPOSUserID(c *gin.Context, userID uint64) {
	c.Set(posUserIDKey, userID)
	c.Set("user_id", userID)
}

func SetPOSStoreID(c *gin.Context, storeID uint64) {
	c.Set(posStoreIDKey, storeID)
	c.Set("store_id", storeID)
}

func SetPOSDeviceMasterID(c *gin.Context, deviceMasterID uint64) {
	c.Set(posDeviceMasterIDKey, deviceMasterID)
	c.Set("device_master_id", deviceMasterID)
}

func SetPOSTillID(c *gin.Context, tillID uint64) {
	c.Set(posTillIDKey, tillID)
	c.Set("till_id", tillID)
}

func SetPOSTillTransactionID(c *gin.Context, tillTransactionID uint64) {
	c.Set(posTillTransactionIDKey, tillTransactionID)
	c.Set("till_transaction_id", tillTransactionID)
}

func SetPOSPermissions(c *gin.Context, permissions map[string]map[string][]string) {
	c.Set(posPermissionsKey, permissions)
}

func POSUserID(c *gin.Context) (uint64, bool) {
	return contextUint64(c, posUserIDKey)
}

func POSStoreID(c *gin.Context) (uint64, bool) {
	return contextUint64(c, posStoreIDKey)
}

func POSDeviceMasterID(c *gin.Context) (*uint64, bool) {
	v, ok := contextUint64(c, posDeviceMasterIDKey)
	if !ok || v == 0 {
		return nil, ok
	}
	return &v, true
}

func POSTillID(c *gin.Context) (*uint64, bool) {
	v, ok := contextUint64(c, posTillIDKey)
	if !ok || v == 0 {
		return nil, ok
	}
	return &v, true
}

func POSTillTransactionID(c *gin.Context) (*uint64, bool) {
	v, ok := contextUint64(c, posTillTransactionIDKey)
	if !ok || v == 0 {
		return nil, ok
	}
	return &v, true
}

func POSPermissions(c *gin.Context) (map[string]map[string][]string, bool) {
	v, ok := c.Get(posPermissionsKey)
	if !ok {
		return nil, false
	}
	permissions, ok := v.(map[string]map[string][]string)
	return permissions, ok
}

func contextUint64(c *gin.Context, key string) (uint64, bool) {
	v, ok := c.Get(key)
	if !ok {
		return 0, false
	}
	val, ok := v.(uint64)
	return val, ok
}
