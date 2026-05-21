-- Create index "provider_slug" to table: "providers"
CREATE UNIQUE INDEX IF NOT EXISTS "provider_slug" ON "providers" ("slug");
