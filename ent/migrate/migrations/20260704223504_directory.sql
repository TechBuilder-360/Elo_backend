-- Modify "role_permissions" table
ALTER TABLE "role_permissions" DROP COLUMN "role_id", DROP COLUMN "permission_id";
-- Create "permission_role_permissions" table
CREATE TABLE "permission_role_permissions" (
  "permission_id" character varying NOT NULL,
  "role_permission_id" character varying NOT NULL,
  PRIMARY KEY ("permission_id", "role_permission_id"),
  CONSTRAINT "permission_role_permissions_permission_id" FOREIGN KEY ("permission_id") REFERENCES "permissions" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "permission_role_permissions_role_permission_id" FOREIGN KEY ("role_permission_id") REFERENCES "role_permissions" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create "role_role_permissions" table
CREATE TABLE "role_role_permissions" (
  "role_id" character varying NOT NULL,
  "role_permission_id" character varying NOT NULL,
  PRIMARY KEY ("role_id", "role_permission_id"),
  CONSTRAINT "role_role_permissions_role_id" FOREIGN KEY ("role_id") REFERENCES "roles" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "role_role_permissions_role_permission_id" FOREIGN KEY ("role_permission_id") REFERENCES "role_permissions" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
