package services

import (
	"context"
	"movies/apperrors"
	"movies/dto"
	"movies/models"
	"movies/repositories"
)

type MovieService struct {
	MR *repositories.MoviesRepo
}

func (MS *MovieService) AddMovie(ctx context.Context, MovieData dto.MovieReq) error {

	if *MovieData.ImdbRating < 0 || *MovieData.ImdbRating > 10 {
		return apperrors.ErrInvalidRatingInput
	}

	Movie := models.Movie{
		Title:        MovieData.Title,
		PosterUrl:    MovieData.PosterUrl,
		ReleaseDate:  &MovieData.ReleaseDate,
		ImdbRating:   MovieData.ImdbRating,
		StoryLine:    MovieData.StoryLine,
		Status:       MovieData.Status,
		Category:     MovieData.Category,
		ProducerName: MovieData.ProducerName,
	}

	err := MS.MR.AddMovie(ctx, Movie)

	return err
}

func (MS *MovieService) GetMovieById(ctx context.Context, movieId int) (dto.MovieDetailsResponse, error) {
	MovieWithActorsData, err := MS.MR.GetMovieById(ctx, movieId)
	if err != nil {
		return dto.MovieDetailsResponse{}, err
	}
	var MovieData dto.MovieDetailsResponse
	MovieData.ImdbRating = MovieWithActorsData.Movie.ImdbRating
	MovieData.PosterUrl = MovieWithActorsData.Movie.PosterUrl
	MovieData.ProducerName = MovieWithActorsData.Movie.ProducerName
	MovieData.Title = MovieWithActorsData.Movie.Title
	MovieData.Id = MovieWithActorsData.Movie.Id
	MovieData.Category = MovieWithActorsData.Movie.Category
	MovieData.StoryLine = MovieWithActorsData.Movie.StoryLine
	MovieData.InternalRating = MovieWithActorsData.Movie.InternalRating
	MovieData.ReleaseDate = MovieWithActorsData.Movie.ReleaseDate
	MovieData.Status = MovieWithActorsData.Movie.Status
	MovieData.Actors = make([]dto.ActorRes, 0, len(MovieWithActorsData.Actors))

	for _, actor := range MovieWithActorsData.Actors {
		MovieData.Actors = append(MovieData.Actors, dto.ActorRes{
			ActorId:       actor.ActorId,
			Name:          actor.Name,
			Role:          actor.Role,
			CharacterName: actor.CharacterName,
		})
	}

	return MovieData, err

}
