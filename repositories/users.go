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

type UserRepo struct {
	DB *database.Database
}

func (UR *UserRepo) AddUser(ctx context.Context, User models.User) (models.User, error) {
	CreateQuery := `INSERT INTO users (email,hashed_password,full_name) VALUES ($1,$2,$3) RETURNING id,role;`
	var UserData models.User
	err := UR.DB.DB.QueryRowContext(ctx, CreateQuery, User.Email, User.HashedPassword, User.FullName).Scan(&UserData.Id, &UserData.Role)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return models.User{}, apperrors.ErrEmailIsUsed
		}
		return models.User{}, err
	}
	UserData.Email = User.Email
	UserData.FullName = User.FullName
	UserData.HashedPassword = User.HashedPassword
	return UserData, nil
}

func (UP *UserRepo) UpdateUser(ctx context.Context, id int, UserData models.UpdateUserData) error {
	var QueryBuilder strings.Builder

	QueryBuilder.WriteString("UPDATE users SET ")

	Args := make([]any, 0)
	ArgID := 1

	if UserData.FullName != "" {
		QueryBuilder.WriteString(fmt.Sprintf("full_name = $%d, ", ArgID))
		Args = append(Args, UserData.FullName)
		ArgID++
	}

	if UserData.Email != "" {
		QueryBuilder.WriteString(fmt.Sprintf("email = $%d, ", ArgID))
		Args = append(Args, UserData.Email)
		ArgID++
	}

	if UserData.NewHashedPassword != "" {
		QueryBuilder.WriteString(fmt.Sprintf("hashed_password = $%d, ", ArgID))
		Args = append(Args, UserData.NewHashedPassword)
		ArgID++
	}

	if len(Args) == 0 {
		return apperrors.ErrEmptyUpdate
	}

	Query := strings.TrimSuffix(QueryBuilder.String(), ", ")

	Query += fmt.Sprintf(
		" WHERE id = $%d;",
		ArgID,
	)

	Args = append(Args, id)

	result, err := UP.DB.DB.ExecContext(ctx, Query, Args...)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return apperrors.ErrEmailIsUsed
		}

		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return apperrors.ErrUserNotFound
	}

	return nil
}

func (UR *UserRepo) GetUserByEmail(ctx context.Context, UserEmail string) (models.User, error) {
	GetUserQuery := `SELECT id, role, hashed_password, full_name FROM users WHERE email = $1;`
	var UserData models.User
	err := UR.DB.DB.QueryRowContext(ctx, GetUserQuery, UserEmail).Scan(&UserData.Id, &UserData.Role, &UserData.HashedPassword, &UserData.FullName)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, apperrors.ErrUserNotFound
		}
		return models.User{}, err
	}

	UserData.Email = UserEmail

	return UserData, nil
}

func (UR *UserRepo) GetUserById(ctx context.Context, UserId int) (models.User, error) {
	GetUserQuery := `SELECT email, role, hashed_password, full_name FROM users WHERE id=$1;`

	var UserData models.User
	err := UR.DB.DB.QueryRowContext(ctx, GetUserQuery, UserId).Scan(&UserData.Email, &UserData.Role, &UserData.HashedPassword, &UserData.FullName)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, apperrors.ErrUserNotFound
		}
		return models.User{}, err
	}

	UserData.Id = UserId

	return UserData, nil
}

func (UR *UserRepo) GetUserProfile(ctx context.Context, userID int) (models.UserProfile, error) {
	query := `SELECT u.full_name,
		COUNT(r.id) OVER () AS total_reviews,
		r.id,
		COALESCE(r.rating, 0),
		COALESCE(r.review, ''),
		COALESCE(m.poster_url, ''),
		COALESCE(m.title, '')
	FROM users AS u
	LEFT JOIN reviews AS r ON r.user_id = u.id
	LEFT JOIN movies AS m ON m.id = r.movie_id
	WHERE u.id = $1
	ORDER BY r.created_at DESC NULLS LAST, r.id DESC NULLS LAST;`

	rows, err := UR.DB.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return models.UserProfile{}, err
	}
	defer rows.Close()

	profile := models.UserProfile{Reviews: make([]models.UserProfileReview, 0)}
	found := false
	for rows.Next() {
		found = true
		var reviewID sql.NullInt64
		var review models.UserProfileReview
		if err := rows.Scan(
			&profile.FullName,
			&profile.TotalReviews,
			&reviewID,
			&review.Rating,
			&review.Review,
			&review.MovieURL,
			&review.MovieName,
		); err != nil {
			return models.UserProfile{}, err
		}

		if reviewID.Valid {
			review.ID = int(reviewID.Int64)
			profile.Reviews = append(profile.Reviews, review)
		}
	}

	if err := rows.Err(); err != nil {
		return models.UserProfile{}, err
	}
	if !found {
		return models.UserProfile{}, apperrors.ErrUserNotFound
	}

	return profile, nil
}

func (UR *UserRepo) DeleteUser(ctx context.Context, userID int) error {
	result, err := UR.DB.DB.ExecContext(ctx, `DELETE FROM users WHERE id = $1;`, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return apperrors.ErrUserNotFound
	}

	return nil
}
