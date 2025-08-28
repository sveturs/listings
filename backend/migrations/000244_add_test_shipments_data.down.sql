-- Удаление тестовых данных

-- Удаляем тестовые данные из BEX shipments
DELETE FROM bex_shipments WHERE tracking_number LIKE 'BEX202411%';

-- Удаляем тестовые данные из Post Express shipments  
DELETE FROM post_express_shipments WHERE tracking_number LIKE 'PE202411%';