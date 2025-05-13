# Run project requirements

generate_ent_schema:
	go generate ./ent

migrate_schema:
	atlas migrate diff migration_name \
      --dir "file://ent/migrate/migrations" \
      --to "ent://ent/schema" \
      --dev-url "docker://postgres/15/test?search_path=public"

migration_compute_hash:
	atlas migrate hash \
	 	--dir "file://ent/migrate/migrations"
