package invoice

import (
	"github.com/gofiber/fiber/v3"
)

// RegisterRoutes registers invoice-related routes
func RegisterRoutes(router fiber.Router, handler *InvoiceHandler) {
	invoiceGroup := router.Group("/invoice")
	// invoiceGroup.Use(middleware.AuthMiddleware(handler.service.GetAuthService()))
	invoiceGroup.Post("/create", handler.CreateInvoice)
}
