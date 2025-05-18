-- Enable UUID extension for generating unique invoice codes
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Invoices: Payment requests with QR codes and disposable wallets
CREATE TABLE IF NOT EXISTS invoices (
    id SERIAL PRIMARY KEY,
    invoice_code UUID UNIQUE DEFAULT uuid_generate_v4(), -- Unique invoice code
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    wallet_id INTEGER NOT NULL REFERENCES wallets(id) ON DELETE RESTRICT, -- Disposable wallet
    currency_id INTEGER NOT NULL REFERENCES currencies(id) ON DELETE RESTRICT,
    amount DECIMAL(20,8) NOT NULL,
    fee_payer VARCHAR(20) NOT NULL DEFAULT 'client', -- client, acceptor
    status VARCHAR(20) NOT NULL, -- pending, paid, expired, failed
    qr_code_url VARCHAR(255), -- URL to QR code
    pdf_url VARCHAR(255), -- URL to PDF invoice
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL, -- Expiration time
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT valid_amount CHECK (amount > 0),
    CONSTRAINT valid_status CHECK (status IN ('pending', 'paid', 'expired', 'failed'))
);
CREATE INDEX IF NOT EXISTS idx_invoices_user_id_status ON invoices(user_id, status);
CREATE INDEX IF NOT EXISTS idx_invoices_invoice_code ON invoices(invoice_code);
CREATE INDEX IF NOT EXISTS idx_invoices_expires_at ON invoices(expires_at);

-- Payment Settings: Per-user payment configurations
CREATE TABLE IF NOT EXISTS payment_settings (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    currency_id INTEGER NOT NULL REFERENCES currencies(id) ON DELETE RESTRICT,
    fee_percentage DECIMAL(5,2) NOT NULL DEFAULT 2.00, -- e.g., 2% fee
    callback_url VARCHAR(255), -- URL for payment confirmation
    success_url VARCHAR(255), -- Redirect after successful payment
    cancel_url VARCHAR(255), -- Redirect after cancellation
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, currency_id),
    CONSTRAINT valid_fee_percentage CHECK (fee_percentage >= 0 AND fee_percentage <= 100)
);
CREATE INDEX IF NOT EXISTS idx_payment_settings_user_id ON payment_settings(user_id);