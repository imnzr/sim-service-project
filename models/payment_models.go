package models

import "time"

// Definisikan status order SIM
type SimOrderStatus string

// Payment Status
const (
	OrderStatusPendingPayment   SimOrderStatus = "PENDING_PAYMENT"
	OrderStatusPaymentFailed    SimOrderStatus = "PAYMENT_FAILED"
	OrderStatusPaid             SimOrderStatus = "PAID"
	OrderStatusProcessingSIM    SimOrderStatus = "PROCESSING_SIM"
	OrderStatusComplete         SimOrderStatus = "COMPLETED"
	OrderStatusFailedSIMService SimOrderStatus = "FAILED_SIM_SERVICE"
	OrderStatusCanceled         SimOrderStatus = "CANCELED"
)

type SimOrder struct {
	Id                uint           `json:"id"`
	UserId            uint           `json:"user_id"`
	Service           string         `json:"service"`
	Country           string         `json:"country"`
	Operator          string         `json:"operator"`
	Price             float64        `json:"price"`
	InvoiceId         string         `json:"invoice_id"`
	SimOrderServiceId uint           `json:"sim_order_service_id"`
	PhoneNumber       string         `json:"phone_number"`
	OTP               string         `json:"otp"`
	Status            SimOrderStatus `json:"status"`
	ErrorMessage      string         `json:"error_message"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
}

type GetUserOrdersRequest struct {
	Status string `query:"status"`
	Limit  int    `query:"limit"`
	Offset int    `query:"offset"`
}

type Product struct {
	Id           uint      `json:"id"`
	Service      string    `json:"service"`
	Country      string    `json:"country"`
	Operator     string    `json:"operator"`
	PriceDefault float64   `json:"price_default"`
	PriceSell    float64   `json:"price_sell"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type ProductResponse struct {
	Id        uint    `json:"id"`
	Service   string  `json:"service"`
	Country   string  `json:"country"`
	Operator  string  `json:"operator"`
	PriceSell float64 `json:"price_sell"`
	IsActive  bool    `json:"is_active"`
}

// Create invoice request untuk membuat invoice baru
type CreateInvoiceRequest struct {
	ReferenceId string `json:"reference_id"`
	TotalAmount string `json:"total_amount"`
	Email       string `json:"email"`
	Description string `json:"description"`
}

// Invoice response adalah response setelah invoice berhasil dibuat
type InvoiceRequest struct {
	ID          string    `json:"id"`
	ReferenceId string    `json:"reference_id"`
	PaymentLink string    `json:"payment_link"`
	Status      string    `json:"status"`
	Amount      float64   `json:"amount"`
	PayerEmail  string    `json:"payer_email"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	ExpiresAt   time.Time `json:"expires_at"`
}

// XenditCallbackRequest adalah struktur untuk menerima callback dari xendit
type XenditCallbackRequest struct {
	ID                       string   `json:"id"`
	ExternalId               string   `json:"external_id"`
	UserId                   string   `json:"user_id"`
	Status                   string   `json:"status"`
	Amount                   float64  `json:"amount"`
	PayerEmail               string   `json:"payer_email"`
	Description              string   `json:"description"`
	InvoiceURL               string   `json:"invoice_url"`
	AvailablePaymentMethod   []string `json:"available_payment_method"`
	CallbackVirtualAccountID *string  `json:"callback_virtual_account_id"`
}
