package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	// Ganti "github.com/ewallet-backend" di bawah ini
	// dengan nama module yang ada di file go.mod project kamu.
	"github.com/ewallet-backend/internal/router"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

var (
	app       http.Handler
	once      sync.Once
	dbPool    *pgxpool.Pool
	rdbClient *redis.Client
)

func initApp() {
	// 1. Set Gin ke mode production
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	// 2. Middleware CORS untuk Vercel
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Koda-X")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	r.Use(gin.Recovery())

	// 3. Koneksi Database PostgreSQL (Mendukung Connection String atau variabel terpisah)
	ctx := context.Background()
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		// Jika di Vercel kamu menyimpannya terpisah seperti di file .env lokalmu
		dbURL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASS"),
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"), // Pastikan port di Vercel terset 6543
			os.Getenv("DB_NAME"),
			os.Getenv("DB_SSL_MODE"),
		)
	}

	var err error
	dbPool, err = pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Printf("Gagal konek ke database: %v", err)
	}

	// 4. Inisialisasi Redis
	rdbClient = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("RDB_ADDR"),
		Password: os.Getenv("RDB_PASS"),
	})

	// 5. Health Check endpoint
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Sakuin Backend is online on Vercel with Database connected!",
		})
	})

	// 6. MENERUSKAN ROUTER ASLI DENGAN DATABASE DAN REDIS
	// Ini otomatis mengaktifkan semua rute (auth, user, transaction) milikmu!
	router.InitRouter(r, dbPool, rdbClient)

	app = r
}

// Handler utama Vercel
func Handler(w http.ResponseWriter, r *http.Request) {
	once.Do(func() {
		initApp()
	})
	app.ServeHTTP(w, r)
}
