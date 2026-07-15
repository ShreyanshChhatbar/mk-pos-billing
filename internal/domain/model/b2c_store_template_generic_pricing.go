package model

import "gorm.io/gorm"

type B2CStoreTemplateGenericPricing struct {
	ID                   uint64         `gorm:"column:id;primaryKey"`
	B2CPricingTemplateID uint64         `gorm:"column:b_2_c_pricing_template_id"`
	ProductID            uint64         `gorm:"column:product_id"`
	SalesPrice           float64        `gorm:"column:sales_price"`
	DeletedAt            gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (B2CStoreTemplateGenericPricing) TableName() string {
	return "b_2_c_store_template_generic_pricing"
}
