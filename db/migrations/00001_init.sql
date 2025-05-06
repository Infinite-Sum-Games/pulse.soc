-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_account(
  id UUID NOT NULL,
  email TEXT NOT NULL,
  first_name TEXT NULL,
  last_name TEXT NULL,
  bounty INT NOT NULL DEFAULT 0,
  refresh_token TEXT NULL,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),

  CONSTRAINT "user_account_pkey" PRIMARY KEY (id)
);
CREATE UNIQUE INDEX idx_user_account__email ON user_account(email);
-- +goose StatementEnd
