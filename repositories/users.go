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

func (UP *UserRepo) UpdateUser(ctx context.Context, id int, UserData models.UpdateUserSData) (models.User, error) {
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
		return models.User{}, apperrors.ErrEmptyUpdate
	}

	Query := strings.TrimSuffix(QueryBuilder.String(), ", ")

	Query += fmt.Sprintf(
		" WHERE id = $%d RETURNING id, full_name, email, hashed_password, role;",
		ArgID,
	)
	Args = append(Args, id)

	var user models.User

	err := UP.DB.DB.QueryRowContext(ctx, Query, Args...).Scan(
		&user.Id,
		&user.FullName,
		&user.Email,
		&user.HashedPassword,
		&user.Role,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return models.User{}, apperrors.ErrUserNotFound
	}

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return models.User{}, apperrors.ErrEmailIsUsed
		}

		return models.User{}, err
	}

	return user, nil
}

func (UR *UserRepo) GetUserByEmail(ctx context.Context, UserEmail string) (models.User, error) {
	GetUserQuery := `SELECT id,role,hashed_password,full_name WHERE email=$1;`
	var UserData models.User
	err := UR.DB.DB.QueryRowContext(ctx, GetUserQuery, UserEmail).Scan(&UserData.Id, UserData.Role, &UserData.HashedPassword, &UserData.FullName)
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
	GetUserQuery := `SELECT email,role,hashed_password,full_name WHERE id=$1;`

	var UserData models.User
	err := UR.DB.DB.QueryRowContext(ctx, GetUserQuery, UserId).Scan(&UserData.Email, UserData.Role, &UserData.HashedPassword, &UserData.FullName)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, apperrors.ErrUserNotFound
		}
		return models.User{}, err
	}

	UserData.Id = UserId

	return UserData, nil
}
