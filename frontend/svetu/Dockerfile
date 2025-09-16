# syntax=docker/dockerfile:1

# Stage 1: Dependencies
FROM node:20-alpine AS deps
RUN apk add --no-cache libc6-compat
WORKDIR /app

# Copy package files
COPY package.json yarn.lock ./

# Install ALL dependencies including devDependencies (needed for build)
RUN yarn install --frozen-lockfile --production=false --network-timeout 100000

# Stage 2: Builder
FROM node:20-alpine AS builder
WORKDIR /app

# Copy dependencies from previous stage
COPY --from=deps /app/node_modules ./node_modules

# Copy application code
COPY . .

# Accept build arguments
ARG NEXT_PUBLIC_API_URL=http://localhost:3000
ARG NEXT_PUBLIC_MINIO_URL=http://localhost:9000
ARG NEXT_PUBLIC_IMAGE_HOSTS=http:localhost:9000,https:svetu.rs:443,http:localhost:3000
ARG NEXT_PUBLIC_IMAGE_PATH_PATTERN=/listings/**
ARG NEXT_PUBLIC_WEBSOCKET_URL=ws://localhost:3000
ARG NEXT_PUBLIC_ENABLE_PAYMENTS=false
ARG NEXT_PUBLIC_MAPBOX_ACCESS_TOKEN=dummy_token_for_build
ARG NEXT_PUBLIC_MAPBOX_TOKEN=dummy_token_for_build
ARG NEXT_PUBLIC_AUTH_SERVICE_URL=http://localhost:28080

# Set build-time environment variables from arguments
ENV NEXT_PUBLIC_API_URL=$NEXT_PUBLIC_API_URL
ENV NEXT_PUBLIC_MINIO_URL=$NEXT_PUBLIC_MINIO_URL
ENV NEXT_PUBLIC_IMAGE_HOSTS=$NEXT_PUBLIC_IMAGE_HOSTS
ENV NEXT_PUBLIC_IMAGE_PATH_PATTERN=$NEXT_PUBLIC_IMAGE_PATH_PATTERN
ENV NEXT_PUBLIC_WEBSOCKET_URL=$NEXT_PUBLIC_WEBSOCKET_URL
ENV NEXT_PUBLIC_ENABLE_PAYMENTS=$NEXT_PUBLIC_ENABLE_PAYMENTS
ENV NEXT_PUBLIC_MAPBOX_ACCESS_TOKEN=$NEXT_PUBLIC_MAPBOX_ACCESS_TOKEN
ENV NEXT_PUBLIC_MAPBOX_TOKEN=$NEXT_PUBLIC_MAPBOX_TOKEN
ENV NEXT_PUBLIC_AUTH_SERVICE_URL=$NEXT_PUBLIC_AUTH_SERVICE_URL

# Build the application
RUN yarn build

# Stage 3: Runner
FROM node:20-alpine AS runner
WORKDIR /app

# Set production environment
ENV NODE_ENV=production

# Create non-root user
RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nextjs

# Copy necessary files from builder (standalone)
COPY --from=builder /app/.next/standalone ./
COPY --from=builder /app/.next/static ./.next/static
COPY --from=builder /app/public ./public

# Copy i18n messages
COPY --from=builder /app/src/messages ./src/messages

# Copy package.json for runtime dependencies info
COPY --from=builder /app/package.json ./package.json

# Create directory for runtime environment script
RUN mkdir -p .next/static/runtime

# Change ownership
RUN chown -R nextjs:nodejs /app

# Switch to non-root user
USER nextjs

# Expose port
EXPOSE 3000

# Set default environment variables for runtime
ENV PORT=3000
ENV HOSTNAME="0.0.0.0"

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=30s --retries=3 \
  CMD node -e "fetch('http://localhost:3000/api/health').then(r => process.exit(r.ok ? 0 : 1)).catch(() => process.exit(1))"

# Start the application
CMD ["node", "server.js"]