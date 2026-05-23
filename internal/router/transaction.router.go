package router

import (
	"github.com/ewallet-backend/internal/controller"
	"github.com/ewallet-backend/internal/middleware"
	"github.com/ewallet-backend/internal/repository"
	"github.com/ewallet-backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterTransactionRouter(router *gin.Engine, db *pgxpool.Pool) {

	// authRepo := repository.NewAuthRepository(db)
	transactionRepo := repository.NewTransactionRepo(db)

	transactionService := service.NewTransactionService(transactionRepo, db)
	transactionController := controller.NewTransactionController(transactionService)

	transactionRouter := router.Group("/transaction", middleware.VerifyToken)

	transactionRouter.GET("/receivers", transactionController.FindReceivers)
}
