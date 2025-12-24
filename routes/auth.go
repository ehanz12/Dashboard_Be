package routes

import (
	"go-api-dashboard/handlers"
	"go-api-dashboard/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetRouteAuth(router fiber.Router) {
	auth := router.Group("/auth")

	auth.Get("/me",middleware.JWTMiddleware ,handlers.Me)
	auth.Post("/refresh", handlers.RefreshToken)
	auth.Post("/logout", handlers.LogoutHandler)
}