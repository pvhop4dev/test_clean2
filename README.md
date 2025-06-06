# Clean Architecture Go with Gin, gRPC, MySQL, and Redis

This is a clean architecture Go project with Gin web framework, gRPC, MySQL, and Redis. It includes user authentication, rate limiting, and CRUD operations for books.

## Features

- User authentication (register, login)
- JWT-based authentication
- Rate limiting
- CRUD operations for books
- gRPC server (WIP)
- MySQL database
- Redis for caching and rate limiting
- Environment-based configuration

## Prerequisites

- Go 1.16+
- MySQL 5.7+
- Redis 6.0+
- protoc (for gRPC)

## Installation

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd clean-arch-go
   ```

2. Copy the example environment file and update the values:
   ```bash
   cp .env.example .env
   ```

3. Install dependencies:
   ```bash
   go mod download
   ```

4. Set up the database:
   ```sql
   CREATE DATABASE clean_arch_go;
   ```

5. Run migrations (auto-migrate is enabled by default):
   ```bash
   # The application will automatically create tables on startup
   ```

## Running the Application

1. Start Redis:
   ```bash
   redis-server
   ```

2. Run the application:
   ```bash
   go run cmd/api/main.go
   ```

3. The application will be available at `http://localhost:8080`

## API Endpoints

### Authentication

- `POST /api/register` - Register a new user
- `POST /api/login` - Login and get JWT token

### Books (Requires Authentication)

- `GET /api/books` - List all books for the authenticated user
- `POST /api/books` - Create a new book
- `GET /api/books/:id` - Get a book by ID
- `PUT /api/books/:id` - Update a book
- `DELETE /api/books/:id` - Delete a book

## Project Structure

```
.
├── cmd/
│   └── api/
│       └── main.go          # Application entry point
├── internal/
│   ├── domain/             # Core business logic
│   │   ├── entities/        # Business entities
│   │   ├── repository/      # Data access interfaces
│   │   └── service/         # Business logic
│   ├── delivery/           # Delivery mechanisms
│   │   ├── http/           # HTTP handlers
│   │   └── grpc/           # gRPC handlers
│   ├── middleware/          # HTTP middleware
│   └── pkg/                 # Reusable packages
│       ├── config/         # Configuration
│       ├── database/       # Database connection
│       └── redis/          # Redis client
├── migrations/             # Database migrations
├── pkg/                    # External packages
│   └── pb/                 # Generated protobuf files
├── .env                   # Environment variables
├── go.mod                 # Go module definition
└── README.md             # This file
```

## Environment Variables

See `.env.example` for all available environment variables.

## Running Tests

```bash
go test ./...
```

## Building the Application

```bash
go build -o bin/api cmd/api/main.go
```

## Deployment

### Docker

1. Build the Docker image:
   ```bash
   docker build -t clean-arch-go .
   ```

2. Run the container:
   ```bash
   docker run -p 8080:8080 --env-file .env clean-arch-go
   ```

### Kubernetes

See the `deploy/` directory for Kubernetes deployment manifests.

## License

MIT
