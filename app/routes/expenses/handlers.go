package expenses

import (
	"github.com/gofiber/fiber/v2"
)

func ShowExpenses(c *fiber.Ctx) error {
	user := c.Locals("user")
	return c.Render("expenses", fiber.Map{
		"Title":       "Expense Ledger",
		"CurrentPage": "expenses",
		"User":        user,
	})
}
