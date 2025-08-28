-- Revert category icons to NULL

UPDATE marketplace_categories SET icon = NULL WHERE slug IN (
    'agriculture', 'automotive', 'business-industrial', 'collectibles-hobby', 
    'electronics', 'fashion', 'food-beverages', 'home-garden', 'industrial', 
    'real-estate', 'services', 'sports-recreation', 'travel-tourism',
    'computers', 'drones-rc', 'home-appliances', 'smartphones', 'smartwatches', 
    'tv-audio', 'vr-ar-equipment', 'accessories', 'costumes', 'mens-clothing', 
    'shoes', 'uniforms', 'vintage-clothing', 'womens-clothing', 'auto-parts', 
    'cars', 'motorcycles', 'apartments', 'commercial-real-estate', 'houses', 
    'land', 'building-materials', 'furniture', 'garden-tools', 'home-decor', 
    'farm-machinery', 'farm-products', 'livestock', 'seeds-fertilizers', 
    'chemical-products', 'industrial-machinery', 'raw-materials', 'safety-equipment', 
    'beverages', 'dairy-products', 'meat-products', 'organic-food', 'beauty-wellness', 
    'business-services', 'construction-services', 'it-services', 'fitness-equipment', 
    'outdoor-sports', 'team-sports', 'winter-sports', 'auto-accessories', 
    'body-parts', 'brake-system', 'cooling-system', 'electrical-parts', 
    'engine-and-parts', 'interior-parts', 'suspension-system', 'tires-and-wheels', 
    'transmission-parts', 'all-season-tires', 'complete-wheels', 'rims', 
    'summer-tires', 'wheel-bolts', 'wheel-covers', 'winter-tires', 'belts-and-chains', 
    'exhaust-system', 'filters', 'oils-and-fluids', 'spark-plugs', 'bumpers', 
    'doors', 'fenders', 'hoods', 'mirrors', 'windows', 'passenger-summer-tires', 
    'suv-summer-tires', 'truck-summer-tires', 'passenger-winter-tires', 
    'suv-winter-tires', 'truck-winter-tires', 'aluminum-rims', 'sport-rims', 
    'steel-rims', 'office-supplies', 'printing-graphics', 'security-safety', 
    'coins-banknotes', 'models-miniatures', 'stamps', 'hotels-accommodation', 
    'tour-guides', 'transport-rides'
);