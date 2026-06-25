package container

import (
	"mk-pos-billing/internal/api/handlers"
	"mk-pos-billing/internal/api/middleware"
	"mk-pos-billing/internal/service"
)

type ServerApp struct {
	SalesInvoiceHandler            *handlers.SalesInvoiceHandler
	DuplicateRequestMiddleware     *middleware.DuplicateRequestMiddleware
	DeviceTokenValidateMiddleware  *middleware.DeviceTokenValidateMiddleware
	POSAuthTokenValidateMiddleware *middleware.POSAuthTokenValidateMiddleware
	MapTillMiddleware              *middleware.MapTillMiddleware
	CacheMasterService             service.CacheMasterService
	SPOSCheckPermissionsMiddleware *middleware.SPOSCheckPermissionsMiddleware
	CheckTillStatusMiddleware      *middleware.CheckTillStatusMiddleware
	CacheTestHandler               *handlers.CacheTestHandler
}
