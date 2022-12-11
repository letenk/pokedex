package domain

import "time"

type User struct {
	ID        string
	Fullname  string
	Username  string
	Password  string
	Role      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
