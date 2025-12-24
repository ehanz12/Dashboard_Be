package handlers

import (
	"go-api-dashboard/database"
	"go-api-dashboard/models"
	"go-api-dashboard/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
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
		return c.Status(400).JSON(fiber.Map{"error": "invalid payload"})
	}

	var user models.User
	database.DB.Where("email = ?", req.Email).First(&user)
	if user.ID == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "invalid email or password"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid email or password"})
	}

	// 🔐 generate tokens
	accessToken, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to generate access token"})
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to generate refresh token"})
	}

	// 🍪 simpan refresh token di cookie
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HTTPOnly: true,
		Secure:   true,     // WAJIB kalau production (https)
		SameSite: "None",   // karena FE & BE beda domain
		Expires:  time.Now().Add(7 * 24 * time.Hour),
	})

	return c.JSON(fiber.Map{
		"access_token": accessToken,
	})
}

func LogoutHandler(c *fiber.Ctx) error {
	//hapus refresh token
	c.Cookie(&fiber.Cookie{
		Name: "refresh_token",
		Value: "",
		HTTPOnly: true,
		Secure: true,
		SameSite: "none",
		Expires: time.Now().Add(-time.Hour),
	})

	return c.JSON(fiber.Map{"message" : "logout success"})
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


func RefreshToken(c *fiber.Ctx) error {
	RefreshToken := c.Cookies("refresh_token")
	if RefreshToken == "" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	token, err := utils.VerifyRefreshToken(RefreshToken)
	if err != nil || !token.Valid {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	claims := token.Claims.(jwt.MapClaims)
	userID := uint(claims["user_id"].(float64))

	var user models.User

	database.DB.Find(&user, userID)

	accessToken, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	} 

	return c.JSON(fiber.Map{"access_token" : accessToken})
}