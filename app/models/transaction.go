package models

import (
	"time"
)

type Transaction struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	CarID     *string   `json:"car_id"`
	ClientID  *string   `json:"client_id"`
	Amount    float64   `json:"amount"`
	Date      time.Time `json:"date"`
	CreatedAt time.Time `json:"created_at"`
}
