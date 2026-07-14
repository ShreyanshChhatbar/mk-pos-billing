package handlers

import (
	"errors"
	"mk-pos-billing/internal/api/middleware"
	"mk-pos-billing/internal/api/request"
	"mk-pos-billing/internal/infrastructure/config"
	"mk-pos-billing/internal/service"
	"mk-pos-billing/pkg/validation"
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
	body := validation.GetValidBodyData[request.CreateOrUpdateSalesInvoiceRequest](c)

	posCtx, ok := h.posContext(c)
	if !ok {
		return
	}

	itemsPresent := body.Items != nil
	items := []request.CreateOrUpdateInvoiceItem{}
	if body.Items != nil {
		items = *body.Items
	}

	input := service.CreateOrUpdateSalesInvoiceInput{
		ID:                id,
		OrganizationID:    h.salesCfg.DefaultOrganizationID,
		StoreID:           posCtx.StoreID,
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
		DeviceMasterID:    posCtx.DeviceMasterID,
		TillID:            posCtx.TillID,
		TillTransactionID: posCtx.TillTransactionID,
		UserID:            posCtx.UserID,
		ItemsPresent:      itemsPresent,
		Items:             make([]service.CreateOrUpdateSalesInvoiceItem, 0, len(items)),
		Payments:          make([]service.CreateOrUpdateSalesInvoicePayment, 0, len(body.Payments)),
	}

	for _, it := range items {
		input.Items = append(input.Items, service.CreateOrUpdateSalesInvoiceItem{
			ProductID:      it.ProductID,
			BatchCode:      it.BatchCode,
			Quantity:       it.Quantity,
			IsFreeProduct:  it.IsFreeProduct,
			ComboProductID: it.ComboProductID,
		})
	}
	for _, p := range body.Payments {
		input.Payments = append(input.Payments, service.CreateOrUpdateSalesInvoicePayment{
			ID:                   p.ID,
			StorePaymentMethodID: p.StorePaymentMethodID,
			Amount:               p.Amount,
			VoucherCode:          p.VoucherCode,
			IsAdvanceRefund:      p.IsAdvanceRefund,
		})
	}

	// zap.L().Info("SPOSCheckPermissionsMiddleware: storeCache", zap.Any("input", input))
	// os.Exit(0)

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

type posContext struct {
	StoreID           uint64
	UserID            uint64
	DeviceMasterID    *uint64
	TillID            *uint64
	TillTransactionID *uint64
}

func (h *SalesInvoiceHandler) posContext(c *gin.Context) (posContext, bool) {
	storeID, ok := middleware.POSStoreID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "type": "Bad Request", "message": "Store Not Found"})
		return posContext{}, false
	}

	userID, ok := middleware.POSUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "Invalid User Token: Unauthorized"})
		return posContext{}, false
	}

	deviceMasterID, _ := middleware.POSDeviceMasterID(c)
	tillID, _ := middleware.POSTillID(c)
	tillTransactionID, _ := middleware.POSTillTransactionID(c)

	return posContext{
		StoreID:           storeID,
		UserID:            userID,
		DeviceMasterID:    deviceMasterID,
		TillID:            tillID,
		TillTransactionID: tillTransactionID,
	}, true
}

func mapCreateSalesInvoiceError(err error) int {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return http.StatusBadRequest
	}
	return http.StatusBadRequest
}
