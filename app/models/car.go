package models

import (
	"time"
)

type Car struct {
	ID            string    `json:"id"`
	Make          string    `json:"make"`
	Model         string    `json:"model"`
	Year          int       `json:"year"`
	Color         string    `json:"color"`
	EngineNumber  string    `json:"engine_number"`
	ChassisNumber string    `json:"chassis_number"`
	Price         float64   `json:"price"`
	PurchasePrice float64   `json:"purchase_price"`
	Sold          bool      `json:"sold"`
	ClientID      *string   `json:"client_id"`
	SellerID      *string   `json:"seller_id"`
	NumberPlate   string    `json:"number_plate"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
