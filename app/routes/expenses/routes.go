package expenses

import (
	"mms/app/routes/auth"

	"github.com/gofiber/fiber/v2"
)

func SetupExpenseRoutes(app *fiber.App) {
	// UI Routes
	app.Get("/expenses", auth.AuthMiddleware, ShowExpenses)

	// API Routes
	expenses := app.Group("/api/expenses")
	expenses.Use(auth.AuthMiddleware)

	expenses.Get("/", GetExpensesAPI)
	expenses.Post("/", CreateExpenseAPI)
	expenses.Get("/car/:car_id", GetCarExpensesAPI)
}
