package service

import (
	"context"

	"github.com/AngelicaNice/HollywoodStarsCRUD/internal/domain"
)

type ActorsRepository interface {
	Create(ctx context.Context, actor domain.Actor) (int64, error)
	GetByID(ctx context.Context, id int64) (domain.Actor, error)
	GetAllActors(ctx context.Context) ([]domain.Actor, error)
	Update(ctx context.Context, id int64, info domain.UpdateActorInfo) error
	Delete(ctx context.Context, id int64) error
}

type Actors struct {
	repo ActorsRepository
}

func NewActors(repo ActorsRepository) *Actors {
	return &Actors{
		repo: repo,
	}
}

func (a *Actors) Create(ctx context.Context, actor domain.Actor) (int64, error) {
	return a.repo.Create(ctx, actor)
}

func (a *Actors) GetByID(ctx context.Context, id int64) (domain.Actor, error) {
	return a.repo.GetByID(ctx, id)
}

func (a *Actors) GetAllActors(ctx context.Context) ([]domain.Actor, error) {
	return a.repo.GetAllActors(ctx)
}

func (a *Actors) Update(ctx context.Context, id int64, info domain.UpdateActorInfo) error {
	return a.repo.Update(ctx, id, info)
}

func (a *Actors) Delete(ctx context.Context, id int64) error {
	return a.repo.Delete(ctx, id)
}
