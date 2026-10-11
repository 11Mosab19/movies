package repositories

import (
	"context"
	"database/sql"
	"errors"
	"movies/apperrors"
	"movies/database"
	"movies/models"

	"github.com/jackc/pgx/v5/pgconn"
)

type ReviewRepo struct {
	DB *database.Database
}

func (RP *ReviewRepo) AddReview(ctx context.Context, Review models.Review) error {
	AddReviewQuery := `INSERT INTO reviews (review,user_id,movie_id,rating) VALUES ($1,$2,$3,$4);`
	_, err := RP.DB.DB.ExecContext(ctx, AddReviewQuery, Review.Review, Review.UserID, Review.MovieId, Review.Rating)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return apperrors.ErrReviewAlreadyExists
		}
		return err
	}
	return nil
}

func (RP *ReviewRepo) GetUserReviews(ctx context.Context, UserID int) ([]models.ReviewWithMovie, error) {
	GetReviewsQuery := `SELECT r.review,r.rating,m.title,m.poster_url 
	FROM movies AS m
	JOIN review AS r ON r.movie_id = m.id
	WHERE r.user_id = $1
	ORDER BY r.created_at DESC;`

	rows, err := RP.DB.DB.QueryContext(ctx, GetReviewsQuery, UserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	UserReviews := make([]models.ReviewWithMovie, 0)

	for rows.Next() {
		var UR models.ReviewWithMovie

		err := rows.Scan(&UR.Review, &UR.Rating, &UR.Title, &UR.PosterUrl)
		if err != nil {
			return nil, err
		}

		UserReviews = append(UserReviews, UR)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return UserReviews, nil
}

func (RP *ReviewRepo) GetReviewWithId(ctx context.Context, ReviewId int) (models.ReviewWithMovie, error) {
	GetReviewQuery := `SELECT r.review,r.rating,m.title,m.poster_url 
	FROM movies AS m
	JOIN review AS r ON r.movie_id = m.id
	WHERE r.id=$1;`

	var UR models.ReviewWithMovie

	err := RP.DB.DB.QueryRowContext(ctx, GetReviewQuery, ReviewId).Scan(&UR.Review, &UR.Rating, &UR.Title, &UR.PosterUrl)

	if errors.Is(err, sql.ErrNoRows) {
		return models.ReviewWithMovie{}, apperrors.ErrMovieNotFound
	}
	if err != nil {
		return models.ReviewWithMovie{}, err
	}
	return UR, nil
}
