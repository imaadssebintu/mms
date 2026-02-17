package clients

import (
	"github.com/gofiber/fiber/v2"
)

func ShowClients(c *fiber.Ctx) error {
	user := c.Locals("user")
	return c.Render("clients", fiber.Map{
		"Title":       "Clients Management",
		"CurrentPage": "clients",
		"User":        user,
	})
}
