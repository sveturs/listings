#!/bin/bash

# Автоматически исправить все неправильные ссылки на модели в swagger аннотациях

echo "Исправление ссылок на модели в swagger аннотациях..."

# Исправление для delivery модуля
sed -i 's/backend_internal_domain_models\.DeliveryProvider/backend_internal_proj_delivery_models.DeliveryProvider/g' /data/hostel-booking-system/backend/internal/proj/delivery/handler/*.go
sed -i 's/backend_internal_domain_models\.ShippingRate/backend_internal_proj_delivery_models.ShippingRate/g' /data/hostel-booking-system/backend/internal/proj/delivery/handler/*.go

# Исправление для VIN модуля
sed -i 's/backend_internal_domain_models\.VINDecodeResult/backend_internal_proj_vin_models.VINDecodeResult/g' /data/hostel-booking-system/backend/internal/proj/vin/handler/*.go
sed -i 's/backend_internal_domain_models\.VINDecodeRequest/backend_internal_proj_vin_models.VINDecodeRequest/g' /data/hostel-booking-system/backend/internal/proj/vin/handler/*.go

# Исправление для recommendations модуля
sed -i 's/backend_internal_domain_models\.RecommendationResponse/backend_internal_proj_recommendations_models.RecommendationResponse/g' /data/hostel-booking-system/backend/internal/proj/recommendations/*.go
sed -i 's/backend_internal_domain_models\.RecommendationRequest/backend_internal_proj_recommendations_models.RecommendationRequest/g' /data/hostel-booking-system/backend/internal/proj/recommendations/*.go

# Исправление для postexpress модуля
sed -i 's/backend_internal_domain_models\.PostExpressShipment/backend_internal_proj_postexpress_models.PostExpressShipment/g' /data/hostel-booking-system/backend/internal/proj/postexpress/handler/*.go
sed -i 's/backend_internal_domain_models\.CreateShipmentRequest/backend_internal_proj_postexpress_models.CreateShipmentRequest/g' /data/hostel-booking-system/backend/internal/proj/postexpress/handler/*.go

# Исправление для orders модуля
sed -i 's/backend_internal_domain_models\.OrderDelivery/backend_internal_proj_orders_models.OrderDelivery/g' /data/hostel-booking-system/backend/internal/proj/orders/handler/*.go
sed -i 's/backend_internal_domain_models\.UpdateOrderRequest/backend_internal_proj_orders_models.UpdateOrderRequest/g' /data/hostel-booking-system/backend/internal/proj/orders/handler/*.go

# Исправление для payments модуля
sed -i 's/backend_internal_domain_models\.PaymentMethod/backend_internal_proj_payments_models.PaymentMethod/g' /data/hostel-booking-system/backend/internal/proj/payments/handler/*.go
sed -i 's/backend_internal_domain_models\.Payment/backend_internal_proj_payments_models.Payment/g' /data/hostel-booking-system/backend/internal/proj/payments/handler/*.go

# Исправление для storefronts модуля
sed -i 's/backend_internal_domain_models\.StorefrontProduct/backend_internal_proj_storefronts_models.StorefrontProduct/g' /data/hostel-booking-system/backend/internal/proj/storefronts/handler/*.go
sed -i 's/backend_internal_domain_models\.CreateProductRequest/backend_internal_proj_storefronts_models.CreateProductRequest/g' /data/hostel-booking-system/backend/internal/proj/storefronts/handler/*.go

# Исправление для reviews модуля
sed -i 's/backend_internal_domain_models\.Review/backend_internal_proj_reviews_models.Review/g' /data/hostel-booking-system/backend/internal/proj/reviews/handler/*.go
sed -i 's/backend_internal_domain_models\.CreateReviewRequest/backend_internal_proj_reviews_models.CreateReviewRequest/g' /data/hostel-booking-system/backend/internal/proj/reviews/handler/*.go

# Исправление для subscriptions модуля
sed -i 's/backend_internal_domain_models\.Subscription/backend_internal_proj_subscriptions_models.Subscription/g' /data/hostel-booking-system/backend/internal/proj/subscriptions/handler/*.go
sed -i 's/backend_internal_domain_models\.CreateSubscriptionRequest/backend_internal_proj_subscriptions_models.CreateSubscriptionRequest/g' /data/hostel-booking-system/backend/internal/proj/subscriptions/handler/*.go

echo "Исправления завершены!"