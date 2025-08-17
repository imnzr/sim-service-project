package xenditpayment

import "time"

// type ChannelPropertiesModels struct {
// 	CustomerName       string `json:"customer_name"`
// 	SuccessRedirectURL string `json:"success_redirect_url"`
// 	FailureRedirectURL string `json:"failure_redirect_url"`
// 	CancelRedirectURL  string `json:"cancel_redirect_url"`
// }

type CreateInvoiceRequest struct {
	ExternalId  string  `json:"external_id"`
	Amount      float64 `json:"amount"`
	PayerEmail  string  `json:"payer_email"`
	Description string  `json:"description"`
}

type ResponsePayment struct {
	Success bool                 `json:"success"`
	Data    DataResponsePayment  `json:"data"`
	Order   OrderResponsePayment `json:"order"`
}

type DataResponsePayment struct {
	CheckoutURL string `json:"checkout_url"`
	InvoiceId   string `json:"invoice_id"`
}

type OrderResponsePayment struct {
	Id        int       `json:"id"`
	UserId    int       `json:"user_id"`
	Service   string    `json:"service"`
	Country   string    `json:"country"`
	Operator  string    `json:"operator"`
	PriceSell float64   `json:"price_sell"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
