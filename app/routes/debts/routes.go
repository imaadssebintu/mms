package debts

import (
	"mms/app/routes/auth"

	"github.com/gofiber/fiber/v2"
)

func SetupDebtRoutes(app *fiber.App) {
	// UI Routes
	app.Get("/debts", auth.AuthMiddleware, ShowDebts)

	// API Routes
	debts := app.Group("/api/debts")
	debts.Use(auth.AuthMiddleware)

	debts.Get("/", GetDebtsAPI)
}
