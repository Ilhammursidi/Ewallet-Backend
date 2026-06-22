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
//	@Failure 		409 {object}	dto.ErrorResponse
//	@Failure		500	{object}	dto.ErrorResponse
//	@Router			/auth/register [post]
func (a *AuthController) Register(ctx *gin.Context) {
	var body dto.RegisterDto
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
		ctx.JSON(http.StatusConflict, dto.ErrorResponse{
			Message: "Email already exist",
			Success: false,
			Error:   "Conflict",
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
//	@Success		200	{object}	dto.Response         "Login Success"
//	@Failure		500	{object}	dto.ErrorResponse	"Internal Server Error"
//	@Failure		401	{object}	dto.ErrorResponse	"Invalid Email or Password"
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
		ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Message: "Error",
			Success: false,
			Error:   "Invalid email or password",
		})
		return
	}
	hasPin, err := a.authService.CekPinUser(ctx.Request.Context(), body.Email)
	log.Println("haspin cont", hasPin)
	ctx.JSON(http.StatusOK, dto.Response{
		Data: dto.LoginResponse{
			Token:  token,
			HasPin: hasPin,
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
//		@Success		201	{object}	dto.Response       "Create PIN Success"
//		@Failure		401 {object}	dto.ErrorResponse	"Invalid Token Claims","User not found"
//		@Failure		400	{object}	dto.ErrorResponse	"validation error"
//		@Router			/auth/enter-pin [patch]
func (a *AuthController) CreatePin(ctx *gin.Context) {
	claims, exists := ctx.Get("claims")
	log.Println("claims exists:", exists)

	if !exists {
		ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Message: "Unauthorized",
			Success: false,
			Error:   "User Not Found",
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

// Request Forgot Password
// @Summary      Get link for forgot password
// @Description  Check email and save temporary token to Redis
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body dto.ForgotPasswordRequest true "Input Email"
// @Success      200 {object} map[string]string "Message: Instruksi terkirim"
// @Failure      400 {object} map[string]string "Error message"
// @Router       /auth/forgot-password [post]
func (a *AuthController) RequestForgotPassword(c *gin.Context) {
	var req dto.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[DEBUG] Gagal mengurai data JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Printf("[DEBUG] Data email berhasil diterima: %s", req.Email)

	err := a.authService.RequestReset(c.Request.Context(), req)
	if err != nil {
		log.Printf("[DEBUG] Error pada layer service: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Message: "Link reset_token has send",
	})
}

// Verify Reset Token
// @Summary      Verify Token Forgot Password
// @Description  Checking Redis whether the email link token is still active
// @Tags         auth
// @Produce      json
// @Param        token query string true "Token from email"
// @Success      200 {object} map[string]string "Message: Token valid"
// @Failure      400 {object} map[string]string "Error message"
// @Router       /auth/verify-reset-token [post]
func (a *AuthController) VerifyResetToken(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "token dibutuhkan"})
		return
	}

	err := a.authService.VerifyToken(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Token valid, silakan ganti password"})
}

// Reset Password
// @Summary      Save New Password
// @Description  Change Password in database and delete token from redis
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body dto.ResetPasswordRequest true "Input New Password"
// @Success      200 {object} map[string]string "Message: Change Password successfull"
// @Failure      400 {object} map[string]string "Error message"
// @Router       /auth/reset-password [post]
func (a *AuthController) ResetPassword(c *gin.Context) {
	var req dto.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := a.authService.ResetPassword(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Change Password Success"})
}
