# Run project requirements

generate-ent-schema:
	 go run -mod=mod entgo.io/ent/cmd/ent generate --feature sql/upsert ./ent/schema

migrate-schema:
	atlas migrate diff directory \
      --dir "file://ent/migrate/migrations" \
      --to "ent://ent/schema" \
      --dev-url "docker://postgres/latest/test?search_path=public"

migration-compute_hash:
	atlas migrate hash \
	 	--dir "file://ent/migrate/migrations"

migrate:
	go run cmd/migration/migration.go

start-server:
	go run cmd/http/server.go 


gqlgen:
	go run github.com/99designs/gqlgen generate