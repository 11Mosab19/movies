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
	ImdbRating     float64
	InternalRating float64
	Status         string
	ReleaseDate    time.Time
	Category       string
	PosterUrl      string
}

type User struct {
	Id             int
	FullName       string
	HashedPassword string
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

type Actor struct {
	Id        int
	Name      string
	Role      string
	birthDate time.Time
	homeLand  string
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
	Actors []Actor
}
