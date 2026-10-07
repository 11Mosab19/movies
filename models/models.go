package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Movie struct {
	Name           string
	StoryLine      string
	ProducerName   string
	ImdbRating     float64
	InternalRating float64
	Status         string
	ReleaseDate    time.Time
	Category       string
}

type User struct {
	FullName       string
	HashedPassword string
	Email          string
	Role           string
}

type Review struct {
	userID int
	review string
	Rating float64
}

type Stars struct {
	Name      string
	birthDate time.Time
	homeLand  string
}

type Claims struct {
	UserId int
	Role   string
	jwt.RegisteredClaims
}
