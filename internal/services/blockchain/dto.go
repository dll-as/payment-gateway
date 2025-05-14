package blockchain

type TRXAccount struct {
	Balance int64 `json:"balance"`
}

type TronResponse struct {
	Data []TRXAccount `json:"data"`
}
