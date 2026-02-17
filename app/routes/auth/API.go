package auth

import (
	"database/sql"
	"log"
	"mms/app/config"
	"mms/app/database"
	"time"

	"github.com/gofiber/fiber/v2"
)

func ShowLoginPage(c *fiber.Ctx) error {
	return c.Render("auth/login", fiber.Map{
		"Title": "Login",
	}, "layouts/empty")
}

func LoginAPI(c *fiber.Ctx) error {
	type LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		if c.Get("Content-Type") == "application/x-www-form-urlencoded" {
			return c.Redirect("/auth/login?error=Invalid+request")
		}
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	user, err := database.GetUserByEmail(config.GetDB(), req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			if c.Get("Content-Type") == "application/x-www-form-urlencoded" {
				return c.Redirect("/auth/login?error=Invalid+credentials")
			}
			return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Database error"})
	}

	if !CheckPasswordHash(req.Password, user.Password) {
		if c.Get("Content-Type") == "application/x-www-form-urlencoded" {
			return c.Redirect("/auth/login?error=Invalid+credentials")
		}
		return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	// For MMS, we'll assume an 'admin' role for now
	roles := []string{"admin"}

	token, err := GenerateJWT(user.ID, user.Email, user.FirstName, user.LastName, roles)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "jwt_token",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
	})

	if c.Get("Content-Type") == "application/x-www-form-urlencoded" {
		return c.Redirect("/dashboard")
	}

	return c.JSON(fiber.Map{
		"message": "Login successful",
		"user":    user,
	})
}

func LogoutAPI(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "jwt_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	})

	return c.Redirect("/auth/login")
}

func ForgotPasswordAPI(c *fiber.Ctx) error {
	type ForgotPasswordRequest struct {
		Email string `json:"email" form:"email"`
	}

	var req ForgotPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Check if user exists
	_, err := database.GetUserByEmail(config.GetDB(), req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(fiber.Map{"error": "Email not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Database error"})
	}

	// Generate reset token
	resetToken, err := GenerateResetToken()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to generate reset token"})
	}

	// Store token in database
	if err := database.CreatePasswordResetToken(config.GetDB(), req.Email, resetToken); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create reset token"})
	}

	// Send reset email
	if err := SendPasswordResetEmail(req.Email, resetToken, c.Get("Host")); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to send reset email"})
	}

	return c.JSON(fiber.Map{"message": "Password reset link sent to your email"})
}

func ResetPasswordAPI(c *fiber.Ctx) error {
	type ResetPasswordRequest struct {
		Token       string `json:"token" form:"token"`
		NewPassword string `json:"new_password" form:"new_password"`
	}

	var req ResetPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Validate token and get email
	email, err := database.ValidatePasswordResetToken(config.GetDB(), req.Token)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid or expired reset token"})
	}

	// Get user by email
	user, err := database.GetUserByEmail(config.GetDB(), email)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "User not found"})
	}

	// Hash new password
	hashedPassword, err := database.HashPassword(req.NewPassword)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to hash password"})
	}

	// Update password
	if err := database.UpdateUserPassword(config.GetDB(), user.ID, hashedPassword); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update password"})
	}

	// Mark token as used
	if err := database.MarkPasswordResetTokenAsUsed(config.GetDB(), req.Token); err != nil {
		log.Printf("Failed to mark token as used: %v", err)
	}

	return c.JSON(fiber.Map{"message": "Password reset successfully"})
}

func ShowForgotPasswordPage(c *fiber.Ctx) error {
	return c.Render("auth/forgot-password", fiber.Map{
		"Title": "Forgot Password",
	}, "layouts/empty")
}

func ShowResetPasswordPage(c *fiber.Ctx) error {
	token := c.Query("token")
	if token == "" {
		return c.Redirect("/auth/login")
	}
	return c.Render("auth/reset-password", fiber.Map{
		"Title": "Reset Password",
		"Token": token,
	}, "layouts/empty")
}

func GetMeAPI(c *fiber.Ctx) error {
	user := c.Locals("user")
	return c.JSON(fiber.Map{"user": user})
}
