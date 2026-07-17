package repository

import (
	"context"
	"time"

	"mk-pos-billing/internal/domain/model"
	"mk-pos-billing/internal/infrastructure/database"

	"gorm.io/gorm"
)

type StoreInventoryRepository struct {
	db *gorm.DB
}

type BatchStockRow struct {
	StoreID      uint64    `gorm:"column:store_id"`
	ProductID    uint64    `gorm:"column:product_id"`
	BatchCode    string    `gorm:"column:batch_code"`
	StoreBatchID uint64    `gorm:"column:store_batch_id"`
	PurchaseRate float64   `gorm:"column:purchase_rate"`
	ExpiryDate   time.Time `gorm:"column:expiry_date"`
	ClosingStock int       `gorm:"column:closing_stock"`
	Key          string    `gorm:"column:key"`
}

func NewStoreInventoryRepository(db *database.DB) *StoreInventoryRepository {
	return &StoreInventoryRepository{db: db.DB}
}

func (r *StoreInventoryRepository) Tx(tx Transaction) *StoreInventoryRepository {
	if tx == nil {
		return r
	}
	if gTx, ok := tx.(*gormTransaction); ok {
		return &StoreInventoryRepository{db: gTx.db}
	}
	return r
}

func (r *StoreInventoryRepository) WithContext(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx)
}

func (r *StoreInventoryRepository) InsertInventoryTransactions(ctx context.Context, txns []model.StoreInventoryTransaction) error {
	if len(txns) == 0 {
		return nil
	}
	return r.WithContext(ctx).Create(&txns).Error
}

func (r *StoreInventoryRepository) FetchBatchStocks(ctx context.Context, storeID uint64, productIDs []uint64, batchCodes []string, ignoreStockCheck bool, expiryCutoff time.Time) (map[string][]BatchStockRow, error) {
	if len(productIDs) == 0 || len(batchCodes) == 0 {
		return map[string][]BatchStockRow{}, nil
	}

	rows := make([]BatchStockRow, 0)
	q := r.WithContext(ctx).
		Table("store_inventories as si").
		Joins("join store_batches on si.store_batch_id = store_batches.id and store_batches.store_id = ?", storeID).
		Joins("join batches on store_batches.batch_id = batches.id").
		Select("si.store_id, si.product_id, si.batch_code, store_batches.id as store_batch_id, store_batches.purchase_rate, batches.expiry_date, si.closing_stock, CONCAT(si.product_id, '_', si.batch_code) as key").
		Where("si.is_active = true").
		Where("si.deleted_at IS NULL").
		Where("si.store_id = ?", storeID).
		Where("si.product_id IN ?", productIDs).
		Where("si.batch_code IN ?", batchCodes).
		Where("batches.expiry_date >= ?", expiryCutoff)

	if !ignoreStockCheck {
		q = q.Where("si.closing_stock > 0")
	}

	if err := q.Order("si.product_id, si.batch_code, store_batches.id").Scan(&rows).Error; err != nil {
		return nil, err
	}

	grouped := make(map[string][]BatchStockRow)
	for _, row := range rows {
		grouped[row.Key] = append(grouped[row.Key], row)
	}
	return grouped, nil
}


