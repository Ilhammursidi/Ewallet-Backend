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
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
	Total_Data int    `json:"total_data"`
	Total_Page int    `json:"total_page"`
	Next_Page  string `json:"next_page"`
	Prev_Page  string `json:"prev_page"`
}

// @name ReceiverListResponse
type ReceiverListResponse struct {
	Items []ReceiverResponse     `json:"items"`
	Meta  PaginationMetaResponse `json:"meta"`
}

type TransactionStatus string
type TransactionType string
type CashflowType string

const (
	StatusPending TransactionStatus = "PENDING"
	StatusSuccess TransactionStatus = "SUCCESS"
	StatusFailed  TransactionStatus = "FAILED"

	TypeTopUp       TransactionType = "TOPUP"
	TypeTransferIn  TransactionType = "TRANSFER_IN"
	TypeTransferOut TransactionType = "TRANSFER_OUT"

	FlowIncome  CashflowType = "income"
	FlowExpense CashflowType = "expense"
)

type CreateTransactionParams struct {
	UserID           int `json:"user_id"`
	ReceiverWalletID int `json:"receiver_wallet_id"`
	PaymentMethodID  int `json:"payment_method_id"`
	Amount           int `json:"amount"`
}

type CreateTopUpDetailParams struct {
	TransactionID   int `json:"transaction_id"`
	WalletID        int `json:"wallet_id"`
	PaymentMethodID int `json:"payment_method_id"`
	OrderAmount     int `json:"order_amount"`
	TaxAmount       int `json:"tax_amount"`
	DeliveryFee     int `json:"delivery_fee"`
	TotalAmount     int `json:"total_amount"`
}

type TopUpHTTPRequest struct {
	PaymentMethodID int `json:"payment_method_id"   validate:"required"`
	OrderAmount     int `json:"order_amount"        validate:"required,min=1"`
	TaxAmount       int `json:"tax_amount"          validate:"min=0"`
	DeliveryFee     int `json:"delivery_fee"        validate:"min=0"`
}

type TopUpServiceRequest struct {
	UserID          int `json:"user_id"`
	PaymentMethodID int `json:"payment_method_id"`
	OrderAmount     int `json:"order_amount"`
	TaxAmount       int `json:"tax_amount"`
	DeliveryFee     int `json:"delivery_fee"`
}

type TopUpResponse struct {
	TransactionID int               `json:"transaction_id"`
	TopUpDetailID int               `json:"topup_detail_id"`
	TotalAmount   int               `json:"total_amount"`
	CreditAmount  int               `json:"credit_amount"`
	Status        TransactionStatus `json:"status"`
}

type CreateTransferParams struct {
	UserID           int `json:"user_id"`
	SenderWalletID   int `json:"sender_wallet_id"`
	ReceiverWalletID int `json:"receiver_wallet_id"`
	Amount           int `json:"amount"`
}

type TransferResponse struct {
	TransactionID int               `json:"transaction_id"`
	Amount        int               `json:"amount"`
	Status        TransactionStatus `json:"status"`
}

type TransferHTTPRequest struct {
	ReceiverID int    `json:"receiver_id" validate:"required"`
	Amount     int    `json:"amount"      validate:"required,min=1"`
	Pin        string `json:"pin" validate:"required,len=6"`
}

type TransferServiceRequest struct {
	UserID     int `json:"user_id"`
	ReceiverID int `json:"receiver_id"`
	Amount     int `json:"amount"`
}
