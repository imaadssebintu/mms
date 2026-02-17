package clients

import (
	"mms/app/config"
	"mms/app/database"
	"mms/app/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func GetClientsAPI(c *fiber.Ctx) error {
	search := c.Query("search", "")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "8"))
	offset := (page - 1) * limit

	clients, total, err := database.GetAllClients(config.GetDB(), search, limit, offset)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch clients"})
	}

	return c.JSON(fiber.Map{
		"clients":  clients,
		"total":    total,
		"has_more": total > page*limit,
	})
}

func GetClientAPI(c *fiber.Ctx) error {
	id := c.Params("id")
	client, err := database.GetClientByID(config.GetDB(), id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Client not found"})
	}
	return c.JSON(client)
}

func CreateClientAPI(c *fiber.Ctx) error {
	client := new(models.Client)
	if err := c.BodyParser(client); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := database.CreateClient(config.GetDB(), client); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create client"})
	}

	return c.Status(201).JSON(client)
}
