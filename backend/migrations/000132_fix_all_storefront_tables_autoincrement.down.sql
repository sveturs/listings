-- Откат миграции: убираем DEFAULT значения
ALTER TABLE storefront_delivery_options ALTER COLUMN id DROP DEFAULT;
ALTER TABLE storefront_hours ALTER COLUMN id DROP DEFAULT;
ALTER TABLE storefront_inventory_movements ALTER COLUMN id DROP DEFAULT;
ALTER TABLE storefront_orders ALTER COLUMN id DROP DEFAULT;
ALTER TABLE storefront_payment_methods ALTER COLUMN id DROP DEFAULT;
ALTER TABLE storefront_product_attributes ALTER COLUMN id DROP DEFAULT;
ALTER TABLE storefront_product_variants ALTER COLUMN id DROP DEFAULT;