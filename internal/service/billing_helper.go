package service

import (
	"fmt"
	"strconv"
	"strings"

	"mk-pos-billing/internal/constants"
	domaincache "mk-pos-billing/internal/domain/cache"
	"mk-pos-billing/internal/domain/model"
)

// TaxCalculationResult holds the calculated GST breakdown for a single line item.
type TaxCalculationResult struct {
	BaseRate float64
	TotalGST float64
	GSTPct   float64
	CGST     float64
	SGST     float64
	IGST     float64
}

// CalculateSalesInvoiceItemTax is the Go equivalent of CommonHelper::calculateSalesInvoiceItemTax.
// It calculates the GST breakdown for a single billing line using a tax-inclusive reverse calculation.
// salesRate is the total line amount (quantity × unit rate), already rounded.
func CalculateSalesInvoiceItemTax(
	store *domaincache.Store,
	product *domaincache.Product,
	salesRate float64,
) (TaxCalculationResult, error) {
	var gstPct float64

	switch store.GstTreatment {
	case constants.GstTreatmentRegular, constants.GstTreatmentUnregistered:
		rawGST := strings.TrimRight(product.GstType, "%")
		parsed, err := strconv.ParseFloat(rawGST, 64)
		if err != nil {
			return TaxCalculationResult{}, fmt.Errorf("invalid gst_type %q for product %d: %w", rawGST, product.ID, err)
		}
		gstPct = parsed

	case constants.GstTreatmentComposition:
		gstPct = 0

	default:
		return TaxCalculationResult{}, fmt.Errorf("invalid gst treatment of store: %q", store.GstTreatment)
	}

	// Tax-inclusive reverse calculation — mirrors Laravel exactly:
	//   baseRate = (salesRate * 100) / (100 + gst)
	baseRate := (salesRate * 100) / (100 + gstPct)
	gstAmount := salesRate - baseRate

	result := TaxCalculationResult{
		BaseRate: baseRate,
		TotalGST: gstAmount,
		GSTPct:   gstPct,
	}

	// Intra-state → split into CGST + SGST; inter-state → IGST
	// interstate is used for wholesale billing only so not needed
	insideState := true // always same supply code for now (bill_from == bill_to)
	if insideState {
		result.CGST = gstAmount / 2
		result.SGST = gstAmount / 2
	} else {
		result.IGST = gstAmount
	}

	return result, nil
}

// SalesRateResult holds the calculated sales rate and discount.
type SalesRateResult struct {
	SalesRate            float64
	DiscountPercentage   float64
	DiscountAmount       float64
	DiscountType         string
	SalesRateBeforePromo float64
}

// GetSalesRate mirrors CommonHelper::getSalesRateCache.
func GetSalesRate(
	organizationID uint64,
	product *domaincache.Product,
	store *domaincache.Store,
	mrp float64,
	defaultOrgID uint64,
	genericPricingMap map[string]model.B2CStoreTemplateGenericPricing,
) (SalesRateResult, error) {

	// In Laravel, Combo Products are checked first (omitted here as per Laravel logic if isComboProduct=false).

	if organizationID == defaultOrgID {
		category := product.B2CProductCategoryID
		discountData, ok := store.ProductCategoriesDiscounts[category]
		if !ok {
			return SalesRateResult{}, fmt.Errorf("no pricing template assigned for category %d", category)
		}

		if product.IsGeneric {
			templateID := discountData.B2CPricingTemplateID
			key := fmt.Sprintf("%d:%d", templateID, product.ID)

			if gp, exists := genericPricingMap[key]; exists {
				return prepareSalesRateData(gp.SalesPrice, mrp), nil
			}
		}

		return CalculateSalesRate(discountData.Mode, discountData.Operator, discountData.PricingCategory, discountData.Value, mrp), nil
	} else {
		// Non-default org B2B Logic (mirrors Laravel)
		category := product.OrganizationCategoryID
		orgDiscounts, ok := store.OrganizationProductCategoriesDiscounts[int(organizationID)]
		if !ok {
			return SalesRateResult{}, fmt.Errorf("no pricing category found for org %d", organizationID)
		}

		discountData, ok := orgDiscounts[category]
		if !ok {
			return SalesRateResult{}, fmt.Errorf("no pricing template assigned for org category %d", category)
		}

		// Note: Laravel code for org generic pricing is commented out in original PHP.
		valFloat, _ := strconv.ParseFloat(discountData.Value, 64)
		return CalculateSalesRate(discountData.Mode, discountData.Operator, discountData.PricingCategory, valFloat, mrp), nil
	}
}

func prepareSalesRateData(salesPrice float64, mrp float64) SalesRateResult {
	discountPercentage := 0.0
	if mrp > 0 {
		discountPercentage = 100 - ((salesPrice * 100) / mrp)
	}

	return SalesRateResult{
		SalesRate:            salesPrice,
		DiscountPercentage:   discountPercentage,
		DiscountAmount:       mrp - salesPrice,
		DiscountType:         "INR", // "INR" is explicitly returned in Laravel's prepareSalesRateData
		SalesRateBeforePromo: salesPrice,
	}
}

// CalculateSalesRate mirrors CommonHelper::calculateSalesRate
func CalculateSalesRate(mode string, operator string, pricingCategory string, value float64, mrp float64) SalesRateResult {
	salesRate := mrp
	discountPercentage := 0.0

	if mode == "INR" {
		if operator == "-" && pricingCategory == "mrp" { // using lower-cased 'mrp' as in Laravel strtolower
			salesRate = mrp - value
			if mrp > 0 {
				discountPercentage = (value * 100) / mrp
			}
		}
	} else if mode == "%" {
		if operator == "-" && pricingCategory == "mrp" {
			salesRate = mrp - (mrp * (value / 100))
			discountPercentage = value
		}
	}

	return SalesRateResult{
		SalesRate:            salesRate,
		DiscountPercentage:   discountPercentage,
		DiscountAmount:       mrp - salesRate,
		DiscountType:         mode,
		SalesRateBeforePromo: salesRate,
	}
}
