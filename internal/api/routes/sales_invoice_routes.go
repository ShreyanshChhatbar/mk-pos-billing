package routes

import (
	"mk-pos-billing/internal/api/handlers"
	"mk-pos-billing/internal/api/middleware"

	"github.com/gin-gonic/gin"
)

func SalesInvoiceRoutes(r *gin.Engine, salesInvoiceHandler *handlers.SalesInvoiceHandler, duplicateMW *middleware.DuplicateRequestMiddleware) {
	v1 := r.Group("/api/v1")

	draft := v1.Group("/sales/sales-invoice/draft")
	draft.Use(duplicateMW.Handle())
	draft.POST("/create", salesInvoiceHandler.CreateSalesInvoice)
	draft.PUT("/:id", salesInvoiceHandler.UpdateSalesInvoice)
}
