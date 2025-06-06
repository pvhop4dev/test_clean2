# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go mod and sum files
COPY go.mod go.sum ./


# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .


# Generate Swagger docs
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init -g cmd/api/main.go

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api

# Final stage
FROM alpine:latest

WORKDIR /app

# Install required packages
RUN apk --no-cache add ca-certificates

# Copy the binary from builder
COPY --from=builder /app/main .
# Copy config files
COPY .env .

# Expose port
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
