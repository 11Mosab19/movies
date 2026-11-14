package services

import (
	"context"
	"movies/apperrors"
	"movies/dto"
	"movies/repositories"
)

type WatchListService struct {
	WlR *repositories.WatchListRepo
}

func (WlS *WatchListService) AddMovieToWatchList(ctx context.Context, movieID int, userId int) error {
	return WlS.WlR.AddMovieToWatchList(ctx, movieID, userId)
}

func (WlS *WatchListService) GetWatchListMovies(ctx context.Context, userID int) ([]dto.WatchListMovieResponse, error) {
	watchListMovies, err := WlS.WlR.GetWatchListMovies(ctx, userID)
	if err != nil {
		return []dto.WatchListMovieResponse{}, err
	}

	responses := make([]dto.WatchListMovieResponse, 0, len(watchListMovies))
	for _, movie := range watchListMovies {
		responses = append(responses, dto.WatchListMovieResponse{
			MovieID:   movie.MovieID,
			Title:     movie.Title,
			PosterURL: movie.PosterURL,
			SetAt:     movie.SetAt,
			Status:    movie.WatchStatus,
		})
	}

	return responses, nil
}

func (WlS *WatchListService) RemoveMovieFromWatchList(ctx context.Context, userID int, movieID int) error {
	return WlS.WlR.RemoveMovieFromWatchList(ctx, userID, movieID)
}

func (WlS *WatchListService) SetWatchListStatus(ctx context.Context, userID int, movieID int, status string) error {
	if status != "listed" && status != "watched" {
		return apperrors.ErrInvalidWLStatus
	}

	return WlS.WlR.SetWatchListStatus(ctx, userID, movieID, status)
}
