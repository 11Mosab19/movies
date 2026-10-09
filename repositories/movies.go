package repositories

import (
	"context"
	"errors"
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
