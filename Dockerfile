# Build stage
FROM golang:1.24.0-alpine AS builder

# Install git (required for some Go modules)
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Runtime stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Copy configuration files
COPY --from=builder /app/configs ./configs
# Use Docker-specific config as the main config
COPY --from=builder /app/configs/config.docker.yaml ./configs/config.yaml

# Copy migration files
COPY --from=builder /app/migrations ./migrations

# Expose port
EXPOSE 8080

# Command to run the executable
CMD ["./main"] 