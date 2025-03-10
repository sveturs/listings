-- Увеличиваем точность для числовых значений в атрибутах
ALTER TABLE listing_attribute_values 
ALTER COLUMN numeric_value TYPE NUMERIC(20,5);