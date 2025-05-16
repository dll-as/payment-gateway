package payment

// import (
// 	"github.com/rezatg/payment-gateway/internal/services/auth"
// 	"github.com/rezatg/payment-gateway/internal/services/blockchain"
// )

// // PaymentService defines the interface for payment operations
// type PaymentService interface {
// 	CheckBalance(currency, address string) (float64, error)
// 	GetAuthService() auth.AuthService
// }

// // paymentService implements PaymentService
// type paymentService struct {
// 	blockchainFactory blockchain.Factory
// 	authService       auth.AuthService
// }

// // NewPaymentService creates a new payment service
// func NewPaymentService(factory blockchain.Factory, authService auth.AuthService) PaymentService {
// 	return &paymentService{
// 		blockchainFactory: factory,
// 		authService:       authService,
// 	}
// }

// // CheckBalance checks the balance for a given address
// func (p *paymentService) CheckBalance(currency, address string) (float64, error) {
// 	svc, err := p.blockchainFactory.GetService(currency)
// 	if err != nil {
// 		return 0, err
// 	}
// 	return svc.GetBalance(address)
// }

// // GetAuthService returns the auth service for middleware usage
// func (p *paymentService) GetAuthService() auth.AuthService {
// 	return p.authService
// }
