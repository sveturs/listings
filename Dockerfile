# Stage 1: Build stage
FROM golang:1.25-alpine AS builder

# Accept GitHub token as build arg (optional - will work without for public repos)
ARG GITHUB_TOKEN

# Install build dependencies
RUN apk add --no-cache git make gcc musl-dev

# Set working directory
WORKDIR /build

# Set GOPRIVATE for sveturs modules
ENV GOPRIVATE=github.com/sveturs/*

# Configure Git credentials if token provided (for private repos)
RUN if [ -n "$GITHUB_TOKEN" ]; then \
        git config --global url."https://${GITHUB_TOKEN}@github.com/".insteadOf "https://github.com/"; \
    fi

# Copy go mod and vendor directory for offline build
COPY go.mod go.sum ./
COPY vendor/ ./vendor/

# Copy source code
COPY . .

# Build the application using vendored dependencies (no network access needed)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor \
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
