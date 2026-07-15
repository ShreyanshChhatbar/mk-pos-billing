package model

import (
	"time"

	"gorm.io/gorm"
)

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
