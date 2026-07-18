-- Modify "businesses" table
ALTER TABLE "businesses" ALTER COLUMN "active" SET DEFAULT true, ALTER COLUMN "disabled" SET DEFAULT false;
-- Create "currencies" table
CREATE TABLE "currencies" (
  "id" character varying NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "deleted_at" timestamptz NULL,
  "name" character varying NOT NULL,
  "symbol" character varying NOT NULL,
  "code" character varying NOT NULL,
  "is_fiat" boolean NOT NULL DEFAULT true,
  "active" boolean NOT NULL DEFAULT true,
  "multipler" bigint NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "currencies_name_key" to table: "currencies"
CREATE UNIQUE INDEX "currencies_name_key" ON "currencies" ("name");
-- Create "wallets" table
CREATE TABLE "wallets" (
  "id" character varying NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "deleted_at" timestamptz NULL,
  "type" character varying NOT NULL DEFAULT 'TREASURY',
  "available_balance" bigint NOT NULL,
  "ledger_balance" bigint NOT NULL,
  "holding_balance" bigint NOT NULL,
  "owner" character varying NOT NULL DEFAULT 'USER',
  "identifier" character varying NOT NULL,
  "active" boolean NOT NULL DEFAULT true,
  "currency_id" character varying NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "wallets_currencies_wallets" FOREIGN KEY ("currency_id") REFERENCES "currencies" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "wallet_identifier" to table: "wallets"
CREATE INDEX "wallet_identifier" ON "wallets" ("identifier");
