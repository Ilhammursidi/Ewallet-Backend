package dto

type ReceiverResponse struct {
	Id          int     `json:"id"`
	Photo       *string `json:"photo"`
	Receiver    *string `json:"receiver"`
	PhoneNumber *string `json:"phone_number"`
}

type PaginationMetaResponse struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

type ReceiverListResponse struct {
	Items []ReceiverResponse     `json:"items"`
	Meta  PaginationMetaResponse `json:"meta"`
}
