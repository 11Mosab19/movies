package routes

import (
	"movies/handlers"
	"movies/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	router *gin.Engine,
	authMiddleware *middlewares.AuthMiddleware,
	authHandler *handlers.AuthHandler,
	userHandler *handlers.UserHandler,
	movieHandler *handlers.MovieHandler,
	actorHandler *handlers.ActorHandler,
	reviewHandler *handlers.ReviewHandler,
	watchListHandler *handlers.WatchListHandler,
) {
	RegisterAuthRoutes(router, authHandler)

	registerPublicRoutes(router, userHandler, movieHandler, actorHandler, reviewHandler)

	userRoutes := router.Group("")
	userRoutes.Use(authMiddleware.LoginRequired())
	registerUserRoutes(userRoutes, userHandler, reviewHandler, watchListHandler)

	adminRoutes := router.Group("")
	adminRoutes.Use(authMiddleware.LoginRequired(), authMiddleware.RequireRole("admin"))
	registerAdminRoutes(adminRoutes, movieHandler, actorHandler)
}

func registerPublicRoutes(
	router gin.IRoutes,
	userHandler *handlers.UserHandler,
	movieHandler *handlers.MovieHandler,
	actorHandler *handlers.ActorHandler,
	reviewHandler *handlers.ReviewHandler,
) {
	router.GET("/users/:id", userHandler.GetUserById)

	router.GET("/movies/search/:title", movieHandler.GetMovieByTitle)
	router.GET("/movies/category/:category", movieHandler.GetMoviesByCategory)
	router.GET("/movies/producer/:producer_name", movieHandler.GetMoviesByProducerName)
	router.GET("/movies/top-rated", movieHandler.GetTopRatedMovies)
	router.GET("/movies/upcoming", movieHandler.GetUpcomingMovies)
	router.GET("/movies/:id", movieHandler.GetMovieById)

	router.GET("/actors", actorHandler.GetActorBySearch)
	router.GET("/actors/:id", actorHandler.GetActorById)

	router.GET("/movies/:id/reviews", reviewHandler.GetMovieReviews)
	router.GET("/reviews/:id", reviewHandler.GetReviewWithId)
	router.GET("/users/:id/reviews", reviewHandler.GetUserReviews)
}

func registerUserRoutes(
	router gin.IRoutes,
	userHandler *handlers.UserHandler,
	reviewHandler *handlers.ReviewHandler,
	watchListHandler *handlers.WatchListHandler,
) {
	router.GET("/users/me", userHandler.GetUserProfile)
	router.PATCH("/users/me", userHandler.UpdateUser)
	router.DELETE("/users/me", userHandler.DeleteUser)

	router.POST("/movies/:id/reviews", reviewHandler.CreateReview)
	router.PATCH("/reviews/:id", reviewHandler.UpdateReview)
	router.DELETE("/reviews/:id", reviewHandler.DeleteReview)

	router.GET("/users/me/watchlist", watchListHandler.GetUserWatchList)
	router.POST("/movies/:id/watchlist", watchListHandler.AddMovieToWL)
	router.PATCH("/movies/:id/watchlist", watchListHandler.SetWatchListStatus)
	router.DELETE("/movies/:id/watchlist", watchListHandler.RemoveMovieFromWL)
}

func registerAdminRoutes(router gin.IRoutes, movieHandler *handlers.MovieHandler, actorHandler *handlers.ActorHandler) {
	router.POST("/movies", movieHandler.AddMovie)
	router.PATCH("/movies/:id", movieHandler.UpdateMovie)
	router.DELETE("/movies/:id", movieHandler.DeleteMovie)

	router.POST("/actors", actorHandler.AddActor)
	router.PATCH("/actors/:id", actorHandler.UpdateActor)
	router.DELETE("/actors/:id", actorHandler.DeleteActor)
}
