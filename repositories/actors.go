package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"movies/apperrors"
	"movies/database"
	"movies/models"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
)

type ActorRepo struct {
	DB *database.Database
}

func (AR *ActorRepo) AddActor(ctx context.Context, Actor models.Actor) error {
	AddActorQuery := `INSERT INTO actors (name,home_land,date_of_birth) VALUES ($1,$2,$3);`
	_, err := AR.DB.DB.ExecContext(ctx, AddActorQuery, Actor.Name, Actor.HomeLand, Actor.BirthDate)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return apperrors.ErrActorAlreadyExists
		}
		return err
	}
	return nil
}

func (AR *ActorRepo) GetActorById(ctx context.Context, ActorId int) (models.ActorWithMovies, error) {
	GetActorQuery := `SELECT a.id,a.name,a.birth_date,a.home_land,m.id,m.title,ms.role
	FROM actors AS a
	LEFT JOIN movie_stars AS ms ON ms.star_id = a.id
	LEFT JOIN movies AS m ON m.id = ms.movie_id
	WHERE a.id = $1;`

	var Actor models.ActorWithMovies
	Actor.Movies = make([]models.MovieInActor, 0)

	rows, err := AR.DB.DB.QueryContext(ctx, GetActorQuery, ActorId)
	if err != nil {
		return models.ActorWithMovies{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var Movie models.MovieInActor
		var MovieId sql.NullInt64
		var MovieTitle sql.NullString
		var Role sql.NullString

		err := rows.Scan(
			&Actor.Actor.Id,
			&Actor.Actor.Name,
			&Actor.Actor.BirthDate,
			&Actor.Actor.HomeLand,
			&MovieId,
			&MovieTitle,
			&Role,
		)
		if err != nil {
			return models.ActorWithMovies{}, err
		}

		if MovieId.Valid {
			Movie.MovieId = int(MovieId.Int64)
			Movie.Title = MovieTitle.String
			Movie.Role = Role.String
			Actor.Movies = append(Actor.Movies, Movie)
		}
	}

	if err := rows.Err(); err != nil {
		return models.ActorWithMovies{}, err
	}

	if Actor.Actor.Id == 0 {
		return models.ActorWithMovies{}, apperrors.ErrActorNotFound
	}

	return Actor, nil
}

func (AR *ActorRepo) GetActorByName(ctx context.Context, ActorName string) ([]models.ActorWithMovies, error) {
	GetActorQuery := `SELECT a.id,a.name,a.birth_date,a.home_land,m.id,m.title,ms.role
	FROM actors AS a
	LEFT JOIN movie_stars AS ms ON ms.star_id = a.id
	LEFT JOIN movies AS m ON m.id = ms.movie_id
	WHERE a.name ILIKE $1
	ORDER BY a.name;`

	rows, err := AR.DB.DB.QueryContext(ctx, GetActorQuery, "%"+ActorName+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ActorsMap := make(map[int]*models.ActorWithMovies)
	var Actors []models.ActorWithMovies
	var ActorOrder []int

	for rows.Next() {
		var Actor models.Actor
		var Movie models.MovieInActor
		var MovieId sql.NullInt64
		var MovieTitle sql.NullString
		var Role sql.NullString

		err := rows.Scan(
			&Actor.Id,
			&Actor.Name,
			&Actor.BirthDate,
			&Actor.HomeLand,
			&MovieId,
			&MovieTitle,
			&Role,
		)
		if err != nil {
			return nil, err
		}

		ActorData, exists := ActorsMap[Actor.Id]
		if !exists {
			ActorData = &models.ActorWithMovies{
				Actor:  Actor,
				Movies: make([]models.MovieInActor, 0),
			}
			ActorsMap[Actor.Id] = ActorData
			ActorOrder = append(ActorOrder, Actor.Id)
		}

		if MovieId.Valid {
			Movie.MovieId = int(MovieId.Int64)
			Movie.Title = MovieTitle.String
			Movie.Role = Role.String
			ActorData.Movies = append(ActorData.Movies, Movie)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(ActorsMap) == 0 {
		return nil, apperrors.ErrActorNotFound
	}

	for _, Id := range ActorOrder {
		Actors = append(Actors, *ActorsMap[Id])
	}

	return Actors, nil
}

func (AR *ActorRepo) DeleteActor(ctx context.Context, ActorId int) error {
	DeleteQuery := `DELETE FROM actors WHERE id = $1;`

	Result, err := AR.DB.DB.ExecContext(ctx, DeleteQuery, ActorId)
	if err != nil {
		return err
	}

	RowsAffected, err := Result.RowsAffected()
	if err != nil {
		return err
	}

	if RowsAffected == 0 {
		return apperrors.ErrActorNotFound
	}

	return nil
}

func (AR *ActorRepo) UpdateActor(ctx context.Context, ActorId int, ActorData models.Actor) error {
	var QueryBuilder strings.Builder

	QueryBuilder.WriteString("UPDATE actors SET ")

	Args := make([]any, 0)
	ArgID := 1

	if ActorData.Name != "" {
		QueryBuilder.WriteString(fmt.Sprintf("name = $%d, ", ArgID))
		Args = append(Args, ActorData.Name)
		ArgID++
	}

	if !ActorData.BirthDate.IsZero() {
		QueryBuilder.WriteString(fmt.Sprintf("birth_date = $%d, ", ArgID))
		Args = append(Args, ActorData.BirthDate)
		ArgID++
	}

	if ActorData.HomeLand != "" {
		QueryBuilder.WriteString(fmt.Sprintf("home_land = $%d, ", ArgID))
		Args = append(Args, ActorData.HomeLand)
		ArgID++
	}

	if len(Args) == 0 {
		return nil
	}

	Query := strings.TrimSuffix(QueryBuilder.String(), ", ")

	Query += fmt.Sprintf(" WHERE id = $%d", ArgID)
	Args = append(Args, ActorId)

	Result, err := AR.DB.DB.ExecContext(ctx, Query, Args...)
	if err != nil {
		return err
	}

	RowsAffected, err := Result.RowsAffected()
	if err != nil {
		return err
	}

	if RowsAffected == 0 {
		return apperrors.ErrActorNotFound
	}

	return nil
}
