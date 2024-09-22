package models_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"filmoteka/internal/app/models"
)

func TestActor_Validate(t *testing.T) {
	testCases := []struct {
		name    string
		a       func() *models.Actor
		isValid bool
	}{
		{
			name: "valid",
			a: func() *models.Actor {
				return models.TestActor(t)
			},
			isValid: true,
		},
		{
			name: "empty name",
			a: func() *models.Actor {
				a := models.TestActor(t)
				a.Name = ""

				return a
			},
			isValid: false,
		},
		{
			name: "empty gender",
			a: func() *models.Actor {
				a := models.TestActor(t)
				a.Gender = ""

				return a
			},
			isValid: false,
		},
		{
			name: "empty birthdate",
			a: func() *models.Actor {
				a := models.TestActor(t)
				a.BirthDate = ""

				return a
			},
			isValid: false,
		},
		{
			name: "invalid birthdate",
			a: func() *models.Actor {
				a := models.TestActor(t)
				// "2000-01-01" required!
				a.BirthDate = "01/01/2000"

				return a
			},
			isValid: false,
		},
		{
			name: "short name",
			a: func() *models.Actor {
				a := models.TestActor(t)
				a.Name = "No"

				return a
			},
			isValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.isValid {
				assert.NoError(t, tc.a().Validate())
			} else {
				assert.Error(t, tc.a().Validate())
			}
		})
	}
}
