package routes

import (
	"go-api-dashboard/handlers"

	"github.com/gofiber/fiber/v2"
)



func SetupRoutes(app *fiber.App) {
	api := app.Group("/api/v1")
	api.Post("/register", handlers.RegisterHandler)
	api.Post("/login", handlers.LoginHandler)

	SetRouteAuth(api)
	SetupTodoRoutes(api)
	SetupCatRoute(api)
	SetupExpenseRoute(api)
}