package model

import "gorm.io/gorm"

type OpenAPIUser struct {
	ID        uint64         `gorm:"column:id;primaryKey"`
	TenantID  uint64         `gorm:"column:tenant_id"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (OpenAPIUser) TableName() string { return "open_api_users" }
