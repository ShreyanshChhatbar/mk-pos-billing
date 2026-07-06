package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"mk-pos-billing/internal/domain/model"
	"mk-pos-billing/internal/domain/repository"
	"mk-pos-billing/internal/infrastructure/cache"
	"mk-pos-billing/internal/infrastructure/config"
	"mk-pos-billing/pkg/constants"
	"sort"
	"strconv"
	"time"

	"gorm.io/gorm"
)

var ErrValidation = errors.New("validation error")

type CreateOrUpdateSalesInvoiceInput struct {
	ID                *uint64
	OrganizationID    uint64
	StoreID           uint64
	BillingUserID     uint64
	CustomerID        *uint64
	CustomerAddressID *uint64
	DoctorID          *uint64
	PatientID         *uint64
	IsHomeDelivery    *bool
	IsConfirmed       *bool
	PromoCode         *string
	CourseDays        *float64
	Notes             *string
	ASMUserID         *uint64
	DeviceMasterID    *uint64
	TillID            *uint64
	TillTransactionID *uint64
	UserID            uint64
	ItemsPresent      bool
	Items             []CreateOrUpdateSalesInvoiceItem
	Payments          []CreateOrUpdateSalesInvoicePayment
}

type CreateOrUpdateSalesInvoiceItem struct {
	ProductID      uint64
	BatchCode      string
	Quantity       int
	SalesRate      *float64
	IsFreeProduct  bool
	ComboProductID *uint64
}

type CreateOrUpdateSalesInvoicePayment struct {
	ID                   *uint64
	StorePaymentMethodID uint64
	Amount               float64
	VoucherCode          *string
	IsAdvanceRefund      bool
}

type CreateOrUpdateSalesInvoiceOutput struct {
	Data    any
	Message string
}

type draftProductJSON struct {
	ProductID          uint64  `json:"product_id"`
	BatchCode          string  `json:"batch_code"`
	ExpiryDate         string  `json:"expiry_date"`
	MRP                float64 `json:"mrp"`
	SalesRate          float64 `json:"sales_rate"`
	BaseRate           float64 `json:"base_rate"`
	Quantity           int     `json:"quantity"`
	GSTPercentage      float64 `json:"gst_percentage"`
	GSTAmount          float64 `json:"gst_amount"`
	BillAmount         float64 `json:"bill_amount"`
	DiscountType       string  `json:"discount_type"`
	DiscountAmount     float64 `json:"discount_amount"`
	DiscountPercentage float64 `json:"discount_percentage"`
	IsFreeProduct      bool    `json:"is_free_product"`
	ComboProductID     *uint64 `json:"combo_product_id"`
}

type salesInvoiceDraftJSONPayload struct {
	IsHomeDelivery  bool               `json:"is_home_delivery"`
	IsActive        bool               `json:"is_active"`
	PaymentStatus   string             `json:"payment_status"`
	DeviceMasterID  *uint64            `json:"device_master_id"`
	PromoCode       *string            `json:"promo_code"`
	Notes           *string            `json:"notes"`
	TotalProducts   int                `json:"total_products"`
	TotalItems      int                `json:"total_items"`
	TotalQuantity   int                `json:"total_quantity"`
	TotalGST        float64            `json:"total_gst"`
	SGST            float64            `json:"sgst"`
	CGST            float64            `json:"cgst"`
	IGST            float64            `json:"igst"`
	DeliveryCharges float64            `json:"delivery_charges"`
	TaxableAmount   float64            `json:"taxable_amount"`
	RoundOff        float64            `json:"round_off"`
	Products        []draftProductJSON `json:"products"`
}

type draftComputedLine struct {
	Item         CreateOrUpdateSalesInvoiceItem
	Batch        model.Batch
	SalesRate    float64
	BillAmount   float64
	TotalAmount  float64
	DiscountType string
	DiscountAmt  float64
	DiscountPct  float64
}

type draftCacheData struct {
	SalesInvoiceDraft     model.SalesInvoiceDraftJSON     `json:"salesInvoiceDraft"`
	CalculatedProductData draftCalculatedProductContainer `json:"calculatedProductData"`
}

type draftCalculatedProductContainer struct {
	Products map[string]draftProductJSON `json:"products"`
}

