package database

import "gorm.io/gorm"

// applyGlobalScopes registers default scopes (e.g., soft-delete filter) on the DB handle.
func applyGlobalScopes(db *gorm.DB) *gorm.DB {
	return db.Scopes(softDeleteScope)
}

func softDeleteScope(tx *gorm.DB) *gorm.DB {
	if tx.Statement == nil || tx.Statement.Schema == nil || tx.Statement.Unscoped {
		return tx
	}

	if tx.Statement.Schema.LookUpField("DeletedAt") != nil {
		return tx.Where(tx.Statement.Quote("deleted_at") + " IS NULL")
	}

	return tx
}
