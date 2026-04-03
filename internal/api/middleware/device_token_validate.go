package middleware

import (
	"fmt"
	"mk-pos-billing/internal/service"
	"mk-pos-billing/pkg/response"
	"net/http"

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

		device_details, err := m.cacheMasterService.GetDeviceCache(c, deviceToken)
		fmt.Println(device_details)

		if err != nil || device_details == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "Device Token is Invalid", "invalid_device_token": true})
			c.Abort()
			return
		}

		isExpired, expiredOk := device_details["is_expired"].(bool)
		isActive, activeOk := device_details["is_active"].(bool)
		// Optional: handle missing or invalid types
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
				"invalid_device_token":      true,
				"is_active_device_token":    isActive,
				"is_expired_device_token":   isExpired,
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
