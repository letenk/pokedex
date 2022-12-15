## up: starts all containers in the background without forcing build
up:
	@echo "Starting Docker images..."
	docker-compose up -d
	@echo "Docker images started!"

## up_build: stops docker-compose (if running), builds projects and starts docker compose
up_build: 
	@echo "Stopping docker images (if running...)"
	docker-compose down
	@echo "Building (when required) and starting docker images..."
	docker-compose up --build -d
	@echo "Docker images built and started!"

## up: starts all containers in the background without forcing build
down:
	@echo "Stopping docker compose..."
	docker-compose down
	@echo "Done!"

## test: Run all test in this app
test:
	@echo "All tests are running..."
	go test -v ./...
	@echo "Test finished"

## test: Run all test with clean cache in this app
test_nocache:
	@echo "Clean all cache..."
	go clean -testcache
	@echo "All tests are running..."
	go test -v ./...
	@echo "Test finished"

## test_cover: Run all test with coverage
test_cover:
	@echo "All test are running with coverage..."
	go test ./... -v -cover

## test: Run all test with clean cache and coverage
test_cover_nocache:
	@echo "Clean all cache..."
	go clean -testcache
	@echo "All tests are running..."
	go test ./... -v -cover
	@echo "Test finished"

# Run app
run: 
	go run main.go

# Create table and seed sample data user, category, types
run-migrate-seed: 
	psql -d postgresql://root:secret@localhost:5432/pokedex -f scripts/db/dump.sql

# Create table and seed sample data user, category, types for database test
run-migrate-seed-test: 
	psql -d postgresql://root:secret@localhost:5432/pokedex_test -f scripts/db/dump.sql

.PHONY: up down test test_nocache test_cover test_cover_nocache run