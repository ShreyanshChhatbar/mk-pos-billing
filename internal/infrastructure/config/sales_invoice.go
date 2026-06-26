package config

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
	return SalesInvoiceConfig{
		DefaultOrganizationID: envUint64OrDefault("DEFAULT_ORGANIZATION_ID", 1),

		DuplicateCheckCacheKey:        envStringOrDefault("DUPLICATE_CHECK_CACHE_KEY", "duplicate_check_cache_key_"),
		DuplicateRequestExpirySeconds: envIntOrDefault("DUPLICATE_REQUEST_CHECK_EXPIRY_IN_SECONDS", 10),

		PrefixDraftBillCache:     envStringOrDefault("PREFIX_DRAFT_BILL_CACHE", "draft_invoices_"),
		PrefixDraftBillCacheTags: envStringOrDefault("PREFIX_DRAFT_BILL_CACHE_TAGS", "draft_invoices_tags_"),
		DraftCacheTTLMinutes:     envIntOrDefault("SALES_INVOICE_DRAFT_CACHE_TTL_MINUTES", 30),

		DeliveryChargeProductWSCode: envIntOrDefault("DELIVERY_CHARGE_PRODUCT_WS_CODE", 3234),

		NarcoticsProductScheduledType: envStringOrDefault(
			"NARCOTICS_PRODUCT_SCHEDULED_TYPE",
			"Schedule H1 - Narcotics",
		),

		DefaultMinimumDaysForBatch: envIntOrDefault("DEFAULT_MINIMUM_DAYS_BATCH_EXPIRY", 30),
		BatchCodeLength:            envIntOrDefault("BATCH_CODE_LENGTH", 25),
	}
}
