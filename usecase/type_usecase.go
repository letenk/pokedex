package usecase

import (
	"context"

	"github.com/letenk/pokedex/models/domain"
	"github.com/letenk/pokedex/repository"
)

type TypeUsecase interface {
	FindAll(ctx context.Context) ([]domain.Type, error)
}

type typeUsecase struct {
	repository repository.TypeRepository
}

func NewUsecaseType(repository repository.TypeRepository) *typeUsecase {
	return &typeUsecase{repository}
}

func (u *typeUsecase) FindAll(ctx context.Context) ([]domain.Type, error) {
	// Find all
	types, err := u.repository.FindAll(ctx)

	if err != nil {
		return types, err
	}

	return types, nil
}
