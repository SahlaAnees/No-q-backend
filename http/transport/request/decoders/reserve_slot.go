package decoders

import (
	"errors"
	"no-q-solution/domain/entities"
	"time"
)

type ReserveSlot struct {
	QueueID    int64     `json:"queue_id" validate:"required"`
	StartTime  time.Time `json:"start_time" validate:"required"`
	EndTime    time.Time `json:"end_time" validate:"required"`
	ReservedBy User      `json:"reserved_by" validate:"required"`
}

type User struct {
	Name  string `json:"name" validate:"required"`
	Phone string `json:"phone" validate:"required"`
	Email string `json:"email"`
}

func (r ReserveSlot) Format() string {
	return `
		{
			"queue_id": 1,
			"start_time": "2023-04-14T10:00:00Z",
			"end_time": "2023-04-14T11:00:00Z",
			"reserved_by": {
				"name": "sahla",
				"phone": "0779497842",
				"email": "sahla@gmail.com"
			}
		}
	`
}

func (r ReserveSlot) Validate() (entities.ReservedSlots, error) {

	reserveSlot := entities.ReservedSlots{}

	reserveSlot.QueueID = r.QueueID
	reserveSlot.StartTime = r.StartTime
	reserveSlot.EndTime = r.EndTime
	reserveSlot.ReservedBy.Name = r.ReservedBy.Name
	reserveSlot.ReservedBy.Phone = r.ReservedBy.Phone
	reserveSlot.ReservedBy.Email = r.ReservedBy.Email

	if len(reserveSlot.ReservedBy.Phone) != 10 {
		return entities.ReservedSlots{}, errors.New("invalid phone number")
	}

	return reserveSlot, nil
}
