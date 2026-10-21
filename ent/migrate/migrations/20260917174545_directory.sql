-- Modify "transactions" table
ALTER TABLE "transactions" DROP COLUMN "invoice_transaction";
-- Modify "invoices" table
ALTER TABLE "invoices" ADD COLUMN "transaction_invoice" character varying NOT NULL, ADD CONSTRAINT "invoices_transactions_invoice" FOREIGN KEY ("transaction_invoice") REFERENCES "transactions" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
