package service

import (
	"fmt"

	"filmoteka/internal"
	"filmoteka/internal/app/models"
	m "filmoteka/internal/restapi/models"
)

// FilmRepository defines the datastore handling Film records.
type FilmRepository interface {
	Create(f m.CreateFilm) (models.Film, error)
	Delete(id string) error
	SearchBy() ([]models.Film, error)
	Find(id string) (models.Film, error)
	Update(id string, f m.UpdateFilm) error
}

// FilmService defines the application service in charge of interacting with Tasks.
type FilmService struct {
	repo FilmRepository
}

// NewFilmService
func NewFilmService(repo FilmRepository) *FilmService {
	return &FilmService{
		repo: repo,
	}
}

// Create stores a new record.
func (s *FilmService) Create(f m.CreateFilm) (models.Film, error) {
	if err := f.Validate(); err != nil {
		return models.Film{}, internal.WrapErrorf(err, internal.ErrorCodeInvalidArgument, "validate film")
	}

	film, err := s.repo.Create(f)
	if err != nil {
		return models.Film{}, fmt.Errorf("repo create: %w", err)
	}

	return film, nil
}

// Delete removes an existing Film from the datastore.
func (s *FilmService) Delete(id string) error {
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("repo delete: %w", err)
	}

	return nil
}

// Find gets an existing Film from the datastore.
func (s *FilmService) Find(id string) (models.Film, error) {
	task, err := s.repo.Find(id)
	if err != nil {
		return models.Film{}, fmt.Errorf("repo find: %w", err)
	}

	return task, nil
}

// Search gets all existing Films from the datastore.
func (s *FilmService) Search() ([]models.Film, error) {
	films, err := s.repo.SearchBy()
	if err != nil {
		return nil, fmt.Errorf("repo find: %w", err)
	}

	return films, nil
}

// Update updates an existing Film in the datastore.
func (s *FilmService) Update(id string, f m.UpdateFilm) error {
	if err := f.Validate(); err != nil {
		return internal.WrapErrorf(err, internal.ErrorCodeInvalidArgument, "validate film")
	}

	if err := s.repo.Update(id, f); err != nil {
		return fmt.Errorf("repo update: %w", err)
	}

	return nil
}
