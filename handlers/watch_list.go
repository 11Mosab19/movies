package handlers

import (
	"movies/apperrors"
	"movies/dto"
	"movies/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type WatchListHandler struct {
	WlS *services.WatchListService
}

// route => movie/:id/watchlist
func (WlH *WatchListHandler) AddMovieToWL(c *gin.Context) {
	ctx := c.Request.Context()

	movieID := c.Param("id")

	movieId, err := strconv.Atoi(movieID)
	if err != nil || movieId <= 0 {
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

	if err := WlH.WlS.AddMovieToWatchList(ctx, movieId, userID); err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.Status(http.StatusCreated)
}

// user/me/watchlist
func (WlH *WatchListHandler) GetUserWatchList(c *gin.Context) {
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

	movies, err := WlH.WlS.GetWatchListMovies(ctx, userID)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, movies)
}

// DELETE/movie/:id/watchlist
func (WlH *WatchListHandler) RemoveMovieFromWL(c *gin.Context) {
	ctx := c.Request.Context()

	movieID := c.Param("id")

	movieId, err := strconv.Atoi(movieID)
	if err != nil || movieId <= 0 {
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

	if err := WlH.WlS.RemoveMovieFromWatchList(ctx, userID, movieId); err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.Status(http.StatusNoContent)
}

func (WlH *WatchListHandler) SetWatchListStatus(c *gin.Context) {
	ctx := c.Request.Context()

	movieID := c.Param("id")

	movieId, err := strconv.Atoi(movieID)
	if err != nil || movieId <= 0 {
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

	var updateData dto.UpdateWLStatus

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.Error(apperrors.ErrBadRequestData)
		c.Abort()
		return
	}

	status := updateData.Status

	err = WlH.WlS.SetWatchListStatus(ctx, userID, movieId, status)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.Status(http.StatusOK)

}
