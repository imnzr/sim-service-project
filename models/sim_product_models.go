package models

type SimProduct struct {
	Id           int     `json:"id"`
	Country      string  `json:"country"`
	Service      string  `json:"service"`
	Operator     string  `json:"operator"`
	PriceDefault float64 `json:"price_default"`
	PriceSell    float64 `json:"price_sell"`
}
