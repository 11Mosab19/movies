package util

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"math/big"
	"movies/apperrors"
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
