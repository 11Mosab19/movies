package handlers

import (
	"movies/apperrors"
	"movies/dto"
	"movies/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ReviewHandler struct {
	RS *services.ReviewService
}

// route => POST/movie/:id/review
func (RH *ReviewHandler) CreateReview(c *gin.Context) {

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

	var reviewData dto.CreateReviewRequest

	err := c.ShouldBindJSON(&reviewData)
	if err != nil {
		c.Error(apperrors.ErrBadRequestData)
		c.Abort()
		return
	}

	movieID := c.Param("id")

	movieId, err := strconv.Atoi(movieID)
	if err != nil || movieId <= 0 {
		c.Error(apperrors.ErrBadRequestData)
		c.Abort()
		return
	}

	err = RH.RS.AddReview(ctx, userID, movieId, reviewData)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.Status(http.StatusCreated)
}

// route => GET/movie/:id/review
func (RH *ReviewHandler) GetMovieReviews(c *gin.Context) {
	ctx := c.Request.Context()

	movieID := c.Param("id")

	movieId, err := strconv.Atoi(movieID)
	if err != nil || movieId <= 0 {
		c.Error(apperrors.ErrBadRequestData)
		c.Abort()
		return
	}

	reviews, err := RH.RS.GetMovieReviews(ctx, movieId)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, reviews)
}

// route => /review/:id
func (RH *ReviewHandler) GetReviewWithId(c *gin.Context) {
	ctx := c.Request.Context()

	reviewID := c.Param("id")

	reviewId, err := strconv.Atoi(reviewID)
	if err != nil || reviewId <= 0 {
		c.Error(apperrors.ErrBadRequestData)
		c.Abort()
		return
	}

	reviewData, err := RH.RS.GetReviewWithId(ctx, reviewId)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, reviewData)
}

// route => user/:id/review
func (RH *ReviewHandler) GetUserReviews(c *gin.Context) {
	ctx := c.Request.Context()

	userID := c.Param("id")

	userId, err := strconv.Atoi(userID)
	if err != nil || userId <= 0 {
		c.Error(apperrors.ErrBadRequestData)
		c.Abort()
		return
	}

	reviews, err := RH.RS.GetUserReviews(ctx, userId)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, reviews)
}

// route => PATCH/review/:id
func (RH *ReviewHandler) UpdateReview(c *gin.Context) {
	ctx := c.Request.Context()

	reviewID := c.Param("id")

	reviewId, err := strconv.Atoi(reviewID)
	if err != nil || reviewId <= 0 {
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

	var updateData dto.UpdateReviewRequest

	err = c.ShouldBindJSON(&updateData)
	if err != nil {
		c.Error(apperrors.ErrBadRequestData)
		c.Abort()
		return
	}

	err = RH.RS.UpdateReview(ctx, userID, reviewId, updateData)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.Status(http.StatusOK)
}

func (RH *ReviewHandler) DeleteReview(c *gin.Context) {
	ctx := c.Request.Context()

	reviewID := c.Param("id")

	reviewId, err := strconv.Atoi(reviewID)
	if err != nil || reviewId <= 0 {
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

	err = RH.RS.DeleteReview(ctx, userID, reviewId)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.Status(http.StatusNoContent)
}
