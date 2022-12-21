package domain

import "time"

type Type struct {
	ID        string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	Monsters  []*Monster `gorm:"many2many:monster_types"`
}
