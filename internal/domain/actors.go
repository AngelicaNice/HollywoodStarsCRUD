package domain

import (
	"errors"
)

var ErrActorNotFound = errors.New("actor not found")

type Actor struct {
	ID         int64   `json:"id"`
	Name       string  `json:"name"`
	Surname    string  `json:"surname"`
	Sex        string  `json:"sex"`
	BirthYear  int     `json:"birth_year"`
	BirthPlace string  `json:"birth_place"`
	RestYear   *int    `json:"rest_year"`
	Language   *string `json:"language"`
}

type UpdateActorInfo struct {
	Name     *string `json:"name"`
	Surname  *string `json:"surname"`
	Sex      *string `json:"sex"`
	RestYear *int    `json:"rest_year"`
	Language *string `json:"language"`
}
