-- Transactions: Records of deposits, withdrawals, and fee transfers
CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    invoice_id INTEGER REFERENCES invoices(id) ON DELETE SET NULL, -- Link to invoice
    user_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
    wallet_id INTEGER REFERENCES wallets(id) ON DELETE SET NULL, -- Source wallet
    currency_id INTEGER NOT NULL REFERENCES currencies(id) ON DELETE RESTRICT,
    amount DECIMAL(20,8) NOT NULL, -- Transaction amount
    fee DECIMAL(20,8) NOT NULL DEFAULT 0, -- Fee (e.g., 2%)
    tx_hash VARCHAR(255), -- Blockchain transaction hash
    fee_tx_id VARCHAR(255), -- Fee transaction ID (to master wallet)
    client_tx_id VARCHAR(255), -- Client transaction ID (to client wallet)
    status VARCHAR(20) NOT NULL, -- pending, confirmed, failed
    type VARCHAR(20) NOT NULL, -- deposit, withdrawal, fee_transfer, client_transfer
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT valid_amount CHECK (amount > 0),
    CONSTRAINT valid_fee CHECK (fee >= 0),
    CONSTRAINT valid_status CHECK (status IN ('pending', 'confirmed', 'failed')),
    CONSTRAINT valid_type CHECK (type IN ('deposit', 'withdrawal', 'fee_transfer', 'client_transfer'))
);
CREATE INDEX IF NOT EXISTS idx_transactions_user_id ON transactions(user_id);
CREATE INDEX IF NOT EXISTS idx_transactions_user_id_status ON transactions(user_id, status);
CREATE INDEX IF NOT EXISTS idx_transactions_status ON transactions(status);
CREATE INDEX IF NOT EXISTS idx_transactions_invoice_id ON transactions(invoice_id);
CREATE INDEX IF NOT EXISTS idx_transactions_tx_hash ON transactions(tx_hash);