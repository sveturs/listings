# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# GitHub token for private repos
ARG GITHUB_TOKEN
RUN git config --global url."https://${GITHUB_TOKEN}@github.com/".insteadOf "https://github.com/"

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Аргументы для версии
ARG GIT_COMMIT=unknown
ARG BUILD_TIME=unknown

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-X main.gitCommit=${GIT_COMMIT} -X main.buildTime=${BUILD_TIME}" \
      -o main ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-X main.gitCommit=${GIT_COMMIT} -X main.buildTime=${BUILD_TIME}" \
      -o reindex ./cmd/reindex
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-X main.gitCommit=${GIT_COMMIT} -X main.buildTime=${BUILD_TIME}" \
      -o migrator ./cmd/migrator

# Final stage
FROM alpine:latest

WORKDIR /app

# Create non-root user
RUN adduser -D -g '' appuser

# Copy binaries from builder
COPY --from=builder /app/main /app/main
COPY --from=builder /app/reindex /app/reindex
COPY --from=builder /app/migrator /app/migrator

# Copy migrations
COPY --from=builder /app/migrations /app/migrations
COPY --from=builder /app/fixtures /app/fixtures

# Create uploads directory and set permissions
RUN mkdir -p uploads && chown -R appuser:appuser uploads

# Switch to non-root user
USER appuser

EXPOSE 3000

CMD ["./main"]
