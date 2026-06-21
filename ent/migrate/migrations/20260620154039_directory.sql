-- Modify "request_verifications" table
ALTER TABLE "request_verifications" ADD COLUMN IF NOT EXISTS "message" character varying NULL;
