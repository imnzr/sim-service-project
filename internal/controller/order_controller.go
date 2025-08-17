package controller

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	xenditpayment "github.com/imnzr/sim-service-project/internal/payment_gateway/xendit_payment"
	"github.com/imnzr/sim-service-project/internal/repository"
	"github.com/imnzr/sim-service-project/internal/service"
	"github.com/imnzr/sim-service-project/models"
)

type OrderController interface {
	CreateOrder(ctx *fiber.Ctx) error
	HandleWebhook(ctx *fiber.Ctx) error
	// CheckOrderServiceStatus(ctx *fiber.Ctx) error
}

type OrderControllerImplement struct {
	simOrderRepo         repository.SimOrderRepository
	simOrderService      service.OrderService
	XenditPaymentService xenditpayment.XenditPayment
}

func NewOrderController(simOrderRepository repository.SimOrderRepository, simOrderService service.OrderService, xenditPaymentService xenditpayment.XenditPayment) OrderController {
	return &OrderControllerImplement{
		simOrderRepo:         simOrderRepository,
		XenditPaymentService: xenditPaymentService,
		simOrderService:      simOrderService,
	}
}

// // CheckOrderServiceStatus implements OrderController.
// func (o *OrderControllerImplement) CheckOrderServiceStatus(ctx *fiber.Ctx) error {
// 	orderIdParams := ctx.Params("orderId")
// 	orderId, err := strconv.Atoi(orderIdParams)
// 	if err != nil {
// 		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"error": "Invalid order ID",
// 		})
// 	}

// 	statusResp, err := o.simOrderService.CheckOrderStatus(context.Background(), orderId)
// 	if err != nil {
// 		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"error": err.Error(),
// 		})
// 	}

// 	return ctx.JSON(statusResp)
// }

// CreateOrder implements OrderController.
func (o *OrderControllerImplement) CreateOrder(ctx *fiber.Ctx) error {
	userID := ctx.Locals("userID").(uint)

	var req struct {
		Service   string  `json:"service"`
		Country   string  `json:"country"`
		Operator  string  `json:"operator"`
		PriceSell float64 `json:"price_sell"`
	}

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	order := &models.SimOrder{
		UserId:    int(userID),
		Service:   req.Service,
		Country:   req.Country,
		Operator:  req.Operator,
		PriceSell: req.PriceSell,
	}

	createOrder, err := o.XenditPaymentService.CreateOrderAndStartPayment(ctx, order)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	responseWeb := xenditpayment.ResponsePayment{
		Success: true,
		Data: xenditpayment.DataResponsePayment{
			CheckoutURL: createOrder,
			InvoiceId:   order.InvoiceId,
		},
		Order: xenditpayment.OrderResponsePayment{
			Id:        order.Id,
			UserId:    order.UserId,
			Service:   order.Service,
			Country:   order.Country,
			Operator:  order.Operator,
			PriceSell: order.PriceSell,
			Status:    order.Status,
		},
	}
	return ctx.Status(fiber.StatusCreated).JSON(responseWeb)
}

func (o *OrderControllerImplement) HandleWebhook(ctx *fiber.Ctx) error {
	var payload struct {
		Id         string `json:"id"`
		ExternalId string `json:"external_id"`
		Status     string `json:"status"`
	}

	if err := ctx.BodyParser(&payload); err != nil {
		log.Println("❌ Gagal parsing payload:", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	log.Printf("📩 Webhook received | ExternalID: %s | Status: %s\n", payload.ExternalId, payload.Status)

	if payload.Status != "PAID" {
		log.Println("ℹ️ Webhook status bukan PAID, abaikan")
		return ctx.SendStatus(fiber.StatusOK)
	}

	// 1. Ambil order berdasarkan invoice_id
	order, err := o.XenditPaymentService.FindyByInvoiceId(ctx, payload.ExternalId)
	if err != nil {
		log.Println("❌ Gagal cari order:", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal cari order dari invoice_id",
		})
	}

	if order.Status == "PAID" {
		log.Println("ℹ️ Order sudah PAID sebelumnya, skip update.")
		return ctx.SendStatus(fiber.StatusOK)
	}

	// 2. Update status ke PAID
	err = o.XenditPaymentService.UpdateStatusByInvoiceId(ctx, payload.ExternalId, "PAID")
	if err != nil {
		log.Println("❌ Gagal update status:", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal update status order ke PAID",
		})
	}
	log.Println("✅ Status order diupdate ke PAID")

	// 3. Pesan nomor ke 5sim
	resultOrder, err := o.simOrderService.BuyNumberFromService(context.Background(), order.Service, order.Country, order.Operator)
	if err != nil {
		log.Println("❌ Gagal beli nomor dari 5sim:", err)
		return ctx.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": "Gagal beli nomor dari 5sim",
		})
	}
	log.Println("✅ Nomor berhasil dipesan dari 5sim")

	if resultOrder == nil {
		log.Println("⚠️ Tidak ada data untuk disimpan (resultOrder nil)")
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Data 5sim tidak tersedia",
		})
	}
	log.Printf("✅ Menyimpan data SIM: OrderID=%v (%T), SimServiceID=%d", order.Id, order.Id, resultOrder.Id)

	// 4. Simpan hasil pembelian ke database
	err = o.simOrderRepo.AttachSimDataService(context.Background(), order.Id, resultOrder)
	if err != nil {
		log.Println("❌ Gagal simpan data SIM:", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal simpan hasil dari 5sim",
		})
	}
	log.Println("✅ Data 5sim berhasil disimpan ke database")

	return ctx.Status(200).JSON(fiber.Map{
		"message": "Webhook processed successfully",
	})
}
