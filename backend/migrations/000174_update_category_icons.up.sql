-- Update category icons for all categories without icons

-- Root categories
UPDATE marketplace_categories SET icon = '🌾' WHERE slug = 'agriculture' AND (icon IS NULL OR icon = '');
UPDATE marketplace_categories SET icon = '🚗' WHERE slug = 'automotive' AND (icon IS NULL OR icon = '');
UPDATE marketplace_categories SET icon = '🏢' WHERE slug = 'business-industrial' AND (icon IS NULL OR icon = '');
UPDATE marketplace_categories SET icon = '🎨' WHERE slug = 'collectibles-hobby' AND (icon IS NULL OR icon = '');
UPDATE marketplace_categories SET icon = '📱' WHERE slug = 'electronics' AND (icon IS NULL OR icon = '');
UPDATE marketplace_categories SET icon = '👕' WHERE slug = 'fashion' AND (icon IS NULL OR icon = '');
UPDATE marketplace_categories SET icon = '🍽️' WHERE slug = 'food-beverages' AND (icon IS NULL OR icon = '');
UPDATE marketplace_categories SET icon = '🏡' WHERE slug = 'home-garden' AND (icon IS NULL OR icon = '');
UPDATE marketplace_categories SET icon = '🏭' WHERE slug = 'industrial' AND (icon IS NULL OR icon = '');
UPDATE marketplace_categories SET icon = '🏠' WHERE slug = 'real-estate' AND (icon IS NULL OR icon = '');
UPDATE marketplace_categories SET icon = '🛠️' WHERE slug = 'services' AND (icon IS NULL OR icon = '');
UPDATE marketplace_categories SET icon = '⚽' WHERE slug = 'sports-recreation' AND (icon IS NULL OR icon = '');
UPDATE marketplace_categories SET icon = '✈️' WHERE slug = 'travel-tourism' AND (icon IS NULL OR icon = '');

-- Electronics subcategories
UPDATE marketplace_categories SET icon = '💻' WHERE slug = 'computers' AND (icon IS NULL OR icon = '');
UPDATE marketplace_categories SET icon = '🚁' WHERE slug = 'drones-rc' AND (icon IS NULL OR icon = '');
UPDATE marketplace_categories SET icon = '🏠' WHERE slug = 'home-appliances' AND (icon IS NULL OR icon = '');
UPDATE marketplace_categories SET icon = '📱' WHERE slug = 'smartphones' AND (icon IS NULL OR icon = '');
UPDATE marketplace_categories SET icon = '⌚' WHERE slug = 'smartwatches' AND (icon IS NULL OR icon = '');
UPDATE marketplace_categories SET icon = '📺' WHERE slug = 'tv-audio' AND (icon IS NULL OR icon = '');
UPDATE marketplace_categories SET icon = '🥽' WHERE slug = 'vr-ar-equipment' AND (icon IS NULL OR icon = '');

-- Fashion subcategories
UPDATE marketplace_categories SET icon = '👔' WHERE slug = 'accessories' AND (icon IS NULL OR icon = '');
UPDATE marketplace_categories SET icon = '🎭' WHERE slug = 'costumes' AND (icon IS NULL OR icon = '');
UPDATE marketplace_categories SET icon = '👔' WHERE slug = 'mens-clothing' AND (icon IS NULL OR icon = '');
UPDATE marketplace_categories SET icon = '👠' WHERE slug = 'shoes' AND (icon IS NULL OR icon = '');
UPDATE marketplace_categories SET icon = '🎖️' WHERE slug = 'uniforms' AND (icon IS NULL OR icon = '');
UPDATE marketplace_categories SET icon = '🏛️' WHERE slug = 'vintage-clothing' AND (icon IS NULL OR icon = '');
UPDATE marketplace_categories SET icon = '👗' WHERE slug = 'womens-clothing' AND (icon IS NULL OR icon = '');

-- Automotive subcategories
UPDATE marketplace_categories SET icon = '🔧' WHERE slug = 'auto-parts' AND (icon IS NULL OR icon = '');
UPDATE marketplace_categories SET icon = '🚙' WHERE slug = 'cars' AND (icon IS NULL OR icon = '');
UPDATE marketplace_categories SET icon = '🏍️' WHERE slug = 'motorcycles' AND (icon IS NULL OR icon = '');

-- Real estate subcategories
UPDATE marketplace_categories SET icon = '🏢' WHERE slug = 'apartments' AND (icon IS NULL OR icon = '');
UPDATE marketplace_categories SET icon = '🏪' WHERE slug = 'commercial-real-estate' AND (icon IS NULL OR icon = '');
UPDATE marketplace_categories SET icon = '🏡' WHERE slug = 'houses' AND (icon IS NULL OR icon = '');
UPDATE marketplace_categories SET icon = '🌍' WHERE slug = 'land' AND (icon IS NULL OR icon = '');

