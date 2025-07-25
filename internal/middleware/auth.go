package middleware

import (
	"fmt"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/imnzr/sim-service-project/config"
	"github.com/imnzr/sim-service-project/internal/service"
)

// Authmiddleware memverifikasi JWT dan menambahkan user ID ke konteks
func AuthMiddleware(userService service.UserService, cfg config.AppConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "authorization header is missing",
			})
		}

		// Periksa format bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid authorization header format. Expected 'Bearer <token>'",
			})
		}
		tokenString := parts[1]

		// validasi menggunakan user service (yang berisi logika validasi token)
		claims, err := userService.ValidateToken(tokenString)
		if err != nil {
			log.Printf("JWT validation failed: %v", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": fmt.Sprintf("invalid or expired token: %v", err),
			})
		}

		// ekstrak user_id dari claims
		userIdFloat, ok := claims["user_id"].(float64)
		if !ok {
			log.Printf("user id not found or invalid type in token claims")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid user ID in token claimms",
			})
		}
		userID := uint(userIdFloat)

		// simpan user id di local fiber context agar bisa di akses di handler
		c.Locals("userID", userID)

		// lanjutkan ke handler berikutnya jika token valid
		return c.Next()
	}
}