type draftResponse struct {
	ID                   uint64             `json:"id"`
	OrganizationID       uint64             `json:"organization_id"`
	IsHomeDelivery       bool               `json:"is_home_delivery"`
	TotalProducts        int                `json:"total_products"`
	TotalItems           int                `json:"total_items"`
	TotalQuantity        int                `json:"total_quantity"`
	PrepaidAmount        float64            `json:"prepaid_amount"`
	TotalInvoiceAmount   float64            `json:"total_invoice_amount"`
	TotalAmountReceived  float64            `json:"total_amount_received"`
	AmountDue            float64            `json:"amount_due"`
	DeliveryCharges      float64            `json:"delivery_charges"`
	TotalMRP             float64            `json:"total_mrp"`
	TotalSavings         float64            `json:"total_savings"`
	TaxableAmount        float64            `json:"taxable_amount"`
	Status               string             `json:"status"`
	PaymentStatus        string             `json:"payment_status"`
	RoundOff             float64            `json:"round_off"`
	TotalBillBeforePromo float64            `json:"total_bill_amount_before_promo"`
	PromoCode            *string            `json:"promo_code"`
	Items                []draftProductJSON `json:"items"`
	Payments             []draftPaymentResp `json:"payments"`
}

type draftPaymentResp struct {
	ID                   uint64  `json:"id"`
	Amount               float64 `json:"amount"`
	VoucherCode          *string `json:"voucher_code"`
	StorePaymentMethodID uint64  `json:"store_payment_method_id"`
	Type                 string  `json:"type"`
}

type SalesInvoiceService struct {
	repo  *repository.SalesInvoiceRepository
	cache *cache.Service
	cfg   config.SalesInvoiceConfig
}

func NewSalesInvoiceService(repo *repository.SalesInvoiceRepository, cacheService *cache.Service, cfg config.SalesInvoiceConfig) *SalesInvoiceService {
	return &SalesInvoiceService{repo: repo, cache: cacheService, cfg: cfg}
}

