package routes

import (
	"mk-pos-billing/internal/api/handlers"
	"mk-pos-billing/internal/api/middleware"

	"github.com/gin-gonic/gin"
)

type RouteConfig struct {
	DuplicateRequestMiddleware     *middleware.DuplicateRequestMiddleware
	DeviceTokenValidateMiddleware  *middleware.DeviceTokenValidateMiddleware
	SPOSCheckPermissionsMiddleware *middleware.SPOSCheckPermissionsMiddleware
	SalesInvoiceHandler            *handlers.SalesInvoiceHandler
	CacheTestHandler               *handlers.CacheTestHandler
}

func RegisterAllRoutes(r *gin.Engine, cfg RouteConfig) {
	r.GET("/health", handlers.HealthCheck)

	v1 := r.Group("/api/v1")

	TestRoutes(v1, cfg.CacheTestHandler)
	SalesInvoiceRoutes(v1, cfg.SalesInvoiceHandler, cfg.DuplicateRequestMiddleware, cfg.DeviceTokenValidateMiddleware, cfg.SPOSCheckPermissionsMiddleware)
}
