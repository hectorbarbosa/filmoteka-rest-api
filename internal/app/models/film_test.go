package models_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"fimoteka/internal/app/models"
)

func TestFilm_Validate(t *testing.T) {
	testCases := []struct {
		name    string
		f       func() *models.Film
		isValid bool
	}{
		{
			name: "valid",
			f: func() *models.Film {
				return models.TestFilm(t)
			},
			isValid: true,
		},
		{
			name: "empty name",
			f: func() *models.Film {
				f := models.TestFilm(t)
				f.Name = ""

				return f
			},
			isValid: false,
		},
		{
			name: "empty description",
			f: func() *models.Film {
				f := models.TestFilm(t)
				f.Description = ""

				return f
			},
			isValid: false,
		},
		{
			name: "empty Release Year",
			f: func() *models.Film {
				f := models.TestFilm(t)
				f.ReleaseYear = 0

				return f
			},
			isValid: false,
		},
		{
			name: "negative Rating",
			f: func() *models.Film {
				f := models.TestFilm(t)
				f.Rating = -2.0

				return f
			},
			isValid: false,
		},
		{
			name: "short name",
			f: func() *models.Film {
				f := models.TestFilm(t)
				f.Name = "A"

				return f
			},
			isValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.isValid {
				assert.NoError(t, tc.f().Validate())
			} else {
				assert.Error(t, tc.f().Validate())
			}
		})
	}
}
