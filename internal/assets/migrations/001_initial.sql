-- +migrate Up

CREATE TABLE customers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(100) UNIQUE NOT NULL,
    password_hash CHAR(60) NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    balance INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE customers_accounts (
    customer_fkey UUID REFERENCES customers(id) ON DELETE RESTRICT,
    account_fkey UUID REFERENCES accounts(id),
    PRIMARY KEY (customer_fkey, account_fkey)
);

CREATE TABLE deposits (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    recepient_fkey UUID REFERENCES accounts(id),
    amount INTEGER NOT NULL,
    atm_signature CHAR(64) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE transfers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sender_fkey UUID REFERENCES accounts(id),
    recepient_fkey UUID REFERENCES accounts(id),
    amount INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE withdrawals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sender_fkey UUID REFERENCES accounts(id),
    amount INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_customers_email ON customers(email);
CREATE INDEX idx_customers_username ON customers(username);

CREATE INDEX idx_transfers_sender ON transfers(sender_fkey);
CREATE INDEX idx_transfers_recipient ON transfers(recepient_fkey);
CREATE INDEX idx_deposits_recipient ON deposits(recepient_fkey);
CREATE INDEX idx_withdrawals_sender ON withdrawals(sender_fkey);

-- +migrate Down

DROP TABLE IF EXISTS withdrawals;
DROP TABLE IF EXISTS transfers;
DROP TABLE IF EXISTS deposits;
DROP TABLE IF EXISTS customers_accounts;
DROP TABLE IF EXISTS accounts;
DROP TABLE IF EXISTS customers;

DROP INDEX IF EXISTS idx_customers_email;
DROP INDEX IF EXISTS idx_customers_username;
DROP INDEX IF EXISTS idx_transfers_sender;
DROP INDEX IF EXISTS idx_transfers_recipient;
DROP INDEX IF EXISTS idx_deposits_recipient;
DROP INDEX IF EXISTS idx_withdrawals_sender;
