package models

import (
	"time"
)

type Debt struct {
	ID               string    `json:"id"`
	ClientID         string    `json:"client_id"`
	CarID            string    `json:"car_id"`
	RemainingBalance float64   `json:"remaining_balance"`
	PaymentDeadline  time.Time `json:"payment_deadline"`
	NextPaymentDate  time.Time `json:"next_payment_date"`
	Status           string    `json:"status"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
