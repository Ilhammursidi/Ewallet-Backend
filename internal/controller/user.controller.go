package controller

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	response "github.com/ewallet-backend/internal/Response"
	"github.com/ewallet-backend/internal/dto"
	"github.com/ewallet-backend/internal/service"
	"github.com/ewallet-backend/pkg"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type UserController struct {
	userService *service.UserService
}

func NewUserController(userService *service.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

// User Profile
//
//		@Summary		Get User Profile
//		@Description	Retrieve detailed profile information for the authenticated user
//		@Tags			users
//		@Accept			json
//		@Produce		json
//	 @Security		ApiKeyAuth
//		@Success		200	{object}	dto.Response	"Success retrieve profile"
//		@Failure		500	{object}	dto.ErrorResponse						"Internal Server Error"
//		@Router			/users/profile [get]
func (u *UserController) GetProfile(ctx *gin.Context) {
	token, _ := ctx.Get("claims")
	claims := token.(*pkg.Claims)
	user, err := u.userService.GetUserProfile(ctx.Request.Context(), claims.Id)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "internal error",
			Success: false,
			Error:   "internal server error",
		})
		return
	}
	ctx.JSON(http.StatusOK, dto.Response{
		Message: "OK",
		Success: true,
		Data: dto.UserProfileResponse{
			Id:           user.Id,
			Fullname:     user.Fullname,
			Email:        user.Email,
			Photo_path:   user.Photo_path,
			Phone_number: user.Phone_number,
		},
	})
}

// Dashboard Info
//
//		@Summary		Get Dashboard Info
//		@Description	Retrieve balance, income, and expense summary for the user dashboard
//		@Tags			users
//		@Accept			json
//		@Produce		json
//	 @Security		ApiKeyAuth
//		@Success		200	{object}	dto.Response	"Success retrieve balance, income, and expense"
//		@Failure		500	{object}	dto.ErrorResponse			"Internal Server Error"
//		@Router			/users/dashboard [get]
func (u *UserController) GetDashboardInfo(ctx *gin.Context) {
	token, _ := ctx.Get("claims")
	claims := token.(*pkg.Claims)

	data, err := u.userService.GetMoneyInfo(ctx.Request.Context(), claims.Id)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "internal error",
			Success: false,
			Error:   "internal server error",
		})
		return
	}
	log.Println(data)
	ctx.JSON(http.StatusOK, dto.Response{
		Message: "get info success",
		Success: true,
		Data:    data,
	})
}

// Edit User Profile
//
//	@Summary		Edit User Profile
//	@Description	Update user profile fields like name, email, or profile photo
//	@Tags			users
//	@Accept			mpfd
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			fullname	formData	string	false	"Full name"
//	@Param			phone		formData	string	false	"Phone number"
//	@Param			photo		formData	file	false	"Profile photo (jpg, jpeg, png, max 2MB)"
//	@Success		200		{object}	dto.UserProfileResponse		"Profile updated successfully"
//	@Failure		400		{object}	dto.ErrorResponse	"Bad Request"
//	@Failure		500		{object}	dto.ErrorResponse	"Internal Server Error"
//	@Router			/users/profile [patch]
func (u *UserController) EditUserProfile(ctx *gin.Context) {
	token, _ := ctx.Get("claims")
	claims := token.(*pkg.Claims)

	var body dto.EditProfileRequest

	if err := ctx.ShouldBind(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "invalid request",
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	var photoPath *string

	if body.Photo_path.Filename != "" {

		const maxUploadSize = 2 * 1024 * 1024

		if body.Photo_path.Size > maxUploadSize {
			ctx.JSON(http.StatusUnprocessableEntity, dto.ErrorResponse{
				Message: "File too large",
				Success: false,
			})
			return
		}

		ext := strings.ToLower(
			filepath.Ext(body.Photo_path.Filename),
		)

		allowedExt := map[string]bool{
			".jpg":  true,
			".jpeg": true,
			".png":  true,
			".webp": true,
			".svg":  true,
		}

		if !allowedExt[ext] {
			ctx.JSON(http.StatusUnprocessableEntity, dto.ErrorResponse{
				Message: "Invalid file format",
				Success: false,
			})
			return
		}

		filename := fmt.Sprintf(
			"user_%d%s",
			time.Now().UnixNano(),
			ext,
		)

		dst := filepath.Join(
			"public",
			"img",
			"profiles",
			filename,
		)

		if err := ctx.SaveUploadedFile(&body.Photo_path, dst); err != nil {
			ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Message: "Failed to save image",
				Success: false,
			})
			return
		}

		generatedURL := "img/profiles/" + filename
		photoPath = &generatedURL
	}

	data, err := u.userService.EditProfile(
		ctx.Request.Context(),
		claims.Id,
		body,
		photoPath,
	)

	if err != nil {
		log.Println(err.Error())

		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "internal error",
			Success: false,
			Error:   "internal server error",
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Message: "Profile updated successfully",
		Success: true,
		Data:    data,
	})
}