-- Home & Garden subcategories
UPDATE marketplace_categories SET icon = '🧱' WHERE slug = 'building-materials' AND (icon IS NULL OR icon = '');
UPDATE marketplace_categories SET icon = '🛋️' WHERE slug = 'furniture' AND (icon IS NULL OR icon = '');
UPDATE marketplace_categories SET icon = '🌱' WHERE slug = 'garden-tools' AND (icon IS NULL OR icon = '');
UPDATE marketplace_categories SET icon = '🖼️' WHERE slug = 'home-decor' AND (icon IS NULL OR icon = '');

-- Agriculture subcategories
UPDATE marketplace_categories SET icon = '🚜' WHERE slug = 'farm-machinery' AND (icon IS NULL OR icon = '');
UPDATE marketplace_categories SET icon = '🌾' WHERE slug = 'farm-products' AND (icon IS NULL OR icon = '');
UPDATE marketplace_categories SET icon = '🐄' WHERE slug = 'livestock' AND (icon IS NULL OR icon = '');
UPDATE marketplace_categories SET icon = '🌱' WHERE slug = 'seeds-fertilizers' AND (icon IS NULL OR icon = '');

-- Industrial subcategories
UPDATE marketplace_categories SET icon = '🧪' WHERE slug = 'chemical-products' AND (icon IS NULL OR icon = '');
UPDATE marketplace_categories SET icon = '⚙️' WHERE slug = 'industrial-machinery' AND (icon IS NULL OR icon = '');
UPDATE marketplace_categories SET icon = '📦' WHERE slug = 'raw-materials' AND (icon IS NULL OR icon = '');
UPDATE marketplace_categories SET icon = '🦺' WHERE slug = 'safety-equipment' AND (icon IS NULL OR icon = '');

-- Food & Beverages subcategories
UPDATE marketplace_categories SET icon = '🍷' WHERE slug = 'beverages' AND (icon IS NULL OR icon = '');
UPDATE marketplace_categories SET icon = '🥛' WHERE slug = 'dairy-products' AND (icon IS NULL OR icon = '');
UPDATE marketplace_categories SET icon = '🥩' WHERE slug = 'meat-products' AND (icon IS NULL OR icon = '');
UPDATE marketplace_categories SET icon = '🥗' WHERE slug = 'organic-food' AND (icon IS NULL OR icon = '');

-- Services subcategories
UPDATE marketplace_categories SET icon = '💅' WHERE slug = 'beauty-wellness' AND (icon IS NULL OR icon = '');
UPDATE marketplace_categories SET icon = '💼' WHERE slug = 'business-services' AND (icon IS NULL OR icon = '');
UPDATE marketplace_categories SET icon = '🏗️' WHERE slug = 'construction-services' AND (icon IS NULL OR icon = '');
UPDATE marketplace_categories SET icon = '💻' WHERE slug = 'it-services' AND (icon IS NULL OR icon = '');

-- Sports & Recreation subcategories
UPDATE marketplace_categories SET icon = '🏋️' WHERE slug = 'fitness-equipment' AND (icon IS NULL OR icon = '');
UPDATE marketplace_categories SET icon = '🏃' WHERE slug = 'outdoor-sports' AND (icon IS NULL OR icon = '');
UPDATE marketplace_categories SET icon = '⚽' WHERE slug = 'team-sports' AND (icon IS NULL OR icon = '');
UPDATE marketplace_categories SET icon = '⛷️' WHERE slug = 'winter-sports' AND (icon IS NULL OR icon = '');

-- Auto parts subcategories
UPDATE marketplace_categories SET icon = '🚗' WHERE slug = 'auto-accessories' AND (icon IS NULL OR icon = '');
UPDATE marketplace_categories SET icon = '🚪' WHERE slug = 'body-parts' AND (icon IS NULL OR icon = '');
UPDATE marketplace_categories SET icon = '🛑' WHERE slug = 'brake-system' AND (icon IS NULL OR icon = '');
UPDATE marketplace_categories SET icon = '❄️' WHERE slug = 'cooling-system' AND (icon IS NULL OR icon = '');
UPDATE marketplace_categories SET icon = '⚡' WHERE slug = 'electrical-parts' AND (icon IS NULL OR icon = '');
UPDATE marketplace_categories SET icon = '🔧' WHERE slug = 'engine-and-parts' AND (icon IS NULL OR icon = '');
UPDATE marketplace_categories SET icon = '🚗' WHERE slug = 'interior-parts' AND (icon IS NULL OR icon = '');
UPDATE marketplace_categories SET icon = '🔧' WHERE slug = 'suspension-system' AND (icon IS NULL OR icon = '');
UPDATE marketplace_categories SET icon = '🛞' WHERE slug = 'tires-and-wheels' AND (icon IS NULL OR icon = '');
UPDATE marketplace_categories SET icon = '⚙️' WHERE slug = 'transmission-parts' AND (icon IS NULL OR icon = '');

