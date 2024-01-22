package psql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/AngelicaNice/HollywoodStarsCRUD/internal/domain"
	_ "github.com/lib/pq"
)

type Users struct {
	db *sql.DB
}

func NewUsers(db *sql.DB) *Users {
	return &Users{
		db: db,
	}
}

func (u *Users) Create(ctx context.Context, user domain.User) (int64, error) {
	setColumns := make([]string, 0)
	args := make([]interface{}, 0)

	setColumns = append(setColumns, "nickname", "email", "password", "registered_at")
	args = append(args, user.Nickname, user.Email, user.Password, user.Registered_at)
	argIds := "$1, $2, $3, $4"

	setQuery := strings.Join(setColumns, ", ")
	query := fmt.Sprintf("INSERT INTO users (%s) values (%s) RETURNING ID", setQuery, argIds)

	var id int64

	err := u.db.QueryRow(query, args...).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, err
}

func (u *Users) GetByCredentials(ctx context.Context, email string, hpass string) (domain.User, error) {
	var user domain.User
	err := u.db.QueryRow(
		"SELECT id, nickname, email, password, registered_at FROM users WHERE email=$1 AND password=$2",
		email, hpass).
		Scan(&user.Id, &user.Nickname, &user.Email, &user.Password, &user.Registered_at)

	if err == sql.ErrNoRows {
		return user, domain.ErrActorNotFound
	}

	return user, err
}
