package model

import (
	"time"

	"gorm.io/gorm"
)

type SalesInvoiceDraftJSON struct {
	ID                    uint64         `gorm:"column:id;primaryKey"`
	StoreID               uint64         `gorm:"column:store_id"`
	OrganizationID        uint64         `gorm:"column:organization_id"`
	BillingUserID         uint64         `gorm:"column:billing_user_id"`
	CustomerID            *uint64        `gorm:"column:customer_id"`
	CustomerAddressID     *uint64        `gorm:"column:customer_address_id"`
	DoctorID              *uint64        `gorm:"column:doctor_id"`
	PatientID             *uint64        `gorm:"column:patient_id"`
	PrescriptionID        *uint64        `gorm:"column:prescription_id"`
	IsHomeDelivery        bool           `gorm:"column:is_home_delivery"`
	Status                string         `gorm:"column:status"`
	PaymentStatus         string         `gorm:"column:payment_status"`
	TotalBillAmount       float64        `gorm:"column:total_bill_amount"`
	TaxableAmount         float64        `gorm:"column:taxable_amount"`
	TotalAmountBeforeDisc float64        `gorm:"column:total_amount_before_discount"`
	DiscountType          string         `gorm:"column:discount_type"`
	DiscountPercentage    float64        `gorm:"column:discount_percentage"`
	DiscountAmount        float64        `gorm:"column:discount_amount"`
	IGST                  float64        `gorm:"column:igst"`
	CGST                  float64        `gorm:"column:cgst"`
	SGST                  float64        `gorm:"column:sgst"`
	PrepaidAmount         float64        `gorm:"column:prepaid_amount"`
	TotalGST              float64        `gorm:"column:total_gst"`
	RoundOff              float64        `gorm:"column:round_off"`
	TotalInvoiceAmount    float64        `gorm:"column:total_invoice_amount"`
	TotalProducts         int            `gorm:"column:total_products"`
	TotalItems            int            `gorm:"column:total_items"`
	TotalQuantity         int            `gorm:"column:total_quantity"`
	TotalAmount           float64        `gorm:"column:total_amount"`
	TotalDiscount         float64        `gorm:"column:total_discount"`
	TotalAmountReceived   float64        `gorm:"column:total_amount_received"`
	TotalBillBeforePromo  float64        `gorm:"column:total_bill_amount_before_promo"`
	PromoCode             *string        `gorm:"column:promo_code"`
	Notes                 *string        `gorm:"column:notes"`
	DraftJSON             []byte         `gorm:"column:draft_json"`
	DeviceMasterID        *uint64        `gorm:"column:device_master_id"`
	TillID                *uint64        `gorm:"column:till_id"`
	TillTransactionID     *uint64        `gorm:"column:till_transaction_id"`
	IsActive              bool           `gorm:"column:is_active"`
	CreatedBy             uint64         `gorm:"column:created_by"`
	UpdatedBy             *uint64        `gorm:"column:updated_by"`
	CreatedAt             time.Time      `gorm:"column:created_at"`
	UpdatedAt             time.Time      `gorm:"column:updated_at"`
	DeletedAt             gorm.DeletedAt `gorm:"column:deleted_at;index"`
	DeletedBy             *uint64        `gorm:"column:deleted_by"`
}

func (SalesInvoiceDraftJSON) TableName() string { return "sales_invoice_draft_jsons" }

type SalesInvoiceDraftPayment struct {
	ID                   uint64         `gorm:"column:id;primaryKey"`
	SalesInvoiceDraftID  uint64         `gorm:"column:sales_invoice_draft_id"`
	StoreID              uint64         `gorm:"column:store_id"`
	StorePaymentMethodID uint64         `gorm:"column:store_payment_method_id"`
	Amount               float64        `gorm:"column:amount"`
	VoucherCode          *string        `gorm:"column:voucher_code"`
	VoucherAmount        *float64       `gorm:"column:voucher_amount"`
	Type                 string         `gorm:"column:type"`
	DeviceMasterID       *uint64        `gorm:"column:device_master_id"`
	CreatedBy            uint64         `gorm:"column:created_by"`
	UpdatedBy            *uint64        `gorm:"column:updated_by"`
	CreatedAt            time.Time      `gorm:"column:created_at"`
	UpdatedAt            time.Time      `gorm:"column:updated_at"`
	DeletedAt            gorm.DeletedAt `gorm:"column:deleted_at;index"`
	DeletedBy            *uint64        `gorm:"column:deleted_by"`
}

