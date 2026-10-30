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

type User struct {
	Id             int
	FullName       string
	HashedPassword string `json:"-"`
	Email          string
	Role           string
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

type Claims struct {
	UserId int
	Role   string
	jwt.RegisteredClaims
}

type ReviewWithMovie struct {
	Review    string
	Rating    float64
	Title     string
	PosterUrl string
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
