package model

import (
	"time"

	"gorm.io/gorm"
)

type StoreBatch struct {
	ID           uint64         `gorm:"column:id;primaryKey"`
	StoreID      uint64         `gorm:"column:store_id"`
	BatchID      uint64         `gorm:"column:batch_id"`
	ProductID    uint64         `gorm:"column:product_id"`
	PurchaseRate float64        `gorm:"column:purchase_rate"`
	HubID        *uint64        `gorm:"column:hub_id"`
	CreatedBy    uint64         `gorm:"column:created_by"`
	UpdatedBy    *uint64        `gorm:"column:updated_by"`
	DeletedBy    *uint64        `gorm:"column:deleted_by"`
	CreatedAt    time.Time      `gorm:"column:created_at"`
	UpdatedAt    time.Time      `gorm:"column:updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (StoreBatch) TableName() string { return "store_batches" }
