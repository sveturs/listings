#!/bin/bash

echo "Финальное исправление swagger типов..."

# Функция для поиска типа в разных местах
find_type_location() {
    local type_name=$1

    # Ищем в domain/models
    if grep -q "type $type_name " /data/hostel-booking-system/backend/internal/domain/models/*.go 2>/dev/null; then
        echo "backend_internal_domain_models.$type_name"
        return
    fi

    # Ищем в domain/logistics
    if grep -q "type $type_name " /data/hostel-booking-system/backend/internal/domain/logistics/*.go 2>/dev/null; then
        echo "backend_internal_domain_logistics.$type_name"
        return
    fi

    # Ищем в domain/behavior
    if grep -q "type $type_name " /data/hostel-booking-system/backend/internal/domain/behavior/*.go 2>/dev/null; then
        echo "backend_internal_domain_behavior.$type_name"
        return
    fi

    # Ищем в domain/search
    if grep -q "type $type_name " /data/hostel-booking-system/backend/internal/domain/search/*.go 2>/dev/null; then
        echo "backend_internal_domain_search.$type_name"
        return
    fi

    echo ""
}

# Исправление balance module
sed -i 's/backend_internal_proj_balance_models\./backend_internal_domain_models./g' /data/hostel-booking-system/backend/internal/proj/balance/handler/*.go

# Исправление users module - многие типы в domain/models
sed -i 's/backend_internal_proj_users_models\.User\b/backend_internal_domain_models.User/g' /data/hostel-booking-system/backend/internal/proj/users/handler/*.go
sed -i 's/backend_internal_proj_users_models\.UpdateUserRequest/backend_internal_domain_models.UpdateUserRequest/g' /data/hostel-booking-system/backend/internal/proj/users/handler/*.go

# Исправление contacts module
sed -i 's/backend_internal_proj_contacts_models\./backend_internal_domain_models./g' /data/hostel-booking-system/backend/internal/proj/contacts/handler/*.go

# Исправление payments module
sed -i 's/backend_internal_proj_payments_models\.Payment\b/backend_internal_domain_models.Payment/g' /data/hostel-booking-system/backend/internal/proj/payments/handler/*.go

# Исправление orders module
sed -i 's/backend_internal_proj_orders_models\.Order\b/backend_internal_domain_models.Order/g' /data/hostel-booking-system/backend/internal/proj/orders/handler/*.go
sed -i 's/backend_internal_proj_orders_models\.OrderItem\b/backend_internal_domain_models.OrderItem/g' /data/hostel-booking-system/backend/internal/proj/orders/handler/*.go

# Исправление marketplace module - многие типы находятся в domain/models
sed -i 's/backend_internal_proj_marketplace_models\.MarketplaceListing/backend_internal_domain_models.MarketplaceListing/g' /data/hostel-booking-system/backend/internal/proj/marketplace/handler/*.go
sed -i 's/backend_internal_proj_marketplace_models\.Category/backend_internal_domain_models.Category/g' /data/hostel-booking-system/backend/internal/proj/marketplace/handler/*.go
sed -i 's/backend_internal_proj_marketplace_models\.CreateListingRequest/backend_internal_domain_models.CreateListingRequest/g' /data/hostel-booking-system/backend/internal/proj/marketplace/handler/*.go
sed -i 's/backend_internal_proj_marketplace_models\.UpdateListingRequest/backend_internal_domain_models.UpdateListingRequest/g' /data/hostel-booking-system/backend/internal/proj/marketplace/handler/*.go

# Исправление cars - находятся в domain/models
sed -i 's/backend_internal_proj_marketplace_models\.CarMake/backend_internal_domain_models.CarMake/g' /data/hostel-booking-system/backend/internal/proj/marketplace/handler/*.go
sed -i 's/backend_internal_proj_marketplace_models\.CarModel/backend_internal_domain_models.CarModel/g' /data/hostel-booking-system/backend/internal/proj/marketplace/handler/*.go
sed -i 's/backend_internal_proj_marketplace_models\.CarGeneration/backend_internal_domain_models.CarGeneration/g' /data/hostel-booking-system/backend/internal/proj/marketplace/handler/*.go
sed -i 's/backend_internal_proj_marketplace_models\.VINDecodeResult/backend_internal_domain_models.VINDecodeResult/g' /data/hostel-booking-system/backend/internal/proj/marketplace/handler/*.go

# Исправление reviews module
sed -i 's/backend_internal_proj_reviews_models\.Review\b/backend_internal_domain_models.Review/g' /data/hostel-booking-system/backend/internal/proj/reviews/handler/*.go
sed -i 's/backend_internal_proj_reviews_models\.CreateReviewRequest/backend_internal_domain_models.CreateReviewRequest/g' /data/hostel-booking-system/backend/internal/proj/reviews/handler/*.go

# Исправление storefronts module
sed -i 's/backend_internal_proj_storefronts_models\.Storefront/backend_internal_domain_models.Storefront/g' /data/hostel-booking-system/backend/internal/proj/storefronts/handler/*.go
sed -i 's/backend_internal_proj_storefronts_models\.CreateStorefrontRequest/backend_internal_domain_models.CreateStorefrontRequest/g' /data/hostel-booking-system/backend/internal/proj/storefronts/handler/*.go

# Исправление subscriptions module
sed -i 's/backend_internal_proj_subscriptions_models\.UserSubscription/backend_internal_domain_models.UserSubscription/g' /data/hostel-booking-system/backend/internal/proj/subscriptions/handler/*.go

# Исправление notifications module
sed -i 's/backend_internal_proj_notifications_models\.Notification/backend_internal_domain_models.Notification/g' /data/hostel-booking-system/backend/internal/proj/notifications/handler/*.go

# Исправление translations module
sed -i 's/backend_internal_proj_translation_models\./backend_internal_domain_models./g' /data/hostel-booking-system/backend/internal/proj/translation_admin/*.go

# Исправление gis module - используют domain/models
sed -i 's/backend_internal_proj_gis_models\./backend_internal_domain_models./g' /data/hostel-booking-system/backend/internal/proj/gis/handler/*.go

# Исправление geocode module
sed -i 's/backend_internal_proj_geocode_models\./backend_internal_domain_models./g' /data/hostel-booking-system/backend/internal/proj/geocode/handler/*.go

# Исправление docserver module
sed -i 's/backend_internal_proj_docserver_models\./backend_internal_domain_models./g' /data/hostel-booking-system/backend/internal/proj/docserver/handler/*.go

# Исправление analytics module
sed -i 's/backend_internal_proj_analytics_models\./backend_internal_domain_models./g' /data/hostel-booking-system/backend/internal/proj/analytics/handler/*.go

echo "Финальные исправления завершены!"