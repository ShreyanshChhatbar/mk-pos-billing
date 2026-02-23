package routes

import (
	"mk-pos-billing/internal/api/handlers"
	"mk-pos-billing/internal/api/middleware"

	"github.com/gin-gonic/gin"
)

type RouteConfig struct {
	SalesInvoiceHandler        *handlers.SalesInvoiceHandler
	DuplicateRequestMiddleware *middleware.DuplicateRequestMiddleware
}

func RegisterAllRoutes(r *gin.Engine, cfg RouteConfig) {
	r.GET("/health", handlers.HealthCheck)
	SalesInvoiceRoutes(r, cfg.SalesInvoiceHandler, cfg.DuplicateRequestMiddleware)
}
