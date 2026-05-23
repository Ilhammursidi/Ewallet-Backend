package dto

// @name CashFlow
type CashFlow struct {
	Balance int `json:"balance"`
	Income  int `json:"income"`
	Expense int `json:"expense"`
}

// @name EditProfileRequest
type EditProfileRequest struct {
	Fullname     string `json:"fullname"`
	Phone_number string `json:"phone"`
	Photo_path   string `json:"photo"`
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
	Pin string `json:"pin"`
}

// @name UserProfileResponse
type UserProfileResponse struct {
	Fullname     string `json:"fullname"`
	Email        string `json:"email"`
	Phone_number string `json:"phone"`
	Photo_path   string `json:"photo"`
}
