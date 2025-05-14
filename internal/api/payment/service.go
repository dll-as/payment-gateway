package payment

import "github.com/rezatg/payment-gateway/internal/services/blockchain"

type PaymentService interface {
	CheckBalance(currency, address string) (float64, error)
}

type paymentService struct{}

func NewPaymentService() PaymentService {
	return &paymentService{}
}

func (p *paymentService) CheckBalance(currency, address string) (float64, error) {
	svc, err := blockchain.GetBlockchainService(currency)
	if err != nil {
		return 0, err
	}

	return svc.GetBalance(address)
}
