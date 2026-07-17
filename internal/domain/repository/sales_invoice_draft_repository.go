package repository

import (
	"context"
	"time"

	"mk-pos-billing/internal/domain/model"
	"mk-pos-billing/internal/infrastructure/database"

	"gorm.io/gorm"
)

type SalesInvoiceDraftRepository interface {
	Tx(tx Transaction) SalesInvoiceDraftRepository
	GetDraftByID(ctx context.Context, id, storeID, organizationID uint64) (*model.SalesInvoiceDraftJSON, error)
	SaveDraft(ctx context.Context, draft *model.SalesInvoiceDraftJSON) error
	CreateDraft(ctx context.Context, draft *model.SalesInvoiceDraftJSON) error
	ReplaceDraftPayments(ctx context.Context, draftID uint64, payments []model.SalesInvoiceDraftPayment, deletedBy uint64) error
	AppendDraftPayments(ctx context.Context, payments []model.SalesInvoiceDraftPayment) error
	GetDraftPayments(ctx context.Context, draftID uint64) ([]model.SalesInvoiceDraftPayment, error)
	MarkDraftInvoiced(ctx context.Context, draftID uint64, updatedBy uint64) error
}

type salesInvoiceDraftRepository struct {
	db *gorm.DB
}

func NewSalesInvoiceDraftRepository(db *database.DB) SalesInvoiceDraftRepository {
	return &salesInvoiceDraftRepository{db: db.DB}
}

func (r *salesInvoiceDraftRepository) Tx(tx Transaction) SalesInvoiceDraftRepository {
	if tx == nil {
		return r
	}
	if gTx, ok := tx.(*gormTransaction); ok {
		return &salesInvoiceDraftRepository{db: gTx.db}
	}
	return r
}

func (r *salesInvoiceDraftRepository) WithContext(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx)
}

func (r *salesInvoiceDraftRepository) GetDraftByID(ctx context.Context, id, storeID, organizationID uint64) (*model.SalesInvoiceDraftJSON, error) {
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

func (r *salesInvoiceDraftRepository) SaveDraft(ctx context.Context, draft *model.SalesInvoiceDraftJSON) error {
	return r.WithContext(ctx).Omit("InvoiceNumber").Updates(draft).Error
}

func (r *salesInvoiceDraftRepository) CreateDraft(ctx context.Context, draft *model.SalesInvoiceDraftJSON) error {
	return r.WithContext(ctx).Omit("InvoiceNumber").Create(draft).Error
}

func (r *salesInvoiceDraftRepository) ReplaceDraftPayments(ctx context.Context, draftID uint64, payments []model.SalesInvoiceDraftPayment, deletedBy uint64) error {
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

func (r *salesInvoiceDraftRepository) AppendDraftPayments(ctx context.Context, payments []model.SalesInvoiceDraftPayment) error {
	if len(payments) == 0 {
		return nil
	}
	return r.WithContext(ctx).Create(&payments).Error
}

func (r *salesInvoiceDraftRepository) GetDraftPayments(ctx context.Context, draftID uint64) ([]model.SalesInvoiceDraftPayment, error) {
	var payments []model.SalesInvoiceDraftPayment
	err := r.WithContext(ctx).
		Where("sales_invoice_draft_id = ? AND deleted_at IS NULL", draftID).
		Find(&payments).Error
	return payments, err
}

func (r *salesInvoiceDraftRepository) MarkDraftInvoiced(ctx context.Context, draftID uint64, updatedBy uint64) error {
	return r.WithContext(ctx).
		Model(&model.SalesInvoiceDraftJSON{}).
		Where("id = ?", draftID).
		Updates(map[string]interface{}{"status": "INVOICED", "updated_by": updatedBy}).Error
}
