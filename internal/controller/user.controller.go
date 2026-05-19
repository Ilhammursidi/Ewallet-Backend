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

func (u *UserController) GetProfile(ctx *gin.Context) {
	token, _ := ctx.Get("claims")
	claims := token.(pkg.Claims)
	user, err := u.userService.GetUserProfile(ctx.Request.Context(), claims.Id)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, dto.Response{
			Message: "internal error",
			Success: false,
			Error:   "internal server error",
		})
		return
	}
	ctx.JSON(http.StatusOK, dto.Response{
		Message: "OK",
		Success: true,
		Data:    user,
	})
}

func (u *UserController) GetDashboardInfo(ctx *gin.Context) {
	token, _ := ctx.Get("claims")
	claims := token.(pkg.Claims)

	data, err := u.userService.GetMoneyInfo(ctx.Request.Context(), claims.Id)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, dto.Response{
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

func (u *UserController) EditUserProfile(ctx *gin.Context) {
	token, _ := ctx.Get("claims")
	claims := token.(pkg.Claims)
	var body dto.EditProfileRequest
	if err := ctx.ShouldBindWith(&body, binding.JSON); err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, dto.Response{
			Message: "internal error",
			Success: false,
			Error:   "internal server error",
		})
		return
	}
	data, err := u.userService.EditProfile(ctx.Request.Context(), claims.Id, body)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, dto.Response{
			Message: "internal error",
			Success: false,
			Error:   "internal server error",
		})
		return
	}
	ctx.JSON(http.StatusCreated, dto.Response{
		Message: "update users success",
		Success: true,
		Data:    data,
	})
}
