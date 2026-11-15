package middlewares

import (
	"errors"
	"movies/apperrors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last()

		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			c.JSON(appErr.Code, gin.H{"message": appErr.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "server error"})
	}
}
