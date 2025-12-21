package handlers

import (
	"go-api-dashboard/database"
	"go-api-dashboard/models"

	"github.com/gofiber/fiber/v2"
)

//payload wajib expense
type ExpReq struct {
	CategoryID *uint `json:"category_id"`
	Amount float64 `json:"amount"`
	Note string `json:"note"`
}

//func untuk create expense / pengeluaran
func CreateExpense(c *fiber.Ctx) error {
	//ambil id user saat login
	userID := c.Locals("user_id").(uint)
	//cek payload apakah sudah benar?
	var req ExpReq
	if err := c.BodyParser(&req);err != nil {
		return c.Status(400).JSON(fiber.Map{"error" : "invalid payload"})
	}

	//cari category apakah ada or tidak
	if req.CategoryID != nil {
		var cat models.Category
		database.DB.Where("id = ?", req.CategoryID).First(&cat)
		if cat.ID == 0 {
			return c.Status(404).JSON(fiber.Map{"error" : "not found category"})
		}
	}

	//created expense
	exp := models.Expense{
		UserID: userID,
		CategoryID: req.CategoryID,
		Amount: req.Amount,
		Note: req.Note,
	}
	//proses database create expense
	if err := database.DB.Create(&exp).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}

	//preload relasi untuk ditampilkan di JSON
	database.DB.Preload("Category").First(&exp, exp.ID)

	return c.Status(201).JSON(fiber.Map{
		"Message" : "Berhasil Create Expense",
		"data" : exp,
	})
}

func UpdateExpense(c *fiber.Ctx) error {
	//ambil user id yang login
	userID := c.Locals("user_id").(uint)
	expID , _ := c.ParamsInt("id")

	//cek payload apakah sudah benar?
	var req ExpReq
	if err := c.BodyParser(&req);err != nil {
		return c.Status(400).JSON(fiber.Map{"invalid" : "error payload"})
	}
	var exp models.Expense
	//cari expense berdasarkan user dan id expense
	database.DB.Where("id = ? AND user_id = ?", expID, userID).First(&exp)
	if exp.ID == 0 {
		return c.Status(404).JSON(fiber.Map{"error" : "Not found expense"})
	}
	//cari category dari request 
	if req.CategoryID != nil {
		var cat models.Category
		database.DB.Where("id = ? AND user_id = ?", req.CategoryID, userID).First(&cat)
		if cat.ID == 0 {
			return c.Status(404).JSON(fiber.Map{"error" : "not found category"})
		}
	}

	//proses lewat model
	exp.CategoryID = req.CategoryID
	exp.Amount = req.Amount
	exp.Note = req.Note

	//simpan ke DB
	database.DB.Save(&exp)

	//preload relasi untuk ditampilkan di JSON
	database.DB.Preload("Category").First(&exp, exp.ID)

	return c.Status(200).JSON(fiber.Map{"Message" : "berhasil update expense", "data" : exp})
}

//method delete
func DeleteExpense(c *fiber.Ctx) error {
	//ambil user id yang lagi login
	userID := c.Locals("user_id").(uint)
	expID, _ := c.ParamsInt("id")

	//cari exp berdasarkan id user dan expense
	var exp models.Expense
	database.DB.Where("id = ? AND user_id = ?", expID, userID).First(&exp)
	//jika tidak ada
	if exp.ID == 0 {
		return c.Status(404).JSON(fiber.Map{"error" : "not found expense"})
	}
	//hapus jika ada
	if err := database.DB.Delete(&exp).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(200).JSON(fiber.Map{
		"message" : "berhasil menghapus data",
		"data" : nil,
	})
}

//get all expense peruser
func ListExpenses(c *fiber.Ctx) error {
	//ambil user yang login (id)
	userID := c.Locals("user_id").(uint)
	//tampung di slice
	var exp []models.Expense
	database.DB.Where("user_id = ?", userID).Find(&exp)

	//preload dengan relasi
	database.DB.Preload("Category").Find(&exp)
	return c.Status(200).JSON(fiber.Map{
		"Message" : "found expense",
		"data" : exp,
	})
}


//list total berdasarkan bulan dan hari
func ExpenseSummary(c *fiber.Ctx) error {
	// ambil user id login
	userID := c.Locals("user_id").(uint)

	var dailyTotal float64
	var monthlyTotal float64

	// ================= DAILY =================
	database.DB.
		Model(&models.Expense{}).
		Where("user_id = ?", userID).
		Where("DATE(created_at) = CURDATE()").
		Select("COALESCE(SUM(amount), 0)").
		Scan(&dailyTotal)

	// ================= MONTHLY =================
	database.DB.
		Model(&models.Expense{}).
		Where("user_id = ?", userID).
		Where("MONTH(created_at) = MONTH(CURDATE())").
		Where("YEAR(created_at) = YEAR(CURDATE())").
		Select("COALESCE(SUM(amount), 0)").
		Scan(&monthlyTotal)

	// ================= RESPONSE =================
	return c.Status(200).JSON(fiber.Map{
		"daily_total":   dailyTotal,
		"monthly_total": monthlyTotal,
		"chart": []fiber.Map{
			{
				"name":  "Hari Ini",
				"value": dailyTotal,
			},
			{
				"name":  "Bulan Ini",
				"value": monthlyTotal,
			},
		},
	})
}
