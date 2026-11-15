package main

import (
	"fmt"
	"log"
	"movies/database"
	"movies/handlers"
	"movies/middlewares"
	"movies/repositories"
	"movies/routes"
	"movies/services"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	db := &database.Database{}
	if err := db.ConnectDatabase(); err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	defer db.CloseConnection()

	userRepo := &repositories.UserRepo{DB: db}
	authRepo := &repositories.AuthRepo{DB: db}
	userService := &services.UserService{UP: userRepo}
	authService := &services.AuthService{US: userRepo, AR: authRepo}
	authHandler := &handlers.AuthHandler{
		AS: authService,
		US: userService,
	}

	router := gin.Default()
	router.Use(middlewares.ErrorMiddleWare())
	routes.RegisterAuthRoutes(router, authHandler)

	if err := router.Run(); err != nil {
		return fmt.Errorf("run HTTP server: %w", err)
	}
	return nil
}
