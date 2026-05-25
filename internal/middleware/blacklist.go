package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/ewallet-backend/internal/dto"
	"github.com/ewallet-backend/internal/repository"
	"github.com/gin-gonic/gin"
)

func Blacklist(authRepo *repository.AuthRepository) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		bearerToken := ctx.GetHeader("Authorization")
		// log.Println("ini", bearerToken)
		token := strings.Split(bearerToken, " ")[1]

		// cek blacklist
		isBlacklist, err := authRepo.IsBlacklist(ctx.Request.Context(), token)
		log.Println("isblack", isBlacklist)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{
				Message: "Error",
				Success: false,
				Error:   "Internal(middleware1) Server Error",
			})
			return
		}
		if isBlacklist {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{
				Message: "Unauthorized",
				Success: false,
				Error:   "Token sudah tidak valid",
			})
			return
		}

		ctx.Next()
	}
}
