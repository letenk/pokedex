package domain

import (
	"time"
)

type Monster struct {
	ID          string
	Name        string
	CategoryID  string
	Description string
	Length      float32
	Weight      uint16
	Hp          uint16
	Attack      uint16
	Defends     uint16
	Speed       uint16
	Catched     bool
	Image       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	TypeID      []string `gorm:"-"`                       // Ignore as field column
	Types       []Type   `gorm:"many2many:monster_types"` // Relation many to many to category
	Category    Category // Relation one to many to category
}
