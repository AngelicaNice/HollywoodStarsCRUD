package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
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

type TokensRepository interface {
	Create(ctx context.Context, rt domain.RefreshToken) error
	Get(ctx context.Context, refreshToken string) (domain.RefreshToken, error)
}

type Users struct {
	repo     UsersRepository
	trepo    TokensRepository
	hash     PasswordHasher
	secret   []byte
	tokenTtl time.Duration
}

func NewUsers(r UsersRepository, tr TokensRepository, ph PasswordHasher, s []byte, tttl time.Duration) *Users {
	return &Users{
		repo:     r,
		trepo:    tr,
		hash:     ph,
		secret:   s,
		tokenTtl: tttl,
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

func (u *Users) GetToken(ctx context.Context, input domain.SignInInput) (string, string, error) {
	hpass, err := u.hash.Hash(input.Password)
	if err != nil {
		return "", "", err
	}

	user, err := u.repo.GetByCredentials(ctx, input.Email, hpass)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", "", domain.ErrUserNotFound
		}

		return "", "", err
	}

	return u.GenerateTokens(ctx, user.Id)
}

func (u *Users) GenerateTokens(ctx context.Context, userId int64) (string, string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   strconv.Itoa(int(userId)),
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(u.tokenTtl).Unix(),
	})

	accessToken, err := token.SignedString(u.secret)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := newRefreshToken()
	if err != nil {
		return "", "", err
	}

	if err := u.trepo.Create(ctx, domain.RefreshToken{
		UserId:    userId,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 30),
	}); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func newRefreshToken() (string, error) {
	b := make([]byte, 32)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	if _, err := r.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
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

func (u *Users) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	rtoken, err := u.trepo.Get(ctx, refreshToken)
	if err != nil {
		return "", "", err
	}

	if rtoken.ExpiresAt.Unix() < time.Now().Unix() {
		return "", "", domain.ErrRefreshTokenExpired
	}

	return u.GenerateTokens(ctx, rtoken.UserId)
}
