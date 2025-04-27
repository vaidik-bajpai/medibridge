package models

type SuccessResponse struct {
	Message string `json:message`
	Status int	`json:"status"`
	Data interface{}	`json:"data,omitempty"`
}

type FailureResponse struct {
	Status int `json:"status"`
	Error string `json:"error"`
}