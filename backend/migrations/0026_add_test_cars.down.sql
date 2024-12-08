-- backend/migrations/0026_add_test_cars.down.sql
DELETE FROM car_feature_links WHERE car_id IN (SELECT id FROM cars WHERE make IN ('Axia', 'Velaris', 'Terra', 'Navis'));
DELETE FROM car_images WHERE car_id IN (SELECT id FROM cars WHERE make IN ('Axia', 'Velaris', 'Terra', 'Navis'));
DELETE FROM cars WHERE make IN ('Axia', 'Velaris', 'Terra', 'Navis');