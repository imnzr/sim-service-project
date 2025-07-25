package main

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/imnzr/sim-service-project/config"
	"github.com/imnzr/sim-service-project/database"
	"github.com/imnzr/sim-service-project/internal/controller"
	"github.com/imnzr/sim-service-project/internal/middleware"
	"github.com/imnzr/sim-service-project/internal/repository"
	"github.com/imnzr/sim-service-project/internal/service"
	"github.com/imnzr/sim-service-project/routes"
)

func main() {
	cfg := config.LoadConfig()

	// Inisialisasi database
	db := database.InitDb(cfg.DatabaseURL)
	defer db.Close()

	// Inisialisasi Repository
	userRepository := repository.NewUserRepository(db)
	// Inisialisasi Service
	userService := service.NewUserService(userRepository, cfg)
	// Inisialisasi Controller
	userController := controller.NewUserController(userService)

	app := fiber.New()

	// Middleware
	authMiddleware := middleware.AuthMiddleware(userService, *cfg)

	// Routes
	routes.SetupUserRoutes(app, userController, authMiddleware)

	log.Printf("Server starting on port %s", cfg.AppPort)
	err := app.Listen(":" + cfg.AppPort)
	if err != nil {
		log.Fatal("Server failed to start: %w", err)
	}
}
