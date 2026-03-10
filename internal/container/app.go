package container

import (
	"mk-pos-billing/internal/api/handlers"
	"mk-pos-billing/internal/api/middleware"
	"mk-pos-billing/internal/service"
)

type ServerApp struct {
	SalesInvoiceHandler        *handlers.SalesInvoiceHandler
	DuplicateRequestMiddleware     *middleware.DuplicateRequestMiddleware
	CacheMasterService             service.CacheMasterService
	// SPOSCheckPermissionsMiddleware *middleware.SPOSCheckPermissionsMiddleware
	CacheTestHandler               *handlers.CacheTestHandler
}
