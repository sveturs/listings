-- Удаляем добавленные переводы для товаров витрин
DELETE FROM translations 
WHERE entity_type = 'storefront_product' 
  AND entity_id IN (225, 226, 227, 228, 229, 230, 231, 232, 233, 234, 235, 236, 237, 238, 239, 240, 241, 242)
  AND is_machine_translated = true;