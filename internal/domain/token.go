package domain

import (
	"errors"
	"time"
)

var ErrRefreshTokenExpired = errors.New("refresh token is expired")

type RefreshToken struct {
	Id        int64
	UserId    int64
	Token     string
	ExpiresAt time.Time
}
