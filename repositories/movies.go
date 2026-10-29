package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"movies/apperrors"
	db "movies/database"
	"movies/models"

	"github.com/jackc/pgx/v5/pgconn"
)

type MoviesRepo struct {
	DB *db.Database
}

func isValidMovieStatus(status string) bool {
	return status == "released" || status == "upcoming"
}

func (MR *MoviesRepo) AddMovie(ctx context.Context, Movie models.Movie) error {
	if !isValidMovieStatus(Movie.Status) {
		return apperrors.ErrInvalidStatusInput
	}

	AddMovieQuery := `INSERT INTO movies (title, story_line, producer_name, category, release_date, poster_url, status, imdb_rating)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8);`
	_, err := MR.DB.DB.ExecContext(ctx, AddMovieQuery, Movie.Title, Movie.StoryLine, Movie.ProducerName, Movie.Category, Movie.ReleaseDate, Movie.PosterUrl, Movie.Status, Movie.ImdbRating)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return apperrors.ErrMovieAlreadyExists
		}
		return err
	}
	return nil
}

func (MR *MoviesRepo) DeleteMovie(ctx context.Context, MovieId int) error {
	RemoveMovieQuery := `DELETE FROM movies WHERE id=$1;`
	result, err := MR.DB.DB.ExecContext(ctx, RemoveMovieQuery, MovieId)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return apperrors.ErrMovieNotFound
	}

	return nil

}

