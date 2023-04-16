package decoders

import (
	"time"
)

type Dates struct {
	Dates []time.Time `json:"dates" validate:"required"`
}

func (d Dates) Format() string {
	return `
		{
			"email": "merchant@example.com",
			"password": "#xsgJ62J"
		}
	`
}

func (d Dates) Validate() ([]time.Time, error) {

	return d.Dates, nil
}
