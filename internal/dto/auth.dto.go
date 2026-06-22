package dto

import "time"

// @name User
type User struct {
	Id           int        `json:"id"`
	Email        string     `json:"email"`
	Password     string     `json:"password,omitempty"`
	Pin          string     `json:"pin,omitempty"`
	Fullname     string     `json:"fullname"`
	Photo_path   string     `json:"photo"`
	Phone_number string     `json:"phone"`
	IsVerified   bool       `json:"is_verified"`
	Created_at   time.Time  `json:"create_at"`
	Updated_at   *time.Time `json:"updated_at"`
	Deleted_at   *time.Time `json:"deleted_at"`
}

// @name NewUser
type NewUser struct {
	Email    string `json:"email" binding:"required,email" example:"ilham@mail.com"`
	Password string `json:"password" binding:"required,min=8" example:"ilham123"`
}

type RegisterDto struct {
	Email    string `json:"email" binding:"required,email" example:"ilham@mail.com"`
	Password string `json:"password" binding:"required,min=8" example:"ilham123"`
}

// @name SetPin
type SetPin struct {
	Pin string `json:"pin" binding:"required,len=6" example:"123456"`
}

// @name LoginResponse
type LoginResponse struct {
	Token  string `json:"token" example:"token..."`
	HasPin bool   `json:"has_pin" example:"false"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email" example:"user@mail.com"`
}

// Input saat user memasukkan password baru
type ResetPasswordRequest struct {
	Token        string `json:"token" binding:"required"`
	PasswordBaru string `json:"password" binding:"required,min=6" example:"rahasia123"`
}
