package request

type CreateOrUpdateSalesInvoiceRequest struct {
	BillingUserID     uint64                       `json:"billing_user_id" binding:"required,gt=0"`
	CustomerID        *uint64                      `json:"customer_id" binding:"omitempty,gt=0"`
	CustomerAddressID *uint64                      `json:"customer_address_id" binding:"omitempty,gt=0"`
	DoctorID          *uint64                      `json:"doctor_id" binding:"omitempty,gt=0"`
	PatientID         *uint64                      `json:"patient_id" binding:"omitempty,gt=0"`
	IsHomeDelivery    *bool                        `json:"is_home_delivery"`
	IsConfirmed       *bool                        `json:"is_confirmed" binding:"required"`
	PromoCode         *string                      `json:"promo_code"`
	CourseDays        *float64                     `json:"course_days" binding:"omitempty,gt=0"`
	Notes             *string                      `json:"notes"`
	ASMUserID         *uint64                      `json:"asm_user_id" binding:"omitempty,gt=0"`
	StoreID           uint64                       `json:"store_id" binding:"required,gt=0"`
	UserID            uint64                       `json:"user_id" binding:"required,gt=0"`
	DeviceMasterID    *uint64                      `json:"device_master_id" binding:"omitempty,gt=0"`
	TillID            *uint64                      `json:"till_id" binding:"omitempty,gt=0"`
	TillTransactionID *uint64                      `json:"till_transaction_id" binding:"omitempty,gt=0"`
	Items             []CreateOrUpdateInvoiceItem  `json:"items"`
	Payments          []CreateOrUpdatePayment      `json:"payments"`
}

type CreateOrUpdateInvoiceItem struct {
	ProductID      uint64   `json:"product_id" binding:"required,gt=0"`
	BatchCode      string   `json:"batch_code" binding:"required,max=32"` // use cfg value if fixed, otherwise service
	Quantity       int      `json:"quantity" binding:"required,gt=0"`
	IsFreeProduct  bool     `json:"is_free_product"`
	ComboProductID *uint64  `json:"combo_product_id" binding:"omitempty,gt=0"`
	SalesRate      *float64 `json:"sales_rate"`
	PriceDelta     any      `json:"priceDelta"`
	BestAlternate  any      `json:"best_alternate"`
	PerTabFrontend any      `json:"per_tab_frontend"`
}

type CreateOrUpdatePayment struct {
	ID                   *uint64 `json:"id" binding:"omitempty,gt=0"`
	StorePaymentMethodID uint64  `json:"store_payment_method_id" binding:"required,gt=0"`
	Amount               float64 `json:"amount" binding:"required,gt=0"`
	VoucherCode          *string `json:"voucher_code"`
	IsAdvanceRefund      bool    `json:"is_advance_refund"`
}

type UpdateSalesInvoiceURI struct {
	ID uint64 `uri:"id" binding:"required"`
}
