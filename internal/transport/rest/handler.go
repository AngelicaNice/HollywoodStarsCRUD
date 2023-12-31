package rest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	_ "github.com/AngelicaNice/HollywoodStarsCRUD/docs"
	"github.com/AngelicaNice/HollywoodStarsCRUD/internal/domain"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

type Actors interface {
	Create(ctx context.Context, actor domain.Actor) (int64, error)
	GetByID(ctx context.Context, id int64) (domain.Actor, error)
	GetAllActors(ctx context.Context) ([]domain.Actor, error)
	Update(ctx context.Context, id int64, info domain.UpdateActorInfo) error
	Delete(ctx context.Context, id int64) error
}

type Users interface {
	Create(ctx context.Context, user domain.SignUpInput) (int64, error)
	GetByID(ctx context.Context, id int64) (domain.User, error)
}

type Handler struct {
	actorsService Actors
	usersService  Users
}

func NewHandler(a Actors, u Users) *Handler {
	return &Handler{
		actorsService: a,
		usersService:  u,
	}
}

func (h *Handler) InitRouter() *gin.Engine {
	r := gin.Default()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	auth := r.Group("/auth")
	{
		auth.Handle(http.MethodPost, "/sign-up", h.SignUp)
	}

	api := r.Group("/actors")
	{
		api.Handle(http.MethodPost, "", h.AddActor)
		api.Handle(http.MethodGet, "", h.GetAllActors)
		api.Handle(http.MethodGet, "/id", h.GetActor)
		api.Handle(http.MethodPut, "/id", h.UpdateActor)
		api.Handle(http.MethodDelete, "/id", h.DeleteActor)
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}

// Auth godoc
//
//	@Summary		SignUp
//	@Description	Registration in system
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			input body domain.SignUpInput true "user's info"
//	@Success		201	{integer} integer 1
//	@Failure		400,404,500 {integer} integer 0
//	@Router			/auth/sign-up [post]
func (h *Handler) SignUp(c *gin.Context) {
	var user domain.SignUpInput

	decoder := json.NewDecoder(c.Request.Body)
	if err := decoder.Decode(&user); err != nil {
		log.WithFields(log.Fields{
			"handler": "SignUp",
			"issue":   "failed unmarshalling request body",
		}).Error(err)
		c.Writer.WriteHeader(http.StatusBadRequest)

		return
	}

	if err := user.Validate(); err != nil {
		log.WithFields(log.Fields{
			"handler": "SignUp",
			"issue":   "wrong params",
		}).Error(err)
		c.Writer.WriteHeader(http.StatusBadRequest)

		return
	}

	if _, err := h.usersService.Create(context.TODO(), user); err != nil {
		log.WithFields(log.Fields{
			"handler": "SignUp",
			"issue":   "internal error",
		}).Error(err)
		c.Writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	c.Writer.WriteHeader(http.StatusCreated)
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
//	@Router			/actors [post]
func (h *Handler) AddActor(c *gin.Context) {
	var actor domain.Actor

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

	response, err := json.Marshal(actors)
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "GetAllActors",
			"issue":   "failed marshaling response body",
		}).Error(err)
		c.Writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	c.Writer.Header().Add("Content-Type", "application/json")
	c.Writer.Write(response)
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

	response, err := json.Marshal(actor)
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	c.Writer.Header().Add("Content-Type", "application/json")
	c.Writer.Write(response)
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
