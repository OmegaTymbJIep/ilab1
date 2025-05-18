-- +migrate Up
CREATE TYPE audit_action_enum AS ENUM (
    'login_success',
    'account_created',
    'account_deleted',
    'deposit_made',
    'withdrawal_made',
    'transfer_made',
    'excel_report_generated'
);

CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id UUID REFERENCES customers(id),
    account_id UUID REFERENCES accounts(id),
    action audit_action_enum NOT NULL,
    details JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_audit_logs_customer_id ON audit_logs(customer_id);
CREATE INDEX idx_audit_logs_account_id ON audit_logs(account_id);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at);

-- +migrate Down
DROP INDEX IF EXISTS idx_audit_logs_created_at;
DROP INDEX IF EXISTS idx_audit_logs_action;
DROP INDEX IF EXISTS idx_audit_logs_account_id;
DROP INDEX IF EXISTS idx_audit_logs_customer_id;

DROP TABLE IF EXISTS audit_logs;

DROP TYPE IF EXISTS audit_action_enum;