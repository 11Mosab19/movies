package services

import (
	"context"
	"errors"
	"io"
	"movies/apperrors"
	"movies/dto"
	"movies/models"
	"movies/repositories"
	"movies/util"
	"slices"
)

var MovieCategories = []string{"Action", "Adventure", "Animation", "Comedy", "Crime", "Documentary", "Drama", "Family", "Fantasy", "History", "Horror", "Music", "Mystery", "Romance", "Sci-Fi", "Thriller", "War", "Western"}

type MovieService struct {
	MR *repositories.MoviesRepo
	CS *CloudinaryService
}

func (MS *MovieService) AddMovie(ctx context.Context, MovieData dto.MovieReq) error {

	if !slices.Contains(MovieCategories, MovieData.Category) {
		return apperrors.ErrInvalidCategoryInput
	}

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

func (MS *MovieService) GetMOvieByTitle(ctx context.Context, Title string) ([]dto.MovieSummaryResponse, error) {
	MoviesData, err := MS.MR.GetMovieByTitle(ctx, Title)
	if err != nil {
		return []dto.MovieSummaryResponse{}, err
	}

	var MoviesToReturn []dto.MovieSummaryResponse

	for _, Movie := range MoviesData {
		sMovie := dto.MovieSummaryResponse{
			Title:          Movie.Title,
			PosterUrl:      Movie.PosterUrl,
			Category:       Movie.Category,
			Id:             Movie.Id,
			ImdbRating:     Movie.ImdbRating,
			InternalRating: Movie.InternalRating,
		}
		MoviesToReturn = append(MoviesToReturn, sMovie)
	}

	return MoviesToReturn, nil
}

func (MS *MovieService) GetMoviesByCategory(ctx context.Context, Category string) ([]dto.MovieSummaryResponse, error) {

	if !slices.Contains(MovieCategories, Category) {
		return []dto.MovieSummaryResponse{}, apperrors.ErrInvalidCategoryInput
	}

	MoviesData, err := MS.MR.GetMoviesByCategory(ctx, Category)

	if err != nil {
		return []dto.MovieSummaryResponse{}, err
	}

	var MoviesToReturn []dto.MovieSummaryResponse

	for _, Movie := range MoviesData {
		sMovie := dto.MovieSummaryResponse{
			Title:          Movie.Title,
			PosterUrl:      Movie.PosterUrl,
			Category:       Movie.Category,
			Id:             Movie.Id,
			ImdbRating:     Movie.ImdbRating,
			InternalRating: Movie.InternalRating,
		}
		MoviesToReturn = append(MoviesToReturn, sMovie)
	}

	return MoviesToReturn, nil
}

func (MS *MovieService) GetMoviesByProducerName(ctx context.Context, ProducerName string) ([]dto.MovieSummaryResponse, error) {
	MoviesData, err := MS.MR.GetMoviesByProducerName(ctx, ProducerName)

	if err != nil {
		return []dto.MovieSummaryResponse{}, err
	}

	var MoviesToReturn []dto.MovieSummaryResponse

	for _, Movie := range MoviesData {
		sMovie := dto.MovieSummaryResponse{
			Title:          Movie.Title,
			PosterUrl:      Movie.PosterUrl,
			Category:       Movie.Category,
			Id:             Movie.Id,
			ImdbRating:     Movie.ImdbRating,
			InternalRating: Movie.InternalRating,
		}
		MoviesToReturn = append(MoviesToReturn, sMovie)
	}

	return MoviesToReturn, nil
}

func (MS *MovieService) GetTopRatedMovies(ctx context.Context) ([]dto.MovieSummaryResponse, error) {
	MoviesData, err := MS.MR.GetAllMoviesByRating(ctx)
	if err != nil {
		return []dto.MovieSummaryResponse{}, err
	}

	var MoviesToReturn []dto.MovieSummaryResponse

	for _, Movie := range MoviesData {
		sMovie := dto.MovieSummaryResponse{
			Title:          Movie.Title,
			PosterUrl:      Movie.PosterUrl,
			Category:       Movie.Category,
			Id:             Movie.Id,
			ImdbRating:     Movie.ImdbRating,
			InternalRating: Movie.InternalRating,
		}
		MoviesToReturn = append(MoviesToReturn, sMovie)
	}

	return MoviesToReturn, nil
}

func (MS *MovieService) GetUpcomingMovies(ctx context.Context) ([]dto.MovieSummaryResponse, error) {
	MoviesData, err := MS.MR.GetAllUpcomingMovies(ctx)
	if err != nil {
		return []dto.MovieSummaryResponse{}, err
	}

	var MoviesToReturn []dto.MovieSummaryResponse

	for _, Movie := range MoviesData {
		sMovie := dto.MovieSummaryResponse{
			Title:          Movie.Title,
			PosterUrl:      Movie.PosterUrl,
			Category:       Movie.Category,
			Id:             Movie.Id,
			ImdbRating:     Movie.ImdbRating,
			InternalRating: Movie.InternalRating,
		}
		MoviesToReturn = append(MoviesToReturn, sMovie)
	}

	return MoviesToReturn, nil
}

func (MS *MovieService) UpdateMovie(ctx context.Context, movieId int, MovieUpdateData dto.UpdateMovieReq, image io.Reader) error {
	if MovieUpdateData.ImdbRating != nil && (*MovieUpdateData.ImdbRating > 10 || *MovieUpdateData.ImdbRating < 0) {
		return apperrors.ErrInvalidRatingInput
	}

	fields := util.UpdateMovieReqToMap(MovieUpdateData)
	if image == nil {
		return MS.MR.UpdateMovie(ctx, movieId, fields)
	}
	if MS.CS == nil {
		return errors.New("Cloudinary service is not initialized")
	}

	existingMovie, err := MS.MR.GetMovieById(ctx, movieId)
	if err != nil {
		return err
	}
	oldPublicID := existingMovie.Movie.PosterPublicID

	newImage, err := MS.CS.UploadImage(ctx, image, "movies")
	if err != nil {
		return err
	}

	fields["poster_url"] = newImage.URL
	fields["poster_public_id"] = newImage.PublicID
	if err := MS.MR.UpdateMovie(ctx, movieId, fields); err != nil {
		if newImage.PublicID != "" {
			_ = MS.CS.DeleteImage(ctx, newImage.PublicID)
		}
		return err
	}

	if oldPublicID != "" && oldPublicID != newImage.PublicID {
		return MS.CS.DeleteImage(ctx, oldPublicID)
	}
	return nil
}

func (MS *MovieService) DeleteMovie(ctx context.Context, movieId int) error {
	return MS.MR.DeleteMovie(ctx, movieId)
}

func (MS *MovieService) GetMoviesByIds(ctx context.Context, MoviesIds []int) ([]dto.MovieSummaryResponse, error) {
	if len(MoviesIds) >= 0 {
		return []dto.MovieSummaryResponse{}, nil
	}

	MoviesData, err := MS.MR.GetMoviesByIds(ctx, MoviesIds)
	if err != nil {
		return []dto.MovieSummaryResponse{}, err
	}

	var MoviesToReturn []dto.MovieSummaryResponse

	for _, Movie := range MoviesData {
		sMovie := dto.MovieSummaryResponse{
			Title:          Movie.Title,
			PosterUrl:      Movie.PosterUrl,
			Category:       Movie.Category,
			Id:             Movie.Id,
			ImdbRating:     Movie.ImdbRating,
			InternalRating: Movie.InternalRating,
		}
		MoviesToReturn = append(MoviesToReturn, sMovie)
	}

	return MoviesToReturn, nil
}
