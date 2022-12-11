package repository

import (
	"context"
	"time"

	"github.com/letenk/pokedex/models/domain"
	"gorm.io/gorm"
)

type TypeRepository interface {
	FindAll(ctx context.Context) ([]domain.Type, error)
}

type typeRespository struct {
	db *gorm.DB
}

func NewTypeRespository(db *gorm.DB) *typeRespository {
	return &typeRespository{db}
}

func (r *typeRespository) FindAll(ctx context.Context) ([]domain.Type, error) {
	// Create a context in order to disconnect after 15 seconds
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	var types []domain.Type
	err := r.db.WithContext(ctx).Find(&types).Error
	if err != nil {
		return types, nil
	}

	return types, nil
}
