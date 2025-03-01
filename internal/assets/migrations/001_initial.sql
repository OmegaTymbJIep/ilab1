-- +migrate Up

CREATE TYPE transaction_type_enum AS ENUM (
    'deposit',
    'withdrawal',
    'transfer'
);

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
    balance BIGINT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE customers_accounts (
    customer_fkey UUID REFERENCES customers(id) ON DELETE RESTRICT,
    account_fkey UUID REFERENCES accounts(id),
    PRIMARY KEY (customer_fkey, account_fkey)
);

CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    type INTEGER,
    sender_fkey UUID REFERENCES accounts(id),
    recipient_fkey UUID REFERENCES accounts(id),
    amount BIGINT NOT NULL,
    atm_signature CHAR(64) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CHECK (
        (type = 0
            AND sender_fkey IS NULL
            AND recipient_fkey IS NOT NULL)
        OR
        (type = 1
            AND sender_fkey IS NOT NULL
            AND recipient_fkey IS NULL)
        OR
        (type = 2
            AND sender_fkey IS NOT NULL
            AND recipient_fkey IS NOT NULL)
    )
);

CREATE UNIQUE INDEX unique_atm_signature_for_deposits ON transactions (atm_signature)
    WHERE type = 0;

CREATE INDEX idx_customers_email ON customers(email);
CREATE INDEX idx_customers_username ON customers(username);

CREATE INDEX idx_transactions_sender ON transactions(sender_fkey);
CREATE INDEX idx_transactions_recipient ON transactions(recipient_fkey);

-- Only index rows where recipient_fkey is relevant
CREATE INDEX idx_transactions_recipient_created_at_deposit_transfer
    ON transactions (recipient_fkey, created_at DESC)
    WHERE type IN (0, 2);

-- Only index rows where sender_fkey is relevant
CREATE INDEX idx_transactions_sender_created_at_withdrawal_transfer
    ON transactions (sender_fkey, created_at DESC)
    WHERE type IN (1, 2);

-- +migrate Down

DROP INDEX IF EXISTS unique_atm_signature_for_deposits;

DROP INDEX IF EXISTS idx_customers_email;
DROP INDEX IF EXISTS idx_customers_username;

DROP INDEX IF EXISTS idx_transactions_sender;
DROP INDEX IF EXISTS idx_transactions_recipient;

DROP INDEX IF EXISTS idx_transactions_recipient_created_at_deposit_transfer;
DROP INDEX IF EXISTS idx_transactions_sender_created_at_withdrawal_transfer;

DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS customers_accounts;
DROP TABLE IF EXISTS accounts;
DROP TABLE IF EXISTS customers;

DROP TYPE IF EXISTS transaction_type_enum;