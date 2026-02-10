package database

import (
	"context"
	"database/sql"

	"gorm.io/gorm"
)

type DB struct {
	DB *gorm.DB
}

func (d *DB) BeginTx(ctx context.Context, opts ...*sql.TxOptions) (*gorm.DB, error) {
	tx := d.DB.WithContext(ctx).Begin(opts...)
	return tx, tx.Error
}
