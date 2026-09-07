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
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	// CORS Middleware
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

	// Health check endpoint
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Ewallet Backend is up and running on Vercel!",
		})
	})

	// Tambahkan rute dasar atau tangani API di sini
	r.Any("/api/*any", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "API endpoint reached successfully",
			"path":    c.Param("any"),
		})
	})

	app = r
}

func Handler(w http.ResponseWriter, r *http.Request) {
	once.Do(func() {
		initApp()
	})
	app.ServeHTTP(w, r)
}
