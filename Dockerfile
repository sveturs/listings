# Stage 1: Build stage
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make gcc musl-dev

# Set working directory
WORKDIR /build

# Copy go mod files first for better layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -X main.Version=$(git describe --tags --always --dirty) -X main.BuildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
    -o /build/bin/listings-service \
    ./cmd/server

# Stage 2: Runtime stage
FROM alpine:3.19

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1000 listings && \
    adduser -D -u 1000 -G listings listings

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/bin/listings-service /app/listings-service

# Copy migrations (if needed for embedded migrations)
COPY --from=builder /build/migrations /app/migrations

# Change ownership
RUN chown -R listings:listings /app

# Switch to non-root user
USER listings

# Expose ports
EXPOSE 50053 8086 9093

# Health check
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD ["/app/listings-service", "healthcheck"]

# Run the application
ENTRYPOINT ["/app/listings-service"]
CMD ["serve"]
