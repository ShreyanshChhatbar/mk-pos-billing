package middleware

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

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

func SetPOSPermissions(c *gin.Context, permissions map[string]interface{}) {
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

func POSPermissions(c *gin.Context) (map[string]interface{}, bool) {
	v, ok := c.Get(posPermissionsKey)
	if !ok {
		return nil, false
	}
	permissions, ok := v.(map[string]interface{})
	return permissions, ok
}

func contextUint64(c *gin.Context, key string) (uint64, bool) {
	v, ok := c.Get(key)
	if !ok {
		return 0, false
	}
	return toUint64(v)
}

func toUint64(value interface{}) (uint64, bool) {
	switch v := value.(type) {
	case uint64:
		return v, v > 0
	case uint:
		return uint64(v), v > 0
	case uint32:
		return uint64(v), v > 0
	case int:
		return uint64(v), v > 0
	case int64:
		return uint64(v), v > 0
	case int32:
		return uint64(v), v > 0
	case float64:
		if v <= 0 {
			return 0, false
		}
		return uint64(v), true
	case float32:
		if v <= 0 {
			return 0, false
		}
		return uint64(v), true
	case json.Number:
		parsed, err := strconv.ParseUint(string(v), 10, 64)
		return parsed, err == nil && parsed > 0
	case string:
		parsed, err := strconv.ParseUint(strings.TrimSpace(v), 10, 64)
		return parsed, err == nil && parsed > 0
	default:
		return 0, false
	}
}

func toBool(value interface{}) (bool, bool) {
	switch v := value.(type) {
	case bool:
		return v, true
	case string:
		parsed, err := strconv.ParseBool(strings.TrimSpace(v))
		return parsed, err == nil
	case int:
		return v != 0, true
	case int64:
		return v != 0, true
	case float64:
		return v != 0, true
	default:
		return false, false
	}
}

func toString(value interface{}) string {
	if value == nil {
		return ""
	}
	return strings.TrimSpace(fmt.Sprint(value))
}

func toTime(value interface{}) (time.Time, bool) {
	switch v := value.(type) {
	case time.Time:
		return v, !v.IsZero()
	case string:
		value := strings.TrimSpace(v)
		if value == "" {
			return time.Time{}, false
		}
		layouts := []string{
			time.RFC3339,
			"2006-01-02T15:04:05",
			"2006-01-02 15:04:05",
			"2006-01-02",
		}
		for _, layout := range layouts {
			if parsed, err := time.ParseInLocation(layout, value, time.Local); err == nil {
				return parsed, true
			}
		}
	}
	return time.Time{}, false
}