-- Tires subcategories
UPDATE marketplace_categories SET icon = '🛞' WHERE slug = 'all-season-tires' AND (icon IS NULL OR icon = '');
UPDATE marketplace_categories SET icon = '🛞' WHERE slug = 'complete-wheels' AND (icon IS NULL OR icon = '');
UPDATE marketplace_categories SET icon = '⭕' WHERE slug = 'rims' AND (icon IS NULL OR icon = '');
UPDATE marketplace_categories SET icon = '☀️' WHERE slug = 'summer-tires' AND (icon IS NULL OR icon = '');
UPDATE marketplace_categories SET icon = '🔩' WHERE slug = 'wheel-bolts' AND (icon IS NULL OR icon = '');
UPDATE marketplace_categories SET icon = '🛞' WHERE slug = 'wheel-covers' AND (icon IS NULL OR icon = '');
UPDATE marketplace_categories SET icon = '❄️' WHERE slug = 'winter-tires' AND (icon IS NULL OR icon = '');

-- Engine parts subcategories
UPDATE marketplace_categories SET icon = '⛓️' WHERE slug = 'belts-and-chains' AND (icon IS NULL OR icon = '');
UPDATE marketplace_categories SET icon = '💨' WHERE slug = 'exhaust-system' AND (icon IS NULL OR icon = '');
UPDATE marketplace_categories SET icon = '🔍' WHERE slug = 'filters' AND (icon IS NULL OR icon = '');
UPDATE marketplace_categories SET icon = '🛢️' WHERE slug = 'oils-and-fluids' AND (icon IS NULL OR icon = '');
UPDATE marketplace_categories SET icon = '🔌' WHERE slug = 'spark-plugs' AND (icon IS NULL OR icon = '');

-- Body parts subcategories
UPDATE marketplace_categories SET icon = '🚗' WHERE slug = 'bumpers' AND (icon IS NULL OR icon = '');
UPDATE marketplace_categories SET icon = '🚪' WHERE slug = 'doors' AND (icon IS NULL OR icon = '');
UPDATE marketplace_categories SET icon = '🚗' WHERE slug = 'fenders' AND (icon IS NULL OR icon = '');
UPDATE marketplace_categories SET icon = '🚗' WHERE slug = 'hoods' AND (icon IS NULL OR icon = '');
UPDATE marketplace_categories SET icon = '🪞' WHERE slug = 'mirrors' AND (icon IS NULL OR icon = '');
UPDATE marketplace_categories SET icon = '🪟' WHERE slug = 'windows' AND (icon IS NULL OR icon = '');

-- Tire types subcategories
UPDATE marketplace_categories SET icon = '🚙' WHERE slug = 'passenger-summer-tires' AND (icon IS NULL OR icon = '');
UPDATE marketplace_categories SET icon = '🚐' WHERE slug = 'suv-summer-tires' AND (icon IS NULL OR icon = '');
UPDATE marketplace_categories SET icon = '🚛' WHERE slug = 'truck-summer-tires' AND (icon IS NULL OR icon = '');
UPDATE marketplace_categories SET icon = '🚙' WHERE slug = 'passenger-winter-tires' AND (icon IS NULL OR icon = '');
UPDATE marketplace_categories SET icon = '🚐' WHERE slug = 'suv-winter-tires' AND (icon IS NULL OR icon = '');
UPDATE marketplace_categories SET icon = '🚛' WHERE slug = 'truck-winter-tires' AND (icon IS NULL OR icon = '');

-- Rims subcategories
UPDATE marketplace_categories SET icon = '⭕' WHERE slug = 'aluminum-rims' AND (icon IS NULL OR icon = '');
UPDATE marketplace_categories SET icon = '🏁' WHERE slug = 'sport-rims' AND (icon IS NULL OR icon = '');
UPDATE marketplace_categories SET icon = '⚙️' WHERE slug = 'steel-rims' AND (icon IS NULL OR icon = '');

-- Business & Industrial subcategories
UPDATE marketplace_categories SET icon = '📋' WHERE slug = 'office-supplies' AND (icon IS NULL OR icon = '');
UPDATE marketplace_categories SET icon = '🖨️' WHERE slug = 'printing-graphics' AND (icon IS NULL OR icon = '');
UPDATE marketplace_categories SET icon = '🔒' WHERE slug = 'security-safety' AND (icon IS NULL OR icon = '');

-- Collectibles & Hobby subcategories
UPDATE marketplace_categories SET icon = '🪙' WHERE slug = 'coins-banknotes' AND (icon IS NULL OR icon = '');
UPDATE marketplace_categories SET icon = '🚂' WHERE slug = 'models-miniatures' AND (icon IS NULL OR icon = '');
UPDATE marketplace_categories SET icon = '📮' WHERE slug = 'stamps' AND (icon IS NULL OR icon = '');

-- Travel & Tourism subcategories
UPDATE marketplace_categories SET icon = '🏨' WHERE slug = 'hotels-accommodation' AND (icon IS NULL OR icon = '');
UPDATE marketplace_categories SET icon = '🗺️' WHERE slug = 'tour-guides' AND (icon IS NULL OR icon = '');
UPDATE marketplace_categories SET icon = '🚌' WHERE slug = 'transport-rides' AND (icon IS NULL OR icon = '');