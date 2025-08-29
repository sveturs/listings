-- Удаляем добавленные переводы для объявлений
DELETE FROM translations 
WHERE entity_type = 'listing' 
  AND entity_id IN (172, 173, 174, 175, 176, 177, 178, 179, 180, 181)
  AND is_machine_translated = true;

-- Удаляем добавленные переводы для витрин
DELETE FROM translations 
WHERE entity_type = 'storefront' 
  AND entity_id IN (1, 18, 19, 20)
  AND is_machine_translated = true;

-- Удаляем добавленные переводы для товаров
DELETE FROM translations 
WHERE entity_type = 'storefront_product' 
  AND entity_id IN (215, 216, 217, 218, 219, 220, 221, 222, 223, 224)
  AND is_machine_translated = true;