# CI/CD Setup - Orders Microservice

**Version:** 1.0.0
**Last Updated:** 2025-11-14

Complete guide for Continuous Integration and Continuous Deployment pipeline for the Orders Microservice.

---

## Table of Contents

1. [Overview](#overview)
2. [Pipeline Architecture](#pipeline-architecture)
3. [GitHub Actions Workflows](#github-actions-workflows)
4. [Test Stages](#test-stages)
5. [Deployment](#deployment)
6. [Monitoring](#monitoring)
7. [Troubleshooting](#troubleshooting)

---

## Overview

The CI/CD pipeline for Orders Microservice ensures:

- **Code Quality**: Automated linting and formatting checks
- **Test Coverage**: Unit, integration, and E2E tests with >80% coverage
- **Performance**: Benchmark tests and regression detection
- **Security**: Vulnerability scanning and dependency checks
- **Automated Deployment**: Push-to-deploy for dev and staging environments
- **Rollback Capability**: Quick rollback to previous version on failures

### Pipeline Goals

1. **Fast Feedback**: Developers get feedback within 10 minutes
2. **High Confidence**: >80% test coverage, all tests passing
3. **Zero Downtime**: Blue-green deployments
4. **Observability**: Metrics, logs, and traces for all deployments

---

## Pipeline Architecture

```
┌─────────────┐
│  Git Push   │
└──────┬──────┘
       │
       v
┌─────────────────────────────────────────────────┐
│              GitHub Actions                      │
├─────────────────────────────────────────────────┤
│                                                  │
│  ┌──────────┐  ┌────────────┐  ┌─────────────┐ │
│  │   Lint   │─>│ Unit Tests │─>│Integration  │ │
│  │  (2min)  │  │   (5min)   │  │Tests (8min) │ │
│  └──────────┘  └────────────┘  └─────────────┘ │
│                                                  │
│       │                 │               │        │
│       └─────────────────┴───────────────┘        │
│                      │                           │
│                      v                           │
│         ┌────────────────────────┐               │
│         │   Build Docker Image   │               │
│         │       (3min)           │               │
│         └────────────┬───────────┘               │
│                      │                           │
│                      v                           │
│         ┌────────────────────────┐               │
│         │ Performance Tests      │               │
│         │ (main branch only)     │               │
│         └────────────┬───────────┘               │
│                      │                           │
│                      v                           │
│         ┌────────────────────────┐               │
│         │  Deploy to Dev/Staging │               │
│         │  (auto on main/develop)│               │
│         └────────────────────────┘               │
└─────────────────────────────────────────────────┘
```

### Workflow Files

- `.github/workflows/ci.yml` - Main CI pipeline (runs on all branches)
- `.github/workflows/orders-service-ci.yml` - Orders-specific CI/CD
- `.github/workflows/coverage.yml` - Coverage reporting

---

## GitHub Actions Workflows

### Main CI Workflow

**File:** `.github/workflows/orders-service-ci.yml`

**Triggers:**
- Push to `main`, `develop`, `feature/orders-*` branches
- Pull requests to `main`, `develop`

**Jobs:**

#### 1. Lint
```yaml
- golangci-lint with custom config
- Timeout: 5 minutes
- Fail fast on critical issues
```

#### 2. Unit Tests
```yaml
- Run unit tests with race detector
- Generate coverage report
- Fail if coverage < 80%
- Upload to Codecov
```

#### 3. Integration Tests (E2E)
```yaml
- Start PostgreSQL, Redis, OpenSearch (services)
- Run database migrations
- Execute E2E tests
- Upload test results
```

#### 4. Performance Tests
```yaml
- Run Go benchmarks
- Compare with baseline
- Fail if regression > 20%
- Upload results (main branch only)
```

#### 5. Build
```yaml
- Build Go binary
- Upload build artifact
```

#### 6. Docker
```yaml
- Build Docker image
- Push to registry (main/develop only)
- Tag with branch-SHA
```

#### 7. Deploy
```yaml
- Deploy to dev (develop branch)
- Deploy to staging (main branch)
- Smoke tests after deployment
```

### Coverage Workflow

**File:** `.github/workflows/coverage.yml`

**Purpose:** Generate detailed coverage reports and badges

**Features:**
- Line-by-line coverage
- Coverage diff for PRs
- Comment on PR with coverage report
- Update coverage badge

---

## Test Stages

### Stage 1: Lint (2 minutes)

**Tools:**
- golangci-lint (aggregates 40+ linters)

**Configuration:** `.golangci.yml`

**Checks:**
- Code formatting (gofmt, goimports)
- Code quality (govet, staticcheck)
- Complexity (gocyclo)
- Security (gosec)
- Spelling (misspell)

**Example:**
```bash
golangci-lint run --timeout=5m
```

**Success Criteria:**
- Zero critical issues
- Warnings are acceptable but should be addressed

---

### Stage 2: Unit Tests (5 minutes)

**Coverage Target:** >80%

**Environment:**
- PostgreSQL 15 (service container)
- Redis 7 (service container)
- Test database: `listings_db_test`

**Commands:**
```bash
go test -v -race -coverprofile=coverage.txt -covermode=atomic -short \
  ./internal/... \
  ./pkg/...
```

**Flags:**
- `-race`: Detect race conditions
- `-short`: Skip slow integration tests
- `-coverprofile`: Generate coverage report

**Success Criteria:**
- All tests pass
- Coverage ≥ 80%
- No race conditions detected

---

### Stage 3: Integration Tests (8 minutes)

**Environment:**
- Full stack with testcontainers
- PostgreSQL, Redis, OpenSearch
- Real gRPC server

**Tests:**
- End-to-end flows (cart → order → payment)
- gRPC handler tests
- Database transaction tests
- Cache integration tests

**Commands:**
```bash
go test -v -timeout=10m -tags=integration \
  ./test/integration/... \
  ./internal/transport/grpc/...
```

**Success Criteria:**
- All E2E scenarios pass
- No flaky tests
- Proper cleanup (no resource leaks)

---

### Stage 4: Performance Tests (10 minutes, main only)

**Purpose:** Detect performance regressions before deployment

**Tests:**
- Go benchmarks (`BenchmarkAddToCart`, etc.)
- Load tests (skipped in CI, run manually)

**Commands:**
```bash
go test -bench=. -benchmem -benchtime=10s \
  -run=^$ \
  ./tests/performance/...
```

**Regression Detection:**
```bash
# Download baseline from previous run
# Compare metrics
# Fail if:
# - Latency increased > 20%
# - Memory allocations increased > 30%
# - Error rate increased
```

**Success Criteria:**
- All benchmarks complete
- No regressions > 20%
- Memory allocations stable

---

### Stage 5: Build (3 minutes)

**Artifacts:**
1. **Binary**: `bin/listings-service`
2. **Docker Image**: `sveturs/listings-service:branch-sha`

**Build Commands:**
```bash
# Go binary
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
  go build -ldflags="-s -w" \
  -o bin/listings-service \
  ./cmd/server/main.go

# Docker image
docker build -t sveturs/listings-service:$SHA .
docker push sveturs/listings-service:$SHA
```

**Optimizations:**
- Multi-stage Docker build
- Layer caching
- Minimal base image (alpine)

---

### Stage 6: Deploy (5 minutes)

**Environments:**

| Branch | Environment | Auto-Deploy | URL |
|--------|-------------|-------------|-----|
| `develop` | Development | Yes | https://dev.vondi.rs |
| `main` | Staging | Yes | https://staging.vondi.rs |
| `main` (manual) | Production | No | https://api.vondi.rs |

**Deployment Strategy:** Blue-Green

**Steps:**
1. Health check existing deployment
2. Deploy new version (green)
3. Wait for health check (30s)
4. Run smoke tests
5. Switch traffic to green
6. Keep blue for 5 minutes (rollback window)
7. Terminate blue

**Rollback:**
```bash
# Automatic rollback if:
# - Health check fails
# - Smoke tests fail
# - Error rate > 5% in first 5 minutes

# Manual rollback
kubectl rollout undo deployment/listings-service
```

---

## Monitoring

### Deployment Metrics

**Tracked:**
- Deployment duration
- Success/failure rate
- Rollback frequency
- Time to deploy
- Time to rollback

**Dashboard:** Grafana (Deployments Dashboard)

### Post-Deployment Checks

**Automated Smoke Tests:**
```bash
# 1. Health check
grpcurl -plaintext api.vondi.rs:50052 grpc.health.v1.Health/Check

# 2. Basic operations
grpcurl -plaintext -d '{"user_id": 1}' \
  api.vondi.rs:50052 listings.v1.OrdersService/GetCart

# 3. Check metrics
curl https://api.vondi.rs/metrics | grep grpc_server_started_total
```

**Manual Verification:**
- Review logs for errors
- Check Grafana for anomalies
- Verify database migrations applied
- Test critical user flows

---

## Secrets Management

### Required Secrets

**GitHub Secrets:**
```
CODECOV_TOKEN              # Coverage reporting
DOCKER_USERNAME            # Docker registry
DOCKER_PASSWORD            # Docker registry
DATABASE_URL_DEV           # Dev database connection
DATABASE_URL_STAGING       # Staging database connection
DEPLOY_SSH_KEY             # Deployment SSH key
SLACK_WEBHOOK_URL          # Notifications
```

**Setup:**
```bash
# Add secret to GitHub
gh secret set CODECOV_TOKEN --body="xxx"

# Verify secrets
gh secret list
```

---

## Troubleshooting

### CI Failures

#### Lint Failures

**Symptom:** golangci-lint job fails

**Diagnosis:**
```bash
# Run locally
cd /p/github.com/sveturs/listings
golangci-lint run --timeout=5m
```

**Solutions:**
- Fix formatting: `gofmt -w .`
- Fix imports: `goimports -w .`
- Address issues reported by linter

#### Test Failures

**Symptom:** Unit or integration tests fail

**Diagnosis:**
```bash
# Run locally
go test -v ./...

# Run with verbose output
go test -v -run TestFailingTest ./path/to/package
```

**Common Issues:**
- Flaky tests (timing issues) → Add proper waits
- Database state pollution → Improve test cleanup
- Race conditions → Use mutexes or channels

#### Coverage Below Threshold

**Symptom:** Coverage check fails (< 80%)

**Diagnosis:**
```bash
# Generate coverage report
go test -coverprofile=coverage.out ./...

# View in browser
go tool cover -html=coverage.out

# Check specific package
go test -cover ./internal/domain
```

**Solutions:**
- Add tests for untested code
- Remove dead code
- Exclude generated code from coverage

#### Performance Regression

**Symptom:** Benchmark comparison fails

**Diagnosis:**
```bash
# Run benchmarks locally
go test -bench=. -benchmem ./tests/performance/

# Compare with baseline
benchstat baseline.txt current.txt
```

**Solutions:**
- Profile with pprof: `go test -bench=. -cpuprofile=cpu.prof`
- Identify hot spots: `go tool pprof cpu.prof`
- Optimize critical paths

#### Docker Build Failure

**Symptom:** Docker build step fails

**Diagnosis:**
```bash
# Build locally
docker build -t test .

# Check build logs
docker build --progress=plain -t test .
```

**Common Issues:**
- Missing dependencies → Update `go.mod`
- Build cache issues → Clear cache
- Base image issues → Update base image

#### Deployment Failure

**Symptom:** Deployment job fails

**Diagnosis:**
```bash
# Check deployment logs
kubectl logs -n production deploy/listings-service

# Check events
kubectl get events -n production --sort-by='.lastTimestamp'

# Check pod status
kubectl describe pod -n production <pod-name>
```

**Common Issues:**
- Health check fails → Check application logs
- Image pull error → Verify image exists
- Resource limits → Adjust CPU/memory requests

---

## Best Practices

### Writing CI-Friendly Tests

1. **Fast Tests**: Unit tests should complete in <1 second
2. **Isolated Tests**: No shared state between tests
3. **Deterministic**: Same input → same output (no randomness)
4. **Idempotent**: Can run multiple times without side effects
5. **Proper Cleanup**: Teardown resources in `defer`

### Optimizing CI Time

1. **Parallel Jobs**: Run independent jobs in parallel
2. **Caching**: Cache dependencies (Go modules, Docker layers)
3. **Fail Fast**: Stop on first failure
4. **Incremental Tests**: Run only affected tests (future)
5. **Pre-merge Checks**: Run expensive tests only on main/develop

### Security

1. **Secrets**: Never commit secrets, use GitHub Secrets
2. **Scanning**: Run vulnerability scans on dependencies
3. **Image Scanning**: Scan Docker images for CVEs
4. **RBAC**: Limit deployment permissions
5. **Audit Logs**: Track all deployments and rollbacks

---

## Notifications

### Slack Integration

**Channels:**
- `#ci-builds` - All CI runs
- `#deployments` - Deployment notifications
- `#alerts` - Failures and rollbacks

**Message Format:**
```
🚀 Deployment Started
Environment: Production
Version: v1.2.3
Commit: abc123f
Author: @developer
Branch: main
```

### Email Notifications

**Recipients:**
- Build failures → Author of commit
- Deployment failures → DevOps team
- Security alerts → Security team

---

## Metrics & SLOs

### CI/CD SLOs

| Metric | Target | Current |
|--------|--------|---------|
| **Build Success Rate** | >95% | 97% |
| **Deployment Success Rate** | >98% | 99% |
| **Mean Time to Deploy** | <15min | 12min |
| **Mean Time to Rollback** | <5min | 3min |
| **Test Coverage** | >80% | 85% |

### Tracking

**Dashboard:** Grafana (CI/CD Metrics)

**Queries:**
```promql
# Deployment success rate
sum(rate(deployments_total{status="success"}[1d]))
/
sum(rate(deployments_total[1d]))

# Mean time to deploy
histogram_quantile(0.50, rate(deployment_duration_seconds_bucket[1d]))
```

---

## References

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Docker Best Practices](https://docs.docker.com/develop/dev-best-practices/)
- [Kubernetes Deployment Strategies](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/)
- [CI/CD Best Practices](https://www.atlassian.com/continuous-delivery/principles/continuous-integration-vs-delivery-vs-deployment)

---

**Document Owner:** DevOps Team
**Last Review:** 2025-11-14
**Next Review:** 2025-12-14
