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

func RegisterUserRouter(router *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	userRouter := router.Group("/users")

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo, rdb)
	userController := controller.NewUserController(userService)

	userRouter.GET("/dashboard", middleware.VerifyToken, userController.GetDashboardInfo)
	userRouter.GET("/profile", middleware.VerifyToken, userController.GetProfile)
	userRouter.PATCH("/profile", middleware.VerifyToken, userController.EditUserProfile)
	{
		userRouter.PATCH("/change-pin", middleware.VerifyToken, userController.EditUserPin)
		userRouter.PATCH("/change-password", middleware.VerifyToken, userController.EditPassword)
	}
	userRouter.GET("/check-pin", middleware.VerifyToken, userController.CheckPin)
	userRouter.GET("/transaction-report", middleware.VerifyToken, userController.TransactionReportGraph)
	userRouter.GET("/transactions", middleware.VerifyToken, userController.GetTransactionHistory)
}
