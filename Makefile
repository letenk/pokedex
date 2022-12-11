# Migrate up
up:
	@echo "Starting Docker images..."
	docker-compose up -d
	@echo "Docker images started!"

down:
	@echo "Stopping docker compose..."
	docker-compose down
	@echo "Done!"
# Running test
test:
	@echo "Testing started..."
	go test -v -cover ./...

# Run app
run: 
	go run main.go

.PHONY: up down test run