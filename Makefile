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

docker-full:
	docker-compose down
	docker-compose up --build -d
	docker-compose logs -f

# Install migration tool
install-migrate:
	go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Install dependencies
dep:
	go mod download

# Install all tools
install-tools:
	go install github.com/golang/mock/mockgen@latest
	go install github.com/swaggo/swag/cmd/swag@latest

generate-mocks:
	mockgen -source=internal/domain/translation/repository.go -destination=internal/mocks/translation/repository.go
	mockgen -source=internal/domain/user/repository.go -destination=internal/mocks/user/repository.go
	mockgen -source=internal/domain/book/repository.go -destination=internal/mocks/book/repository.go

generate: swagger generate-mocks

# Install dependencies
dep:
	go mod download

# Install all tools and generate mocks
install: install-tools generate-mocks

.PHONY: build run test clean swagger migrate-up migrate-down docker-build docker-up docker-down install-tools generate-mocks generate dep install
tools: install-migrate
	go install github.com/swaggo/swag/cmd/swag@latest

# Setup development environment
setup: tools dep
	@echo "Development environment is ready!"
