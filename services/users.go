package services

import (
	"context"
	"errors"
	"movies/apperrors"
	"movies/dto"
	"movies/models"
	"movies/repositories"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	UP *repositories.UserRepo
}

func (US *UserService) Login(ctx context.Context, LoginData dto.LoginRequest) (dto.UserRes, error) {
	User, err := US.UP.GetUserByEmail(ctx, LoginData.Email)
	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			return dto.UserRes{}, apperrors.ErrInvalidCredentials
		}
		return dto.UserRes{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(User.HashedPassword), []byte(LoginData.Password))
	if err != nil {
		return dto.UserRes{}, apperrors.ErrInvalidCredentials
	}

	var ResUser = dto.UserRes{
		Id:       User.Id,
		FullName: User.FullName,
		Role:     User.Role,
	}
	return ResUser, err
}

func (US *UserService) RegisterUser(ctx context.Context, registerData dto.RegisterUserReq) (dto.UserRes, error) {

	if registerData.Password != registerData.ConfirmationPassword {
		return dto.UserRes{}, apperrors.ErrWrongConfirmationPassword
	}

	var user models.User
	HashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerData.Password), bcrypt.DefaultCost)
	if err != nil {
		return dto.UserRes{}, apperrors.ErrHashing
	}
	registerData.Password = string(HashedPassword)

	user.Email = registerData.Email
	user.FullName = registerData.FullName
	user.HashedPassword = registerData.Password

	user, err = US.UP.AddUser(ctx, user)
	if err != nil {
		return dto.UserRes{}, err
	}

	var ResUser = dto.UserRes{
		Id:       user.Id,
		FullName: user.FullName,
		Role:     user.Role,
	}

	return ResUser, nil

}

func (US *UserService) UpdateUserData(ctx context.Context, UpdateUserData dto.UpdateUserRequest, UserId int) error {

	User, err := US.UP.GetUserById(ctx, UserId)
	if err != nil {
		return err
	}

	var UpdateData models.UpdateUserData

	UpdateData.FullName = UpdateUserData.FullName
	UpdateData.Email = UpdateUserData.Email

	if UpdateUserData.NewPassword != "" {

		if UpdateUserData.OldPassword == "" {
			return apperrors.ErrBadRequestData
		}

		err := bcrypt.CompareHashAndPassword([]byte(User.HashedPassword), []byte(UpdateUserData.OldPassword))

		if err != nil {
			return apperrors.ErrWrongUserPassword
		}

		HashedPassword, err := bcrypt.GenerateFromPassword([]byte(UpdateUserData.NewPassword), bcrypt.DefaultCost)

		if err != nil {
			return apperrors.ErrHashing
		}

		UpdateData.NewHashedPassword = string(HashedPassword)
	}

	err = US.UP.UpdateUser(ctx, UserId, UpdateData)
	if err != nil {
		return err
	}
	return nil
}

func (US *UserService) GetUserProfile(ctx context.Context, userID int) (dto.UserProfileResponse, error) {
	profile, err := US.UP.GetUserProfile(ctx, userID)
	if err != nil {
		return dto.UserProfileResponse{}, err
	}

	response := dto.UserProfileResponse{
		FullName:     profile.FullName,
		Reviews:      make([]dto.UserProfileReview, 0, len(profile.Reviews)),
		TotalReviews: profile.TotalReviews,
	}
	for _, review := range profile.Reviews {
		response.Reviews = append(response.Reviews, dto.UserProfileReview{
			ID:        review.ID,
			Rating:    review.Rating,
			Review:    review.Review,
			MovieURL:  review.MovieURL,
			MovieName: review.MovieName,
		})
	}

	return response, nil
}

func (US *UserService) VerifyPassword(ctx context.Context, userID int, password string) error {
	user, err := US.UP.GetUserById(ctx, userID)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password)); err != nil {
		return apperrors.ErrWrongUserPassword
	}

	return nil
}

func (US *UserService) DeleteUser(ctx context.Context, userID int, request dto.DeleteUserReq) error {
	if err := US.VerifyPassword(ctx, userID, request.Password); err != nil {
		return err
	}

	return US.UP.DeleteUser(ctx, userID)
}
