package response

import (
	"net/http"

	"github.com/ewallet-backend/internal/dto"
	"github.com/gin-gonic/gin"
)

// Status 500 - Internal Server Error
func JSONInternalServerError(ctx *gin.Context) {
	ctx.JSON(http.StatusInternalServerError, dto.Response{
		Message: "Error",
		Error:   "Internal Server Error",
	})
}

// Status 400 - Bad Request
func JSONBadRequest(ctx *gin.Context) {
	ctx.JSON(http.StatusBadRequest, dto.Response{
		Message: "Invalid Request Payload",
		Error:   "Bad Request",
	})
}

// Status 401 - Unauthorized
func JSONUnauthorized(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusUnauthorized, dto.Response{
		Message: message,
		Error:   "Unauthorized",
	})
}

// Status 409 - Conflict
func JSONDuplicate(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusConflict, dto.Response{
		Message: message,
		Error:   "Conflict",
	})
}

// Status 422 - Unprocessable Entity
func JSONUnprocessableEntity(ctx *gin.Context, errors map[string]string) {
	ctx.JSON(http.StatusUnprocessableEntity, dto.Response{
		Message: "Unprocessable Entity",
		Errors:  errors,
	})
}

// Status 200 - OK
func JSONSuccess(ctx *gin.Context, data any, message string) {
	ctx.JSON(http.StatusOK, dto.Response{
		Data:    data,
		Message: message,
		Success: true,
	})
}

// Status 201 - Created
func JSONCreated(ctx *gin.Context, data any, message string) {
	ctx.JSON(http.StatusCreated, dto.Response{
		Data:    data,
		Message: message,
		Success: true,
	})
}

// Status 204 - No Content
func JSONNoContent(ctx *gin.Context) {
	ctx.JSON(http.StatusNoContent, dto.Response{
		Message: "No Content",
		Success: true,
	})
}
