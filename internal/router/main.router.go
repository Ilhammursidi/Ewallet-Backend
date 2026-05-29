package router

import (
	_ "github.com/ewallet-backend/docs"
	"github.com/ewallet-backend/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouter(router *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	router.Use(middleware.CORSMiddleware)
	router.Static("/img", "public/img")
	// swagger docs
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	RegisterAuthRouter(router, db, rdb)
	RegisterUserRouter(router, db, rdb)
	RegisterTransactionRouter(router, db, rdb)
}
