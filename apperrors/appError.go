package apperrors

import "net/http"

type AppError struct {
	Code    int
	Message string
}

func (ap *AppError) Error() string {
	return ap.Message
}

var (
	ErrUserNotFound = &AppError{
		Code:    http.StatusNotFound,
		Message: "user not found",
	}

	ErrEmailIsUsed = &AppError{
		Code:    http.StatusConflict,
		Message: "this email is used",
	}

	ErrWrongUserPassword = &AppError{
		Code:    http.StatusUnauthorized,
		Message: "wrong user password",
	}

	ErrInvalidCredentials = &AppError{
		Code:    http.StatusUnauthorized,
		Message: "invalid email or password",
	}

	ErrHashing = &AppError{
		Code:    http.StatusInternalServerError,
		Message: "error while hashing",
	}

	ErrInvalidStatusInput = &AppError{
		Code:    http.StatusBadRequest,
		Message: "invalid status input",
	}

	ErrUnauthorized = &AppError{
		Code:    http.StatusUnauthorized,
		Message: "not allowed",
	}

	ErrWrongConfirmationPassword = &AppError{
		Code:    http.StatusBadRequest,
		Message: "wrong confirmation password",
	}

	ErrUnsupportedAuthenticationMethod = &AppError{
		Code:    http.StatusUnauthorized,
		Message: "not supported",
	}

	ErrBadRequestData = &AppError{
		Code:    http.StatusBadRequest,
		Message: "bad request",
	}

	ErrMovieNotFound = &AppError{
		Code:    http.StatusNotFound,
		Message: "movie not found",
	}

	ErrStarNotFound = &AppError{
		Code:    http.StatusNotFound,
		Message: "star not found",
	}

	ErrProducerNotFound = &AppError{
		Code:    http.StatusNotFound,
		Message: "producer not found",
	}

	ErrMovieAlreadyExists = &AppError{
		Code:    http.StatusConflict,
		Message: "movie already exists",
	}

	ErrEmptyUpdate = &AppError{
		Code:    http.StatusBadRequest,
		Message: "no fields provided for update",
	}

	ErrReviewAlreadyExists = &AppError{
		Code:    http.StatusConflict,
		Message: "review already exists",
	}

	ErrReviewNotFound = &AppError{
		Code:    http.StatusNotFound,
		Message: "review not found",
	}

	ErrActorAlreadyExists = &AppError{
		Code:    http.StatusConflict,
		Message: "actor already exists",
	}

	ErrActorNotFound = &AppError{
		Code:    http.StatusNotFound,
		Message: "actor not found",
	}

	ErrGenerateRefreshToken = &AppError{
		Code:    http.StatusInternalServerError,
		Message: "error while generating refresh token",
	}

	ErrGenerateAccessToken = &AppError{
		Code:    http.StatusInternalServerError,
		Message: "error while generating access token",
	}

	ErrInvalidRefreshToken = &AppError{
		Code:    http.StatusUnauthorized,
		Message: "invalid refresh token",
	}

	ErrRefreshTokenExpired = &AppError{
		Code:    http.StatusUnauthorized,
		Message: "refresh token expired",
	}

	ErrMovieAlreadyInWatchList = &AppError{
		Code:    http.StatusConflict,
		Message: "movie is already in watch list",
	}

	ErrWatchListEntryNotFound = &AppError{
		Code:    http.StatusNotFound,
		Message: "watch list entry not found",
	}

	ErrInvalidRatingInput = &AppError{
		Code:    http.StatusBadRequest,
		Message: "Invalid rating input",
	}
)
