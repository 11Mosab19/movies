package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"movies/apperrors"
	"movies/database"
	"movies/models"

	"github.com/jackc/pgx/v5/pgconn"
)

type ActorRepo struct {
	DB *database.Database
}

func (AR *ActorRepo) AddActor(ctx context.Context, Actor models.Actor) error {
	AddActorQuery := `INSERT INTO actors (name,home_land,date_of_birth,image_url,image_public_id) VALUES ($1,$2,$3,$4,$5);`
	_, err := AR.DB.DB.ExecContext(ctx, AddActorQuery, Actor.Name, Actor.HomeLand, Actor.BirthDate, Actor.ImageURL, Actor.ImagePublicID)
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
	GetActorQuery := `SELECT a.id,a.name,a.date_of_birth,a.home_land,COALESCE(a.image_url, ''),COALESCE(a.image_public_id, ''),m.id,m.title,ma.role,ma.character_name
	FROM actors AS a
	LEFT JOIN movie_actors AS ma ON ma.actor_id = a.id
	LEFT JOIN movies AS m ON m.id = ma.movie_id
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
		var CharacterName sql.NullString

		err := rows.Scan(
			&Actor.Actor.Id,
			&Actor.Actor.Name,
			&Actor.Actor.BirthDate,
			&Actor.Actor.HomeLand,
			&Actor.Actor.ImageURL,
			&Actor.Actor.ImagePublicID,
			&MovieId,
			&MovieTitle,
			&Role,
			&CharacterName,
		)
		if err != nil {
			return models.ActorWithMovies{}, err
		}

		if MovieId.Valid {
			Movie.MovieId = int(MovieId.Int64)
			Movie.Title = MovieTitle.String
			Movie.Role = Role.String
			Movie.CharacterName = CharacterName.String
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
	GetActorQuery := `SELECT a.id,a.name,a.date_of_birth,a.home_land,COALESCE(a.image_url, ''),COALESCE(a.image_public_id, ''),m.id,m.title,ma.role,ma.character_name
	FROM actors AS a
	LEFT JOIN movie_actors AS ma ON ma.actor_id = a.id
	LEFT JOIN movies AS m ON m.id = ma.movie_id
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
		var CharacterName sql.NullString

		err := rows.Scan(
			&Actor.Id,
			&Actor.Name,
			&Actor.BirthDate,
			&Actor.HomeLand,
			&Actor.ImageURL,
			&Actor.ImagePublicID,
			&MovieId,
			&MovieTitle,
			&Role,
			&CharacterName,
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
			Movie.CharacterName = CharacterName.String
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

func (AR *ActorRepo) UpdateActor(ctx context.Context, ActorId int, fields map[string]interface{}) error {
	if len(fields) == 0 {
		return apperrors.ErrEmptyUpdate
	}

	AllowedFields := map[string]bool{
		"name":            true,
		"home_land":       true,
		"date_of_birth":   true,
		"image_url":       true,
		"image_public_id": true,
	}

	UpdateActorQuery := "UPDATE actors SET "
	args := make([]interface{}, 0, len(fields)+1)
	i := 1

	for field, value := range fields {
		if !AllowedFields[field] {
			return apperrors.ErrBadRequestData
		}
		if i > 1 {
			UpdateActorQuery += ", "
		}
		UpdateActorQuery += fmt.Sprintf("%s = $%d", field, i)
		args = append(args, value)
		i++
	}

	UpdateActorQuery += fmt.Sprintf(" WHERE id = $%d", i)
	args = append(args, ActorId)

	Result, err := AR.DB.DB.ExecContext(ctx, UpdateActorQuery, args...)
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
