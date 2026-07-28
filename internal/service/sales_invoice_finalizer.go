package service

import (
	"context"
	"fmt"
	"time"

	"mk-pos-billing/internal/constants"
	"mk-pos-billing/internal/domain/model"
	"mk-pos-billing/internal/domain/repository"
)

func (s *SalesInvoiceService) finalizeInvoice(
	ctx context.Context,
	invoiceRepo repository.SalesInvoiceRepository,
	inventoryRepo repository.StoreInventoryRepository,
	draftRepo repository.SalesInvoiceDraftRepository,
	input CreateOrUpdateSalesInvoiceInput,
	draft *model.SalesInvoiceDraftJSON,
	lines []draftComputedLine,
) (uint64, error) {
	invoice, err := s.buildAndSaveInvoice(ctx, invoiceRepo, input, draft)
	if err != nil {
		return 0, err
	}

	details, err := s.processInvoiceDetailsAndInventory(ctx, invoiceRepo, inventoryRepo, input, invoice, lines)
	if err != nil {
		return 0, err
	}

	if err := s.processInvoiceTaxDetails(ctx, invoiceRepo, input, invoice.ID, invoice.StoreID, details, lines); err != nil {
		return 0, err
	}

	if err := s.processInvoicePayments(ctx, invoiceRepo, draftRepo, input, invoice.ID, invoice.StoreID, draft.ID); err != nil {
		return 0, err
	}

	if err := draftRepo.MarkDraftInvoiced(ctx, draft.ID, input.UserID); err != nil {
		return 0, err
	}

	return invoice.ID, nil
}

func (s *SalesInvoiceService) buildAndSaveInvoice(
	ctx context.Context,
	invoiceRepo repository.SalesInvoiceRepository,
	input CreateOrUpdateSalesInvoiceInput,
	draft *model.SalesInvoiceDraftJSON,
) (model.SalesInvoice, error) {
	draftPayload := draft.DraftJSON

	var gstTreatment *string
	if draftPayload.GSTTreatment != "" {
		gstTreatment = &draftPayload.GSTTreatment
	}
	var placeOfSupplyCode *string
	if draftPayload.PlaceOfSupplyCode != "" {
		placeOfSupplyCode = &draftPayload.PlaceOfSupplyCode
	}

	invoice := model.SalesInvoice{
		OrganizationID:              draft.OrganizationID,
		SalesInvoiceDraftID:         draft.ID,
		StoreID:                     draft.StoreID,
		BillingUserID:               draft.BillingUserID,
		CustomerID:                  draft.CustomerID,
		CustomerAddressID:           draft.CustomerAddressID,
		DoctorID:                    draft.DoctorID,
		PatientID:                   draft.PatientID,
		OrderType:                   constants.SalesPaymentTypeSales,
		IsHomeDelivery:              draftPayload.IsHomeDelivery,
		PaymentStatus:               &draft.PaymentStatus,
		TotalBillAmount:             draftPayload.TotalBillAmount,
		TaxableAmount:               draftPayload.TaxableAmount,
		TotalAmountBeforeDisc:       draftPayload.TotalBillAmount,
		DiscountType:                "INR",
		DiscountPercentage:          0,
		DiscountAmount:              draftPayload.TotalAmount - draftPayload.TotalBillAmount,
		IGST:                        draftPayload.IGST,
		CGST:                        draftPayload.CGST,
		SGST:                        draftPayload.SGST,
		PrepaidAmount:               draft.PrepaidAmount,
		TotalGST:                    draftPayload.TotalGST,
		RoundOff:                    draftPayload.RoundOff,
		TotalInvoiceAmount:          draftPayload.TotalInvoiceAmount,
		TotalProducts:               draftPayload.TotalProducts,
		TotalItems:                  draftPayload.TotalItems,
		TotalQuantity:               draftPayload.TotalQuantity,
		TotalAmount:                 draftPayload.TotalAmount,
		TotalDiscount:               draft.TotalDiscount,
		TotalAmountReceived:         draft.TotalAmountReceived,
		TotalBillBeforePromo:        draftPayload.TotalBillAmountBeforePromo,
		PromoCode:                   draftPayload.PromoCode,
		Notes:                       draftPayload.Notes,
		IsActive:                    true,
		VoucherDiscountAmount:       0,
		LoyaltyPoints:               0,
		LoyaltyProgramDiscount:      draft.LoyaltyProgramDiscount,
		IsLoyaltyUpdated:            false,
		IsCustomerUpdated:           false,
		IsSyncBill:                  false,
		IsPrescriptionRequired:      draftPayload.IsPrescriptionRequired,
		DeliveryCharges:             draftPayload.DeliveryCharges,
		IsPrescriptionUploaded:      false,
		PrescriptionUploadedBy:      nil,
		GSTTreatment:                gstTreatment,
		GSTNumber:                   draftPayload.GSTNumber,
		CINNumber:                   draftPayload.CINNumber,
		PlaceOfSupplyCode:           placeOfSupplyCode,
		EditedUserID:                nil,
		WhatsAppSentCount:           0,
		IsRecommendationDataUpdated: false,
		PrescriptionID:              draft.PrescriptionID,
		CreatedBy:                   input.UserID,
		DeviceMasterID:              draftPayload.DeviceMasterID,
	}
	err := invoiceRepo.CreateInvoice(ctx, &invoice)
	return invoice, err
}

