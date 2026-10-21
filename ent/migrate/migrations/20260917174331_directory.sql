-- Modify "transactions" table
ALTER TABLE "transactions" DROP CONSTRAINT "transactions_invoices_transaction", ALTER COLUMN "invoice_transaction" DROP NOT NULL, ADD CONSTRAINT "transactions_invoices_transaction" FOREIGN KEY ("invoice_transaction") REFERENCES "invoices" ("id") ON UPDATE NO ACTION ON DELETE SET NULL;
