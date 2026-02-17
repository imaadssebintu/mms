package dashboard

import (
	"mms/app/routes/auth"

	"github.com/gofiber/fiber/v2"
)

func SetupDashboardRoutes(app *fiber.App) {
	app.Get("/dashboard", auth.AuthMiddleware, ShowDashboard)

	api := app.Group("/api")
	api.Use(auth.AuthMiddleware)
	api.Get("/dashboard-stats", GetDashboardStatsAPI)
	api.Get("/car-inventory", GetCarInventoryAPI)
	api.Get("/active-debts", GetActiveDebtsAPI)
}

func ShowDashboard(c *fiber.Ctx) error {
	user := c.Locals("user")
	return c.Render("dashboard/index", fiber.Map{
		"Title":       "Dashboard",
		"CurrentPage": "dashboard",
		"User":        user,
	})
}
