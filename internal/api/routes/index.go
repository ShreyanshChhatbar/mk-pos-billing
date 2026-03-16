package routes

import (
	"mk-pos-billing/internal/api/handlers"
	"mk-pos-billing/internal/api/middleware"

	"github.com/gin-gonic/gin"
)

type RouteConfig struct {
	SalesInvoiceHandler            *handlers.SalesInvoiceHandler
	DuplicateRequestMiddleware     *middleware.DuplicateRequestMiddleware
	SPOSCheckPermissionsMiddleware *middleware.SPOSCheckPermissionsMiddleware
	CacheTestHandler               *handlers.CacheTestHandler
}

func RegisterAllRoutes(r *gin.Engine, cfg RouteConfig) {
	r.GET("/health", handlers.HealthCheck)

	TestRoutes(r, cfg.CacheTestHandler)
	SalesInvoiceRoutes(r, cfg.SalesInvoiceHandler, cfg.DuplicateRequestMiddleware)
}
