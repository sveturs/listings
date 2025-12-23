-- Migration: 20251222120000_category_ai_mapping
-- Purpose: Create AI category mapping table to resolve Claude category names to DB slugs
-- Fixes Bug #1: No mapping between AI categories and DB slugs

-- Create table for AI category mapping
CREATE TABLE IF NOT EXISTS category_ai_mapping (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ai_category_name VARCHAR(100) NOT NULL UNIQUE,
    target_category_id UUID NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    confidence_boost NUMERIC(3,2) DEFAULT 0.15 CHECK (confidence_boost >= 0 AND confidence_boost <= 1.0),
    priority INTEGER DEFAULT 100 CHECK (priority > 0),
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
    notes TEXT
);

-- Create indexes for performance
CREATE UNIQUE INDEX IF NOT EXISTS idx_ai_mapping_name ON category_ai_mapping(ai_category_name);
CREATE INDEX IF NOT EXISTS idx_ai_mapping_target ON category_ai_mapping(target_category_id);
CREATE INDEX IF NOT EXISTS idx_ai_mapping_priority ON category_ai_mapping(priority DESC);

-- Insert AI category mappings for ALL major categories
-- Fashion & Jewelry (Часы, украшения, аксессуары)
INSERT INTO category_ai_mapping (ai_category_name, target_category_id, confidence_boost, priority, notes) VALUES
('Fashion', (SELECT id FROM categories WHERE slug = 'nakit-i-satovi' LIMIT 1), 0.20, 90, 'General fashion category'),
('Jewelry', (SELECT id FROM categories WHERE slug = 'nakit-i-satovi' LIMIT 1), 0.25, 100, 'Jewelry and accessories'),
('Jewellery', (SELECT id FROM categories WHERE slug = 'nakit-i-satovi' LIMIT 1), 0.25, 100, 'British spelling variant'),
('Watches', (SELECT id FROM categories WHERE slug = 'nakit-i-satovi' LIMIT 1), 0.30, 100, 'Watches category'),
('Watch', (SELECT id FROM categories WHERE slug = 'nakit-i-satovi' LIMIT 1), 0.30, 95, 'Singular form'),
('Accessories', (SELECT id FROM categories WHERE slug = 'nakit-i-satovi' LIMIT 1), 0.15, 80, 'Fashion accessories'),

-- Smartwatches (умные часы) - специальный маппинг на подкатегорию
('Smartwatch', (SELECT id FROM categories WHERE slug = 'pametni-satovi' LIMIT 1), 0.30, 100, 'Smartwatches category'),
('Smartwatches', (SELECT id FROM categories WHERE slug = 'pametni-satovi' LIMIT 1), 0.30, 100, 'Plural form'),
('Wearables', (SELECT id FROM categories WHERE slug = 'pametni-satovi' LIMIT 1), 0.25, 95, 'Wearable devices'),
('Fitness Tracker', (SELECT id FROM categories WHERE slug = 'pametni-satovi' LIMIT 1), 0.25, 90, 'Fitness tracking devices'),
('Smart Band', (SELECT id FROM categories WHERE slug = 'pametni-satovi' LIMIT 1), 0.25, 90, 'Smart bands'),

-- Electronics (вся электроника)
('Electronics', (SELECT id FROM categories WHERE slug = 'elektronika' LIMIT 1), 0.25, 100, 'General electronics'),
('Electronic', (SELECT id FROM categories WHERE slug = 'elektronika' LIMIT 1), 0.25, 95, 'Singular form'),
('Gadgets', (SELECT id FROM categories WHERE slug = 'elektronika' LIMIT 1), 0.20, 85, 'Tech gadgets'),
('Technology', (SELECT id FROM categories WHERE slug = 'elektronika' LIMIT 1), 0.20, 85, 'Technology products'),

-- Computers & Equipment
('Computers', (SELECT id FROM categories WHERE slug = 'elektronika' LIMIT 1), 0.30, 100, 'Computers and laptops'),
('Computer', (SELECT id FROM categories WHERE slug = 'elektronika' LIMIT 1), 0.30, 95, 'Singular form'),
('Laptops', (SELECT id FROM categories WHERE slug = 'elektronika' LIMIT 1), 0.30, 95, 'Laptop computers'),
('PC', (SELECT id FROM categories WHERE slug = 'elektronika' LIMIT 1), 0.25, 90, 'Personal computers'),

-- Mobile Phones
('Mobile Phones', (SELECT id FROM categories WHERE slug = 'elektronika' LIMIT 1), 0.30, 100, 'Mobile phones'),
('Phones', (SELECT id FROM categories WHERE slug = 'elektronika' LIMIT 1), 0.25, 95, 'General phones'),
('Smartphones', (SELECT id FROM categories WHERE slug = 'elektronika' LIMIT 1), 0.30, 100, 'Smart phones'),
('Smartphone', (SELECT id FROM categories WHERE slug = 'elektronika' LIMIT 1), 0.30, 95, 'Singular form'),

