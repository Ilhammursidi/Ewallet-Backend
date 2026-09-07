package handler

import (
	"context"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/ewallet-backend/internal/router" // Sesuaikan dengan nama module di go.mod kamu
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

var (
	app  http.Handler
	once sync.Once
)

func initApp() {
	// 1. Set Gin ke mode Release untuk production di Vercel
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	// 2. Ambil environment variables untuk koneksi Database & Redis
	dbURL := os.Getenv("DATABASE_URL") // Atau susun dari DB_HOST, DB_USER, DB_PASS, dll jika terpisah
	if dbURL == "" {
		// Fallback jika pakai variabel terpisah
		dbHost := os.Getenv("DB_HOST")
		dbPort := os.Getenv("DB_PORT")
		dbUser := os.Getenv("DB_USER")
		dbPass := os.Getenv("DB_PASS")
		dbName := os.Getenv("DB_NAME")
		dbSSL := os.Getenv("DB_SSLMODE")
		dbURL = "postgres://" + dbUser + ":" + dbPass + "@" + dbHost + ":" + dbPort + "/" + dbName + "?sslmode=" + dbSSL
	}

	ctx := context.Background()

	// 3. Inisialisasi koneksi PostgreSQL Pool (Supabase)
	dbPool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("Gagal terhubung ke database: %v", err)
	}

	// 4. Inisialisasi koneksi Redis (opsional, sesuaikan jika ada env Redis)
	redisURL := os.Getenv("REDIS_URL")
	var rdb *redis.Client
	if redisURL != "" {
		opt, err := redis.ParseURL(redisURL)
		if err == nil {
			rdb = redis.NewClient(opt)
		}
	}
	if rdb == nil {
		// Fallback instance kosong jika Redis belum diset
		rdb = redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	}

	// 5. Panggil fungsi InitRouter yang sudah ada sejak dulu!
	router.InitRouter(r, dbPool, rdb)

	app = r
}

// Handler utama yang dipanggil oleh Vercel
func Handler(w http.ResponseWriter, r *http.Request) {
	once.Do(func() {
		initApp()
	})
	app.ServeHTTP(w, r)
}
