package middleware

import (
	"fmt"
	"mk-pos-billing/internal/service"
	"net/http"
	"strings"

	"mk-pos-billing/pkg/response"

	"github.com/gin-gonic/gin"
)

type SPOSCheckPermissionsMiddleware struct {
	cacheMaster service.CacheMasterService
}

func NewSPOSCheckPermissionsMiddleware(cacheMaster service.CacheMasterService) *SPOSCheckPermissionsMiddleware {
	return &SPOSCheckPermissionsMiddleware{
		cacheMaster: cacheMaster,
	}
}

func getOperation(method string) string {
	switch strings.ToUpper(method) {
	case http.MethodGet:
		return "READ"
	case http.MethodPost:
		return "CREATE"
	case http.MethodPut, http.MethodPatch:
		return "UPDATE"
	case http.MethodDelete:
		return "DELETE"
	default:
		return "UNKNOWN"
	}
}

func (m *SPOSCheckPermissionsMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Get and Validate Store Header
		store := c.GetHeader("store")
		if store == "" {
			response.Error(c, http.StatusBadRequest, "Bad Request", "Store Not Found")
			c.Abort()
			return
		}

		// 2. Extract Authorization Header
		var token string
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			response.Error(c, http.StatusUnauthorized, "Unauthorized", "Unauthorized")
			c.Abort()
			return
		}

		token = strings.TrimPrefix(authHeader, "Bearer ")

		// 3. Get User Auth Cache via CacheMasterService
		userCacheMaster, err := m.cacheMaster.GetUserAuthCache(c.Request.Context(), token)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, "Internal Server Error", "Something went wrong")
			c.Abort()
			return
		}

		if len(userCacheMaster) == 0 {
			response.Error(c, http.StatusUnauthorized, "Unauthorized", "Unauthorized, you don`t have any permissions")
			c.Abort()
			return
		}

		userId := fmt.Sprint(userCacheMaster["user_id"])
		var permissions map[string]interface{}
		if perms, ok := userCacheMaster["permissions"].(map[string]interface{}); ok {
			permissions = perms
		}

		if len(permissions) == 0 {
			response.MiddlewareError(c, http.StatusUnauthorized, "Unauthorized", "Unauthorized, Request for access to higher authorities!", map[string]interface{}{
				"store_access_required": true,
			})
			c.Abort()
			return
		}

		// 4. Extract segments from URL path
		// E.g. /api/v1/module/submodule/approve -> segments = ["", "api", "v1", "module", "submodule", "approve"]
		path := strings.TrimSpace(c.Request.URL.Path)
		segments := strings.Split(path, "/")
		
		// In Laravel, segment(1) is the first part after the domain name.
		// In Gin, strings.Split("/a/b", "/") returns ["", "a", "b"].
		// We need to map Laravel's `segment(3)` logic reliably. 
		// Assuming base API structure: /api/v1/{module}/{subModule}/{approvalPermission?}
		// segment(1) = api, segment(2) = v1, segment(3) = module...
		
		var module, subModule, approvalPermission string
		// Adding +1 offset because strings.Split on a leading slash has an empty element at index 0.
		if len(segments) > 3 {
			module = strings.ToLower(segments[3])
		}
		if len(segments) > 4 {
			subModule = strings.ToLower(segments[4])
		}
		if len(segments) > 5 {
			approvalPermission = strings.ToLower(segments[5])
		}

		operation := getOperation(c.Request.Method)

		// 5. Check Module Permission
		modulePerms, hasModule := permissions[module]
		if !hasModule {
			response.MiddlewareError(c, http.StatusUnauthorized, "Unauthorized", "Unauthorized, You don't have access to this module!", map[string]interface{}{
				"is_permission_required": true,
			})
			c.Abort()
			return
		}

		// 6. Check SubModule Permission
		modulePermsMap, ok := modulePerms.(map[string]interface{})
		if !ok {
			response.MiddlewareError(c, http.StatusUnauthorized, "Unauthorized", "Unauthorized, You don't have access to this module format!", map[string]interface{}{
				"is_permission_required": true,
			})
			c.Abort()
			return
		}

		subModulePerms, hasSubModule := modulePermsMap[subModule]
		if !hasSubModule {
			response.MiddlewareError(c, http.StatusUnauthorized, "Unauthorized", "Unauthorized, You don't have access to this sub-module!", map[string]interface{}{
				"is_permission_required": true,
			})
			c.Abort()
			return
		}

		// 7. Check Action (Operation)
		subModuleActionsList, ok := subModulePerms.([]interface{})
		if !ok {
			response.MiddlewareError(c, http.StatusUnauthorized, "Unauthorized", "Unauthorized, Invalid sub-module permissions format!", map[string]interface{}{
				"is_permission_required": true,
			})
			c.Abort()
			return
		}

		hasActionAccess := false
		for _, action := range subModuleActionsList {
			if fmt.Sprint(action) == operation {
				hasActionAccess = true
				break
			}
		}

		if !hasActionAccess {
			response.MiddlewareError(c, http.StatusUnauthorized, "Unauthorized", "Unauthorized, You don't have access to this action!", map[string]interface{}{
				"is_permission_required": true,
			})
			c.Abort()
			return
		}

		// 8. Special Approval Permission Check
		if approvalPermission == "approve" {
			hasApproveAccess := false
			for _, action := range subModuleActionsList {
				if fmt.Sprint(action) == "APPROVE" {
					hasApproveAccess = true
					break
				}
			}
			if !hasApproveAccess {
				response.MiddlewareError(c, http.StatusUnauthorized, "Unauthorized", "Unauthorized, You don't have access to this action!", map[string]interface{}{
					"is_permission_required": true,
				})
				c.Abort()
				return
			}
		}

		// 9. Merge data into Context
		c.Set("store_id", store)
		c.Set("login_user_id", userId)

		c.Next()
	}
}
