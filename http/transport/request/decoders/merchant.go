package decoders

import (
	"no-q-solution/domain/entities"
)

type Merchant struct {
	Name      string `json:"name" validate:"required"`
	Category  string `json:"category" validate:"required"`
	Email     string `json:"email" validate:"required"`
	Password  string `json:"password" validate:"required"`
	Facebook  string `json:"facebook"`
	Instagram string `json:"instagram"`
	Website   string `json:"website"`
}

func (m Merchant) Format() string {
	return `
		{
			"name": "merchant",
			"category": "medical",
			"email": "merchant@example.com",
			"password": "#xsgJ62J",
			"facebook": "fb.com/merchant",
			"instagram": "insta.com/merchant",
			"website": "merchant.com"
		}
	`
}

func (m Merchant) Validate() (entities.Merchant, error) {

	merchant := entities.Merchant{}

	merchant.Name = m.Name
	merchant.Category = m.Category
	merchant.Email = m.Email
	merchant.Password = m.Password
	merchant.Facebook = m.Facebook
	merchant.Instagram = m.Instagram
	merchant.Website = m.Website

	return merchant, nil
}
