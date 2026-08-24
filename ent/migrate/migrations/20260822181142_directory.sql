-- Modify "nuban_static_accounts" table
ALTER TABLE "nuban_static_accounts" DROP COLUMN "address";
-- Create index "wallet_type" to table: "wallets"
CREATE INDEX "wallet_type" ON "wallets" ("type");
-- Create "transactions" table
CREATE TABLE "transactions" (
  "id" character varying NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "deleted_at" timestamptz NULL,
  "reference" character varying NOT NULL,
  "external_reference" character varying NOT NULL,
  "type" character varying NOT NULL,
  "channel" character varying NOT NULL,
  "status" character varying NOT NULL DEFAULT 'PENDING',
  "currency" character varying NOT NULL,
  "summary" character varying NULL,
  "provider" character varying NULL,
  "amount" bigint NOT NULL DEFAULT 0,
  "fee" bigint NOT NULL DEFAULT 0,
  "transaction_wallet" character varying NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "transactions_wallets_wallet" FOREIGN KEY ("transaction_wallet") REFERENCES "wallets" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "transactions_external_reference_key" to table: "transactions"
CREATE UNIQUE INDEX "transactions_external_reference_key" ON "transactions" ("external_reference");
-- Create index "transactions_reference_key" to table: "transactions"
CREATE UNIQUE INDEX "transactions_reference_key" ON "transactions" ("reference");
-- Create "ledgers" table
CREATE TABLE "ledgers" (
  "id" character varying NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "deleted_at" timestamptz NULL,
  "debit" bigint NOT NULL DEFAULT 0,
  "credit" bigint NOT NULL DEFAULT 0,
  "current_balance" bigint NOT NULL DEFAULT 0,
  "previous_balance" bigint NOT NULL DEFAULT 0,
  "is_reversal" boolean NOT NULL DEFAULT false,
  "transaction_id" character varying NOT NULL,
  "reversal_transaction_id" character varying NULL,
  "wallet_id" character varying NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "ledgers_transactions_ledger_entries" FOREIGN KEY ("transaction_id") REFERENCES "transactions" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "ledgers_transactions_reversal_ledger_entries" FOREIGN KEY ("reversal_transaction_id") REFERENCES "transactions" ("id") ON UPDATE NO ACTION ON DELETE SET NULL,
  CONSTRAINT "ledgers_wallets_ledger_entries" FOREIGN KEY ("wallet_id") REFERENCES "wallets" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create "nuban_deposits" table
CREATE TABLE "nuban_deposits" (
  "id" character varying NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "deleted_at" timestamptz NULL,
  "recipient_account_name" character varying NOT NULL,
  "recipient_account_number" character varying NOT NULL,
  "sender_account_name" character varying NOT NULL,
  "sender_account_number" character varying NOT NULL,
  "sender_bank_name" character varying NOT NULL,
  "sender_bank_code" character varying NOT NULL,
  "narration" character varying NOT NULL,
  "amount" bigint NOT NULL DEFAULT 0,
  "session_id" character varying NOT NULL,
  "provider_reference" character varying NOT NULL,
  "provider" character varying NOT NULL,
  "type" character varying NOT NULL DEFAULT 'STATIC',
  "transaction_id" character varying NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "nuban_deposits_transactions_nuban_deposit" FOREIGN KEY ("transaction_id") REFERENCES "transactions" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "nuban_deposits_transaction_id_key" to table: "nuban_deposits"
CREATE UNIQUE INDEX "nuban_deposits_transaction_id_key" ON "nuban_deposits" ("transaction_id");
-- Create "nuban_dynamic_accounts" table
CREATE TABLE "nuban_dynamic_accounts" (
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
  "expiration" timestamptz NOT NULL,
  "state" character varying NOT NULL DEFAULT 'OPEN',
  "wallet_nuban_dynamic_account" character varying NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "nuban_dynamic_accounts_wallets_nuban_dynamic_account" FOREIGN KEY ("wallet_nuban_dynamic_account") REFERENCES "wallets" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create "nuban_transfers" table
CREATE TABLE "nuban_transfers" (
  "id" character varying NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "deleted_at" timestamptz NULL,
  "recipient_account_name" character varying NOT NULL,
  "recipient_account_number" character varying NOT NULL,
  "sender_account_name" character varying NOT NULL,
  "sender_account_number" character varying NOT NULL,
  "sender_bank_name" character varying NOT NULL,
  "sender_bank_code" character varying NOT NULL,
  "amount" bigint NOT NULL DEFAULT 0,
  "session_id" character varying NOT NULL,
  "provider_reference" character varying NOT NULL,
  "provider" character varying NOT NULL,
  "transaction_nuban_transfer" character varying NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "nuban_transfers_transactions_nuban_transfer" FOREIGN KEY ("transaction_nuban_transfer") REFERENCES "transactions" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "nuban_transfers_transaction_nuban_transfer_key" to table: "nuban_transfers"
CREATE UNIQUE INDEX "nuban_transfers_transaction_nuban_transfer_key" ON "nuban_transfers" ("transaction_nuban_transfer");
-- Modify "currencies" table
ALTER TABLE "currencies" ADD COLUMN "logo" character varying NULL;
-- Create "stablecoin_networks" table
CREATE TABLE "stablecoin_networks" (
  "id" character varying NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "deleted_at" timestamptz NULL,
  "name" character varying NOT NULL,
  "slug" character varying NOT NULL,
  "logo_url" character varying NULL,
  "active" boolean NOT NULL DEFAULT true,
  PRIMARY KEY ("id")
);
-- Create "stablecoin_supported_networks" table
CREATE TABLE "stablecoin_supported_networks" (
  "id" character varying NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "deleted_at" timestamptz NULL,
  "can_send" boolean NOT NULL DEFAULT false,
  "can_receive" boolean NOT NULL DEFAULT true,
  "active" boolean NOT NULL DEFAULT true,
  "coin_id" character varying NOT NULL,
  "network_id" character varying NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "stablecoin_supported_networks__a83c556fee24394d04793f0a02423e85" FOREIGN KEY ("network_id") REFERENCES "stablecoin_networks" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "stablecoin_supported_networks__cdf2fca02e321c8d4043699501362c95" FOREIGN KEY ("coin_id") REFERENCES "currencies" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "stablecoinsupportednetwork_network_id_coin_id" to table: "stablecoin_supported_networks"
CREATE UNIQUE INDEX "stablecoinsupportednetwork_network_id_coin_id" ON "stablecoin_supported_networks" ("network_id", "coin_id");
-- Create "stablecoin_wallets" table
CREATE TABLE "stablecoin_wallets" (
  "id" character varying NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "deleted_at" timestamptz NULL,
  "coin" character varying NOT NULL,
  "network" character varying NOT NULL,
  "address" character varying NOT NULL,
  "provider_reference" character varying NOT NULL,
  "provider" character varying NOT NULL,
  "coin_id" character varying NOT NULL,
  "network_id" character varying NOT NULL,
  "stablecoin_supported_network_stablecoin_wallets" character varying NULL,
  "wallet_stablecoin_wallets" character varying NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "stablecoin_wallets_currencies_stablecoin_currencies" FOREIGN KEY ("coin_id") REFERENCES "currencies" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "stablecoin_wallets_stablecoin__e44924c2a6e0240f695a39c5f43834e9" FOREIGN KEY ("stablecoin_supported_network_stablecoin_wallets") REFERENCES "stablecoin_supported_networks" ("id") ON UPDATE NO ACTION ON DELETE SET NULL,
  CONSTRAINT "stablecoin_wallets_stablecoin_networks_stablecoin_networks" FOREIGN KEY ("network_id") REFERENCES "stablecoin_networks" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "stablecoin_wallets_wallets_stablecoin_wallets" FOREIGN KEY ("wallet_stablecoin_wallets") REFERENCES "wallets" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create "stablecoin_deposits" table
CREATE TABLE "stablecoin_deposits" (
  "id" character varying NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "deleted_at" timestamptz NULL,
  "coin" character varying NOT NULL,
  "network" character varying NOT NULL,
  "address" character varying NOT NULL,
  "provider_reference" character varying NOT NULL,
  "provider" character varying NOT NULL,
  "stablecoin_wallet_id" character varying NOT NULL,
  "transaction_id" character varying NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "stablecoin_deposits_stablecoin_wallets_stablecoin_deposits" FOREIGN KEY ("stablecoin_wallet_id") REFERENCES "stablecoin_wallets" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "stablecoin_deposits_transactions_stablecoin_deposit" FOREIGN KEY ("transaction_id") REFERENCES "transactions" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "stablecoin_deposits_transaction_id_key" to table: "stablecoin_deposits"
CREATE UNIQUE INDEX "stablecoin_deposits_transaction_id_key" ON "stablecoin_deposits" ("transaction_id");
-- Create "stablecoin_withdrawals" table
CREATE TABLE "stablecoin_withdrawals" (
  "id" character varying NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "deleted_at" timestamptz NULL,
  "coin" character varying NOT NULL,
  "network" character varying NOT NULL,
  "destination_address" character varying NOT NULL,
  "amount" bigint NOT NULL DEFAULT 0,
  "provider_reference" character varying NULL,
  "provider" character varying NOT NULL,
  "transaction_id" character varying NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "stablecoin_withdrawals_transactions_stablecoin_withdrawal" FOREIGN KEY ("transaction_id") REFERENCES "transactions" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "stablecoin_withdrawals_transaction_id_key" to table: "stablecoin_withdrawals"
CREATE UNIQUE INDEX "stablecoin_withdrawals_transaction_id_key" ON "stablecoin_withdrawals" ("transaction_id");
