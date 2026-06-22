package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"strconv"
	"time"

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
	rkey := fmt.Sprintf("ilhammursidi:profile:%d", id)
	log.Println("gpp kok:", rkey)
	cache, err := u.rdb.Get(ctx, rkey).Result()
	if err == nil {
		var cachedProfile dto.User
		if err := json.Unmarshal([]byte(cache), &cachedProfile); err == nil {
			log.Println("profile: cache hit")
			return cachedProfile, nil
		}
	} else if err != redis.Nil {
		return dto.User{}, nil
	}

	log.Println("profile: cache miss")

	res, err := u.userRepo.GetProfileId(ctx, id)
	if err != nil {
		return dto.User{}, err
	}

	userResponse := dto.User{
		Id:           res.Id,
		Email:        res.Email,
		Fullname:     *res.Fullname,
		Photo_path:   *res.Photo_path,
		Phone_number: *res.Phone_number,
		Created_at:   *res.Created_at,
		Updated_at:   res.Updated_at,
	}

	jsonData, err := json.Marshal(userResponse)
	if err != nil {
		log.Println("failed to marshal user for cache:", err)
	} else {
		err = u.rdb.Set(ctx, rkey, jsonData, 10*time.Hour).Err()
		if err != nil {
			log.Println("failed to save redis:", err)
		}
	}

	return userResponse, nil
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

	if user.Fullname != nil {
		fullname = *user.Fullname
	}

	if user.Phone_number != nil {
		phone = *user.Phone_number
	}

	if user.Photo_path != nil {
		photo = *user.Photo_path
	}

	rkey := fmt.Sprintf("ilhammursidi:profile:%d", id)
	_ = u.rdb.Del(ctx, rkey).Err()

	return dto.EditProfileResponse{
		Fullname:     fullname,
		Email:        user.Email,
		Phone_number: phone,
		Photo_path:   photo,
	}, nil
}

func (u *UserService) EditPin(
	ctx context.Context,
	id int,
	req dto.EditUserPinRequest,
) error {

	user, err := u.userRepo.GetUserByIdUser(ctx, id)
	if err != nil {
		return err
	}

	var hc pkg.HashConfig
	hc.UseRecommended()

	err = hc.Compare(req.OldPin, user.Pin)
	if err != nil {
		return errors.New("old pin invalid")
	}
	err = hc.Compare(req.NewPin, user.Pin)
	if err == nil {
		return errors.New("the new pin cannot be the same as the old pin")
	}

	hashedPin := hc.GenHash(req.NewPin)

	return u.userRepo.EditUserPin(ctx, id, &hashedPin)
}

func (u *UserService) EditPassword(ctx context.Context, id int, req dto.EditPasswordRequest) error {
	currentUser, err := u.userRepo.GetUserByIdUser(ctx, id)
	if err != nil {
		return errors.New("user not found")
	}
	log.Println("apakah cocok", currentUser.Password)

	var hc pkg.HashConfig
	hc.UseRecommended()

	err = hc.Compare(req.OldPassword, currentUser.Password)
	if err != nil {
		return errors.New("old password is wrong")
	}

	err = hc.Compare(req.NewPassword, currentUser.Password)
	if err == nil {
		return errors.New("The new password cannot be the same as the old password")
	}

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

func (u *UserService) TransactionHistory(ctx context.Context, id int, req dto.TransactionHistoryRequest) ([]dto.TransactionHistoryDTO, dto.PaginationMetaResponse, error) {
	page := 1
	if req.Page != "" {
		if p, err := strconv.Atoi(req.Page); err == nil && p > 0 {
			page = p
		}
	}
	log.Println("TransactionHistory - calling repo, id:", id, "req:", req)
	data, err := u.userRepo.GetTransactionHistory(ctx, id, req)
	if err != nil {
		return nil, dto.PaginationMetaResponse{}, err
	}
	log.Println("TransactionHistory - repo success, len:", len(data))
	if len(data) == 0 {
		return []dto.TransactionHistoryDTO{}, dto.PaginationMetaResponse{}, nil
	}
	totalTransactions, err := u.userRepo.GetAllTransactionId(ctx, id)
	totalData := totalTransactions
	limit := 10
	totalPage := int(math.Ceil(float64(totalData) / float64(limit)))

	var NextPage string
	Domain := "http://localhost:8081"
	if page < totalPage {
		Next := page + 1
		NextPage = fmt.Sprintf("%s/transactions/?search=%s&page=%d", Domain, req.Search, Next)
	} else {
		NextPage = ""
	}

	var PrevPage string
	if page > totalPage {
		Prev := page - 1
		PrevPage = fmt.Sprintf("%s/transactions/?search=%s&page=%d", Domain, req.Search, Prev)
	} else {
		PrevPage = ""
	}

	meta := dto.PaginationMetaResponse{
		Page:       page,
		Total_Data: totalData,
		Total_Page: totalPage,
		Next_Page:  NextPage,
		Prev_Page:  PrevPage,
	}

	return data, meta, nil
}
