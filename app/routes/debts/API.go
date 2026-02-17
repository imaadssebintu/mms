package debts

import (
	"mms/app/config"
	"mms/app/database"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func GetDebtsAPI(c *fiber.Ctx) error {
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
