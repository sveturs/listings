# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux go build -o reindex ./cmd/reindex

# Final stage
FROM alpine:latest

WORKDIR /app

# Create non-root user
RUN adduser -D -g '' appuser

# Copy binaries from builder
COPY --from=builder /app/main /app/main
COPY --from=builder /app/reindex /app/reindex

# Create uploads directory and set permissions
RUN mkdir -p uploads && chown -R appuser:appuser uploads

# Switch to non-root user
USER appuser

EXPOSE 3000

CMD ["./main"]
