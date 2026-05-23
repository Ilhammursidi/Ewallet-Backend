package router

import (
	"github.com/ewallet-backend/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitRouter(router *gin.Engine, db *pgxpool.Pool) {
	router.Use(middleware.CORSMiddleware)

	RegisterAuthRouter(router, db)
	RegisterUserRouter(router, db)
	RegisterTransactionRouter(router, db)
}
