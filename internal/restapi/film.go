package restapi

import (
	"encoding/json"
	"fimoteka/internal"
	"fimoteka/internal/app/models"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// FilmService ...
type FilmService interface {
	// By(args internal.SearchParams) (internal.SearchResults, error)
	Create(f CreateFilm) (models.Film, error)
	Delete(id string) error
	Search() ([]models.Film, error)
	Find(id string) (models.Film, error)
	Update(id string, f UpdateFilm) error
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
	r.HandleFunc("/films", h.search).Methods(http.MethodGet)
	r.HandleFunc("/films/{id}", h.find).Methods(http.MethodGet)
	r.HandleFunc("/films/{id}", h.update).Methods(http.MethodPut)
	r.HandleFunc("/films/{id}", h.delete).Methods(http.MethodDelete)
}

func (h *FilmHandler) create(w http.ResponseWriter, r *http.Request) {
	var req CreateFilm
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		e := internal.WrapErrorf(err, internal.ErrorCodeInvalidArgument, "json decoder")
		msg := fmt.Errorf("invalid request %w", e)
		renderErrorResponse(w, msg.Error(), msg)
		return
	}

	defer r.Body.Close()

	film, err := h.svc.Create(req)
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

	if err := h.svc.Delete(id); err != nil {
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

func (h *FilmHandler) search(w http.ResponseWriter, r *http.Request) {
	films, err := h.svc.Search()
	if err != nil {
		msg := fmt.Errorf("search failed: %w", err)
		renderErrorResponse(w, msg.Error(), msg)
		return
	}

	renderResponse(w,
		films,
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
	var req UpdateFilm
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		e := internal.WrapErrorf(err, internal.ErrorCodeInvalidArgument, "json decoder")
		msg := fmt.Errorf("invalid request %w", e)
		renderErrorResponse(w, msg.Error(), msg)
		return
	}

	defer r.Body.Close()

	id := mux.Vars(r)["id"] // NOTE: Safe to ignore error, because it's always defined.

	err := h.svc.Update(id, req)
	if err != nil {
		msg := fmt.Errorf("update failed: %w", err)
		renderErrorResponse(w, msg.Error(), msg)
		return
	}

	renderResponse(w, &struct{}{}, http.StatusOK)
}
