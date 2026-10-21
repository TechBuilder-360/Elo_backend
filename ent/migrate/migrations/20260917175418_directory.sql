-- Modify "nuban_transfers" table
ALTER TABLE "nuban_transfers" ALTER COLUMN "session_id" DROP NOT NULL, ALTER COLUMN "provider_reference" DROP NOT NULL;
