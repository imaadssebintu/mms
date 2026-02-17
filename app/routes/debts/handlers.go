package debts

import (
	"github.com/gofiber/fiber/v2"
)

func ShowDebts(c *fiber.Ctx) error {
	user := c.Locals("user")
	return c.Render("debts", fiber.Map{
		"Title":       "Receivables Ledger",
		"CurrentPage": "debts",
		"User":        user,
	})
}
