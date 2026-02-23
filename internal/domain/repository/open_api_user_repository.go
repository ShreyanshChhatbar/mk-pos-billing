package repository

import (
	"context"
	"errors"
	"mk-pos-billing/internal/domain/model"
	"mk-pos-billing/internal/infrastructure/database"

	"gorm.io/gorm"
)

type OpenApiUserRepository struct {
	db *database.DB
}

func NewOpenApiUserRepository(db *database.DB) *OpenApiUserRepository {
	return &OpenApiUserRepository{db: db}
}

func (r *OpenApiUserRepository) FindByID(ctx context.Context, id uint64) (*model.OpenAPIUser, error) {
	var user model.OpenAPIUser
	err := r.db.DB.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("open api user not found")
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}
