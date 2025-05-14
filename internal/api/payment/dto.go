package payment

type balanceRequest struct {
	Currency string `json:"currency"`
	Address  string `json:"address"`
}
