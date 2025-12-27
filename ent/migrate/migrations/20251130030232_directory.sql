-- Modify "businesses" table
ALTER TABLE "businesses" ALTER COLUMN "logo" DROP NOT NULL, ALTER COLUMN "website" DROP NOT NULL, ALTER COLUMN "disable_reason" DROP NOT NULL, ALTER COLUMN "verified_at" DROP NOT NULL;
-- Modify "managers" table
ALTER TABLE "managers" ALTER COLUMN "disable_reason" DROP NOT NULL, ALTER COLUMN "disabled_at" DROP NOT NULL;
-- Modify "users" table
ALTER TABLE "users" ALTER COLUMN "middle_name" DROP NOT NULL, ALTER COLUMN "display_name" DROP NOT NULL, ALTER COLUMN "email_verified_at" DROP NOT NULL, ALTER COLUMN "phone_number" DROP NOT NULL, ALTER COLUMN "avatar" DROP NOT NULL;
