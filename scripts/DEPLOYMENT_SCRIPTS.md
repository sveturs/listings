# Production Deployment Scripts

Comprehensive deployment automation with Blue-Green strategy and zero downtime.

## Quick Start

```bash
# 1. Configure
cp .env.deploy.example .env.deploy
nano .env.deploy

# 2. Validate
./validate-deployment.sh

# 3. Deploy
./deploy-to-prod.sh
```

## Scripts

### `deploy-to-prod.sh` - Main Deployment
Blue-Green deployment with canary rollout.

```bash
./deploy-to-prod.sh              # Full deployment
./deploy-to-prod.sh --dry-run    # Test run
./deploy-to-prod.sh --verbose    # Debug mode
./deploy-to-prod.sh --skip-tests # Emergency (not recommended)
```

Duration: 25-35 minutes

### `rollback-prod.sh` - Emergency Rollback
Instant rollback to previous stable version.

```bash
./rollback-prod.sh --reason "high error rate"
./rollback-prod.sh --reason "migration failed" --restore-db
```

Duration: 1-10 minutes

### `smoke-tests.sh` - Automated Testing
Validates critical endpoints.

```bash
./smoke-tests.sh --host production.example.com --port 80
./smoke-tests.sh --json          # CI/CD mode
./smoke-tests.sh --verbose       # Debug mode
```

Duration: 2-3 minutes

### `validate-deployment.sh` - Pre-flight Checks
Validates environment before deployment.

```bash
./validate-deployment.sh
./validate-deployment.sh --verbose
```

Duration: 30 seconds

### `traffic-split.sh` - Traffic Control
Manages Blue-Green traffic distribution.

```bash
./traffic-split.sh --green-weight 0    # 100% Blue
./traffic-split.sh --green-weight 50   # 50/50 split
./traffic-split.sh --green-weight 100  # 100% Green
```

Duration: < 1 minute

### `deployment-report.sh` - Report Generator
Creates comprehensive deployment reports.

```bash
./deployment-report.sh --deployment-id deploy-20250105-143022
```

Duration: < 10 seconds

## Documentation

- [Deployment Guide](../docs/DEPLOYMENT.md) - Full deployment procedures
- [Rollback Guide](../docs/ROLLBACK.md) - Emergency procedures
- [Configuration](.env.deploy.example) - Environment settings

## Support

**Slack:** `#deployments`
**On-Call:** See `.env.deploy`
