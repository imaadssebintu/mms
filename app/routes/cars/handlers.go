package cars

import (
	"mms/app/config"
	"mms/app/database"

	"github.com/gofiber/fiber/v2"
)

func ShowCars(c *fiber.Ctx) error {
	user := c.Locals("user")
	return c.Render("cars", fiber.Map{
		"Title":       "Fleet Management",
		"CurrentPage": "cars",
		"User":        user,
	})
}

func ShowNewCar(c *fiber.Ctx) error {
	user := c.Locals("user")
	return c.Render("cars_form", fiber.Map{
		"Title":       "Add New Vehicle",
		"CurrentPage": "cars",
		"User":        user,
		"IsEdit":      false,
	})
}

func ShowEditCar(c *fiber.Ctx) error {
	id := c.Params("id")
	car, err := database.GetCarByID(config.GetDB(), id)
	if err != nil {
		return c.Redirect("/cars")
	}

	user := c.Locals("user")
	return c.Render("cars_form", fiber.Map{
		"Title":       "Edit Vehicle",
		"CurrentPage": "cars",
		"User":        user,
		"Car":         car,
		"IsEdit":      true,
	})
}
