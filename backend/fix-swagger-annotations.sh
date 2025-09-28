#!/bin/bash

# Скрипт для исправления неправильных ссылок на типы в Swagger аннотациях

echo "Исправление Swagger аннотаций в backend/internal/proj/marketplace/handler/..."

# Функция для безопасной замены
safe_replace() {
    local file=$1
    local old=$2
    local new=$3

    if grep -q "$old" "$file"; then
        sed -i "s|$old|$new|g" "$file"
        echo "  ✓ Исправлено: $old -> $new"
    fi
}

# 1. unified_attributes.go
echo "Обработка unified_attributes.go..."
file="backend/internal/proj/marketplace/handler/unified_attributes.go"
safe_replace "$file" "backend_internal_proj_marketplace_models.UnifiedAttribute" "backend_internal_domain_models.UnifiedAttribute"
safe_replace "$file" "backend_internal_proj_marketplace_models.UnifiedCategoryAttribute" "backend_internal_domain_models.UnifiedCategoryAttribute"
safe_replace "$file" "backend_internal_proj_marketplace_models.UnifiedAttributeValue" "backend_internal_domain_models.UnifiedAttributeValue"

# 2. admin_categories.go
echo "Обработка admin_categories.go..."
file="backend/internal/proj/marketplace/handler/admin_categories.go"
safe_replace "$file" "backend_internal_proj_marketplace_models.AttributeGroup" "backend_internal_domain_models.AttributeGroup"
safe_replace "$file" "backend_internal_proj_marketplace_models.MarketplaceCategory" "backend_internal_domain_models.MarketplaceCategory"

# 3. custom_components.go
echo "Обработка custom_components.go..."
file="backend/internal/proj/marketplace/handler/custom_components.go"
safe_replace "$file" "backend_internal_proj_marketplace_models.CreateCustomComponentRequest" "backend_internal_domain_models.CreateCustomComponentRequest"
safe_replace "$file" "backend_internal_proj_marketplace_models.CustomUIComponent" "backend_internal_domain_models.CustomUIComponent"
safe_replace "$file" "backend_internal_proj_marketplace_models.UpdateCustomComponentRequest" "backend_internal_domain_models.UpdateCustomComponentRequest"
safe_replace "$file" "backend_internal_proj_marketplace_models.CreateComponentUsageRequest" "backend_internal_domain_models.CreateComponentUsageRequest"
safe_replace "$file" "backend_internal_proj_marketplace_models.CustomUIComponentUsage" "backend_internal_domain_models.CustomUIComponentUsage"
safe_replace "$file" "backend_internal_proj_marketplace_models.CreateTemplateRequest" "backend_internal_domain_models.CreateTemplateRequest"
safe_replace "$file" "backend_internal_proj_marketplace_models.ComponentTemplate" "backend_internal_domain_models.ComponentTemplate"

# 4. translation.go
echo "Обработка translation.go..."
file="backend/internal/proj/marketplace/handler/translation.go"
safe_replace "$file" "backend_internal_proj_marketplace_models.Translation" "backend_internal_domain_models.Translation"

# 5. listings.go
echo "Обработка listings.go..."
file="backend/internal/proj/marketplace/handler/listings.go"
safe_replace "$file" "backend_internal_proj_marketplace_models.PriceHistoryEntry" "backend_internal_domain_models.PriceHistoryEntry"
safe_replace "$file" "backend_internal_proj_marketplace_models.MarketplaceListing" "backend_internal_domain_models.MarketplaceListing"
safe_replace "$file" "backend_internal_proj_marketplace_models.CreateListingRequest" "backend_internal_domain_models.CreateListingRequest"
safe_replace "$file" "backend_internal_proj_marketplace_models.UpdateListingRequest" "backend_internal_domain_models.UpdateListingRequest"

# 6. categories.go
echo "Обработка categories.go..."
file="backend/internal/proj/marketplace/handler/categories.go"
safe_replace "$file" "backend_internal_proj_marketplace_models.MarketplaceCategory" "backend_internal_domain_models.MarketplaceCategory"

# 7. admin_variant_attributes.go
echo "Обработка admin_variant_attributes.go..."
file="backend/internal/proj/marketplace/handler/admin_variant_attributes.go"
safe_replace "$file" "backend_internal_proj_marketplace_models.ProductVariantAttribute" "backend_internal_domain_models.ProductVariantAttribute"

# 8. variant_mappings.go
echo "Обработка variant_mappings.go..."
file="backend/internal/proj/marketplace/handler/variant_mappings.go"
safe_replace "$file" "backend_internal_proj_marketplace_models.UnifiedAttribute" "backend_internal_domain_models.UnifiedAttribute"
safe_replace "$file" "backend_internal_proj_marketplace_models.VariantAttributeMapping" "backend_internal_domain_models.VariantAttributeMapping"
safe_replace "$file" "backend_internal_proj_marketplace_models.VariantAttributeMappingCreateRequest" "backend_internal_domain_models.VariantAttributeMappingCreateRequest"
safe_replace "$file" "backend_internal_proj_marketplace_models.VariantAttributeMappingUpdateRequest" "backend_internal_domain_models.VariantAttributeMappingUpdateRequest"

# 9. chat.go
echo "Обработка chat.go..."
file="backend/internal/proj/marketplace/handler/chat.go"
safe_replace "$file" "backend_internal_proj_marketplace_models.MarketplaceChat" "backend_internal_domain_models.MarketplaceChat"
safe_replace "$file" "backend_internal_proj_marketplace_models.CreateMessageRequest" "backend_internal_domain_models.CreateMessageRequest"
safe_replace "$file" "backend_internal_proj_marketplace_models.MarketplaceMessage" "backend_internal_domain_models.MarketplaceMessage"
safe_replace "$file" "backend_internal_proj_marketplace_models.MarkAsReadRequest" "backend_internal_domain_models.MarkAsReadRequest"
safe_replace "$file" "backend_internal_proj_marketplace_models.ChatAttachment" "backend_internal_domain_models.ChatAttachment"

# 10. order_handler.go
echo "Обработка order_handler.go..."
file="backend/internal/proj/marketplace/handler/order_handler.go"
safe_replace "$file" "backend_internal_proj_marketplace_models.MarketplaceOrder" "backend_internal_domain_models.MarketplaceOrder"

# 11. variant_attributes.go
echo "Обработка variant_attributes.go..."
file="backend/internal/proj/marketplace/handler/variant_attributes.go"
safe_replace "$file" "backend_internal_proj_marketplace_models.ProductVariantAttribute" "backend_internal_domain_models.ProductVariantAttribute"

echo "✅ Все исправления применены!"