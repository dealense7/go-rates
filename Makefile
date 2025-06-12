-include .env

# ------------- MIGRATIONS
migrate-up:
	@GOOSE_DRIVER=$(DB_DRIVER) GOOSE_DBSTRING="$(DB_STRING)" goose up -dir $(MIGRATION_FOLDER)

migrate-down:
	@GOOSE_DRIVER=$(DB_DRIVER) GOOSE_DBSTRING="$(DB_STRING)" goose down -dir $(MIGRATION_FOLDER)

migrate-status:
	@GOOSE_DRIVER=$(DB_DRIVER) GOOSE_DBSTRING="$(DB_STRING)" goose status -dir database/migrations

# ------------ PREPARE
seed:
	@echo "==> Running seed script…"
	@go run cmd/seed/main.go

prepare: migrate-up seed
	@echo "✅ Database ready (migrations applied & seed run)."
