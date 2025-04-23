CREATE TABLE accounts (
    id SERIAL PRIMARY KEY,
    balance BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- some values for testing

INSERT INTO accounts (balance, created_at, updated_at) VALUES
  (15000, NOW(), NOW()),
  (7500, NOW(), NOW()),
  (30000, NOW(), NOW()),
  (0, NOW(), NOW()),
  (500, NOW(), NOW());