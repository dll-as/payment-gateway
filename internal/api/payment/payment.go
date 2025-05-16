package payment

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/rezatg/payment-gateway/internal/models"
	"github.com/rezatg/payment-gateway/internal/services/auth"
	"github.com/rezatg/payment-gateway/internal/services/blockchain"
	"github.com/rezatg/payment-gateway/pkg/errors"
	"github.com/rezatg/payment-gateway/pkg/logger"
)

// PaymentHandler handles payment-related requests
type PaymentHandler struct {
	blockchainFactory blockchain.Factory
	authService       auth.AuthService
	validator         *validator.Validate
}

// New creates a new PaymentHandler
func New(factory blockchain.Factory, authService auth.AuthService) *PaymentHandler {
	return &PaymentHandler{
		blockchainFactory: factory,
		authService:       authService,
		validator:         validator.New(),
	}
}

// CheckBalance checks the balance of a wallet
func (h *PaymentHandler) CheckBalance(c fiber.Ctx) error {
	var req models.BalanceRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "Invalid input",
			Code:    fiber.StatusBadRequest,
			Details: err.Error(),
		})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "Validation failed",
			Code:    fiber.StatusBadRequest,
			Details: err.Error(),
		})
	}

	svc, err := h.blockchainFactory.GetService(req.Currency)
	if err != nil {
		logger.Error("Failed to get blockchain service", err, "currency", req.Currency)
		return errors.HandleAPIError(c, err)
	}

	balance, err := svc.GetBalance(req.Address)
	if err != nil {
		logger.Error("Failed to check balance", err, "currency", req.Currency, "address", req.Address)
		return errors.HandleAPIError(c, err)
	}

	logger.Info("Balance checked", "currency", req.Currency, "address", req.Address, "balance", balance)
	return c.JSON(models.BalanceResponse{
		Balance:  balance,
		Currency: req.Currency,
	})
}

// GetAuthService provides auth validation for middleware
func (h *PaymentHandler) GetAuthService() auth.AuthService {
	return h.authService
}
