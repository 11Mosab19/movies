package middlewares

import (
	"movies/apperrors"
	"movies/services"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	AS *services.AuthService
}

func (AM *AuthMiddleware) LoginRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.Error(apperrors.ErrUnauthorized)
			c.Abort()
			return
		}
		if !strings.HasPrefix(header, "Bearer ") {
			c.Error(apperrors.ErrUnsupportedAuthenticationMethod)
			c.Abort()
			return
		}
		tokenString := strings.TrimPrefix(header, "Bearer ")
		claims, err := AM.AS.VerifyToken(tokenString)
		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}
		c.Set("userId", claims.UserId)
		c.Set("Role", claims.Role)
		c.Next()

	}
}

func (AM *AuthMiddleware) RequireRole(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleValue, exists := c.Get("Role")
		if !exists {
			c.Error(apperrors.ErrUnauthorized)
			c.Abort()
			return
		}

		role, ok := roleValue.(string)
		if !ok || role == "" {
			c.Error(apperrors.ErrUnauthorized)
			c.Abort()
			return
		}

		if role != requiredRole {
			c.Error(apperrors.ErrForbidden)
			c.Abort()
			return
		}

		c.Next()
	}
}
