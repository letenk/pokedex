package usecase

import (
	"context"

	"github.com/letenk/pokedex/models/domain"
	"github.com/letenk/pokedex/repository"
)

type CategoryUsecase interface {
	FindAll(ctx context.Context) ([]domain.Category, error)
}

type categoryUsecase struct {
	repository repository.CategoryRepository
}

func NewUsecaseCategory(repository repository.CategoryRepository) *categoryUsecase {
	return &categoryUsecase{repository}
}

func (u *categoryUsecase) FindAll(ctx context.Context) ([]domain.Category, error) {
	// Find all
	categories, err := u.repository.FindAll(ctx)

	if err != nil {
		return categories, err
	}

	return categories, nil
}
