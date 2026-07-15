package model

import (
	"time"

	"gorm.io/gorm"
)

type StoreInventoryTransaction struct {
	ID              uint64         `gorm:"column:id;primaryKey;autoIncrement"`
	StoreID         uint64         `gorm:"column:store_id"`
	ProductID       uint64         `gorm:"column:product_id"`
	BatchCode       string         `gorm:"column:batch_code"`
	StoreBatchID    uint64         `gorm:"column:store_batch_id"`
	ExpiryDate      time.Time      `gorm:"column:expiry_date"`
	Quantity        int            `gorm:"column:quantity"`
	Rate            float64        `gorm:"column:rate"`
	TotalAmount     float64        `gorm:"column:total_amount"`
	VoucherType     string         `gorm:"column:voucher_type"`
	VoucherID       uint64         `gorm:"column:voucher_id"`
	CreatedBy       uint64         `gorm:"column:created_by"`
	UpdatedBy       *uint64        `gorm:"column:updated_by"`
	DeletedBy       *uint64        `gorm:"column:deleted_by"`
	CreatedAt       time.Time      `gorm:"column:created_at"`
	UpdatedAt       time.Time      `gorm:"column:updated_at"`
	TransactionTime time.Time      `gorm:"column:transaction_time"`
	DeletedAt       gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (StoreInventoryTransaction) TableName() string { return "store_inventory_transactions" }
