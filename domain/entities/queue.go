package entities

import "time"

type Queue struct {
	ID               int64
	MerchantID       int64
	Name             string
	Interval         int // minutes
	StartTime        time.Time
	EndTime          time.Time
	IsAvailable      bool
	UnavailableDates []time.Time
	ReservedSlots    []ReservedSlots
	CreatedAt        time.Time
}

type UnavailableDates struct {
	QueueID int64
	Dates   []time.Time
}

type ReservedSlots struct {
	TokenNo    int64
	QueueID    int64
	StartTime  time.Time
	EndTime    time.Time
	ReservedBy User
	CreatedAt  time.Time
}

type User struct {
	ID        int64
	Name      string
	Phone     string
	Email     string
	CreatedAt time.Time
}
