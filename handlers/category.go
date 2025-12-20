package handlers

import (
	"go-api-dashboard/database"
	"go-api-dashboard/models"

	"github.com/gofiber/fiber/v2"
)

//membuat struct request form category
type CatReq struct {
	Name string `json:"name"`
}

func CreateCategory(c *fiber.Ctx) error {
	//mendapatkan user_id dari login
	UserID := c.Locals("user_id").(uint)
	var req CatReq

	//pengecekan jika payload tidak sesuai
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error" : "invalid request"})
	}
	//membuat add dari model
	cat := models.Category{
		Name: req.Name,
		UserID: UserID,
	}

	//tambahkan ke database dan membuat error handlernya
	if err := database.DB.Create(&cat).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}
	
	
	return c.Status(201).JSON(fiber.Map{
		"Message" : "Category Berhasil Dibuat",
		"data" : cat,
	})
	
}

func UpdateCategory(c *fiber.Ctx) error {
	//ambil id user dari login
	userID := c.Locals("user_id").(uint)
	catID, _ := c.ParamsInt("id")
	//cek payload
	var req CatReq
	if err := c.BodyParser(&req);err != nil {
		return c.Status(400).JSON(fiber.Map{"error" : "invalid payload"})
	}

	//mencari category berdasarkan id
	var cat models.Category
	database.DB.Where("id = ? AND user_id = ?", catID, userID).First(&cat)
	if cat.ID == 0 {
		return c.Status(404).JSON(fiber.Map{"error" : "not found category"})
	}

	//proses update category
	cat.Name = req.Name

	//simpan ke database
	if err := database.DB.Save(&cat).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(200).JSON(fiber.Map{
		"Message" : "Update category berhasil",
		"data" : cat,
	})
}

func DeleteCategory(c *fiber.Ctx) error {
	UserID := c.Locals("user_id").(uint)
	CatID , _ := c.ParamsInt("id")

	//mencari category
	var cat models.Category
	database.DB.Where("id = ? AND user_id = ?", CatID, UserID).First(&cat)

	if cat.ID == 0 {
		return c.Status(404).JSON(fiber.Map{"error" : "not found category !"})
	}

	//delete category
	database.DB.Delete(&cat)

	return c.Status(200).JSON(fiber.Map{
		"Message" :"Berhasil menghapus Category",
		"data" : nil,
	})
}

func ListCategory(c *fiber.Ctx) error {
	UserID := c.Locals("user_id").(uint)

	//cari category berdasarkan id user login
	var cat []models.Category
	//cari category di database
	database.DB.Where("user_id = ?", UserID).Find(&cat)

	return c.Status(200).JSON(fiber.Map{"Message" : "category found", "data": cat})
}