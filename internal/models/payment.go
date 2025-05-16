package models

// BalanceRequest represents the request payload for checking balance
// @Description Request to check wallet balance
type BalanceRequest struct {
	Currency string `json:"currency" validate:"required,oneof=ETH TRX BTC" example:"ETH"`
	Address  string `json:"address" validate:"required" example:"0x1234567890abcdef1234567890abcdef12345678"`
}

// BalanceResponse represents the response payload for balance check
// @Description Response containing wallet balance
type BalanceResponse struct {
	Balance  float64 `json:"balance" example:"10.5"`
	Currency string  `json:"currency" example:"ETH"`
}
