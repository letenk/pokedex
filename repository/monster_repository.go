package repository

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/jackc/pgconn"
	"github.com/letenk/pokedex/models/domain"
	"github.com/letenk/pokedex/models/web"
	"gorm.io/gorm"
)

type MonsterRepository interface {
	Create(ctx context.Context, monster domain.Monster) (domain.Monster, error)
	FindAll(ctx context.Context, reqQuery web.MonsterQueryRequest) ([]domain.Monster, error)
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

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		// Create monster
		err := tx.WithContext(ctx).Create(&monster).Error
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				if pgErr.Code == "23503" {
					return errors.New("invalid category id or type id, please check valid id in each of their list")
				}
			}
			return err
		}

		// Insert relation monster id and type id
		for _, typeID := range monster.TypeID {
			monsterTag := new(domain.MonsterType)
			monsterTag.MonsterID = monster.ID
			monsterTag.TypeID = typeID
			err = tx.WithContext(ctx).Create(&monsterTag).Error
			if err != nil {
				var pgErr *pgconn.PgError
				if errors.As(err, &pgErr) {
					if pgErr.Code == "23503" {
						return errors.New("invalid category id or type id, please check valid id in each of their list")
					}
				}
				return err
			}
		}
		return nil
	})

	if err != nil {
		return monster, err
	}

	return monster, nil
}

func (r *monsterRespository) FindAll(ctx context.Context, reqQuery web.MonsterQueryRequest) ([]domain.Monster, error) {
	// Create a context in order to disconnect
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	// Cancel context after all process ends
	defer cancel()

	var monsters []domain.Monster

	db := r.db.WithContext(ctx).Model(&domain.Monster{})

	if reqQuery.Name != "" {
		db = db.Where("lower(monsters.name) LIKE lower(?)", "%"+reqQuery.Name+"%")
	}

	if reqQuery.Catched != "" {
		boolCatched, err := strconv.ParseBool(reqQuery.Catched)
		if err != nil {
			return monsters, err
		}
		db = db.Where("catched = ?", boolCatched)
	}

	// For use query parameter order, sort must not empty
	if reqQuery.Order != "" {
		if reqQuery.Sort != "" && reqQuery.Order != "" {
			sortQuery := fmt.Sprintf("%s %s", reqQuery.Sort, reqQuery.Order)
			db.Order(sortQuery)
		} else {
			return monsters, errors.New("for use order, query parameter sort is required")
		}
	}

	if reqQuery.Sort != "" {
		db.Order(reqQuery.Sort)
	}

	if len(reqQuery.Types) != 0 {
		err := db.Preload("Category").Preload("Types").Joins("inner join monster_types mt on mt.monster_id = monsters.id ").Joins("inner join types t on t.id = mt.type_id ").Where("t.id IN ?", reqQuery.Types).Group("monsters.id").Find(&monsters).Error
		if err != nil {
			return monsters, err
		}

		return monsters, nil
	} else {
		err := db.Preload("Category").Preload("Types").Find(&monsters).Error
		if err != nil {
			return monsters, err
		}

		return monsters, nil
	}

}
