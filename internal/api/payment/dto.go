package payment

type (
	generateAddressRequest struct {
		Currency string `json:"currency" validate:"required,oneof=ETH TRX BTC" example:"ETH"`
	}

	generateAddressResponse struct {
		Address             string `json:"address"`
		EncryptedPrivateKey string `json:"encrypted_privateKey"`
		Currency            string `json:"currency" example:"ETH"`
	}
)
