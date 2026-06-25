package repository

import (
	"context"
	"errors"
	"fmt"
	"mk-pos-billing/internal/domain/model"
	"mk-pos-billing/internal/infrastructure/database"
	"strings"
	"time"

	"gorm.io/gorm"
)

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

type ProductValidationData struct {
	ID                uint64 `gorm:"column:id"`
	SalesUnit         int    `gorm:"column:sales_unit"`
	WSCode            int    `gorm:"column:ws_code"`
	IsMSPProduct      bool   `gorm:"column:is_msp_product"`
	ScheduledTypeCode string `gorm:"column:scheduled_type_code"`
}

type SalesInvoiceRepository struct {
	db *gorm.DB
}

func NewSalesInvoiceRepository(db *database.DB) *SalesInvoiceRepository {
	return &SalesInvoiceRepository{db: db.DB}
}

func (r *SalesInvoiceRepository) Tx(tx *gorm.DB) *SalesInvoiceRepository {
	if tx == nil {
		return r
	}
	return &SalesInvoiceRepository{db: tx}
}

func (r *SalesInvoiceRepository) WithContext(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx)
}

func (r *SalesInvoiceRepository) GetDraftByID(ctx context.Context, id, storeID, organizationID uint64) (*model.SalesInvoiceDraftJSON, error) {
	var draft model.SalesInvoiceDraftJSON
	err := r.WithContext(ctx).
		Where("id = ? AND store_id = ? AND organization_id = ?", id, storeID, organizationID).
		Where("status IN ?", []string{"DRAFT", "PAYMENT_PENDING"}).
		First(&draft).Error
	if err != nil {
		return nil, err
	}
	return &draft, nil
}

func (r *SalesInvoiceRepository) SaveDraft(ctx context.Context, draft *model.SalesInvoiceDraftJSON) error {
	return r.WithContext(ctx).Save(draft).Error
}

func (r *SalesInvoiceRepository) CreateDraft(ctx context.Context, draft *model.SalesInvoiceDraftJSON) error {
	return r.WithContext(ctx).Create(draft).Error
}

func (r *SalesInvoiceRepository) ReplaceDraftPayments(ctx context.Context, draftID uint64, payments []model.SalesInvoiceDraftPayment, deletedBy uint64) error {
	if err := r.WithContext(ctx).
		Model(&model.SalesInvoiceDraftPayment{}).
		Where("sales_invoice_draft_id = ? AND deleted_at IS NULL", draftID).
		Updates(map[string]interface{}{"deleted_at": time.Now(), "deleted_by": deletedBy}).Error; err != nil {
		return err
	}
	if len(payments) == 0 {
		return nil
	}
	return r.WithContext(ctx).Create(&payments).Error
}

func (r *SalesInvoiceRepository) AppendDraftPayments(ctx context.Context, payments []model.SalesInvoiceDraftPayment) error {
	if len(payments) == 0 {
		return nil
	}
	return r.WithContext(ctx).Create(&payments).Error
}

func (r *SalesInvoiceRepository) GetDraftPayments(ctx context.Context, draftID uint64) ([]model.SalesInvoiceDraftPayment, error) {
	var payments []model.SalesInvoiceDraftPayment
	err := r.WithContext(ctx).
		Where("sales_invoice_draft_id = ? AND deleted_at IS NULL", draftID).
		Find(&payments).Error
	return payments, err
}

func (r *SalesInvoiceRepository) CreateInvoice(ctx context.Context, invoice *model.SalesInvoice) error {
	return r.WithContext(ctx).Create(invoice).Error
}

func (r *SalesInvoiceRepository) CreateInvoiceDetails(ctx context.Context, details []model.SalesInvoiceDetail) error {
	if len(details) == 0 {
		return nil
	}
	return r.WithContext(ctx).Create(&details).Error
}

func (r *SalesInvoiceRepository) CreateInvoicePayments(ctx context.Context, payments []model.SalesInvoicePayment) error {
	if len(payments) == 0 {
		return nil
	}
	return r.WithContext(ctx).Create(&payments).Error
}

func (r *SalesInvoiceRepository) MarkDraftInvoiced(ctx context.Context, draftID uint64, updatedBy uint64) error {
	return r.WithContext(ctx).
		Model(&model.SalesInvoiceDraftJSON{}).
		Where("id = ?", draftID).
		Updates(map[string]interface{}{"status": "INVOICED", "updated_by": updatedBy}).Error
}

