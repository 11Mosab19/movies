package services

import (
	"context"
	"errors"
	"movies/apperrors"
	"movies/dto"
	"movies/models"
	"movies/repositories"
	"movies/util"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	US repositories.UserRepo
	AR repositories.AuthRepo
}

func GetConfig() string {
	SecretKey := os.Getenv("SECRET_KEY")
	return SecretKey
}

func (AS *AuthService) GenerateAccessToken(ctx context.Context, UserData models.User) (dto.Token, error) {
	claims := models.Claims{
		UserId: UserData.Id,
		Role:   UserData.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 90)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	MySecretKey := []byte(GetConfig())
	tokenString, err := token.SignedString(MySecretKey)
	if err != nil {
		return dto.Token{}, apperrors.ErrGenerateAccessToken
	}
	return dto.Token{Key: tokenString}, nil
}

func (AS *AuthService) VerifyToken(TokenString string) (*models.Claims, error) {
	MySecretKey := []byte(GetConfig())
	VerifiedClaims := &models.Claims{}
	Returned, err := jwt.ParseWithClaims(TokenString, VerifiedClaims, func(t *jwt.Token) (any, error) {
		if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, apperrors.ErrUnsupportedAuthenticationMethod
		}
		return MySecretKey, nil
	})
	if err != nil {
		return &models.Claims{}, apperrors.ErrUnauthorized
	}
	if !Returned.Valid {
		return &models.Claims{}, apperrors.ErrUnauthorized
	}
	return VerifiedClaims, nil
}

func (AS *AuthService) CreateSession(ctx context.Context, user models.User) (dto.Token, string, error) {
	accessToken, err := AS.GenerateAccessToken(ctx, user)
	if err != nil {
		return dto.Token{}, "", err
	}

	refreshToken, err := util.GenerateRefreshTokenString()
	if err != nil {
		return dto.Token{}, "", err
	}

	if err := AS.AR.CreateRefreshToken(ctx, user.Id, util.HashRefreshToken(refreshToken)); err != nil {
		return dto.Token{}, "", err
	}

	return accessToken, refreshToken, nil
}

func (AS *AuthService) Refresh(ctx context.Context, refreshToken string) (dto.Token, string, error) {
	tokenHash := util.HashRefreshToken(refreshToken)
	token, err := AS.AR.GetByHash(ctx, tokenHash)
	if err != nil {
		return dto.Token{}, "", err
	}

	if !token.ExpiresAt.After(time.Now()) {
		if err := AS.AR.DeleteRefreshToken(ctx, tokenHash); err != nil {
			return dto.Token{}, "", err
		}
		return dto.Token{}, "", apperrors.ErrRefreshTokenExpired
	}

	user, err := AS.US.GetUserById(ctx, token.UserID)
	if err != nil {
		return dto.Token{}, "", err
	}

	accessToken, err := AS.GenerateAccessToken(ctx, user)
	if err != nil {
		return dto.Token{}, "", err
	}

	newRefreshToken, err := util.GenerateRefreshTokenString()
	if err != nil {
		return dto.Token{}, "", err
	}

	if err := AS.AR.RotateRefreshToken(ctx, tokenHash, util.HashRefreshToken(newRefreshToken)); err != nil {
		return dto.Token{}, "", err
	}

	return accessToken, newRefreshToken, nil
}

func (AS *AuthService) Logout(ctx context.Context, refreshToken string) error {
	err := AS.AR.DeleteRefreshToken(ctx, util.HashRefreshToken(refreshToken))
	if errors.Is(err, apperrors.ErrInvalidRefreshToken) {
		return nil
	}
	return err
}
