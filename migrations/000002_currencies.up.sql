-- Currencies: Supported cryptocurrencies (e.g., USDT, BTC)
CREATE TABLE IF NOT EXISTS currencies (
    id SERIAL PRIMARY KEY,
    code VARCHAR(10) UNIQUE NOT NULL, -- e.g., USDT, TRX
    name VARCHAR(50) NOT NULL, -- e.g., Tether, Tron
    decimals INTEGER NOT NULL, -- e.g., 6 for USDT
    network VARCHAR(50) NOT NULL, -- e.g., TRON, ETHEREUM
    contract_address VARCHAR(255), -- e.g., USDT TRC-20 contract
    is_active BOOLEAN DEFAULT TRUE,
    CONSTRAINT valid_decimals CHECK (decimals >= 0)
);