func (s *SalesInvoiceService) CreateOrUpdate(ctx context.Context, input CreateOrUpdateSalesInvoiceInput) (*CreateOrUpdateSalesInvoiceOutput, error) {
	input.OrganizationID = s.cfg.DefaultOrganizationID

	if err := s.validateInput(ctx, input); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrValidation, err.Error())
	}

	var out *CreateOrUpdateSalesInvoiceOutput
	err := s.repo.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := s.repo.Tx(tx)
		cacheKey, cacheTag, combinedCacheKey := s.buildDraftCacheKeys(input.StoreID, input.ID)

		cachedData, _ := s.getDraftCache(ctx, combinedCacheKey)

		draft, err := s.loadOrCreateDraft(ctx, txRepo, input, cachedData)
		if err != nil {
			return err
		}

		if input.ID != nil && !input.ItemsPresent {
			existingPayments, err := txRepo.GetDraftPayments(ctx, draft.ID)
			if err != nil {
				return err
			}
			_, draftPayments := buildDraftPayments(input, draft.ID)
			if err := txRepo.AppendDraftPayments(ctx, draftPayments); err != nil {
				return err
			}
			if len(draftPayments) > 0 {
				existingPayments, err = txRepo.GetDraftPayments(ctx, draft.ID)
				if err != nil {
					return err
				}
			}
			var draftPayload salesInvoiceDraftJSONPayload
			if len(draft.DraftJSON) > 0 {
				_ = json.Unmarshal(draft.DraftJSON, &draftPayload)
			}
			draftPayload.Products = []draftProductJSON{}
			draftPayload.TotalProducts = 0
			draftPayload.TotalItems = 0
			draftPayload.TotalQuantity = 0
			draftPayload.TotalGST = 0
			draftPayload.SGST = 0
			draftPayload.CGST = 0
			draftPayload.IGST = 0
			draftPayload.DeliveryCharges = 0
			draftPayload.TaxableAmount = 0
			draftPayload.RoundOff = 0

			raw, _ := json.Marshal(draftPayload)
			draft.DraftJSON = raw
			draft.TotalProducts = 0
			draft.TotalItems = 0
			draft.TotalQuantity = 0
			draft.TotalAmount = 0
			draft.TotalBillAmount = 0
			draft.TotalInvoiceAmount = 0
			draft.TotalAmountReceived = totalReceivedFromDraftPayments(existingPayments)
			draft.TotalDiscount = 0
			draft.RoundOff = 0
			draft.Status = resolveDraftStatus(len(existingPayments) > 0 || len(draftPayments) > 0)
			draft.PaymentStatus = resolvePaymentStatus(draft.TotalAmountReceived, 0)
			draft.UpdatedBy = &input.UserID
			if err := txRepo.SaveDraft(ctx, draft); err != nil {
				return err
			}

			cacheData := draftCacheData{SalesInvoiceDraft: *draft, CalculatedProductData: draftCalculatedProductContainer{Products: map[string]draftProductJSON{}}}
			_ = s.cache.SetJSON(ctx, combinedCacheKey, cacheData, time.Duration(s.cfg.DraftCacheTTLMinutes)*time.Minute)

			out = &CreateOrUpdateSalesInvoiceOutput{Data: s.buildDraftResponse(*draft, nil, existingPayments), Message: "Sales Invoice (Draft) Updated Successfully"}
			return nil
		}

		lines, totalProducts, totalItems, totalQty, totalBill, totalAmount, err := s.calculateDraftLines(ctx, txRepo, input, cachedData)
		if err != nil {
			return err
		}

		totalInvoice := math.Round(totalBill)
		roundOff := totalInvoice - totalBill
		totalDiscount := totalAmount - totalBill

		existingPayments, err := txRepo.GetDraftPayments(ctx, draft.ID)
		if err != nil {
			return err
		}
		newPaymentTotal, draftPayments := buildDraftPayments(input, draft.ID)
		totalReceived := round2(totalReceivedFromDraftPayments(existingPayments) + newPaymentTotal)
		paymentStatus := resolvePaymentStatus(totalReceived, totalInvoice)
		status := resolveDraftStatus(len(existingPayments) > 0 || len(draftPayments) > 0)

		isHomeDelivery := false
		if input.IsHomeDelivery != nil {
			isHomeDelivery = *input.IsHomeDelivery
		}

		draftPayload := salesInvoiceDraftJSONPayload{
			IsHomeDelivery:  isHomeDelivery,
			IsActive:        true,
			PaymentStatus:   paymentStatus,
			DeviceMasterID:  input.DeviceMasterID,
			PromoCode:       input.PromoCode,
			Notes:           input.Notes,
			TotalProducts:   totalProducts,
			TotalItems:      totalItems,
			TotalQuantity:   totalQty,
			TotalGST:        0,
			SGST:            0,
			CGST:            0,
			IGST:            0,
			DeliveryCharges: 0,
			TaxableAmount:   totalBill,
			RoundOff:        roundOff,
			Products:        mapLinesToDraftProducts(lines),
		}
		rawDraftJSON, err := json.Marshal(draftPayload)
		if err != nil {
			return err
		}

		now := time.Now()
		draft.StoreID = input.StoreID
		draft.OrganizationID = input.OrganizationID
		draft.BillingUserID = input.BillingUserID
		draft.CustomerID = input.CustomerID
		draft.CustomerAddressID = input.CustomerAddressID
		draft.DoctorID = input.DoctorID
		draft.PatientID = input.PatientID
		draft.IsHomeDelivery = isHomeDelivery
		draft.Status = status
		draft.PaymentStatus = paymentStatus
		draft.TotalBillAmount = totalBill
		draft.TaxableAmount = totalBill
		draft.TotalAmountBeforeDisc = totalAmount
		draft.DiscountType = "INR"
		draft.DiscountPercentage = 0
		draft.DiscountAmount = totalDiscount
		draft.IGST = 0
		draft.CGST = 0
		draft.SGST = 0
		draft.PrepaidAmount = 0
		draft.TotalGST = 0
		draft.RoundOff = roundOff
		draft.TotalInvoiceAmount = totalInvoice
		draft.TotalProducts = totalProducts
		draft.TotalItems = totalItems
		draft.TotalQuantity = totalQty
		draft.TotalAmount = totalAmount
		draft.TotalDiscount = totalDiscount
		draft.TotalAmountReceived = totalReceived
		draft.TotalBillBeforePromo = totalBill
		draft.PromoCode = input.PromoCode
		draft.Notes = input.Notes
		draft.DraftJSON = rawDraftJSON
		draft.DeviceMasterID = input.DeviceMasterID
		draft.TillID = input.TillID
		draft.TillTransactionID = input.TillTransactionID
		draft.IsActive = true
		draft.UpdatedBy = &input.UserID
		draft.UpdatedAt = now

		if draft.ID == 0 {
			draft.CreatedBy = input.UserID
			draft.CreatedAt = now
			if err := txRepo.CreateDraft(ctx, draft); err != nil {
				return err
			}
			cacheKey, _, combinedCacheKey = s.buildDraftCacheKeys(input.StoreID, &draft.ID)
			_ = cacheKey
			_ = cacheTag
		} else {
			if err := txRepo.SaveDraft(ctx, draft); err != nil {
				return err
			}
		}

		for i := range draftPayments {
			draftPayments[i].SalesInvoiceDraftID = draft.ID
		}
		if err := txRepo.AppendDraftPayments(ctx, draftPayments); err != nil {
			return err
		}

		persistedPayments, _ := txRepo.GetDraftPayments(ctx, draft.ID)
		cacheData := draftCacheData{
			SalesInvoiceDraft:     *draft,
			CalculatedProductData: draftCalculatedProductContainer{Products: s.linesToProductsMap(lines)},
		}
		_ = s.cache.SetJSON(ctx, combinedCacheKey, cacheData, time.Duration(s.cfg.DraftCacheTTLMinutes)*time.Minute)

		msg := "Sales Invoice (Draft) Created Successfully"
		if input.ID != nil && *input.ID > 0 {
			msg = "Sales Invoice (Draft) Updated Successfully"
		}

		if !derefBool(input.IsConfirmed) {
			out = &CreateOrUpdateSalesInvoiceOutput{Data: s.buildDraftResponse(*draft, mapLinesToDraftProducts(lines), persistedPayments), Message: msg}
			return nil
		}

		if !hasAdvanceRefundPayment(persistedPayments) && round2(draft.TotalAmountReceived+draft.PrepaidAmount) != round2(draft.TotalInvoiceAmount) {
			return fmt.Errorf("entered amount is more then the bill amount, please enter proper amount")
		}

		invoiceID, err := s.finalizeInvoice(ctx, txRepo, input, draft, lines)
		if err != nil {
			return err
		}

		// cache final state as draft slot with updated status
		draft.Status = constants.SalesInvoiceStatusInvoiced
		cacheData.SalesInvoiceDraft = *draft
		_ = s.cache.SetJSON(ctx, combinedCacheKey, cacheData, time.Duration(s.cfg.DraftCacheTTLMinutes)*time.Minute)

		out = &CreateOrUpdateSalesInvoiceOutput{Data: invoiceID, Message: "Sales Invoice Created Successfully"}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (s *SalesInvoiceService) validateInput(ctx context.Context, input CreateOrUpdateSalesInvoiceInput) error {
	if !s.repo.ExistsActiveUser(ctx, input.BillingUserID) {
		return errors.New("billing_user_id is invalid")
	}

	hasPayments := len(input.Payments) > 0
	if derefBool(input.IsConfirmed) || hasPayments {
		if input.CustomerID == nil || !s.repo.ExistsActiveCustomer(ctx, *input.CustomerID) {
			return errors.New("customer_id is invalid")
		}
	}

	if derefBool(input.IsConfirmed) {
		if input.PatientID == nil || input.CustomerID == nil || !s.repo.ExistsPatientForCustomer(ctx, *input.PatientID, *input.CustomerID) {
			return errors.New("patient_id is invalid for customer")
		}
	}

	if hasPayments {
		if input.DoctorID == nil || !s.repo.ExistsDoctor(ctx, *input.DoctorID) {
			return errors.New("doctor_id is invalid")
		}
	}

	if hasPayments && derefBool(input.IsHomeDelivery) {
		if input.CustomerAddressID == nil || input.CustomerID == nil || !s.repo.ExistsCustomerAddress(ctx, *input.CustomerAddressID, *input.CustomerID, true) {
			return errors.New("customer_address_id is invalid")
		}
	}

	if derefBool(input.IsConfirmed) {
		if len(input.Items) == 0 {
			return errors.New("items are required when is_confirmed is true")
		}
		if len(input.Payments) == 0 {
			return errors.New("payments are required when is_confirmed is true")
		}
	}

	productIDs := make([]uint64, 0, len(input.Items))
	for i, item := range input.Items {
		if s.cfg.BatchCodeLength > 0 && len(item.BatchCode) > s.cfg.BatchCodeLength {
			return errors.New("items." + strconv.Itoa(i) + ".batch_code length is invalid")
		}

		if !s.repo.ExistsActiveProduct(ctx, item.ProductID) {
			return errors.New("items." + strconv.Itoa(i) + ".product_id is invalid")
		}
		if !s.repo.ExistsBatchCode(ctx, item.BatchCode) {
			return errors.New("items." + strconv.Itoa(i) + ".batch_code is invalid")
		}

		pData, err := s.repo.GetProductValidationData(ctx, item.ProductID)
		if err != nil {
			return err
		}
		if pData == nil {
			return errors.New("items." + strconv.Itoa(i) + ".product_id is invalid")
		}
		if pData.SalesUnit > 0 && item.Quantity%pData.SalesUnit != 0 {
			return errors.New("the quantity must be a multiple of the product's sales unit")
		}
		if pData.IsMSPProduct && item.Quantity != 1 {
			return errors.New("the quantity for a loyalty program product must be one")
		}
		if pData.WSCode == s.cfg.DeliveryChargeProductWSCode && item.Quantity != 1 {
			return errors.New("the quantity for the product for delivery charge must be one")
		}

		productIDs = append(productIDs, item.ProductID)
	}

	if derefBool(input.IsConfirmed) && s.repo.HasNarcoticsProduct(ctx, productIDs, s.cfg.NarcoticsProductScheduledType) {
		if input.CourseDays == nil {
			return errors.New("course_days is required for narcotics products")
		}
	}

	for i, p := range input.Payments {
		if !s.repo.ExistsStorePaymentMethod(ctx, p.StorePaymentMethodID, input.StoreID, input.OrganizationID) {
			return errors.New("payments." + strconv.Itoa(i) + ".store_payment_method_id is invalid")
		}
	}

	if input.ASMUserID != nil && !s.repo.ExistsActiveUser(ctx, *input.ASMUserID) {
		return errors.New("invalid passkey")
	}

	return nil
}
func (s *SalesInvoiceService) loadOrCreateDraft(ctx context.Context, repo *repository.SalesInvoiceRepository, input CreateOrUpdateSalesInvoiceInput, cachedData *draftCacheData) (*model.SalesInvoiceDraftJSON, error) {
	if cachedData != nil && cachedData.SalesInvoiceDraft.ID != 0 {
		draft := cachedData.SalesInvoiceDraft
		return &draft, nil
	}

	if input.ID != nil && *input.ID > 0 {
		return repo.GetDraftByID(ctx, *input.ID, input.StoreID, input.OrganizationID)
	}
	return &model.SalesInvoiceDraftJSON{}, nil
}

