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
	TypeID      []string `gorm:"-"` // Ignore as field column
	Type        []*Type  `gorm:"many2many:monster_type"`
}
