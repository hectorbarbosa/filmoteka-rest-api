package models

import (
	"github.com/go-playground/validator/v10"
)

type Film struct {
	Id          int     `json:"id"`
	Name        string  `json:"name" validate:"required,min=2,max=150"`
	Description string  `json:"description" validate:"required,min=5,max=500"`
	ReleaseYear uint16  `json:"release_year" validate:"required,gte=1900,lte=2030"`
	Rating      float32 `json:"rating" validate:"required,gte=0,lte=10"`
}

func (f *Film) Validate() error {
	validate := validator.New()
	if err := validate.Struct(f); err != nil {
		return err
	}
	return nil
}
