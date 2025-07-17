MIGRATE = migrate -path ./internal/database/migrations -database $(DATABASE_URL)

migrate-up:
	@echo "Running migrations on: $(DATABASE_URL)"
	$(MIGRATE) up

migrate-down:
	@echo "Running migrations on: $(DATABASE_URL)"
	$(MIGRATE) down 1

migrate-new:
	migrate create -ext sql -dir ./internal/database/migrations -seq $(name)
