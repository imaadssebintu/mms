package clients

import (
	"mms/app/routes/auth"

	"github.com/gofiber/fiber/v2"
)

func SetupClientRoutes(app *fiber.App) {
	app.Get("/clients", auth.AuthMiddleware, ShowClients)

	clients := app.Group("/api/clients")
	clients.Use(auth.AuthMiddleware)

	clients.Get("/", GetClientsAPI)
	clients.Get("/:id", GetClientAPI)
	clients.Post("/", CreateClientAPI)
}
