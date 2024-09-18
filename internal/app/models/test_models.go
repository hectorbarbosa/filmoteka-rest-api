package models

import "testing"

// TestFilm ...
func TestFilm(t *testing.T) *Film {
	t.Helper()

	return &Film{
		Name:        "Film Test 1",
		Description: "Description 1",
		ReleaseYear: 2002,
		Rating:      7.5,
	}
}

// TestActor ...
func TestActor(t *testing.T) *Actor {
	t.Helper()

	return &Actor{
		Name:      "First Actor",
		Gender:    "M",
		BirthDate: "1995-02-07",
	}
}
