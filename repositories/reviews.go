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

func (RP *ReviewRepo) UpdateReview(ctx context.Context, UserId int, ReviewId int, ReviewData models.Review) error {
	var QueryBuilder strings.Builder

	QueryBuilder.WriteString("UPDATE review SET ")

	Args := make([]any, 0)
	ArgID := 1

	if ReviewData.Review != "" {
		QueryBuilder.WriteString(fmt.Sprintf("review = $%d, ", ArgID))
		Args = append(Args, ReviewData.Review)
		ArgID++
	}

	if ReviewData.Rating > 0 {
		QueryBuilder.WriteString(fmt.Sprintf("rating = $%d, ", ArgID))
		Args = append(Args, ReviewData.Rating)
		ArgID++
	}

	if len(Args) == 0 {
		return nil
	}

	Query := strings.TrimSuffix(QueryBuilder.String(), ", ")

	Query += fmt.Sprintf(" WHERE id = $%d AND user_id = $%d", ArgID, ArgID+1)

	Args = append(Args, ReviewId, UserId)

	_, err := RP.DB.DB.ExecContext(ctx, Query, Args...)
	if err != nil {
		return err
	}

	return nil
}

func (RP *ReviewRepo) DeleteReview(ctx context.Context, UserId int, ReviewId int) error {
	DeleteQuery := `DELETE FROM reviews WHERE user_id=$1 AND id=$2;`

	result, err := RP.DB.DB.ExecContext(ctx, DeleteQuery, UserId, ReviewId)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return apperrors.ErrReviewNotFound
	}

	return nil
}
