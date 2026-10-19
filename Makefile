include .envrc
MIGRATIONS_DIR = ./cmd/migrate/migrations

.PHONY: migrate-create
migrate-create:
	@echo "Creating migration..."
	@migrate create -seq -ext sql -dir $(MIGRATIONS_DIR) $(filter-out $@,$(MAKECMDGOALS))

.PHONY: migrate-up
migrate-up:
	@echo "Running migrations..."
	@migrate -path $(MIGRATIONS_DIR) -database $(DB_ADDR) up

.PHONY: migrate-down
migrate-down:
	@echo "Rolling back migrations..."
	@migrate -path $(MIGRATIONS_DIR) -database $(DB_ADDR) down $(filter-out $@,$(MAKECMDGOALS))

.PHONY: migrate-force
migrate-force:
	@echo "Forcing migration version $(version)..."
	@migrate -path $(MIGRATIONS_DIR) -database $(DB_ADDR) force $(version)

.PHONY: seed
seed:
	@echo "Seeding database..."
	@go run ./cmd/migrate/seed/main.go

.PHONY: gen-docs
gen-docs:
	@echo "Generating API documentation..."
	@swag init -g main.go -d cmd/api,internal && swag fmt