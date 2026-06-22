package middleware

import (
	"net/http"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware(ctx *gin.Context) {
	allowedOrigin := []string{"http://localhost:5173", "http://localhost:8081", "http://localhost:81"}
	currentOrigin := ctx.GetHeader("Origin")
	if slices.Contains(allowedOrigin, currentOrigin) {
		ctx.Header("Access-Control-Allow-Origin", currentOrigin)
	}

	allowedHeaders := []string{"Content-Type", "X-koda-X", "Authorization", "Origin", "Accept"}
	ctx.Header("Access-Control-Allow-Headers", strings.Join(allowedHeaders, ", "))

	allowedMethods := []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch, http.MethodOptions}
	ctx.Header("Access-Control-Allow-Methods", strings.Join(allowedMethods, ", "))

	ctx.Header("Access-Control-Allow-Credentials", "true")

	if ctx.Request.Method == http.MethodOptions {
		ctx.AbortWithStatus(http.StatusNoContent)
		return
	}

	ctx.Next()
}
