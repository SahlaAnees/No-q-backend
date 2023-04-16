package decoders

import "no-q-solution/domain/entities"

type Login struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (l Login) Format() string {
	return `
		{
			"email": "merchant@example.com",
			"password": "#xsgJ62J"
		}
	`
}

func (l Login) Validate() (entities.Login, error) {

	login := entities.Login{}

	login.Email = l.Email
	login.Password = l.Password

	return login, nil
}
