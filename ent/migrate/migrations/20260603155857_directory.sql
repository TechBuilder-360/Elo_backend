-- Modify "business_documents" table
ALTER TABLE "business_documents" ADD COLUMN "type" character varying NOT NULL;
-- Modify "user_documents" table
ALTER TABLE "user_documents" DROP CONSTRAINT "user_documents_users_user_documents", DROP COLUMN "business_user_documents", ALTER COLUMN "user_user_documents" SET NOT NULL, ADD CONSTRAINT "user_documents_users_user_documents" FOREIGN KEY ("user_user_documents") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Create "verifications" table
CREATE TABLE "verifications" (
  "id" character varying NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "deleted_at" timestamptz NULL,
  "provider" character varying NOT NULL,
  "verification_type" character varying NOT NULL,
  "status" character varying NOT NULL DEFAULT 'IN_PROGRESS',
  "reference_id" character varying NOT NULL,
  "provider_reference" character varying NOT NULL,
  "metadata" jsonb NOT NULL,
  "verified_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create "verification_business" table
CREATE TABLE "verification_business" (
  "verification_id" character varying NOT NULL,
  "business_id" character varying NOT NULL,
  PRIMARY KEY ("verification_id", "business_id"),
  CONSTRAINT "verification_business_business_id" FOREIGN KEY ("business_id") REFERENCES "businesses" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "verification_business_verification_id" FOREIGN KEY ("verification_id") REFERENCES "verifications" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create "verification_user" table
CREATE TABLE "verification_user" (
  "verification_id" character varying NOT NULL,
  "user_id" character varying NOT NULL,
  PRIMARY KEY ("verification_id", "user_id"),
  CONSTRAINT "verification_user_user_id" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "verification_user_verification_id" FOREIGN KEY ("verification_id") REFERENCES "verifications" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
