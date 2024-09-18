package restapi

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
)

type CreateFilm struct {
	Name        string  `json:"name" validate:"required,min=2,max=150"`
	Description string  `json:"description" validate:"required,min=5,max=500"`
	ReleaseYear uint16  `json:"release_year" validate:"required,gte=1900,lte=2030"`
	Rating      float32 `json:"rating" validate:"required,gte=0,lte=10"`
}

func (f *CreateFilm) Validate() error {
	validate := validator.New()
	if err := validate.Struct(f); err != nil {
		return err
	}
	return nil
}

type UpdateFilm struct {
	Name        string  `json:"name" validate:"required,min=2,max=150"`
	Description string  `json:"description" validate:"required,min=5,max=500"`
	ReleaseYear uint16  `json:"release_year" validate:"required,gte=1900,lte=2030"`
	Rating      float32 `json:"rating" validate:"required,gte=0,lte=10"`
}

func (f *UpdateFilm) Validate() error {
	validate := validator.New()
	if err := validate.Struct(f); err != nil {
		return err
	}
	return nil
}

type CreateActor struct {
	Name      string `json:"name" validate:"required,min=3,max=100"`
	Gender    string `json:"gender" validate:"required,len=1"`
	BirthDate string `json:"birth_date" validate:"required"`
}

func (a *CreateActor) Validate() error {
	validate := validator.New()
	if err := validate.Struct(a); err != nil {
		return err
	}
	if _, err := time.Parse("2006-01-02", a.BirthDate); err != nil {
		return fmt.Errorf("valid date format 2006-01-02")
	}
	if a.Gender != "M" && a.Gender != "F" {
		return fmt.Errorf("valid gender values: 'M' or 'F'")
	}

	return nil
}

type UpdateActor struct {
	Name      string `json:"name" validate:"required,min=3,max=100"`
	Gender    string `json:"gender" validate:"required,len=1"`
	BirthDate string `json:"birth_date" validate:"required"`
}

func (a *UpdateActor) Validate() error {
	validate := validator.New()
	if err := validate.Struct(a); err != nil {
		return err
	}
	if _, err := time.Parse("2006-01-02", a.BirthDate); err != nil {
		return fmt.Errorf("valid date format 2006-01-02")
	}
	if a.Gender != "M" && a.Gender != "F" {
		return fmt.Errorf("valid gender values: 'M' or 'F'")
	}

	return nil
}
