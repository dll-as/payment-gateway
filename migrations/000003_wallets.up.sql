-- Wallets: User and master wallets (persistent and disposable)
CREATE TABLE IF NOT EXISTS wallets (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE, -- NULL for master/disposable wallets
    currency_id INTEGER NOT NULL REFERENCES currencies(id) ON DELETE RESTRICT,
    address VARCHAR(255) UNIQUE NOT NULL,
    encrypted_private_key TEXT, -- Encrypted private key (AES or KMS)
    trx_balance DECIMAL(20,8) DEFAULT 0, -- TRX balance for network fees
    is_master BOOLEAN DEFAULT FALSE, -- Master wallet for fees
    is_disposable BOOLEAN DEFAULT FALSE, -- Disposable wallet for invoices
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, currency_id, is_master, is_disposable),
    CONSTRAINT valid_trx_balance CHECK (trx_balance >= 0)
);
CREATE INDEX IF NOT EXISTS idx_wallets_user_id ON wallets(user_id);
CREATE INDEX IF NOT EXISTS idx_wallets_user_id_currency_id ON wallets(user_id, currency_id);
CREATE INDEX IF NOT EXISTS idx_wallets_address ON wallets(address);