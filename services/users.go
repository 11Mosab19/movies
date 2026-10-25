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

func (US *UserService) UpdateUserData(ctx context.Context, UpdateUserData dto.UpdateUserRequest, UserId int) (dto.UserRes, error) {

	User, err := US.UP.GetUserById(ctx, UserId)
	if err != nil {
		return dto.UserRes{}, err
	}

	var UpdateData models.UpdateUserData

	UpdateData.FullName = UpdateUserData.FullName
	UpdateData.Email = UpdateUserData.Email

	if UpdateUserData.NewPassword != "" {

		if UpdateUserData.OldPassword == "" {
			return dto.UserRes{}, apperrors.ErrBadRequestData
		}

		err := bcrypt.CompareHashAndPassword([]byte(User.HashedPassword), []byte(UpdateUserData.OldPassword))

		if err != nil {
			return dto.UserRes{}, apperrors.ErrWrongUserPassword
		}

		HashedPassword, err := bcrypt.GenerateFromPassword([]byte(UpdateUserData.NewPassword), bcrypt.DefaultCost)

		if err != nil {
			return dto.UserRes{}, apperrors.ErrHashing
		}

		UpdateData.NewHashedPassword = string(HashedPassword)
	}

	UpdatedUser, err := US.UP.UpdateUser(ctx, UserId, UpdateData)
	if err != nil {
		return dto.UserRes{}, err
	}

	var ResUser = dto.UserRes{
		Id:       UpdatedUser.Id,
		FullName: UpdatedUser.FullName,
		Role:     UpdatedUser.Role,
	}

	return ResUser, nil
}
