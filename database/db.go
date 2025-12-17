package database

import (
	"fmt"
	"go-api-dashboard/config"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

//variable DB untuk menyimpan koneksi database
var DB *gorm.DB

func ConnectDB() {
	//ambil config database dari config package
	cfg := config.ConfDB
	//lakukan koneksi ke database menggunakan cfg
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBhost,
		cfg.DBPort,
		cfg.DBName,
	)

	//inisialisasi variabel DB dengan koneksi database yang berhasil
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Gagal Terhubung ke Database :", err)
	}
	DB = db
	fmt.Println("✅ Berhasil Terhubung ke Database")
}