func (r *SalesInvoiceRepository) ReduceInventoryStock(ctx context.Context, storeID, productID, storeBatchID uint64, batchCode string, quantity int) error {
	result := r.WithContext(ctx).
		Model(&model.StoreInventory{}).
		Where("store_id = ? AND product_id = ? AND store_batch_id = ? AND batch_code = ? AND closing_stock >= ? AND deleted_at IS NULL", storeID, productID, storeBatchID, batchCode, quantity).
		Update("closing_stock", gorm.Expr("closing_stock - ?", quantity))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("inventory row not found for product %d batch %s", productID, batchCode)
	}
	return nil
}

func (r *SalesInvoiceRepository) FetchBatchStocks(ctx context.Context, storeID uint64, productIDs []uint64, batchCodes []string, ignoreStockCheck bool, expiryCutoff time.Time) (map[string][]BatchStockRow, error) {
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

func (r *SalesInvoiceRepository) FetchBatchMeta(ctx context.Context, productIDs []uint64, batchCodes []string) (map[string]model.Batch, error) {
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

func (r *SalesInvoiceRepository) ExistsActiveUser(ctx context.Context, id uint64) bool {
	return r.exists(ctx, "users", "id = ? AND is_active = true AND deleted_at IS NULL", id)
}

func (r *SalesInvoiceRepository) ExistsActiveCustomer(ctx context.Context, id uint64) bool {
	return r.exists(ctx, "customers", "id = ? AND is_active = true AND deleted_at IS NULL", id)
}

func (r *SalesInvoiceRepository) ExistsPatientForCustomer(ctx context.Context, patientID, customerID uint64) bool {
	return r.exists(ctx, "patients", "id = ? AND customer_id = ? AND is_active = true AND deleted_at IS NULL", patientID, customerID)
}

func (r *SalesInvoiceRepository) ExistsDoctor(ctx context.Context, doctorID uint64) bool {
	return r.exists(ctx, "doctors", "id = ? AND is_active = true AND deleted_at IS NULL", doctorID)
}

func (r *SalesInvoiceRepository) ExistsCustomerAddress(ctx context.Context, customerAddressID, customerID uint64, requirePinCode bool) bool {
	query := "id = ? AND customer_id = ? AND is_active = true AND deleted_at IS NULL"
	if requirePinCode {
		query += " AND pin_code IS NOT NULL"
	}
	return r.exists(ctx, "customer_addresses", query, customerAddressID, customerID)
}

func (r *SalesInvoiceRepository) ExistsActiveProduct(ctx context.Context, productID uint64) bool {
	return r.exists(ctx, "products", "id = ? AND is_active = true AND deleted_at IS NULL", productID)
}

func (r *SalesInvoiceRepository) ExistsBatchCode(ctx context.Context, batchCode string) bool {
	return r.exists(ctx, "batches", "batch_number = ? AND deleted_at IS NULL", batchCode)
}

func (r *SalesInvoiceRepository) ExistsStorePaymentMethod(ctx context.Context, storePaymentMethodID, storeID, organizationID uint64) bool {
	return r.exists(ctx, "store_payment_methods", "id = ? AND store_id = ? AND organization_id = ? AND is_active = true AND deleted_at IS NULL", storePaymentMethodID, storeID, organizationID)
}

func (r *SalesInvoiceRepository) GetProductValidationData(ctx context.Context, productID uint64) (*ProductValidationData, error) {
	var data ProductValidationData
	err := r.WithContext(ctx).
		Table("products").
		Select("id, sales_unit, ws_code, is_msp_product, scheduled_type_code").
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

func (r *SalesInvoiceRepository) containsNarcotics(ctx context.Context, productIDs []uint64, narcoticsLabel string) (bool, error) {
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

func (r *SalesInvoiceRepository) HasNarcoticsProduct(ctx context.Context, productIDs []uint64, narcoticsLabel string) bool {
	has, err := r.containsNarcotics(ctx, productIDs, narcoticsLabel)
	if err != nil {
		return false
	}
	return has
}

func (r *SalesInvoiceRepository) exists(ctx context.Context, table, query string, args ...interface{}) bool {
	var count int64
	if err := r.WithContext(ctx).Table(table).Where(query, args...).Count(&count).Error; err != nil {
		return false
	}
	return count > 0
}
