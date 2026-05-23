package dto

// @name ReceiverResponse
type ReceiverResponse struct {
	Id          int     `json:"id" example:"0"`
	Photo       *string `json:"photo" example:"photo"`
	Receiver    *string `json:"receiver" example:"name"`
	PhoneNumber *string `json:"phone_number" example:"08xx"`
}

// @name PaginationMetaResponse
type PaginationMetaResponse struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

// @name ReceiverListResponse
type ReceiverListResponse struct {
	Items []ReceiverResponse     `json:"items"`
	Meta  PaginationMetaResponse `json:"meta"`
}
