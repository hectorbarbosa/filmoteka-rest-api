package service

import (
	"fmt"

	"fimoteka/internal"
	"fimoteka/internal/app/models"
	"fimoteka/internal/restapi"
)

// ActorRepository defines the datastore handling Actor records.
type ActorRepository interface {
	Create(f restapi.CreateActor) (models.Actor, error)
	Delete(id string) error
	SearchBy() ([]models.Actor, error)
	Find(id string) (models.Actor, error)
	Update(id string, f restapi.UpdateActor) error
}

// Task defines the application service in charge of interacting with Tasks.
type ActorService struct {
	repo ActorRepository
}

// NewActorService
func NewActorService(repo ActorRepository) *ActorService {
	return &ActorService{
		repo: repo,
	}
}

// Create stores a new record.
func (s *ActorService) Create(a restapi.CreateActor) (models.Actor, error) {
	if err := a.Validate(); err != nil {
		return models.Actor{}, internal.WrapErrorf(err, internal.ErrorCodeInvalidArgument, "validate actor")
	}

	actor, err := s.repo.Create(a)
	if err != nil {
		return models.Actor{}, fmt.Errorf("repo create: %w", err)
	}

	return actor, nil
}

// Delete removes an existing Actor from the datastore.
func (s *ActorService) Delete(id string) error {
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("repo delete: %w", err)
	}

	return nil
}

// Find gets an existing Actor from the datastore.
func (s *ActorService) Find(id string) (models.Actor, error) {
	task, err := s.repo.Find(id)
	if err != nil {
		return models.Actor{}, fmt.Errorf("repo find: %w", err)
	}

	return task, nil
}

// Search gets all existing Actors from the datastore.
func (s *ActorService) Search() ([]models.Actor, error) {
	actors, err := s.repo.SearchBy()
	if err != nil {
		return nil, fmt.Errorf("repo find: %w", err)
	}

	return actors, nil
}

// Update updates an existing Actor in the datastore.
func (s *ActorService) Update(id string, a restapi.UpdateActor) error {
	if err := a.Validate(); err != nil {
		return internal.WrapErrorf(err, internal.ErrorCodeInvalidArgument, "validate actor")
	}

	if err := s.repo.Update(id, a); err != nil {
		return fmt.Errorf("repo update: %w", err)
	}

	return nil
}
