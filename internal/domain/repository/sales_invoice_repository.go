package repository

import (
	"context"
	"mk-pos-billing/internal/domain/model"
	"mk-pos-billing/internal/infrastructure/database"

	"gorm.io/gorm"
)



type SalesInvoiceRepository interface {
	Tx(tx Transaction) SalesInvoiceRepository
	CreateInvoice(ctx context.Context, invoice *model.SalesInvoice) error
	CreateInvoiceDetails(ctx context.Context, details *[]model.SalesInvoiceDetail) error
	CreateInvoiceTaxDetails(ctx context.Context, taxDetails *[]model.SalesInvoiceTaxDetail) error
	CreateInvoicePayments(ctx context.Context, payments []model.SalesInvoicePayment) error
}

type salesInvoiceRepository struct {
	db *gorm.DB
}

func NewSalesInvoiceRepository(db *database.DB) SalesInvoiceRepository {
	return &salesInvoiceRepository{db: db.DB}
}

func (r *salesInvoiceRepository) Tx(tx Transaction) SalesInvoiceRepository {
	if tx == nil {
		return r
	}
	if gTx, ok := tx.(*gormTransaction); ok {
		return &salesInvoiceRepository{db: gTx.db}
	}
	return r
}

func (r *salesInvoiceRepository) WithContext(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx)
}

func (r *salesInvoiceRepository) CreateInvoice(ctx context.Context, invoice *model.SalesInvoice) error {
	return r.WithContext(ctx).Create(invoice).Error
}

func (r *salesInvoiceRepository) CreateInvoiceDetails(ctx context.Context, details *[]model.SalesInvoiceDetail) error {
	if details == nil || len(*details) == 0 {
		return nil
	}
	return r.WithContext(ctx).Create(details).Error
}

func (r *salesInvoiceRepository) CreateInvoiceTaxDetails(ctx context.Context, taxDetails *[]model.SalesInvoiceTaxDetail) error {
	if taxDetails == nil || len(*taxDetails) == 0 {
		return nil
	}
	return r.WithContext(ctx).Create(taxDetails).Error
}

func (r *salesInvoiceRepository) CreateInvoicePayments(ctx context.Context, payments []model.SalesInvoicePayment) error {
	if len(payments) == 0 {
		return nil
	}
	return r.WithContext(ctx).Create(&payments).Error
}
