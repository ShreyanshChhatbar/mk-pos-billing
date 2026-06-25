package routes

import (
	"mk-pos-billing/internal/api/handlers"
	"mk-pos-billing/internal/api/middleware"

	"github.com/gin-gonic/gin"
)

func SalesInvoiceRoutes(
	r *gin.RouterGroup,
	salesInvoiceHandler *handlers.SalesInvoiceHandler,
	duplicateMW *middleware.DuplicateRequestMiddleware,
	deviceTokenValidateMiddleware *middleware.DeviceTokenValidateMiddleware,
	posAuthTokenValidateMiddleware *middleware.POSAuthTokenValidateMiddleware,
	mapTillMiddleware *middleware.MapTillMiddleware,
	sposCheckPermissionsMiddleware *middleware.SPOSCheckPermissionsMiddleware,
	checkTillStatusMiddleware *middleware.CheckTillStatusMiddleware,
) {
	draft := r.Group("/sales/sales-invoice/draft")
	draft.Use(
		deviceTokenValidateMiddleware.Handle(),
		posAuthTokenValidateMiddleware.Handle(),
		mapTillMiddleware.Handle(),
		sposCheckPermissionsMiddleware.Handle(),
		checkTillStatusMiddleware.Handle(),
	)
	draft.POST("/create", duplicateMW.Handle(), salesInvoiceHandler.CreateSalesInvoice)
	draft.PUT("/:id", duplicateMW.Handle(), salesInvoiceHandler.UpdateSalesInvoice)
}
