package dto

// @name Response
type Response struct {
	Message string `json:"message" example:"Error Message"`
	Data    any    `json:"data,omitempty"`
	Success bool   `json:"isSuccess" example:"false"`
	Error   string `json:"error,omitempty" example:"internal server error"`
	Errors  any    `json:"errors,omitempty"`
}

// @name BaseResponse
type BaseResponse struct {
	Message string `json:"message"`
	Success bool   `json:"isSuccess"`
}

// @name ErrorResponse
type ErrorResponse struct {
	Message string `json:"message" example:"error message"`
	Success bool   `json:"isSuccess" example:"false"`
	Error   string `json:"error" example:"error"`
}

type ResponseSuccess struct {
	Status  string `json:"status" example:"success"`
	Message string `json:"message" example:"Welcome, John doe"`
}
