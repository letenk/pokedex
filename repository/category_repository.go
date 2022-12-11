package repository

import (
	"context"
	"time"

	"github.com/letenk/pokedex/models/domain"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	FindAll(ctx context.Context) ([]domain.Category, error)
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *categoryRepository {
	return &categoryRepository{db}
}

func (r *categoryRepository) FindAll(ctx context.Context) ([]domain.Category, error) {
	// Create a context in order to disconnect after 15 seconds
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	var categories []domain.Category
	err := r.db.WithContext(ctx).Find(&categories).Error
	if err != nil {
		return categories, nil
	}

	return categories, nil
}