func (SalesInvoiceDraftPayment) TableName() string { return "sales_invoice_draft_payments" }

type SalesInvoice struct {
	ID                    uint64         `gorm:"column:id;primaryKey"`
	OrganizationID        uint64         `gorm:"column:organization_id"`
	SalesInvoiceDraftID   uint64         `gorm:"column:sales_invoice_draft_id"`
	StoreID               uint64         `gorm:"column:store_id"`
	BillingUserID         uint64         `gorm:"column:billing_user_id"`
	CustomerID            *uint64        `gorm:"column:customer_id"`
	CustomerAddressID     *uint64        `gorm:"column:customer_address_id"`
	DoctorID              *uint64        `gorm:"column:doctor_id"`
	PatientID             *uint64        `gorm:"column:patient_id"`
	OrderType             string         `gorm:"column:order_type"`
	IsHomeDelivery        bool           `gorm:"column:is_home_delivery"`
	TotalBillAmount       float64        `gorm:"column:total_bill_amount"`
	TaxableAmount         float64        `gorm:"column:taxable_amount"`
	TotalAmountBeforeDisc float64        `gorm:"column:total_amount_before_discount"`
	DiscountType          string         `gorm:"column:discount_type"`
	DiscountPercentage    float64        `gorm:"column:discount_percentage"`
	DiscountAmount        float64        `gorm:"column:discount_amount"`
	IGST                  float64        `gorm:"column:igst"`
	CGST                  float64        `gorm:"column:cgst"`
	SGST                  float64        `gorm:"column:sgst"`
	PrepaidAmount         float64        `gorm:"column:prepaid_amount"`
	TotalGST              float64        `gorm:"column:total_gst"`
	RoundOff              float64        `gorm:"column:round_off"`
	TotalInvoiceAmount    float64        `gorm:"column:total_invoice_amount"`
	TotalProducts         int            `gorm:"column:total_products"`
	TotalItems            int            `gorm:"column:total_items"`
	TotalQuantity         int            `gorm:"column:total_quantity"`
	TotalAmount           float64        `gorm:"column:total_amount"`
	TotalDiscount         float64        `gorm:"column:total_discount"`
	TotalAmountReceived   float64        `gorm:"column:total_amount_received"`
	TotalBillBeforePromo  float64        `gorm:"column:total_bill_amount_before_promo"`
	PromoCode             *string        `gorm:"column:promo_code"`
	Notes                 *string        `gorm:"column:notes"`
	IsActive              bool           `gorm:"column:is_active"`
	CreatedBy             uint64         `gorm:"column:created_by"`
	UpdatedBy             *uint64        `gorm:"column:updated_by"`
	DeviceMasterID        *uint64        `gorm:"column:device_master_id"`
	TillID                *uint64        `gorm:"column:till_id"`
	TillTransactionID     *uint64        `gorm:"column:till_transaction_id"`
	CreatedAt             time.Time      `gorm:"column:created_at"`
	UpdatedAt             time.Time      `gorm:"column:updated_at"`
	DeletedAt             gorm.DeletedAt `gorm:"column:deleted_at;index"`
	DeletedBy             *uint64        `gorm:"column:deleted_by"`
}

func (SalesInvoice) TableName() string { return "sales_invoices" }

