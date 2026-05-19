package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitRouter(router *gin.Engine, db *pgxpool.Pool) {
	RegisterAuthRouter(router, db)
	RegisterUserRouter(router, db)
}
