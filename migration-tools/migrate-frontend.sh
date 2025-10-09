#!/bin/bash
set -euo pipefail

# Цвета для вывода
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

FRONTEND_DIR="/data/hostel-booking-system/frontend/svetu"

echo -e "${GREEN}=== Frontend Migration: marketplace → c2c, storefronts → b2c ===${NC}\n"

# Функция для безопасного переименования
safe_mv() {
    local src="$1"
    local dst="$2"

    if [ -e "$src" ]; then
        echo -e "${YELLOW}Moving:${NC} $src → $dst"
        mkdir -p "$(dirname "$dst")"
        mv "$src" "$dst"
        return 0
    else
        echo -e "${RED}Warning:${NC} Source not found: $src"
        return 1
    fi
}

cd "$FRONTEND_DIR"

echo -e "\n${GREEN}Phase 1: Переименование директорий${NC}"
echo "================================================"

# Components
safe_mv "src/components/marketplace" "src/components/c2c" || true
safe_mv "src/components/storefronts" "src/components/b2c" || true
safe_mv "src/components/Storefront" "src/components/B2C" || true

# App routes
safe_mv "src/app/[locale]/marketplace" "src/app/[locale]/c2c" || true
safe_mv "src/app/[locale]/storefronts" "src/app/[locale]/b2c" || true
safe_mv "src/app/[locale]/create-storefront" "src/app/[locale]/create-b2c-store" || true
safe_mv "src/app/[locale]/debug-marketplace" "src/app/[locale]/debug-c2c" || true
safe_mv "src/app/[locale]/admin/storefronts" "src/app/[locale]/admin/b2c-stores" || true
safe_mv "src/app/[locale]/admin/storefront-products" "src/app/[locale]/admin/b2c-products" || true
safe_mv "src/app/[locale]/profile/storefronts" "src/app/[locale]/profile/b2c-stores" || true
safe_mv "src/app/[locale]/examples/ideal-marketplace" "src/app/[locale]/examples/ideal-c2c" || true
safe_mv "src/app/[locale]/examples/storefront-dashboard" "src/app/[locale]/examples/b2c-dashboard" || true
safe_mv "src/app/[locale]/examples/storefront-import" "src/app/[locale]/examples/b2c-import" || true

echo -e "\n${GREEN}Phase 2: Переименование файлов${NC}"
echo "================================================"

# Services
safe_mv "src/services/marketplace.ts" "src/services/c2c.ts" || true
safe_mv "src/services/marketplaceApi.ts" "src/services/c2cApi.ts" || true
safe_mv "src/services/marketplaceOrders.ts" "src/services/c2cOrders.ts" || true
safe_mv "src/services/storefrontApi.ts" "src/services/b2cStoreApi.ts" || true
safe_mv "src/services/storefrontProducts.ts" "src/services/b2cProducts.ts" || true
safe_mv "src/services/ai/storefronts.service.ts" "src/services/ai/b2c.service.ts" || true

# Types
safe_mv "src/types/marketplace.ts" "src/types/c2c.ts" || true
safe_mv "src/types/storefront.ts" "src/types/b2c.ts" || true

# Store slices
safe_mv "src/store/slices/storefrontSlice.ts" "src/store/slices/b2cStoreSlice.ts" || true

# Hooks
safe_mv "src/hooks/useStorefronts.ts" "src/hooks/useB2CStores.ts" || true

# Contexts
safe_mv "src/contexts/CreateStorefrontContext.tsx" "src/contexts/CreateB2CStoreContext.tsx" || true

# Components
safe_mv "src/components/MarketplaceCard.tsx" "src/components/C2CCard.tsx" || true

# i18n files (all locales)
for locale in en ru sr; do
    safe_mv "src/messages/$locale/marketplace.json" "src/messages/$locale/c2c.json" || true
    safe_mv "src/messages/$locale/marketplace.home.json" "src/messages/$locale/c2c.home.json" || true
    safe_mv "src/messages/$locale/storefronts.json" "src/messages/$locale/b2c.json" || true
    safe_mv "src/messages/$locale/create_storefront.json" "src/messages/$locale/create_b2c_store.json" || true
done

echo -e "\n${GREEN}Phase 3: Обновление импортов и ссылок в коде${NC}"
echo "================================================"

# Найти все TypeScript/JavaScript файлы
echo "Обновление импортов в .ts/.tsx/.js/.jsx файлах..."

