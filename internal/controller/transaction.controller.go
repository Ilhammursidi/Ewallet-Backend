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

// TopUp godoc
// @Summary      Top up wallet balance
// @Description  Add balance to user's wallet using a payment method
// @Tags         transaction
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        request  body      dto.TopUpHTTPRequest  true  "Top up request payload"
// @Success      201      {object}  dto.TopUpResponse     "Top up successful"
// @Failure      400      {object}  dto.ErrorResponse     "Invalid request body or amount"
// @Failure      401      {object}  dto.ErrorResponse     "Unauthorized"
// @Failure      500      {object}  dto.ErrorResponse     "Internal server error"
// @Router       /transaction/topup [post]
func (tc *TransactionController) TopUp(ctx *gin.Context) {
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
			Message: "Error Unauthorized",
			Success: false,
			Error:   "Unauthorized",
		})
		return
	}

	var req dto.TopUpHTTPRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "invalid request body",
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	result, err := tc.transactionService.TopUp(ctx.Request.Context(), dto.TopUpServiceRequest{
		UserID:          claims.Id,
		PaymentMethodID: req.PaymentMethodID,
		OrderAmount:     req.OrderAmount,
		TaxAmount:       req.TaxAmount,
		DeliveryFee:     req.DeliveryFee,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: err.Error(),
			Success: false,
			Error:   "Internal Server Error",
		})
		return
	}

	ctx.JSON(http.StatusCreated, result)
}

// Transfer godoc
// @Summary      Transfer balance to another wallet
// @Description  Transfer balance from sender wallet to receiver wallet
// @Tags         transaction
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        request  body      dto.TransferHTTPRequest  true  "Transfer request payload"
// @Success      201      {object}  dto.TransferResponse     "Transfer successful"
// @Failure      400      {object}  dto.ErrorResponse        "Invalid request or insufficient balance"
// @Failure      401      {object}  dto.ErrorResponse        "Unauthorized"
// @Failure      500      {object}  dto.ErrorResponse        "Internal server error"
// @Router       /transaction/transfer [post]
func (tc *TransactionController) Transfer(ctx *gin.Context) {
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

	var req dto.TransferHTTPRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "invalid request body"})
		return
	}

	result, err := tc.transactionService.Transfer(ctx.Request.Context(), dto.TransferServiceRequest{
		UserID:     claims.Id,
		ReceiverID: req.ReceiverID,
		Amount:     req.Amount,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, result)
}
