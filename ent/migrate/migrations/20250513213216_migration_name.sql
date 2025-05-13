-- Rename a column from "business_business_manager" to "business_manager"
ALTER TABLE "managers" RENAME COLUMN "business_business_manager" TO "business_manager";
-- Rename a column from "user_user_manager" to "user_manager"
ALTER TABLE "managers" RENAME COLUMN "user_user_manager" TO "user_manager";
-- Modify "managers" table
ALTER TABLE "managers" DROP CONSTRAINT "managers_businesses_business_manager", DROP CONSTRAINT "managers_users_user_manager", ADD CONSTRAINT "managers_businesses_manager" FOREIGN KEY ("business_manager") REFERENCES "businesses" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT "managers_users_manager" FOREIGN KEY ("user_manager") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Rename a column from "business_business_social" to "business_social"
ALTER TABLE "socials" RENAME COLUMN "business_business_social" TO "business_social";
-- Modify "socials" table
ALTER TABLE "socials" DROP CONSTRAINT "socials_businesses_business_social", ADD CONSTRAINT "socials_businesses_social" FOREIGN KEY ("business_social") REFERENCES "businesses" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
