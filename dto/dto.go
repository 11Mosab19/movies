package dto

import "time"

type CreateUserReq struct {
	Email                string `email:"json" binding:"required,email"`
	Password             string `password:"json" binding:"required,min=8"`
	ConfirmationPassword string `json:"confirmation_password" binding:"required,eqfield=Password"`
	fullName             string `full_name:"json" binding:"required"`
}

type LoginRequest struct {
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type UpdateUserRequest struct {
	NewPassword string `json:"new_password" binding:"min=8"`
	OldPassword string `json:"old_password"`
	FullName    string `json:"full_name"`
	Email       string `json:"email" binding:"email"`
}

type SetPasswordRequest struct {
	NewPassword          string `json:"new_password" binding:"min=8"`
	ConfirmationPassword string `json:"confirmation_password" binding:"required,eqfield=NewPassword"`
}

type Token struct {
	Key string `json:"token"`
}

type AddMovieReq struct {
	Name         string    `name:"json" binding:"required"`
	StoryLine    string    `story_line:"json" binding:"required"`
	ImdbRating   float64   `imdb_rating:"json"`
	Status       string    `status:"json" binding:"required"`
	ReleaseDate  time.Time `release_date:"json" binding:"required"`
	ProducerName string    `producer_name:"json" binding:"required"`
	Category     string    `category:"json" binding:"required"`
}

type AddStarReq struct {
	Name        string    `name:"json" binding:"required"`
	HomeLand    string    `home_land:"json" binding:"required"`
	DateOfBirth time.Time `date_of_birth:"json" binding:"required"`
}
