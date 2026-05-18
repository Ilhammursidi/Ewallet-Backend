package router

import (
	"github.com/ewallet-backend/internal/controller"
	"github.com/ewallet-backend/internal/repository"
	"github.com/ewallet-backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func ListUserRouter(router *gin.Engine, db *pgxpool.Pool) {
	authRouter := router.Group("/auth")

	authRepo := repository.NewAuthRepository(db)
	authService := service.NewAuthService(authRepo)
	authController := controller.NewUserController(authService)

	authRouter.GET("", authController.GetAll)
	authRouter.POST("/register", authController.Register)
}