func (s *SalesInvoiceService) calculateDraftLines(ctx context.Context, repo *repository.SalesInvoiceRepository, input CreateOrUpdateSalesInvoiceInput, cachedData *draftCacheData) ([]draftComputedLine, int, int, int, float64, float64, error) {
	productIDs, batchCodes := extractUniqueFromItems(input.Items)
	batchMeta, err := repo.FetchBatchMeta(ctx, productIDs, batchCodes)
	if err != nil {
		return nil, 0, 0, 0, 0, 0, err
	}
	stockRows, err := repo.FetchBatchStocks(ctx, input.StoreID, productIDs, batchCodes, false, s.batchExpiryCutoff())
	if err != nil {
		return nil, 0, 0, 0, 0, 0, err
	}

	lines := make([]draftComputedLine, 0, len(input.Items))
	totalBill := 0.0
	totalAmount := 0.0
	totalQty := 0

	for _, item := range input.Items {
		key := fmt.Sprintf("%d_%s", item.ProductID, item.BatchCode)
		batch, ok := batchMeta[key]
		if !ok {
			return nil, 0, 0, 0, 0, 0, fmt.Errorf("batch not found for product %d batch %s", item.ProductID, item.BatchCode)
		}

		available := 0
		for _, r := range stockRows[key] {
			available += r.ClosingStock
		}
		if available < item.Quantity {
			return nil, 0, 0, 0, 0, 0, fmt.Errorf("requested quantity (%d) for batch (%s) is more than available quantity (%d)", item.Quantity, item.BatchCode, available)
		}

		salesRate := batch.MRP
		if item.SalesRate != nil {
			salesRate = *item.SalesRate
		} else if cachedData != nil && cachedData.CalculatedProductData.Products != nil {
			if p, ok := cachedData.CalculatedProductData.Products[key]; ok {
				salesRate = p.SalesRate
			}
		}

		if salesRate-batch.MRP > 0.1 {
			return nil, 0, 0, 0, 0, 0, fmt.Errorf("the sales rate of batch (%s) is greater than MRP", item.BatchCode)
		}

		billAmount := round2(float64(item.Quantity) * salesRate)
		totalLineAmount := round2(float64(item.Quantity) * batch.MRP)
		discountAmt := round2(batch.MRP - salesRate)
		discountPct := 0.0
		if batch.MRP > 0 {
			discountPct = round2(((batch.MRP - salesRate) / batch.MRP) * 100)
		}

		lines = append(lines, draftComputedLine{
			Item:         item,
			Batch:        batch,
			SalesRate:    salesRate,
			BillAmount:   billAmount,
			TotalAmount:  totalLineAmount,
			DiscountType: "INR",
			DiscountAmt:  discountAmt,
			DiscountPct:  discountPct,
		})
		totalBill += billAmount
		totalAmount += totalLineAmount
		totalQty += item.Quantity
	}

	return lines, len(productIDs), len(batchCodes), totalQty, round2(totalBill), round2(totalAmount), nil
}

