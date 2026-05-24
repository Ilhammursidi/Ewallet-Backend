package controller

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/ewallet-backend/internal/dto"
	"github.com/ewallet-backend/internal/service"
	"github.com/ewallet-backend/pkg"
	"github.com/gin-gonic/gin"
)

type TransactionController struct {
	transactionService *service.TransactionService
}

func NewTransactionController(transactionService *service.TransactionService) *TransactionController {
	return &TransactionController{
		transactionService: transactionService,
	}
}

// Find Receivers
//
//	@Summary		Find Receivers
//	@Description	Get list of transaction receivers with pagination and search filter
//	@Tags			transaction
//	@Accept			json
//	@Produce		json
//	@Security       ApiKeyAuth
//	@Param			page 	query		int					false	"Page number (default 1)" default(1)
//	@Param			limit	query		int					false	"Items per page (default 10)" default(10)
//	@Param			search	query		string				false	"Search by receiver name or phone number"
//	@Success		200		{object}	dto.Response				"Success get receivers list"
//	@Failure		400		{object}	dto.ErrorResponse			"Bad Request"
//	@Failure		401		{object}	dto.ErrorResponse			"Unauthorized"
//	@Router			/transaction/receivers [get]
func (tc *TransactionController) FindReceivers(ctx *gin.Context) {
	token, exists := ctx.Get("claims")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Message: "Claims not exists",
			Success: false,
			Error:   "Unauthorized",
		})
		return
	}

	claims, ok := token.(*pkg.Claims)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Message: "Error Uanuthorized",
			Success: false,
			Error:   "Unauthorized",
		})
		return
	}

	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Bad Request",
			Success: false,
			Error:   "Invalid limit parameter",
		})
		return
	}

	limit, err := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Bad Request",
			Success: false,
			Error:   "Invalid limit parameter",
		})
		return
	}

	search := strings.TrimSpace(ctx.DefaultQuery("search", ""))

	res, err := tc.transactionService.FindReceivers(ctx.Request.Context(), claims.Id, search, page, limit)
	if err != nil {
		log.Println("Error: ", err.Error())
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Bad Request",
			Success: false,
			Error:   "Invalid limit parameter",
		})
		return
	}
	ctx.JSON(http.StatusOK, dto.Response{
		Message: "OK",
		Success: true,
		Data:    res,
	})
}

// func (tc *TransactionController) CreateTopUp
