# Listings Service - Production Deployment

## Automated Deployment

**Полностью автоматический CD pipeline!**

При merge PR в `main`:
1. Docker build → Harbor registry
2. Kubernetes rollout → terra.vondi.rs
3. Health check → https://listings.vondi.rs/health
4. Auto-rollback при ошибке

---

## Quick Reference

| Parameter | Value |
|-----------|-------|
| **Image** | `registry.vondi.rs/library/listings` |
| **Deployment** | `listings-service` |
| **Namespace** | `production` |
| **Health URL** | https://listings.vondi.rs/health |
| **Workflow** | `.github/workflows/deploy-production.yml` |

---

## GitHub Secrets Required

```bash
gh secret set HARBOR_USERNAME -R vondi-global/listings
gh secret set HARBOR_PASSWORD -R vondi-global/listings
gh secret set K8S_SSH_KEY -R vondi-global/listings < ~/.ssh/id_ed25519
```

---

## Manual Deployment

### Via GitHub Actions

```bash
gh workflow run deploy-production.yml --repo vondi-global/listings
```

### Via Kubernetes

```bash
# Update to latest
ssh agrouser@terra.vondi.rs "sudo kubectl set image deployment/listings-service \
  listings-service=registry.vondi.rs/library/listings:latest -n production"

# Update to specific version (rollback)
ssh agrouser@terra.vondi.rs "sudo kubectl set image deployment/listings-service \
  listings-service=registry.vondi.rs/library/listings:abc1234 -n production"
```

---

## Monitoring

### Logs

```bash
ssh agrouser@terra.vondi.rs "sudo kubectl logs -f deployment/listings-service -n production"
```

### Health Check

```bash
curl https://listings.vondi.rs/health
```

### Metrics

```bash
curl https://listings.vondi.rs/metrics
```

---

## Rollback

```bash
# Auto: При провале health check - автоматический откат
# Manual:
ssh agrouser@terra.vondi.rs "sudo kubectl rollout undo deployment/listings-service -n production"
```

---

## Troubleshooting

### Check Deployment Status

```bash
ssh agrouser@terra.vondi.rs "sudo kubectl get deployments -n production | grep listings"
```

### Check Pods

```bash
ssh agrouser@terra.vondi.rs "sudo kubectl get pods -n production | grep listings"
```

### Describe Pod

```bash
ssh agrouser@terra.vondi.rs "sudo kubectl describe pod <POD_NAME> -n production"
```

---

**Full Documentation:** `/p/github.com/vondi-global/vondi/DEPLOY.md`
