package models

import "time"

// Invoice represents an invoice in the system
type Invoice struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Amount      float64   `json:"amount"`
	Currency    string    `json:"currency"`
	Note        string    `json:"note"`
	Status      string    `json:"status"`
	PaymentLink string    `json:"payment_link"`
	CreatedAt   time.Time `json:"created_at"`
}

// CreateInvoiceRequest represents the request payload for creating an invoice
type CreateInvoiceRequest struct {
	Amount   float64 `json:"amount" validate:"required,gt=0" example:"100.50"`
	Currency string  `json:"currency" validate:"required,oneof=ETH TRX BTC USD" example:"ETH"`
	Note     string  `json:"note,omitempty" example:"Payment for services"`
}

// InvoiceResponse represents the response payload for an invoice
type InvoiceResponse struct {
	ID          string  `json:"id" example:"inv_123456789"`
	Amount      float64 `json:"amount" example:"100.50"`
	Currency    string  `json:"currency" example:"ETH"`
	Status      string  `json:"status" example:"pending"`
	PaymentLink string  `json:"payment_link" example:"https://gateway.example.com/pay/inv_123456789"`
	CreatedAt   string  `json:"created_at" example:"2025-05-16T10:30:00Z"`
}
