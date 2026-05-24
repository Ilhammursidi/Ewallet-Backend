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

// @name TopupRequest
type TopupRequest struct {
	SubTotal        uint `json:"sub_total" binding:"required"`
	PaymentMethodId int  `json:"payment_method_id" binding:"required"`
	Order           int  `json:""`
	Delivery        int
	Tax             int
}
