package service

import (
	"context"
	"errors"
	"log"

	"github.com/ewallet-backend/internal/dto"
	"github.com/ewallet-backend/internal/repository"
	"github.com/ewallet-backend/pkg"
)

type AuthService struct {
	authRepo      *repository.AuthRepository
	blacklistRepo *repository.BlacklistRepository
}

func NewAuthService(authRepo *repository.AuthRepository, blacklistRepo *repository.BlacklistRepository) *AuthService {
	return &AuthService{
		authRepo:      authRepo,
		blacklistRepo: blacklistRepo,
	}
}

func (a *AuthService) RegisterUser(ctx context.Context, user dto.NewUser) (dto.User, error) {
	var hc pkg.HashConfig
	hc.UseRecommended()
	hashedPwd := hc.GenHash(user.Password)
	newUser, err := a.authRepo.NewUser(ctx, user.Email, hashedPwd)
	if err != nil {
		return dto.User{}, err
	}
	return dto.User{
		Id:         newUser.Id,
		Email:      newUser.Email,
		Password:   newUser.Password,
		Created_at: *newUser.Created_at,
	}, nil
}

func (a *AuthService) LoginUser(ctx context.Context, user dto.NewUser) (string, error) {
	login, err := a.authRepo.GetUserByEmail(ctx, user.Email)
	if err != nil {
		return "", err
	}
	var hc pkg.HashConfig

	if err := hc.Compare(user.Password, login.Password); err != nil {
		return "", err
	}
	claims := pkg.NewClaims(login.Id, user.Email)
	token, err := claims.GenJWT()
	if err != nil {
		return "", err
	}
	return token, nil
}

func (a *AuthService) CreatePin(ctx context.Context, userId int, body dto.SetPin) error {
	if body.Pin != body.ConfirmPin {
		return errors.New("pin and pin confirmation are not the same")
	}

	var hc pkg.HashConfig
	hc.UseRecommended()
	hashedPin := hc.GenHash(body.Pin)

	return a.authRepo.UpdatePin(ctx, userId, hashedPin)
}

func (a *AuthService) LogoutUser(ctx context.Context, token string) error {
	log.Println("service :", token)
	return a.authRepo.AddToBlackList(ctx, token)
}
