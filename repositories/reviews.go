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

func (RP *ReviewRepo) GetUserReviews(ctx context.Context, UserID int) ([]models.ReviewsForUserWithMovie, error) {
	GetReviewsQuery := `SELECT r.review,r.rating,m.title,m.poster_url 
	FROM movies AS m
	JOIN reviews AS r ON r.movie_id = m.id
	WHERE r.user_id = $1
	ORDER BY r.created_at DESC;`

	rows, err := RP.DB.DB.QueryContext(ctx, GetReviewsQuery, UserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	UserReviews := make([]models.ReviewsForUserWithMovie, 0)

	for rows.Next() {
		var UR models.ReviewsForUserWithMovie

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

func (RP *ReviewRepo) GetReviewWithId(ctx context.Context, reviewID int) (models.DetailedReview, error) {
	getReviewQuery := `SELECT r.id, r.movie_id, r.user_id, u.full_name, m.title,
		m.imdb_rating, r.rating, m.poster_url, r.review, r.created_at, r.updated_at
	FROM reviews AS r
	JOIN users AS u ON u.id = r.user_id
	JOIN movies AS m ON m.id = r.movie_id
	WHERE r.id = $1;`

	var review models.DetailedReview

	err := RP.DB.DB.QueryRowContext(ctx, getReviewQuery, reviewID).Scan(
		&review.ReviewID,
		&review.MovieID,
		&review.UserID,
		&review.UserName,
		&review.MovieTitle,
		&review.IMDBRating,
		&review.UserRating,
		&review.PosterURL,
		&review.Review,
		&review.CreatedAt,
		&review.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return models.DetailedReview{}, apperrors.ErrReviewNotFound
	}
	if err != nil {
		return models.DetailedReview{}, err
	}
	return review, nil
}

func (RP *ReviewRepo) GetMovieReviews(ctx context.Context, movieID int) ([]models.ReviewForMovie, error) {
	getMovieReviewsQuery := `SELECT u.full_name, r.rating, r.review
	FROM reviews AS r
	JOIN users AS u ON u.id = r.user_id
	WHERE r.movie_id = $1
	ORDER BY r.created_at DESC;`

	rows, err := RP.DB.DB.QueryContext(ctx, getMovieReviewsQuery, movieID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reviews := make([]models.ReviewForMovie, 0)
	for rows.Next() {
		var review models.ReviewForMovie
		if err := rows.Scan(&review.UserName, &review.Rating, &review.Review); err != nil {
			return nil, err
		}

		reviews = append(reviews, review)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return reviews, nil
}

func (RP *ReviewRepo) UpdateReview(ctx context.Context, UserId int, ReviewId int, ReviewData models.UpdateReviewData) error {
	var QueryBuilder strings.Builder

	QueryBuilder.WriteString("UPDATE reviews SET ")

	Args := make([]any, 0)
	ArgID := 1

	if ReviewData.Review != nil {
		QueryBuilder.WriteString(fmt.Sprintf("review = $%d, ", ArgID))
		Args = append(Args, *ReviewData.Review)
		ArgID++
	}

	if ReviewData.Rating != nil {
		if *ReviewData.Rating < 0 || *ReviewData.Rating > 10 {
			return apperrors.ErrBadRequestData
		}
		QueryBuilder.WriteString(fmt.Sprintf("rating = $%d, ", ArgID))
		Args = append(Args, *ReviewData.Rating)
		ArgID++
	}

	if len(Args) == 0 {
		return apperrors.ErrEmptyUpdate
	}

	Query := strings.TrimSuffix(QueryBuilder.String(), ", ")

	Query += fmt.Sprintf(" WHERE id = $%d AND user_id = $%d", ArgID, ArgID+1)

	Args = append(Args, ReviewId, UserId)

	result, err := RP.DB.DB.ExecContext(ctx, Query, Args...)
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
