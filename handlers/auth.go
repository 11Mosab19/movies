package handlers

import (
	"movies/apperrors"
	"movies/dto"
	"movies/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	AS *services.AuthService
	US *services.UserService
}

func (AH *AuthHandler) Register(c *gin.Context) {
	ctx := c.Request.Context()

	var registerData dto.RegisterUserReq
	if err := c.ShouldBindJSON(&registerData); err != nil {
		c.Error(apperrors.ErrBadRequestData)
		c.Abort()
		return
	}

	userRes, err := AH.US.RegisterUser(ctx, registerData)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusCreated, userRes)
}

func (AH *AuthHandler) Login(c *gin.Context) {
	ctx := c.Request.Context()

	var loginData dto.LoginRequest

	err := c.ShouldBindJSON(&loginData)
	if err != nil {
		c.Error(apperrors.ErrBadRequestData)
		c.Abort()
		return
	}
	userRes, err := AH.US.Login(ctx, loginData)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	accessToken, refreshToken, err := AH.AS.CreateSession(ctx, userRes)

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	loginRes := dto.LoginResponse{
		AccessToken:  accessToken.Key,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
	}

	c.JSON(http.StatusOK, loginRes)

}

func (AH *AuthHandler) Refresh(c *gin.Context) {
	ctx := c.Request.Context()

	var refreshReq dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&refreshReq); err != nil {
		c.Error(apperrors.ErrBadRequestData)
		c.Abort()
		return
	}

	accessToken, refreshToken, err := AH.AS.Refresh(ctx, refreshReq.RefreshToken)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	refreshRes := dto.LoginResponse{
		AccessToken:  accessToken.Key,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
	}

	c.JSON(http.StatusOK, refreshRes)
}

func (AH *AuthHandler) Logout(c *gin.Context) {
	ctx := c.Request.Context()

	var logoutReq dto.RefreshTokenRequest

	if err := c.ShouldBindJSON(&logoutReq); err != nil {
		c.Error(apperrors.ErrBadRequestData)
		c.Abort()
		return
	}

	err := AH.AS.Logout(ctx, logoutReq.RefreshToken)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully",
	})
}
