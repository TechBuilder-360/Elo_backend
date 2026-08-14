-- Modify "businesses" table
ALTER TABLE "businesses" ADD COLUMN "tax_identification_number" character varying NULL;
-- Create index "currency_name_code" to table: "currencies"
CREATE UNIQUE INDEX "currency_name_code" ON "currencies" ("name", "code");
