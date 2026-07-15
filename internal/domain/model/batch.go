package model

import (
	"time"

	"gorm.io/gorm"
)

type Batch struct {
	ID           uint64    `gorm:"column:id;primaryKey"`
	ProductID    uint64    `gorm:"column:product_id"`
	WMSProductID *uint64   `gorm:"column:wms_product_id"`
	WMSBatchID   *uint64   `gorm:"column:wms_batch_id"`
	BatchNumber  string    `gorm:"column:batch_number"`
	MfgDate      *string   `gorm:"column:mfg_date"`
	ExpiryDate   time.Time `gorm:"column:expiry_date"`
	MRP          float64   `gorm:"column:mrp"`
	OldMRP       *float64  `gorm:"column:old_mrp"`
	PTR          *float64  `gorm:"column:ptr"`
	CAS          *float64  `gorm:"column:cas"`
	CTR          *float64  `gorm:"column:ctr"`
	CreatedBy    uint64    `gorm:"column:created_by"`
	UpdatedBy    *uint64   `gorm:"column:updated_by"`
	DeletedBy    *uint64   `gorm:"column:deleted_by"`
	CreatedAt    time.Time `gorm:"column:created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (Batch) TableName() string { return "batches" }
