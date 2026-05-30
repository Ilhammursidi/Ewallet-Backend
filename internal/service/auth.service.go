package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/ewallet-backend/internal/dto"
	"github.com/ewallet-backend/internal/repository"
	"github.com/ewallet-backend/pkg"
	"github.com/redis/go-redis/v9"
)

type PasswordService interface {
	RequestReset(ctx context.Context, req dto.ForgotPasswordRequest) error
	VerifyToken(ctx context.Context, token string) error
	ResetPassword(ctx context.Context, req dto.ResetPasswordRequest) error
}

type AuthService struct {
	authRepo      *repository.AuthRepository
	blacklistRepo *repository.BlacklistRepository
	rdb           *redis.Client
	cacheRepo     repository.CacheRepository
}

func NewAuthService(authRepo *repository.AuthRepository, blacklistRepo *repository.BlacklistRepository, rdb *redis.Client, cacheRepo repository.CacheRepository) *AuthService {
	return &AuthService{
		authRepo:      authRepo,
		blacklistRepo: blacklistRepo,
		rdb:           rdb,
		cacheRepo:     cacheRepo,
	}
}

func (a *AuthService) RegisterUser(ctx context.Context, user dto.NewUser) (dto.User, error) {
	rkey := "ilhammursidi:register:" + user.Email
	cache, err := a.rdb.Get(ctx, rkey).Result()
	log.Println("ini cache", cache)
	if err == nil {
		var cachedUser dto.User
		if err := json.Unmarshal([]byte(cache), &cachedUser); err == nil {
			log.Println("register: cache hit")
			return cachedUser, nil
		}
	} else if err != redis.Nil {
		return dto.User{}, err
	}

	log.Println("register: cache miss")

	var hc pkg.HashConfig
	hc.UseRecommended()
	hashedPwd := hc.GenHash(user.Password)

	newUser, err := a.authRepo.NewUser(ctx, user.Email, hashedPwd)
	if err != nil {
		return dto.User{}, err
	}

	userResponse := dto.User{
		Id:         newUser.Id,
		Email:      newUser.Email,
		Password:   newUser.Password,
		Created_at: *newUser.Created_at,
	}

	jsonData, err := json.Marshal(userResponse)
	if err != nil {
		log.Println("failed to marshal user for cache:", err)
	} else {
		err = a.rdb.Set(ctx, rkey, jsonData, 30*time.Minute).Err()
		if err != nil {
			log.Println("failed to save to redis:", err)
		}
	}

	return userResponse, nil
}

func (a *AuthService) LoginUser(ctx context.Context, user dto.NewUser) (string, error) {
	login, err := a.authRepo.GetUserByEmail(ctx, user.Email)
	if err != nil {
		return "", err
	}

	var hc pkg.HashConfig
	hc.UseRecommended()

	if err := hc.Compare(user.Password, login.Password); err != nil {
		return "", errors.New("invalid email or password")
	}

	claims := pkg.NewClaims(login.Id, user.Email)
	token, err := claims.GenJWT()
	if err != nil {
		return "", err
	}

	sessionKey := "ilhammursidi:session"
	err = a.rdb.Set(ctx, sessionKey, token, 24*time.Hour).Err()
	if err != nil {
		log.Println("failed to save session to redis:", err)
	}

	return token, nil
}

func (a *AuthService) CekPinUser(ctx context.Context, email string) (bool, error) {
	return a.authRepo.CekPinUser(ctx, email)
}

func (a *AuthService) CreatePin(ctx context.Context, userId int, body dto.SetPin) error {
	var hc pkg.HashConfig
	hc.UseRecommended()
	hashedPin := hc.GenHash(body.Pin)

	return a.authRepo.UpdatePin(ctx, userId, hashedPin)
}

func (a *AuthService) LogoutUser(ctx context.Context, token string) error {
	blacklistKey := "ilhammursidi:blacklist:" + token

	fmt.Println("ini di auth serv redis", blacklistKey)
	err := a.rdb.Set(ctx, blacklistKey, "true", 10*time.Hour).Err()
	if err != nil {
		return err
	}
	return nil
	// return a.authRepo.AddToBlackList(ctx, token)

}

func (a *AuthService) RequestReset(ctx context.Context, req dto.ForgotPasswordRequest) error {
	user, err := a.authRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return err // User tidak terdaftar
	}

	// Buat token random yang aman
	b := make([]byte, 16)
	rand.Read(b)
	token := hex.EncodeToString(b)

	err = a.cacheRepo.SaveResetToken(ctx, token, user.Id, 15*time.Minute)
	if err != nil {
		return err
	}

	return nil
}

func (a *AuthService) VerifyToken(ctx context.Context, token string) error {
	userID, err := a.cacheRepo.GetUserIDByToken(ctx, token)
	if err != nil {
		return err
	}
	if userID == 0 {
		return errors.New("token sudah kedaluwarsa atau tidak valid")
	}
	return nil
}

func (a *AuthService) ResetPassword(ctx context.Context, req dto.ResetPasswordRequest) error {
	userID, err := a.cacheRepo.GetUserIDByToken(ctx, req.Token)
	if err != nil || userID == 0 {
		return errors.New("token tidak valid atau sudah kedaluwarsa")
	}

	var hc pkg.HashConfig
	hc.UseRecommended()
	hashedPwd := hc.GenHash(req.PasswordBaru)
	// // 2. Hash password baru
	// hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.PasswordBaru), bcrypt.DefaultCost)
	// if err != nil {
	// 	return err
	// }

	err = a.authRepo.UpdatePassword(ctx, userID, hashedPwd)
	if err != nil {
		return err
	}

	_ = a.cacheRepo.DeleteResetToken(ctx, req.Token)

	return nil
}
