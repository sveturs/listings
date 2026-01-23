# Primary Image Fix Analysis

**Date:** 2026-01-11
**Status:** RESOLVED

## Problem Description

When changing the primary image for a product on vondi.rs:
- Works on product detail page (shows correct image)
- Does NOT work on main page/listing cards (shows old image)

## Root Cause Analysis

### Issue 1: Missing OpenSearch Reindexing

**Location:** `internal/service/listings/service.go:SetProductImagePrimary`

**Problem:** The function updated the database but did NOT reindex the product in OpenSearch.

**Fix Applied (PR #23):**
```go
// After DB update, trigger async reindex
if s.indexer != nil && productID > 0 {
    go func() {
        // ... reindex product in OpenSearch
    }()
}
```

### Issue 2: Deployment Did Not Restart Pods

**Root Cause:** The deployment workflow used `kubectl set image` with `latest` tag. When the tag name doesn't change, Kubernetes doesn't trigger a pod restart even if the image content has changed.

**Timeline:**
- Pods created: 16:33:54 UTC (4 hours before deployment)
- Deployment ran: 20:21:30 UTC
- Pods NOT restarted (same tag name `latest`)

**Evidence from logs:**
```
SetProductImagePrimary called
SetProductImagePrimary completed
# NO reindexing message - old code still running
```

**Immediate Fix:** Manual pod restart
```bash
kubectl rollout restart deployment/listings-service -n production
```

**Permanent Fix (PR #25):** Changed workflow to use SHA tag
```yaml
# Before (broken):
$HARBOR_REGISTRY/$IMAGE_NAME:latest

# After (fixed):
SHA="${{ steps.git_sha.outputs.sha }}"
$HARBOR_REGISTRY/$IMAGE_NAME:$SHA
```

## Verification Steps

1. Check pod age after deployment:
```bash
kubectl get pods -n production -l app=listings-service --no-headers
```

2. Check for reindexing in logs:
```bash
kubectl logs -n production deployment/listings-service | grep "reindexed.*primary"
```

3. Test on vondi.rs:
- Change primary image on product
- Verify image updates on main page listing cards

## Lessons Learned

1. **Always use unique image tags for K8s deployments** - SHA or build number, never just `latest`
2. **Verify pod restart after deployment** - check pod age matches deployment time
3. **Check logs for expected behavior** - missing log entries indicate old code running
