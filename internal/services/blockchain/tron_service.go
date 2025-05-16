package blockchain

import (
	"fmt"

	"github.com/bytedance/sonic"
	"github.com/valyala/fasthttp"
)

const tronGridAPI = "https://api.trongrid.io/v1/accounts"

type TronService struct{}

func NewTronService() *TronService {
	return &TronService{}
}

func (t *TronService) GenerateAddress() (string, error) {
	return "TXabc123...", nil
}

func (t *TronService) GetBalance(address string) (float64, error) {
	statusCode, body, err := fasthttp.Get(nil, fmt.Sprintf("%s/%s", tronGridAPI, address))
	if err != nil {
		return 0, err
	}

	if statusCode != fasthttp.StatusOK {
		return 0, fmt.Errorf("error retrieving account information, HTTP status: %d", statusCode)
	}

	var response TronResponse
	if err := sonic.Unmarshal(body, &response); err != nil || len(response.Data) == 0 {
		return 0, err
	}

	return float64(response.Data[0].Balance) / 1_000_000, nil
}

func (t *TronService) SendTransaction(fromPrivKey, toAddress string, amount float64) (string, error) {
	return "trx_hash_123", nil
}

func (t *TronService) CheckTxConfirmed(txHash string) (bool, error) {
	// call TronGrid to check confirmations
	return true, nil
}
