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
	CreatedAt   time.Time `json:"created_at"`
}
