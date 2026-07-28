package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"mk-pos-billing/internal/constants"
	domaincache "mk-pos-billing/internal/domain/cache"
	"mk-pos-billing/internal/domain/model"
	"mk-pos-billing/internal/domain/repository"
	"mk-pos-billing/internal/infrastructure/cache"
	"mk-pos-billing/internal/infrastructure/config"
	"sort"
	"strconv"
	"time"
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
	Items             []CreateOrUpdateSalesInvoiceItem
	Payments          []CreateOrUpdateSalesInvoicePayment
}

type CreateOrUpdateSalesInvoiceItem struct {
	ProductID      uint64
	BatchCode      string
	Quantity       int
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

type draftComputedLine struct {
	Item                 CreateOrUpdateSalesInvoiceItem
	Batch                model.Batch
	SalesRate            float64
	SalesRateBeforePromo float64
	BaseRate             float64
	BillAmount           float64
	TotalAmount          float64
	DiscountType         string
	DiscountAmt          float64
	DiscountPct          float64
	GSTPct               float64
	GSTAmount            float64
	CGST                 float64
	SGST                 float64
	IGST                 float64
	HSNCode              string
}

type DraftCalculationResult struct {
	Lines         []draftComputedLine
	TotalProducts int
	TotalItems    int
	TotalQty      int
	TotalBill     float64
	TotalAmount   float64
	TaxableAmount float64
	TotalGST      float64
	TotalCGST     float64
	TotalSGST     float64
	TotalIGST     float64
}

type draftCacheData struct {
	SalesInvoiceDraft     model.SalesInvoiceDraftJSON     `json:"salesInvoiceDraft"     php:"salesInvoiceDraft"`
	CalculatedProductData draftCalculatedProductContainer `json:"calculatedProductData" php:"calculatedProductData"`
}

type draftCalculatedProductContainer struct {
	Products map[string]model.DraftProductJSON `json:"products" php:"products"`
}

type draftResponse struct {
	ID                   uint64                   `json:"id"`
	OrganizationID       uint64                   `json:"organization_id"`
	IsHomeDelivery       bool                     `json:"is_home_delivery"`
	TotalProducts        int                      `json:"total_products"`
	TotalItems           int                      `json:"total_items"`
	TotalQuantity        int                      `json:"total_quantity"`
	PrepaidAmount        float64                  `json:"prepaid_amount"`
	TotalInvoiceAmount   float64                  `json:"total_invoice_amount"`
	TotalAmountReceived  float64                  `json:"total_amount_received"`
	AmountDue            float64                  `json:"amount_due"`
	DeliveryCharges      float64                  `json:"delivery_charges"`
	TotalMRP             float64                  `json:"total_mrp"`
	TotalSavings         float64                  `json:"total_savings"`
	TaxableAmount        float64                  `json:"taxable_amount"`
	Status               string                   `json:"status"`
	PaymentStatus        string                   `json:"payment_status"`
	RoundOff             float64                  `json:"round_off"`
	TotalBillBeforePromo float64                  `json:"total_bill_amount_before_promo"`
	PromoCode            *string                  `json:"promo_code"`
	Items                []model.DraftProductJSON `json:"items"`
	Payments             []draftPaymentResp       `json:"payments"`
}

type draftPaymentResp struct {
	ID                   uint64  `json:"id"`
	Amount               float64 `json:"amount"`
	VoucherCode          *string `json:"voucher_code"`
	StorePaymentMethodID uint64  `json:"store_payment_method_id"`
	Type                 string  `json:"type"`
}

type SalesInvoiceService struct {
	invoiceRepo   repository.SalesInvoiceRepository
	draftRepo     repository.SalesInvoiceDraftRepository
	inventoryRepo repository.StoreInventoryRepository
	productRepo   repository.ProductRepository
	masterRepo    repository.MasterDataRepository
	cache         *cache.Service
	cfg           config.SalesInvoiceConfig
	productCache  domaincache.ProductCache
	storeCache    domaincache.StoreCache
	txManager     repository.TransactionManager
}

func NewSalesInvoiceService(
	invoiceRepo repository.SalesInvoiceRepository,
	draftRepo repository.SalesInvoiceDraftRepository,
	inventoryRepo repository.StoreInventoryRepository,
	productRepo repository.ProductRepository,
	masterRepo repository.MasterDataRepository,
	cacheService *cache.Service,
	cfg config.SalesInvoiceConfig,
	productCache domaincache.ProductCache,
	storeCache domaincache.StoreCache,
	txManager repository.TransactionManager,
) *SalesInvoiceService {
	return &SalesInvoiceService{
		invoiceRepo:   invoiceRepo,
		draftRepo:     draftRepo,
		inventoryRepo: inventoryRepo,
		productRepo:   productRepo,
		masterRepo:    masterRepo,
		cache:         cacheService,
		cfg:           cfg,
		productCache:  productCache,
		storeCache:    storeCache,
		txManager:     txManager,
	}
}

func (s *SalesInvoiceService) CreateOrUpdate(ctx context.Context, input CreateOrUpdateSalesInvoiceInput) (*CreateOrUpdateSalesInvoiceOutput, error) {
	input.OrganizationID = s.cfg.DefaultOrganizationID

	if err := s.validateInput(ctx, input); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrValidation, err.Error())
	}

	var out *CreateOrUpdateSalesInvoiceOutput
	err := s.txManager.Do(ctx, func(tx repository.Transaction) error {
		txDraftRepo := s.draftRepo.Tx(tx)
		txInvoiceRepo := s.invoiceRepo.Tx(tx)
		txInventoryRepo := s.inventoryRepo.Tx(tx)
		txProductRepo := s.productRepo.Tx(tx)
		combinedCacheKey := s.buildDraftCacheKeys(input.StoreID, input.ID)

		cachedData, _ := s.getDraftCache(ctx, combinedCacheKey)

		draft, err := s.loadOrCreateDraft(ctx, txDraftRepo, input, cachedData)
		if err != nil {
			return err
		}

		if input.ID != nil && len(input.Items) == 0 {
			out, err = s.handlePaymentOnlyUpdate(ctx, txDraftRepo, input, draft, combinedCacheKey)
			return err
		}

		calcResult, err := s.calculateDraftLines(ctx, txInventoryRepo, txProductRepo, input, cachedData)
		if err != nil {
			return err
		}

		existingPayments, err := txDraftRepo.GetDraftPayments(ctx, draft.ID)
		if err != nil {
			return err
		}
		newPaymentTotal, draftPayments := buildDraftPayments(input, draft.ID)

		state := s.computeDraftState(calcResult, existingPayments, draftPayments, newPaymentTotal)

		s.updateDraftMetadata(draft, input, state, calcResult)

		draft.DraftJSON, err = s.buildDraftPayload(ctx, input, calcResult, state)
		if err != nil {
			return err
		}

		persistedPayments, err := s.persistDraft(ctx, txDraftRepo, input, draft, calcResult, draftPayments, combinedCacheKey)
		if err != nil {
			return err
		}

		msg := "Sales Invoice (Draft) Created Successfully"
		if input.ID != nil && *input.ID > 0 {
			msg = "Sales Invoice (Draft) Updated Successfully"
		}

		if !derefBool(input.IsConfirmed) {
			out = &CreateOrUpdateSalesInvoiceOutput{Data: s.buildDraftResponse(*draft, mapLinesToDraftProducts(calcResult.Lines), persistedPayments), Message: msg}
			return nil
		}

		invoiceID, err := s.finalizeDraft(ctx, txInvoiceRepo, txInventoryRepo, txDraftRepo, input, draft, calcResult, persistedPayments, combinedCacheKey)
		if err != nil {
			return err
		}

		out = &CreateOrUpdateSalesInvoiceOutput{Data: invoiceID, Message: "Sales Invoice Created Successfully"}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return out, nil
}

type draftStateTotals struct {
	TotalInvoice  float64
	RoundOff      float64
	TotalDiscount float64
	TotalReceived float64
	PaymentStatus string
	DraftStatus   string
}

func (s *SalesInvoiceService) validateInput(ctx context.Context, input CreateOrUpdateSalesInvoiceInput) error {
	if !s.masterRepo.ExistsActiveUser(ctx, input.BillingUserID) {
		return errors.New("billing_user_id is invalid")
	}

	hasPayments := len(input.Payments) > 0
	if derefBool(input.IsConfirmed) || hasPayments {
		if input.CustomerID == nil || !s.masterRepo.ExistsActiveCustomer(ctx, *input.CustomerID) {
			return errors.New("customer_id is invalid")
		}
	}

	if derefBool(input.IsConfirmed) {
		if input.PatientID == nil || input.CustomerID == nil || !s.masterRepo.ExistsPatientForCustomer(ctx, *input.PatientID, *input.CustomerID) {
			return errors.New("patient_id is invalid for customer")
		}
	}

	if hasPayments {
		if input.DoctorID == nil || !s.masterRepo.ExistsDoctor(ctx, *input.DoctorID) {
			return errors.New("doctor_id is invalid")
		}
	}

	if hasPayments && derefBool(input.IsHomeDelivery) {
		if input.CustomerAddressID == nil || input.CustomerID == nil || !s.masterRepo.ExistsCustomerAddress(ctx, *input.CustomerAddressID, *input.CustomerID, true) {
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

		if !s.productRepo.ExistsActiveProduct(ctx, item.ProductID) {
			return errors.New("items." + strconv.Itoa(i) + ".product_id is invalid")
		}
		if !s.productRepo.ExistsBatchCode(ctx, item.BatchCode) {
			return errors.New("items." + strconv.Itoa(i) + ".batch_code is invalid")
		}

		pData, err := s.productRepo.GetProductValidationData(ctx, item.ProductID)
		if err != nil {
			return err
		}
		if pData == nil {
			return errors.New("items." + strconv.Itoa(i) + ".product_id is invalid")
		}
		pCacheData, err := s.productCache.Get(ctx, int(item.ProductID))
		if err != nil {
			return err
		}
		if pData.SalesUnit > 0 && item.Quantity%pData.SalesUnit != 0 {
			return errors.New("the quantity must be a multiple of the product's sales unit")
		}
		if pCacheData.IsMspProduct && item.Quantity != 1 {
			return errors.New("the quantity for a loyalty program product must be one")
		}
		if pData.WSCode == s.cfg.DeliveryChargeProductWSCode && item.Quantity != 1 {
			return errors.New("the quantity for the product for delivery charge must be one")
		}

		productIDs = append(productIDs, item.ProductID)
	}

	if derefBool(input.IsConfirmed) && s.productRepo.HasNarcoticsProduct(ctx, productIDs, s.cfg.NarcoticsProductScheduledType) {
		if input.CourseDays == nil {
			return errors.New("course_days is required for narcotics products")
		}
	}

	for i, p := range input.Payments {
		if !s.masterRepo.ExistsStorePaymentMethod(ctx, p.StorePaymentMethodID, input.StoreID, input.OrganizationID) {
			return errors.New("payments." + strconv.Itoa(i) + ".store_payment_method_id is invalid")
		}
	}

	if input.ASMUserID != nil && !s.masterRepo.ExistsActiveUser(ctx, *input.ASMUserID) {
		return errors.New("invalid passkey")
	}

	return nil
}
func (s *SalesInvoiceService) buildDraftCacheKeys(storeID uint64, draftID *uint64) string {
	id := uint64(0)
	if draftID != nil {
		id = *draftID
	}
	key := s.cfg.PrefixDraftBillCache + strconv.FormatUint(id, 10)
	tag := s.cfg.PrefixDraftBillCacheTags + strconv.FormatUint(storeID, 10)
	return tag + ":" + key
}

func (s *SalesInvoiceService) getDraftCache(ctx context.Context, combinedKey string) (*draftCacheData, error) {
	var data draftCacheData
	err := s.cache.GetPHPSerialized(ctx, combinedKey, &data)
	if err != nil {
		return nil, err
	}
	if data.CalculatedProductData.Products == nil {
		data.CalculatedProductData.Products = map[string]model.DraftProductJSON{}
	}
	return &data, nil
}

func (s *SalesInvoiceService) loadOrCreateDraft(ctx context.Context, draftRepo repository.SalesInvoiceDraftRepository, input CreateOrUpdateSalesInvoiceInput, cachedData *draftCacheData) (*model.SalesInvoiceDraftJSON, error) {
	if cachedData != nil && cachedData.SalesInvoiceDraft.ID != 0 {
		draft := cachedData.SalesInvoiceDraft
		return &draft, nil
	}

	if input.ID != nil && *input.ID > 0 {
		return draftRepo.GetDraftByID(ctx, *input.ID, input.StoreID, input.OrganizationID)
	}
	return &model.SalesInvoiceDraftJSON{}, nil
}

func (s *SalesInvoiceService) handlePaymentOnlyUpdate(
	ctx context.Context,
	txDraftRepo repository.SalesInvoiceDraftRepository,
	input CreateOrUpdateSalesInvoiceInput,
	draft *model.SalesInvoiceDraftJSON,
	combinedCacheKey string,
) (*CreateOrUpdateSalesInvoiceOutput, error) {
	existingPayments, err := txDraftRepo.GetDraftPayments(ctx, draft.ID)
	if err != nil {
		return nil, err
	}
	_, draftPayments := buildDraftPayments(input, draft.ID)
	if err := txDraftRepo.AppendDraftPayments(ctx, draftPayments); err != nil {
		return nil, err
	}
	if len(draftPayments) > 0 {
		existingPayments, err = txDraftRepo.GetDraftPayments(ctx, draft.ID)
		if err != nil {
			return nil, err
		}
	}

	draft.DraftJSON.Products = []model.DraftProductJSON{}
	draft.DraftJSON.TotalProducts = 0
	draft.DraftJSON.TotalItems = 0
	draft.DraftJSON.TotalQuantity = 0
	draft.DraftJSON.TotalGST = 0
	draft.DraftJSON.SGST = 0
	draft.DraftJSON.CGST = 0
	draft.DraftJSON.IGST = 0
	draft.DraftJSON.DeliveryCharges = 0
	draft.DraftJSON.TaxableAmount = 0
	draft.DraftJSON.RoundOff = 0
	draft.TotalAmount = 0
	draft.TotalBillAmount = 0
	draft.TotalInvoiceAmount = 0
	draft.TotalAmountReceived = totalReceivedFromDraftPayments(existingPayments)
	draft.TotalDiscount = 0
	draft.RoundOff = 0
	draft.Status = resolveDraftStatus(len(existingPayments) > 0 || len(draftPayments) > 0)
	draft.PaymentStatus = resolvePaymentStatus(draft.TotalAmountReceived, 0)
	draft.UpdatedBy = &input.UserID
	if err := txDraftRepo.SaveDraft(ctx, draft); err != nil {
		return nil, err
	}

	cacheData := draftCacheData{SalesInvoiceDraft: *draft, CalculatedProductData: draftCalculatedProductContainer{Products: map[string]model.DraftProductJSON{}}}
	_ = s.cache.SetPHPSerialized(ctx, combinedCacheKey, cacheData, time.Duration(s.cfg.DraftCacheTTLMinutes)*time.Minute)

	out := &CreateOrUpdateSalesInvoiceOutput{Data: s.buildDraftResponse(*draft, nil, existingPayments), Message: "Sales Invoice (Draft) Updated Successfully"}
	return out, nil
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

func (s *SalesInvoiceService) buildDraftResponse(draft model.SalesInvoiceDraftJSON, items []model.DraftProductJSON, payments []model.SalesInvoiceDraftPayment) draftResponse {
	payload := draft.DraftJSON
	if items == nil {
		items = payload.Products
	}
	pRes := make([]draftPaymentResp, 0, len(payments))
	for _, p := range payments {
		pRes = append(pRes, draftPaymentResp{ID: p.ID, Amount: p.Amount, VoucherCode: p.VoucherCode, StorePaymentMethodID: p.StorePaymentMethodID, Type: p.Type})
	}

	return draftResponse{
		ID:                   draft.ID,
		OrganizationID:       draft.OrganizationID,
		IsHomeDelivery:       payload.IsHomeDelivery,
		TotalProducts:        payload.TotalProducts,
		TotalItems:           payload.TotalItems,
		TotalQuantity:        payload.TotalQuantity,
		PrepaidAmount:        round2(draft.PrepaidAmount),
		TotalInvoiceAmount:   round2(draft.TotalInvoiceAmount),
		TotalAmountReceived:  round2(draft.TotalAmountReceived),
		AmountDue:            round2(draft.TotalInvoiceAmount - draft.PrepaidAmount - draft.TotalAmountReceived),
		DeliveryCharges:      payload.DeliveryCharges,
		TotalMRP:             round2(draft.TotalAmount),
		TotalSavings:         round2(draft.TotalDiscount),
		TaxableAmount:        round2(payload.TaxableAmount),
		Status:               draft.Status,
		PaymentStatus:        draft.PaymentStatus,
		RoundOff:             round2(payload.RoundOff),
		TotalBillBeforePromo: round2(payload.TotalBillAmountBeforePromo),
		PromoCode:            payload.PromoCode,
		Items:                items,
		Payments:             pRes,
	}
}

func (s *SalesInvoiceService) calculateDraftLines(ctx context.Context, inventoryRepo repository.StoreInventoryRepository, productRepo repository.ProductRepository, input CreateOrUpdateSalesInvoiceInput, cachedData *draftCacheData) (DraftCalculationResult, error) {
	var emptyResult DraftCalculationResult

	productIDs, batchCodes := extractUniqueFromItems(input.Items)
	batchMeta, err := productRepo.FetchBatchMeta(ctx, productIDs, batchCodes)
	if err != nil {
		return emptyResult, err
	}
	stockRows, err := inventoryRepo.FetchBatchStocks(ctx, input.StoreID, productIDs, batchCodes, false, s.batchExpiryCutoff())
	if err != nil {
		return emptyResult, err
	}

	storeCache, err := s.storeCache.Get(ctx, int(input.StoreID))
	if err != nil {
		return emptyResult, fmt.Errorf("store cache unavailable: %w", err)
	}

	// Pre-fetch generic pricing templates for B2C default org
	var templateIDs []uint64
	var genericProductIDs []uint64
	if input.OrganizationID == s.cfg.DefaultOrganizationID {
		for _, item := range input.Items {
			productCache, err := s.productCache.Get(ctx, int(item.ProductID))
			if err != nil {
				continue
			}
			if productCache.IsGeneric {
				category := productCache.B2CProductCategoryID
				if discountData, ok := storeCache.ProductCategoriesDiscounts[category]; ok {
					templateIDs = append(templateIDs, uint64(discountData.B2CPricingTemplateID))
					genericProductIDs = append(genericProductIDs, item.ProductID)
				}
			}
		}
	}

	genericPricingMap, err := productRepo.FetchGenericPricings(ctx, templateIDs, genericProductIDs)
	if err != nil {
		return emptyResult, fmt.Errorf("failed to fetch generic pricings: %w", err)
	}

	result := DraftCalculationResult{
		Lines: make([]draftComputedLine, 0, len(input.Items)),
	}

	for _, item := range input.Items {
		key := fmt.Sprintf("%d_%s", item.ProductID, item.BatchCode)
		batch, ok := batchMeta[key]
		if !ok {
			return emptyResult, fmt.Errorf("batch not found for product %d batch %s", item.ProductID, item.BatchCode)
		}

		available := 0
		for _, r := range stockRows[key] {
			available += r.ClosingStock
		}
		if available < item.Quantity {
			return emptyResult, fmt.Errorf("requested quantity (%d) for batch (%s) is more than available quantity (%d)", item.Quantity, item.BatchCode, available)
		}

		productCache, err := s.productCache.Get(ctx, int(item.ProductID))
		if err != nil {
			return emptyResult, fmt.Errorf("product cache unavailable for product %d: %w", item.ProductID, err)
		}

		salesRateResult, err := GetSalesRate(input.OrganizationID, productCache, storeCache, batch.MRP, s.cfg.DefaultOrganizationID, genericPricingMap)
		if err != nil {
			return emptyResult, fmt.Errorf("failed to calculate sales rate for product %d: %w", item.ProductID, err)
		}

		salesRate := salesRateResult.SalesRate

		if salesRate-batch.MRP > 0.1 {
			return emptyResult, fmt.Errorf("the sales rate of batch (%s) is greater than MRP", item.BatchCode)
		}

		billAmount := round2(float64(item.Quantity) * salesRate)
		totalLineAmount := round2(float64(item.Quantity) * batch.MRP)
		discountAmt := round2(salesRateResult.DiscountAmount)
		discountPct := round2(salesRateResult.DiscountPercentage)

		tax, err := CalculateSalesInvoiceItemTax(storeCache, productCache, billAmount)
		if err != nil {
			return emptyResult, fmt.Errorf("tax calculation failed for product %d: %w", item.ProductID, err)
		}

		result.TotalBill += billAmount
		result.TotalAmount += totalLineAmount
		result.TotalQty += item.Quantity
		result.TotalGST += tax.TotalGST
		result.TotalCGST += tax.CGST
		result.TotalSGST += tax.SGST
		result.TotalIGST += tax.IGST

		result.Lines = append(result.Lines, draftComputedLine{
			Item:                 item,
			Batch:                batch,
			SalesRate:            salesRate,
			SalesRateBeforePromo: salesRateResult.SalesRateBeforePromo,
			BaseRate:             round2(tax.BaseRate),
			BillAmount:           billAmount,
			TotalAmount:          totalLineAmount,
			DiscountType:         salesRateResult.DiscountType,
			DiscountAmt:          discountAmt,
			DiscountPct:          discountPct,
			GSTPct:               tax.GSTPct,
			GSTAmount:            round2(tax.TotalGST),
			CGST:                 round2(tax.CGST),
			SGST:                 round2(tax.SGST),
			IGST:                 round2(tax.IGST),
			HSNCode:              productCache.HsnCode,
		})
	}

	result.TotalProducts = len(productIDs)
	result.TotalItems = len(result.Lines)
	result.TaxableAmount = round2(result.TotalBill - result.TotalGST)

	return result, nil
}

func (s *SalesInvoiceService) batchExpiryCutoff() time.Time {
	return time.Now().AddDate(0, 0, s.cfg.DefaultMinimumDaysForBatch).Truncate(24 * time.Hour)
}

func (s *SalesInvoiceService) computeDraftState(
	calcResult DraftCalculationResult,
	existingPayments []model.SalesInvoiceDraftPayment,
	draftPayments []model.SalesInvoiceDraftPayment,
	newPaymentTotal float64,
) draftStateTotals {
	totalInvoice := math.Round(calcResult.TotalBill)
	totalReceived := round2(totalReceivedFromDraftPayments(existingPayments) + newPaymentTotal)

	return draftStateTotals{
		TotalInvoice:  totalInvoice,
		RoundOff:      totalInvoice - calcResult.TotalBill,
		TotalDiscount: calcResult.TotalAmount - calcResult.TotalBill,
		TotalReceived: totalReceived,
		PaymentStatus: resolvePaymentStatus(totalReceived, totalInvoice),
		DraftStatus:   resolveDraftStatus(len(existingPayments) > 0 || len(draftPayments) > 0),
	}
}

func (s *SalesInvoiceService) updateDraftMetadata(
	draft *model.SalesInvoiceDraftJSON,
	input CreateOrUpdateSalesInvoiceInput,
	state draftStateTotals,
	calcResult DraftCalculationResult,
) {
	draft.StoreID = input.StoreID
	draft.OrganizationID = input.OrganizationID
	draft.BillingUserID = input.BillingUserID
	draft.CustomerID = input.CustomerID
	draft.CustomerAddressID = input.CustomerAddressID
	draft.DoctorID = input.DoctorID
	draft.PatientID = input.PatientID
	draft.Status = state.DraftStatus
	draft.PaymentStatus = state.PaymentStatus
	draft.TotalBillAmount = calcResult.TotalBill
	draft.PrepaidAmount = 0
	draft.RoundOff = state.RoundOff
	draft.TotalInvoiceAmount = state.TotalInvoice
	draft.TotalAmount = calcResult.TotalAmount
	draft.TotalDiscount = state.TotalDiscount
	draft.TotalAmountReceived = state.TotalReceived
	draft.TillID = input.TillID
	draft.TillTransactionID = input.TillTransactionID
	draft.UpdatedBy = &input.UserID
	draft.UpdatedAt = time.Now()
}

func (s *SalesInvoiceService) buildDraftPayload(
	ctx context.Context,
	input CreateOrUpdateSalesInvoiceInput,
	calcResult DraftCalculationResult,
	state draftStateTotals,
) (model.DraftJSONPayload, error) {
	isHomeDelivery := false
	if input.IsHomeDelivery != nil {
		isHomeDelivery = *input.IsHomeDelivery
	}

	storeInfo, err := s.storeCache.Get(ctx, int(input.StoreID))
	if err != nil {
		return model.DraftJSONPayload{}, err
	}

	draftPayload := model.DraftJSONPayload{
		IsHomeDelivery:             isHomeDelivery,
		IsActive:                   true,
		PaymentStatus:              state.PaymentStatus,
		DeviceMasterID:             input.DeviceMasterID,
		PromoCode:                  input.PromoCode,
		Notes:                      input.Notes,
		TotalProducts:              calcResult.TotalProducts,
		TotalItems:                 calcResult.TotalItems,
		TotalQuantity:              calcResult.TotalQty,
		TotalGST:                   calcResult.TotalGST,
		SGST:                       calcResult.TotalSGST,
		CGST:                       calcResult.TotalCGST,
		IGST:                       calcResult.TotalIGST,
		DeliveryCharges:            0,
		TaxableAmount:              calcResult.TaxableAmount,
		RoundOff:                   state.RoundOff,
		CINNumber:                  storeInfo.CinNumber,
		GSTNumber:                  &storeInfo.GstNumber,
		GSTTreatment:               storeInfo.GstTreatment,
		PlaceOfSupplyCode:          storeInfo.PlaceOfSupplyCode,
		TotalAmount:                calcResult.TotalAmount,
		TotalBillAmount:            calcResult.TotalBill,
		TotalInvoiceAmount:         state.TotalInvoice,
		TotalBillAmountBeforePromo: calcResult.TotalBill,
		Products:                   mapLinesToDraftProducts(calcResult.Lines),
	}
	return draftPayload, nil
}

func mapLinesToDraftProducts(lines []draftComputedLine) []model.DraftProductJSON {
	products := make([]model.DraftProductJSON, 0, len(lines))
	for _, l := range lines {
		var taxDetails []model.DraftTaxDetailJSON
		if l.IGST > 0 {
			taxDetails = append(taxDetails, model.DraftTaxDetailJSON{
				TaxName:   constants.TaxNameIGST,
				TaxType:   constants.TaxTypeAmount,
				TaxRate:   l.GSTPct,
				TaxAmount: l.IGST,
			})
		} else if l.CGST > 0 || l.SGST > 0 {
			taxDetails = append(taxDetails, model.DraftTaxDetailJSON{
				TaxName:   constants.TaxNameSGST,
				TaxType:   constants.TaxTypeAmount,
				TaxRate:   l.GSTPct / 2,
				TaxAmount: l.SGST,
			})
			taxDetails = append(taxDetails, model.DraftTaxDetailJSON{
				TaxName:   constants.TaxNameCGST,
				TaxType:   constants.TaxTypeAmount,
				TaxRate:   l.GSTPct / 2,
				TaxAmount: l.CGST,
			})
		}

		products = append(products, model.DraftProductJSON{
			ProductID:                     l.Item.ProductID,
			BatchCode:                     l.Item.BatchCode,
			ExpiryDate:                    l.Batch.ExpiryDate.Format(time.DateOnly),
			MRP:                           l.Batch.MRP,
			SalesRate:                     l.SalesRate,
			BaseRate:                      l.BaseRate,
			Quantity:                      l.Item.Quantity,
			GSTPercentage:                 l.GSTPct,
			GSTAmount:                     l.GSTAmount,
			BillAmount:                    l.BillAmount,
			DiscountType:                  l.DiscountType,
			DiscountAmount:                l.DiscountAmt,
			DiscountPercentage:            l.DiscountPct,
			IsFreeProduct:                 l.Item.IsFreeProduct,
			ComboProductID:                l.Item.ComboProductID,
			TotalAmount:                   round2(l.Batch.MRP * float64(l.Item.Quantity)),
			SalesRateBeforePromo:          l.SalesRate,
			DiscountAmountBeforePromo:     l.DiscountAmt,
			DiscountPercentageBeforePromo: l.DiscountPct,
			TaxDetails:                    taxDetails,
		})
	}
	return products
}

func (s *SalesInvoiceService) persistDraft(
	ctx context.Context,
	txDraftRepo repository.SalesInvoiceDraftRepository,
	input CreateOrUpdateSalesInvoiceInput,
	draft *model.SalesInvoiceDraftJSON,
	calcResult DraftCalculationResult,
	draftPayments []model.SalesInvoiceDraftPayment,
	combinedCacheKey string,
) ([]model.SalesInvoiceDraftPayment, error) {
	if draft.ID == 0 {
		draft.CreatedBy = input.UserID
		draft.CreatedAt = draft.UpdatedAt
		if err := txDraftRepo.CreateDraft(ctx, draft); err != nil {
			return nil, err
		}
		combinedCacheKey = s.buildDraftCacheKeys(input.StoreID, &draft.ID)
	} else {
		if err := txDraftRepo.SaveDraft(ctx, draft); err != nil {
			return nil, err
		}
	}

	for i := range draftPayments {
		draftPayments[i].SalesInvoiceDraftID = draft.ID
	}
	if err := txDraftRepo.AppendDraftPayments(ctx, draftPayments); err != nil {
		return nil, err
	}

	persistedPayments, _ := txDraftRepo.GetDraftPayments(ctx, draft.ID)
	cacheData := draftCacheData{
		SalesInvoiceDraft:     *draft,
		CalculatedProductData: draftCalculatedProductContainer{Products: s.linesToProductsMap(calcResult.Lines)},
	}
	_ = s.cache.SetPHPSerialized(ctx, combinedCacheKey, cacheData, time.Duration(s.cfg.DraftCacheTTLMinutes)*time.Minute)

	return persistedPayments, nil
}

func (s *SalesInvoiceService) linesToProductsMap(lines []draftComputedLine) map[string]model.DraftProductJSON {
	products := mapLinesToDraftProducts(lines)
	result := make(map[string]model.DraftProductJSON, len(products))
	for _, p := range products {
		k := fmt.Sprintf("%d_%s", p.ProductID, p.BatchCode)
		result[k] = p
	}
	return result
}

func (s *SalesInvoiceService) finalizeDraft(
	ctx context.Context,
	txInvoiceRepo repository.SalesInvoiceRepository,
	txInventoryRepo repository.StoreInventoryRepository,
	txDraftRepo repository.SalesInvoiceDraftRepository,
	input CreateOrUpdateSalesInvoiceInput,
	draft *model.SalesInvoiceDraftJSON,
	calcResult DraftCalculationResult,
	persistedPayments []model.SalesInvoiceDraftPayment,
	combinedCacheKey string,
) (uint64, error) {
	if !hasAdvanceRefundPayment(persistedPayments) && round2(draft.TotalAmountReceived+draft.PrepaidAmount) != round2(draft.TotalInvoiceAmount) {
		return 0, fmt.Errorf("entered amount is more then the bill amount, please enter proper amount")
	}

	invoiceID, err := s.finalizeInvoice(ctx, txInvoiceRepo, txInventoryRepo, txDraftRepo, input, draft, calcResult.Lines)
	if err != nil {
		return 0, err
	}

	// cache final state as draft slot with updated status
	draft.Status = constants.SalesInvoiceStatusInvoiced
	cacheData := draftCacheData{
		SalesInvoiceDraft:     *draft,
		CalculatedProductData: draftCalculatedProductContainer{Products: s.linesToProductsMap(calcResult.Lines)},
	}
	_ = s.cache.SetPHPSerialized(ctx, combinedCacheKey, cacheData, time.Duration(s.cfg.DraftCacheTTLMinutes)*time.Minute)

	return invoiceID, nil
}

func hasAdvanceRefundPayment(payments []model.SalesInvoiceDraftPayment) bool {
	for _, payment := range payments {
		if payment.Type == constants.SalesPaymentTypeAdvanceRefund {
			return true
		}
	}
	return false
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
