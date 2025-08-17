package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/imnzr/sim-service-project/internal/controller"
)

func SetupUserRoutes(app *fiber.App, userControlller controller.UserController, authMiddleware fiber.Handler) {
	authGroup := app.Group("/auth")
	{
		authGroup.Post("/register", userControlller.Register)
		authGroup.Post("/login", userControlller.Login)
		authGroup.Get("/profile", authMiddleware, userControlller.GetProfile)
	}
}

func SetupProductRoutes(app *fiber.App, productController controller.ProductController, authMiddleware fiber.Handler) {
	productGroup := app.Group("/product")
	{
		productGroup.Get("/services", productController.GetProductAvailable)
		productGroup.Post("/sync-services", productController.SyncFromSimServices)
		// purchase
		// status order
		// order otp
	}
}

func SetupSimOrderRoutes(app *fiber.App, controller controller.OrderController, authMiddleware fiber.Handler) {
	orderGroup := app.Group("/sim-order")
	orderGroup.Post("/create", authMiddleware, controller.CreateOrder)
	// orderGroup.Get("/status/:orderId", authMiddleware, controller.CheckOrderServiceStatus)
	orderGroup.Post("/webhook", controller.HandleWebhook)
}
