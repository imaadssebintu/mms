package dashboard

import (
	"mms/app/config"
	"mms/app/database"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func GetDashboardStatsAPI(c *fiber.Ctx) error {
	stats, err := database.GetDashboardStats(config.GetDB())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch stats"})
	}

	// Placeholder for dynamic stats (percentage increases)
	stats["new_cars_last_month"] = 10
	stats["value_increase_percentage"] = 5.2
	stats["revenue_percentage"] = 8.3
	stats["debt_percentage"] = 2.1

	return c.JSON(stats)
}

func GetCarInventoryAPI(c *fiber.Ctx) error {
	search := c.Query("search", "")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset := (page - 1) * limit

	// For dashboard inventory, we likely only want cars that are NOT sold
	sold := false
	cars, total, err := database.GetAllCars(config.GetDB(), search, &sold, limit, offset)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch car inventory"})
	}

	return c.JSON(fiber.Map{
		"cars":     cars,
		"total":    total,
		"has_more": total > page*limit,
	})
}

func GetActiveDebtsAPI(c *fiber.Ctx) error {
	search := c.Query("search", "")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "8"))
	offset := (page - 1) * limit

	debts, total, err := database.GetDebts(config.GetDB(), search, limit, offset)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch debts"})
	}

	return c.JSON(fiber.Map{
		"debts":    debts,
		"total":    total,
		"has_more": total > page*limit,
	})
}