func (s *SalesInvoiceService) finalizeInvoice(
	ctx context.Context,
	repo *repository.SalesInvoiceRepository,
	input CreateOrUpdateSalesInvoiceInput,
	draft *model.SalesInvoiceDraftJSON,
	lines []draftComputedLine,
) (uint64, error) {
	invoice := model.SalesInvoice{
		OrganizationID:        draft.OrganizationID,
		SalesInvoiceDraftID:   draft.ID,
		StoreID:               draft.StoreID,
		BillingUserID:         draft.BillingUserID,
		CustomerID:            draft.CustomerID,
		CustomerAddressID:     draft.CustomerAddressID,
		DoctorID:              draft.DoctorID,
		PatientID:             draft.PatientID,
		OrderType:             constants.SalesPaymentTypeSales,
		IsHomeDelivery:        draft.IsHomeDelivery,
		TotalBillAmount:       draft.TotalBillAmount,
		TaxableAmount:         draft.TaxableAmount,
		TotalAmountBeforeDisc: draft.TotalAmountBeforeDisc,
		DiscountType:          draft.DiscountType,
		DiscountPercentage:    draft.DiscountPercentage,
		DiscountAmount:        draft.DiscountAmount,
		IGST:                  draft.IGST,
		CGST:                  draft.CGST,
		SGST:                  draft.SGST,
		PrepaidAmount:         draft.PrepaidAmount,
		TotalGST:              draft.TotalGST,
		RoundOff:              draft.RoundOff,
		TotalInvoiceAmount:    draft.TotalInvoiceAmount,
		TotalProducts:         draft.TotalProducts,
		TotalItems:            draft.TotalItems,
		TotalQuantity:         draft.TotalQuantity,
		TotalAmount:           draft.TotalAmount,
		TotalDiscount:         draft.TotalDiscount,
		TotalAmountReceived:   draft.TotalAmountReceived,
		TotalBillBeforePromo:  draft.TotalBillBeforePromo,
		PromoCode:             draft.PromoCode,
		Notes:                 draft.Notes,
		IsActive:              true,
		CreatedBy:             input.UserID,
		DeviceMasterID:        draft.DeviceMasterID,
		TillID:                draft.TillID,
		TillTransactionID:     draft.TillTransactionID,
	}
	if err := repo.CreateInvoice(ctx, &invoice); err != nil {
		return 0, err
	}

	productIDs, batchCodes := extractProductIDsAndBatchCodes(lines)
	batchStocks, err := repo.FetchBatchStocks(ctx, input.StoreID, productIDs, batchCodes, false, s.batchExpiryCutoff())
	if err != nil {
		return 0, err
	}

	details := make([]model.SalesInvoiceDetail, 0)
	for _, line := range lines {
		key := fmt.Sprintf("%d_%s", line.Item.ProductID, line.Item.BatchCode)
		allocations, err := fulfillBatches(batchStocks[key], line.Item.Quantity, line.Item.ProductID, line.Item.BatchCode)
		if err != nil {
			return 0, err
		}

		for _, alloc := range allocations {
			detail := model.SalesInvoiceDetail{
				SalesInvoiceID:     invoice.ID,
				StoreID:            invoice.StoreID,
				ProductID:          line.Item.ProductID,
				StoreBatchID:       &alloc.StoreBatchID,
				PurchaseRate:       alloc.PurchaseRate,
				BatchCode:          alloc.BatchCode,
				ExpiryDate:         alloc.ExpiryDate,
				MRP:                line.Batch.MRP,
				SalesRate:          line.SalesRate,
				BaseRate:           line.SalesRate,
				BillAmount:         round2(float64(alloc.QuantityTaken) * line.SalesRate),
				Quantity:           alloc.QuantityTaken,
				DiscountType:       line.DiscountType,
				DiscountPercentage: line.DiscountPct,
				DiscountAmount:     line.DiscountAmt,
				GSTPercentage:      0,
				GSTAmount:          0,
				TotalAmount:        round2(float64(alloc.QuantityTaken) * line.Batch.MRP),
				CreatedBy:          input.UserID,
				DeviceMasterID:     input.DeviceMasterID,
			}
			details = append(details, detail)

			if err := repo.ReduceInventoryStock(ctx, input.StoreID, line.Item.ProductID, alloc.StoreBatchID, line.Item.BatchCode, alloc.QuantityTaken); err != nil {
				return 0, err
			}
		}
	}

	if err := repo.CreateInvoiceDetails(ctx, details); err != nil {
		return 0, err
	}

	draftPayments, err := repo.GetDraftPayments(ctx, draft.ID)
	if err != nil {
		return 0, err
	}

	invoicePayments := make([]model.SalesInvoicePayment, 0, len(draftPayments))
	for _, p := range draftPayments {
		invoicePayments = append(invoicePayments, model.SalesInvoicePayment{
			SalesInvoiceID:       invoice.ID,
			StoreID:              invoice.StoreID,
			StorePaymentMethodID: p.StorePaymentMethodID,
			Amount:               p.Amount,
			VoucherCode:          p.VoucherCode,
			VoucherAmount:        p.VoucherAmount,
			Type:                 p.Type,
			IsActive:             true,
			CreatedBy:            p.CreatedBy,
			DeviceMasterID:       p.DeviceMasterID,
			TillID:               input.TillID,
			TillTransactionID:    input.TillTransactionID,
		})
	}
	if err := repo.CreateInvoicePayments(ctx, invoicePayments); err != nil {
		return 0, err
	}

	if err := repo.MarkDraftInvoiced(ctx, draft.ID, input.UserID); err != nil {
		return 0, err
	}

	return invoice.ID, nil
}

