package xenditpayment

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/imnzr/sim-service-project/internal/repository"
	"github.com/imnzr/sim-service-project/models"
	"github.com/xendit/xendit-go/v7"
	"github.com/xendit/xendit-go/v7/invoice"
)

type XenditPayment interface {
	CreateOrderAndStartPayment(ctx *fiber.Ctx, order *models.SimOrder) (string, error)
	FindyByInvoiceId(ctx *fiber.Ctx, invoiceId string) (*models.SimOrder, error)
	UpdateStatusByInvoiceId(ctx *fiber.Ctx, invoiceId, status string) error
}

type XenditPaymentImplement struct {
	simOrderRepo repository.SimOrderRepository
	userRepo     repository.UserRepository
}

type XenditPaymentRequest struct{}

func NewXenditPayment(userRepo repository.UserRepository, simOrderRepo repository.SimOrderRepository) XenditPayment {
	return &XenditPaymentImplement{
		simOrderRepo: simOrderRepo,
		userRepo:     userRepo,
	}
}

func (x *XenditPaymentImplement) CreateOrderAndStartPayment(ctx *fiber.Ctx, order *models.SimOrder) (string, error) {
	UserId := ctx.Locals("userID").(uint)

	user, err := x.userRepo.GetUserById(context.Background(), UserId)
	if err != nil {
		return "", fmt.Errorf("failed to get user by ID: %w", err)
	}

	order.Status = "PENDING"
	order.InvoiceId = fmt.Sprintf("INV-%d", time.Now().UnixNano())

	xenditKey := os.Getenv("XENDIT_API_KEY")
	if xenditKey == "" {
		return "", fmt.Errorf("XENDIT_API_KEY environment variable is not set")
	}

	xenditClient := xendit.NewClient(xenditKey)

	orderId, err := x.simOrderRepo.Create(context.Background(), order)
	if err != nil {
		return "", fmt.Errorf("failed to create order in repository: %w", err)
	}
	order.Id = orderId

	// Build invoice data
	email := user.Email
	description := fmt.Sprintf("Payment for %s service in %s by %s", order.Service, order.Country, order.Operator)

	createInvoiceRequest := *invoice.NewCreateInvoiceRequest(order.InvoiceId, order.PriceSell)
	createInvoiceRequest.PayerEmail = &email
	createInvoiceRequest.Description = &description

	resp, res, errXendit := xenditClient.InvoiceApi.CreateInvoice(context.Background()).
		CreateInvoiceRequest(createInvoiceRequest).
		Execute()

	if errXendit != nil {
		// Hanya masuk sini jika err benar-benar tidak nil
		fmt.Fprintf(os.Stderr, "Error creating invoice: %v\n", errXendit)

		if res != nil {
			fmt.Fprintf(os.Stderr, "HTTP response: %v\n", res)
		}
		return "", fmt.Errorf("failed to create invoice: %w", errXendit)
	}

	log.Printf("Invoice created successfully: %s\n", resp.InvoiceUrl)
	return resp.InvoiceUrl, nil
}

// FindyByInvoiceId implements XenditPayment.
func (x *XenditPaymentImplement) FindyByInvoiceId(ctx *fiber.Ctx, invoiceId string) (*models.SimOrder, error) {
	order, err := x.simOrderRepo.GetByInvoiceId(context.Background(), invoiceId)
	if err != nil {
		return nil, fmt.Errorf("failed to find order by invoice ID: %w", err)
	}
	if order == nil {
		return nil, fmt.Errorf("order not found for invoice ID: %s", invoiceId)
	}

	return order, nil
}

// UpdateStatusByInvoiceId implements XenditPayment.
func (x *XenditPaymentImplement) UpdateStatusByInvoiceId(ctx *fiber.Ctx, invoiceId string, status string) error {
	order, err := x.simOrderRepo.GetByInvoiceId(ctx.Context(), invoiceId)
	if err != nil {
		return fmt.Errorf("failed to find order by invoice ID: %w", err)
	}
	if order == nil {
		return fmt.Errorf("order not found for invoice ID: %s", invoiceId)
	}
	if err := x.simOrderRepo.UpdateStatusByInvoiceId(ctx.Context(), invoiceId, status); err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	log.Printf("Order status updated successfully: %s to %s\n", invoiceId, status)
	return nil
}
