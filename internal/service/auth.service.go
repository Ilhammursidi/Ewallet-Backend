package service

import (
	"context"
	"errors"

	"github.com/ewallet-backend/internal/dto"
	"github.com/ewallet-backend/internal/repository"
	"github.com/ewallet-backend/pkg"
)

type AuthService struct {
	authRepo *repository.AuthRepository
}

func NewAuthService(authRepo *repository.AuthRepository) *AuthService {
	return &AuthService{
		authRepo: authRepo,
	}
}

func (a *AuthService) PrintUser(ctx context.Context) ([]dto.User, error) {
	res, err := a.authRepo.FetchUserList(ctx)
	if err != nil {
		return nil, err
	}
	var userList []dto.User
	for _, user := range res {
		userList = append(userList, dto.User{
			Id:         user.Id,
			Email:      user.Email,
			Password:   user.Password,
			Created_at: user.Created_at,
			Updated_at: user.Updated_at,
			Deleted_at: user.Deleted_at,
		})
	}
	return userList, nil
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
		Created_at: newUser.Created_at,
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
		return errors.New("pin dan konfirmasi pin tidak sama")
	}

	var hc pkg.HashConfig
	hc.UseRecommended()
	hashedPin := hc.GenHash(body.Pin)

	return a.authRepo.UpdatePin(ctx, userId, hashedPin)
}

func (a *AuthService) LogoutUser(ctx context.Context, token string) error {
	return a.authRepo.AddToBlackList(ctx, token)
}
