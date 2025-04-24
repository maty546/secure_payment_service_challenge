CREATE TABLE accounts (
    id TEXT PRIMARY KEY,
    balance BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

--some values for testing

INSERT INTO accounts (id,balance, created_at, updated_at) VALUES
  ('user1',15000, NOW(), NOW()),
  ('user2',7500, NOW(), NOW()),
  ('user3',30000, NOW(), NOW());