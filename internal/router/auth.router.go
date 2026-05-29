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

func RegisterAuthRouter(router *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	authRouter := router.Group("/auth")

	authRepo := repository.NewAuthRepository(db)
	blacklistRepo := repository.NewBlacklistRepository(db)
	cacheRepo := repository.NewCacheRepository(rdb)

	authService := service.NewAuthService(authRepo, blacklistRepo, rdb, cacheRepo)
	authController := controller.NewAuthController(authService)

	authRouter.POST("/register", authController.Register)
	authRouter.POST("/login", authController.Login)
	authRouter.PATCH("/enter-pin", middleware.Blacklist(authRepo), middleware.VerifyToken, authController.CreatePin)
	authRouter.POST("/logout", middleware.Blacklist(authRepo), middleware.VerifyToken, authController.Logout)
	authRouter.POST("/forgot-password", authController.RequestForgotPassword)
	authRouter.POST("/verify-reset-token", authController.VerifyResetToken)
	authRouter.POST("/reset-password", authController.ResetPassword)
}
