-- Modify "wallets" table
ALTER TABLE "wallets" ALTER COLUMN "available_balance" SET DEFAULT 0, ALTER COLUMN "ledger_balance" SET DEFAULT 0, ALTER COLUMN "holding_balance" SET DEFAULT 0;
