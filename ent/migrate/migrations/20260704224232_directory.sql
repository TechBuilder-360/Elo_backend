-- Modify "role_permissions" table
ALTER TABLE "role_permissions" ADD COLUMN "permission_id" character varying NOT NULL, ADD COLUMN "role_id" character varying NOT NULL, ADD CONSTRAINT "role_permissions_permissions_role_permissions" FOREIGN KEY ("permission_id") REFERENCES "permissions" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT "role_permissions_roles_role_permissions" FOREIGN KEY ("role_id") REFERENCES "roles" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Create index "rolepermission_role_id_permission_id" to table: "role_permissions"
CREATE UNIQUE INDEX "rolepermission_role_id_permission_id" ON "role_permissions" ("role_id", "permission_id");
-- Drop "permission_role_permissions" table
DROP TABLE "permission_role_permissions";
-- Drop "role_role_permissions" table
DROP TABLE "role_role_permissions";
