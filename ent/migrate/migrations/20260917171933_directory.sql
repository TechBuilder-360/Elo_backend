-- Drop index "nubanbank_cbn_code" from table: "nuban_banks"
DROP INDEX "nubanbank_cbn_code";
-- Create index "nubanbank_cbn_code_slug" to table: "nuban_banks"
CREATE UNIQUE INDEX "nubanbank_cbn_code_slug" ON "nuban_banks" ("cbn_code", "slug");