func (s *SalesInvoiceService) buildDraftCacheKeys(storeID uint64, draftID *uint64) (string, string, string) {
	id := uint64(0)
	if draftID != nil {
		id = *draftID
	}
	key := s.cfg.PrefixDraftBillCache + strconv.FormatUint(id, 10)
	tag := s.cfg.PrefixDraftBillCacheTags + strconv.FormatUint(storeID, 10)
	return key, tag, tag + ":" + key
}

func (s *SalesInvoiceService) batchExpiryCutoff() time.Time {
	return time.Now().AddDate(0, 0, s.cfg.DefaultMinimumDaysForBatch).Truncate(24 * time.Hour)
}

func (s *SalesInvoiceService) getDraftCache(ctx context.Context, combinedKey string) (*draftCacheData, error) {
	var data draftCacheData
	err := s.cache.GetJSON(ctx, combinedKey, &data)
	if err != nil {
		return nil, err
	}
	if data.CalculatedProductData.Products == nil {
		data.CalculatedProductData.Products = map[string]draftProductJSON{}
	}
	return &data, nil
}

func (s *SalesInvoiceService) linesToProductsMap(lines []draftComputedLine) map[string]draftProductJSON {
	result := make(map[string]draftProductJSON, len(lines))
	for _, l := range lines {
		k := fmt.Sprintf("%d_%s", l.Item.ProductID, l.Item.BatchCode)
		result[k] = draftProductJSON{
			ProductID:          l.Item.ProductID,
			BatchCode:          l.Item.BatchCode,
			ExpiryDate:         l.Batch.ExpiryDate.Format(time.DateOnly),
			MRP:                l.Batch.MRP,
			SalesRate:          l.SalesRate,
			BaseRate:           l.SalesRate,
			Quantity:           l.Item.Quantity,
			GSTPercentage:      0,
			GSTAmount:          0,
			BillAmount:         l.BillAmount,
			DiscountType:       l.DiscountType,
			DiscountAmount:     l.DiscountAmt,
			DiscountPercentage: l.DiscountPct,
			IsFreeProduct:      l.Item.IsFreeProduct,
			ComboProductID:     l.Item.ComboProductID,
		}
	}
	return result
}

