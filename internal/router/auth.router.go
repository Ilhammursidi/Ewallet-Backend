package router

import (
	"github.com/ewallet-backend/internal/controller"
	"github.com/ewallet-backend/internal/middleware"
	"github.com/ewallet-backend/internal/repository"
	"github.com/ewallet-backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterAuthRouter(router *gin.Engine, db *pgxpool.Pool) {
	authRouter := router.Group("/auth")

	authRepo := repository.NewAuthRepository(db)
	blacklistRepo := repository.NewBlacklistRepository(db)

	authService := service.NewAuthService(authRepo, blacklistRepo)
	authController := controller.NewAuthController(authService)

	// authRouter.GET("", authController.GetAll)
	authRouter.POST("/register", authController.Register)
	authRouter.POST("/login", authController.Login)
	authRouter.PATCH("/enter-pin", middleware.VerifyToken, authController.CreatePin)
	authRouter.POST("/logout", middleware.VerifyToken, middleware.Blacklist(blacklistRepo), authController.Logout)
}