-- Clothing & Footwear
('Clothing', (SELECT id FROM categories WHERE slug = 'odeca-i-obuca' LIMIT 1), 0.25, 100, 'Clothing and apparel'),
('Clothes', (SELECT id FROM categories WHERE slug = 'odeca-i-obuca' LIMIT 1), 0.25, 95, 'Plural form'),
('Footwear', (SELECT id FROM categories WHERE slug = 'odeca-i-obuca' LIMIT 1), 0.20, 90, 'Shoes and footwear'),
('Apparel', (SELECT id FROM categories WHERE slug = 'odeca-i-obuca' LIMIT 1), 0.20, 85, 'Apparel and garments'),
('Shoes', (SELECT id FROM categories WHERE slug = 'odeca-i-obuca' LIMIT 1), 0.25, 90, 'Footwear'),

-- Automotive
('Automotive', (SELECT id FROM categories WHERE slug = 'automobilizam' LIMIT 1), 0.25, 100, 'Automotive category'),
('Cars', (SELECT id FROM categories WHERE slug = 'automobilizam' LIMIT 1), 0.30, 100, 'Cars and vehicles'),
('Car', (SELECT id FROM categories WHERE slug = 'automobilizam' LIMIT 1), 0.30, 95, 'Singular form'),
('Auto Parts', (SELECT id FROM categories WHERE slug = 'automobilizam' LIMIT 1), 0.25, 95, 'Auto parts and accessories'),
('Vehicles', (SELECT id FROM categories WHERE slug = 'automobilizam' LIMIT 1), 0.20, 90, 'Vehicles'),
('Vehicle', (SELECT id FROM categories WHERE slug = 'automobilizam' LIMIT 1), 0.20, 85, 'Singular form'),

-- Home & Garden
('Home & Garden', (SELECT id FROM categories WHERE slug = 'dom-i-basta' LIMIT 1), 0.25, 100, 'Home and garden'),
('Home', (SELECT id FROM categories WHERE slug = 'dom-i-basta' LIMIT 1), 0.20, 85, 'Home products'),
('Garden', (SELECT id FROM categories WHERE slug = 'dom-i-basta' LIMIT 1), 0.20, 85, 'Garden products'),
('Furniture', (SELECT id FROM categories WHERE slug = 'dom-i-basta' LIMIT 1), 0.20, 90, 'Furniture'),
('Home Decor', (SELECT id FROM categories WHERE slug = 'dom-i-basta' LIMIT 1), 0.20, 85, 'Home decoration'),

-- Appliances
('Appliances', (SELECT id FROM categories WHERE slug = 'kucni-aparati' LIMIT 1), 0.25, 100, 'Home appliances'),
('Home Appliances', (SELECT id FROM categories WHERE slug = 'kucni-aparati' LIMIT 1), 0.30, 100, 'Home appliances'),
('Kitchen Appliances', (SELECT id FROM categories WHERE slug = 'kucni-aparati' LIMIT 1), 0.25, 95, 'Kitchen appliances'),

-- Books & Media
('Books', (SELECT id FROM categories WHERE slug = 'knjige-i-mediji' LIMIT 1), 0.25, 100, 'Books'),
('Book', (SELECT id FROM categories WHERE slug = 'knjige-i-mediji' LIMIT 1), 0.25, 95, 'Singular form'),
('Media', (SELECT id FROM categories WHERE slug = 'knjige-i-mediji' LIMIT 1), 0.20, 90, 'Media products'),
('Music', (SELECT id FROM categories WHERE slug = 'knjige-i-mediji' LIMIT 1), 0.20, 85, 'Music media'),

-- Sports & Outdoors
('Sports', (SELECT id FROM categories WHERE slug = 'sport-i-turizam' LIMIT 1), 0.25, 100, 'Sports equipment'),
('Sport', (SELECT id FROM categories WHERE slug = 'sport-i-turizam' LIMIT 1), 0.25, 95, 'Singular form'),
('Outdoors', (SELECT id FROM categories WHERE slug = 'sport-i-turizam' LIMIT 1), 0.20, 90, 'Outdoor equipment'),
('Fitness', (SELECT id FROM categories WHERE slug = 'sport-i-turizam' LIMIT 1), 0.20, 85, 'Fitness equipment'),

-- Toys & Kids
('Toys', (SELECT id FROM categories WHERE slug = 'za-bebe-i-decu' LIMIT 1), 0.25, 100, 'Toys'),
('Toy', (SELECT id FROM categories WHERE slug = 'za-bebe-i-decu' LIMIT 1), 0.25, 95, 'Singular form'),
('Kids', (SELECT id FROM categories WHERE slug = 'za-bebe-i-decu' LIMIT 1), 0.20, 90, 'Kids products'),
('Children', (SELECT id FROM categories WHERE slug = 'za-bebe-i-decu' LIMIT 1), 0.20, 90, 'Children products'),
('Baby', (SELECT id FROM categories WHERE slug = 'za-bebe-i-decu' LIMIT 1), 0.25, 95, 'Baby products'),