func (s *SalesInvoiceService) buildDraftResponse(draft model.SalesInvoiceDraftJSON, items []draftProductJSON, payments []model.SalesInvoiceDraftPayment) draftResponse {
	if items == nil {
		var payload salesInvoiceDraftJSONPayload
		if err := json.Unmarshal(draft.DraftJSON, &payload); err == nil {
			items = payload.Products
		}
	}
	pRes := make([]draftPaymentResp, 0, len(payments))
	for _, p := range payments {
		pRes = append(pRes, draftPaymentResp{ID: p.ID, Amount: p.Amount, VoucherCode: p.VoucherCode, StorePaymentMethodID: p.StorePaymentMethodID, Type: p.Type})
	}

	return draftResponse{
		ID:                   draft.ID,
		OrganizationID:       draft.OrganizationID,
		IsHomeDelivery:       draft.IsHomeDelivery,
		TotalProducts:        draft.TotalProducts,
		TotalItems:           draft.TotalItems,
		TotalQuantity:        draft.TotalQuantity,
		PrepaidAmount:        round2(draft.PrepaidAmount),
		TotalInvoiceAmount:   round2(draft.TotalInvoiceAmount),
		TotalAmountReceived:  round2(draft.TotalAmountReceived),
		AmountDue:            round2(draft.TotalInvoiceAmount - draft.PrepaidAmount - draft.TotalAmountReceived),
		DeliveryCharges:      0,
		TotalMRP:             round2(draft.TotalAmount),
		TotalSavings:         round2(draft.TotalDiscount),
		TaxableAmount:        round2(draft.TaxableAmount),
		Status:               draft.Status,
		PaymentStatus:        draft.PaymentStatus,
		RoundOff:             round2(draft.RoundOff),
		TotalBillBeforePromo: round2(draft.TotalBillBeforePromo),
		PromoCode:            draft.PromoCode,
		Items:                items,
		Payments:             pRes,
	}
}

func mapLinesToDraftProducts(lines []draftComputedLine) []draftProductJSON {
	products := make([]draftProductJSON, 0, len(lines))
	for _, l := range lines {
		products = append(products, draftProductJSON{
			ProductID:          l.Item.ProductID,
			BatchCode:          l.Item.BatchCode,
			ExpiryDate:         l.Batch.ExpiryDate.Format(time.DateOnly),
			MRP:                l.Batch.MRP,
			SalesRate:          l.SalesRate,
			BaseRate:           l.SalesRate,
			Quantity:           l.Item.Quantity,
			GSTPercentage:      0,
			GSTAmount:          0,
			BillAmount:         l.BillAmount,
			DiscountType:       l.DiscountType,
			DiscountAmount:     l.DiscountAmt,
			DiscountPercentage: l.DiscountPct,
			IsFreeProduct:      l.Item.IsFreeProduct,
			ComboProductID:     l.Item.ComboProductID,
		})
	}
	return products
}

