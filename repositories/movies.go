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

func (MR *MoviesRepo) AddMovie(ctx context.Context, Movie models.Movie) error {
	AddMovieQuery := `INSERT  INTO movies (title,story_line,producer_name,category,release_date,status,imdb_rating,internal_rating)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8);`
	_, err := MR.DB.DB.ExecContext(ctx, AddMovieQuery, Movie.Title, Movie.StoryLine, Movie.ProducerName, Movie.Category, Movie.ReleaseDate, Movie.Status, Movie.ImdbRating, Movie.InternalRating)
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
		return nil
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
	FROM movies ORDER BY imdb_rating DESC;`

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

func (MR *MoviesRepo) GetMovieById(ctx context.Context, MovieId int) (models.Movie, error) {
	GetMovieQuery := `SELECT id,title,story_line,producer_name,category,release_date,poster_url,status,imdb_rating,internal_rating
	FROM movies WHERE id=$1;`

	var Movie models.Movie

	err := MR.DB.DB.QueryRowContext(ctx, GetMovieQuery, MovieId).Scan(&MovieId, &Movie.Title, &Movie.StoryLine, &Movie.ProducerName, &Movie.Category, &Movie.ReleaseDate, &Movie.PosterUrl, &Movie.Status, &Movie.ImdbRating, &Movie.InternalRating)
	if errors.Is(err, sql.ErrNoRows) {
		return models.Movie{}, apperrors.ErrMovieNotFound
	}
	if err != nil {
		return models.Movie{}, err
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
	ORDER BY imdb_rating DESC LIMIT 10;`

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
		"title":           true,
		"story_line":      true,
		"producer_name":   true,
		"category":        true,
		"release_date":    true,
		"poster_url":      true,
		"status":          true,
		"imdb_rating":     true,
		"internal_rating": true,
	}

	UpdateMovieQuery := "UPDATE movies SET "
	args := make([]interface{}, 0, len(fields)+1)

	i := 1

	for field, value := range fields {
		if !AllowedFields[field] {
			return fmt.Errorf("invalid field: %s", field)
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
