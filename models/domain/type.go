package domain

import "time"

type Type struct {
	ID        string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
