-- Modify "stablecoin_networks" table
ALTER TABLE "stablecoin_networks" ALTER COLUMN "logo_url" SET NOT NULL, ALTER COLUMN "logo_url" SET DEFAULT '';
-- Create index "stablecoin_networks_slug_key" to table: "stablecoin_networks"
CREATE UNIQUE INDEX "stablecoin_networks_slug_key" ON "stablecoin_networks" ("slug");
