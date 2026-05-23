package controller

import (
	"log"
	"net/http"

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
//		@Tags			user
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
//		@Tags			user
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
	ctx.JSON(http.StatusOK, dto.Response{
		Message: "get info success",
		Success: true,
		Data:    data,
	})
}

// Edit User Profile
//
//		@Summary		Edit User Profile
//		@Description	Update user profile fields like name, email, or profile photo URL
//		@Tags			user
//		@Accept			json
//		@Produce		json
//	 	@Security		ApiKeyAuth
//		@Param			body	body		dto.EditProfileRequest					true	"Updated profile payload"
//		@Success		201		{object}	dto.Response	"Profile updated successfully"
//		@Failure		500		{object}	dto.ErrorResponse						"Internal Server Error"
//		@Router			/users/profile [patch]
func (u *UserController) EditUserProfile(ctx *gin.Context) {
	token, _ := ctx.Get("claims")
	claims := token.(*pkg.Claims)
	var body dto.EditProfileRequest
	if err := ctx.ShouldBindWith(&body, binding.JSON); err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "internal error",
			Success: false,
			Error:   "internal server error",
		})
		return
	}
	data, err := u.userService.EditProfile(ctx.Request.Context(), claims.Id, body)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "internal error",
			Success: false,
			Error:   "internal server error",
		})
		return
	}
	ctx.JSON(http.StatusCreated, dto.Response{
		Message: "Profile updated successfully",
		Success: true,
		Data: dto.UserProfileResponse{
			Fullname:     data.Fullname,
			Email:        data.Email,
			Photo_path:   data.Photo_path,
			Phone_number: data.Phone_number,
		},
	})
}

// Edit User Pin
//
//		@Summary		Change User PIN
//		@Description	Change old transaction PIN to a new 6-digit secure PIN
//		@Tags			user
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
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "internal error",
			Success: false,
			Error:   "internal server error",
		})
		return
	}
	if err := u.userService.EditPin(ctx.Request.Context(), claims.Id, body); err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "internal error",
			Success: false,
			Error:   "internal server error",
		})
		return
	}
	ctx.JSON(http.StatusCreated, dto.Response{
		Message: "update usersPin success",
		Success: true,
		Data:    "",
	})
}

// Edit Password
//
//	@Summary		Change User Password
//	@Description	Change current password to a new password for account security
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Secutiry		ApiKeyAuth
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
			Message: "internal error",
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

// CheckPin godoc
//
//		@Summary		Check PIN Status
//		@Description	Verify whether the user has already set a transaction PIN or not
//		@Tags			user
//		@Accept			json
//		@Produce		json
//	    @Security		ApiKeyAuth
//		@Success		200	{object}	dto.Response	"Success check PIN"
//		@Failure		500	{object}	dto.ErrorResponse					"Internal Server Error"
//		@Router			/users/check-pin [get]
func (u *UserController) CheckPin(ctx *gin.Context) {
	token, _ := ctx.Get("claims")
	claims := token.(*pkg.Claims)
	data, err := u.userService.CheckUserPin(ctx.Request.Context(), claims.Id)
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
		Message: "Check PIN success",
		Success: true,
		Data: dto.CheckPinResponse{
			Pin: data.Pin,
		},
	})
}
