package main

import (
	"go-api-dashboard/config"
	"go-api-dashboard/database"
	"go-api-dashboard/routes"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	//Load Env Variables
	config.LoadEnv()
	//Connect to Database
	database.ConnectDB()

	//Setup Routes
	app := fiber.New()
	routes.SetupRoutes(app)

	log.Fatal("Server is Running in Port !", app.Listen(":3000"))
}