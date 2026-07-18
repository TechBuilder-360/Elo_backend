-- Drop index "wallet_identifier_currency_id_type" from table: "wallets"
DROP INDEX "wallet_identifier_currency_id_type";
-- Create index "wallet_identifier_currency_id_type" to table: "wallets"
CREATE UNIQUE INDEX "wallet_identifier_currency_id_type" ON "wallets" ("identifier", "currency_id", "type");
