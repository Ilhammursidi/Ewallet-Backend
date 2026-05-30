package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/ewallet-backend/internal/dto"
	// "github.com/ewallet-backend/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func Blacklist(rdb *redis.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		bearerToken := ctx.GetHeader("Authorization")
		// log.Println("ini", bearerToken)
		token := strings.Split(bearerToken, " ")[1]
		if len(token) < 2 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{
				Message: "Unauthorized",
				Success: false,
				Error:   "invalid token format",
			})
			return
		}

		rkey := "ilhammursidi:blacklist:" + token

		isBlacklist, err := rdb.Get(ctx, rkey).Result()

		if err != nil {
			if err == redis.Nil {
				log.Println("cache miss == aman")
				ctx.Next()
				return
			}

			log.Println("Error Redis:", err.Error())
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{
				Message: "Error",
				Success: false,
				Error:   "Internal Server Error (Gagal validasi session)",
			})
			return
		}

		log.Println("cache hit == sudah di block", isBlacklist)

		ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{
			Message: "Unauthorized",
			Success: false,
			Error:   "token sudah di blacklist",
		})
	}
}
