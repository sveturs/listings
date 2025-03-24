-- Добавление новых полей в таблицу user_storefronts
ALTER TABLE user_storefronts 
  ADD COLUMN phone VARCHAR(50),
  ADD COLUMN email VARCHAR(255),
  ADD COLUMN website VARCHAR(255),
  ADD COLUMN address VARCHAR(255),
  ADD COLUMN city VARCHAR(100),
  ADD COLUMN country VARCHAR(100),
  ADD COLUMN latitude DECIMAL(10,8),
  ADD COLUMN longitude DECIMAL(11,8);