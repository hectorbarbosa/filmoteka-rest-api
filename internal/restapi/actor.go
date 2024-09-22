package restapi

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"filmoteka/internal"
	"filmoteka/internal/app/models"
	m "filmoteka/internal/restapi/models"
)

//go:generate mockgen -source=actor.go -destination=mock_restapi/mock.go

// ActorService
type ActorService interface {
	Create(a m.CreateActor) (models.Actor, error)
	Delete(id string) error
	Search() ([]models.Actor, error)
	Find(id string) (models.Actor, error)
	Update(id string, a m.UpdateActor) error
}

// ActorHandler
type ActorHandler struct {
	svc ActorService
}

// NewActorHandler ...
func NewActorHandler(svc ActorService) *ActorHandler {
	return &ActorHandler{
		svc: svc,
	}
}

func (h *ActorHandler) Register(r *mux.Router) {
	r.HandleFunc("/actors", h.create).Methods(http.MethodPost)
	r.HandleFunc("/actors", h.search).Methods(http.MethodGet)
	r.HandleFunc("/actors/{id}", h.find).Methods(http.MethodGet)
	r.HandleFunc("/actors/{id}", h.update).Methods(http.MethodPut)
	r.HandleFunc("/actors/{id}", h.delete).Methods(http.MethodDelete)
}

func (h *ActorHandler) create(w http.ResponseWriter, r *http.Request) {
	var req m.CreateActor
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		e := internal.WrapErrorf(err, internal.ErrorCodeInvalidArgument, "json decoder")
		msg := fmt.Errorf("invalid request %w", e)
		renderErrorResponse(w, msg.Error(), msg)
		return
	}

	defer r.Body.Close()

	actor, err := h.svc.Create(req)
	if err != nil {
		// fmt.Println(err)
		msg := fmt.Errorf("create failed: %w", err)
		renderErrorResponse(w, msg.Error(), msg)
		return
	}

	renderResponse(w,
		actor,
		http.StatusCreated)
}

func (h *ActorHandler) delete(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"] // NOTE: Safe to ignore error, because it's always defined.

	if err := h.svc.Delete(id); err != nil {
		msg := fmt.Errorf("delete failed: %w", err)
		renderErrorResponse(w, msg.Error(), msg)
		return
	}

	renderResponse(w, struct{}{}, http.StatusOK)
}

// ReadTasksResponse defines the response returned back after searching one task.
type ReadActorResponse struct {
	Film models.Actor `json:"actor"`
}

func (h *ActorHandler) search(w http.ResponseWriter, r *http.Request) {
	actors, err := h.svc.Search()
	if err != nil {
		msg := fmt.Errorf("search failed: %w", err)
		renderErrorResponse(w, msg.Error(), msg)
		return
	}

	renderResponse(w,
		actors,
		http.StatusOK)
}

func (h *ActorHandler) find(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"] // NOTE: Safe to ignore error, because it's always defined.

	actor, err := h.svc.Find(id)
	if err != nil {
		msg := fmt.Errorf("find failed: %w", err)
		renderErrorResponse(w, msg.Error(), msg)
		return
	}

	renderResponse(w,
		actor,
		http.StatusOK)
}

func (h *ActorHandler) update(w http.ResponseWriter, r *http.Request) {
	var req m.UpdateActor
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
