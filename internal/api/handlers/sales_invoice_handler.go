package handlers

import (
	"errors"
	"mk-pos-billing/internal/api/request"
	"mk-pos-billing/internal/infrastructure/config"
	"mk-pos-billing/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SalesInvoiceHandler struct {
	salesInvoiceService *service.SalesInvoiceService
	salesCfg            config.SalesInvoiceConfig
}

func NewSalesInvoiceHandler(salesInvoiceService *service.SalesInvoiceService, salesCfg config.SalesInvoiceConfig) *SalesInvoiceHandler {
	return &SalesInvoiceHandler{salesInvoiceService: salesInvoiceService, salesCfg: salesCfg}
}

func (h *SalesInvoiceHandler) CreateSalesInvoice(c *gin.Context) {
	h.createOrUpdate(c, nil)
}

func (h *SalesInvoiceHandler) UpdateSalesInvoice(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "type": "Bad Request", "message": "Validation Error", "errors": gin.H{"id": "id must be a positive number"}})
		return
	}
	h.createOrUpdate(c, &id)
}

func (h *SalesInvoiceHandler) createOrUpdate(c *gin.Context, id *uint64) {
	var body request.CreateOrUpdateSalesInvoiceRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "type": "Bad Request", "message": "Validation Error", "errors": gin.H{"body": err.Error()}})
		return
	}

	input := service.CreateOrUpdateSalesInvoiceInput{
		ID:                id,
		OrganizationID:    h.salesCfg.DefaultOrganizationID,
		StoreID:           body.StoreID,
		BillingUserID:     body.BillingUserID,
		CustomerID:        body.CustomerID,
		CustomerAddressID: body.CustomerAddressID,
		DoctorID:          body.DoctorID,
		PatientID:         body.PatientID,
		IsHomeDelivery:    body.IsHomeDelivery,
		IsConfirmed:       body.IsConfirmed,
		PromoCode:         body.PromoCode,
		CourseDays:        body.CourseDays,
		Notes:             body.Notes,
		ASMUserID:         body.ASMUserID,
		DeviceMasterID:    body.DeviceMasterID,
		TillID:            body.TillID,
		TillTransactionID: body.TillTransactionID,
		UserID:            body.UserID,
		Items:             make([]service.CreateOrUpdateSalesInvoiceItem, 0, len(body.Items)),
		Payments:          make([]service.CreateOrUpdateSalesInvoicePayment, 0, len(body.Payments)),
	}

	for _, it := range body.Items {
		input.Items = append(input.Items, service.CreateOrUpdateSalesInvoiceItem{
			ProductID:      it.ProductID,
			BatchCode:      it.BatchCode,
			Quantity:       it.Quantity,
			SalesRate:      it.SalesRate,
			IsFreeProduct:  it.IsFreeProduct,
			ComboProductID: it.ComboProductID,
		})
	}
	for _, p := range body.Payments {
		input.Payments = append(input.Payments, service.CreateOrUpdateSalesInvoicePayment{
			StorePaymentMethodID: p.StorePaymentMethodID,
			Amount:               p.Amount,
			VoucherCode:          p.VoucherCode,
			IsAdvanceRefund:      p.IsAdvanceRefund,
		})
	}

	result, err := h.salesInvoiceService.CreateOrUpdate(c.Request.Context(), input)
	if err != nil {
		status := mapCreateSalesInvoiceError(err)
		if errors.Is(err, service.ErrValidation) {
			c.JSON(status, gin.H{"code": 400, "type": "Bad Request", "message": "Validation Error", "errors": gin.H{"error": err.Error()}})
			return
		}
		c.JSON(status, gin.H{"code": status, "data": []any{}, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": result.Data, "message": result.Message})
}

func mapCreateSalesInvoiceError(err error) int {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return http.StatusBadRequest
	}
	return http.StatusBadRequest
}
