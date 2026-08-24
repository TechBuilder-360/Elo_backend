-- Modify "stablecoin_wallets" table
ALTER TABLE "stablecoin_wallets" ADD COLUMN "disabled" boolean NOT NULL DEFAULT false;
