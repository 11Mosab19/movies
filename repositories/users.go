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

func (UP *UserRepo) UpdateUserFullName(ctx context.Context, id int, NewName string) (models.User, error) {
	var user models.User
	UpdateQuery := `UPDATE users SET full_name = $1 WHERE id=$2 RETURNING email,hashed_password,role;`
	err := UP.DB.DB.QueryRowContext(ctx, UpdateQuery, NewName, id).Scan(&user.Email, &user.HashedPassword, &user.Role)
	if err == sql.ErrNoRows {
		return models.User{}, apperrors.ErrUserNotFound
	} else if err != nil {
		return models.User{}, err
	}
	user.FullName = NewName
	user.Id = id
	return user, nil
}

func (UP *UserRepo) UpdateUserEmail(ctx context.Context, id int, email string) (models.User, error) {
	var user models.User

	UpdateQuery := `UPDATE users SET email = $1 WHERE id=$2 RETURNING full_name,hashed_password,role;`
	err := UP.DB.DB.QueryRowContext(ctx, UpdateQuery, email, id).Scan(&user.FullName, &user.HashedPassword, &user.Role)
	if err == sql.ErrNoRows {
		return models.User{}, apperrors.ErrUserNotFound
	} else if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return models.User{}, apperrors.ErrEmailIsUsed
			}
		}
		return models.User{}, err
	}
	user.Email = email
	user.Id = id
	return user, nil
}

func (UP *UserRepo) UpdateUserHashedPassword(ctx context.Context, NewHashedPassword string, id int) (models.User, error) {
	var user models.User
	UpdateQuery := `UPDATE users SET hashed_password=$1 WHERE id=$2 RETURNING email,full_name,role;`
	err := UP.DB.DB.QueryRowContext(ctx, UpdateQuery, NewHashedPassword, id).Scan(&user.Email, &user.FullName, &user.Role)
	if err == sql.ErrNoRows {
		return models.User{}, apperrors.ErrUserNotFound
	}
	if err != nil {
		return models.User{}, err
	}
	user.Id = id
	user.HashedPassword = NewHashedPassword
	return user, nil
}
