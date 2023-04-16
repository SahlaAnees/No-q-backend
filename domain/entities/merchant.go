package entities

import "time"

type Merchant struct {
	ID        int64
	Name      string
	Category  string
	Email     string
	Password  string
	Facebook  string
	Instagram string
	Website   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
