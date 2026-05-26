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

// User Register
//
//	@Summary		Register a user
//	@Description	create user account
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param 			body body dto.NewUser true "register payload"
//	@Success		201	{object}	dto.Response
//	@Failure		500	{object}	dto.ErrorResponse
//	@Router			/auth/register [post]
func (a *AuthController) Register(ctx *gin.Context) {
	var body dto.NewUser
	if err := ctx.ShouldBindWith(&body, binding.JSON); err != nil {
		log.Println("Error: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Error",
			Success: false,
			Error:   "Internal Server Error",
		})
		return
	}
	res, err := a.authService.RegisterUser(ctx.Request.Context(), body)
	if err != nil {
		log.Println("Error: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Error",
			Success: false,
			Error:   "Email already register",
		})
		return
	}
	log.Println("response data", res)
	ctx.JSON(http.StatusCreated, dto.Response{
		Data:    "",
		Message: "Register Success",
		Success: true,
	})
}

// User Login
//
//	@Summary		Login a user
//	@Description	Login into user account
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param 			body body dto.NewUser true "login payload"
//	@Success		200	{object}	dto.Response
//	@Failure		500	{object}	dto.ErrorResponse
//	@Router			/auth/login [post]
func (a *AuthController) Login(ctx *gin.Context) {
	var body dto.NewUser
	if err := ctx.ShouldBindWith(&body, binding.JSON); err != nil {
		log.Println("Error: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Error",
			Success: false,
			Error:   "Internal Server Error",
		})
		return
	}
	token, err := a.authService.LoginUser(ctx.Request.Context(), body)
	if err != nil {
		log.Println("Error: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Error",
			Success: false,
			Error:   "Invalid email or password",
		})
		return
	}
	log.Println("response data", token)
	ctx.JSON(http.StatusOK, dto.Response{
		Data: dto.LoginResponse{
			Token: token,
		},
		Message: "Login Success",
		Success: true,
	})
}

// Create PIN
//
//		@Summary		Create a PIN
//		@Description	Create PIN for user account
//		@Tags			auth
//		@Accept			json
//		@Produce		json
//	    @security       ApiKeyAuth
//		@Param 			body body dto.SetPin true "set PIN payload"
//		@Success		201	{object}	dto.Response
//		@Failure		401 {object}	dto.ErrorResponse
//		@Failure		400	{object}	dto.ErrorResponse
//		@Router			/auth/enter-pin [patch]
func (a *AuthController) CreatePin(ctx *gin.Context) {
	claims, exists := ctx.Get("claims")
	log.Println("claims exists:", exists)

	if !exists {
		ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Message: "Unauthorized",
			Success: false,
			Error:   "user tidak ditemukan",
		})
		return
	}

	userClaims, ok := claims.(*pkg.Claims)
	log.Println("cast ok:", ok)
	log.Println("userId:", userClaims.Id)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Message: "Unauthorized",
			Success: false,
			Error:   "invalid token claims",
		})
		return
	}

	var body dto.SetPin
	if err := ctx.ShouldBindWith(&body, binding.JSON); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Validation Error",
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	if err := a.authService.CreatePin(ctx.Request.Context(), userClaims.Id, body); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Error",
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, dto.Response{
		Message: "PIN created successfully",
		Success: true,
		Data:    nil,
	})
}

// Logout User
//
//		@Summary		Logout User
//		@Description	Revoke user session and invalidate the JWT token
//		@Tags			auth
//		@Accept			json
//		@Produce		json
//	 	@Security		ApiKeyAuth
//		@Success		200				{object}	dto.Response	"Logout Success"
//		@Failure		500				{object}	dto.ErrorResponse	"Internal Server Error"
//		@Router			/auth/logout [post]
func (a *AuthController) Logout(ctx *gin.Context) {
	bearerToken := ctx.GetHeader("Authorization")
	log.Println("bearerToken logout", bearerToken)
	token := strings.Split(bearerToken, " ")[1]
	log.Println("logout_controller : ", token)

	if err := a.authService.LogoutUser(ctx.Request.Context(), token); err != nil {
		log.Println("c logout: ", err)
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Logout Failed",
			Success: false,
			Error:   "Internal() Server Error",
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Message: "Logout Success",
		Success: true,
		Data:    nil,
	})
}
