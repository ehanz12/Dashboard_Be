package main

import (
	"fmt"
	"go-api-dashboard/config"
	"go-api-dashboard/database"
	"go-api-dashboard/routes"

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

	fmt.Println("Server is Running in Port !", app.Listen(":3000"))
}