-- +migrate Up
ALTER TABLE accounts
    ADD COLUMN is_deleted BOOLEAN NOT NULL DEFAULT FALSE;

CREATE INDEX idx_accounts_is_deleted ON accounts(is_deleted);

-- +migrate Down
DROP INDEX IF EXISTS idx_accounts_is_deleted;
ALTER TABLE accounts DROP COLUMN IF EXISTS is_deleted;