func buildDraftPayments(input CreateOrUpdateSalesInvoiceInput, draftID uint64) (float64, []model.SalesInvoiceDraftPayment) {
	now := time.Now()
	totalReceived := 0.0
	payments := make([]model.SalesInvoiceDraftPayment, 0, len(input.Payments))
	for _, p := range input.Payments {
		if p.ID != nil && *p.ID > 0 {
			continue
		}
		amount := round2(p.Amount)
		ptype := constants.SalesPaymentTypeSales
		if p.IsAdvanceRefund {
			amount = -amount
			ptype = constants.SalesPaymentTypeAdvanceRefund
		} else {
			totalReceived += amount
		}
		payments = append(payments, model.SalesInvoiceDraftPayment{
			SalesInvoiceDraftID:  draftID,
			StoreID:              input.StoreID,
			StorePaymentMethodID: p.StorePaymentMethodID,
			Amount:               amount,
			VoucherCode:          p.VoucherCode,
			Type:                 ptype,
			DeviceMasterID:       input.DeviceMasterID,
			CreatedBy:            input.UserID,
			CreatedAt:            now,
			UpdatedAt:            now,
		})
	}
	return round2(totalReceived), payments
}

func totalReceivedFromDraftPayments(payments []model.SalesInvoiceDraftPayment) float64 {
	total := 0.0
	for _, payment := range payments {
		if payment.Type == constants.SalesPaymentTypeAdvanceRefund {
			continue
		}
		total += payment.Amount
	}
	return round2(total)
}

func hasAdvanceRefundPayment(payments []model.SalesInvoiceDraftPayment) bool {
	for _, payment := range payments {
		if payment.Type == constants.SalesPaymentTypeAdvanceRefund {
			return true
		}
	}
	return false
}

func resolveDraftStatus(hasPayments bool) string {
	if hasPayments {
		return constants.SalesInvoiceStatusPaymentPending
	}
	return constants.SalesInvoiceStatusDraft
}

func resolvePaymentStatus(totalReceived, invoiceTotal float64) string {
	totalReceived = round2(totalReceived)
	invoiceTotal = round2(invoiceTotal)
	switch {
	case totalReceived == 0:
		return constants.SalesPaymentStatusDraft
	case totalReceived >= invoiceTotal:
		return constants.SalesPaymentStatusFullyPaid
	default:
		return constants.SalesPaymentStatusPartiallyPaid
	}
}

func fulfillBatches(rows []repository.BatchStockRow, required int, productID uint64, batchCode string) ([]batchAllocation, error) {
	remaining := required
	allocs := make([]batchAllocation, 0)
	for _, row := range rows {
		if remaining <= 0 {
			break
		}
		take := minInt(row.ClosingStock, remaining)
		if take <= 0 {
			continue
		}
		allocs = append(allocs, batchAllocation{StoreBatchID: row.StoreBatchID, PurchaseRate: row.PurchaseRate, BatchCode: row.BatchCode, ExpiryDate: row.ExpiryDate, QuantityTaken: take})
		remaining -= take
	}
	if remaining > 0 {
		return nil, fmt.Errorf("not enough stock for product %d batch %s", productID, batchCode)
	}
	return allocs, nil
}

type batchAllocation struct {
	StoreBatchID  uint64
	PurchaseRate  float64
	BatchCode     string
	ExpiryDate    time.Time
	QuantityTaken int
}

func extractUniqueFromItems(items []CreateOrUpdateSalesInvoiceItem) ([]uint64, []string) {
	productSet := map[uint64]struct{}{}
	batchSet := map[string]struct{}{}
	for _, item := range items {
		productSet[item.ProductID] = struct{}{}
		batchSet[item.BatchCode] = struct{}{}
	}
	productIDs := make([]uint64, 0, len(productSet))
	for id := range productSet {
		productIDs = append(productIDs, id)
	}
	sort.Slice(productIDs, func(i, j int) bool { return productIDs[i] < productIDs[j] })

	batchCodes := make([]string, 0, len(batchSet))
	for code := range batchSet {
		batchCodes = append(batchCodes, code)
	}
	sort.Strings(batchCodes)
	return productIDs, batchCodes
}

func extractProductIDsAndBatchCodes(lines []draftComputedLine) ([]uint64, []string) {
	productSet := map[uint64]struct{}{}
	batchSet := map[string]struct{}{}
	for _, line := range lines {
		productSet[line.Item.ProductID] = struct{}{}
		batchSet[line.Item.BatchCode] = struct{}{}
	}
	productIDs := make([]uint64, 0, len(productSet))
	for id := range productSet {
		productIDs = append(productIDs, id)
	}
	sort.Slice(productIDs, func(i, j int) bool { return productIDs[i] < productIDs[j] })
	batchCodes := make([]string, 0, len(batchSet))
	for code := range batchSet {
		batchCodes = append(batchCodes, code)
	}
	sort.Strings(batchCodes)
	return productIDs, batchCodes
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}

func derefBool(v *bool) bool {
	if v == nil {
		return false
	}
	return *v
}
