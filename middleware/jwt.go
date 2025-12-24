package middleware

import (
	"go-api-dashboard/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTMiddleware(c *fiber.Ctx) error {
	// JWT middleware implementation
	auth := c.Get("Authorization")
	if auth == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing or invalid JWT"})
	}
	if !strings.HasPrefix(auth, "Bearer ") {
		return c.Status(401).JSON(fiber.Map{"error" : "invalid token format"})
	}
	tokenStr := strings.TrimPrefix(auth, "Bearer ")//menghilangkan "Bearer " dari token
	token, err := utils.VerifyAccessToken(tokenStr)//memvalidasi token
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired JWT"})
	}

	claims, ok := token.Claims.(jwt.MapClaims)//mengambil claims dari token
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid JWT claims"})
	}
	userID, ok := claims["user_id"].(float64)//mengambil user_id dari claims
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid JWT claims"})
	}
	c.Locals("user_id", uint(userID))//menyimpan user_id di context

	return c.Next()
}