func (s *SalesInvoiceService) processInvoiceDetailsAndInventory(
	ctx context.Context,
	invoiceRepo repository.SalesInvoiceRepository,
	inventoryRepo repository.StoreInventoryRepository,
	input CreateOrUpdateSalesInvoiceInput,
	invoice model.SalesInvoice,
	lines []draftComputedLine,
) ([]model.SalesInvoiceDetail, error) {
	productIDs, batchCodes := extractProductIDsAndBatchCodes(lines)
	batchStocks, err := inventoryRepo.FetchBatchStocks(ctx, input.StoreID, productIDs, batchCodes, false, s.batchExpiryCutoff())
	if err != nil {
		return nil, err
	}

	details := make([]model.SalesInvoiceDetail, 0)
	txns := make([]model.StoreInventoryTransaction, 0)
	for _, line := range lines {
		key := fmt.Sprintf("%d_%s", line.Item.ProductID, line.Item.BatchCode)
		allocations, err := fulfillBatches(batchStocks[key], line.Item.Quantity, line.Item.ProductID, line.Item.BatchCode)
		if err != nil {
			return nil, err
		}

		for _, alloc := range allocations {
			detail := model.SalesInvoiceDetail{
				SalesInvoiceID:       invoice.ID,
				StoreID:              invoice.StoreID,
				ProductID:            line.Item.ProductID,
				StoreBatchID:         &alloc.StoreBatchID,
				PurchaseRate:         alloc.PurchaseRate,
				BatchCode:            alloc.BatchCode,
				ExpiryDate:           alloc.ExpiryDate,
				MRP:                  line.Batch.MRP,
				SalesRate:            line.SalesRate,
				BaseRate:             line.BaseRate,
				BillAmount:           round2(float64(alloc.QuantityTaken) * line.SalesRate),
				Quantity:             alloc.QuantityTaken,
				DiscountType:         line.DiscountType,
				DiscountPercentage:   line.DiscountPct,
				DiscountAmount:       line.DiscountAmt,
				GSTPercentage:        line.GSTPct,
				GSTAmount:            line.GSTAmount,
				TotalAmount:          round2(float64(alloc.QuantityTaken) * line.Batch.MRP),
				AmountBeforeDiscount: line.SalesRate,
				SalesRateBeforePromo: &line.SalesRateBeforePromo,
				CreatedBy:            input.UserID,
				DeviceMasterID:       input.DeviceMasterID,
				HSNCode:              &line.HSNCode,
				IsFreeProduct:        &line.Item.IsFreeProduct,
			}
			details = append(details, detail)

			txn := model.StoreInventoryTransaction{
				StoreID:         input.StoreID,
				ProductID:       line.Item.ProductID,
				BatchCode:       line.Item.BatchCode,
				StoreBatchID:    alloc.StoreBatchID,
				ExpiryDate:      alloc.ExpiryDate,
				Quantity:        -alloc.QuantityTaken,
				Rate:            line.SalesRate,
				TotalAmount:     round2(float64(-alloc.QuantityTaken) * line.SalesRate),
				VoucherType:     "SALES_INVOICE",
				VoucherID:       invoice.ID,
				CreatedBy:       input.UserID,
				TransactionTime: time.Now(),
			}
			txns = append(txns, txn)
		}
	}

	if err := invoiceRepo.CreateInvoiceDetails(ctx, &details); err != nil {
		return nil, err
	}

	if err := inventoryRepo.InsertInventoryTransactions(ctx, txns); err != nil {
		return nil, err
	}

	return details, nil
}

