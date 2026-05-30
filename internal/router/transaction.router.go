package router

import (
	"github.com/ewallet-backend/internal/controller"
	"github.com/ewallet-backend/internal/middleware"
	"github.com/ewallet-backend/internal/repository"
	"github.com/ewallet-backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func RegisterTransactionRouter(router *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {

	transactionRepo := repository.NewTransactionRepo(db)

	transactionService := service.NewTransactionService(transactionRepo, db, rdb)
	transactionController := controller.NewTransactionController(transactionService)

	transactionRouter := router.Group("/transaction", middleware.Blacklist(rdb), middleware.VerifyToken)

	transactionRouter.GET("/receivers", transactionController.FindReceivers)
	transactionRouter.POST("/topup", transactionController.TopUp)
	transactionRouter.POST("/transfer", transactionController.Transfer)
}
