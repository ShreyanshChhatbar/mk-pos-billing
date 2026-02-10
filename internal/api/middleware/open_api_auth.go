package middleware

import (
	"fmt"
	"time"

	"mk-pos-billing/internal/domain/repository"
	"mk-pos-billing/internal/infrastructure/cache"
	"mk-pos-billing/internal/infrastructure/config"

	"github.com/gin-gonic/gin"
)

type OpenApiAuthMiddleware struct {
	authCfg config.AuthConfig
	cache   *cache.Service
	repo    *repository.OpenApiUserRepository
}

func NewOpenApiAuthMiddleware(authCfg config.AuthConfig, cache *cache.Service, repo *repository.OpenApiUserRepository) *OpenApiAuthMiddleware {
	return &OpenApiAuthMiddleware{authCfg: authCfg, cache: cache, repo: repo}
}

func (m *OpenApiAuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractBearer(c.GetHeader("Authorization"))
		if token == "" {
			respondUnauthorized(c, "missing token")
			return
		}

		claims, exp, err := parseAccessToken(token, m.authCfg.JWTSecret)
		if err != nil {
			respondUnauthorized(c, "invalid token")
			return
		}

		scope, ok := claims["scope"].(string)
		if !ok || scope != "open_api_user" {
			respondUnauthorized(c, "invalid token scope")
			return
		}

		userID, ok := claims["user_id"].(float64)
		if !ok || userID == 0 {
			respondUnauthorized(c, "invalid token claims")
			return
		}
		uid := uint64(userID)

		// Validate against cache
		cacheKey := m.cache.BuildKey("open_api", "token", "user", fmt.Sprintf("%d", uid), "access")
		storedToken, err := m.cache.Get(c.Request.Context(), cacheKey)
		if err != nil || storedToken != token {
			respondUnauthorized(c, "token invalid or expired")
			return
		}

		// Fetch user to get Tenant ID
		user, err := m.repo.FindByID(c.Request.Context(), uid)
		if err != nil {
			respondUnauthorized(c, "user not found")
			return
		}
		tenantID := user.TenantID

		setAuthContext(c, AuthContext{
			UserID:    uid,
			TenantID:  &tenantID,
			ExpiresAt: time.Unix(exp, 0),
			Claims:    claims,
		})
		c.Next()
	}
}
