package models

// SuccessResponse represents a standard success response.
// @Description Standard success response format with an optional data field.
type SuccessResponse struct {
	Message string `json:"message" example:"Operation successful"`
	Status  int    `json:"status" example:"200"`
	Data    any    `json:"data,omitempty"`
}

// FailureResponse represents a standard failure response.
// @Description Standard error response format with status and error message.
type FailureResponse struct {
	Status int    `json:"status" example:"400"`
	Error  string `json:"error" example:"Bad Request"`
}
