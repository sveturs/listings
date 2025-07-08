-- Rename city back to address_city and country back to address_country in marketplace_listings table
ALTER TABLE marketplace_listings 
  RENAME COLUMN city TO address_city;

ALTER TABLE marketplace_listings 
  RENAME COLUMN country TO address_country;