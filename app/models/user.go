package models

import (
	"time"
)

type User struct {
	ID             string    `json:"id"`
	Email          string    `json:"email"`
	Password       string    `json:"password,omitempty"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	EmailConfirmed bool      `json:"email_confirmed"`
	CompanyName    string    `json:"company_name"`
	Location       string    `json:"location"`
	Phone1         string    `json:"phone1"`
	Phone2         string    `json:"phone2"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
