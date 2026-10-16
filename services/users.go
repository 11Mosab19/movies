package services

import (
	"context"
	"movies/apperrors"
	"movies/dto"
	"movies/models"
	"movies/repositories"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	UP *repositories.UserRepo
}

func (US *UserService) Login(ctx context.Context, LoginData dto.LoginRequest) (models.User, error) {
	User, err := US.UP.GetUserByEmail(ctx, LoginData.Email)
	if err != nil {
		return models.User{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(User.HashedPassword), []byte(LoginData.Password))
	if err != nil {
		return models.User{}, apperrors.ErrBadRequestData
	}
	return User, err
}

func (US *UserService) RegisterUser(ctx context.Context, registerData dto.RegisterUserReq) (models.User, error) {
	var user models.User
	HashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerData.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, apperrors.HashingErr
	}
	registerData.Password = string(HashedPassword)

	user.Email = registerData.Email
	user.FullName = registerData.FullName
	user.HashedPassword = registerData.Password

	user, err = US.UP.AddUser(ctx, user)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (US *UserService) UpdateUserData(ctx context.Context, UpdateUserData dto.UpdateUserRequest, UserId int) (models.User, error) {

	User, err := US.UP.GetUserById(ctx, UserId)
	if err != nil {
		return models.User{}, err
	}

	var UpdateData models.UpdateUserSData

	UpdateData.FullName = UpdateUserData.FullName
	UpdateData.Email = UpdateUserData.Email

	if UpdateUserData.NewPassword != "" {

		if UpdateUserData.OldPassword == "" {
			return models.User{}, apperrors.ErrBadRequestData
		}

		err := bcrypt.CompareHashAndPassword([]byte(User.HashedPassword), []byte(UpdateUserData.OldPassword))

		if err != nil {
			return models.User{}, apperrors.ErrWrongUserPassword
		}

		HashedPassword, err := bcrypt.GenerateFromPassword([]byte(UpdateUserData.NewPassword), bcrypt.DefaultCost)

		if err != nil {
			return models.User{}, apperrors.HashingErr
		}

		UpdateData.NewHashedPassword = string(HashedPassword)
	}

	UpdatedUser, err := US.UP.UpdateUser(ctx, UserId, UpdateData)
	if err != nil {
		return models.User{}, err
	}

	return UpdatedUser, nil
}
