package controller

import (
	"log"
	"net/http"
	"strings"

	"github.com/ewallet-backend/internal/dto"
	"github.com/ewallet-backend/internal/service"
	"github.com/ewallet-backend/pkg"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type AuthController struct {
	authService *service.AuthService
}

func NewAuthController(authService *service.AuthService) *AuthController {
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

func (a *AuthController) Login(ctx *gin.Context) {
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
	token, err := a.authService.LoginUser(ctx.Request.Context(), body)
	if err != nil {
		log.Println("Error: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, dto.Response{
			Message: "Error",
			Success: false,
			Error:   "Internal Server Error",
		})
		return
	}
	log.Println("response data", token)
	ctx.JSON(http.StatusOK, dto.Response{
		Data: gin.H{
			"token": token,
		},
		Message: "Login Success",
		Success: true,
	})
}

func (a *AuthController) CreatePin(ctx *gin.Context) {
	claims, exists := ctx.Get("claims")
	log.Println("claims exists:", exists)

	if !exists {
		ctx.JSON(http.StatusUnauthorized, dto.Response{
			Message: "Unauthorized",
			Success: false,
			Error:   "user tidak ditemukan",
		})
		return
	}

	userClaims, ok := claims.(pkg.Claims)
	log.Println("cast ok:", ok)
	log.Println("userId:", userClaims.Id)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, dto.Response{
			Message: "Unauthorized",
			Success: false,
			Error:   "invalid token claims",
		})
		return
	}

	var body dto.SetPin
	if err := ctx.ShouldBindWith(&body, binding.JSON); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Message: "Validation Error",
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	if err := a.authService.CreatePin(ctx.Request.Context(), userClaims.Id, body); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Message: "Error",
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Message: "PIN berhasil dibuat",
		Success: true,
	})
}

func (a *AuthController) Logout(ctx *gin.Context) {
	bearerToken := ctx.GetHeader("Authorization")
	token := strings.Split(bearerToken, " ")[1]

	if err := a.authService.LogoutUser(ctx.Request.Context(), token); err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.Response{
			Message: "Logout Failed",
			Success: false,
			Error:   "Internal Server Error",
		})
		return
	}
	ctx.JSON(http.StatusOK, dto.Response{
		Message: "Logout Success",
		Success: true,
	})
}
