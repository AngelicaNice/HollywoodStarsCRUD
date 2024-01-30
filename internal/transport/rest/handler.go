package rest

import (
	"context"
	"net/http"

	_ "github.com/AngelicaNice/HollywoodStarsCRUD/docs"
	"github.com/AngelicaNice/HollywoodStarsCRUD/internal/domain"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

type Actors interface {
	Create(ctx context.Context, actor domain.ActorInput) (int64, error)
	GetByID(ctx context.Context, id int64) (domain.Actor, error)
	GetAllActors(ctx context.Context) ([]domain.Actor, error)
	Update(ctx context.Context, id int64, info domain.UpdateActorInfo) error
	Delete(ctx context.Context, id int64) error
}

type Users interface {
	Create(ctx context.Context, user domain.SignUpInput) (int64, error)
	GetToken(ctx context.Context, input domain.SignInInput) (string, string, error)
	ParseToken(ctx context.Context, token string) (int64, error)
	RefreshToken(ctx context.Context, refreshToken string) (string, string, error)
	GenerateTokens(ctx context.Context, id int64) (string, string, error)
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
		auth.Handle(http.MethodPost, "/sign-in", h.SignIn)
		auth.Handle(http.MethodGet, "/refresh", h.Refresh)
	}

	api := r.Group("/actors").Use(authMiddleware(h))
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
