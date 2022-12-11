package repository

import (
	"context"
	"time"

	"github.com/letenk/pokedex/models/domain"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindByUsername(ctx context.Context, username string) (domain.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{db}
}

func (r *userRepository) FindByUsername(ctx context.Context, username string) (domain.User, error) {
	// Create a context in order to disconnect after 15 seconds
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	var user domain.User

	err := r.db.WithContext(ctx).Where("username = ?", username).Find(&user).Error
	if err != nil {
		return user, err
	}

	return user, nil
}
