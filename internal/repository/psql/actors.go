package psql

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/AngelicaNice/HollywoodStarsCRUD/internal/domain"
	_ "github.com/lib/pq"
)

type Actors struct {
	db *sql.DB
}

func NewActors(db *sql.DB) *Actors {
	return &Actors{
		db: db,
	}
}

func (a *Actors) Create(ctx context.Context, actor domain.Actor) (int64, error) {
	setColumns := make([]string, 0)
	args := make([]interface{}, 0)

	setColumns = append(setColumns, "name", "surname", "sex", "birth_year", "birth_place")
	args = append(args, actor.Name, actor.Surname, actor.Sex, actor.BirthYear, actor.BirthPlace)
	argId := 5
	argIds := "$1, $2, $3, $4, $5"

	if actor.RestYear != nil {
		setColumns = append(setColumns, "rest_year")
		args = append(args, *actor.RestYear)
		argId++
		argIds = argIds + ", $" + strconv.Itoa(argId)
	}

	if actor.Language != nil {
		setColumns = append(setColumns, "language")
		args = append(args, *actor.Language)
		argId++
		argIds = argIds + ", $" + strconv.Itoa(argId)
	}

	setQuery := strings.Join(setColumns, ", ")
	query := fmt.Sprintf("INSERT INTO actors (%s) values (%s) RETURNING ID", setQuery, argIds)

	var id int64

	err := a.db.QueryRow(query, args...).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, err
}

func (a *Actors) GetByID(ctx context.Context, id int64) (domain.Actor, error) {
	var actor domain.Actor
	err := a.db.QueryRow("SELECT id, name, surname, sex, birth_year, birth_place, rest_year, language FROM actors WHERE id=$1", id).
		Scan(&actor.ID, &actor.Name, &actor.Surname, &actor.Sex, &actor.BirthYear,
			&actor.BirthPlace, &actor.RestYear, &actor.Language)

	if err == sql.ErrNoRows {
		return actor, domain.ErrActorNotFound
	}

	return actor, err
}

func (a *Actors) GetAllActors(ctx context.Context) ([]domain.Actor, error) {
	rows, err := a.db.Query("SELECT * FROM actors")
	if err == nil {
		defer rows.Close()
	}

	if err != nil {
		return nil, err
	}

	actors := make([]domain.Actor, 0)

	for rows.Next() {
		var actor domain.Actor
		if err := rows.Scan(&actor.ID, &actor.Name, &actor.Surname, &actor.Sex, &actor.BirthYear, &actor.BirthPlace, &actor.RestYear, &actor.Language); err != nil {
			return nil, err
		}

		actors = append(actors, actor)
	}

	return actors, rows.Err()
}

func (a *Actors) Update(ctx context.Context, id int64, inp domain.UpdateActorInfo) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if inp.Name != nil {
		setValues = append(setValues, fmt.Sprintf("name=$%d", argId))
		args = append(args, *inp.Name)
		argId++
	}

	if inp.Surname != nil {
		setValues = append(setValues, fmt.Sprintf("surname=$%d", argId))
		args = append(args, *inp.Surname)
		argId++
	}

	if inp.Sex != nil {
		setValues = append(setValues, fmt.Sprintf("sex=$%d", argId))
		args = append(args, *inp.Sex)
		argId++
	}

	if inp.RestYear != nil {
		setValues = append(setValues, fmt.Sprintf("rest_year=$%d", argId))
		args = append(args, *inp.RestYear)
		argId++
	}

	if inp.Language != nil {
		setValues = append(setValues, fmt.Sprintf("language=$%d", argId))
		args = append(args, *inp.Language)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE actors SET %s WHERE id=$%d", setQuery, argId)

	args = append(args, id)
	_, err := a.db.Exec(query, args...)

	if err == sql.ErrNoRows {
		return domain.ErrActorNotFound
	}

	return err
}

func (a *Actors) Delete(ctx context.Context, id int64) error {
	_, err := a.db.Exec("DELETE FROM actors WHERE id=$1", id)
	if err == sql.ErrNoRows {
		return domain.ErrActorNotFound
	}

	return err
}
