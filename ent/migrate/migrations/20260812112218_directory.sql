-- Rename a column from "multipler" to "multiplier"
ALTER TABLE "currencies" RENAME COLUMN "multipler" TO "multiplier";
-- Create "ledger_owners" table
CREATE TABLE "ledger_owners" (
  "id" character varying NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "deleted_at" timestamptz NULL,
  "type" character varying NOT NULL,
  "status" character varying NOT NULL DEFAULT 'ACTIVE',
  "business_owner" character varying NULL,
  "user_owner" character varying NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "ledger_owners_businesses_owner" FOREIGN KEY ("business_owner") REFERENCES "businesses" ("id") ON UPDATE NO ACTION ON DELETE SET NULL,
  CONSTRAINT "ledger_owners_users_owner" FOREIGN KEY ("user_owner") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE SET NULL
);
-- Create index "ledger_owners_business_owner_key" to table: "ledger_owners"
CREATE UNIQUE INDEX "ledger_owners_business_owner_key" ON "ledger_owners" ("business_owner");
-- Create index "ledger_owners_user_owner_key" to table: "ledger_owners"
CREATE UNIQUE INDEX "ledger_owners_user_owner_key" ON "ledger_owners" ("user_owner");
-- Create "vaults" table
CREATE TABLE "vaults" (
  "id" character varying NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "deleted_at" timestamptz NULL,
  "type" character varying NOT NULL DEFAULT 'TREASURY',
  "status" character varying NOT NULL DEFAULT 'ACTIVE',
  "metadata" jsonb NULL,
  "ledger_owner_vaults" character varying NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "vaults_ledger_owners_vaults" FOREIGN KEY ("ledger_owner_vaults") REFERENCES "ledger_owners" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Modify "wallets" table
ALTER TABLE "wallets" ALTER COLUMN "type" SET DEFAULT 'FIAT', DROP COLUMN "owner", DROP COLUMN "identifier", ADD COLUMN "vault_wallets" character varying NOT NULL, ADD CONSTRAINT "wallets_vaults_wallets" FOREIGN KEY ("vault_wallets") REFERENCES "vaults" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Create index "wallet_currency_id_vault_wallets" to table: "wallets"
CREATE UNIQUE INDEX "wallet_currency_id_vault_wallets" ON "wallets" ("currency_id", "vault_wallets");
-- Create "nuban_static_accounts" table
CREATE TABLE "nuban_static_accounts" (
  "id" character varying NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "deleted_at" timestamptz NULL,
  "provider" character varying NOT NULL,
  "provider_reference" character varying NOT NULL,
  "account_number" character varying NOT NULL,
  "account_name" character varying NOT NULL,
  "bank_name" character varying NOT NULL,
  "bank_code" character varying NOT NULL,
  "address" character varying NULL,
  "state" character varying NOT NULL DEFAULT 'OPEN',
  "wallet_nuban_static_account" character varying NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "nuban_static_accounts_wallets_nuban_static_account" FOREIGN KEY ("wallet_nuban_static_account") REFERENCES "wallets" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
