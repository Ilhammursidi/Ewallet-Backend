package dto

import "time"

type User struct {
	Id           int        `json:"id"`
	Email        string     `json:"email"`
	Password     string     `json:"password,omitempty"`
	Pin          string     `json:"pin"`
	Fullname     string     `json:"fullname"`
	Photo_path   string     `json:"photo"`
	Phone_number string     `json:"phone"`
	IsVerified   bool       `json:"is_verified"`
	Created_at   time.Time  `json:"create_at"`
	Updated_at   *time.Time `json:"updated_at"`
	Deleted_at   *time.Time `json:"deleted_at"`
}

type NewUser struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}
