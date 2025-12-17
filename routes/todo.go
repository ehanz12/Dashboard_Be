package routes

import (
	"go-api-dashboard/handlers"
	"go-api-dashboard/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupTodoRoutes(router fiber.Router) {
	todo := router.Group("/todos")

	todo.Post("/", middleware.JWTMiddleware, handlers.CreateTodoHandler)
	todo.Get("/", middleware.JWTMiddleware, handlers.GetTodosHandler)
	todo.Patch("/:id", middleware.JWTMiddleware, handlers.UpdateTodoHandler)
	todo.Delete("/:id", middleware.JWTMiddleware, handlers.DeleteTodoHandler)
	todo.Patch("/set_status/:id", middleware.JWTMiddleware, handlers.SetTodoHandler)
}