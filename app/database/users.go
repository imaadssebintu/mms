package database

import (
	"database/sql"
	"mms/app/models"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func GetUserByEmail(db *sql.DB, email string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, email, password, first_name, last_name, email_confirmed, company_name, location, phone1, phone2, created_at, updated_at 
			  FROM users WHERE email = $1`

	err := db.QueryRow(query, email).Scan(
		&user.ID, &user.Email, &user.Password, &user.FirstName, &user.LastName,
		&user.EmailConfirmed, &user.CompanyName, &user.Location, &user.Phone1, &user.Phone2,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}
	return user, nil
}

func UpdateUserPassword(db *sql.DB, userID string, newHash string) error {
	query := "UPDATE users SET password = $1, updated_at = NOW() WHERE id = $2"
	_, err := db.Exec(query, newHash, userID)
	return err
}

func CreatePasswordResetToken(db *sql.DB, email, token string) error {
	// Delete any existing tokens for this email
	_, _ = db.Exec("DELETE FROM password_reset_tokens WHERE email = $1", email)

	query := "INSERT INTO password_reset_tokens (email, token, expires_at) VALUES ($1, $2, NOW() + INTERVAL '1 hour')"
	_, err := db.Exec(query, email, token)
	return err
}

func ValidatePasswordResetToken(db *sql.DB, token string) (string, error) {
	var email string
	query := "SELECT email FROM password_reset_tokens WHERE token = $1 AND used = false AND expires_at > NOW()"
	err := db.QueryRow(query, token).Scan(&email)
	if err != nil {
		return "", err
	}
	return email, nil
}

func MarkPasswordResetTokenAsUsed(db *sql.DB, token string) error {
	query := "UPDATE password_reset_tokens SET used = true WHERE token = $1"
	_, err := db.Exec(query, token)
	return err
}

func CreateUser(db *sql.DB, user *models.User) error {
	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		return err
	}

	query := `INSERT INTO users (email, password, first_name, last_name, company_name, location, phone1, phone2, updated_at)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW())
			  RETURNING id, created_at, updated_at`

	return db.QueryRow(query, user.Email, hashedPassword, user.FirstName, user.LastName, user.CompanyName, user.Location, user.Phone1, user.Phone2).Scan(
		&user.ID, &user.CreatedAt, &user.UpdatedAt,
	)
}
