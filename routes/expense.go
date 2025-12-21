package routes

import (
	"go-api-dashboard/handlers"
	"go-api-dashboard/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupExpenseRoute(router fiber.Router) {
	exp := router.Group("/expense")

	exp.Post("/", middleware.JWTMiddleware, handlers.CreateExpense)
	exp.Get("/", middleware.JWTMiddleware, handlers.ListExpenses)
	exp.Patch("/:id", middleware.JWTMiddleware, handlers.UpdateExpense)
	exp.Delete("/:id", middleware.JWTMiddleware, handlers.DeleteExpense)
	exp.Get("/summary", middleware.JWTMiddleware, handlers.ExpenseSummary)
}