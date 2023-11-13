package service

import (
	"context"
	"time"

	"github.com/AngelicaNice/HollywoodStarsCRUD/internal/domain"
)

type UsersRepository interface {
	Create(ctx context.Context, user domain.User) (int64, error)
	GetByID(ctx context.Context, id int64) (domain.User, error)
}

type PasswordHasher interface {
	Hash(password string) (string, error)
}

type Users struct {
	repo UsersRepository
	hash PasswordHasher
}

func NewUsers(repo UsersRepository, h PasswordHasher) *Users {
	return &Users{
		repo: repo,
		hash: h,
	}
}

func (u *Users) Create(ctx context.Context, user domain.SignUpInput) (int64, error) {
	pass, err := u.hash.Hash(user.Password)
	if err != nil {
		return 0, err
	}

	return u.repo.Create(ctx, domain.User{
		Nickname:      user.Nickname,
		Email:         user.Email,
		Password:      pass,
		Registered_at: time.Now(),
	})
}

func (u *Users) GetByID(ctx context.Context, id int64) (domain.User, error) {
	return u.repo.GetByID(ctx, id)
}
