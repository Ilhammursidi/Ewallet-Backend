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
	authRepo := repository.NewAuthRepository(db)

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo, rdb)
	userController := controller.NewUserController(userService)

	userRouter.GET("/dashboard", middleware.Blacklist(authRepo), middleware.VerifyToken, userController.GetDashboardInfo)
	userRouter.GET("/profile", middleware.Blacklist(authRepo), middleware.VerifyToken, userController.GetProfile)
	userRouter.PATCH("/profile", middleware.Blacklist(authRepo), middleware.VerifyToken, userController.EditUserProfile)
	userRouter.PATCH("/profile/change-password", middleware.Blacklist(authRepo), middleware.VerifyToken, userController.EditPassword)
	userRouter.PATCH("/profile/change-pin", middleware.Blacklist(authRepo), middleware.VerifyToken, userController.EditUserPin)
	userRouter.GET("/check-pin", middleware.Blacklist(authRepo), middleware.VerifyToken, userController.CheckPin)
	userRouter.GET("/transaction-report", middleware.Blacklist(authRepo), middleware.VerifyToken, userController.TransactionReportGraph)
	userRouter.GET("/transactions", middleware.Blacklist(authRepo), middleware.VerifyToken, userController.GetTransactionHistory)
}
