-- Добавляем поколения для Audi A4 (model_id = 58)
INSERT INTO car_generations (model_id, name, slug, year_start, year_end, is_active, sort_order)
VALUES 
    (58, 'B9', 'b9', 2015, NULL, true, 1),
    (58, 'B8', 'b8', 2008, 2015, true, 2),
    (58, 'B7', 'b7', 2004, 2008, true, 3),
    (58, 'B6', 'b6', 2000, 2006, true, 4),
    (58, 'B5', 'b5', 1994, 2001, true, 5);