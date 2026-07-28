package constants

const (
	SalesInvoiceStatusDraft          = "DRAFT"
	SalesInvoiceStatusPaymentPending = "PAYMENT_PENDING"
	SalesInvoiceStatusInvoiced       = "INVOICED"

	SalesPaymentStatusDraft         = "DRAFT"
	SalesPaymentStatusPartiallyPaid = "PARTIALLY_PAID"
	SalesPaymentStatusFullyPaid     = "FULLY_PAID"

	SalesPaymentTypeSales         = "SALES"
	SalesPaymentTypeAdvanceRefund = "ADVANCE_REFUND"
)

const (
	GstTreatmentRegular      = "Regular"
	GstTreatmentUnregistered = "Unregistered"
	GstTreatmentComposition  = "Composition"

	TaxNameIGST = "igst"
	TaxNameCGST = "cgst"
	TaxNameSGST = "sgst"

	TaxTypeAmount     = "AMT"
	TaxTypePercentage = "PCT"
)
