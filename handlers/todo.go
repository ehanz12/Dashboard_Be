package handlers

import (
	"go-api-dashboard/database"
	"go-api-dashboard/models"

	"github.com/gofiber/fiber/v2"
)

type TodoReq struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
	Status      string `json:"status" validate:"omitempty,oneof=pending in_progress completed"`
}

//func untuk membuat todo baru
func CreateTodoHandler(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	var req TodoReq

	//validasi payload
	if err := c.BodyParser(&req);err != nil {
		return c.Status(400).JSON(fiber.Map{"invalid" : "invalid payload"})
	}

	if req.Status == "" {
		req.Status = "pending"
	}
	//membuat todo baru
	todo := models.Todo{
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		UserID:      userID,
	}
	//simpan todo ke database
	if err := database.DB.Create(&todo).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}

	//ini untuk menampilkan user di response
	database.DB.Preload("User").First(&todo, todo.ID)
	return c.Status(201).JSON(fiber.Map{
		"Message" : "Berhasil Membuat Todo",
		"data" : todo,
	})
}

//func untuk mendapatkan semua todo dari user yang sedang login
func GetTodosHandler(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	var todos []models.Todo
	//mengambil todo dari database berdasarkan user_id
	if err := database.DB.Preload("User").Where("user_id = ?", userID).Find(&todos).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error" : err.Error()})
	}
	return c.Status(200).JSON(fiber.Map{
		"Message" : "Berhasil Mendapatkan Todos",
		"data" : todos,
	})
}

//func untuk mengupdate todo berdasarkan id
func UpdateTodoHandler(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	todoID, _ := c.ParamsInt("id")
	var req TodoReq

	//validasi payload
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"invalid" : "invalid payload"})
	}
	var todo models.Todo
	//mencari todo berdasarkan id dan user_id
	if err := database.DB.Where("id = ? AND user_id = ?", todoID, userID).First(&todo).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error" : "Todo not found"})
	}
	//mengupdate todo
	todo.Title = req.Title
	todo.Description = req.Description
	todo.Status = req.Status
	if err := database.DB.Save(&todo).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error" : err.Error()})
	}
	//ini untuk menampilkan user di response
	database.DB.Preload("User").First(&todo, todo.ID)
	return c.Status(200).JSON(fiber.Map{
		"Message" : "Berhasil Mengupdate Todo",
		"data" : todo,
	})
}

//func untuk menghapus todo berdasarkan id
func DeleteTodoHandler(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	todoID, _ := c.ParamsInt("id")
	var todo models.Todo
	//mencari todo berdasarkan id dan user_id
	if err := database.DB.Where("id = ? AND user_id = ?", todoID, userID).First(&todo).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error" : "Todo not found"})
	}
	//menghapus todo
	if err := database.DB.Delete(&todo).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error" : err.Error()})
	}
	return c.Status(200).JSON(fiber.Map{
		"Message" : "Berhasil Menghapus Todo",
	})
}

//func untuk set status todo
func SetTodoHandler(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	todoID, _ := c.ParamsInt("id")
	var req struct {
		Status string `json:"status" validate:"required,oneof=pending in_progress completed"`
	}
	//validasi payload
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"invalid" : "invalid payload"})
	}
	var todo models.Todo
	//mencari todo berdasarkan id dan user_id
	if err := database.DB.Where("id = ? AND user_id = ?", todoID, userID).First(&todo).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error" : "Todo not found"})
	}
	//mengupdate status todo
	todo.Status = req.Status
	if err := database.DB.Save(&todo).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error" : err.Error()})
	}
	//ini untuk menampilkan user di response
	database.DB.Preload("User").First(&todo, todo.ID)
	return c.Status(200).JSON(fiber.Map{
		"Message" : "Berhasil Mengupdate Status Todo",
		"data" : todo,
	})
}