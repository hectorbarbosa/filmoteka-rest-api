package postgresql

import (
	"github.com/go-playground/validator/v10"
)

type DbFilm struct {
	Id          string  `validate:"required,min=2,max=150"`
	Name        string  `validate:"required,min=2,max=150"`
	Description string  `validate:"required,min=5,max=500"`
	ReleaseYear uint16  `validate:"required,gte=1900,lte=2030"`
	Rating      float32 `validate:"required,gte=0,lte=10"`
}

func (f *DbFilm) Validate() error {
	validate := validator.New()
	if err := validate.Struct(f); err != nil {
		return err
	}
	return nil
}
