-- Modify "fee_rules" table
ALTER TABLE "fee_rules" ADD COLUMN "min_fee" bigint NULL, ADD COLUMN "max_fee" bigint NULL;
-- Drop index "transactions_reference_key" from table: "transactions"
DROP INDEX "transactions_reference_key";
-- Modify "ledgers" table
ALTER TABLE "ledgers" DROP CONSTRAINT "ledgers_transactions_ledger_entries", ADD CONSTRAINT "ledgers_transactions_ledger_entry" FOREIGN KEY ("transaction_id") REFERENCES "transactions" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Create index "ledgers_transaction_id_key" to table: "ledgers"
CREATE UNIQUE INDEX "ledgers_transaction_id_key" ON "ledgers" ("transaction_id");
