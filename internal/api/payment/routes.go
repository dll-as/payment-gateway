package payment

import "github.com/gofiber/fiber/v3"

func RegisterRoutes(router fiber.Router) {
	service := NewPaymentService()
	handler := NewPaymentHandler(service)

	router.Post("/check-balance", handler.CheckBalance)
}
