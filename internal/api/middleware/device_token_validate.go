package middleware

import (
	"mk-pos-billing/internal/service"
	"mk-pos-billing/pkg/response"
	"net/http"
	"strings"

	"mk-pos-billing/internal/infrastructure/config"

	"github.com/gin-gonic/gin"
)

type DeviceTokenValidateMiddleware struct {
	cacheMasterService service.CacheMasterService
	cachePrefixesCfg   config.CachePrefixes
}

func NewDeviceTokenValidateMiddleware(cacheMasterService service.CacheMasterService, cachePrefixesCfg config.CachePrefixes) *DeviceTokenValidateMiddleware {
	return &DeviceTokenValidateMiddleware{
		cacheMasterService: cacheMasterService,
		cachePrefixesCfg:   cachePrefixesCfg,
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

		deviceDetails, err := m.cacheMasterService.GetDeviceCache(c, deviceToken)

		if err != nil || deviceDetails == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "Device Token is Invalid", "invalid_device_token": true})
			c.Abort()
			return
		}

		isExpired, expiredOk := toBool(deviceDetails["is_expired"])
		isActive, activeOk := toBool(deviceDetails["is_active"])
		if !expiredOk || !activeOk {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "Invalid device details format",
			})
			c.Abort()
			return
		}

		if isExpired || !isActive {
			response.Error(c, http.StatusUnauthorized, "Device is either expired or inactive", map[string]interface{}{
				"invalid_device_token":    true,
				"is_active_device_token":  isActive,
				"is_expired_device_token": isExpired,
			})
			c.Abort()
			return
		}

		if deviceMasterID, ok := toUint64(deviceDetails["device_id"]); ok {
			SetPOSDeviceMasterID(c, deviceMasterID)
		}

		deviceType := strings.ToUpper(toString(deviceDetails["device_type"]))
		if c.GetHeader("store") == "" && deviceType != "TAB" {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "Store Not Found"})
			c.Abort()
			return
		}
		c.Next()
	}
}
