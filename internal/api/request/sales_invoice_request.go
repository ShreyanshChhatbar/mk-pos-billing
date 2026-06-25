package request

type CreateOrUpdateSalesInvoiceRequest struct {
	BillingUserID     uint64                       `json:"billing_user_id"`
	CustomerID        *uint64                      `json:"customer_id"`
	CustomerAddressID *uint64                      `json:"customer_address_id"`
	DoctorID          *uint64                      `json:"doctor_id"`
	PatientID         *uint64                      `json:"patient_id"`
	IsHomeDelivery    *bool                        `json:"is_home_delivery"`
	IsConfirmed       *bool                        `json:"is_confirmed"`
	PromoCode         *string                      `json:"promo_code"`
	CourseDays        *float64                     `json:"course_days"`
	Notes             *string                      `json:"notes"`
	ASMUserID         *uint64                      `json:"asm_user_id"`
	StoreID           uint64                       `json:"store_id"`
	UserID            uint64                       `json:"user_id"`
	DeviceMasterID    *uint64                      `json:"device_master_id"`
	TillID            *uint64                      `json:"till_id"`
	TillTransactionID *uint64                      `json:"till_transaction_id"`
	Items             *[]CreateOrUpdateInvoiceItem `json:"items"`
	Payments          []CreateOrUpdatePayment      `json:"payments"`
}

type CreateOrUpdateInvoiceItem struct {
	ProductID      uint64   `json:"product_id"`
	BatchCode      string   `json:"batch_code"`
	Quantity       int      `json:"quantity"`
	IsFreeProduct  bool     `json:"is_free_product"`
	ComboProductID *uint64  `json:"combo_product_id"`
	SalesRate      *float64 `json:"sales_rate"`
	PriceDelta     any      `json:"priceDelta"`
	BestAlternate  any      `json:"best_alternate"`
	PerTabFrontend any      `json:"per_tab_frontend"`
}

type CreateOrUpdatePayment struct {
	ID                   *uint64 `json:"id"`
	StorePaymentMethodID uint64  `json:"store_payment_method_id"`
	Amount               float64 `json:"amount"`
	VoucherCode          *string `json:"voucher_code"`
	IsAdvanceRefund      bool    `json:"is_advance_refund"`
}

type UpdateSalesInvoiceURI struct {
	ID uint64 `uri:"id"`
}
