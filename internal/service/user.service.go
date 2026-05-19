package service

import (
	"context"

	"github.com/ewallet-backend/internal/dto"
	"github.com/ewallet-backend/internal/repository"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (u *UserService) GetUserProfile(ctx context.Context, id int) (dto.User, error) {
	res, err := u.userRepo.GetProfileId(ctx, id)
	if err != nil {
		return dto.User{}, err
	}
	return dto.User{
		Id:         res.Id,
		Email:      res.Email,
		Fullname:   res.Fullname,
		Photo_path: res.Photo_path,
		Phone_number: res.Phone_number,
		Created_at: res.Created_at,
		Updated_at: res.Updated_at,
	}, nil
}

func (u *UserService) GetMoneyInfo(ctx context.Context, id int) (dto.CashFlow, error) {
	dashboard, err := u.userRepo.GetMoneyAccountInfo(ctx, id)
	if err != nil {
		return dto.CashFlow{}, err
	}
	return dto.CashFlow{
		Balance: dashboard.Balance,
		Income:  dashboard.Income,
		Expense: dashboard.Expense,
	}, nil
}

func (u *UserService) EditProfile(ctx context.Context, id int, req dto.EditProfileRequest) (dto.EditProfileResponse, error) {
	user, err := u.userRepo.EditProfile(ctx, id, &req.Fullname, &req.Phone_number, &req.Photo_path)
	if err != nil {
		return dto.EditProfileResponse{}, err
	}
	return dto.EditProfileResponse{
		Fullname:     user.Fullname,
		Email:        user.Email,
		Phone_number: user.Phone_number,
		Photo_path:   user.Photo_path,
	}, nil
}
