package main

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/imnzr/sim-service-project/config"
	"github.com/imnzr/sim-service-project/database"
	"github.com/imnzr/sim-service-project/internal/controller"
	"github.com/imnzr/sim-service-project/internal/middleware"
	xenditpayment "github.com/imnzr/sim-service-project/internal/payment_gateway/xendit_payment"
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
	userProduct := repository.NewProductRepository(db)
	orderRepository := repository.NewOrderRepository(db)

	// Inisialisasi Service
	userService := service.NewUserService(userRepository, cfg)
	orderService := service.NewOrderService(orderRepository, db)
	productService := service.NewProductService(userProduct, *cfg)
	xenditService := xenditpayment.NewXenditPayment(userRepository, orderRepository)

	// Inisialisasi Controller
	userController := controller.NewUserController(userService)
	productController := controller.NewProductController(productService)
	xenditController := controller.NewOrderController(orderRepository, orderService, xenditService)

	app := fiber.New()

	// Middleware
	authMiddleware := middleware.AuthMiddleware(userService, *cfg)

	// Routes
	routes.SetupUserRoutes(app, userController, authMiddleware)
	routes.SetupProductRoutes(app, productController, authMiddleware)
	routes.SetupSimOrderRoutes(app, xenditController, authMiddleware)

	log.Printf("Server starting on port %s", cfg.AppPort)
	err := app.Listen(":" + cfg.AppPort)
	if err != nil {
		log.Fatal("Server failed to start: %w", err)
	}
}
