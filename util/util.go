package util

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"math/big"
	"movies/apperrors"
	"movies/dto"
	"strings"
)

const (
	alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	length   = 32
)

func GenerateRefreshTokenString() (string, error) {
	var sb strings.Builder
	alphabetLen := big.NewInt(int64(len(alphabet)))

	for i := 0; i < length; i++ {
		index, err := rand.Int(rand.Reader, alphabetLen)
		if err != nil {
			return "", apperrors.ErrGenerateRefreshToken
		}
		sb.WriteByte(alphabet[index.Int64()])
	}
	return sb.String(), nil
}

func HashRefreshToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func UpdateMovieReqToMap(req dto.UpdateMovieReq) map[string]interface{} {
	fields := make(map[string]interface{})

	if req.Title != nil {
		fields["title"] = *req.Title
	}

	if req.StoryLine != nil {
		fields["story_line"] = *req.StoryLine
	}

	if req.ProducerName != nil {
		fields["producer_name"] = *req.ProducerName
	}

	if req.Category != nil {
		fields["category"] = *req.Category
	}

	if req.ReleaseDate != nil {
		fields["release_date"] = *req.ReleaseDate
	}

	if req.PosterURL != nil {
		fields["poster_url"] = *req.PosterURL
	}

	if req.Status != nil {
		fields["status"] = *req.Status
	}

	if req.ImdbRating != nil {
		fields["imdb_rating"] = *req.ImdbRating
	}

	return fields
}