func (s *SalesInvoiceService) processInvoiceTaxDetails(
	ctx context.Context,
	invoiceRepo repository.SalesInvoiceRepository,
	input CreateOrUpdateSalesInvoiceInput,
	invoiceID uint64,
	storeID uint64,
	details []model.SalesInvoiceDetail,
	lines []draftComputedLine,
) error {
	lineMap := make(map[string]*draftComputedLine, len(lines))
	for i := range lines {
		l := &lines[i]
		k := fmt.Sprintf("%d_%s", l.Item.ProductID, l.Item.BatchCode)
		lineMap[k] = l
	}

	var allTaxDetails []model.SalesInvoiceTaxDetail
	for _, detail := range details {
		k := fmt.Sprintf("%d_%s", detail.ProductID, detail.BatchCode)
		if matchedLine, ok := lineMap[k]; ok {
			var taxDetails []model.SalesInvoiceTaxDetail
			now := time.Now()
			if matchedLine.IGST > 0 {
				taxAmount := round2(float64(detail.Quantity) / float64(matchedLine.Item.Quantity) * matchedLine.IGST)
				if taxAmount > 0 {
					taxDetails = append(taxDetails, model.SalesInvoiceTaxDetail{
						SalesInvoiceDetailID: detail.ID,
						StoreID:              storeID,
						TaxType:              constants.TaxTypeAmount,
						TaxRate:              matchedLine.GSTPct,
						TaxName:              constants.TaxNameIGST,
						TaxAmount:            taxAmount,
						IsActive:             true,
						CreatedBy:            input.UserID,
						CreatedAt:            now,
					})
				}
			} else if matchedLine.CGST > 0 || matchedLine.SGST > 0 {
				sgstAmount := round2(float64(detail.Quantity) / float64(matchedLine.Item.Quantity) * matchedLine.SGST)
				if sgstAmount > 0 {
					taxDetails = append(taxDetails, model.SalesInvoiceTaxDetail{
						SalesInvoiceDetailID: detail.ID,
						StoreID:              storeID,
						TaxType:              constants.TaxTypeAmount,
						TaxRate:              matchedLine.GSTPct / 2,
						TaxName:              constants.TaxNameSGST,
						TaxAmount:            sgstAmount,
						IsActive:             true,
						CreatedBy:            input.UserID,
						CreatedAt:            now,
					})
				}
				cgstAmount := round2(float64(detail.Quantity) / float64(matchedLine.Item.Quantity) * matchedLine.CGST)
				if cgstAmount > 0 {
					taxDetails = append(taxDetails, model.SalesInvoiceTaxDetail{
						SalesInvoiceDetailID: detail.ID,
						StoreID:              storeID,
						TaxType:              constants.TaxTypeAmount,
						TaxRate:              matchedLine.GSTPct / 2,
						TaxName:              constants.TaxNameCGST,
						TaxAmount:            cgstAmount,
						IsActive:             true,
						CreatedBy:            input.UserID,
						CreatedAt:            now,
					})
				}
			}
			allTaxDetails = append(allTaxDetails, taxDetails...)
		}
	}

	if len(allTaxDetails) > 0 {
		return invoiceRepo.CreateInvoiceTaxDetails(ctx, &allTaxDetails)
	}

	return nil
}

func (s *SalesInvoiceService) processInvoicePayments(
	ctx context.Context,
	invoiceRepo repository.SalesInvoiceRepository,
	draftRepo repository.SalesInvoiceDraftRepository,
	input CreateOrUpdateSalesInvoiceInput,
	invoiceID uint64,
	storeID uint64,
	draftID uint64,
) error {
	draftPayments, err := draftRepo.GetDraftPayments(ctx, draftID)
	if err != nil {
		return err
	}

	invoicePayments := make([]model.SalesInvoicePayment, 0, len(draftPayments))
	for _, p := range draftPayments {
		invoicePayments = append(invoicePayments, model.SalesInvoicePayment{
			SalesInvoiceID:       invoiceID,
			StoreID:              storeID,
			StorePaymentMethodID: p.StorePaymentMethodID,
			Amount:               p.Amount,
			VoucherCode:          p.VoucherCode,
			VoucherAmount:        p.VoucherAmount,
			Type:                 p.Type,
			CreatedBy:            p.CreatedBy,
			DeviceMasterID:       p.DeviceMasterID,
			TillID:               input.TillID,
			TillTransactionID:    input.TillTransactionID,
		})
	}
	return invoiceRepo.CreateInvoicePayments(ctx, invoicePayments)
}
