package router

import (
	_ "github.com/ewallet-backend/docs"
	"github.com/ewallet-backend/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouter(router *gin.Engine, db *pgxpool.Pool) {
	router.Use(middleware.CORSMiddleware)
	router.Static("/img", "public/img")
	// swagger docs
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	RegisterAuthRouter(router, db)
	RegisterUserRouter(router, db)
	RegisterTransactionRouter(router, db)
}
