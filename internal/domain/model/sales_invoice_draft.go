package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// DraftTaxDetailJSON represents a single GST component (IGST, CGST, or SGST)
// stored inside the draft_json column's products[].tax_details array.
type DraftTaxDetailJSON struct {
	TaxName   string  `json:"tax_name"   php:"tax_name"`
	TaxRate   float64 `json:"tax_rate"   php:"tax_rate"`
	TaxType   string  `json:"tax_type"   php:"tax_type"`
	TaxAmount float64 `json:"tax_amount" php:"tax_amount"`
}

// DraftProductJSON is a single product line as stored inside the draft_json column.
type DraftProductJSON struct {
	ProductID                     uint64               `json:"product_id"                       php:"product_id"`
	BatchCode                     string               `json:"batch_code"                       php:"batch_code"`
	ExpiryDate                    string               `json:"expiry_date"                      php:"expiry_date"`
	MRP                           float64              `json:"mrp"                              php:"mrp"`
	SalesRate                     float64              `json:"sales_rate"                       php:"sales_rate"`
	BaseRate                      float64              `json:"base_rate"                        php:"base_rate"`
	Quantity                      int                  `json:"quantity"                         php:"quantity"`
	GSTPercentage                 float64              `json:"gst_percentage"                   php:"gst_percentage"`
	GSTAmount                     float64              `json:"gst_amount"                       php:"gst_amount"`
	BillAmount                    float64              `json:"bill_amount"                      php:"bill_amount"`
	DiscountType                  string               `json:"discount_type"                    php:"discount_type"`
	DiscountAmount                float64              `json:"discount_amount"                  php:"discount_amount"`
	DiscountPercentage            float64              `json:"discount_percentage"              php:"discount_percentage"`
	IsFreeProduct                 bool                 `json:"is_free_product"                  php:"is_free_product"`
	ComboProductID                *uint64              `json:"combo_product_id"                 php:"combo_product_id"`
	HSNCode                       string               `json:"hsn_code"                         php:"hsn_code"`
	CreatedBy                     uint64               `json:"created_by"                       php:"created_by"`
	TaxDetails                    []DraftTaxDetailJSON `json:"tax_details"                      php:"tax_details"`
	TotalAmount                   float64              `json:"total_amount"                     php:"total_amount"`
	BatchQuantity                 int                  `json:"batch_quantity"                   php:"batch_quantity"`
	DeviceMasterID                *uint64              `json:"device_master_id"                 php:"device_master_id"`
	IsAdvanceOrder                bool                 `json:"is_advance_order"                 php:"is_advance_order"`
	IsComboProduct                bool                 `json:"is_combo_product"                 php:"is_combo_product"`
	OrderedQuantity               *int                 `json:"ordered_quantity"                 php:"ordered_quantity"`
	ProductLocation               string               `json:"product_location"                 php:"product_location"`
	AvailableQuantity             int                  `json:"available_quantity"               php:"available_quantity"`
	SalesRateBeforePromo          float64              `json:"sales_rate_before_promo"          php:"sales_rate_before_promo"`
	IsPrescriptionRequired        bool                 `json:"is_prescription_required"         php:"is_prescription_required"`
	DiscountAmountBeforePromo     float64              `json:"discount_amount_before_promo"     php:"discount_amount_before_promo"`
	DiscountPercentageBeforePromo float64              `json:"discount_percentage_before_promo" php:"discount_percentage_before_promo"`
	IsEditable                    bool                 `json:"is_editable"                      php:"is_editable"`
}

// DraftJSONPayload is the structured representation of the draft_json column
// in sales_invoice_draft_jsons. It implements driver.Valuer and sql.Scanner so
// GORM can read and write it automatically without manual json.Marshal/Unmarshal
// calls in the service layer.
type DraftJSONPayload struct {
	IsHomeDelivery             bool               `json:"is_home_delivery"              php:"is_home_delivery"`
	IsActive                   bool               `json:"is_active"                     php:"is_active"`
	PaymentStatus              string             `json:"payment_status"                php:"payment_status"`
	DeviceMasterID             *uint64            `json:"device_master_id"              php:"device_master_id"`
	PromoCode                  *string            `json:"promo_code"                    php:"promo_code"`
	Notes                      *string            `json:"notes"                         php:"notes"`
	TotalProducts              int                `json:"total_products"                php:"total_products"`
	TotalItems                 int                `json:"total_items"                   php:"total_items"`
	TotalQuantity              int                `json:"total_quantity"                php:"total_quantity"`
	TotalGST                   float64            `json:"total_gst"                     php:"total_gst"`
	SGST                       float64            `json:"sgst"                          php:"sgst"`
	CGST                       float64            `json:"cgst"                          php:"cgst"`
	IGST                       float64            `json:"igst"                          php:"igst"`
	DeliveryCharges            float64            `json:"delivery_charges"              php:"delivery_charges"`
	TaxableAmount              float64            `json:"taxable_amount"                php:"taxable_amount"`
	RoundOff                   float64            `json:"round_off"                     php:"round_off"`
	CINNumber                  *string            `json:"cin_number"                    php:"cin_number"`
	GSTNumber                  *string            `json:"gst_number"                    php:"gst_number"`
	GSTTreatment               string             `json:"gst_treatment"                 php:"gst_treatment"`
	TotalAmount                float64            `json:"total_amount"                  php:"total_amount"`
	TotalBillAmount            float64            `json:"total_bill_amount"             php:"total_bill_amount"`
	TotalInvoiceAmount         float64            `json:"total_invoice_amount"          php:"total_invoice_amount"`
	TotalBillAmountBeforePromo float64            `json:"total_bill_amount_before_promo" php:"total_bill_amount_before_promo"`
	PlaceOfSupplyCode          string             `json:"place_of_supply_code"          php:"place_of_supply_code"`
	IsPrescriptionRequired     bool               `json:"is_prescription_required"      php:"is_prescription_required"`
	Products                   []DraftProductJSON `json:"products"                      php:"products"`
}

