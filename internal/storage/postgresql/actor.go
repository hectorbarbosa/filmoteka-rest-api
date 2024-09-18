package postgresql

import (
	"database/sql"
	"strings"

	"fimoteka/internal"
	"fimoteka/internal/app/models"
	"fimoteka/internal/restapi"
)

// Film represents the repository used for interacting with Film records.
type ActorRepository struct {
	db *sql.DB
}

// NewFilm instantiates the Film repository.
func NewActor(db *sql.DB) *ActorRepository {
	return &ActorRepository{
		db: db,
	}
}

// Create inserts a new Actor record.
func (r *ActorRepository) Create(a restapi.CreateActor) (models.Actor, error) {
	var id int
	if err := r.db.QueryRow(
		"INSERT INTO actors (name, gender, birth_date) VALUES ($1, $2, $3) RETURNING id;",
		a.Name,
		a.Gender,
		a.BirthDate,
	).Scan(&id); err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			return models.Actor{}, internal.WrapErrorf(err, internal.ErrorCodeUniqueConstraints, "insert actor")
		} else {
			return models.Actor{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "insert actor")
		}
	}

	return models.Actor{
		Id:        id,
		Name:      a.Name,
		Gender:    a.Gender,
		BirthDate: a.BirthDate,
	}, nil
}

// Delete deletes the existing record matching the id.
func (r *ActorRepository) Delete(id string) error {
	result, err := r.db.Exec("DELETE FROM actors WHERE id=$1;", id)
	if err != nil {
		return internal.WrapErrorf(err, internal.ErrorCodeUnknown, "delete actor")
	}

	deletedRows, err := result.RowsAffected()
	if err != nil {
		return internal.WrapErrorf(err, internal.ErrorCodeUnknown, "delete actor")
	}
	if deletedRows == 0 {
		return internal.WrapErrorf(err, internal.ErrorCodeNotFound, "delete actor")
	}

	return nil
}

func (r *ActorRepository) SearchBy() ([]models.Actor, error) {
	a := models.Actor{}
	actors := make([]models.Actor, 0)
	rows, err := r.db.Query(
		"SELECT id, name, gender, birth_date FROM actors;")
	if err != nil {
		return nil, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "search by")
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&a.Id,
			&a.Name,
			&a.Gender,
			&a.BirthDate,
		)
		if err != nil {
			return nil, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "search by")
		}
		actors = append(actors, a)
	}

	return actors, nil
}

func (r *ActorRepository) Find(id string) (models.Actor, error) {
	a := models.Actor{}
	if err := r.db.QueryRow(
		"SELECT id, name, gender, birth_date rating FROM actors WHERE id=$1",
		id,
	).Scan(
		&a.Id,
		&a.Name,
		&a.Gender,
		&a.BirthDate,
	); err != nil {
		switch err {
		case sql.ErrNoRows:
			return models.Actor{}, internal.WrapErrorf(err, internal.ErrorCodeNotFound, "find actor")
		default:
			return models.Actor{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "find actor")
		}
	}

	return a, nil
}

func (r *ActorRepository) Update(id string, a restapi.UpdateActor) error {
	result, err := r.db.Exec(
		"UPDATE actors SET name=$1, gender=$2, birth_date=$3 WHERE id=$4;",
		a.Name,
		a.Gender,
		a.BirthDate,
		id,
	)
	if err != nil {
		return internal.WrapErrorf(err, internal.ErrorCodeUnknown, "update actor")
	}

	updatedRows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if updatedRows == 0 {
		return internal.WrapErrorf(err, internal.ErrorCodeNotFound, "update actor")
	}

	return nil
}
