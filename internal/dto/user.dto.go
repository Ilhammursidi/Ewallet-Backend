package dto

import (
	"mime/multipart"
	"time"
)

// @name CashFlow
type CashFlow struct {
	Balance int `json:"balance"`
	Income  int `json:"income"`
	Expense int `json:"expense"`
}

// @name EditProfileRequest
type EditProfileRequest struct {
	Fullname     *string               `form:"fullname"`
	Phone_number *string               `form:"phone"`
	Photo_path   *multipart.FileHeader `form:"photo" binding:"omitempty"`
}

// @name EditProfileResponse
type EditProfileResponse struct {
	Fullname     string `json:"fullname"`
	Email        string `json:"email"`
	Phone_number string `json:"phone"`
	Photo_path   string `json:"photo"`
}

// @name EditUserPinRequest
type EditUserPinRequest struct {
	OldPin     string `json:"oldpin" binding:"required,len=6"`
	NewPin     string `json:"newpin" binding:"required,len=6"`
	ConfirmPin string `json:"confirm_pin" binding:"required,len=6"`
}

// @name EditPasswordRequest
type EditPasswordRequest struct {
	OldPassword     string `json:"oldpassword" binding:"required,min=7"`
	NewPassword     string `json:"newpassword" binding:"required,min=7"`
	ConfrimPassword string `json:"confirm_password" binding:"required,min=7"`
}

// @name CheckPinResponse
type CheckPinResponse struct {
	HasPin bool `json:"has_pin"`
}

// @name UserProfileResponse
type UserProfileResponse struct {
	Fullname     string `json:"fullname"`
	Email        string `json:"email"`
	Phone_number string `json:"phone"`
	Photo_path   string `json:"photo"`
}

// @name TransactionReport
type TransactionReportDTO struct {
	Period  *time.Time `db:"period"`
	Income  int        `db:"income"`
	Expense int        `db:"expense"`
}

// @name TransactionReportRequest
type TransactionReportRequest struct {
	Period string `form:"period" binding:"required,oneof=week month year" example:"month"`
}

// @name TransactionHistoryRequest
type TransactionHistoryRequest struct {
	Page   string `form:"page" default:"1" example:"1"`
	Search string `form:"search" example:"naruto"`
}

type PaginationMetaData struct {
	TotalPages int    `json:"total_page,omitempty"`
	TotalData  int    `json:"total_data,omitempty"`
	PageSize   int    `json:"page_size,omitempty"`
	NextLink   string `json:"next_page,omitempty"`
	PrevLink   string `json:"prev_page,omitempty"`
}

type GetTransactionHistory struct {
	TransactionID     int       `json:"transaction_id"`
	Amount            int       `json:"amount"`
	Flow_type         string    `json:"type"`
	Type              string    `json:"activity_type"`
	Status            string    `json:"status"`
	CreatedAt         time.Time `json:"created_at"`
	Description       string    `json:"transfer_description"`
	ReceiverName      string    `json:"receiver_name"`
	PaymentMethodName string    `json:"payment_method_name"`
	TotalCount        int       `json:"total_count"`
}

type TransactionHistoryDTO struct {
	TransactionID     int       `json:"transaction_id"`
	Amount            int       `json:"amount"`
	Type              string    `json:"type"`
	FlowType          *string   `json:"flow_type"`
	Status            string    `json:"status"`
	CreatedAt         time.Time `json:"created_at"`
	PaymentMethodName string    `json:"payment_method_name"`
	ReceiverID        *int      `json:"receiver_id"`
	ReceiverName      string    `json:"receiver_name"`
	SenderName        string    `json:"sender_name"`
	TotalCount        int       `json:"total_count"`
}
