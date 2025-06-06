.PHONY: build run test clean swagger migrate-up migrate-down docker-build docker-up docker-down

# Build the application
build:
	go build -o bin/api cmd/api/main.go

# Run the application
run:
	go run cmd/api/main.go

# Run tests
test:
	go test -v ./...

# Clean build files
clean:
	rm -rf bin/

# Generate Swagger documentation
swagger:
	swag init -g cmd/api/main.go

# Database migrations
migrate-up:
	migrate -path migrations -database "mysql://root:password@tcp(localhost:3306)/clean_arch_go" -verbose up

migrate-down:
	migrate -path migrations -database "mysql://root:password@tcp(localhost:3306)/clean_arch_go" -verbose down

# Docker commands
docker-build:
	docker-compose build

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

# Install migration tool
install-migrate:
	go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Install dependencies
dep:
	go mod download

# Install all tools
tools: install-migrate
	go install github.com/swaggo/swag/cmd/swag@latest

# Setup development environment
setup: tools dep
	@echo "Development environment is ready!"
