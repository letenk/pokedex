package repository

import (
	"context"
	"time"

	"github.com/letenk/pokedex/models/domain"
	"gorm.io/gorm"
)

type MonsterRepository interface {
	Create(ctx context.Context, monster domain.Monster) (domain.Monster, error)
}

type monsterRespository struct {
	db *gorm.DB
}

func NewMonsterRespository(db *gorm.DB) *monsterRespository {
	return &monsterRespository{db}
}

func (r *monsterRespository) Create(ctx context.Context, monster domain.Monster) (domain.Monster, error) {
	// Create a context in order to disconnect
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	// Cancel context after all process ends
	defer cancel()

	r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		// Create monster
		err := tx.WithContext(ctx).Create(&monster).Error
		if err != nil {
			return err
		}

		// Insert relation monster id and type id
		for _, typeID := range monster.TypeID {
			monsterTag := new(domain.MonsterType)
			monsterTag.MonsterID = monster.ID
			monsterTag.TypeID = typeID
			err = tx.WithContext(ctx).Create(&monsterTag).Error
			if err != nil {
				return err
			}
		}
		return nil
	})

	return monster, nil
}
