package main

import (
	"context"
	"errors"
	"fmt"
	"movies/database"
	"movies/handlers"
	"movies/logging"
	"movies/middlewares"
	"movies/repositories"
	"movies/routes"
	"movies/services"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := run(); err != nil {
		logging.Error("application", "startup", err)
		os.Exit(1)
	}
}

func run() error {
	if err := loadEnvironment(); err != nil {
		return err
	}

	db := &database.Database{}
	if err := db.ConnectDatabase(); err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	defer func() {
		if err := db.CloseConnection(); err != nil {
			logging.Error("database", "close connection", err)
		}
	}()

	cloudinaryService, err := services.NewCloudinaryService()
	if err != nil {
		return fmt.Errorf("initialize Cloudinary: %w", err)
	}

	userRepo := &repositories.UserRepo{DB: db}
	authRepo := &repositories.AuthRepo{DB: db}
	movieRepo := &repositories.MoviesRepo{DB: db}
	actorRepo := &repositories.ActorRepo{DB: db}
	reviewRepo := &repositories.ReviewRepo{DB: db}
	watchListRepo := &repositories.WatchListRepo{DB: db}

	userService := &services.UserService{UP: userRepo}
	authService := &services.AuthService{US: userRepo, AR: authRepo}
	movieService := &services.MovieService{MR: movieRepo, CS: cloudinaryService}
	actorService := &services.ActorService{AR: actorRepo, CS: cloudinaryService}
	reviewService := &services.ReviewService{RR: reviewRepo}
	watchListService := &services.WatchListService{WlR: watchListRepo}

	authHandler := &handlers.AuthHandler{
		AS: authService,
		US: userService,
	}
	userHandler := &handlers.UserHandler{US: userService, AS: authService}
	movieHandler := &handlers.MovieHandler{MS: movieService}
	actorHandler := &handlers.ActorHandler{AS: actorService}
	reviewHandler := &handlers.ReviewHandler{RS: reviewService}
	watchListHandler := &handlers.WatchListHandler{WlS: watchListService}
	authMiddleware := &middlewares.AuthMiddleware{AS: authService}

	router := gin.Default()
	router.Use(middlewares.ErrorMiddleWare())
	routes.RegisterRoutes(
		router,
		authMiddleware,
		authHandler,
		userHandler,
		movieHandler,
		actorHandler,
		reviewHandler,
		watchListHandler,
	)

	return serve(router)
}

func loadEnvironment() error {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("load environment configuration: %w", err)
	}
	return nil
}

func serve(router http.Handler) error {
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	//syscall.SIGTERM => req to close server
	//os.Interrupt => ctrl + c

	signalContext, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	serverErrors := make(chan error, 1)
	go func() {
		serverErrors <- server.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return fmt.Errorf("run HTTP server: %w", err)
	case <-signalContext.Done():
		shutdownContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownContext); err != nil {
			return fmt.Errorf("shutdown HTTP server: %w", err)
		}
		return nil
	}
}
