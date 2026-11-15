package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Movie struct {
	Id             int
	Title          string
	StoryLine      string
	ProducerName   string
	ImdbRating     *float64
	InternalRating float64
	Status         string
	ReleaseDate    *time.Time
	Category       string
	PosterUrl      string
	PosterPublicID string
}

type WatchListMovie struct {
	MovieID     int
	Title       string
	PosterURL   string
	SetAt       time.Time
	WatchStatus string
}

type User struct {
	Id             int
	FullName       string
	HashedPassword string `json:"-"`
	Email          string
	Role           string
}

type UserProfile struct {
	FullName     string
	TotalReviews int
	Reviews      []UserProfileReview
}

type UserProfileReview struct {
	ID        int
	Rating    float64
	Review    string
	MovieURL  string
	MovieName string
}

type Review struct {
	Id      int
	UserID  int
	Review  string
	Rating  float64
	MovieId int
}

type UpdateReviewData struct {
	Review *string
	Rating *float64
}

type Actor struct {
	Id            int
	Name          string
	BirthDate     time.Time
	HomeLand      string
	ImageURL      string
	ImagePublicID string
	Movies        []Movie
}

type ActorSummary struct {
	Id       int
	Name     string
	ImageURL string
}

type Claims struct {
	UserId int
	Role   string
	jwt.RegisteredClaims
}

type ReviewsForUserWithMovie struct {
	Review    string
	Rating    float64
	Title     string
	PosterUrl string
}

type ReviewForMovie struct {
	UserName string
	Rating   float64
	Review   string
}
type DetailedReview struct {
	ReviewID   int
	MovieID    int
	UserID     int
	UserName   string
	MovieTitle string
	IMDBRating *float64
	UserRating float64
	PosterURL  string
	Review     string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type MovieWithActors struct {
	Movie  Movie
	Actors []ActorInMovie
}

type ActorInMovie struct {
	ActorId       int
	Name          string
	MovieId       int
	Role          string
	CharacterName string
}

type ActorWithMovies struct {
	Actor  Actor
	Movies []MovieInActor
}

type MovieInActor struct {
	MovieId       int
	Title         string
	Role          string
	CharacterName string
}

type UpdateUserData struct {
	NewHashedPassword string
	FullName          string
	Email             string
}

type RefreshToken struct {
	Id          int
	UserID      int
	ExpiresAt   time.Time
	Created_at  time.Time
	HashedToken string
}
