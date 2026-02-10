package middleware

import (
	"strings"
	"time"

	"mk-pos-billing/internal/infrastructure/config"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const authContextKey = "rbac.auth_context"

type AuthContext struct {
	UserID    uint64
	TenantID  *uint64
	ExpiresAt time.Time
	Claims    jwt.MapClaims
}

func GetAuthContext(c *gin.Context) (AuthContext, bool) {
	if v, ok := c.Get(authContextKey); ok {
		if ctx, ok := v.(AuthContext); ok {
			return ctx, true
		}
	}
	return AuthContext{}, false
}

func setAuthContext(c *gin.Context, ctx AuthContext) {
	c.Set(authContextKey, ctx)
}

type AuthMiddleware struct {
	authCfg config.AuthConfig
}

func NewAuthMiddleware(authCfg config.AuthConfig) *AuthMiddleware {
	return &AuthMiddleware{authCfg: authCfg}
}

func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
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

		userID, ok := claims["user_id"].(float64)
		if !ok || userID == 0 {
			respondUnauthorized(c, "invalid token claims")
			return
		}

		setAuthContext(c, AuthContext{
			UserID:    uint64(userID),
			ExpiresAt: time.Unix(exp, 0),
			Claims:    claims,
		})
		c.Next()
	}
}

func extractBearer(header string) string {
	if header == "" {
		return ""
	}
	const prefix = "bearer "
	if len(header) >= len(prefix) && strings.EqualFold(header[:len(prefix)], prefix) {
		return strings.TrimSpace(header[len(prefix):])
	}
	return ""
}

func parseAccessToken(tokenStr, secret string) (jwt.MapClaims, int64, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return nil, 0, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, 0, jwt.ErrTokenInvalidClaims
	}
	exp, ok := claims["expires_at"].(float64)
	if !ok {
		return nil, 0, jwt.ErrTokenInvalidClaims
	}
	return claims, int64(exp), nil
}
