package rest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	_ "github.com/AngelicaNice/HollywoodStarsCRUD/docs"
	"github.com/AngelicaNice/HollywoodStarsCRUD/internal/domain"
	"github.com/AngelicaNice/TTLCache/cache"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
)

var Cache = cache.New()

const TTL = 15 * time.Second

type Actors interface {
	Create(ctx context.Context, actor domain.Actor) (int64, error)
	GetByID(ctx context.Context, id int64) (domain.Actor, error)
	GetAllActors(ctx context.Context) ([]domain.Actor, error)
	Update(ctx context.Context, id int64, info domain.UpdateActorInfo) error
	Delete(ctx context.Context, id int64) error
}

type Handler struct {
	actorsService Actors
}

func NewHandler(a Actors) *Handler {
	return &Handler{
		actorsService: a,
	}
}

func (h *Handler) InitRouter() *mux.Router {
	r := mux.NewRouter()
	r.Use(loggingMiddleware)

	api := r.PathPrefix("/actors").Subrouter()
	{
		a := api.HandleFunc("/id/{id:[0-9]+}", h.GetActor).Methods(http.MethodGet)
		s, _ := a.GetPathTemplate()
		fmt.Println(s)
		a = api.HandleFunc("/id/{id:[0-9]+}", h.UpdateActor).Methods(http.MethodPut)
		s, _ = a.GetPathRegexp()
		fmt.Println(s)
		a = api.HandleFunc("/id/{id:[0-9]+}", h.DeleteActor).Methods(http.MethodDelete)
		s, _ = a.GetPathRegexp()
		fmt.Println(s)
		a = api.HandleFunc("", h.AddActor).Methods(http.MethodPost)
		s, _ = a.GetPathRegexp()
		fmt.Println(s)
		a = api.HandleFunc("", h.GetAllActors).Methods(http.MethodGet)
		s, _ = a.GetPathRegexp()
		fmt.Println(s)
	}
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	return r
}

// Auth godoc
//
//	@Summary		Add actor
//	@Description	add actor info
//	@Tags			actor
//	@Accept			json
//	@Produce		json
//	@Param			input body domain.Actor true "actor's info"
//	@Success		201	{integer} integer 1
//	@Failure		400,404,500 {integer} integer 0
//	@Router			/ [post]
func (h *Handler) AddActor(w http.ResponseWriter, r *http.Request) {
	var actor domain.Actor

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&actor); err != nil {
		log.WithFields(log.Fields{
			"handler": "AddActor",
			"issue":   "failed unmarshalling request body",
		}).Error(err)
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	id, err := h.actorsService.Create(context.TODO(), actor)
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "AddActor",
			"issue":   "internal error",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusCreated)

	Cache.Set(fmt.Sprintf("%d", id), actor, TTL)
}

// Auth godoc
//
//	@Summary		Get all actors
//	@Description	get all actors info
//	@Tags			actor
//	@Accept			json
//	@Produce		json
//	@Success		200	{integer} integer 1
//	@Failure		400,404,500 {integer} integer 0
//	@Router			/ [get]
func (h *Handler) GetAllActors(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetAllActors")

	actors, err := h.actorsService.GetAllActors(context.TODO())
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "GetAllActors",
			"issue":   "internal error",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	response, err := json.Marshal(actors)
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "GetAllActors",
			"issue":   "failed marshaling response body",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(response)

	for _, a := range actors {
		Cache.Set(fmt.Sprintf("%d", a.ID), a, TTL)
	}
}

