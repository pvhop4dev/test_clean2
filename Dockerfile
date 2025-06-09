# Build stage
FROM golang:1.24.2-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Install Swag and generate Swagger docs
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init -g cmd/api/main.go

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/main ./cmd/api

# Final stage
FROM alpine:3.18

WORKDIR /app

# Install SSL certificates and timezone data
RUN apk --no-cache add ca-certificates tzdata

# Copy the binary from builder
COPY --from=builder /go/bin/main /app/
# Copy config files
COPY --from=builder /app/.env .
# Copy Swagger files
COPY --from=builder /app/docs ./docs

# Expose port
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
