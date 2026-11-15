package handlers

import (
	"movies/apperrors"
	"movies/dto"
	"movies/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	US *services.UserService
	AS *services.AuthService
}

func (UH *UserHandler) GetUserProfile(c *gin.Context) {
	ctx := c.Request.Context()

	userId, exists := c.Get("userId")
	if !exists {
		c.Error(apperrors.ErrUnauthorized)
		c.Abort()
		return
	}

	userID, ok := userId.(int)
	if !ok {
		c.Error(apperrors.ErrUnauthorized)
		c.Abort()
		return
	}

	userProfileData, err := UH.US.GetUserProfile(ctx, userID)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, userProfileData)
}

func (UH *UserHandler) UpdateUser(c *gin.Context) {
	var updateUserReq dto.UpdateUserRequest

	err := c.ShouldBindJSON(&updateUserReq)
	if err != nil {
		c.Error(apperrors.ErrBadRequestData)
		c.Abort()
		return
	}

	userId, exists := c.Get("userId")

	if !exists {
		c.Error(apperrors.ErrUnauthorized)
		c.Abort()
		return
	}

	userID, ok := userId.(int)
	if !ok {
		c.Error(apperrors.ErrUnauthorized)
		c.Abort()
		return
	}

	ctx := c.Request.Context()

	err = UH.US.UpdateUserData(ctx, updateUserReq, userID)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.Status(http.StatusNoContent)
}

func (UH *UserHandler) DeleteUser(c *gin.Context) {
	ctx := c.Request.Context()
	var deleteUserReq dto.DeleteUserReq
	if err := c.ShouldBindJSON(&deleteUserReq); err != nil {
		c.Error(apperrors.ErrBadRequestData)
		c.Abort()
		return
	}

	userId, exists := c.Get("userId")
	if !exists {
		c.Error(apperrors.ErrUnauthorized)
		c.Abort()
		return
	}

	userID, ok := userId.(int)
	if !ok {
		c.Error(apperrors.ErrUnauthorized)
		c.Abort()
		return
	}

	if err := UH.US.DeleteUser(ctx, userID, deleteUserReq); err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.Status(http.StatusNoContent)
}

func (UH *UserHandler) GetUserById(c *gin.Context) {
	ctx := c.Request.Context()
	UserId := c.Param("id")

	userId, err := strconv.Atoi(UserId)

	if err != nil || userId <= 0 {
		c.Error(apperrors.ErrBadRequestData)
		c.Abort()
		return
	}

	userData, err := UH.US.GetUserProfile(ctx, userId)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, userData)
}
