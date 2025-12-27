-- Create "businesses" table
CREATE TABLE "businesses" (
  "id" character varying NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "deleted_at" timestamptz NULL,
  "category" character varying NOT NULL DEFAULT 'others',
  "name" character varying NOT NULL,
  "logo" character varying NOT NULL,
  "email" character varying NOT NULL,
  "website" character varying NOT NULL,
  "disabled" boolean NOT NULL DEFAULT true,
  "disabled_at" timestamptz NOT NULL,
  "disable_reason" character varying NOT NULL,
  "verified" boolean NOT NULL DEFAULT false,
  "verified_at" timestamptz NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "businesses_name_key" to table: "businesses"
CREATE UNIQUE INDEX "businesses_name_key" ON "businesses" ("name");
-- Create "users" table
CREATE TABLE "users" (
  "id" character varying NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "deleted_at" timestamptz NULL,
  "first_name" character varying NOT NULL,
  "last_name" character varying NOT NULL,
  "middle_name" character varying NOT NULL,
  "display_name" character varying NOT NULL,
  "email_address" character varying NOT NULL,
  "email_verified" boolean NOT NULL DEFAULT false,
  "email_verified_at" timestamptz NOT NULL,
  "phone_number" character varying NOT NULL,
  "avatar" character varying NOT NULL,
  "disabled" boolean NOT NULL DEFAULT false,
  "tier" smallint NOT NULL DEFAULT 0,
  PRIMARY KEY ("id")
);
-- Create index "users_email_address_key" to table: "users"
CREATE UNIQUE INDEX "users_email_address_key" ON "users" ("email_address");
-- Create "managers" table
CREATE TABLE "managers" (
  "id" character varying NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "deleted_at" timestamptz NULL,
  "disabled" boolean NOT NULL DEFAULT false,
  "disable_reason" character varying NOT NULL,
  "disabled_at" timestamptz NOT NULL,
  "business_id" character varying NOT NULL,
  "user_id" character varying NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "managers_businesses_manages" FOREIGN KEY ("business_id") REFERENCES "businesses" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "managers_users_manages" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create "permissions" table
CREATE TABLE "permissions" (
  "id" character varying NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "deleted_at" timestamptz NULL,
  "name" character varying NOT NULL,
  "description" character varying NOT NULL,
  PRIMARY KEY ("id")
);
-- Create "roles" table
CREATE TABLE "roles" (
  "id" character varying NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "deleted_at" timestamptz NULL,
  "name" character varying NOT NULL,
  "description" character varying NOT NULL,
  PRIMARY KEY ("id")
);
-- Create "role_permissions" table
CREATE TABLE "role_permissions" (
  "role_id" character varying NOT NULL,
  "permission_id" character varying NOT NULL,
  PRIMARY KEY ("role_id", "permission_id"),
  CONSTRAINT "role_permissions_permission_id" FOREIGN KEY ("permission_id") REFERENCES "permissions" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "role_permissions_role_id" FOREIGN KEY ("role_id") REFERENCES "roles" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create "socials" table
CREATE TABLE "socials" (
  "id" character varying NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "deleted_at" timestamptz NULL,
  "name" character varying NOT NULL,
  "url" character varying NOT NULL,
  "business_id" character varying NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "socials_businesses_socials" FOREIGN KEY ("business_id") REFERENCES "businesses" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
