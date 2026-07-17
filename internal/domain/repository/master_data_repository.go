package repository

import (
	"context"

	"mk-pos-billing/internal/infrastructure/database"

	"gorm.io/gorm"
)

// MasterDataRepository acts as a temporary dumping ground for simple existence checks
// against global entities (Users, Customers, Doctors, Store configs) since we don't
// have dedicated domains/services for them in this bounded context yet.
type MasterDataRepository interface {
	Tx(tx Transaction) MasterDataRepository
	ExistsActiveUser(ctx context.Context, id uint64) bool
	ExistsActiveCustomer(ctx context.Context, id uint64) bool
	ExistsPatientForCustomer(ctx context.Context, patientID, customerID uint64) bool
	ExistsDoctor(ctx context.Context, doctorID uint64) bool
	ExistsCustomerAddress(ctx context.Context, customerAddressID, customerID uint64, requirePinCode bool) bool
	ExistsStorePaymentMethod(ctx context.Context, storePaymentMethodID, storeID, organizationID uint64) bool
}

type masterDataRepository struct {
	db *gorm.DB
}

func NewMasterDataRepository(db *database.DB) MasterDataRepository {
	return &masterDataRepository{db: db.DB}
}

func (r *masterDataRepository) Tx(tx Transaction) MasterDataRepository {
	if tx == nil {
		return r
	}
	if gTx, ok := tx.(*gormTransaction); ok {
		return &masterDataRepository{db: gTx.db}
	}
	return r
}

func (r *masterDataRepository) WithContext(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx)
}

func (r *masterDataRepository) ExistsActiveUser(ctx context.Context, id uint64) bool {
	return r.exists(ctx, "users", "id = ? AND is_active = true AND deleted_at IS NULL", id)
}

func (r *masterDataRepository) ExistsActiveCustomer(ctx context.Context, id uint64) bool {
	return r.exists(ctx, "customers", "id = ? AND is_active = true AND deleted_at IS NULL", id)
}

func (r *masterDataRepository) ExistsPatientForCustomer(ctx context.Context, patientID, customerID uint64) bool {
	return r.exists(ctx, "patients", "id = ? AND customer_id = ? AND is_active = true AND deleted_at IS NULL", patientID, customerID)
}

func (r *masterDataRepository) ExistsDoctor(ctx context.Context, doctorID uint64) bool {
	return r.exists(ctx, "doctors", "id = ? AND is_active = true AND deleted_at IS NULL", doctorID)
}

func (r *masterDataRepository) ExistsCustomerAddress(ctx context.Context, customerAddressID, customerID uint64, requirePinCode bool) bool {
	query := "id = ? AND customer_id = ? AND is_active = true AND deleted_at IS NULL"
	if requirePinCode {
		query += " AND pin_code IS NOT NULL"
	}
	return r.exists(ctx, "customer_addresses", query, customerAddressID, customerID)
}

func (r *masterDataRepository) ExistsStorePaymentMethod(ctx context.Context, storePaymentMethodID, storeID, organizationID uint64) bool {
	return r.exists(ctx, "store_payment_methods", "id = ? AND store_id = ? AND organization_id = ? AND is_active = true AND deleted_at IS NULL", storePaymentMethodID, storeID, organizationID)
}

func (r *masterDataRepository) exists(ctx context.Context, table, query string, args ...interface{}) bool {
	var count int64
	if err := r.WithContext(ctx).Table(table).Where(query, args...).Count(&count).Error; err != nil {
		return false
	}
	return count > 0
}
