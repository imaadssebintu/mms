package models

import (
	"time"
)

type Car struct {
	ID          string    `json:"id"`
	Make        string    `json:"make"`
	Model       string    `json:"model"`
	Year        int       `json:"year"`
	Price       float64   `json:"price"`
	Sold        bool      `json:"sold"`
	ClientID    *string   `json:"client_id"`
	NumberPlate string    `json:"number_plate"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
