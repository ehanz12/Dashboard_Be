package main

import (
	"fmt"
	"go-api-dashboard/config"
	"go-api-dashboard/database"
	"go-api-dashboard/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	//Load Env Variables
	config.LoadEnv()
	//Connect to Database
	database.ConnectDB()

	//Setup Routes
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173,https://dashboardreihan.vercel.app,https://www.reihan.biz.id", // url yang boleh akses
		AllowMethods: "GET,POST,PUT,DELETE,PATCH,OPTIONS",// method yang boleh dilakukan
		AllowHeaders: "Origin, Content-Type, Accept, Authorization", //content-type header wajib
		AllowCredentials: false, //jika pake jwt
	}))
	routes.SetupRoutes(app)

	fmt.Println("Server is Running in Port !", app.Listen(":3000"))
}