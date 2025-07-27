package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/imnzr/sim-service-project/internal/service"
)

type ProductController interface {
	GetProductAvailable(controller *fiber.Ctx) error
	SyncFromSimServices(controller *fiber.Ctx) error
}

type ProductControllerImplementation struct {
	ProductService service.ProductService
}

func NewProductController(productService service.ProductService) ProductController {
	return &ProductControllerImplementation{
		ProductService: productService,
	}
}

// GetProductAvailable implements ProductController.
func (p *ProductControllerImplementation) GetProductAvailable(controller *fiber.Ctx) error {
	service := controller.Params("service")
	country := controller.Params("country")
	operator := controller.Params("operator")

	if service == "" || country == "" || operator == "" {
		return controller.Status(500).JSON(fiber.Map{
			"message": "service, country, and operator are required",
		})
	}

	products, err := p.ProductService.GetProductAvailable(service, country, operator)
	if err != nil {
		return controller.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return controller.Status(200).JSON(fiber.Map{
		"service":  service,
		"country":  country,
		"products": products,
	})
}

// SyncFromSimServices implements ProductController.
func (p *ProductControllerImplementation) SyncFromSimServices(controller *fiber.Ctx) error {
	err := p.ProductService.SyncFromSimServices(controller.Context())
	if err != nil {
		return controller.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return controller.Status(200).JSON(fiber.Map{
		"message": "product berhasil disinkronasi dari service",
	})
}
