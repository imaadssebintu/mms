package auth

import (
	"mms/app/config"
	"mms/app/database"
	"mms/app/models"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func SetupAuthRoutes(app *fiber.App) {
	auth := app.Group("/auth")

	auth.Get("/login", ShowLoginPage)
	auth.Post("/login", LoginAPI)
	auth.Post("/logout", LogoutAPI)

	auth.Get("/forgot-password", ShowForgotPasswordPage)
	auth.Post("/forgot-password", ForgotPasswordAPI)
	auth.Get("/reset-password", ShowResetPasswordPage)
	auth.Post("/reset-password", ResetPasswordAPI)

	protected := auth.Group("/me")
	protected.Use(AuthMiddleware)
	protected.Get("/", GetMeAPI)
}

func AuthMiddleware(c *fiber.Ctx) error {
	tokenString := c.Cookies("jwt_token")

	if tokenString == "" {
		authHeader := c.Get("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		}
	}

	if tokenString == "" {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	claims, err := ValidateJWT(tokenString)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid token"})
	}

	// Fetch full user from database to get all fields (CompanyName, etc.)
	db := config.GetDB()
	user, err := database.GetUserByEmail(db, claims.Email)
	if err != nil {
		// Fallback to JWT data if DB fail, but ideally we want full data
		user = &models.User{
			ID:        claims.UserID,
			Email:     claims.Email,
			FirstName: claims.FirstName,
			LastName:  claims.LastName,
		}
	}

	c.Locals("user_id", user.ID)
	c.Locals("user_email", user.Email)
	c.Locals("user", user)

	return c.Next()
}
