package rest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/AngelicaNice/HollywoodStarsCRUD/internal/domain"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// Auth godoc
//
//	@Summary		Add actor
//	@Security 		ApiKeyAuth
//	@Description	add actor info
//	@Tags			actor
//	@Accept			json
//	@Produce		json
//	@Param			input body domain.ActorInput true "actor's info"
//	@Success		201	{integer} integer 1
//	@Failure		400,404,500 {integer} integer 0
//	@Router			/actors [post]
func (h *Handler) AddActor(c *gin.Context) {
	var actor domain.ActorInput

	decoder := json.NewDecoder(c.Request.Body)
	if err := decoder.Decode(&actor); err != nil {
		log.WithFields(log.Fields{
			"handler": "AddActor",
			"issue":   "failed unmarshalling request body",
		}).Error(err)
		c.Writer.WriteHeader(http.StatusBadRequest)

		return
	}

	if _, err := h.actorsService.Create(context.TODO(), actor); err != nil {
		log.WithFields(log.Fields{
			"handler": "AddActor",
			"issue":   "internal error",
		}).Error(err)
		c.Writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	c.Writer.WriteHeader(http.StatusCreated)
}

// Auth godoc
//
//	@Summary		Get all actors
//	@Security 		ApiKeyAuth
//	@Description	get all actors info
//	@Tags			actor
//	@Accept			json
//	@Produce		json
//	@Success		200	{integer} integer 1
//	@Failure		400,404,500 {integer} integer 0
//	@Router			/actors [get]
func (h *Handler) GetAllActors(c *gin.Context) {
	actors, err := h.actorsService.GetAllActors(context.TODO())
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "GetAllActors",
			"issue":   "internal error",
		}).Error(err)
		c.Writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	encoder := json.NewEncoder(c.Writer)

	if err := encoder.Encode(&actors); err != nil {
		log.WithFields(log.Fields{
			"handler": "GetAllActors",
			"issue":   "failed marshaling response body",
		}).Error(err)
		c.Writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	c.Writer.Header().Add("Content-Type", "application/json")
}

// Auth godoc
//
//	@Summary		Get actor by id
//	@Security 		ApiKeyAuth
//	@Description	Get actor info by id
//	@Tags			actor
//	@Accept			json
//	@Produce		json
//	@Param			id	query	int	false 	"int valid"	minimum(1)
//	@Success		200	{integer} integer 1
//	@Failure		400,404,500 {integer} integer 0
//	@Router			/actors/id [get]
func (h *Handler) GetActor(c *gin.Context) {
	id, err := getIdFromRequest(c.Request)
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "GetActor",
			"issue":   "failed reading request param",
		}).Error(err)
		c.Writer.WriteHeader(http.StatusBadRequest)

		return
	}

	actor, err := h.actorsService.GetByID(context.TODO(), id)
	if err != nil {
		if errors.Is(err, domain.ErrActorNotFound) {
			issue := fmt.Sprintf("actor with id=%d not found", id)
			log.WithFields(log.Fields{
				"handler": "GetActor",
				"issue":   issue,
			}).Error(err)
			c.Writer.WriteHeader(http.StatusBadRequest)

			return
		}

		c.Writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	encoder := json.NewEncoder(c.Writer)

	if err := encoder.Encode(&actor); err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	c.Writer.Header().Add("Content-Type", "application/json")
}

// Auth godoc
//
//	@Summary		Update actor by id
//	@Security 		ApiKeyAuth
//	@Description	Update actor info by id
//	@Tags			actor
//	@Accept			json
//	@Produce		json
//	@Param			id	query	int	false 	"int valid"	minimum(1)
//	@Param			input body domain.UpdateActorInfo true "new actor's info"
//	@Success		200	{integer} integer 1
//	@Failure		400,404,500 {integer} integer 0
//	@Router			/actors/id [put]
func (h *Handler) UpdateActor(c *gin.Context) {
	id, err := getIdFromRequest(c.Request)
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "UpdateActor",
			"issue":   "failed reading request param",
		}).Error(err)
		c.Writer.WriteHeader(http.StatusBadRequest)

		return
	}

	var src domain.UpdateActorInfo

	decoder := json.NewDecoder(c.Request.Body)
	if err := decoder.Decode(&src); err != nil {
		log.WithFields(log.Fields{
			"handler": "UpdateActor",
			"issue":   "bad request",
		}).Error(err)
		c.Writer.WriteHeader(http.StatusBadRequest)

		return
	}

	if err = h.actorsService.Update(context.TODO(), id, src); err != nil {
		if errors.Is(err, domain.ErrActorNotFound) {
			issue := fmt.Sprintf("actor with id=%d not found", id)
			log.WithFields(log.Fields{
				"handler": "UpdateActor",
				"issue":   issue,
			}).Error(err)
			c.Writer.WriteHeader(http.StatusBadRequest)

			return
		}

		log.WithFields(log.Fields{
			"handler": "UpdateActor",
			"issue":   "internal error",
		}).Error(err)
		c.Writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	c.Writer.WriteHeader(http.StatusOK)
}

// Auth godoc
//
//	@Summary		Delete actor by id
//	@Security 		ApiKeyAuth
//	@Description	Delete actor info by id
//	@Tags			actor
//	@Accept			json
//	@Produce		json
//	@Param			id	query	int	false 	"int valid"	minimum(1)
//	@Success		200	{integer} integer 1
//	@Failure		400,404,500 {integer} integer 0
//	@Router			/actors/id [delete]
func (h *Handler) DeleteActor(c *gin.Context) {
	id, err := getIdFromRequest(c.Request)
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "DeleteActor",
			"issue":   "failed reading request param",
		}).Error(err)
		c.Writer.WriteHeader(http.StatusBadRequest)

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
			c.Writer.WriteHeader(http.StatusBadRequest)

			return
		}

		log.WithFields(log.Fields{
			"handler": "DeleteActor",
			"issue":   "internal error",
		}).Error(err)
		c.Writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	c.Writer.WriteHeader(http.StatusOK)
}

func getIdFromRequest(r *http.Request) (int64, error) {
	param := r.URL.RawQuery
	param = param[3:]

	id, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return 0, err
	}

	if id == 0 {
		return 0, errors.New("actor id can't be 0")
	}

	return id, nil
}
