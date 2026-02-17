package main

import (
	"log"
	"mms/app/config"
	"mms/app/database"
	"mms/app/routes/auth"
	"mms/app/routes/cars"
	"mms/app/routes/clients"
	"mms/app/routes/dashboard"
	"mms/app/routes/debts"
	"mms/app/routes/expenses"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"
)

func main() {
	// Initialize Configuration and Database
	config.InitDB()
	db := config.GetDB()

	// Run Migrations
	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	// Initialize Template Engine
	engine := html.New("./app/templates", ".html")
	engine.AddFunc("substr", func(s string, start, end int) string {
		if len(s) == 0 {
			return ""
		}
		if start < 0 {
			start = 0
		}
		if end > len(s) {
			end = len(s)
		}
		if start > end {
			return ""
		}
		return s[start:end]
	})
	engine.Reload(true) // Development mode

	// Create Fiber App
	app := fiber.New(fiber.Config{
		AppName:     "MMS - Motor Management System",
		Views:       engine,
		ViewsLayout: "layouts/main",
	})

	// Add Middleware
	app.Use(logger.New())
	app.Use(recover.New())

	// Register Routes
	auth.SetupAuthRoutes(app)
	dashboard.SetupDashboardRoutes(app)
	cars.SetupCarRoutes(app)
	clients.SetupClientRoutes(app)
	expenses.SetupExpenseRoutes(app)
	debts.SetupDebtRoutes(app)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/auth/login")
	})

	// Start Server
	log.Fatal(app.Listen(":8082"))
}
