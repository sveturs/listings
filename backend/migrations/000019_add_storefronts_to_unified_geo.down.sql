-- Удаляем витрины из unified_geo
DELETE FROM unified_geo WHERE source_type = 'storefront';

-- Удаляем товары витрин из unified_geo
DELETE FROM unified_geo WHERE source_type = 'storefront_product';