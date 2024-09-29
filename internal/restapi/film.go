package restapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"filmoteka/internal"
	"filmoteka/internal/app/models"
	m "filmoteka/internal/restapi/models"
)

// FilmService ...
//
//go:generate mockgen -source=film.go -destination=mock_restapi/mockfilm.go
type FilmService interface {
	// By(args internal.SearchParams) (internal.SearchResults, error)
	Create(ctx context.Context, f m.CreateFilm) (models.Film, error)
	Delete(ctx context.Context, id string) error
	Search(
		ctx context.Context,
		name string,
		description string,
		releaseYear uint16,
		rating float32,
	) ([]models.Film, error)
	FindAll() ([]models.Film, error)
	Find(id string) (models.Film, error)
	Update(ctx context.Context, id string, f m.UpdateFilm) error
}

// FilmHandler ...
type FilmHandler struct {
	svc FilmService
}

// NewFilmHandler ...
func NewFilmHandler(svc FilmService) *FilmHandler {
	return &FilmHandler{
		svc: svc,
	}
}

func (h *FilmHandler) Register(r *mux.Router) {
	r.HandleFunc("/films", h.create).Methods(http.MethodPost)
	r.HandleFunc("/films/search", h.search).Methods(http.MethodGet)
	r.HandleFunc("/films", h.findAll).Methods(http.MethodGet)
	r.HandleFunc("/films/{id}", h.find).Methods(http.MethodGet)
	r.HandleFunc("/films/{id}", h.update).Methods(http.MethodPut)
	r.HandleFunc("/films/{id}", h.delete).Methods(http.MethodDelete)
}

func (h *FilmHandler) create(w http.ResponseWriter, r *http.Request) {
	var req m.CreateFilm
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		e := internal.WrapErrorf(err, internal.ErrorCodeInvalidArgument, "json decoder")
		msg := fmt.Errorf("invalid request %w", e)
		renderErrorResponse(w, msg.Error(), msg)
		return
	}

	defer r.Body.Close()

	film, err := h.svc.Create(r.Context(), req)
	if err != nil {
		// fmt.Println(err)
		msg := fmt.Errorf("create failed: %w", err)
		renderErrorResponse(w, msg.Error(), msg)
		return
	}

	renderResponse(w,
		film,
		http.StatusCreated)
}

func (h *FilmHandler) delete(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"] // NOTE: Safe to ignore error, because it's always defined.

	if err := h.svc.Delete(r.Context(), id); err != nil {
		msg := fmt.Errorf("delete failed: %w", err)
		renderErrorResponse(w, msg.Error(), msg)
		return
	}

	renderResponse(w, struct{}{}, http.StatusOK)
}

// ReadTasksResponse defines the response returned back after searching one task.
type ReadFilmsResponse struct {
	Film models.Film `json:"film"`
}

type SearchTasksRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	ReleaseYear uint16  `json:"release_year"`
	Rating      float32 `json:"rating"`
}

func (h *FilmHandler) search(w http.ResponseWriter, r *http.Request) {
	var req SearchTasksRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderErrorResponse(
			w,
			"invalid request",
			internal.WrapErrorf(err, internal.ErrorCodeInvalidArgument, "json decoder"))
		return
	}

	defer r.Body.Close()

	films, err := h.svc.Search(
		r.Context(),
		req.Name,
		req.Description,
		req.ReleaseYear,
		req.Rating,
	)
	if err != nil {
		msg := fmt.Errorf("search failed: %w", err)
		renderErrorResponse(w, msg.Error(), msg)
		return
	}

	renderResponse(w,
		films,
		http.StatusOK)
}

func (h *FilmHandler) findAll(w http.ResponseWriter, r *http.Request) {
	film, err := h.svc.FindAll()
	if err != nil {
		msg := fmt.Errorf("find failed: %w", err)
		renderErrorResponse(w, msg.Error(), msg)
		return
	}

	renderResponse(w,
		film,
		http.StatusOK)
}

func (h *FilmHandler) find(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"] // NOTE: Safe to ignore error, because it's always defined.

	film, err := h.svc.Find(id)
	if err != nil {
		msg := fmt.Errorf("find failed: %w", err)
		renderErrorResponse(w, msg.Error(), msg)
		return
	}

	renderResponse(w,
		film,
		http.StatusOK)
}

func (h *FilmHandler) update(w http.ResponseWriter, r *http.Request) {
	var req m.UpdateFilm
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		e := internal.WrapErrorf(err, internal.ErrorCodeInvalidArgument, "json decoder")
		msg := fmt.Errorf("invalid request %w", e)
		renderErrorResponse(w, msg.Error(), msg)
		return
	}

	defer r.Body.Close()

	id := mux.Vars(r)["id"] // NOTE: Safe to ignore error, because it's always defined.

	err := h.svc.Update(r.Context(), id, req)
	if err != nil {
		msg := fmt.Errorf("update failed: %w", err)
		renderErrorResponse(w, msg.Error(), msg)
		return
	}

	renderResponse(w, &struct{}{}, http.StatusOK)
}
