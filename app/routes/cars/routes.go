package cars

import (
	"mms/app/routes/auth"

	"github.com/gofiber/fiber/v2"
)

func SetupCarRoutes(app *fiber.App) {
	// UI Routes
	app.Get("/cars", auth.AuthMiddleware, ShowCars)
	app.Get("/cars/new", auth.AuthMiddleware, ShowNewCar)
	app.Get("/cars/edit/:id", auth.AuthMiddleware, ShowEditCar)

	// API Routes
	cars := app.Group("/api/cars")
	cars.Use(auth.AuthMiddleware)

	cars.Get("/", GetCarsAPI)
	cars.Get("/:id", GetCarAPI)
	cars.Post("/", CreateCarAPI)
}
