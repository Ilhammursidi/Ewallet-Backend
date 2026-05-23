package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ewallet-backend/internal/config"
	"github.com/ewallet-backend/internal/router"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// @title 						Backend koda 7
// @version						1.0
// @description					Backend created by Koda using Gin

// @license.name				MIT

// @host						localhost:8081
// @BasePath					/

// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description					Bearer token used for authorization
func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading env. \ncause: %s", err.Error())
	}
	app := gin.Default()

	db, err := config.ConnectPsql()
	if err != nil {
		log.Fatalf("DB connection error. \ncause: %s", err.Error())
	}
	defer db.Close()
	log.Println("DB Connected")

	router.InitRouter(app, db)
	app.Run(fmt.Sprintf("%s:%s", os.Getenv("APP_HOST"), os.Getenv("APP_PORT")))
}
