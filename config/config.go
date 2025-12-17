package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBName     string
	DBPassword string
	DBUser     string
	DBhost     string
	DBPort     string
}

var ConfDB *Config

func LoadEnv() {
	//cek env apakah ada di project?
	if err := godotenv.Load(); err != nil {
		log.Println("❓ file Env Not Found In Project !")
	}

	//inisialisasi ConfDB dengan data dari env
	ConfDB = &Config{
		DBName:     os.Getenv("DB_NAME"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBUser:     os.Getenv("DB_USER"),
		DBhost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
	}
}
