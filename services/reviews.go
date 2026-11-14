package services

import (
	"context"
	"movies/apperrors"
	"movies/dto"
	"movies/models"
	"movies/repositories"
	"strings"
)

type ReviewService struct {
	RR *repositories.ReviewRepo
}

func (RS *ReviewService) AddReview(ctx context.Context, userID int, movieID int, reviewData dto.CreateReviewRequest) error {
	if strings.TrimSpace(reviewData.Review) == "" {
		return apperrors.ErrBadRequestData
	}

	if reviewData.Rating == nil || *reviewData.Rating < 0 || *reviewData.Rating > 10 {
		return apperrors.ErrInvalidRatingInput
	}

	review := models.Review{
		UserID:  userID,
		MovieId: movieID,
		Review:  strings.TrimSpace(reviewData.Review),
		Rating:  *reviewData.Rating,
	}

	return RS.RR.AddReview(ctx, review)
}

func (RS *ReviewService) GetUserReviews(ctx context.Context, userID int) ([]dto.UserReviewResponse, error) {
	reviews, err := RS.RR.GetUserReviews(ctx, userID)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.UserReviewResponse, 0, len(reviews))
	for _, review := range reviews {
		responses = append(responses, dto.UserReviewResponse{
			Review:     review.Review,
			Rating:     review.Rating,
			MovieTitle: review.Title,
			PosterURL:  review.PosterUrl,
		})
	}

	return responses, nil
}

func (RS *ReviewService) GetReviewWithId(ctx context.Context, reviewID int) (dto.DetailedReviewResponse, error) {
	review, err := RS.RR.GetReviewWithId(ctx, reviewID)
	if err != nil {
		return dto.DetailedReviewResponse{}, err
	}

	return dto.DetailedReviewResponse{
		ReviewID:   review.ReviewID,
		MovieID:    review.MovieID,
		UserID:     review.UserID,
		UserName:   review.UserName,
		MovieTitle: review.MovieTitle,
		IMDBRating: review.IMDBRating,
		UserRating: review.UserRating,
		PosterURL:  review.PosterURL,
		Review:     review.Review,
		CreatedAt:  review.CreatedAt,
		UpdatedAt:  review.UpdatedAt,
	}, nil
}

func (RS *ReviewService) GetMovieReviews(ctx context.Context, movieID int) ([]dto.MovieReviewResponse, error) {
	reviews, err := RS.RR.GetMovieReviews(ctx, movieID)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.MovieReviewResponse, 0, len(reviews))
	for _, review := range reviews {
		responses = append(responses, dto.MovieReviewResponse{
			UserName: review.UserName,
			Rating:   review.Rating,
			Review:   review.Review,
		})
	}

	return responses, nil
}

func (RS *ReviewService) UpdateReview(ctx context.Context, userID int, reviewID int, reviewData dto.UpdateReviewRequest) error {
	if reviewData.Review == nil && reviewData.Rating == nil {
		return apperrors.ErrEmptyUpdate
	}

	if reviewData.Review != nil {
		trimmedReview := strings.TrimSpace(*reviewData.Review)
		if trimmedReview == "" {
			return apperrors.ErrBadRequestData
		}
		reviewData.Review = &trimmedReview
	}

	if reviewData.Rating != nil && (*reviewData.Rating < 0 || *reviewData.Rating > 10) {
		return apperrors.ErrInvalidRatingInput
	}

	return RS.RR.UpdateReview(ctx, userID, reviewID, models.UpdateReviewData{
		Review: reviewData.Review,
		Rating: reviewData.Rating,
	})
}

func (RS *ReviewService) DeleteReview(ctx context.Context, userID int, reviewID int) error {
	return RS.RR.DeleteReview(ctx, userID, reviewID)
}