func (MR *MoviesRepo) GetMovieByTitle(ctx context.Context, MovieTitle string) ([]models.Movie, error) {
	GetMovieQuery := `SELECT id,title,story_line,producer_name,category,release_date,poster_url,status,imdb_rating,internal_rating
	FROM movies WHERE title ILIKE $1;`

	rows, err := MR.DB.DB.QueryContext(ctx, GetMovieQuery, "%"+MovieTitle+"%")
	if err != nil {
		return []models.Movie{}, err
	}
	defer rows.Close()

	movies := make([]models.Movie, 0)

	for rows.Next() {
		var Movie models.Movie
		err := rows.Scan(&Movie.Id, &Movie.Title, &Movie.StoryLine, &Movie.ProducerName, &Movie.Category, &Movie.ReleaseDate, &Movie.PosterUrl, &Movie.Status, &Movie.ImdbRating, &Movie.InternalRating)
		if err != nil {
			return []models.Movie{}, err
		}
		movies = append(movies, Movie)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return movies, nil

}

func (MR *MoviesRepo) GetAllMoviesByRating(ctx context.Context) ([]models.Movie, error) {
	GetMovieQuery := `SELECT id,title,story_line,producer_name,category,release_date,poster_url,status,imdb_rating,internal_rating
	FROM movies ORDER BY imdb_rating DESC NULLS LAST;`

	rows, err := MR.DB.DB.QueryContext(ctx, GetMovieQuery)
	if err != nil {
		return []models.Movie{}, err
	}
	defer rows.Close()

	movies := make([]models.Movie, 0)

	for rows.Next() {
		var Movie models.Movie
		err := rows.Scan(&Movie.Id, &Movie.Title, &Movie.StoryLine, &Movie.ProducerName, &Movie.Category, &Movie.ReleaseDate, &Movie.PosterUrl, &Movie.Status, &Movie.ImdbRating, &Movie.InternalRating)
		if err != nil {
			return []models.Movie{}, err
		}
		movies = append(movies, Movie)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return movies, nil

}

func (MR *MoviesRepo) GetMovieById(ctx context.Context, MovieId int) (models.MovieWithActors, error) {
	GetMovieQuery := `SELECT m.id,m.title,m.story_line,m.producer_name,m.category,m.release_date,m.poster_url,m.status,m.imdb_rating,m.internal_rating,
	ma.actor_id,a.name,ma.movie_id,ma.role,ma.character_name
	FROM movies AS m
	LEFT JOIN movie_actors AS ma ON ma.movie_id = m.id
	LEFT JOIN actors AS a ON a.id = ma.actor_id
	WHERE m.id=$1;`

	var Movie models.MovieWithActors
	Movie.Actors = make([]models.ActorInMovie, 0)

	rows, err := MR.DB.DB.QueryContext(ctx, GetMovieQuery, MovieId)
	if err != nil {
		return models.MovieWithActors{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var Actor models.ActorInMovie
		var ActorId sql.NullInt64
		var ActorName sql.NullString
		var ActorMovieId sql.NullInt64
		var Role sql.NullString
		var CharacterName sql.NullString

		err := rows.Scan(
			&Movie.Movie.Id, &Movie.Movie.Title, &Movie.Movie.StoryLine,
			&Movie.Movie.ProducerName, &Movie.Movie.Category,
			&Movie.Movie.ReleaseDate, &Movie.Movie.PosterUrl,
			&Movie.Movie.Status, &Movie.Movie.ImdbRating, &Movie.Movie.InternalRating,
			&ActorId, &ActorName, &ActorMovieId, &Role, &CharacterName,
		)
		if err != nil {
			return models.MovieWithActors{}, err
		}

		if ActorId.Valid {
			Actor.ActorId = int(ActorId.Int64)
			Actor.Name = ActorName.String
			Actor.MovieId = int(ActorMovieId.Int64)
			Actor.Role = Role.String
			Actor.CharacterName = CharacterName.String
			Movie.Actors = append(Movie.Actors, Actor)
		}
	}

	if err := rows.Err(); err != nil {
		return models.MovieWithActors{}, err
	}

	if Movie.Movie.Id == 0 {
		return models.MovieWithActors{}, apperrors.ErrMovieNotFound
	}

	return Movie, nil
}

func (MR *MoviesRepo) GetMoviesByProducerName(ctx context.Context, ProducerName string) ([]models.Movie, error) {
	GetMovieQuery := `SELECT id,title,story_line,producer_name,category,release_date,poster_url,status,imdb_rating,internal_rating
	FROM movies WHERE producer_name ILIKE $1;`

	rows, err := MR.DB.DB.QueryContext(ctx, GetMovieQuery, "%"+ProducerName+"%")
	if err != nil {
		return []models.Movie{}, err
	}
	defer rows.Close()

	movies := make([]models.Movie, 0)

	for rows.Next() {
		var Movie models.Movie
		err := rows.Scan(&Movie.Id, &Movie.Title, &Movie.StoryLine, &Movie.ProducerName, &Movie.Category, &Movie.ReleaseDate, &Movie.PosterUrl, &Movie.Status, &Movie.ImdbRating, &Movie.InternalRating)
		if err != nil {
			return []models.Movie{}, err
		}
		movies = append(movies, Movie)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return movies, nil

}

func (MR *MoviesRepo) GetMoviesByCategory(ctx context.Context, Category string) ([]models.Movie, error) {
	GetMovieQuery := `SELECT id,title,story_line,producer_name,category,release_date,poster_url,status,imdb_rating,internal_rating
	FROM movies WHERE category ILIKE $1;`

	rows, err := MR.DB.DB.QueryContext(ctx, GetMovieQuery, "%"+Category+"%")
	if err != nil {
		return []models.Movie{}, err
	}
	defer rows.Close()

	movies := make([]models.Movie, 0)

	for rows.Next() {
		var Movie models.Movie
		err := rows.Scan(&Movie.Id, &Movie.Title, &Movie.StoryLine, &Movie.ProducerName, &Movie.Category, &Movie.ReleaseDate, &Movie.PosterUrl, &Movie.Status, &Movie.ImdbRating, &Movie.InternalRating)
		if err != nil {
			return []models.Movie{}, err
		}
		movies = append(movies, Movie)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return movies, nil
}

func (MR *MoviesRepo) GetUpcomingMoviesLimited(ctx context.Context) ([]models.Movie, error) {
	GetMoviesQuery := `SELECT id,title, story_line, producer_name,category,release_date,poster_url,status,imdb_rating,internal_rating FROM movies
    WHERE release_date > CURRENT_DATE OR status = 'upcoming'
    ORDER BY release_date ASC NULLS LAST LIMIT 10;`

	rows, err := MR.DB.DB.QueryContext(ctx, GetMoviesQuery)
	if err != nil {
		return []models.Movie{}, err
	}
	defer rows.Close()

	movies := make([]models.Movie, 0)

	for rows.Next() {
		var Movie models.Movie
		err := rows.Scan(&Movie.Id, &Movie.Title, &Movie.StoryLine, &Movie.ProducerName, &Movie.Category, &Movie.ReleaseDate, &Movie.PosterUrl, &Movie.Status, &Movie.ImdbRating, &Movie.InternalRating)
		if err != nil {
			return []models.Movie{}, err
		}
		movies = append(movies, Movie)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return movies, nil
}

func (MR *MoviesRepo) GetAllUpcomingMovies(ctx context.Context) ([]models.Movie, error) {
	GetMoviesQuery := `SELECT id,title, story_line, producer_name,category,release_date,poster_url,status,imdb_rating,internal_rating FROM movies
    WHERE release_date > CURRENT_DATE OR status = 'upcoming'
    ORDER BY release_date ASC NULLS LAST;`

	rows, err := MR.DB.DB.QueryContext(ctx, GetMoviesQuery)
	if err != nil {
		return []models.Movie{}, err
	}
	defer rows.Close()

	movies := make([]models.Movie, 0)

	for rows.Next() {
		var Movie models.Movie
		err := rows.Scan(&Movie.Id, &Movie.Title, &Movie.StoryLine, &Movie.ProducerName, &Movie.Category, &Movie.ReleaseDate, &Movie.PosterUrl, &Movie.Status, &Movie.ImdbRating, &Movie.InternalRating)
		if err != nil {
			return []models.Movie{}, err
		}
		movies = append(movies, Movie)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return movies, nil
}

func (MR *MoviesRepo) GetTopRatedMovies(ctx context.Context) ([]models.Movie, error) {
	GetMoviesQuery := `SELECT id,title,story_line,producer_name,category,release_date,poster_url,status,imdb_rating,internal_rating FROM movies
	ORDER BY imdb_rating DESC NULLS LAST LIMIT 10;`

	rows, err := MR.DB.DB.QueryContext(ctx, GetMoviesQuery)
	if err != nil {
		return []models.Movie{}, err
	}
	defer rows.Close()

	movies := make([]models.Movie, 0)

	for rows.Next() {
		var Movie models.Movie
		err := rows.Scan(&Movie.Id, &Movie.Title, &Movie.StoryLine, &Movie.ProducerName, &Movie.Category, &Movie.ReleaseDate, &Movie.PosterUrl, &Movie.Status, &Movie.ImdbRating, &Movie.InternalRating)
		if err != nil {
			return []models.Movie{}, err
		}
		movies = append(movies, Movie)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return movies, nil
}

func (MR *MoviesRepo) GetMoviesByIds(ctx context.Context, Ids []int) ([]models.Movie, error) {
	GetMoviesQuery := `SELECT id,title,story_line,producer_name,category,release_date,poster_url,status,imdb_rating,internal_rating FROM movies 
	WHERE id = ANY($1) ORDER BY array_position($1::int[], id);`

	rows, err := MR.DB.DB.QueryContext(ctx, GetMoviesQuery, Ids)
	if err != nil {
		return []models.Movie{}, err
	}
	defer rows.Close()

	movies := make([]models.Movie, 0)

	for rows.Next() {
		var Movie models.Movie
		err := rows.Scan(&Movie.Id, &Movie.Title, &Movie.StoryLine, &Movie.ProducerName, &Movie.Category, &Movie.ReleaseDate, &Movie.PosterUrl, &Movie.Status, &Movie.ImdbRating, &Movie.InternalRating)
		if err != nil {
			return []models.Movie{}, err
		}
		movies = append(movies, Movie)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return movies, nil

}

// dynamic query
func (MR *MoviesRepo) UpdateMovie(ctx context.Context, id int, fields map[string]interface{}) error {
	if len(fields) == 0 {
		return apperrors.ErrEmptyUpdate
	}

	AllowedFields := map[string]bool{
		"title":         true,
		"story_line":    true,
		"producer_name": true,
		"category":      true,
		"release_date":  true,
		"poster_url":    true,
		"status":        true,
		"imdb_rating":   true,
	}

	UpdateMovieQuery := "UPDATE movies SET "
	args := make([]interface{}, 0, len(fields)+1)

	i := 1

	for field, value := range fields {
		if !AllowedFields[field] {
			return fmt.Errorf("invalid field: %s", field)
		}
		if field == "status" {
			status, ok := value.(string)
			if !ok || !isValidMovieStatus(status) {
				return apperrors.ErrInvalidStatusInput
			}
		}

		if i > 1 {
			UpdateMovieQuery += ", "
		}

		UpdateMovieQuery += fmt.Sprintf("%s = $%d", field, i)
		args = append(args, value)
		i++
	}

	UpdateMovieQuery += fmt.Sprintf(" WHERE id = $%d", i)
	args = append(args, id)

	result, err := MR.DB.DB.ExecContext(ctx, UpdateMovieQuery, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return apperrors.ErrMovieNotFound
	}

	return nil
}
