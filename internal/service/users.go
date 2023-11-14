package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/AngelicaNice/HollywoodStarsCRUD/internal/domain"
	"github.com/golang-jwt/jwt"
)

type UsersRepository interface {
	Create(ctx context.Context, user domain.User) (int64, error)
	GetByCredentials(ctx context.Context, email string, hpass string) (domain.User, error)
}

type PasswordHasher interface {
	Hash(password string) (string, error)
}

type Users struct {
	repo     UsersRepository
	hash     PasswordHasher
	secret   []byte
	tokenTtl time.Duration
}

func NewUsers(repo UsersRepository, ph PasswordHasher, secret []byte, token_ttl time.Duration) *Users {
	return &Users{
		repo:     repo,
		hash:     ph,
		secret:   secret,
		tokenTtl: token_ttl,
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

func (u *Users) GetToken(ctx context.Context, input domain.SignInInput) (string, error) {
	hp, err := u.hash.Hash(input.Password)
	if err != nil {
		return "", err
	}

	user, err := u.repo.GetByCredentials(ctx, input.Email, hp)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", domain.ErrUserNotFound
		}

		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   strconv.Itoa(int(user.Id)),
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(u.tokenTtl).Unix(),
	})

	return token.SignedString(u.secret)
}

func (u *Users) ParseToken(ctx context.Context, token string) (int64, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return "", fmt.Errorf("unexpected method %v", token.Header["alg"])
		}

		return u.secret, nil
	})

	if err != nil {
		return 0, err
	}

	if !t.Valid {
		return 0, err
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid claims")
	}

	subj, ok := claims["sub"].(string)
	if !ok {
		return 0, errors.New("invalid sub")
	}

	id, err := strconv.Atoi(subj)
	if err != nil {
		return 0, errors.New("invalid id")
	}

	return int64(id), nil
}
