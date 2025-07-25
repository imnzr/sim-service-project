package controller

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/imnzr/sim-service-project/internal/service"
	"github.com/imnzr/sim-service-project/models"
)

type UserController interface {
	Register(controller *fiber.Ctx) error
	Login(controller *fiber.Ctx) error
	GetProfile(controller *fiber.Ctx) error
}

type UserControllerImplement struct {
	userService service.UserService
}

// GetProfile implements UserController.
func (u *UserControllerImplement) GetProfile(controller *fiber.Ctx) error {
	// ambil userID dari context yang disisipkan oleh middleware
	userID := controller.Locals("userID")
	if userID == nil {
		controller.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "user id not found in context",
		})
	}
	id, ok := userID.(uint)
	if !ok {
		controller.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "invalid user ID type in context",
		})
	}

	userProfile, err := u.userService.GetUserProfile(controller.Context(), id)
	if err != nil {
		controller.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return controller.Status(fiber.StatusOK).JSON(userProfile)

}

// Login implements UserController.
func (u *UserControllerImplement) Login(controller *fiber.Ctx) error {
	var req models.LoginRequest

	if err := controller.BodyParser(&req); err != nil {
		return controller.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	resp, err := u.userService.Login(controller.Context(), &req)
	if err != nil {
		log.Printf("error logging user: %v", err)
		return controller.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return controller.Status(200).JSON(resp)
}

// Register implements UserController.
func (u *UserControllerImplement) Register(controller *fiber.Ctx) error {
	var user models.RegisterUser
	if err := controller.BodyParser(&user); err != nil {
		return controller.Status(500).JSON(fiber.Map{
			"error": "invalid request",
		})
	}
	if err := u.userService.Register(controller.Context(), &user); err != nil {
		return controller.Status(500).JSON(fiber.Map{
			"error": "error register user",
		})
	}

	return controller.Status(200).JSON(fiber.Map{
		"success": "user created successfully",
	})
}

func NewUserController(user service.UserService) UserController {
	return &UserControllerImplement{
		userService: user,
	}
}
