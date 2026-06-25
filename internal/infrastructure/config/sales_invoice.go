package config

import (
	"os"
	"strconv"
)

type SalesInvoiceConfig struct {
	DefaultOrganizationID         uint64
	DuplicateCheckCacheKey        string
	DuplicateRequestExpirySeconds int
	PrefixDraftBillCache          string
	PrefixDraftBillCacheTags      string
	DraftCacheTTLMinutes          int
	DeliveryChargeProductWSCode   int
	NarcoticsProductScheduledType string
	DefaultMinimumDaysForBatch    int
	BatchCodeLength               int
}

func LoadSalesInvoiceConfig() SalesInvoiceConfig {
	defaultOrg := uint64(1)
	if v := os.Getenv("DEFAULT_ORGANIZATION_ID"); v != "" {
		if i, err := strconv.ParseUint(v, 10, 64); err == nil {
			defaultOrg = i
		}
	}

	dupTTL := 10
	if v := os.Getenv("DUPLICATE_REQUEST_CHECK_EXPIRY_IN_SECONDS"); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			dupTTL = i
		}
	}

	draftTTL := 30
	if v := os.Getenv("SALES_INVOICE_DRAFT_CACHE_TTL_MINUTES"); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			draftTTL = i
		}
	}

	deliveryWS := 3234
	if v := os.Getenv("DELIVERY_CHARGE_PRODUCT_WS_CODE"); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			deliveryWS = i
		}
	}

	minimumDaysForBatch := 30
	if v := os.Getenv("DEFAULT_MINIMUM_DAYS_BATCH_EXPIRY"); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			minimumDaysForBatch = i
		}
	}

	batchCodeLength := 25
	if v := os.Getenv("BATCH_CODE_LENGTH"); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			batchCodeLength = i
		}
	}

	return SalesInvoiceConfig{
		DefaultOrganizationID:         defaultOrg,
		DuplicateCheckCacheKey:        envStringOrDefault("DUPLICATE_CHECK_CACHE_KEY", "duplicate_check_cache_key_"),
		DuplicateRequestExpirySeconds: dupTTL,
		PrefixDraftBillCache:          envStringOrDefault("PREFIX_DRAFT_BILL_CACHE", "draft_invoices_"),
		PrefixDraftBillCacheTags:      envStringOrDefault("PREFIX_DRAFT_BILL_CACHE_TAGS", "draft_invoices_tags_"),
		DraftCacheTTLMinutes:          draftTTL,
		DeliveryChargeProductWSCode:   deliveryWS,
		NarcoticsProductScheduledType: envStringOrDefault("NARCOTICS_PRODUCT_SCHEDULED_TYPE", "Schedule H1 - Narcotics"),
		DefaultMinimumDaysForBatch:    minimumDaysForBatch,
		BatchCodeLength:               batchCodeLength,
	}
}
