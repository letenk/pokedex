package usecase

import (
	"context"

	"github.com/letenk/pokedex/models/domain"
	"github.com/letenk/pokedex/models/web"
	"github.com/letenk/pokedex/repository"
)

type MonsterUsecase interface {
	Create(ctx context.Context, monster web.MonsterCreateRequest, imagePathLocation string) (domain.Monster, error)
}

type monsterUsecase struct {
	repository repository.MonsterRepository
}

func NewUsecaseMonster(repository repository.MonsterRepository) *monsterUsecase {
	return &monsterUsecase{repository}
}

func (u *monsterUsecase) Create(ctx context.Context, req web.MonsterCreateRequest, imagePathLocation string) (domain.Monster, error) {

	// Passing data request into object monster
	monster := domain.Monster{
		Name:        req.Name,
		CategoryID:  req.CategoryID,
		Description: req.Description,
		Length:      req.Length,
		Weight:      req.Weight,
		Hp:          req.Hp,
		Attack:      req.Attack,
		Defends:     req.Defends,
		Speed:       req.Speed,
		Image:       imagePathLocation,
		TypeID:      req.TypeID,
	}

	// Create
	monster, err := u.repository.Create(ctx, monster)
	if err != nil {
		return monster, err
	}

	return monster, nil
}
