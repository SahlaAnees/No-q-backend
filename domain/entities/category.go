package entities

import "time"

type Category struct {
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
