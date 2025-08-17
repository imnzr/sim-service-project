package models

import "time"

type SimOrder struct {
	Id                int       `json:"id"`
	UserId            int       `json:"user_id"`
	Email             string    `json:"email"`
	Service           string    `json:"service"`
	Country           string    `json:"country"`
	Operator          string    `json:"operator"`
	PriceSell         float64   `json:"price_sell"`
	InvoiceId         string    `json:"invoice_id"`
	SimOrderServiceId *int      `json:"sim_order_service_id"`
	PhoneNumber       *string   `json:"phone_number"`
	OTP               *string   `json:"otp"`
	Status            string    `json:"status"`
	ErrorMessage      *string   `json:"error_message"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type CreateSimOrderRequest struct {
	Service   string  `json:"service"`
	Country   string  `json:"country"`
	Operator  string  `json:"operator"`
	PriceSell float64 `json:"price_sell"`
}

type OrderNumberFromService struct {
	Id          int64     `json:"id"`
	Service     string    `json:"service"`
	Country     string    `json:"country"`
	Operator    string    `json:"operator"`
	Price       float64   `json:"price"`
	PhoneNumber string    `json:"phone_number"`
	Date        time.Time `json:"date"`
	CreatedAt   time.Time `json:"created_at"`
	ExpiredAt   time.Time `json:"expired_at"`
}

type ResponsOrderFromService struct {
	Id        int     `json:"id"`
	Phone     string  `json:"phone"`
	Operator  string  `json:"operator"`
	Product   string  `json:"product"`
	Price     float64 `json:"price"`
	Status    string  `json:"status"`
	Expires   string  `json:"expires"`    // atau time.Time jika pakai parsing waktu
	SMS       any     `json:"sms"`        // bisa pakai []SMS jika detailnya tahu
	CreatedAt string  `json:"created_at"` // atau time.Time
	Country   string  `json:"country"`
}

// type ResponsOrderFromService struct {
// 	Id        int       `json:"id"`
// 	CreatedAt time.Time `json:"created_at"`
// 	Phone     string    `json:"phone"`
// 	Product   string    `json:"product"`
// 	Price     string    `json:"price"`
// 	Status    string    `json:"status"`
// 	ExpiredAt time.Time `json:"expired_at"`
// 	SMS       []struct {
// 		CreatedAt time.Time `json:"created_at"`
// 		Date      time.Time `json:"date"`
// 		Sender    string    `json:"sender"`
// 		Text      string    `json:"text"`
// 		Code      string    `json:"code"`
// 	} `json:"sms"`
// 	Forwading       bool   `json:"forwading"`
// 	ForwadingNumber string `json:"forwading_number"`
// 	Country         string `json:"country"`
// }