// Value implements driver.Valuer — serializes to JSON for DB writes.
func (p DraftJSONPayload) Value() (driver.Value, error) {
	b, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	return string(b), nil
}

// Scan implements sql.Scanner — deserializes from DB JSON column on reads.
func (p *DraftJSONPayload) Scan(src any) error {
	var b []byte
	switch v := src.(type) {
	case []byte:
		b = v
	case string:
		b = []byte(v)
	case nil:
		return nil
	default:
		return fmt.Errorf("DraftJSONPayload: unsupported type %T", src)
	}
	return json.Unmarshal(b, p)
}

// SalesInvoiceDraftJSON is the GORM model for the sales_invoice_draft_jsons table.
type SalesInvoiceDraftJSON struct {
	ID                     uint64           `gorm:"column:id;primaryKey"               php:"id"`
	InvoiceNumber          string           `gorm:"column:invoice_number"              php:"invoice_number"`
	StoreID                uint64           `gorm:"column:store_id"                    php:"store_id"`
	OrganizationID         uint64           `gorm:"column:organization_id"             php:"organization_id"`
	BillingUserID          uint64           `gorm:"column:billing_user_id"             php:"billing_user_id"`
	CustomerID             *uint64          `gorm:"column:customer_id"                 php:"customer_id"`
	CustomerAddressID      *uint64          `gorm:"column:customer_address_id"         php:"customer_address_id"`
	DoctorID               *uint64          `gorm:"column:doctor_id"                   php:"doctor_id"`
	PatientID              *uint64          `gorm:"column:patient_id"                  php:"patient_id"`
	PrescriptionID         *uint64          `gorm:"column:prescription_id"             php:"prescription_id"`
	Status                 string           `gorm:"column:status"                      php:"status"`
	PaymentStatus          string           `gorm:"column:payment_status"              php:"payment_status"`
	TotalBillAmount        float64          `gorm:"column:total_bill_amount"           php:"total_bill_amount"`
	PrepaidAmount          float64          `gorm:"column:prepaid_amount"              php:"prepaid_amount"`
	RoundOff               float64          `gorm:"column:round_off"                   php:"round_off"`
	TotalInvoiceAmount     float64          `gorm:"column:total_invoice_amount"        php:"total_invoice_amount"`
	TotalAmount            float64          `gorm:"column:total_amount"                php:"total_amount"`
	TotalDiscount          float64          `gorm:"column:total_discount"              php:"total_discount"`
	TotalAmountReceived    float64          `gorm:"column:total_amount_received"       php:"total_amount_received"`
	LoyaltyProgramDiscount int64            `gorm:"column:loyalty_program_discount"    php:"loyalty_program_discount"`
	CancelReason           *string          `gorm:"column:cancel_reason"               php:"cancel_reason"`
	DraftJSON              DraftJSONPayload `gorm:"column:draft_json"                  php:"draft_json"`
	TillID                 *uint64          `gorm:"column:till_id"                     php:"till_id"`
	TillTransactionID      *uint64          `gorm:"column:till_transaction_id"         php:"till_transaction_id"`
	CreatedBy              uint64           `gorm:"column:created_by"                  php:"created_by"`
	UpdatedBy              *uint64          `gorm:"column:updated_by"                  php:"updated_by"`
	CreatedAt              time.Time        `gorm:"column:created_at"                  php:"created_at"`
	UpdatedAt              time.Time        `gorm:"column:updated_at"                  php:"updated_at"`
	DeletedAt              gorm.DeletedAt   `gorm:"column:deleted_at;index"            php:"deleted_at"`
	DeletedBy              *uint64          `gorm:"column:deleted_by"                  php:"deleted_by"`
}

func (SalesInvoiceDraftJSON) TableName() string { return "sales_invoice_draft_jsons" }
