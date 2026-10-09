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
	Id     int
	userID int
	review string
	Rating float64
}

type Stars struct {
	Id        int
	Name      string
	birthDate time.Time
	homeLand  string
}

type Claims struct {
	UserId int
	Role   string
	jwt.RegisteredClaims
}
