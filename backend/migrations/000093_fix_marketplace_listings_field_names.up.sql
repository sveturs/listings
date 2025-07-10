-- Rename address_city to city and address_country to country in marketplace_listings table
ALTER TABLE marketplace_listings 
  RENAME COLUMN address_city TO city;

ALTER TABLE marketplace_listings 
  RENAME COLUMN address_country TO country;