-- Beauty & Health
('Beauty', (SELECT id FROM categories WHERE slug = 'lepota-i-zdravlje' LIMIT 1), 0.25, 100, 'Beauty products'),
('Health', (SELECT id FROM categories WHERE slug = 'lepota-i-zdravlje' LIMIT 1), 0.20, 90, 'Health products'),
('Cosmetics', (SELECT id FROM categories WHERE slug = 'lepota-i-zdravlje' LIMIT 1), 0.25, 95, 'Cosmetics'),
('Skincare', (SELECT id FROM categories WHERE slug = 'lepota-i-zdravlje' LIMIT 1), 0.20, 85, 'Skincare products'),

-- Pet Supplies
('Pet Supplies', (SELECT id FROM categories WHERE slug = 'kucni-ljubimci' LIMIT 1), 0.25, 100, 'Pet supplies'),
('Pets', (SELECT id FROM categories WHERE slug = 'kucni-ljubimci' LIMIT 1), 0.20, 90, 'Pet products'),
('Pet', (SELECT id FROM categories WHERE slug = 'kucni-ljubimci' LIMIT 1), 0.20, 85, 'Singular form'),

-- Musical Instruments
('Musical Instruments', (SELECT id FROM categories WHERE slug = 'muzicki-instrumenti' LIMIT 1), 0.30, 100, 'Musical instruments'),
('Music Instruments', (SELECT id FROM categories WHERE slug = 'muzicki-instrumenti' LIMIT 1), 0.25, 95, 'Alternative phrasing'),

-- Art & Crafts
('Art', (SELECT id FROM categories WHERE slug = 'umetnost-i-rukotvorine' LIMIT 1), 0.25, 100, 'Art products'),
('Crafts', (SELECT id FROM categories WHERE slug = 'umetnost-i-rukotvorine' LIMIT 1), 0.25, 95, 'Craft products'),
('Handmade', (SELECT id FROM categories WHERE slug = 'umetnost-i-rukotvorine' LIMIT 1), 0.20, 90, 'Handmade items'),

-- Office Supplies
('Office Supplies', (SELECT id FROM categories WHERE slug = 'kancelarijski-materijal' LIMIT 1), 0.25, 100, 'Office supplies'),
('Office', (SELECT id FROM categories WHERE slug = 'kancelarijski-materijal' LIMIT 1), 0.20, 85, 'Office products'),
('Stationery', (SELECT id FROM categories WHERE slug = 'kancelarijski-materijal' LIMIT 1), 0.25, 95, 'Stationery'),

-- Industrial & Tools
('Industrial', (SELECT id FROM categories WHERE slug = 'industrija-i-alati' LIMIT 1), 0.25, 100, 'Industrial products'),
('Tools', (SELECT id FROM categories WHERE slug = 'industrija-i-alati' LIMIT 1), 0.25, 95, 'Tools'),
('Tool', (SELECT id FROM categories WHERE slug = 'industrija-i-alati' LIMIT 1), 0.25, 90, 'Singular form'),

-- Food & Beverages
('Food', (SELECT id FROM categories WHERE slug = 'hrana-i-pice' LIMIT 1), 0.25, 100, 'Food products'),
('Beverages', (SELECT id FROM categories WHERE slug = 'hrana-i-pice' LIMIT 1), 0.25, 95, 'Beverages'),
('Drinks', (SELECT id FROM categories WHERE slug = 'hrana-i-pice' LIMIT 1), 0.20, 90, 'Drinks'),

-- Services
('Services', (SELECT id FROM categories WHERE slug = 'usluge' LIMIT 1), 0.25, 100, 'Services'),
('Service', (SELECT id FROM categories WHERE slug = 'usluge' LIMIT 1), 0.25, 95, 'Singular form');

-- Trigger for updated_at
CREATE OR REPLACE FUNCTION update_ai_mapping_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_ai_mapping_updated_at
    BEFORE UPDATE ON category_ai_mapping
    FOR EACH ROW
    EXECUTE FUNCTION update_ai_mapping_updated_at();

-- Add comment
COMMENT ON TABLE category_ai_mapping IS 'Mapping table for AI-detected category names to database category IDs';
COMMENT ON COLUMN category_ai_mapping.ai_category_name IS 'Category name as returned by Claude AI (e.g., Fashion, Electronics)';
COMMENT ON COLUMN category_ai_mapping.target_category_id IS 'Target category ID in categories table';
COMMENT ON COLUMN category_ai_mapping.confidence_boost IS 'Confidence score boost (0.0-1.0) to add to detection';
COMMENT ON COLUMN category_ai_mapping.priority IS 'Priority for conflict resolution (higher = preferred)';
