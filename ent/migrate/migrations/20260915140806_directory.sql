-- Modify "ledgers" table
ALTER TABLE "ledgers" ADD COLUMN "affected_balance" character varying NOT NULL;
-- Modify "services" table
ALTER TABLE "services" DROP COLUMN "min", DROP COLUMN "max", DROP COLUMN "fee", ADD COLUMN "min_value" bigint NOT NULL DEFAULT 0, ADD COLUMN "max_value" bigint NOT NULL DEFAULT 0;
-- Create "fee_rules" table
CREATE TABLE "fee_rules" (
  "id" character varying NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "deleted_at" timestamptz NULL,
  "name" character varying NOT NULL,
  "fee_type" character varying NOT NULL,
  "calculation_method" character varying NOT NULL,
  "flat_amount" bigint NULL DEFAULT 0,
  "percentage_rate" double precision NULL DEFAULT 0,
  "currency" character varying NOT NULL,
  "is_active" boolean NOT NULL DEFAULT true,
  "effective_from" timestamptz NOT NULL,
  "effective_to" timestamptz NULL,
  "min_value" bigint NULL,
  "max_value" bigint NULL,
  "service_id" character varying NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fee_rules_services_fee_rules" FOREIGN KEY ("service_id") REFERENCES "services" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create "fee_tiers" table
CREATE TABLE "fee_tiers" (
  "id" character varying NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "deleted_at" timestamptz NULL,
  "min_amount" bigint NOT NULL DEFAULT 0,
  "max_amount" bigint NULL,
  "charge_type" character varying NOT NULL,
  "flat_amount" bigint NULL DEFAULT 0,
  "percentage_rate" double precision NULL DEFAULT 0,
  "priority" bigint NOT NULL DEFAULT 0,
  "fee_rule_tiers" character varying NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fee_tiers_fee_rules_tiers" FOREIGN KEY ("fee_rule_tiers") REFERENCES "fee_rules" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create "subscription_plans" table
CREATE TABLE "subscription_plans" (
  "id" character varying NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "deleted_at" timestamptz NULL,
  "name" character varying NOT NULL,
  "slug" character varying NOT NULL,
  "description" text NULL,
  "billing_interval" character varying NOT NULL,
  "interval_count" bigint NOT NULL DEFAULT 1,
  "amount" bigint NOT NULL,
  "currency" character varying NOT NULL,
  "trial_period_days" bigint NOT NULL DEFAULT 0,
  "features" jsonb NULL,
  "max_seats" bigint NULL,
  "is_active" boolean NOT NULL DEFAULT true,
  "is_public" boolean NOT NULL DEFAULT true,
  PRIMARY KEY ("id")
);
-- Create index "subscription_plans_slug_key" to table: "subscription_plans"
CREATE UNIQUE INDEX "subscription_plans_slug_key" ON "subscription_plans" ("slug");
-- Create "subscriptions" table
CREATE TABLE "subscriptions" (
  "id" character varying NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "deleted_at" timestamptz NULL,
  "status" character varying NOT NULL DEFAULT 'TRIALING',
  "locked_amount" bigint NOT NULL,
  "locked_currency" character varying NOT NULL,
  "current_period_start" timestamptz NOT NULL,
  "current_period_end" timestamptz NOT NULL,
  "trial_ends_at" timestamptz NULL,
  "cancelled_at" timestamptz NULL,
  "cancel_at_period_end" timestamptz NULL,
  "business_subscriptions" character varying NOT NULL,
  "subscription_plan_subscriptions" character varying NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "subscriptions_businesses_subscriptions" FOREIGN KEY ("business_subscriptions") REFERENCES "businesses" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "subscriptions_subscription_plans_subscriptions" FOREIGN KEY ("subscription_plan_subscriptions") REFERENCES "subscription_plans" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create "invoices" table
CREATE TABLE "invoices" (
  "id" character varying NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "deleted_at" timestamptz NULL,
  "reference" character varying NOT NULL,
  "status" character varying NOT NULL,
  "amount" bigint NOT NULL,
  "currency" character varying NOT NULL,
  "period_start" timestamptz NOT NULL,
  "period_end" timestamptz NOT NULL,
  "due_at" timestamptz NOT NULL,
  "paid_at" timestamptz NULL,
  "subscription_invoices" character varying NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "invoices_subscriptions_invoices" FOREIGN KEY ("subscription_invoices") REFERENCES "subscriptions" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "invoices_reference_key" to table: "invoices"
CREATE UNIQUE INDEX "invoices_reference_key" ON "invoices" ("reference");
-- Create "nuban_banks" table
CREATE TABLE "nuban_banks" (
  "id" character varying NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "deleted_at" timestamptz NULL,
  "name" character varying NOT NULL,
  "cbn_code" character varying NOT NULL,
  "short_name" character varying NOT NULL,
  "slug" character varying NOT NULL,
  "active" boolean NOT NULL DEFAULT true,
  PRIMARY KEY ("id")
);
-- Create index "nuban_banks_cbn_code_key" to table: "nuban_banks"
CREATE UNIQUE INDEX "nuban_banks_cbn_code_key" ON "nuban_banks" ("cbn_code");
-- Create index "nuban_banks_name_key" to table: "nuban_banks"
CREATE UNIQUE INDEX "nuban_banks_name_key" ON "nuban_banks" ("name");
-- Create index "nuban_banks_short_name_key" to table: "nuban_banks"
CREATE UNIQUE INDEX "nuban_banks_short_name_key" ON "nuban_banks" ("short_name");
-- Create index "nuban_banks_slug_key" to table: "nuban_banks"
CREATE UNIQUE INDEX "nuban_banks_slug_key" ON "nuban_banks" ("slug");
-- Create index "nubanbank_cbn_code" to table: "nuban_banks"
CREATE UNIQUE INDEX "nubanbank_cbn_code" ON "nuban_banks" ("cbn_code");
-- Create "nuban_beneficiaries" table
CREATE TABLE "nuban_beneficiaries" (
  "id" character varying NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "deleted_at" timestamptz NULL,
  "account_name" character varying NOT NULL,
  "account_number" character varying NOT NULL,
  "bank_name" character varying NOT NULL,
  "nickname" character varying NULL,
  "disabled" boolean NOT NULL DEFAULT false,
  "owner_id" character varying NOT NULL,
  "nuban_bank_id" character varying NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "nuban_beneficiaries_ledger_owners_nuban_beneficiaries" FOREIGN KEY ("owner_id") REFERENCES "ledger_owners" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "nuban_beneficiaries_nuban_banks_beneficiary_bank" FOREIGN KEY ("nuban_bank_id") REFERENCES "nuban_banks" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "nubanbeneficiary_owner_id" to table: "nuban_beneficiaries"
CREATE INDEX "nubanbeneficiary_owner_id" ON "nuban_beneficiaries" ("owner_id");
-- Create index "nubanbeneficiary_owner_id_account_number_nuban_bank_id" to table: "nuban_beneficiaries"
CREATE UNIQUE INDEX "nubanbeneficiary_owner_id_account_number_nuban_bank_id" ON "nuban_beneficiaries" ("owner_id", "account_number", "nuban_bank_id");
-- Modify "nuban_transfers" table
ALTER TABLE "nuban_transfers" DROP COLUMN "sender_account_name", DROP COLUMN "sender_account_number", DROP COLUMN "sender_bank_name", DROP COLUMN "sender_bank_code", ADD COLUMN "nuban_bank_id" character varying NOT NULL, ADD CONSTRAINT "nuban_transfers_nuban_banks_bank" FOREIGN KEY ("nuban_bank_id") REFERENCES "nuban_banks" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Create "nuban_beneficiary_nuban_beneficiaries" table
CREATE TABLE "nuban_beneficiary_nuban_beneficiaries" (
  "nuban_beneficiary_id" character varying NOT NULL,
  "nuban_transfer_id" character varying NOT NULL,
  PRIMARY KEY ("nuban_beneficiary_id", "nuban_transfer_id"),
  CONSTRAINT "nuban_beneficiary_nuban_beneficiaries_nuban_beneficiary_id" FOREIGN KEY ("nuban_beneficiary_id") REFERENCES "nuban_beneficiaries" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "nuban_beneficiary_nuban_beneficiaries_nuban_transfer_id" FOREIGN KEY ("nuban_transfer_id") REFERENCES "nuban_transfers" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create "nuban_blacklists" table
CREATE TABLE "nuban_blacklists" (
  "id" character varying NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "deleted_at" timestamptz NULL,
  "account_name" character varying NOT NULL,
  "account_number" character varying NOT NULL,
  "bank_name" character varying NOT NULL,
  "reason" character varying NOT NULL,
  "disabled" boolean NOT NULL DEFAULT true,
  "nuban_bank_id" character varying NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "nuban_blacklists_nuban_banks_blacklist" FOREIGN KEY ("nuban_bank_id") REFERENCES "nuban_banks" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "nubanblacklist_account_number_nuban_bank_id" to table: "nuban_blacklists"
CREATE UNIQUE INDEX "nubanblacklist_account_number_nuban_bank_id" ON "nuban_blacklists" ("account_number", "nuban_bank_id");
-- Create "provider_banks" table
CREATE TABLE "provider_banks" (
  "id" bigint NOT NULL GENERATED BY DEFAULT AS IDENTITY,
  "provider_code" character varying NOT NULL,
  "provider" character varying NOT NULL,
  "active" boolean NOT NULL DEFAULT true,
  "nuban_bank_id" character varying NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "provider_banks_nuban_banks_provider_banks" FOREIGN KEY ("nuban_bank_id") REFERENCES "nuban_banks" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "providerbank_provider_nuban_bank_id" to table: "provider_banks"
CREATE UNIQUE INDEX "providerbank_provider_nuban_bank_id" ON "provider_banks" ("provider", "nuban_bank_id");
-- Create index "providerbank_provider_provider_code" to table: "provider_banks"
CREATE UNIQUE INDEX "providerbank_provider_provider_code" ON "provider_banks" ("provider", "provider_code");
-- Modify "transactions" table
ALTER TABLE "transactions" DROP COLUMN "fee", DROP COLUMN "transaction_wallet", ADD COLUMN "net_amount" bigint NULL, ADD COLUMN "entry" character varying NOT NULL DEFAULT 'DEBIT', ADD COLUMN "invoice_transaction" character varying NOT NULL, ADD CONSTRAINT "transactions_invoices_transaction" FOREIGN KEY ("invoice_transaction") REFERENCES "invoices" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Create index "transactions_invoice_transaction_key" to table: "transactions"
CREATE UNIQUE INDEX "transactions_invoice_transaction_key" ON "transactions" ("invoice_transaction");
-- Create "transaction_fees" table
CREATE TABLE "transaction_fees" (
  "id" character varying NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "deleted_at" timestamptz NULL,
  "type" character varying NOT NULL,
  "amount" bigint NOT NULL,
  "currency" character varying NOT NULL,
  "rate" double precision NULL,
  "base_amount" bigint NULL,
  "service_identifier" character varying NULL,
  "transaction_fees" character varying NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "transaction_fees_transactions_fees" FOREIGN KEY ("transaction_fees") REFERENCES "transactions" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Drop "business_features" table
DROP TABLE "business_features";
