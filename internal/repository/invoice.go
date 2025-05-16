package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/rezatg/payment-gateway/internal/models"
	"github.com/rezatg/payment-gateway/pkg/errors"
	"github.com/rezatg/payment-gateway/pkg/logger"
)

// InvoiceRepository handles invoice database operations
type InvoiceRepository interface {
	Create(ctx context.Context, invoice *models.Invoice) error
	FindByID(ctx context.Context, id string) (*models.Invoice, error)
}

type invoiceRepository struct {
	db *sql.DB
}

// NewInvoiceRepository creates a new InvoiceRepository
func NewInvoiceRepository(db *sql.DB) InvoiceRepository {
	return &invoiceRepository{db}
}

// Create creates a new invoice
func (r *invoiceRepository) Create(ctx context.Context, invoice *models.Invoice) error {
	invoice.ID = uuid.New().String()
	invoice.Status = "pending"
	invoice.CreatedAt = time.Now()

	query := `
		INSERT INTO invoices (id, user_id, amount, currency, note, status, payment_link, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.ExecContext(ctx, query,
		invoice.ID, invoice.UserID, invoice.Amount, invoice.Currency, invoice.Note,
		invoice.Status, invoice.PaymentLink, invoice.CreatedAt,
	)
	if err != nil {
		logger.Error("Failed to create invoice", err, "user_id", invoice.UserID)
		return errors.NewInternalServerError("Failed to create invoice", err.Error())
	}

	logger.Info("Invoice created", "invoice_id", invoice.ID)
	return nil
}

// FindByID finds an invoice by ID
func (r *invoiceRepository) FindByID(ctx context.Context, id string) (*models.Invoice, error) {
	var invoice models.Invoice
	query := `
		SELECT id, user_id, amount, currency, note, status, payment_link, created_at
		FROM invoices WHERE id = $1
	`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&invoice.ID, &invoice.UserID, &invoice.Amount, &invoice.Currency,
		&invoice.Note, &invoice.Status, &invoice.PaymentLink, &invoice.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		logger.Error("Failed to find invoice", err, "invoice_id", id)
		return nil, errors.NewInternalServerError("Failed to find invoice", err.Error())
	}

	return &invoice, nil
}
