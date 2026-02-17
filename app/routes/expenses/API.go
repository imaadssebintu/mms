package expenses

import (
	"mms/app/config"
	"mms/app/database"
	"mms/app/models"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func CreateExpenseAPI(c *fiber.Ctx) error {
	expense := new(models.Expense)
	if err := c.BodyParser(expense); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if expense.Date.IsZero() {
		expense.Date = time.Now()
	}

	if err := database.CreateExpense(config.GetDB(), expense); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create expense"})
	}

	return c.Status(201).JSON(expense)
}

func GetCarExpensesAPI(c *fiber.Ctx) error {
	carID := c.Params("car_id")
	expenses, err := database.GetExpensesByCarID(config.GetDB(), carID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch expenses"})
	}
	return c.JSON(expenses)
}

func GetExpensesAPI(c *fiber.Ctx) error {
	search := c.Query("search", "")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	offset := (page - 1) * limit

	expenses, total, err := database.GetAllExpenses(config.GetDB(), search, limit, offset)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch expenses"})
	}

	return c.JSON(fiber.Map{
		"expenses": expenses,
		"total":    total,
		"has_more": total > page*limit,
	})
}
