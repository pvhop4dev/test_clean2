# Project Architecture Documentation

## Clean Architecture Overview
This project follows the Clean Architecture principles, which separates the application into distinct layers with clear dependencies and responsibilities.

## Project Structure
```
internal/
├── domain/         # Core business logic and entities
├── application/    # Use cases and business rules
├── delivery/       # HTTP/GRPC handlers and presentation layer
├── middleware/     # Cross-cutting concerns
└── pkg/           # Shared utilities and helpers

pkg/               # External packages and protocols
└── pb/            # Protocol Buffers definitions
```

## Layer Dependencies
- Domain layer: No external dependencies
- Application layer: Depends only on Domain layer
- Delivery layer: Depends on Application layer
- Middleware: Depends on Application layer

## Key Components

### Domain Layer
Contains core business entities and interfaces:
- Entities: Business objects
- Interfaces: Contracts for repositories and external services
- Value Objects: Immutable objects representing domain concepts

### Application Layer
Contains use cases and business rules:
- Use Cases: Implement business rules
- Interfaces: Contracts for external services
- DTOs: Data Transfer Objects

### Delivery Layer
Handles HTTP/GRPC requests:
- Handlers: Request/response processing
- Routers: Request routing
- Presenters: Response formatting

### Middleware
Handles cross-cutting concerns:
- Authentication
- Authorization
- Logging
- Rate limiting

## Dependency Injection
The project uses constructor injection for all dependencies:
```go
// Example of dependency injection
func NewTranslationUsecase(repo translation.Repository) TranslationUsecase {
    return &translationUsecase{
        repo: repo,
    }
}
```

## Testing
- Unit tests: Test individual components in isolation
- Integration tests: Test component interactions
- Mocks: Use generated mocks for external dependencies

## Best Practices
1. Keep domain layer free of external dependencies
2. Use interfaces for all external dependencies
3. Follow SOLID principles
4. Keep business logic in application layer
5. Use dependency injection for all components

## Future Improvements
1. Add more use cases in application layer
2. Implement more middleware components
3. Add comprehensive testing
4. Improve error handling
5. Add more documentation
