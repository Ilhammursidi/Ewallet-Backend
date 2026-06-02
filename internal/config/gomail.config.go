package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config menampung semua konfigurasi aplikasi
type Config struct {
	AppPort       string
	FrontendURL   string
	SMTPHost      string
	SMTPPort      string
	SMTPUser      string
	SMTPPassword  string
	SMTPFromEmail string
}

// LoadConfig memuat file .env dan mengembalikan objek Config
func LoadConfig() *Config {
	// Memuat file .env sekali di awal aplikasi berjalan
	err := godotenv.Load()
	if err != nil {
		log.Println("Peringatan: File .env tidak ditemukan, menggunakan env sistem")
	}

	return &Config{
		AppPort:       getEnv("PORT", "8081"),
		FrontendURL:   getEnv("FRONTEND_URL", "http://localhost:5173"),
		SMTPHost:      os.Getenv("SMTP_HOST"),
		SMTPPort:      os.Getenv("SMTP_PORT"),
		SMTPUser:      os.Getenv("SMTP_USER"),
		SMTPPassword:  os.Getenv("SMTP_PASSWORD"),
		SMTPFromEmail: os.Getenv("SMTP_FROM_EMAIL"),
	}
}

// Helper untuk memberikan nilai default jika env kosong
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
