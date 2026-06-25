package middleware

import (
	"mk-pos-billing/internal/service"
	"mk-pos-billing/pkg/response"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type POSAuthTokenValidateMiddleware struct {
	cacheMaster service.CacheMasterService
}

func NewPOSAuthTokenValidateMiddleware(cacheMaster service.CacheMasterService) *POSAuthTokenValidateMiddleware {
	return &POSAuthTokenValidateMiddleware{cacheMaster: cacheMaster}
}

func (m *POSAuthTokenValidateMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractPOSBearerToken(c.GetHeader("Authorization"))
		if token == "" {
			response.Error(c, http.StatusUnauthorized, "Invalid User Token: Unauthorized")
			c.Abort()
			return
		}

		userCache, err := m.cacheMaster.GetUserAuthCache(c.Request.Context(), token)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, "Something went wrong")
			c.Abort()
			return
		}
		if len(userCache) == 0 {
			response.Error(c, http.StatusUnauthorized, "Unauthorized, you don`t have any permissions")
			c.Abort()
			return
		}

		userID, ok := toUint64(userCache["user_id"])
		if !ok {
			response.Error(c, http.StatusUnauthorized, "Unauthorized, you don`t have any permissions")
			c.Abort()
			return
		}

		permissions, _ := userCache["permissions"].(map[string]interface{})
		if len(permissions) == 0 {
			response.Error(c, http.StatusUnauthorized, "Unauthorized, you don`t have any permissions")
			c.Abort()
			return
		}

		SetPOSUserID(c, userID)
		SetPOSPermissions(c, permissions)
		c.Next()
	}
}

func extractPOSBearerToken(header string) string {
	if header == "" || !strings.HasPrefix(header, "Bearer ") {
		return ""
	}
	return strings.TrimSpace(strings.TrimPrefix(header, "Bearer "))
}
