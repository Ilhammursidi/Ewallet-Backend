package handler

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	app  http.Handler
	once sync.Once
)

func initApp() {
	// Set Gin ke mode production
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	// Tambahkan middleware dasar atau langsung definisikan endpoint di sini
	// Tambahkan middleware CORS manual di index.go untuk Vercel
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*") // Atau ganti dengan domain frontend Vercel kamu
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Koda-X")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	r.Use(gin.Recovery())

	// Contoh rute dasar agar backend aktif dan merespons
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Sakuin Backend is online on Vercel!",
		})
	})

	// Tambahkan rute e-wallet kamu di sini jika diperlukan tanpa import folder internal

	app = r
}

// Handler utama Vercel
func Handler(w http.ResponseWriter, r *http.Request) {
	once.Do(func() {
		initApp()
	})
	app.ServeHTTP(w, r)
}
