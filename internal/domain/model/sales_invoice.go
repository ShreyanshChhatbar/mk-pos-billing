package model

import (
	"time"

	"gorm.io/gorm"
)

type SalesInvoice struct {
	ID                          uint64         `gorm:"column:id;primaryKey"`
	InvoiceNumber               *string        `gorm:"column:invoice_number"`
	OrganizationID              uint64         `gorm:"column:organization_id"`
	SalesInvoiceDraftID         uint64         `gorm:"column:sales_invoice_draft_id"`
	StoreID                     uint64         `gorm:"column:store_id"`
	BillingUserID               uint64         `gorm:"column:billing_user_id"`
	CustomerID                  *uint64        `gorm:"column:customer_id"`
	CustomerAddressID           *uint64        `gorm:"column:customer_address_id"`
	DoctorID                    *uint64        `gorm:"column:doctor_id"`
	PatientID                   *uint64        `gorm:"column:patient_id"`
	OrderType                   string         `gorm:"column:order_type"`
	IsHomeDelivery              bool           `gorm:"column:is_home_delivery"`
	PaymentStatus               *string        `gorm:"column:payment_status"`
	TotalBillAmount             float64        `gorm:"column:total_bill_amount"`
	TaxableAmount               float64        `gorm:"column:taxable_amount"`
	TotalAmountBeforeDisc       float64        `gorm:"column:total_amount_before_discount"`
	DiscountType                string         `gorm:"column:discount_type"`
	DiscountPercentage          float64        `gorm:"column:discount_percentage"`
	DiscountAmount              float64        `gorm:"column:discount_amount"`
	IGST                        float64        `gorm:"column:igst"`
	CGST                        float64        `gorm:"column:cgst"`
	SGST                        float64        `gorm:"column:sgst"`
	PrepaidAmount               float64        `gorm:"column:prepaid_amount"`
	TotalGST                    float64        `gorm:"column:total_gst"`
	RoundOff                    float64        `gorm:"column:round_off"`
	TotalInvoiceAmount          float64        `gorm:"column:total_invoice_amount"`
	TotalProducts               int            `gorm:"column:total_products"`
	TotalItems                  int            `gorm:"column:total_items"`
	TotalQuantity               int            `gorm:"column:total_quantity"`
	TotalAmount                 float64        `gorm:"column:total_amount"`
	TotalDiscount               float64        `gorm:"column:total_discount"`
	TotalAmountReceived         float64        `gorm:"column:total_amount_received"`
	TotalBillBeforePromo        float64        `gorm:"column:total_bill_amount_before_promo"`
	PromoCode                   *string        `gorm:"column:promo_code"`
	Notes                       *string        `gorm:"column:notes"`
	IsActive                    bool           `gorm:"column:is_active"`
	VoucherDiscountAmount       float64        `gorm:"column:voucher_discount_amount"`
	LoyaltyPoints               float64        `gorm:"column:loyalty_points"` // numeric(10,1)
	LoyaltyProgramDiscount      int64          `gorm:"column:loyalty_program_discount"`
	IsLoyaltyUpdated            bool           `gorm:"column:is_loyalty_updated"`
	IsCustomerUpdated           bool           `gorm:"column:is_customer_updated"`
	IsSyncBill                  bool           `gorm:"column:is_sync_bill"`
	IsPrescriptionRequired      bool           `gorm:"column:is_prescription_required"`
	DeliveryCharges             float64        `gorm:"column:delivery_charges"`
	IsPrescriptionUploaded      bool           `gorm:"column:is_prescription_uploaded"`
	PrescriptionUploadedBy      *uint64        `gorm:"column:prescription_uploaded_by"`
	GSTTreatment                *string        `gorm:"column:gst_treatment"`
	GSTNumber                   *string        `gorm:"column:gst_number"`
	CINNumber                   *string        `gorm:"column:cin_number"`
	PlaceOfSupplyCode           *string        `gorm:"column:place_of_supply_code"`
	EditedUserID                *uint64        `gorm:"column:edited_user_id"`
	WhatsAppSentCount           int            `gorm:"column:whatsapp_sent_count"`
	IsRecommendationDataUpdated bool           `gorm:"column:is_recommendation_data_updated"`
	PrescriptionID              *uint64        `gorm:"column:prescription_id"`
	IsLPEnrolled                *bool          `gorm:"column:is_lp_enrolled"`
	CreatedBy                   uint64         `gorm:"column:created_by"`
	UpdatedBy                   *uint64        `gorm:"column:updated_by"`
	DeviceMasterID              *uint64        `gorm:"column:device_master_id"`
	CreatedAt                   time.Time      `gorm:"column:created_at"`
	UpdatedAt                   time.Time      `gorm:"column:updated_at"`
	DeletedAt                   gorm.DeletedAt `gorm:"column:deleted_at;index"`
	DeletedBy                   *uint64        `gorm:"column:deleted_by"`
}

func (SalesInvoice) TableName() string { return "sales_invoices" }
