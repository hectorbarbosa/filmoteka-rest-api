package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	esv7 "github.com/elastic/go-elasticsearch/v7"
	esv7api "github.com/elastic/go-elasticsearch/v7/esapi"
	"go.opentelemetry.io/otel/trace"

	"filmoteka/internal"
	"filmoteka/internal/app/models"
)

// Film represents the repository used for interacting with Film records.
type FilmSearchRepo struct {
	client *esv7.Client
	index  string
}

type indexedFilm struct {
	Id          int     `json:"id"`
	Name        string  `json:"name" validate:"required,min=2,max=150"`
	Description string  `json:"description" validate:"required,min=5,max=500"`
	ReleaseYear uint16  `json:"release_year" validate:"required,gte=1900,lte=2030"`
	Rating      float32 `json:"rating" validate:"required,gte=0,lte=10"`
}

// NewFilmSerchRepo instantiates the FilmSearchRepo repository.
func NewFilmSearchRepo(client *esv7.Client) *FilmSearchRepo {
	return &FilmSearchRepo{
		client: client,
		index:  "films",
	}
}

// Index creates or updates a film in an index.
func (f *FilmSearchRepo) Index(ctx context.Context, film models.Film) error {
	span := trace.SpanFromContext(ctx)
	defer span.End()

	body := indexedFilm{
		Id:          film.Id,
		Name:        film.Name,
		Description: film.Description,
		ReleaseYear: film.ReleaseYear,
		Rating:      film.Rating,
	}

	var buf bytes.Buffer

	if err := json.NewEncoder(&buf).Encode(body); err != nil {
		return internal.WrapErrorf(err, internal.ErrorCodeUnknown, "json.NewEncoder.Encode")
	}

	req := esv7api.IndexRequest{
		Index:      f.index,
		Body:       &buf,
		DocumentID: strconv.Itoa(film.Id),
		Refresh:    "true",
	}

	resp, err := req.Do(ctx, f.client)
	if err != nil {
		return internal.WrapErrorf(err, internal.ErrorCodeUnknown, "IndexRequest.Do")
	}
	defer resp.Body.Close()

	if resp.IsError() {
		return internal.NewErrorf(internal.ErrorCodeUnknown, "IndexRequest.Do %d", resp.StatusCode)
	}

	io.Copy(io.Discard, resp.Body)

	return nil
}

// Delete removes a film from the index.
func (t *FilmSearchRepo) Delete(ctx context.Context, id string) error {
	span := trace.SpanFromContext(ctx)
	defer span.End()

	req := esv7api.DeleteRequest{
		Index:      t.index,
		DocumentID: id,
	}

	resp, err := req.Do(ctx, t.client)
	if err != nil {
		return internal.WrapErrorf(err, internal.ErrorCodeUnknown, "DeleteRequest.Do")
	}
	defer resp.Body.Close()

	if resp.IsError() {
		return internal.NewErrorf(internal.ErrorCodeUnknown, "DeleteRequest.Do %d", resp.StatusCode)
	}

	io.Copy(io.Discard, resp.Body)

	return nil
}

// Search returns films matching a query.
func (t *FilmSearchRepo) Search(
	ctx context.Context,
	name *string,
	description *string,
	releaseYear *uint16,
	rating *float32,
) ([]models.Film, error) {
	span := trace.SpanFromContext(ctx)
	defer span.End()

	if name == nil && description == nil && releaseYear == nil && rating == nil {
		return nil, nil
	}

	should := make([]interface{}, 0, 3)

	if name != nil {
		should = append(should, map[string]interface{}{
			"match": map[string]interface{}{
				"name": *name,
			},
		})
	}

	if description != nil {
		should = append(should, map[string]interface{}{
			"match": map[string]interface{}{
				"description": *description,
			},
		})
	}

	if releaseYear != nil {
		should = append(should, map[string]interface{}{
			"match": map[string]interface{}{
				"release_year": *releaseYear,
			},
		})
	}

	if rating != nil {
		should = append(should, map[string]interface{}{
			"match": map[string]interface{}{
				"rating": *rating,
			},
		})
	}

	var query map[string]interface{}

	if len(should) > 1 {
		query = map[string]interface{}{
			"query": map[string]interface{}{
				"bool": map[string]interface{}{
					"should": should,
				},
			},
		}
	} else {
		query = map[string]interface{}{
			"query": should[0],
		}
	}

	fmt.Printf("%#v\n", query)

	var buf bytes.Buffer

	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "json.NewEncoder.Encode")
	}

	req := esv7api.SearchRequest{
		Index: []string{t.index},
		Body:  &buf,
	}

	resp, err := req.Do(ctx, t.client)
	if err != nil {
		return nil, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "SearchRequest.Do")
	}
	defer resp.Body.Close()

	if resp.IsError() {
		return nil, internal.NewErrorf(internal.ErrorCodeUnknown, "SearchRequest.Do %d", resp.StatusCode)
	}

	var hits struct {
		Hits struct {
			Hits []struct {
				Source indexedFilm `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&hits); err != nil {
		fmt.Println("Error here", err)
		return nil, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "json.NewDecoder.Decode")
	}

	res := make([]models.Film, len(hits.Hits.Hits))

	for i, hit := range hits.Hits.Hits {
		res[i].Id = hit.Source.Id
		res[i].Name = hit.Source.Name
		res[i].Description = hit.Source.Description
		res[i].ReleaseYear = hit.Source.ReleaseYear
		res[i].Rating = hit.Source.Rating
	}

	return res, nil
}
