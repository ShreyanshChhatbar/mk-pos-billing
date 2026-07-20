package middleware

import (
	domaincache "mk-pos-billing/internal/domain/cache"
	"mk-pos-billing/pkg/response"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type DeviceTokenValidateMiddleware struct {
	deviceCache domaincache.DeviceCache
}

func NewDeviceTokenValidateMiddleware(deviceCache domaincache.DeviceCache) *DeviceTokenValidateMiddleware {
	return &DeviceTokenValidateMiddleware{
		deviceCache: deviceCache,
	}
}

func (m *DeviceTokenValidateMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Get Token From Header ( key : device )
		deviceToken := c.GetHeader("device")
		if deviceToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "Invalid Device Token: Unauthorized", "invalid_device_token": true})
			c.Abort()
			return
		}

		deviceDetails, err := m.deviceCache.Get(c, deviceToken)

		if err != nil || deviceDetails == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "Device Token is Invalid", "invalid_device_token": true})
			c.Abort()
			return
		}

		if deviceDetails.IsExpired || !deviceDetails.IsActive {
			response.Error(c, http.StatusUnauthorized, "Device is either expired or inactive", map[string]interface{}{
				"invalid_device_token":    true,
				"is_active_device_token":  deviceDetails.IsActive,
				"is_expired_device_token": deviceDetails.IsExpired,
			})
			c.Abort()
			return
		}

		if deviceDetails.DeviceID > 0 {
			SetPOSDeviceMasterID(c, uint64(deviceDetails.DeviceID))
		}

		deviceType := strings.ToUpper(deviceDetails.DeviceType)
		if c.GetHeader("store") == "" && deviceType != "TAB" {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "Store Not Found"})
			c.Abort()
			return
		}
		c.Next()
	}
}
