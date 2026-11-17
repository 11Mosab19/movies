package middlewares

import (
	"errors"
	"movies/apperrors"
	"movies/logging"

	"github.com/gin-gonic/gin"
)

func ErrorMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}
		if c.Writer.Written() {
			return
		}

		err := c.Errors.Last().Err

		var appErr *apperrors.AppError
		if errors.As(err, &appErr) && appErr.Code < 500 && !hasInternalCause(err) {
			c.JSON(appErr.Code, gin.H{"message": appErr.Message})
			return
		}

		logging.Error("http", c.Request.Method+" "+c.FullPath(), err)
		c.JSON(apperrors.ErrInternalServerError.Code, gin.H{"message": apperrors.ErrInternalServerError.Message})
	}
}

func hasInternalCause(err error) bool {
	if err == nil {
		return false
	}

	if _, ok := err.(*apperrors.AppError); ok {
		return false
	}

	if joined, ok := err.(interface{ Unwrap() []error }); ok {
		for _, wrapped := range joined.Unwrap() {
			if hasInternalCause(wrapped) {
				return true
			}
		}
		return false
	}

	if wrapped, ok := err.(interface{ Unwrap() error }); ok {
		return hasInternalCause(wrapped.Unwrap())
	}

	return true
}
