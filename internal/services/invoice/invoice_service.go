package invoice

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rezatg/payment-gateway/config"
	"github.com/rezatg/payment-gateway/internal/models"
	"github.com/rezatg/payment-gateway/internal/repository"
	"github.com/rezatg/payment-gateway/internal/services/auth"
	"github.com/rezatg/payment-gateway/pkg/errors"
	"github.com/rezatg/payment-gateway/pkg/logger"
)

// InvoiceService defines the interface for invoice operations
type InvoiceService interface {
	CreateInvoice(ctx context.Context, userID string, amount float64, currency, note string) (*models.InvoiceResponse, error)
	GetAuthService() auth.AuthService
}

// invoiceService implements InvoiceService
type invoiceService struct {
	repo        repository.InvoiceRepository
	authService auth.AuthService
}

// New creates a new InvoiceService
func New(repo repository.InvoiceRepository, authService auth.AuthService) InvoiceService {
	return &invoiceService{
		repo:        repo,
		authService: authService,
	}
}

// CreateInvoice creates a new invoice
func (s *invoiceService) CreateInvoice(ctx context.Context, userID string, amount float64, currency, note string) (*models.InvoiceResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Validate inputs
	if amount <= 0 {
		return nil, errors.NewBadRequestError("Amount must be greater than zero")
	}
	if !isValidCurrency(currency) {
		return nil, errors.NewBadRequestError("Invalid currency")
	}

	// Generate payment link
	paymentLink := fmt.Sprintf("%s/pay/%s", config.GetEnv("BASE_URL", "https://gateway.example.com"), uuid.New().String())

	// Create invoice
	newInvoice := &models.Invoice{
		UserID:      userID,
		Amount:      amount,
		Currency:    currency,
		Note:        note,
		PaymentLink: paymentLink,
	}

	if err := s.repo.Create(ctx, newInvoice); err != nil {
		logger.Error("Failed to create invoice", err, "user_id", userID)
		return nil, err
	}

	logger.Info("Invoice created", "invoice_id", newInvoice.ID, "user_id", userID)
	return &models.InvoiceResponse{
		ID:          newInvoice.ID,
		Amount:      newInvoice.Amount,
		Currency:    newInvoice.Currency,
		Status:      newInvoice.Status,
		PaymentLink: newInvoice.PaymentLink,
		CreatedAt:   newInvoice.CreatedAt.Format(time.RFC3339),
	}, nil
}

// GetAuthService returns the auth service for middleware usage
func (s *invoiceService) GetAuthService() auth.AuthService {
	return s.authService
}

// isValidCurrency checks if the currency is supported
func isValidCurrency(currency string) bool {
	supportedCurrencies := []string{"ETH", "TRX", "BTC", "USD"}
	for _, c := range supportedCurrencies {
		if c == currency {
			return true
		}
	}
	return false
}
