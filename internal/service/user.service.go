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

	return dto.EditProfileResponse{
		Fullname:     user.Fullname,
		Email:        user.Email,
		Phone_number: user.Phone_number,
		Photo_path:   user.Photo_path,
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

func (u *UserService) CheckUserPin(ctx context.Context, id int) (dto.CheckPinResponse, error) {
	user, err := u.userRepo.CheckPin(ctx, id)
	if err != nil {
		return dto.CheckPinResponse{}, err
	}
	return dto.CheckPinResponse{
		HasPin: user != nil,
	}, nil
}

func (u *UserService) GetTransactionReport(ctx context.Context, id int, req dto.TransactionReportRequest) ([]dto.TransactionReportDTO, error) {
	data, err := u.userRepo.GetTransactionReport(ctx, id, req.Period)
	if err != nil {
		return nil, err
	}
	var transactions []dto.TransactionReportDTO
	for _, transaction := range data {
		transactions = append(transactions, dto.TransactionReportDTO{
			Period:  transaction.Period,
			Income:  transaction.Income,
			Expense: transaction.Expense,
		})
	}
	return transactions, nil
}

func (u *UserService) TransactionHistory(ctx context.Context, id int, req dto.TransactionHistoryRequest) ([]dto.GetTransactionHistory, dto.PaginationMetaData, error) {
	data, err := u.userRepo.GetTransactionHistory(ctx, id, req)
	if err != nil {
		return nil, dto.PaginationMetaData{}, err
	}
	log.Println("apakah work:", data)
	if len(data) == 0 {
		return []dto.GetTransactionHistory{}, dto.PaginationMetaData{}, nil
	}

	totalData := data[0].TotalCount
	limit := 10
	totalPage := int(math.Ceil(float64(totalData) / float64(limit)))

	page, err := strconv.Atoi(req.Page)
	if err != nil {
		return nil, dto.PaginationMetaData{}, err
	}

	var users []dto.GetTransactionHistory
	prevLink := fmt.Sprintf("/transactions/?search=%s&page=%d", req.Search, page-1)
	nextLink := fmt.Sprintf("/transactions/?search=%s&page=%d", req.Search, page+1)
	for _, user := range data {
		users = append(users, dto.GetTransactionHistory{
			TransactionID:     user.TransactionID,
			Amount:            user.Amount,
			Flow_type:         user.Flow_type,
			Type:              user.Type,
			Status:            user.Status,
			CreatedAt:         user.CreatedAt,
			Description:       user.Description,
			ReceiverName:      user.ReceiverName,
			PaymentMethodName: user.PaymentMethodName,
		})
	}
	metaDataPagination := dto.PaginationMetaData{
		TotalPages: totalPage,
		TotalData:  totalData,
		NextLink:   nextLink,
		PrevLink:   prevLink,
	}
	return users, metaDataPagination, nil
}
