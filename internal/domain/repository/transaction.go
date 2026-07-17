package repository

import (
	"context"

	"mk-pos-billing/internal/infrastructure/database"

	"gorm.io/gorm"
)

// Transaction represents an opaque database transaction.
// It hides the underlying implementation from the service layer.
type Transaction interface {
	isTransaction() // unexported marker method
}

type gormTransaction struct {
	db *gorm.DB
}

func (t *gormTransaction) isTransaction() {}

// TransactionManager coordinates database transactions.
type TransactionManager interface {
	Do(ctx context.Context, fn func(tx Transaction) error) error
}

type gormTransactionManager struct {
	db *gorm.DB
}

// NewTransactionManager creates a new transaction manager.
func NewTransactionManager(db *database.DB) TransactionManager {
	return &gormTransactionManager{db: db.DB}
}

func (tm *gormTransactionManager) Do(ctx context.Context, fn func(tx Transaction) error) error {
	return tm.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(&gormTransaction{db: tx})
	})
}
