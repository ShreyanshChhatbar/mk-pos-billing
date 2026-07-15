package model

import (
	"time"

	"gorm.io/gorm"
)

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
