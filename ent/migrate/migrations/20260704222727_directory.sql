-- Modify "users" table
ALTER TABLE "users" ADD COLUMN "date_of_birth" timestamptz NULL;
-- Modify "verifications" table
ALTER TABLE "verifications" ADD COLUMN "number" character varying NULL;
-- Create index "verification_number_verification_type" to table: "verifications"
CREATE UNIQUE INDEX "verification_number_verification_type" ON "verifications" ("number", "verification_type");
-- Modify "businesses" table
ALTER TABLE "businesses" ALTER COLUMN "about" SET NOT NULL, ADD COLUMN "country_of_incorporation" character varying NULL, ADD COLUMN "registration_number" character varying NULL, ADD COLUMN "on_site" boolean NOT NULL DEFAULT false, ADD COLUMN "live" boolean NOT NULL DEFAULT false, ADD COLUMN "verification_status" character varying NOT NULL DEFAULT 'UNVERIFIED', ADD COLUMN "registered_by" character varying NULL, ADD CONSTRAINT "businesses_users_registered_businesses" FOREIGN KEY ("registered_by") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE SET NULL;
-- Create index "business_registration_number_name" to table: "businesses"
CREATE UNIQUE INDEX "business_registration_number_name" ON "businesses" ("registration_number", "name");
-- Create "business_locations" table
CREATE TABLE "business_locations" (
  "id" character varying NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "deleted_at" timestamptz NULL,
  "is_head_office" boolean NOT NULL DEFAULT false,
  "name" character varying NOT NULL,
  "address" character varying NOT NULL,
  "city" character varying NOT NULL,
  "state" character varying NOT NULL,
  "country" character varying NOT NULL,
  "zip_code" character varying NOT NULL,
  "latitude" double precision NULL,
  "longitude" double precision NULL,
  "active" boolean NOT NULL DEFAULT true,
  "business_locations" character varying NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "business_locations_businesses_locations" FOREIGN KEY ("business_locations") REFERENCES "businesses" ("id") ON UPDATE NO ACTION ON DELETE SET NULL
);
-- Create index "roles_name_key" to table: "roles"
CREATE UNIQUE INDEX "roles_name_key" ON "roles" ("name");
-- Modify "managers" table
ALTER TABLE "managers" ADD COLUMN "is_owner" boolean NOT NULL DEFAULT false, ADD COLUMN "role_id" character varying NOT NULL, ADD CONSTRAINT "managers_roles_managers" FOREIGN KEY ("role_id") REFERENCES "roles" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Create index "manager_user_id_business_id_role_id" to table: "managers"
CREATE UNIQUE INDEX "manager_user_id_business_id_role_id" ON "managers" ("user_id", "business_id", "role_id");
-- Create index "permissions_name_key" to table: "permissions"
CREATE UNIQUE INDEX "permissions_name_key" ON "permissions" ("name");
-- Modify "role_permissions" table
ALTER TABLE "role_permissions" DROP CONSTRAINT "role_permissions_pkey", DROP CONSTRAINT "role_permissions_permission_id", DROP CONSTRAINT "role_permissions_role_id", ADD COLUMN "id" character varying NOT NULL, ADD COLUMN "created_at" timestamptz NOT NULL, ADD COLUMN "updated_at" timestamptz NOT NULL, ADD COLUMN "deleted_at" timestamptz NULL, ADD PRIMARY KEY ("id"), ADD CONSTRAINT "role_permissions_permissions_role_permissions" FOREIGN KEY ("permission_id") REFERENCES "permissions" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT "role_permissions_roles_role_permissions" FOREIGN KEY ("role_id") REFERENCES "roles" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
