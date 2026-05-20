-- Create "business_features" table
CREATE TABLE "business_features" (
  "id" character varying NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "deleted_at" timestamptz NULL,
  "name" character varying NOT NULL,
  "identifier" character varying NOT NULL,
  "require_subscription" boolean NOT NULL DEFAULT false,
  "active" boolean NOT NULL DEFAULT true,
  "min" bigint NOT NULL DEFAULT 0,
  "max" bigint NOT NULL DEFAULT 0,
  "fee" jsonb NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "business_features_identifier_key" to table: "business_features"
CREATE UNIQUE INDEX "business_features_identifier_key" ON "business_features" ("identifier");
-- Create index "business_features_name_key" to table: "business_features"
CREATE UNIQUE INDEX "business_features_name_key" ON "business_features" ("name");
-- Create "services" table
CREATE TABLE "services" (
  "id" character varying NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "deleted_at" timestamptz NULL,
  "name" character varying NOT NULL,
  "identifier" character varying NOT NULL,
  "require_subscription" boolean NOT NULL DEFAULT false,
  "active" boolean NOT NULL DEFAULT true,
  "min" bigint NOT NULL DEFAULT 0,
  "max" bigint NOT NULL DEFAULT 0,
  "fee" jsonb NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "services_identifier_key" to table: "services"
CREATE UNIQUE INDEX "services_identifier_key" ON "services" ("identifier");
-- Create index "services_name_key" to table: "services"
CREATE UNIQUE INDEX "services_name_key" ON "services" ("name");
