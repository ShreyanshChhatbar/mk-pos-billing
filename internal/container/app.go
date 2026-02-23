package container

import (
	"mk-pos-billing/internal/api/handlers"
	"mk-pos-billing/internal/api/middleware"
)

type ServerApp struct {
	SalesInvoiceHandler        *handlers.SalesInvoiceHandler
	DuplicateRequestMiddleware *middleware.DuplicateRequestMiddleware
}