// Auth godoc
//
//	@Summary		Get actor by id
//	@Description	Get actor info by id
//	@Tags			actor
//	@Accept			json
//	@Produce		json
//	@Param			id	query	int	false 	"int valid"	minimum(1)
//	@Success		200	{integer} integer 1
//	@Failure		400,404,500 {integer} integer 0
//	@Router			/id [get]
func (h *Handler) GetActor(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetActor")

	id, err := getIdFromRequest(r)
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "GetActor",
			"issue":   "failed reading request param",
		}).Error(err)
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	actor, err := Cache.Get(fmt.Sprintf("%d", id))

	if err != nil {
		actor, err = h.actorsService.GetByID(context.TODO(), id)
		if err != nil {
			if errors.Is(err, domain.ErrActorNotFound) {
				issue := fmt.Sprintf("actor with id=%d not found", id)
				log.WithFields(log.Fields{
					"handler": "GetActor",
					"issue":   issue,
				}).Error(err)
				w.WriteHeader(http.StatusBadRequest)

				return
			}

			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		Cache.Set(fmt.Sprintf("%d", id), actor, TTL)
	} else {
		fmt.Println("actor got from cache")
	}

	response, err := json.Marshal(actor)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(response)
}

// Auth godoc
//
//	@Summary		Update actor by id
//	@Description	Update actor info by id
//	@Tags			actor
//	@Accept			json
//	@Produce		json
//	@Param			id	query	int	false 	"int valid"	minimum(1)
//	@Param			input body domain.UpdateActorInfo true "new actor's info"
//	@Success		200	{integer} integer 1
//	@Failure		400,404,500 {integer} integer 0
//	@Router			/id [put]
func (h *Handler) UpdateActor(w http.ResponseWriter, r *http.Request) {
	id, err := getIdFromRequest(r)
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "UpdateActor",
			"issue":   "failed reading request param",
		}).Error(err)
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	var src domain.UpdateActorInfo

	decoder := json.NewDecoder(r.Body)
	if err = decoder.Decode(&src); err != nil {
		log.WithFields(log.Fields{
			"handler": "UpdateActor",
			"issue":   "bad request",
		}).Error(err)
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	if err = h.actorsService.Update(context.TODO(), id, src); err != nil {
		if errors.Is(err, domain.ErrActorNotFound) {
			issue := fmt.Sprintf("actor with id=%d not found", id)
			log.WithFields(log.Fields{
				"handler": "UpdateActor",
				"issue":   issue,
			}).Error(err)
			w.WriteHeader(http.StatusBadRequest)

			return
		}

		log.WithFields(log.Fields{
			"handler": "UpdateActor",
			"issue":   "internal error",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)

	if _, err := Cache.Get(fmt.Sprintf("%d", id)); err == nil {
		Cache.Delete(fmt.Sprintf("%d", id))
	}
}

// Auth godoc
//
//	@Summary		Delete actor by id
//	@Description	Delete actor info by id
//	@Tags			actor
//	@Accept			json
//	@Produce		json
//	@Param			id	query	int	false 	"int valid"	minimum(1)
//	@Success		200	{integer} integer 1
//	@Failure		400,404,500 {integer} integer 0
//	@Router			/id [delete]
func (h *Handler) DeleteActor(w http.ResponseWriter, r *http.Request) {
	id, err := getIdFromRequest(r)
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "DeleteActor",
			"issue":   "failed reading request param",
		}).Error(err)
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	err = h.actorsService.Delete(context.TODO(), id)
	if err != nil {
		if errors.Is(err, domain.ErrActorNotFound) {
			issue := fmt.Sprintf("actor with id=%d not found", id)
			log.WithFields(log.Fields{
				"handler": "DeleteActor",
				"issue":   issue,
			}).Error(err)
			w.WriteHeader(http.StatusBadRequest)

			return
		}

		log.WithFields(log.Fields{
			"handler": "DeleteActor",
			"issue":   "internal error",
		}).Error(err)

		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)

	if _, err := Cache.Get(fmt.Sprintf("%d", id)); err == nil {
		Cache.Delete(fmt.Sprintf("%d", id))
	}
}

func getIdFromRequest(r *http.Request) (int64, error) {
	param := r.URL.RawQuery
	param = param[3:]

	id, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return 0, err
	}

	if id == 0 {
		return 0, errors.New("id can't be 0")
	}

	return id, nil
}
