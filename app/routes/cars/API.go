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
	// Define a custom request struct to handle nested seller info
	type CreateCarRequest struct {
		models.Car
		SellerID      string `json:"seller_id"`
		SellerName    string `json:"seller_name"`
		SellerEmail   string `json:"seller_email"`
		SellerPhone   string `json:"seller_phone"`
		SellerAddress string `json:"seller_address"`
		IsNewSeller   bool   `json:"is_new_seller"`
	}

	req := new(CreateCarRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// copy base fields to car model
	car := &req.Car

	// 1. Check if plate already exists to prevent 500 unique constraint error
	if car.NumberPlate != "" {
		db := config.GetDB()
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM cars WHERE number_plate = $1", car.NumberPlate).Scan(&count)
		if err == nil && count > 0 {
			return c.Status(400).JSON(fiber.Map{"error": "Asset with this Number Plate already exists in the system"})
		}
	}

	// 2. Handle Seller Logic
	if req.IsNewSeller && req.SellerName != "" {
		// PREVENT DUPLICATES: Check if client with same name and phone exists
		existing, err := database.GetClientByNameAndPhone(config.GetDB(), req.SellerName, req.SellerPhone)
		if err == nil && existing != nil {
			// Reuse existing client
			car.SellerID = &existing.ID
		} else {
			// Create new client
			seller := &models.Client{
				Name:    req.SellerName,
				Email:   req.SellerEmail,
				Phone:   req.SellerPhone,
				Address: req.SellerAddress,
				Notes:   "Auto-created as Seller for Car: " + car.Make + " " + car.Model,
			}

			if err := database.CreateClient(config.GetDB(), seller); err != nil {
				return c.Status(500).JSON(fiber.Map{"error": "Failed to create seller record"})
			}
			car.SellerID = &seller.ID
		}
	} else if req.SellerID != "" {
		car.SellerID = &req.SellerID
	}

	// 3. Create the Car
	if err := database.CreateCar(config.GetDB(), car); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to finalize asset registration: " + err.Error()})
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
