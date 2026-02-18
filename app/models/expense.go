package models

import (
	"time"
)

type Expense struct {
	ID          string    `json:"id"`
	Date        time.Time `json:"date"`
	Amount      float64   `json:"amount"`
	Category    string    `json:"category"`
	Description string    `json:"description"`
	CarID       *string   `json:"car_id"`
	CarPlate    string    `json:"car_plate"`
	ClientID    *string   `json:"client_id"`
	ClientName  string    `json:"client_name"`
	CreatedAt   time.Time `json:"created_at"`
}
