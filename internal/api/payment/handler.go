package payment

import (
	"github.com/gofiber/fiber/v3"
)

type PaymentHandler struct {
	service PaymentService
}

func NewPaymentHandler(service PaymentService) *PaymentHandler {
	return &PaymentHandler{service}
}

func (h *PaymentHandler) CheckBalance(c fiber.Ctx) error {
	var req balanceRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	balance, err := h.service.CheckBalance(req.Currency, req.Address)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"balance":  balance,
		"currency": req.Currency,
	})
}