type SalesInvoiceDetail struct {
	ID                 uint64         `gorm:"column:id;primaryKey"`
	SalesInvoiceID     uint64         `gorm:"column:sales_invoice_id"`
	StoreID            uint64         `gorm:"column:store_id"`
	ProductID          uint64         `gorm:"column:product_id"`
	StoreBatchID       *uint64        `gorm:"column:store_batch_id"`
	PurchaseRate       float64        `gorm:"column:purchase_rate"`
	BatchCode          string         `gorm:"column:batch_code"`
	ExpiryDate         time.Time      `gorm:"column:expiry_date"`
	MRP                float64        `gorm:"column:mrp"`
	SalesRate          float64        `gorm:"column:sales_rate"`
	BaseRate           float64        `gorm:"column:base_rate"`
	BillAmount         float64        `gorm:"column:bill_amount"`
	Quantity           int            `gorm:"column:quantity"`
	DiscountType       string         `gorm:"column:discount_type"`
	DiscountPercentage float64        `gorm:"column:discount_percentage"`
	DiscountAmount     float64        `gorm:"column:discount_amount"`
	GSTPercentage      float64        `gorm:"column:gst_percentage"`
	GSTAmount          float64        `gorm:"column:gst_amount"`
	TotalAmount        float64        `gorm:"column:total_amount"`
	CreatedBy          uint64         `gorm:"column:created_by"`
	DeviceMasterID     *uint64        `gorm:"column:device_master_id"`
	CreatedAt          time.Time      `gorm:"column:created_at"`
	UpdatedAt          time.Time      `gorm:"column:updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (SalesInvoiceDetail) TableName() string { return "sales_invoice_details" }

type SalesInvoicePayment struct {
	ID                   uint64         `gorm:"column:id;primaryKey"`
	SalesInvoiceID       uint64         `gorm:"column:sales_invoice_id"`
	StoreID              uint64         `gorm:"column:store_id"`
	StorePaymentMethodID uint64         `gorm:"column:store_payment_method_id"`
	Amount               float64        `gorm:"column:amount"`
	VoucherCode          *string        `gorm:"column:voucher_code"`
	VoucherAmount        *float64       `gorm:"column:voucher_amount"`
	VoucherType          *string        `gorm:"column:voucher_type"`
	Type                 string         `gorm:"column:type"`
	IsActive             bool           `gorm:"column:is_active"`
	CreatedBy            uint64         `gorm:"column:created_by"`
	UpdatedBy            *uint64        `gorm:"column:updated_by"`
	DeviceMasterID       *uint64        `gorm:"column:device_master_id"`
	TillID               *uint64        `gorm:"column:till_id"`
	TillTransactionID    *uint64        `gorm:"column:till_transaction_id"`
	CreatedAt            time.Time      `gorm:"column:created_at"`
	UpdatedAt            time.Time      `gorm:"column:updated_at"`
	DeletedAt            gorm.DeletedAt `gorm:"column:deleted_at;index"`
	DeletedBy            *uint64        `gorm:"column:deleted_by"`
}

func (SalesInvoicePayment) TableName() string { return "sales_invoice_payments" }

type StoreInventory struct {
	ID           uint64 `gorm:"column:id"`
	StoreID      uint64 `gorm:"column:store_id"`
	ProductID    uint64 `gorm:"column:product_id"`
	BatchCode    string `gorm:"column:batch_code"`
	StoreBatchID uint64 `gorm:"column:store_batch_id"`
	ClosingStock int    `gorm:"column:closing_stock"`
}

func (StoreInventory) TableName() string { return "store_inventories" }

type StoreBatch struct {
	ID           uint64  `gorm:"column:id;primaryKey"`
	StoreID      uint64  `gorm:"column:store_id"`
	BatchID      uint64  `gorm:"column:batch_id"`
	PurchaseRate float64 `gorm:"column:purchase_rate"`
}

func (StoreBatch) TableName() string { return "store_batches" }

type Batch struct {
	ID          uint64    `gorm:"column:id;primaryKey"`
	ProductID   uint64    `gorm:"column:product_id"`
	BatchNumber string    `gorm:"column:batch_number"`
	MRP         float64   `gorm:"column:mrp"`
	OldMRP      *float64  `gorm:"column:old_mrp"`
	ExpiryDate  time.Time `gorm:"column:expiry_date"`
}

func (Batch) TableName() string { return "batches" }
