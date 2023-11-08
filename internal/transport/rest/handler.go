package rest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/AngelicaNice/HollywoodStarsCRUD/internal/domain"

	"github.com/gin-gonic/gin"
)

type Actors interface {
	Create(ctx context.Context, actor domain.Actor) error
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

func (h *Handler) InitRouter() *gin.Engine {
	r := gin.Default()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	api := r.Group("/actors")
	{
		api.Handle(http.MethodPost, "", h.AddActor)
		api.Handle(http.MethodGet, "", h.GetAllActors)
		api.Handle(http.MethodGet, "/id", h.GetActor)
		api.Handle(http.MethodPut, "/id", h.UpdateActor)
		api.Handle(http.MethodDelete, "/id", h.DeleteActor)
	}

	return r
}

func (h *Handler) AddActor(c *gin.Context) {
	reqBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	var actor domain.Actor
	if err = json.Unmarshal(reqBytes, &actor); err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Println(actor)

	err = h.actorsService.Create(context.TODO(), actor)
	if err != nil {
		log.Println("addActor() error:", err)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.Writer.WriteHeader(http.StatusCreated)
}

func (h *Handler) GetAllActors(c *gin.Context) {
	actors, err := h.actorsService.GetAllActors(context.TODO())
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(actors)
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.Writer.Header().Add("Content-Type", "application/json")
	c.Writer.Write(response)
}

func (h *Handler) GetActor(c *gin.Context) {
	param := c.Request.URL.RawQuery
	param = param[3:]
	id, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	actor, err := h.actorsService.GetByID(context.TODO(), id)
	if err != nil {
		if errors.Is(err, domain.ActorNotFound) {
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

func (h *Handler) UpdateActor(c *gin.Context) {
	param := c.Request.URL.RawQuery
	param = param[3:]
	id, err := strconv.ParseInt(param, 10, 64)

	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	reqBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	var inp domain.UpdateActorInfo
	if err = json.Unmarshal(reqBytes, &inp); err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.actorsService.Update(context.TODO(), id, inp)
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
}

func (h *Handler) DeleteActor(c *gin.Context) {
	id, err := getIdFromRequest(c.Request)
	if err != nil {
		log.Println("DeleteActor() error:", err)
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.actorsService.Delete(context.TODO(), id)
	if err != nil {
		log.Println("DeleteActor() error:", err)
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
