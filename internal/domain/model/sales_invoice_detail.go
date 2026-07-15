package model

import (
	"time"

	"gorm.io/gorm"
)

type SalesInvoiceDetail struct {
	ID                     uint64         `gorm:"column:id;primaryKey"`
	SalesInvoiceID         uint64         `gorm:"column:sales_invoice_id"`
	StoreID                uint64         `gorm:"column:store_id"`
	ProductID              uint64         `gorm:"column:product_id"`
	StoreBatchID           *uint64        `gorm:"column:store_batch_id"`
	PurchaseRate           float64        `gorm:"column:purchase_rate"`
	BatchCode              string         `gorm:"column:batch_code"`
	ExpiryDate             time.Time      `gorm:"column:expiry_date"`
	MRP                    float64        `gorm:"column:mrp"`
	SalesRate              float64        `gorm:"column:sales_rate"`
	BaseRate               float64        `gorm:"column:base_rate"`
	BillAmount             float64        `gorm:"column:bill_amount"`
	Quantity               int            `gorm:"column:quantity"`
	DiscountType           string         `gorm:"column:discount_type"`
	DiscountPercentage     float64        `gorm:"column:discount_percentage"`
	DiscountAmount         float64        `gorm:"column:discount_amount"`
	GSTPercentage          float64        `gorm:"column:gst_percentage"`
	GSTAmount              float64        `gorm:"column:gst_amount"`
	TotalAmount            float64        `gorm:"column:total_amount"`
	EffectiveQuantity      *int           `gorm:"column:effective_quantity"`
	VoucherDiscountAmount  *float64       `gorm:"column:voucher_discount_amount"`
	AmountBeforeDiscount   *float64       `gorm:"column:amount_before_discount"`
	LoyaltyPoints          *float64       `gorm:"column:loyalty_points"`
	LoyaltyProgramDiscount *float64       `gorm:"column:loyalty_program_discount"`
	SalesRateBeforePromo   *float64       `gorm:"column:sales_rate_before_promo"`
	IsAdvanceOrder         *bool          `gorm:"column:is_advance_order"`
	TillTransactionID      *uint64        `gorm:"column:till_transaction_id"`
	TillID                 *uint64        `gorm:"column:till_id"`
	IsDropshipOrder        *bool          `gorm:"column:is_dropship_order"`
	HSNCode                *string        `gorm:"column:hsn_code"`
	IsFreeProduct          *bool          `gorm:"column:is_free_product"`
	DosageForm             *string        `gorm:"column:dosage_form"`
	MisReportingCategoryID *uint64        `gorm:"column:mis_reporting_category_id"`
	CreatedBy              uint64         `gorm:"column:created_by"`
	UpdatedBy              *uint64        `gorm:"column:updated_by"`
	DeviceMasterID         *uint64        `gorm:"column:device_master_id"`
	CreatedAt              time.Time      `gorm:"column:created_at"`
	UpdatedAt              time.Time      `gorm:"column:updated_at"`
	DeletedAt              gorm.DeletedAt `gorm:"column:deleted_at;index"`
	DeletedBy              *uint64        `gorm:"column:deleted_by"`
}

func (SalesInvoiceDetail) TableName() string { return "sales_invoice_details" }
