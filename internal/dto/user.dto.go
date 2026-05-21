package dto

type CashFlow struct {
	Balance int `json:"balance"`
	Income  int `json:"income"`
	Expense int `json:"expense"`
}

type EditProfileRequest struct {
	Fullname     string `json:"fullname"`
	Phone_number string `json:"phone"`
	Photo_path   string `json:"photo"`
}

type EditProfileResponse struct {
	Fullname     string `json:"fullname"`
	Email        string `json:"email"`
	Phone_number string `json:"phone"`
	Photo_path   string `json:"photo"`
}

type EditUserPinRequest struct {
	OldPin string `json:"oldpin"`
	NewPin string `json:"newpin" binding:"required,min=6"`
}
