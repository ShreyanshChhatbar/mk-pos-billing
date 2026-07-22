package model

import (
	"time"

	"gorm.io/gorm"
)

type SalesInvoiceTaxDetail struct {
	ID                   uint64         `gorm:"column:id;primaryKey"`
	SalesInvoiceDetailID uint64         `gorm:"column:sales_invoice_detail_id"`
	StoreID              uint64         `gorm:"column:store_id"`
	TaxType              string         `gorm:"column:tax_type"`
	TaxRate              float64        `gorm:"column:tax_rate"`
	TaxName              string         `gorm:"column:tax_name"`
	TaxAmount            float64        `gorm:"column:tax_amount"`
	IsActive             bool           `gorm:"column:is_active"`
	CreatedBy            uint64         `gorm:"column:created_by"`
	UpdatedBy            *uint64        `gorm:"column:updated_by"`
	CreatedAt            time.Time      `gorm:"column:created_at"`
	UpdatedAt            time.Time      `gorm:"column:updated_at"`
	DeletedAt            gorm.DeletedAt `gorm:"column:deleted_at;index"`
	DeletedBy            *uint64        `gorm:"column:deleted_by"`
}

func (SalesInvoiceTaxDetail) TableName() string { return "sales_invoice_tax_details" }