# Marketplace → C2C
fd -e ts -e tsx -e js -e jsx . src -x sed -i \
    -e "s|from ['\"]@/services/marketplace'|from '@/services/c2c'|g" \
    -e "s|from ['\"]@/services/marketplaceApi'|from '@/services/c2cApi'|g" \
    -e "s|from ['\"]@/services/marketplaceOrders'|from '@/services/c2cOrders'|g" \
    -e "s|from ['\"]@/types/marketplace'|from '@/types/c2c'|g" \
    -e "s|from ['\"]@/components/marketplace|from '@/components/c2c|g" \
    -e "s|from ['\"]@/components/MarketplaceCard'|from '@/components/C2CCard'|g" \
    {}

# Storefronts → B2C
fd -e ts -e tsx -e js -e jsx . src -x sed -i \
    -e "s|from ['\"]@/services/storefrontApi'|from '@/services/b2cStoreApi'|g" \
    -e "s|from ['\"]@/services/storefrontProducts'|from '@/services/b2cProducts'|g" \
    -e "s|from ['\"]@/services/ai/storefronts.service'|from '@/services/ai/b2c.service'|g" \
    -e "s|from ['\"]@/types/storefront'|from '@/types/b2c'|g" \
    -e "s|from ['\"]@/components/storefronts|from '@/components/b2c|g" \
    -e "s|from ['\"]@/components/Storefront|from '@/components/B2C|g" \
    -e "s|from ['\"]@/hooks/useStorefronts'|from '@/hooks/useB2CStores'|g" \
    -e "s|from ['\"]@/contexts/CreateStorefrontContext'|from '@/contexts/CreateB2CStoreContext'|g" \
    -e "s|from ['\"]@/store/slices/storefrontSlice'|from '@/store/slices/b2cStoreSlice'|g" \
    {}

echo -e "\n${GREEN}Phase 4: Обновление URL маршрутов${NC}"
echo "================================================"

fd -e ts -e tsx -e js -e jsx . src -x sed -i \
    -e "s|/marketplace|/c2c|g" \
    -e "s|/storefronts|/b2c|g" \
    -e "s|/create-storefront|/create-b2c-store|g" \
    -e "s|/admin/storefront-products|/admin/b2c-products|g" \
    -e "s|href=['\"].*marketplace|href=\"/c2c|g" \
    -e "s|href=['\"].*storefronts|href=\"/b2c|g" \
    {}

echo -e "\n${GREEN}Phase 5: Обновление API endpoints${NC}"
echo "================================================"

fd -e ts -e tsx -e js -e jsx . src -x sed -i \
    -e "s|/api/v1/marketplace|/api/v1/c2c|g" \
    -e "s|/api/v1/admin/marketplace|/api/v1/admin/c2c|g" \
    -e "s|/api/v1/storefronts|/api/v1/b2c/stores|g" \
    -e "s|/api/v1/admin/storefronts|/api/v1/admin/b2c/stores|g" \
    -e "s|/api/v1/storefront-products|/api/v1/b2c/products|g" \
    -e "s|/api/v1/admin/storefront-products|/api/v1/admin/b2c/products|g" \
    {}

echo -e "\n${GREEN}Phase 6: Обновление названий типов${NC}"
echo "================================================"

# Type names в TypeScript файлах
fd -e ts -e tsx . src -x sed -i \
    -e "s/MarketplaceListing/C2CListing/g" \
    -e "s/MarketplaceItem/C2CItem/g" \
    -e "s/MarketplaceOrder/C2COrder/g" \
    -e "s/MarketplaceFilters/C2CFilters/g" \
    -e "s/StorefrontProduct/B2CProduct/g" \
    -e "s/StorefrontOrder/B2COrder/g" \
    -e "s/: Storefront\b/: B2CStore/g" \
    -e "s/<Storefront>/<B2CStore>/g" \
    {}

echo -e "\n${GREEN}Phase 7: Обновление i18n ключей${NC}"
echo "================================================"

# Обновить ключи переводов в компонентах
fd -e ts -e tsx . src -x sed -i \
    -e "s|t('marketplace\.|t('c2c.|g" \
    -e "s|t(\"marketplace\.|t(\"c2c.|g" \
    -e "s|t('storefronts\.|t('b2c.|g" \
    -e "s|t(\"storefronts\.|t(\"b2c.|g" \
    -e "s|t('create_storefront\.|t('create_b2c_store.|g" \
    -e "s|t(\"create_storefront\.|t(\"create_b2c_store.|g" \
    {}

echo -e "\n${GREEN}✅ Frontend migration complete!${NC}"
echo ""
echo "Следующие шаги:"
echo "1. Проверь изменения: git status"
echo "2. Запусти сборку: cd frontend/svetu && yarn build"
echo "3. Запусти линтер: yarn lint"
echo "4. Протестируй приложение"
