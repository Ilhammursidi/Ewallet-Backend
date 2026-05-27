package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"strconv"

	"github.com/ewallet-backend/internal/dto"
	"github.com/ewallet-backend/internal/repository"
	"github.com/ewallet-backend/pkg"
	"github.com/redis/go-redis/v9"
)

type UserService struct {
	userRepo *repository.UserRepository
	rdb      *redis.Client
}

func NewUserService(userRepo *repository.UserRepository, rdb *redis.Client) *UserService {
	return &UserService{
		userRepo: userRepo,
		rdb:      rdb,
	}
}

func (u *UserService) GetUserProfile(ctx context.Context, id int) (dto.User, error) {
	res, err := u.userRepo.GetProfileId(ctx, id)
	if err != nil {
		return dto.User{}, err
	}
	return dto.User{
		Id:           res.Id,
		Email:        res.Email,
		Fullname:     *res.Fullname,
		Photo_path:   *res.Photo_path,
		Phone_number: *res.Phone_number,
		Created_at:   *res.Created_at,
		Updated_at:   res.Updated_at,
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

func (u *UserService) EditProfile(ctx context.Context, id int, req dto.EditProfileRequest, photoPath *string) (dto.EditProfileResponse, error) {
	user, err := u.userRepo.EditProfile(ctx, id, req.Fullname, req.Phone_number, photoPath)
	if err != nil {
		return dto.EditProfileResponse{}, err
	}
	var fullname, phone, photo string

	if req.Fullname != nil {
		fullname = *user.Fullname
	}

	if req.Phone_number != nil {
		phone = *user.Phone_number
	}

	if req.Photo_path != nil {
		photo = *user.Photo_path
	}

	return dto.EditProfileResponse{
		Fullname:     fullname,
		Email:        user.Email,
		Phone_number: phone,
		Photo_path:   photo,
	}, nil
}

func (u *UserService) EditPin(ctx context.Context, id int, req dto.EditUserPinRequest) error {
	return u.userRepo.EditUserPin(ctx, id, &req.NewPin)
}

func (u *UserService) EditPassword(ctx context.Context, id int, req dto.EditPasswordRequest) error {
	if req.NewPassword != req.ConfrimPassword {
		return errors.New("passwords are not the same")
	}

	var hc pkg.HashConfig
	hc.UseRecommended()
	hashedPassword := hc.GenHash(req.NewPassword)

	return u.userRepo.EditPassword(ctx, id, &hashedPassword)
}

func (u *UserService) GetTransactionReport(ctx context.Context, id int, req dto.TransactionReportRequest) ([]dto.TransactionReportDTO, error) {
	validPeriods := map[string]bool{
		"week":  true,
		"month": true,
		"year":  true,
	}

	if !validPeriods[req.Period] {
		return nil, fmt.Errorf("invalid period: %s (week/month/year)", req.Period)
	}

	data, err := u.userRepo.GetTransactionReport(ctx, id, req.Period)
	if err != nil {
		return nil, fmt.Errorf("GetTransactionReport: %w", err)
	}

	return data, nil
}

func (u *UserService) TransactionHistory(ctx context.Context, id int, req dto.TransactionHistoryRequest) ([]dto.TransactionHistoryDTO, dto.PaginationMetaData, error) {
	page := 1
	if req.Page != "" {
		if p, err := strconv.Atoi(req.Page); err == nil && p > 0 {
			page = p
		}
	}
	log.Println("TransactionHistory - calling repo, id:", id, "req:", req)
	data, err := u.userRepo.GetTransactionHistory(ctx, id, req)
	if err != nil {
		return nil, dto.PaginationMetaData{}, err
	}
	log.Println("TransactionHistory - repo success, len:", len(data))
	if len(data) == 0 {
		return []dto.TransactionHistoryDTO{}, dto.PaginationMetaData{}, nil
	}

	totalData := data[0].TotalCount
	limit := 10
	totalPage := int(math.Ceil(float64(totalData) / float64(limit)))

	prevLink := fmt.Sprintf("/transactions/?search=%s&page=%d", req.Search, page-1)
	nextLink := fmt.Sprintf("/transactions/?search=%s&page=%d", req.Search, page+1)

	meta := dto.PaginationMetaData{
		TotalPages: totalPage,
		TotalData:  totalData,
		NextLink:   nextLink,
		PrevLink:   prevLink,
	}

	return data, meta, nil
}
