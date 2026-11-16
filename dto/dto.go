package dto

import (
	"mime/multipart"
	"time"
)

// Authentication DTOs are used by registration, login, and refresh-token endpoints.
type RegisterUserReq struct {
	Email                string `json:"email" binding:"required,email"`
	Password             string `json:"password" binding:"required,min=8"`
	ConfirmationPassword string `json:"confirmation_password" binding:"required,eqfield=Password"`
	FullName             string `json:"full_name" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type SetPasswordRequest struct {
	NewPassword          string `json:"new_password" binding:"required,min=8"`
	ConfirmationPassword string `json:"confirmation_password" binding:"required,eqfield=NewPassword"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type LogoutRequest = RefreshTokenRequest

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
}

// Token is used internally while creating an access token.
type Token struct {
	Key string `json:"token"`
}

// User DTOs are used by account and profile endpoints.
type UpdateUserRequest struct {
	NewPassword string `json:"new_password" binding:"omitempty,min=8"`
	OldPassword string `json:"old_password"`
	FullName    string `json:"full_name"`
	Email       string `json:"email" binding:"omitempty,email"`
}

type DeleteUserReq struct {
	Password string `json:"password" binding:"required"`
}

type UserRes struct {
	Id       int    `json:"id"`
	Role     string `json:"role"`
	FullName string `json:"full_name"`
}

type UserProfileResponse struct {
	FullName     string              `json:"full_name"`
	Reviews      []UserProfileReview `json:"reviews"`
	TotalReviews int                 `json:"total_reviews"`
}

type UserProfileReview struct {
	ID        int     `json:"id"`
	Rating    float64 `json:"rating"`
	Review    string  `json:"review"`
	MovieURL  string  `json:"movie_url"`
	MovieName string  `json:"movie_name"`
}

// Movie DTOs are used by multipart movie endpoints and movie responses.
type MovieReq struct {
	Title        string                `form:"title" binding:"required"`
	StoryLine    string                `form:"story_line" binding:"required"`
	ImdbRating   *float64              `form:"imdb_rating"`
	Status       string                `form:"status" binding:"required"`
	ReleaseDate  time.Time             `form:"release_date" time_format:"2006-01-02" binding:"required"`
	ProducerName string                `form:"producer_name" binding:"required"`
	Category     string                `form:"category" binding:"required"`
	Image        *multipart.FileHeader `form:"image" binding:"required"`
}

type UpdateMovieReq struct {
	Title        *string               `form:"title"`
	StoryLine    *string               `form:"story_line"`
	ProducerName *string               `form:"producer_name"`
	Category     *string               `form:"category"`
	ReleaseDate  *time.Time            `form:"release_date" time_format:"2006-01-02"`
	Status       *string               `form:"status"`
	ImdbRating   *float64              `form:"imdb_rating"`
	Image        *multipart.FileHeader `form:"image"`
}

type MovieDetailsResponse struct {
	Id             int               `json:"id"`
	Title          string            `json:"title"`
	StoryLine      string            `json:"story_line"`
	ProducerName   string            `json:"producer_name"`
	ImdbRating     *float64          `json:"imdb_rating"`
	InternalRating float64           `json:"internal_rating"`
	Status         string            `json:"status"`
	ReleaseDate    *time.Time        `json:"release_date"`
	Category       string            `json:"category"`
	PosterUrl      string            `json:"poster_url"`
	Actors         []ActorInMovieRes `json:"actors"`
}

type MovieSummaryResponse struct {
	Id             int      `json:"id"`
	Title          string   `json:"title"`
	PosterUrl      string   `json:"poster_url"`
	Category       string   `json:"category"`
	ImdbRating     *float64 `json:"imdb_rating"`
	InternalRating float64  `json:"internal_rating"`
}

type ActorInMovieRes struct {
	ActorId       int    `json:"actor_id"`
	Name          string `json:"name"`
	ImageURL      string `json:"image_url"`
	Role          string `json:"role"`
	CharacterName string `json:"character_name"`
}

// Actor DTOs are used by multipart actor endpoints and actor responses.
type AddActorReq struct {
	Name        string                `form:"name" binding:"required"`
	HomeLand    string                `form:"home_land" binding:"required"`
	DateOfBirth time.Time             `form:"date_of_birth" time_format:"2006-01-02" binding:"required"`
	Image       *multipart.FileHeader `form:"image" binding:"required"`
}

type UpdateActorReq struct {
	Name        *string               `form:"name"`
	HomeLand    *string               `form:"home_land"`
	DateOfBirth *time.Time            `form:"date_of_birth" time_format:"2006-01-02"`
	Image       *multipart.FileHeader `form:"image"`
}

type ActorDetailsResponse struct {
	Id        int                    `json:"id"`
	Name      string                 `json:"name"`
	BirthDate time.Time              `json:"date_of_birth"`
	HomeLand  string                 `json:"home_land"`
	ImageURL  string                 `json:"image_url"`
	Movies    []MovieInActorResponse `json:"movies"`
}

type ActorSummaryResponse struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
}

type MovieInActorResponse struct {
	MovieId       int    `json:"movie_id"`
	Title         string `json:"title"`
	Role          string `json:"role"`
	CharacterName string `json:"character_name"`
}

// Watch list DTOs describe movies returned for a user's watch list.
type WatchListMovieResponse struct {
	MovieID   int       `json:"movie_id"`
	Title     string    `json:"title"`
	PosterURL string    `json:"poster_url"`
	SetAt     time.Time `json:"set_at"`
	Status    string    `json:"status"`
}

// Review DTOs are used by review endpoints and review responses.
type CreateReviewRequest struct {
	Review string   `json:"review" binding:"required"`
	Rating *float64 `json:"rating" binding:"required,gte=0,lte=10"`
}

type UpdateReviewRequest struct {
	Review *string  `json:"review"`
	Rating *float64 `json:"rating" binding:"omitempty,gte=0,lte=10"`
}

type UserReviewResponse struct {
	Review     string  `json:"review"`
	Rating     float64 `json:"rating"`
	MovieTitle string  `json:"movie_title"`
	PosterURL  string  `json:"poster_url"`
}

type MovieReviewResponse struct {
	UserName string  `json:"user_name"`
	Rating   float64 `json:"rating"`
	Review   string  `json:"review"`
}

type DetailedReviewResponse struct {
	ReviewID   int       `json:"review_id"`
	MovieID    int       `json:"movie_id"`
	UserID     int       `json:"user_id"`
	UserName   string    `json:"user_name"`
	MovieTitle string    `json:"movie_title"`
	IMDBRating *float64  `json:"imdb_rating"`
	UserRating float64   `json:"user_rating"`
	PosterURL  string    `json:"poster_url"`
	Review     string    `json:"review"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
