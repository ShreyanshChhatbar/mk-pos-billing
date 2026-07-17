package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"mk-pos-billing/internal/domain/model"
	"mk-pos-billing/internal/infrastructure/database"

	"gorm.io/gorm"
)

type ProductRepository interface {
	Tx(tx Transaction) ProductRepository
	ExistsActiveProduct(ctx context.Context, productID uint64) bool
	ExistsBatchCode(ctx context.Context, batchCode string) bool
	GetProductValidationData(ctx context.Context, productID uint64) (*ProductValidationData, error)
	HasNarcoticsProduct(ctx context.Context, productIDs []uint64, narcoticsLabel string) bool
	FetchGenericPricings(ctx context.Context, templateIDs []uint64, productIDs []uint64) (map[string]model.B2CStoreTemplateGenericPricing, error)
	FetchBatchMeta(ctx context.Context, productIDs []uint64, batchCodes []string) (map[string]model.Batch, error)
}

type productRepository struct {
	db *gorm.DB
}

type ProductValidationData struct {
	ID                uint64 `gorm:"column:id"`
	SalesUnit         int    `gorm:"column:sales_unit"`
	WSCode            int    `gorm:"column:ws_code"`
	ScheduledTypeCode string `gorm:"column:scheduled_type_code"`
}

func NewProductRepository(db *database.DB) ProductRepository {
	return &productRepository{db: db.DB}
}

func (r *productRepository) Tx(tx Transaction) ProductRepository {
	if tx == nil {
		return r
	}
	if gTx, ok := tx.(*gormTransaction); ok {
		return &productRepository{db: gTx.db}
	}
	return r
}

func (r *productRepository) WithContext(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx)
}

func (r *productRepository) ExistsActiveProduct(ctx context.Context, productID uint64) bool {
	return r.exists(ctx, "products", "id = ? AND is_active = true AND deleted_at IS NULL", productID)
}

func (r *productRepository) ExistsBatchCode(ctx context.Context, batchCode string) bool {
	return r.exists(ctx, "batches", "batch_number = ? AND deleted_at IS NULL", batchCode)
}

func (r *productRepository) GetProductValidationData(ctx context.Context, productID uint64) (*ProductValidationData, error) {
	var data ProductValidationData
	err := r.WithContext(ctx).
		Table("products").
		Select("id, sales_unit, ws_code, scheduled_type_code").
		Where("id = ? AND is_active = true AND deleted_at IS NULL", productID).
		First(&data).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if data.SalesUnit <= 0 {
		data.SalesUnit = 1
	}
	return &data, nil
}

func (r *productRepository) containsNarcotics(ctx context.Context, productIDs []uint64, narcoticsLabel string) (bool, error) {
	if len(productIDs) == 0 {
		return false, nil
	}
	var count int64
	err := r.WithContext(ctx).
		Table("products").
		Where("id IN ? AND is_active = true AND deleted_at IS NULL", productIDs).
		Where("LOWER(scheduled_type_code) = ?", strings.ToLower(narcoticsLabel)).
		Count(&count).Error
	return count > 0, err
}

func (r *productRepository) HasNarcoticsProduct(ctx context.Context, productIDs []uint64, narcoticsLabel string) bool {
	has, err := r.containsNarcotics(ctx, productIDs, narcoticsLabel)
	if err != nil {
		return false
	}
	return has
}

func (r *productRepository) FetchGenericPricings(ctx context.Context, templateIDs []uint64, productIDs []uint64) (map[string]model.B2CStoreTemplateGenericPricing, error) {
	if len(templateIDs) == 0 || len(productIDs) == 0 {
		return map[string]model.B2CStoreTemplateGenericPricing{}, nil
	}

	var pricings []model.B2CStoreTemplateGenericPricing
	if err := r.WithContext(ctx).
		Where("b_2_c_pricing_template_id IN ?", templateIDs).
		Where("product_id IN ?", productIDs).
		Where("deleted_at IS NULL").
		Find(&pricings).Error; err != nil {
		return nil, err
	}

	pricingMap := make(map[string]model.B2CStoreTemplateGenericPricing, len(pricings))
	for _, p := range pricings {
		key := fmt.Sprintf("%d:%d", p.B2CPricingTemplateID, p.ProductID)
		pricingMap[key] = p
	}
	return pricingMap, nil
}

func (r *productRepository) FetchBatchMeta(ctx context.Context, productIDs []uint64, batchCodes []string) (map[string]model.Batch, error) {
	if len(productIDs) == 0 || len(batchCodes) == 0 {
		return map[string]model.Batch{}, nil
	}

	batches := make([]model.Batch, 0)
	if err := r.WithContext(ctx).
		Where("product_id IN ?", productIDs).
		Where("batch_number IN ?", batchCodes).
		Find(&batches).Error; err != nil {
		return nil, err
	}

	batchMap := make(map[string]model.Batch, len(batches))
	for _, b := range batches {
		key := fmt.Sprintf("%d_%s", b.ProductID, b.BatchNumber)
		batchMap[key] = b
	}
	return batchMap, nil
}

func (r *productRepository) exists(ctx context.Context, table, query string, args ...interface{}) bool {
	var count int64
	if err := r.WithContext(ctx).Table(table).Where(query, args...).Count(&count).Error; err != nil {
		return false
	}
	return count > 0
}