// Edit User Pin
//
//		@Summary		Change User PIN
//		@Description	Change old transaction PIN to a new 6-digit secure PIN
//		@Tags			users
//		@Accept			json
//		@Produce		json
//	 	@Security		ApiKeyAuth
//		@Param			body	body		dto.EditUserPinRequest	true	"Old and New PIN payload"
//		@Success		201		{object}	dto.Response		            "PIN updated successfully"
//		@Failure		500		{object}	dto.ErrorResponse		        "Internal Server Error"
//		@Router			/users/profile/change-pin [patch]
func (u *UserController) EditUserPin(ctx *gin.Context) {
	token, _ := ctx.Get("claims")
	claims := token.(*pkg.Claims)
	var body dto.EditUserPinRequest
	if err := ctx.ShouldBindWith(&body, binding.JSON); err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "invalid PIN data",
			Success: false,
			Error:   err.Error(),
		})
		return
	}
	if err := u.userService.EditPin(ctx.Request.Context(), claims.Id, body); err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Failed to edit PIN",
			Success: false,
			Error:   "error",
		})
		return
	}
	ctx.JSON(http.StatusCreated, dto.Response{
		Message: "update userPin success",
		Success: true,
		Data:    "",
	})
}

// Edit Password
//
//	@Summary		Change User Password
//	@Description	Change current password to a new password for account security
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			body			body		dto.EditPasswordRequest		true	"Old and New password payload"
//	@Success		201				{object}	dto.Response			"Password updated successfully"
//	@Failure		500				{object}	dto.ErrorResponse			"Internal Server Error"
//	@Router			/users/profile/change-password [patch]
func (u *UserController) EditPassword(ctx *gin.Context) {
	token, _ := ctx.Get("claims")
	claims := token.(*pkg.Claims)
	var body dto.EditPasswordRequest
	if err := ctx.ShouldBindWith(&body, binding.JSON); err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "internal error",
			Success: false,
			Error:   "internal server error",
		})
		return
	}
	if err := u.userService.EditPassword(ctx.Request.Context(), claims.Id, body); err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: err.Error(),
			Success: false,
			Error:   "internal server error",
		})
		return
	}
	ctx.JSON(http.StatusCreated, dto.Response{
		Message: "update password success",
		Success: true,
		Data:    "",
	})
}

// Get Transaction Report
//
//		@Summary		Get Transaction Report
//		@Description	Get aggregated income and expense report for the authenticated user, grouped by the selected time period
//		@Tags			users
//		@Accept			json
//		@Produce		json
//		@Param			period  query     string  true  "Time period to group by"  Enums(week, month, year)  example("month")
//	    @Security		ApiKeyAuth
//		@Success		200	{object}	dto.Response		"Get Data Successs"
//		@Failure		500	{object}	dto.ErrorResponse	"Internal Server Error"
//		@Router			/users/transaction-report [get]
func (uc *UserController) TransactionReportGraph(ctx *gin.Context) {
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

	timePeriod := ctx.DefaultQuery("period", "week")

	data, err := uc.userService.GetTransactionReport(ctx.Request.Context(), claims.Id, dto.TransactionReportRequest{
		Period: timePeriod,
	})
	if err != nil {
		if strings.Contains(err.Error(), "invalid period") {
			ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Message: err.Error(),
				Success: false,
				Error:   "Bad Request",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: err.Error(),
			Success: false,
			Error:   "Internal Server Error",
		})
		return
	}

	// log.Println(data)
	ctx.JSON(http.StatusOK, dto.Response{
		Message: "Transaction report retrieved successfully",
		Success: true,
		Data:    data,
	})
}

// GetTransactionHistory
//
// @Summary      Get transaction history
// @Description  Get a paginated list of the authenticated user's transaction history, optionally filtered by a search keyword
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        search  query     string  false  "Search by receiver name, payment method, or description"  example("Naruto")
// @Param        page    query     string  false  "Page number"                                               example("1")
// @Success      200     {object}  dto.ResponseSuccess{data=[]dto.GetTransactionHistory}
// @Failure      400     {object}  dto.Response
// @Failure      500     {object}  dto.Response
// @Router       /users/transactions [get]
func (u *UserController) GetTransactionHistory(ctx *gin.Context) {
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

	var body dto.TransactionHistoryRequest
	if err := ctx.ShouldBindQuery(&body); err != nil {
		response.JSONBadRequest(ctx)
		return
	}

	data, metaData, err := u.userService.TransactionHistory(ctx.Request.Context(), claims.Id, body)
	if err != nil {
		log.Println("TransactionHistory error:", err)
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: err.Error(),
			Success: false,
			Error:   "Internal Server Error",
		})
		return
	}

	response.SuccessWithMetaData(ctx, 200, "Get data success", data, metaData)
}

func (u *UserController) GetPin(ctx *gin.Context) {

}
