package middleware

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/ewallet-backend/internal/dto"
	"github.com/ewallet-backend/pkg"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func VerifyToken(ctx *gin.Context) {
	log.Println("VerifyToken dipanggil") // tambah ini
	bearerToken := ctx.GetHeader("Authorization")
	log.Println("bearerToken:", bearerToken) // tambah ini
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

	claims := &pkg.Claims{}
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
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, dto.Response{
			Message: "Error",
			Success: false,
			Error:   "Internal Server Error",
		})
		return
	}

	ctx.Set("claims", claims)
	ctx.Next()
}
