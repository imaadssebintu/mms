package cars

import (
	"mms/app/config"
	"mms/app/database"
	"mms/app/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func GetCarsAPI(c *fiber.Ctx) error {
	search := c.Query("search", "")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	offset := (page - 1) * limit

	var sellable *bool
	if s := c.Query("sellable"); s != "" {
		b := s == "true"
		sellable = &b
	}

	cars, total, err := database.GetAllCars(config.GetDB(), search, sellable, limit, offset)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch cars"})
	}

	return c.JSON(fiber.Map{
		"cars":     cars,
		"total":    total,
		"has_more": total > page*limit,
	})
}

func CreateCarAPI(c *fiber.Ctx) error {
	car := new(models.Car)
	if err := c.BodyParser(car); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := database.CreateCar(config.GetDB(), car); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create car"})
	}

	return c.Status(201).JSON(car)
}

func GetCarAPI(c *fiber.Ctx) error {
	id := c.Params("id")
	car, err := database.GetCarByID(config.GetDB(), id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Car not found"})
	}
	return c.JSON(car)
}
