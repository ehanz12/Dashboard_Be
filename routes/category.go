package routes

import (
	"go-api-dashboard/handlers"
	"go-api-dashboard/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupCatRoute(router fiber.Router) {
	cat := router.Group("/category")

	cat.Get("/", middleware.JWTMiddleware, handlers.ListCategory)
	cat.Post("/", middleware.JWTMiddleware, handlers.CreateCategory)
	cat.Patch("/:id", middleware.JWTMiddleware, handlers.UpdateCategory)
	cat.Delete("/:id", middleware.JWTMiddleware, handlers.DeleteCategory)
}