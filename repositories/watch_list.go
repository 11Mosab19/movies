package repositories

import (
	"context"
	"errors"
	"movies/apperrors"
	"movies/database"
	"movies/models"

	"github.com/jackc/pgx/v5/pgconn"
)

type WatchListRepo struct {
	DB *database.Database
}

func (WlR *WatchListRepo) AddMovieToWatchList(ctx context.Context, MovieId int, UserId int) error {
	AddMovieQuery := `INSERT INTO watch_list (movie_id,user_id) VALUES ($1,$2);`
	_, err := WlR.DB.DB.ExecContext(ctx, AddMovieQuery, MovieId, UserId)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return apperrors.ErrMovieAlreadyExists
		}
		return err
	}
	return nil
}

func (WlR *WatchListRepo) SetWatchListStatus(ctx context.Context, UserID int, MovieId int, status string) error {
	UpdateQuery := `UPDATE watch_list SET watch_status = $1 WHERE user_id = $2 AND movie_id = $3;`

	result, err := WlR.DB.DB.ExecContext(ctx, UpdateQuery, status, UserID, MovieId)
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

func (WlR *WatchListRepo) RemoveMovieFromWatchList(ctx context.Context, UserID int, MovieId int) error {
	UpdateQuery := `DELETE FROM watch_list WHERE user_id = $1 AND movie_id = $2;`

	result, err := WlR.DB.DB.ExecContext(ctx, UpdateQuery, UserID, MovieId)
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

func (WlR *WatchListRepo) GetWatchListMovies(ctx context.Context, UserId int) ([]models.Movie, error) {
	GetWatchListMoviesQuery := `SELECT m.id,m.title,m.story_line,m.producer_name,m.category,m.release_date,m.poster_url,m.status,m.imdb_rating,m.internal_rating
	FROM watch_list AS wl 
	JOIN movies AS m ON m.id = wl.movie_id
	WHERE wl.user_id = $1
	ORDER BY wl.movie_id;`

	rows, err := WlR.DB.DB.QueryContext(ctx, GetWatchListMoviesQuery, UserId)

	if err != nil {
		return []models.Movie{}, err
	}

	defer rows.Close()
	movies := make([]models.Movie, 0)

	for rows.Next() {
		var movie models.Movie
		err := rows.Scan(&movie.Id, &movie.Title, &movie.StoryLine, &movie.ProducerName, &movie.Category, &movie.ReleaseDate, &movie.PosterUrl, &movie.Status, &movie.ImdbRating, &movie.InternalRating)
		if err != nil {
			return []models.Movie{}, err
		}

		movies = append(movies, movie)
	}

	if err := rows.Err(); err != nil {
		return []models.Movie{}, err
	}

	return movies, nil
}
