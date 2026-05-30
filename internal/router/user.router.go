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

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo, rdb)
	userController := controller.NewUserController(userService)
	userRouter := router.Group("/users", middleware.Blacklist(rdb), middleware.VerifyToken)

	userRouter.GET("/dashboard", userController.GetDashboardInfo)
	userRouter.GET("/profile", userController.GetProfile)
	userRouter.PATCH("/profile", userController.EditUserProfile)
	userRouter.PATCH("/profile/change-password", userController.EditPassword)
	userRouter.PATCH("/profile/change-pin", userController.EditUserPin)
	userRouter.GET("/transaction-report", userController.TransactionReportGraph)
	userRouter.GET("/transactions", userController.GetTransactionHistory)
}
