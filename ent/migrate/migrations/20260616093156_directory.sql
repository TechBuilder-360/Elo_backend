-- Modify "verifications" table
ALTER TABLE "verifications" ADD COLUMN "data" jsonb NOT NULL, ADD COLUMN "is_obsolate" boolean NOT NULL DEFAULT false;
-- Create "request_verifications" table
CREATE TABLE "request_verifications" (
  "id" character varying NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "deleted_at" timestamptz NULL,
  "reference_id" character varying NOT NULL,
  "verification_type" character varying NOT NULL,
  "provider" character varying NOT NULL,
  "link" character varying NOT NULL,
  "provider_link" character varying NOT NULL,
  "status" character varying NOT NULL DEFAULT 'PENDING',
  PRIMARY KEY ("id")
);
-- Create index "request_verifications_reference_id_key" to table: "request_verifications"
CREATE UNIQUE INDEX "request_verifications_reference_id_key" ON "request_verifications" ("reference_id");
-- Create "request_verification_business" table
CREATE TABLE "request_verification_business" (
  "request_verification_id" character varying NOT NULL,
  "business_id" character varying NOT NULL,
  PRIMARY KEY ("request_verification_id", "business_id"),
  CONSTRAINT "request_verification_business_business_id" FOREIGN KEY ("business_id") REFERENCES "businesses" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "request_verification_business_request_verification_id" FOREIGN KEY ("request_verification_id") REFERENCES "request_verifications" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create "request_verification_user" table
CREATE TABLE "request_verification_user" (
  "request_verification_id" character varying NOT NULL,
  "user_id" character varying NOT NULL,
  PRIMARY KEY ("request_verification_id", "user_id"),
  CONSTRAINT "request_verification_user_request_verification_id" FOREIGN KEY ("request_verification_id") REFERENCES "request_verifications" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "request_verification_user_user_id" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
