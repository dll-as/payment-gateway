package models

// ErrorResponse represents a standard error response
// @Description Standard error response for API
type ErrorResponse struct {
	Error   string `json:"error" example:"Invalid credentials"`
	Code    int    `json:"code" example:"401"`
	Details string `json:"details,omitempty" example:"Email or password is incorrect"`
}
