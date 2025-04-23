CREATE TABLE transfers (
    id SERIAL PRIMARY KEY,
    from_account_id INTEGER NOT NULL REFERENCES accounts(id),
    to_account_id INTEGER NOT NULL REFERENCES accounts(id),
    amount BIGINT NOT NULL,
    status TEXT NOT NULL,
    external_id TEXT UNIQUE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- some values for testing

INSERT INTO transfers (from_account_id, to_account_id, amount, status, external_id, created_at, updated_at) VALUES
  (1, 2, 5000, 'pending', 'ext-001', NOW(),NOW()),
  (2, 3, 12000, 'completed', 'ext-002', NOW(),NOW()),
  (3, 1, 2000, 'failed', 'ext-003', NOW(),NOW()),
  (4, 5, 1000, 'pending', 'ext-004', NOW(),NOW()),
  (5, 2, 1500, 'completed', 'ext-005', NOW(),NOW());