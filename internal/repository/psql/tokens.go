package psql

import (
	"context"
	"database/sql"

	"github.com/AngelicaNice/HollywoodStarsCRUD/internal/domain"
	_ "github.com/lib/pq"
)

type Tokens struct {
	db *sql.DB
}

func NewTokens(db *sql.DB) *Tokens {
	return &Tokens{
		db: db,
	}
}

func (t *Tokens) Create(ctx context.Context, rt domain.RefreshToken) error {
	_, err := t.db.Exec(
		"INSERT INTO refresh_tokens (id, user_id, token, expires_at) values ($1, $2, $3, $4)",
		rt.Id, rt.UserId, rt.Token, rt.ExpiresAt)

	return err
}

func (t *Tokens) Get(ctx context.Context, token string) (domain.RefreshToken, error) {
	var rt domain.RefreshToken
	if err := t.db.QueryRow(
		"SELECT id, user_id, token, expires_at FROM refresh_tokens WHERE token=$1", token).Scan(
		&rt.Id, &rt.UserId, &rt.Token, &rt.ExpiresAt); err != nil {
		return rt, err
	}

	_, err := t.db.Exec("DELETE FROM refresh_tokens WHERE user_id=$1", rt.UserId)

	return rt, err
}
