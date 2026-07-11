-- Modify "business_documents" table
ALTER TABLE "business_documents" ADD COLUMN "created_at" timestamptz NOT NULL, ADD COLUMN "updated_at" timestamptz NOT NULL, ADD COLUMN "deleted_at" timestamptz NULL;
ALTER TABLE "business_documents" DROP CONSTRAINT "business_documents_pkey";
ALTER TABLE "business_documents" ALTER COLUMN "id" DROP IDENTITY IF EXISTS;
ALTER TABLE "business_documents" ALTER COLUMN "id" TYPE character varying USING "id"::character varying;
ALTER TABLE "business_documents" ADD PRIMARY KEY ("id");
-- Create "kyb_documents" table
CREATE TABLE "kyb_documents" (
  "id" character varying NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "deleted_at" timestamptz NULL,
  "name" character varying NOT NULL,
  "required" boolean NOT NULL DEFAULT true,
  "active" boolean NOT NULL DEFAULT true,
  PRIMARY KEY ("id")
);
-- Create index "kyb_documents_name_key" to table: "kyb_documents"
CREATE UNIQUE INDEX "kyb_documents_name_key" ON "kyb_documents" ("name");
-- Create "kyb_document_kyb_documents" table
CREATE TABLE "kyb_document_kyb_documents" (
  "kyb_document_id" character varying NOT NULL,
  "business_document_id" character varying NOT NULL,
  PRIMARY KEY ("kyb_document_id", "business_document_id"),
  CONSTRAINT "kyb_document_kyb_documents_business_document_id" FOREIGN KEY ("business_document_id") REFERENCES "business_documents" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "kyb_document_kyb_documents_kyb_document_id" FOREIGN KEY ("kyb_document_id") REFERENCES "kyb_documents" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Modify "businesses" table
ALTER TABLE "businesses" ADD COLUMN "cover_image" character varying NULL;
-- Create "kyb_messages" table
CREATE TABLE "kyb_messages" (
  "id" character varying NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "deleted_at" timestamptz NULL,
  "message" character varying NOT NULL,
  "status" character varying NOT NULL DEFAULT 'OPEN',
  "business_kyb_messages" character varying NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "kyb_messages_businesses_kyb_messages" FOREIGN KEY ("business_kyb_messages") REFERENCES "businesses" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
