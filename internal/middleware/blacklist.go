package middleware

import (
	"net/http"
	"strings"

	"github.com/ewallet-backend/internal/dto"
	"github.com/ewallet-backend/internal/repository"
	"github.com/gin-gonic/gin"
)

func Blacklist(blacklistRepo *repository.BlacklistRepository) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		bearerToken := ctx.GetHeader("Authorization")
		token := strings.Split(bearerToken, " ")[1]

		// cek blacklist
		isBlacklist, err := blacklistRepo.IsBlacklist(ctx.Request.Context(), token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, dto.Response{
				Message: "Error",
				Success: false,
				Error:   "Internal Server Error",
			})
			return
		}
		if isBlacklist {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.Response{
				Message: "Unauthorized",
				Success: false,
				Error:   "Token sudah tidak valid",
			})
			return
		}

		ctx.Next()
	}
}
