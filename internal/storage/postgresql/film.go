package postgresql

import (
	"database/sql"
	"strings"

	"fimoteka/internal"
	"fimoteka/internal/app/models"
	"fimoteka/internal/restapi"
)

// Film represents the repository used for interacting with Film records.
type FilmRepository struct {
	db *sql.DB
}

// NewFilm instantiates the Film repository.
func NewFilm(db *sql.DB) *FilmRepository {
	return &FilmRepository{
		db: db,
	}
}

// Create inserts a new Film record.
func (r *FilmRepository) Create(f restapi.CreateFilm) (models.Film, error) {
	var id int
	if err := r.db.QueryRow(
		"INSERT INTO films (name, description, release_year, rating) VALUES ($1, $2, $3, $4) RETURNING id;",
		f.Name,
		f.Description,
		f.ReleaseYear,
		f.Rating,
	).Scan(&id); err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			return models.Film{}, internal.WrapErrorf(err, internal.ErrorCodeUniqueConstraints, "insert film")
		} else {
			return models.Film{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "insert film")
		}
	}

	return models.Film{
		Id:          id,
		Name:        f.Name,
		Description: f.Description,
		ReleaseYear: f.ReleaseYear,
		Rating:      f.Rating,
	}, nil
}

// Delete deletes the existing record matching the id.
func (r *FilmRepository) Delete(id string) error {
	result, err := r.db.Exec("DELETE FROM films WHERE id=$1;", id)
	if err != nil {
		return internal.WrapErrorf(err, internal.ErrorCodeUnknown, "delete film")
	}

	deletedRows, err := result.RowsAffected()
	if err != nil {
		return internal.WrapErrorf(err, internal.ErrorCodeUnknown, "delete film")
	}
	if deletedRows == 0 {
		return internal.WrapErrorf(err, internal.ErrorCodeNotFound, "delete film")
	}

	return nil
}

func (r *FilmRepository) SearchBy() ([]models.Film, error) {
	f := &models.Film{}
	films := make([]models.Film, 0)
	rows, err := r.db.Query(
		"SELECT id, name, description, release_year, rating FROM films;")
	if err != nil {
		return nil, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "search by")
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&f.Id,
			&f.Name,
			&f.Description,
			&f.ReleaseYear,
			&f.Rating,
		)
		if err != nil {
			return nil, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "search by")
		}
		films = append(films, *f)
	}

	return films, nil
}

func (r *FilmRepository) Find(id string) (models.Film, error) {
	f := models.Film{}
	if err := r.db.QueryRow(
		"SELECT id, name, description, release_year, rating FROM films WHERE id=$1",
		id,
	).Scan(
		&f.Id,
		&f.Name,
		&f.Description,
		&f.ReleaseYear,
		&f.Rating,
	); err != nil {
		switch err {
		case sql.ErrNoRows:
			return models.Film{}, internal.WrapErrorf(err, internal.ErrorCodeNotFound, "find film")
		default:
			return models.Film{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "find film")
		}
	}

	return f, nil
}

func (r *FilmRepository) Update(id string, f restapi.UpdateFilm) error {
	result, err := r.db.Exec(
		"UPDATE films SET name=$1, description=$2, release_year=$3, rating=$4 WHERE id=$5;",
		f.Name,
		f.Description,
		f.ReleaseYear,
		f.Rating,
		id,
	)
	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			return internal.WrapErrorf(err, internal.ErrorCodeUniqueConstraints, "update film")
		} else {
			return internal.WrapErrorf(err, internal.ErrorCodeUnknown, "update film")
		}
	}

	updatedRows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if updatedRows == 0 {
		return internal.WrapErrorf(err, internal.ErrorCodeNotFound, "update film")
	}

	return nil
}
