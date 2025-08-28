-- Исправление автоинкремента для всех таблиц storefront
-- Эта миграция исправляет проблему с отсутствующими DEFAULT значениями

-- storefront_delivery_options
ALTER TABLE storefront_delivery_options 
ALTER COLUMN id SET DEFAULT nextval('storefront_delivery_options_id_seq');
ALTER SEQUENCE storefront_delivery_options_id_seq OWNED BY storefront_delivery_options.id;
SELECT setval('storefront_delivery_options_id_seq', COALESCE((SELECT MAX(id) FROM storefront_delivery_options), 1), true);

-- storefront_hours
ALTER TABLE storefront_hours 
ALTER COLUMN id SET DEFAULT nextval('storefront_hours_id_seq');
ALTER SEQUENCE storefront_hours_id_seq OWNED BY storefront_hours.id;
SELECT setval('storefront_hours_id_seq', COALESCE((SELECT MAX(id) FROM storefront_hours), 1), true);

-- storefront_inventory_movements
ALTER TABLE storefront_inventory_movements 
ALTER COLUMN id SET DEFAULT nextval('storefront_inventory_movements_id_seq');
ALTER SEQUENCE storefront_inventory_movements_id_seq OWNED BY storefront_inventory_movements.id;
SELECT setval('storefront_inventory_movements_id_seq', COALESCE((SELECT MAX(id) FROM storefront_inventory_movements), 1), true);

-- storefront_orders
ALTER TABLE storefront_orders 
ALTER COLUMN id SET DEFAULT nextval('storefront_orders_id_seq');
ALTER SEQUENCE storefront_orders_id_seq OWNED BY storefront_orders.id;
SELECT setval('storefront_orders_id_seq', COALESCE((SELECT MAX(id) FROM storefront_orders), 1), true);

-- storefront_payment_methods
ALTER TABLE storefront_payment_methods 
ALTER COLUMN id SET DEFAULT nextval('storefront_payment_methods_id_seq');
ALTER SEQUENCE storefront_payment_methods_id_seq OWNED BY storefront_payment_methods.id;
SELECT setval('storefront_payment_methods_id_seq', COALESCE((SELECT MAX(id) FROM storefront_payment_methods), 1), true);

-- storefront_product_attributes
ALTER TABLE storefront_product_attributes 
ALTER COLUMN id SET DEFAULT nextval('storefront_product_attributes_id_seq');
ALTER SEQUENCE storefront_product_attributes_id_seq OWNED BY storefront_product_attributes.id;
SELECT setval('storefront_product_attributes_id_seq', COALESCE((SELECT MAX(id) FROM storefront_product_attributes), 1), true);

-- storefront_product_variants
ALTER TABLE storefront_product_variants 
ALTER COLUMN id SET DEFAULT nextval('storefront_product_variants_id_seq');
ALTER SEQUENCE storefront_product_variants_id_seq OWNED BY storefront_product_variants.id;
SELECT setval('storefront_product_variants_id_seq', COALESCE((SELECT MAX(id) FROM storefront_product_variants), 1), true);