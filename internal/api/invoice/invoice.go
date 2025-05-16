package invoice

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/rezatg/payment-gateway/internal/models"
	"github.com/rezatg/payment-gateway/internal/services/auth"
	"github.com/rezatg/payment-gateway/internal/services/invoice"
	"github.com/rezatg/payment-gateway/pkg/errors"
	"github.com/rezatg/payment-gateway/pkg/logger"
)

// InvoiceHandler handles invoice-related requests
type InvoiceHandler struct {
	service   invoice.InvoiceService
	validator *validator.Validate
}

// New creates a new InvoiceHandler
func New(service invoice.InvoiceService) *InvoiceHandler {
	return &InvoiceHandler{
		service:   service,
		validator: validator.New(),
	}
}

// CreateInvoice creates a new invoice
func (h *InvoiceHandler) CreateInvoice(c fiber.Ctx) error {
	var req models.CreateInvoiceRequest
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

	userID := c.Locals("userID").(string)
	inv, err := h.service.CreateInvoice(c.Context(), userID, req.Amount, req.Currency, req.Note)
	if err != nil {
		logger.Error("Failed to create invoice", err, "user_id", userID)
		return errors.HandleAPIError(c, err)
	}

	logger.Info("Invoice created", "invoice_id", inv.ID, "user_id", userID)
	return c.Status(fiber.StatusCreated).JSON(inv)
}

// GetAuthService provides auth validation for middleware
func (h *InvoiceHandler) GetAuthService() auth.AuthService {
	return h.service.GetAuthService()
}
