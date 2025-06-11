.PHONY: help build run test clean docker-up docker-down dev

# Default target
help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Build the application
	go build -o bin/app ./cmd

run: ## Run the application
	go run ./cmd

test: ## Run tests
	go test -v ./...

test-coverage: ## Run tests with coverage
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

clean: ## Clean build artifacts
	rm -rf bin/
	rm -f coverage.out coverage.html

dev: ## Run in development mode with hot reload (requires air)
	air

install-dev-tools: ## Install development tools
	go install github.com/air-verse/air@latest
	go install github.com/google/wire/cmd/wire@latest

docker-up: ## Start services with Docker Compose
	docker-compose up -d

docker-down: ## Stop services with Docker Compose
	docker-compose down

docker-build: ## Build Docker image
	docker build -t boilerplate-go .

docker-run: ## Run Docker container
	docker run -p 8080:8080 boilerplate-go

migrate-up: ## Run database migrations up
	@echo "Add your migration command here"

migrate-down: ## Run database migrations down
	@echo "Add your migration command here"

seed: ## Seed the database
	@echo "Add your seed command here"

lint: ## Run linter
	golangci-lint run

fmt: ## Format code
	go fmt ./...

deps: ## Download dependencies
	go mod download
	go mod tidy

wire: ## Generate Wire dependency injection code
	cd internal/wire && wire

wire-check: ## Check Wire dependency injection setup
	cd internal/wire && wire check

wire-clean: ## Clean generated Wire code and regenerate
	rm -f internal/wire/wire_gen.go
	cd internal/wire && wire
