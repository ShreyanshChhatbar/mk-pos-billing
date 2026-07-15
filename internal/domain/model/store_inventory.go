package model

import (
	"time"

	"gorm.io/gorm"
)

type StoreInventory struct {
	ID                          uint64         `gorm:"column:id"`
	StoreInventoryTransactionID *uint64        `gorm:"column:store_inventory_transaction_id"`
	StoreID                     uint64         `gorm:"column:store_id"`
	ProductID                   uint64         `gorm:"column:product_id"`
	BatchCode                   string         `gorm:"column:batch_code"`
	StoreBatchID                uint64         `gorm:"column:store_batch_id"`
	OpeningStock                int            `gorm:"column:opening_stock"`
	Quantity                    int            `gorm:"column:quantity"`
	ClosingStock                int            `gorm:"column:closing_stock"`
	IsActive                    bool           `gorm:"column:is_active"`
	CreatedBy                   uint64         `gorm:"column:created_by"`
	UpdatedBy                   *uint64        `gorm:"column:updated_by"`
	DeletedBy                   *uint64        `gorm:"column:deleted_by"`
	CreatedAt                   time.Time      `gorm:"column:created_at"`
	UpdatedAt                   time.Time      `gorm:"column:updated_at"`
	DeletedAt                   gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (StoreInventory) TableName() string { return "store_inventories" }
