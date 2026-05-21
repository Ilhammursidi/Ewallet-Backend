package router

import (
	"github.com/ewallet-backend/internal/controller"
	"github.com/ewallet-backend/internal/middleware"
	"github.com/ewallet-backend/internal/repository"
	"github.com/ewallet-backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterUserRouter(router *gin.Engine, db *pgxpool.Pool) {
	userRouter := router.Group("/users")

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userController := controller.NewUserController(userService)

	userRouter.GET("/profile", middleware.VerifyToken, userController.GetProfile)
	userRouter.GET("/dashboard", middleware.VerifyToken, userController.GetDashboardInfo)
	userRouter.PATCH("/profile", middleware.VerifyToken, userController.EditUserProfile)
	userRouter.PATCH("/profile/change-pin", middleware.VerifyToken, userController.EditUserPin)
	userRouter.PATCH("/profile/change-password", middleware.VerifyToken, userController.EditPassword)
}
