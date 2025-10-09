#!/bin/bash
echo "=== Checking old references ==="
echo "Backend marketplace: $(grep -r "marketplace" backend/internal/proj --include="*.go" | wc -l)"
echo "Backend storefronts: $(grep -r "storefronts" backend/internal/proj --include="*.go" | wc -l)"
echo "Frontend marketplace: $(grep -r "marketplace" frontend/svetu/src --include="*.ts*" | wc -l)"
echo "Frontend storefronts: $(grep -r "storefronts" frontend/svetu/src --include="*.ts*" | wc -l)"
