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
