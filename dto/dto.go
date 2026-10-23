package dto

import "time"

type RegisterUserReq struct {
	Email                string `json:"email" binding:"required,email"`
	Password             string `json:"password" binding:"required,min=8"`
	ConfirmationPassword string `json:"confirmation_password" binding:"required,eqfield=Password"`
	FullName             string `json:"full_name" binding:"required"`
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

type MovieReq struct {
	Title        string    `json:"title" binding:"required"`
	StoryLine    string    `json:"story_line" binding:"required"`
	ImdbRating   float64   `json:"imdb_rating"`
	Status       string    `json:"status" binding:"required"`
	ReleaseDate  time.Time `json:"release_date" binding:"required"`
	ProducerName string    `json:"producer_name" binding:"required"`
	Category     string    `json:"category" binding:"required"`
	PosterUrl    string    `json:"poster_url" binding:"required,url"`
}

type AddStarReq struct {
	Name        string    `json:"name" binding:"required"`
	HomeLand    string    `json:"home_land" binding:"required"`
	DateOfBirth time.Time `json:"date_of_birth" binding:"required"`
}

type MovieDetailsResponse struct {
	Id             int        `json:"id"`
	Title          string     `json:"title"`
	StoryLine      string     `json:"story_line"`
	ProducerName   string     `json:"producer_name"`
	ImdbRating     *float64   `json:"imdb_rating"`
	InternalRating float64    `json:"internal_rating"`
	Status         string     `json:"status"`
	ReleaseDate    *time.Time `json:"release_date"`
	Category       string     `json:"category"`
	PosterUrl      string     `json:"poster_url"`
	Actors         []ActorRes `json:"actors"`
}

type ActorRes struct {
	ActorId       int    `json:"actor_id"`
	Name          string `json:"name"`
	Role          string `json:"role"`
	CharacterName string `json:"character_name"`
}

type MovieSummaryResponse struct {
	Id             int      `json:"id"`
	Title          string   `json:"title"`
	PosterUrl      string   `json:"poster_url"`
	Category       string   `json:"category"`
	ImdbRating     *float64 `json:"imdb_rating"`
	InternalRating float64  `json:"internal_rating"`
}

type UpdateMovieReq struct {
}
