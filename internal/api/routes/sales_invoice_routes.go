package routes

import (
	"mk-pos-billing/internal/api/handlers"
	"mk-pos-billing/internal/api/middleware"

	"github.com/gin-gonic/gin"
)

func SalesInvoiceRoutes(r *gin.RouterGroup, salesInvoiceHandler *handlers.SalesInvoiceHandler, duplicateMW *middleware.DuplicateRequestMiddleware, deviceTokenValidateMiddleware *middleware.DeviceTokenValidateMiddleware, sposCheckPermissionsMiddleware *middleware.SPOSCheckPermissionsMiddleware) {
	draft := r.Group("/sales/sales-invoice/draft")
	draft.Use(deviceTokenValidateMiddleware.Handle(), sposCheckPermissionsMiddleware.Handle())
	draft.POST("/create", duplicateMW.Handle(), salesInvoiceHandler.CreateSalesInvoice)
	draft.PUT("/:id", salesInvoiceHandler.UpdateSalesInvoice)
}
