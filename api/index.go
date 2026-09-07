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

	// Health Check aman
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Ewallet Backend is online!",
		})
	})

	// Contoh endpoint login & signup manual agar tidak 404/500
	r.POST("/api/auth/login", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Login endpoint ready"})
	})

	r.POST("/api/auth/signup", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Signup endpoint ready"})
	})

	app = r
}

func Handler(w http.ResponseWriter, r *http.Request) {
	once.Do(func() {
		initApp()
	})
	app.ServeHTTP(w, r)
}
