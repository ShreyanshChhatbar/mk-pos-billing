package middleware

import (
	"fmt"
	// "go/printer"
	"mk-pos-billing/internal/service"
	"net/http"
	"strconv"
	"strings"

	"mk-pos-billing/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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
		return "LIST"
	case http.MethodPost:
		return "ADD"
	case http.MethodPut, http.MethodPatch:
		return "EDIT"
	case http.MethodDelete:
		return "DELETE"
	default:
		return "UNKNOWN"
	}
}

func (m *SPOSCheckPermissionsMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		storeHeader := strings.TrimSpace(c.GetHeader("store"))
		if storeHeader == "" {
			response.Error(c, http.StatusBadRequest, "Store Not Found")
			c.Abort()
			return
		}
		storeID, err := strconv.ParseUint(storeHeader, 10, 64)
		zap.L().Info("SPOSCheckPermissionsMiddleware: storeCache", zap.Any("storeID", storeID), zap.Uint64("storeID", storeID))
		// print("------------store_id-----------------00000000000000000000", storeID)
		if err != nil || storeID == 0 {
			response.Error(c, http.StatusBadRequest, "Invalid Store")
			c.Abort()
			return
		}

		storeCache, err := m.cacheMaster.GetStoreCache(c.Request.Context(), int(storeID), false)
		// print("storeCache", storeCache)
		zap.L().Info("SPOSCheckPermissionsMiddleware: storeCache", zap.Any("storeCache", storeCache), zap.Uint64("storeID", storeID))
		if err != nil || len(storeCache) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "type": "Bad Request", "message": "Invalid Store"})
			c.Abort()
			return
		}

		userID, userOK := POSUserID(c)
		permissions, permsOK := POSPermissions(c)

		if !userOK || !permsOK {
			token := extractPOSBearerToken(c.GetHeader("Authorization"))
			if token == "" {
				response.Error(c, http.StatusUnauthorized, "Invalid User Token: Unauthorized")
				c.Abort()
				return
			}
			userCacheMaster, err := m.cacheMaster.GetUserAuthCache(c.Request.Context(), token)
			if err != nil {
				response.Error(c, http.StatusInternalServerError, "Something went wrong")
				c.Abort()
				return
			}
			if len(userCacheMaster) == 0 {
				response.Error(c, http.StatusUnauthorized, "Unauthorized, you don`t have any permissions")
				c.Abort()
				return
			}
			userID, userOK = toUint64(userCacheMaster["user_id"])
			permissions, permsOK = userCacheMaster["permissions"].(map[string]interface{})
		}

		if !userOK || len(permissions) == 0 {
			response.MiddlewareError(c, http.StatusUnauthorized, "Unauthorized", "Unauthorized, Request for access to higher authorities!", map[string]interface{}{"store_access_required": true})
			c.Abort()
			return
		}

		module, subModule, approvalPermission := laravelSegments(c.Request.URL.Path)
		operation := getOperation(c.Request.Method)

		modulePerms, hasModule := permissions[module]
		if !hasModule {
			response.MiddlewareError(c, http.StatusUnauthorized, "Unauthorized", "Unauthorized, You don't have access to this module!", map[string]interface{}{
				"is_permission_required": true,
			})
			c.Abort()
			return
		}

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

		hasOperation, validActionList := hasAction(subModulePerms, operation)
		if !validActionList {
			response.MiddlewareError(c, http.StatusUnauthorized, "Unauthorized", "Unauthorized, Invalid sub-module permissions format!", map[string]interface{}{
				"is_permission_required": true,
			})
			c.Abort()
			return
		}
		zap.L().Info("SPOSCheckPermissionsMiddleware: permissionsCheck", zap.Any("module", module), zap.Any("subModule", subModule), zap.Any("operation", operation), zap.Any("hasOperation", hasOperation), zap.Any("approvalPermission", approvalPermission))
		if !hasOperation {
			response.MiddlewareError(c, http.StatusUnauthorized, "Unauthorized", "Unauthorized, You don't have access to this action!", map[string]interface{}{
				"is_permission_required": true,
			})
			c.Abort()
			return
		}

		hasApprove, _ := hasAction(subModulePerms, "APPROVE")
		if approvalPermission == "approve" && !hasApprove {
			response.MiddlewareError(c, http.StatusUnauthorized, "Unauthorized", "Unauthorized, You don't have access to this action!", map[string]interface{}{
				"is_permission_required": true,
			})
			c.Abort()
			return
		}

		SetPOSStoreID(c, storeID)
		SetPOSUserID(c, userID)

		c.Next()
	}
}

func laravelSegments(path string) (string, string, string) {
	segments := strings.Split(strings.Trim(path, "/"), "/")
	segment := func(pos int) string {
		index := pos - 1
		if index >= 0 && index < len(segments) {
			return strings.ToLower(segments[index])
		}
		return ""
	}
	return segment(3), segment(4), segment(5)
}

func hasAction(actions interface{}, expected string) (bool, bool) {
	switch list := actions.(type) {
	case []interface{}:
		for _, action := range list {
			if fmt.Sprint(action) == expected {
				return true, true
			}
		}
		return false, true
	case []string:
		for _, action := range list {
			if action == expected {
				return true, true
			}
		}
		return false, true
	}
	return false, false
}
