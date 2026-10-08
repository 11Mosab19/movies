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

	HashingErr = &AppError{
		Code:    http.StatusInternalServerError,
		Message: "error while hashing",
	}

	ErrInvalidStatusInput = &AppError{
		Code:    http.StatusBadRequest,
		Message: "invalid update statue",
	}

	ErrUnauthorized = &AppError{
		Code:    http.StatusUnauthorized,
		Message: "not allowed",
	}

	ErrWrongConfirmationPassword = &AppError{
		Code:    http.StatusBadRequest,
		Message: "wrong confirmation password",
	}

	ErrNotSupportedAuthenticateMethod = &AppError{
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
)
