package model

import (
	"time"

	"gorm.io/gorm"
)

type SalesInvoiceDraftJSON struct {
	ID                     uint64         `gorm:"column:id;primaryKey"`
	InvoiceNumber          string         `gorm:"column:invoice_number"`
	StoreID                uint64         `gorm:"column:store_id"`
	OrganizationID         uint64         `gorm:"column:organization_id"`
	BillingUserID          uint64         `gorm:"column:billing_user_id"`
	CustomerID             *uint64        `gorm:"column:customer_id"`
	CustomerAddressID      *uint64        `gorm:"column:customer_address_id"`
	DoctorID               *uint64        `gorm:"column:doctor_id"`
	PatientID              *uint64        `gorm:"column:patient_id"`
	PrescriptionID         *uint64        `gorm:"column:prescription_id"`
	Status                 string         `gorm:"column:status"`
	PaymentStatus          string         `gorm:"column:payment_status"`
	TotalBillAmount        float64        `gorm:"column:total_bill_amount"`
	PrepaidAmount          float64        `gorm:"column:prepaid_amount"`
	RoundOff               float64        `gorm:"column:round_off"`
	TotalInvoiceAmount     float64        `gorm:"column:total_invoice_amount"`
	TotalAmount            float64        `gorm:"column:total_amount"`
	TotalDiscount          float64        `gorm:"column:total_discount"`
	TotalAmountReceived    float64        `gorm:"column:total_amount_received"`
	LoyaltyProgramDiscount int64          `gorm:"column:loyalty_program_discount"`
	CancelReason           *string        `gorm:"column:cancel_reason"`
	DraftJSON              []byte         `gorm:"column:draft_json"`
	TillID                 *uint64        `gorm:"column:till_id"`
	TillTransactionID      *uint64        `gorm:"column:till_transaction_id"`
	CreatedBy              uint64         `gorm:"column:created_by"`
	UpdatedBy              *uint64        `gorm:"column:updated_by"`
	CreatedAt              time.Time      `gorm:"column:created_at"`
	UpdatedAt              time.Time      `gorm:"column:updated_at"`
	DeletedAt              gorm.DeletedAt `gorm:"column:deleted_at;index"`
	DeletedBy              *uint64        `gorm:"column:deleted_by"`
}

func (SalesInvoiceDraftJSON) TableName() string { return "sales_invoice_draft_jsons" }
