package controller

import (
	"log"
	"net/http"

	"github.com/ewallet-backend/internal/dto"
	"github.com/ewallet-backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type AuthController struct {
	authService *service.AuthService
}

func NewUserController(authService *service.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

func (u *AuthController) GetAll(ctx *gin.Context) {
	userlist, err := u.authService.PrintUser(ctx.Request.Context())
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, dto.Response{
			Message: "Internal Error",
			Success: false,
			Error:   "internal server error",
		})
		return
	}
	ctx.JSON(http.StatusOK, dto.Response{
		Data:    userlist,
		Message: "OK",
		Success: true,
	})
}

func (a *AuthController) Register(ctx *gin.Context) {
	var body dto.NewUser
	if err := ctx.ShouldBindWith(&body, binding.JSON); err != nil {
		log.Println("Error: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, dto.Response{
			Message: "Error",
			Success: false,
			Error:   "Internal Server Error",
		})
		return
	}
	res, err := a.authService.RegisterUser(ctx.Request.Context(), body)
	if err != nil {
		log.Println("Error: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, dto.Response{
			Message: "Error",
			Success: false,
			Error:   "Internal Server Error",
		})
		return
	}
	log.Println("response data", res)
	ctx.JSON(http.StatusCreated, dto.Response{
		Data:    res,
		Message: "Register Success",
		Success: true,
	})
}

func (a *AuthController) GetUser(ctx gin.Context) {

}
