-- Drop index "wallet_identifier" from table: "wallets"
DROP INDEX "wallet_identifier";
-- Create index "wallet_identifier_currency_id_type" to table: "wallets"
CREATE INDEX "wallet_identifier_currency_id_type" ON "wallets" ("identifier", "currency_id", "type");
