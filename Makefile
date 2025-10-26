ifneq (,$(wildcard .env))
    include .env
    export
endif

build:
	go build -v ./cmd/
dev:
	go run cmd/main.go
swagger-docs:
	swag init -g cmd/main.go --output docs --parseDependency --parseInternal
MIGRATION_NAME ?=

migration:
	@if [ -z "$(MIGRATION_NAME)" ]; then echo "Error: MIGRATION_NAME is required"; exit 1; fi
	@test -d $(MIGRATION_DIR) || (echo "Error: $(MIGRATION_DIR) does not exist" && exit 1)
	goose create -dir $(MIGRATION_DIR) $(MIGRATION_NAME) sql

migrate:
	@if [ -z "$(DATABASE_URL)" ]; then echo "Error: DATABASE_URL is required"; exit 1; fi
	@test -d $(MIGRATION_DIR) || (echo "Error: $(MIGRATION_DIR) does not exist" && exit 1)
	goose -dir $(MIGRATION_DIR) postgres "$(DATABASE_URL)" up

migrate_down:
	@if [ -z "$(DATABASE_URL)" ]; then echo "Error: DATABASE_URL is required"; exit 1; fi
	@test -d $(MIGRATION_DIR) || (echo "Error: $(MIGRATION_DIR) does not exist" && exit 1)
	goose -dir $(MIGRATION_DIR) postgres "$(DATABASE_URL)" down