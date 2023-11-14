package domain

import (
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

var ErrUserNotFound = errors.New("user not found")

func init() {
	validate = validator.New()
}

type User struct {
	Id            int64     `json:"id"`
	Nickname      string    `json:"nickname"`
	Email         string    `json:"email"`
	Password      string    `json:"password"`
	Registered_at time.Time `json:"registered_at"`
}

type SignUpInput struct {
	Nickname string `json:"nickname" validate:"required,gte=2"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,gte=6"`
}

func (input SignUpInput) Validate() error {
	return validate.Struct(input)
}

type SignInInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,gte=6"`
}

func (input SignInInput) Validate() error {
	return validate.Struct(input)
}
