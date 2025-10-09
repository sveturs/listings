#!/bin/bash
set -e

REPO_ROOT="/data/hostel-booking-system"
cd "$REPO_ROOT"

echo "🔄 Starting automatic renaming..."

# Backend: Go modules
echo "📦 Renaming Go modules..."
git mv backend/internal/proj/marketplace backend/internal/proj/c2c
git mv backend/internal/proj/storefronts backend/internal/proj/b2c

# Backend: Update imports
echo "📝 Updating Go imports..."
/usr/bin/find backend -name "*.go" -type f -exec sed -i \
  's|internal/proj/marketplace|internal/proj/c2c|g' {} +
/usr/bin/find backend -name "*.go" -type f -exec sed -i \
  's|internal/proj/storefronts|internal/proj/b2c|g' {} +

# Frontend: Components
echo "🎨 Renaming frontend components..."
git mv frontend/svetu/src/components/marketplace frontend/svetu/src/components/c2c || true
git mv frontend/svetu/src/components/storefronts frontend/svetu/src/components/b2c || true

# Frontend: Routes
echo "🛤️  Renaming frontend routes..."
git mv frontend/svetu/src/app/\[locale\]/marketplace frontend/svetu/src/app/\[locale\]/c2c || true
git mv frontend/svetu/src/app/\[locale\]/storefronts frontend/svetu/src/app/\[locale\]/b2c || true

# Frontend: i18n
echo "🌐 Renaming translation files..."
for lang in en ru sr; do
  git mv frontend/svetu/src/messages/$lang/marketplace.json \
         frontend/svetu/src/messages/$lang/c2c.json || true
  git mv frontend/svetu/src/messages/$lang/storefronts.json \
         frontend/svetu/src/messages/$lang/b2c.json || true
done

echo "✅ Automatic renaming complete!"
echo "⚠️  Run validation script to check for remaining references"
