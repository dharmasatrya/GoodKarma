package config

import (
	"fmt"
	"log"

	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// InitDB menginisialisasi koneksi ke database PostgreSQL menggunakan GORM.
// Koneksi dibuat menggunakan konfigurasi yang sudah ditentukan, baik untuk lingkungan lokal maupun Supabase.
func InitDatabase() *gorm.DB {
	// // Konfigurasi untuk koneksi ke database
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// Membuat DSN (Data Source Name) untuk koneksi ke PostgreSQL
	// Format: "host=<host> user=<user> password=<password> dbname=<dbname> port=<port> sslmode=disable TimeZone=Asia/Shanghai"
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		dbHost,
		dbUser,
		dbPassword,
		dbName,
		dbPort,
	)

	// Membuka koneksi ke database menggunakan GORM
	var err error
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, // Disables prepared statement caching
	}), &gorm.Config{})

	if err != nil {
		// Jika koneksi gagal, log error dan keluar dari aplikasi
		log.Fatalf("Failed to connect to the database: %v", err)
		return nil
	}

	// Mengembalikan objek koneksi database GORM yang berhasil dibuka
	return db
}
