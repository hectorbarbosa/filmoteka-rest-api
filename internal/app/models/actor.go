package models

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
)

type Actor struct {
	Id        int    `json:"id"`
	Name      string `json:"name" validate:"required,min=3,max=100"`
	Gender    string `json:"gender" validate:"required,len=1"`
	BirthDate string `json:"birth_date" validate:"required"`
}

func (a *Actor) Validate() error {
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
