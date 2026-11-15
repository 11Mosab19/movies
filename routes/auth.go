package routes

import (
	"movies/handlers"

	"github.com/gin-gonic/gin"
)

// RegisterAuthRoutes attaches the authentication endpoints to the root router.
func RegisterAuthRoutes(router gin.IRoutes, authHandler *handlers.AuthHandler) {
	router.POST("/auth/register", authHandler.Register)
	router.POST("/auth/login", authHandler.Login)
	router.POST("/auth/refresh", authHandler.Refresh)
	router.POST("/auth/logout", authHandler.Logout)
}
