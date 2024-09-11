include .env

# make dev - to start development server
# Set APP_ENV to "development" and run the application
dev:
	@export APP_ENV=development && $(MAKE) run

# make prod - to start production server
# Set APP_ENV to "production" and run the application
prod:
	@export APP_ENV=production && $(MAKE) run

migration_up:
	@migrate -database ${DATABASE_URL}?sslmode=disable -path db/migrations up

migration_down:
	@migrate -database ${DATABASE_URL}?sslmode=disable -path db/migrations down

run:
	@echo "Running in [$(APP_ENV)] mode"
	@go run main.go

