package middleware

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/ewallet-backend/internal/dto"
	"github.com/ewallet-backend/internal/repository"
	"github.com/ewallet-backend/pkg"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func VerifyToken(authRepo *repository.AuthRepository) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		bearerToken := ctx.GetHeader("Authorization")
		if bearerToken == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.Response{
				Message: "Unauthorized access, please login",
				Success: false,
				Error:   "Unauthorized(bearer) access, please login",
			})
			return
		}
		splittedBearer := strings.Split(bearerToken, " ")
		if len(splittedBearer) != 2 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.Response{
				Message: "Unauthorized access, please login",
				Success: false,
				Error:   "invalid token",
			})
			return
		}
		token := splittedBearer[1]

		var claims pkg.Claims
		if err := claims.VerifyJWT(token); err != nil {
			log.Println("Error: ", err.Error())
			if errors.Is(err, jwt.ErrTokenInvalidIssuer) || errors.Is(err, jwt.ErrTokenExpired) {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.Response{
					Message: "Unauthorized Access, Please Login",
					Success: false,
					Error:   err.Error(),
				})
				return
			}
			isBlacklist, err := authRepo.IsBlackList(ctx.Request.Context(), token)
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
			ctx.Set("claims", claims)
			ctx.Next()
		}
	}
}
