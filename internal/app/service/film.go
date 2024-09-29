package service

import (
	"context"
	"fmt"

	"filmoteka/internal"
	"filmoteka/internal/app/models"
	m "filmoteka/internal/restapi/models"
)

// FilmRepository defines the datastore handling Film records.
type FilmRepository interface {
	Create(f m.CreateFilm) (models.Film, error)
	Delete(id string) error
	FindAll() ([]models.Film, error)
	Find(id string) (models.Film, error)
	Update(id string, f m.UpdateFilm) error
}

// FilmSearchRepository defines the datastore handling persisting Searchable Film records.
type FilmSearchRepository interface {
	Delete(ctx context.Context, id string) error
	Index(ctx context.Context, film models.Film) error
	Search(
		ctx context.Context,
		name *string,
		description *string,
		releaseYear *uint16,
		rating *float32,
	) ([]models.Film, error)
}

// FilmService defines the application service in charge of interacting with Tasks.
type FilmService struct {
	repo   FilmRepository
	search FilmSearchRepository
}

// NewFilmService
func NewFilmService(repo FilmRepository, search FilmSearchRepository) *FilmService {
	return &FilmService{
		repo:   repo,
		search: search,
	}
}

// Search gets all existing Films from the datastore.
func (s *FilmService) Search(
	ctx context.Context,
	name string,
	description string,
	releaseYear uint16,
	rating float32,
) ([]models.Film, error) {
	films, err := s.search.Search(ctx, &name, &description, &releaseYear, &rating)
	if err != nil {
		return nil, fmt.Errorf("search: %w", err)
	}

	return films, nil
}

// Create stores a new record.
func (s *FilmService) Create(ctx context.Context, f m.CreateFilm) (models.Film, error) {
	if err := f.Validate(); err != nil {
		return models.Film{}, internal.WrapErrorf(err, internal.ErrorCodeInvalidArgument, "validate film")
	}

	film, err := s.repo.Create(f)
	if err != nil {
		return models.Film{}, fmt.Errorf("repo create: %w", err)
	}

	_ = s.search.Index(ctx, film) // Ignoring errors on purpose

	return film, nil
}

// Delete removes an existing Film from the datastore.
func (s *FilmService) Delete(ctx context.Context, id string) error {
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("repo delete: %w", err)
	}

	_ = s.search.Delete(ctx, id) // Ignoring errors on purpose

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

// FindAll gets all existing Films from the datastore.
func (s *FilmService) FindAll() ([]models.Film, error) {
	films, err := s.repo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("repo find: %w", err)
	}

	return films, nil
}

// Update updates an existing Film in the datastore.
func (s *FilmService) Update(ctx context.Context, id string, f m.UpdateFilm) error {
	if err := f.Validate(); err != nil {
		return internal.WrapErrorf(err, internal.ErrorCodeInvalidArgument, "validate film")
	}

	if err := s.repo.Update(id, f); err != nil {
		return fmt.Errorf("repo update: %w", err)
	}

	film, err := s.repo.Find(id)
	if err == nil {
		_ = s.search.Index(ctx, film) // Ignoring errors on purpose
	}
	return nil
}
