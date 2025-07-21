-- Обновление опций для атрибутов типа select
-- Формат: {"values": ["option1", "option2", ...]}

-- Condition (состояние)
UPDATE category_attributes 
SET options = '{"values": ["new", "used", "refurbished", "damaged"]}'
WHERE name = 'condition' AND attribute_type = 'select';

-- Brand (бренд)
UPDATE category_attributes 
SET options = '{"values": ["Apple", "Samsung", "Sony", "LG", "Dell", "HP", "Lenovo", "Asus", "Acer", "Other"]}'
WHERE name = 'brand' AND attribute_type = 'select';

-- Color (цвет)
UPDATE category_attributes 
SET options = '{"values": ["black", "white", "silver", "gold", "blue", "red", "green", "yellow", "purple", "other"]}'
WHERE name = 'color' AND attribute_type = 'select';

-- Storage (память)
UPDATE category_attributes 
SET options = '{"values": ["16GB", "32GB", "64GB", "128GB", "256GB", "512GB", "1TB", "2TB"]}'
WHERE name = 'storage' AND attribute_type = 'select';

-- Operating System (операционная система)
UPDATE category_attributes 
SET options = '{"values": ["iOS", "Android", "Windows", "macOS", "Linux", "Other"]}'
WHERE name = 'operating_system' AND attribute_type = 'select';

-- Gender (пол)
UPDATE category_attributes 
SET options = '{"values": ["male", "female", "unisex", "boys", "girls"]}'
WHERE name = 'gender' AND attribute_type = 'select';

-- Size (размер одежды)
UPDATE category_attributes 
SET options = '{"values": ["XS", "S", "M", "L", "XL", "XXL", "XXXL"]}'
WHERE name = 'clothing_size' AND attribute_type = 'select';

-- Shoe Size (размер обуви)
UPDATE category_attributes 
SET options = '{"values": ["35", "36", "37", "38", "39", "40", "41", "42", "43", "44", "45", "46"]}'
WHERE name = 'shoe_size' AND attribute_type = 'select';

-- Fuel Type (тип топлива)
UPDATE category_attributes 
SET options = '{"values": ["petrol", "diesel", "electric", "hybrid", "lpg", "cng"]}'
WHERE name = 'fuel_type' AND attribute_type = 'select';

-- Transmission (коробка передач)
UPDATE category_attributes 
SET options = '{"values": ["manual", "automatic", "semi-automatic", "cvt"]}'
WHERE name = 'transmission' AND attribute_type = 'select';

-- Property Type (тип недвижимости)
UPDATE category_attributes 
SET options = '{"values": ["apartment", "house", "land", "commercial", "garage", "other"]}'
WHERE name = 'property_type' AND attribute_type = 'select';

-- Pet Type (тип питомца)
UPDATE category_attributes 
SET options = '{"values": ["dog", "cat", "bird", "fish", "reptile", "small_animal", "other"]}'
WHERE name = 'pet_type' AND attribute_type = 'select';