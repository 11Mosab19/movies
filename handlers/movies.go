package handlers

import (
	"io"
	"movies/apperrors"
	"movies/dto"
	"movies/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MovieHandler struct {
	MS *services.MovieService
}

func (MH *MovieHandler) AddMovie(c *gin.Context) {
	ctx := c.Request.Context()

	var movieData dto.MovieReq

	if err := c.ShouldBind(&movieData); err != nil {
		c.Error(apperrors.ErrBadRequestData)
		c.Abort()
		return
	}

	image, err := movieData.Image.Open()
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	defer image.Close()

	if err := MH.MS.AddMovie(ctx, movieData, image); err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.Status(http.StatusCreated)

}

func (MH *MovieHandler) GetMovieById(c *gin.Context) {
	ctx := c.Request.Context()

	movieID := c.Param("id")

	movieId, err := strconv.Atoi(movieID)
	if err != nil || movieId <= 0 {
		c.Error(apperrors.ErrBadRequestData)
		c.Abort()
		return
	}

	movieData, err := MH.MS.GetMovieById(ctx, movieId)

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, movieData)
}

func (MH *MovieHandler) GetMovieByTitle(c *gin.Context) {
	ctx := c.Request.Context()

	title := c.Param("title")

	movies, err := MH.MS.GetMovieByTitle(ctx, title)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, movies)
}

func (MH *MovieHandler) GetMoviesByCategory(c *gin.Context) {
	ctx := c.Request.Context()

	category := c.Param("category")

	movies, err := MH.MS.GetMoviesByCategory(ctx, category)

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, movies)
}

func (MH *MovieHandler) GetMoviesByProducerName(c *gin.Context) {
	ctx := c.Request.Context()

	producerName := c.Param("producer_name")

	movies, err := MH.MS.GetMoviesByProducerName(ctx, producerName)

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, movies)
}

func (MH *MovieHandler) GetTopRatedMovies(c *gin.Context) {
	ctx := c.Request.Context()

	movies, err := MH.MS.GetTopRatedMovies(ctx)

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, movies)
}

func (MH *MovieHandler) GetUpcomingMovies(c *gin.Context) {
	ctx := c.Request.Context()

	movies, err := MH.MS.GetUpcomingMovies(ctx)

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, movies)
}

func (MH *MovieHandler) UpdateMovie(c *gin.Context) {
	ctx := c.Request.Context()

	MovieId := c.Param("id")

	movieId, err := strconv.Atoi(MovieId)

	if err != nil || movieId <= 0 {
		c.Error(apperrors.ErrBadRequestData)
		c.Abort()
		return
	}

	var updateMovieData dto.UpdateMovieReq

	if err := c.ShouldBind(&updateMovieData); err != nil {
		c.Error(apperrors.ErrBadRequestData)
		c.Abort()
		return
	}

	var image io.Reader
	if updateMovieData.Image != nil {
		file, err := updateMovieData.Image.Open()
		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}
		defer file.Close()
		image = file
	}

	err = MH.MS.UpdateMovie(ctx, movieId, updateMovieData, image)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.Status(http.StatusNoContent)
}

func (MH *MovieHandler) DeleteMovie(c *gin.Context) {
	ctx := c.Request.Context()

	MovieId := c.Param("id")

	movieId, err := strconv.Atoi(MovieId)

	if err != nil || movieId <= 0 {
		c.Error(apperrors.ErrBadRequestData)
		c.Abort()
		return
	}

	if err = MH.MS.DeleteMovie(ctx, movieId); err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.Status(http.StatusNoContent)
}
