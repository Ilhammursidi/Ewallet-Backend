package service

import (
	"context"

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

func (a AuthService) GetUserProfile(ctx context.Context, id int) (dto.User, error) {
	user, err := a.authRepo.GetUserById(ctx, id)
	if err != nil {
		return dto.User{}, err
	}
	return dto.User{
		Id:         user.Id,
		Email:      user.Email,
		Created_at: user.Created_at,
	}, nil
}
