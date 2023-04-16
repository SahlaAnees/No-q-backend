package decoders

import (
	"no-q-solution/domain/entities"
	"time"
)

type Queue struct {
	Name      string    `json:"name" validate:"required"`
	Interval  int       `json:"interval" validate:"required"`
	StartTime time.Time `json:"start_time" validate:"required"`
	EndTime   time.Time `json:"end_time" validate:"required"`
}

func (q Queue) Format() string {
	return `
		{
			"name": "xyz",
			"interval": 30,
			"start_time": "2023-04-14T10:00:00Z",
			"end_time": "2023-04-14T11:00:00Z"
		}
	`
}

func (q Queue) Validate() (entities.Queue, error) {

	queue := entities.Queue{}

	queue.Name = q.Name
	queue.Interval = q.Interval
	queue.StartTime = q.StartTime
	queue.EndTime = q.EndTime

	return queue, nil
}
