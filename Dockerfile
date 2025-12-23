# Listings Service Dockerfile
# Build optimized multi-stage image for production

# Stage 1: Builder
FROM golang:1.25-alpine AS builder

# Accept GitHub token as build arg (optional for private repos)
ARG GITHUB_TOKEN

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata make gcc musl-dev

# Set working directory
WORKDIR /build

# Set GOPRIVATE for vondi-global modules
ENV GOPRIVATE=github.com/vondi-global/*

# Configure Git credentials if token provided
RUN if [ -n "$GITHUB_TOKEN" ]; then \
        git config --global url."https://${GITHUB_TOKEN}@github.com/".insteadOf "https://github.com/"; \
    fi

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary
# CGO_ENABLED=0 for static binary
# -ldflags="-s -w" to strip debug info
# -mod=mod to ignore vendor directory
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -mod=mod \
    -ldflags="-s -w -X main.Version=$(git describe --tags --always --dirty 2>/dev/null || echo 'unknown') -X main.BuildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
    -o /build/bin/listings-service \
    ./cmd/server/main.go

# Create keys directory if it doesn't exist
RUN mkdir -p /build/keys

# Stage 2: Runtime
FROM alpine:3.19

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1000 vondi && \
    adduser -D -u 1000 -G vondi vondi

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/bin/listings-service /app/listings-service

# Copy migrations
COPY --from=builder /build/migrations /app/migrations

# Copy auth service public key directory (may be empty)
COPY --from=builder /build/keys /app/keys

# Change ownership
RUN chown -R vondi:vondi /app

# Switch to non-root user
USER vondi

# Expose ports
# gRPC port
EXPOSE 50053
# HTTP port
EXPOSE 8084
# Metrics port
EXPOSE 9093

# Health check
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD ["/app/listings-service", "healthcheck"] || wget --no-verbose --tries=1 --spider http://localhost:8084/health || exit 1

# Run the application
ENTRYPOINT ["/app/listings-service"]
CMD ["serve"]
