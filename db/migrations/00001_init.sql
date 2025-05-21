-- +gooseUp
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE SCHEMA IF NOT EXISTS "public";
-- +goose StatementEnd

-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS maintainers(
  id SERIAL NOT NULL,
  ghUsername TEXT NOT NULL UNIQUE,
  full_name TEXT NOT NULL,

  CONSTRAINT "maintainers_pkey" PRIMARY KEY (id)
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_account(
  id SERIAL NOT NULL,
  email TEXT NOT NULL,
  ghId TEXT,
  ghUsername TEXT NOT NULL,
  status BOOLEAN DEFAULT true,
  bounty INT NOT NULL DEFAULT 0,
  refresh_token TEXT,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),

  CONSTRAINT "user_account_pkey" PRIMARY KEY (id)
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_onboarding(
  id SERIAL NOT NULL,
  email TEXT NOT NULL,
  ghUsername TEXT NOT NULL,
  otp TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT NOW(),
  expiry_at TIMESTAMP NOT NULL,

  CONSTRAINT "user_onboarding_pkey" PRIMARY KEY (id)
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS repository(
  id UUID NOT NULL,
  name TEXT NOT NULL,
  description TEXT NOT NULL,
  url TEXT NOT NULL,
  maintainers TEXT[],
  tags TEXT[],
  is_internal BOOLEAN DEFAULT false,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),

  CONSTRAINT "repository_pkey" PRIMARY KEY (id)
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS issues(
  id UUID NOT NULL,
  title TEXT NOT NULL,
  repoId UUID NOT NULL,
  url TEXT NOT NULL UNIQUE,
  resolved BOOLEAN DEFAULT false,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),

  CONSTRAINT "issues_pkey" PRIMARY KEY (id)
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS issue_claims(
  id SERIAL NOT NULL,
  ghUsername TEXT NOT NULL,
  issue_id UUID NOT NULL,
  claimed_on TIMESTAMP NOT NULL,
  elapsed_on TIMESTAMP NOT NULL,

  CONSTRAINT "issue_claims_pkey" PRIMARY KEY (id),
  CONSTRAINT "issue_claims_ghUsername_fkey"
    FOREIGN KEY (ghUsername)
      REFERENCES user_account(ghUsername)
        ON DELETE RESTRICT
        ON UPDATE CASCADE
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS bounty_log(
  id SERIAL NOT NULL,
  ghUsername TEXT NOT NULL,
  dispatchedBy TEXT NOT NULL,
  proof_url TEXT NOT NULL,
  repo_id UUID NOT NULL,
  amount INTEGER NOT NULL,
  created_at TIMESTAMP DEFAULT NOW(),

  CONSTRAINT "bounty_log_pkey" PRIMARY KEY (id),
  CONSTRAINT "bounty_log_dispatch_fkey"
    FOREIGN KEY (dispatched_by)
      REFERENCES maintainers(ghUsername)
        ON DELETE RESTRICT
        ON UPDATE CASCADE,
  CONSTRAINT "bounty_log_ghUsername_fkey"
    FOREIGN KEY (ghUsername)
      REFERENCES user_account(ghUsername)
        ON DELETE RESTRICT
        ON UPDATE CASCADE,
  CONSTRAINT "bounty_log_repo_id_fkey"
    FOREIGN KEY (repo_id)
      REFERENCES repository(id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS badge_info(
  id UUID NOT NULL,
  name TEXT NOT NULL,
  description TEXT NOT NULL,
  won_by INTEGER DEFAULT 0,

  CONSTRAINT "badge_info_pkey" PRIMARY KEY (id)
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS badge_dispatch(
  id SERIAL NOT NULL,
  ghUsername TEXT NOT NULL,
  badge_id UUID NOT NULL,
  awarded_at TIMESTAMP DEFAULT NOW(),

  CONSTRAINT "badge_dispatch_pkey" PRIMARY KEY (id),
  CONSTRAINT "badge_dispatch_uniqueness_validator" UNIQUE (ghUsername, badge_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP SCHEMA IF EXISTS public;
-- +goose StatementEnd
