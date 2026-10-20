package repositories

import (
	"context"
	"database/sql"
	"movies/apperrors"
	"movies/database"
	"movies/models"
	"time"
)

type AuthRepo struct {
	DB *database.Database
}

func (AR *AuthRepo) CreateRefreshToken(ctx context.Context, UserId int, HashedToken string) error {
	ExpireAt := time.Now().AddDate(0, 1, 0)

	AddRefreshQuery := `INSERT INTO refresh_tokens (user_id,token_hash,expires_at) VALUES ($1,$2,$3);`
	_, err := AR.DB.DB.ExecContext(ctx, AddRefreshQuery, UserId, HashedToken, ExpireAt)

	return err
}

func (AR *AuthRepo) GetByHash(ctx context.Context, tokenHash string) (models.RefreshToken, error) {
	query := `SELECT user_id, expires_at FROM refresh_tokens WHERE token_hash = $1`

	var token models.RefreshToken

	err := AR.DB.DB.QueryRowContext(ctx, query, tokenHash).Scan(&token.UserID, &token.ExpiresAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.RefreshToken{}, apperrors.ErrInvalidRefreshToken
		}
		return models.RefreshToken{}, err
	}

	return token, nil
}

func (AR *AuthRepo) DeleteRefreshToken(ctx context.Context, HashedToken string) error {
	query := `DELETE FROM refresh_tokens WHERE token_hash = $1;`

	result, err := AR.DB.DB.ExecContext(ctx, query, HashedToken)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return apperrors.ErrInvalidRefreshToken
	}

	return nil
}

func (AR *AuthRepo) RotateRefreshToken(ctx context.Context, oldTokenHash string, newTokenHash string) error {
	expiresAt := time.Now().AddDate(0, 1, 0)
	query := `UPDATE refresh_tokens
		SET token_hash = $1, expires_at = $2, created_at = CURRENT_TIMESTAMP
		WHERE token_hash = $3 AND expires_at > CURRENT_TIMESTAMP;`

	result, err := AR.DB.DB.ExecContext(ctx, query, newTokenHash, expiresAt, oldTokenHash)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return apperrors.ErrInvalidRefreshToken
	}

	return nil
}
