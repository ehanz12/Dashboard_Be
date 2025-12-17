package handlers

import (
	"go-api-dashboard/database"
	"go-api-dashboard/models"
	"go-api-dashboard/utils"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// Payload JSON apa ajah yang diizinkan ketika register dan login
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Func handler untuk register
func RegisterHandler(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"Invalid": "error payload"})
	}

	var exits models.User
	database.DB.Where("email = ?", req.Email).First(&exits)
	if exits.ID != 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Email already exists"})
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to hash password"})
	}
	user := models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hash),
	}
	database.DB.Create(&user)
	return c.Status(201).JSON(fiber.Map{"message": "User created successfully", "data": user})
}

// Func handler untuk login
func LoginHandler(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"Invalid" : "error payload"})
	}

	var user models.User
	database.DB.Where("email = ?", req.Email).First(&user)
	if user.ID == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid email or password"})
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid email or password"})
	}

	token, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to generate token"})
	}
	return c.Status(200).JSON(fiber.Map{"message": "Login successful", "token": token})
}


func Me(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)//mengambil user_id dari context

	var user models.User// ambil data user dari database
	database.DB.First(&user, userID)// mencari user berdasarkan user_id
	if user.ID == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	return c.Status(200).JSON(fiber.Map{"data": user})
}