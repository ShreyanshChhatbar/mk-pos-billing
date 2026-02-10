package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"mk-pos-billing/internal/infrastructure/cache"
	"mk-pos-billing/internal/service"
	"mk-pos-billing/pkg/utils"

	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
)

const requestContextKey = "rbac.request_context"

var headerValidationRules = map[string]struct {
	missing string
	invalid string
}{
	"X-Tenant-ID": {
		missing: "X-Tenant-ID is required",
		invalid: "X-Tenant-ID must be a positive number",
	},
}

const headerValidationMessage = "Validation failed"

type PermissionMiddleware struct {
	cache *cache.Service
}

func NewPermissionMiddleware(cache *cache.Service) *PermissionMiddleware {
	return &PermissionMiddleware{cache: cache}
}

type RequestContext struct {
	UserID     uint64
	TenantID   uint64
	FacilityID *uint64
}

func GetRequestContext(c *gin.Context) (RequestContext, bool) {
	if v, ok := c.Get(requestContextKey); ok {
		if ctx, ok := v.(RequestContext); ok {
			return ctx, true
		}
	}
	return RequestContext{}, false
}

func setRequestContext(c *gin.Context, ctx RequestContext) {
	c.Set(requestContextKey, ctx)
}

func (m *PermissionMiddleware) RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authCtx, ok := GetAuthContext(c)
		if !ok {
			respondUnauthorized(c, "authentication required")
			return
		}

		tenantID, ok := getHeaderUint64(c, "X-Tenant-ID", true)
		if !ok {
			return
		}

		key := m.cache.BuildKey("rbac", strconv.FormatUint(authCtx.UserID, 10))
		var bundle service.CachedRBAC
		if err := m.cache.GetJSON(c.Request.Context(), key, &bundle); err != nil {
			respondForbidden(c, "Failed to load permissions")
			return
		}

		tenant, ok := bundle.Tenants[tenantID]
		if !ok {
			respondForbidden(c, "Tenant access denied")
			return
		}

		if bundle.CurrentFacilityID == nil {
			// If no facility selected, check if any role for this tenant doesn't require facility login
			hasPerm := false
			for _, role := range tenant.Roles {
				if !role.IsFacilityLoginRequired && roleHasPermission(role, permission) {
					hasPerm = true
					break
				}
			}
			if !hasPerm {
				respondForbidden(c, "Facility selection required")
				return
			}
		} else {
			facilityID := *bundle.CurrentFacilityID
			facility, ok := tenant.Facilities[facilityID]
			if !ok {
				respondForbidden(c, "Facility access denied")
				return
			}

			role, ok := tenant.Roles[facility.RoleID]
			if !ok {
				respondForbidden(c, "Role access denied")
				return
			}

			if role.IsFacilityLoginRequired {
				if bundle.CurrentFacilityCheckIn == nil || bundle.CurrentFacilityCheckIn.CheckOutTime != nil {
					respondForbidden(c, "Facility check-in required")
					return
				}
			}

			if !roleHasPermission(role, permission) {
				respondForbidden(c, "Permission denied")
				return
			}
		}

		ctx := RequestContext{
			UserID:     authCtx.UserID,
			TenantID:   tenantID,
			FacilityID: bundle.CurrentFacilityID,
		}
		setRequestContext(c, ctx)
		enrichSentryScope(c, ctx)
		c.Next()
	}
}

func roleHasPermission(role service.RoleRBAC, permission string) bool {
	if role.PermissionIndex == nil {
		return false
	}
	return role.PermissionIndex[permission]
}

func getHeaderUint64(c *gin.Context, header string, required bool) (uint64, bool) {
	value := strings.TrimSpace(c.GetHeader(header))
	if value == "" {
		if required {
			msg := headerValidationRules[header].missing
			if msg == "" {
				msg = header + " is required"
			}
			respondHeaderError(c, headerValidationMessage, msg)
			return 0, false
		}
		return 0, true
	}
	n, err := strconv.ParseUint(value, 10, 64)
	if err != nil || n == 0 {
		msg := headerValidationRules[header].invalid
		if msg == "" {
			msg = "Invalid " + header
		}
		respondHeaderError(c, headerValidationMessage, msg)
		return 0, false
	}
	return n, true
}

func respondHeaderError(c *gin.Context, apiMessage, detail string) {
	utils.RespondError(c, http.StatusBadRequest, apiMessage, detail)
	c.Abort()
}

func respondUnauthorized(c *gin.Context, reason string) {
	utils.RespondError(c, http.StatusUnauthorized, "Unauthorized", reason)
	c.Abort()
}

func respondForbidden(c *gin.Context, message string) {
	utils.RespondError(c, http.StatusForbidden, message, nil)
	c.Abort()
}

func enrichSentryScope(c *gin.Context, ctx RequestContext) {
	if hub := sentrygin.GetHubFromContext(c); hub != nil {
		scope := hub.Scope()
		scope.SetTag("tenant_id", fmt.Sprintf("%d", ctx.TenantID))
		if ctx.FacilityID != nil {
			scope.SetTag("facility_id", fmt.Sprintf("%d", *ctx.FacilityID))
		}
		scope.SetTag("user_id", fmt.Sprintf("%d", ctx.UserID))
	}
}
