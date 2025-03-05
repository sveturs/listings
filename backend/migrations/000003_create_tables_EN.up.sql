-- English translations for the translations table
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text, is_machine_translated, is_verified, created_at, updated_at) VALUES
-- Main categories
('category', 1, 'en', 'name', 'Real Estate', true, true, NOW(), NOW()),
('category', 2, 'en', 'name', 'Automobiles', true, true, NOW(), NOW()),
('category', 3, 'en', 'name', 'Electronics', true, true, NOW(), NOW()),
('category', 4, 'en', 'name', 'Home and Apartment', true, true, NOW(), NOW()),
('category', 5, 'en', 'name', 'Garden', true, true, NOW(), NOW()),
('category', 6, 'en', 'name', 'Hobbies and Recreation', true, true, NOW(), NOW()),
('category', 7, 'en', 'name', 'Animals', true, true, NOW(), NOW()),
('category', 8, 'en', 'name', 'Ready Business and Equipment', true, true, NOW(), NOW()),
('category', 9, 'en', 'name', 'Other', true, true, NOW(), NOW()),
('category', 10, 'en', 'name', 'Jobs', true, true, NOW(), NOW()),
('category', 11, 'en', 'name', 'Clothing, Footwear, Accessories', true, true, NOW(), NOW()),
('category', 12, 'en', 'name', 'Children''s Clothing and Footwear, Accessories', true, true, NOW(), NOW()),

-- Real Estate (subcategories)
('category', 13, 'en', 'name', 'Apartments', true, true, NOW(), NOW()),
('category', 14, 'en', 'name', 'Houses', true, true, NOW(), NOW()),
('category', 15, 'en', 'name', 'Commercial Space', true, true, NOW(), NOW()),
('category', 16, 'en', 'name', 'Land', true, true, NOW(), NOW()),
('category', 17, 'en', 'name', 'Garages', true, true, NOW(), NOW()),
('category', 18, 'en', 'name', 'Rooms', true, true, NOW(), NOW()),
('category', 19, 'en', 'name', 'International Real Estate', true, true, NOW(), NOW()),
('category', 20, 'en', 'name', 'Commercial Property Sales', true, true, NOW(), NOW()),
('category', 21, 'en', 'name', 'Commercial Property Rental', true, true, NOW(), NOW()),

-- Automobile subcategories
('category', 22, 'en', 'name', 'Passenger Cars', true, true, NOW(), NOW()),
('category', 23, 'en', 'name', 'Trucks', true, true, NOW(), NOW()),
('category', 24, 'en', 'name', 'Special Equipment', true, true, NOW(), NOW()),
('category', 25, 'en', 'name', 'Agricultural Machinery', true, true, NOW(), NOW()),
('category', 26, 'en', 'name', 'Car and Special Equipment Rental', true, true, NOW(), NOW()),
('category', 27, 'en', 'name', 'Motorcycles and Motor Equipment', true, true, NOW(), NOW()),
('category', 28, 'en', 'name', 'Watercraft', true, true, NOW(), NOW()),
('category', 29, 'en', 'name', 'Spare Parts and Accessories', true, true, NOW(), NOW()),

-- Trucks
('category', 30, 'en', 'name', 'Trucks', true, true, NOW(), NOW()),
('category', 31, 'en', 'name', 'Semi-trailers', true, true, NOW(), NOW()),
('category', 32, 'en', 'name', 'Light Commercial Vehicles', true, true, NOW(), NOW()),
('category', 33, 'en', 'name', 'Buses', true, true, NOW(), NOW()),

-- Special Equipment
('category', 34, 'en', 'name', 'Excavators', true, true, NOW(), NOW()),
('category', 35, 'en', 'name', 'Loaders', true, true, NOW(), NOW()),
('category', 36, 'en', 'name', 'Backhoe Loaders', true, true, NOW(), NOW()),
('category', 37, 'en', 'name', 'Mobile Cranes', true, true, NOW(), NOW()),
('category', 38, 'en', 'name', 'Concrete Mixers', true, true, NOW(), NOW()),
('category', 39, 'en', 'name', 'Road Rollers', true, true, NOW(), NOW()),
('category', 40, 'en', 'name', 'Street Sweepers', true, true, NOW(), NOW()),
('category', 41, 'en', 'name', 'Garbage Trucks', true, true, NOW(), NOW()),
('category', 42, 'en', 'name', 'Aerial Work Platforms', true, true, NOW(), NOW()),
('category', 43, 'en', 'name', 'Bulldozers', true, true, NOW(), NOW()),
('category', 44, 'en', 'name', 'Motor Graders', true, true, NOW(), NOW()),
('category', 45, 'en', 'name', 'Drilling Rigs', true, true, NOW(), NOW()),

-- Agricultural Machinery
('category', 46, 'en', 'name', 'Tractors', true, true, NOW(), NOW()),
('category', 47, 'en', 'name', 'Mini Tractors', true, true, NOW(), NOW()),
('category', 48, 'en', 'name', 'Balers', true, true, NOW(), NOW()),
('category', 49, 'en', 'name', 'Harrows', true, true, NOW(), NOW()),
('category', 50, 'en', 'name', 'Mowers', true, true, NOW(), NOW()),

-- Garden (subcategories)
('category', 51, 'en', 'name', 'Garden Furniture', true, true, NOW(), NOW()),
('category', 52, 'en', 'name', 'Garden Tools', true, true, NOW(), NOW()),
('category', 53, 'en', 'name', 'Garden Plants', true, true, NOW(), NOW()),
('category', 54, 'en', 'name', 'Seeds and Seedlings', true, true, NOW(), NOW()),
('category', 55, 'en', 'name', 'BBQs and Equipment', true, true, NOW(), NOW()),
('category', 56, 'en', 'name', 'Pools and Equipment', true, true, NOW(), NOW()),
('category', 57, 'en', 'name', 'Irrigation', true, true, NOW(), NOW()),
('category', 58, 'en', 'name', 'Composting', true, true, NOW(), NOW()),
('category', 59, 'en', 'name', 'Greenhouses', true, true, NOW(), NOW()),
('category', 60, 'en', 'name', 'Fertilizers and Soil', true, true, NOW(), NOW()),

-- Hobbies and Recreation (subcategories)
('category', 101, 'en', 'name', 'Musical Instruments', true, true, NOW(), NOW()),
('category', 102, 'en', 'name', 'Books', true, true, NOW(), NOW()),
('category', 103, 'en', 'name', 'Sports Equipment', true, true, NOW(), NOW()),
('category', 104, 'en', 'name', 'Travel', true, true, NOW(), NOW()),
('category', 105, 'en', 'name', 'Collecting', true, true, NOW(), NOW()),
('category', 106, 'en', 'name', 'Art Objects', true, true, NOW(), NOW()),
('category', 107, 'en', 'name', 'Toys', true, true, NOW(), NOW()),
('category', 108, 'en', 'name', 'Bicycles', true, true, NOW(), NOW()),
('category', 109, 'en', 'name', 'Hunting and Fishing', true, true, NOW(), NOW()),
('category', 110, 'en', 'name', 'Camping', true, true, NOW(), NOW()),
('category', 111, 'en', 'name', 'Antiques', true, true, NOW(), NOW()),
('category', 112, 'en', 'name', 'Handicrafts', true, true, NOW(), NOW()),
('category', 113, 'en', 'name', 'Event Tickets', true, true, NOW(), NOW()),

-- Musical Instruments (subcategories)
('category', 114, 'en', 'name', 'String Instruments', true, true, NOW(), NOW()),
('category', 115, 'en', 'name', 'Pianos and Keyboards', true, true, NOW(), NOW()),
('category', 116, 'en', 'name', 'Percussion', true, true, NOW(), NOW()),
('category', 117, 'en', 'name', 'Wind Instruments', true, true, NOW(), NOW()),
('category', 118, 'en', 'name', 'Accordions and Harmoniums', true, true, NOW(), NOW()),
('category', 119, 'en', 'name', 'Audio Equipment', true, true, NOW(), NOW()),
('category', 120, 'en', 'name', 'Instrument Accessories', true, true, NOW(), NOW()),

-- Animals
('category', 200, 'en', 'name', 'Dogs', true, true, NOW(), NOW()),
('category', 201, 'en', 'name', 'Cats', true, true, NOW(), NOW()),
('category', 202, 'en', 'name', 'Birds', true, true, NOW(), NOW()),
('category', 203, 'en', 'name', 'Aquarium', true, true, NOW(), NOW()),
('category', 204, 'en', 'name', 'Other Animals', true, true, NOW(), NOW()),

-- Car and Special Equipment Rental
('category', 205, 'en', 'name', 'Cars', true, true, NOW(), NOW()),
('category', 206, 'en', 'name', 'Lifting Equipment', true, true, NOW(), NOW()),
('category', 207, 'en', 'name', 'Earthmoving Equipment', true, true, NOW(), NOW()),
('category', 208, 'en', 'name', 'Municipal Equipment', true, true, NOW(), NOW()),
('category', 209, 'en', 'name', 'Road Construction Equipment', true, true, NOW(), NOW()),
('category', 210, 'en', 'name', 'Commercial Vehicles', true, true, NOW(), NOW()),
('category', 211, 'en', 'name', 'Loading Equipment', true, true, NOW(), NOW()),
('category', 212, 'en', 'name', 'Attachments', true, true, NOW(), NOW()),
('category', 213, 'en', 'name', 'Trailers', true, true, NOW(), NOW()),
('category', 214, 'en', 'name', 'Agricultural Machinery', true, true, NOW(), NOW()),
('category', 215, 'en', 'name', 'Motorhomes', true, true, NOW(), NOW()),

-- Motorcycles and Motor Equipment
('category', 216, 'en', 'name', 'All-Terrain Vehicles', true, true, NOW(), NOW()),
('category', 217, 'en', 'name', 'Karting', true, true, NOW(), NOW()),
('category', 218, 'en', 'name', 'ATVs and Buggies', true, true, NOW(), NOW()),
('category', 219, 'en', 'name', 'Mopeds', true, true, NOW(), NOW()),
('category', 220, 'en', 'name', 'Scooters', true, true, NOW(), NOW()),
('category', 221, 'en', 'name', 'Motorcycles', true, true, NOW(), NOW()),
('category', 222, 'en', 'name', 'Snowmobiles', true, true, NOW(), NOW()),

-- Watercraft
('category', 223, 'en', 'name', 'Rowing Boats', true, true, NOW(), NOW()),
('category', 224, 'en', 'name', 'Kayaks', true, true, NOW(), NOW()),
('category', 225, 'en', 'name', 'Jet Skis', true, true, NOW(), NOW()),
('category', 226, 'en', 'name', 'Boats and Yachts', true, true, NOW(), NOW()),
('category', 227, 'en', 'name', 'Motorboats and Motors', true, true, NOW(), NOW()),

-- Spare Parts and Accessories
('category', 228, 'en', 'name', 'Spare Parts', true, true, NOW(), NOW()),
('category', 229, 'en', 'name', 'Tires, Rims and Wheels', true, true, NOW(), NOW()),
('category', 230, 'en', 'name', 'Audio and Video Equipment', true, true, NOW(), NOW()),
('category', 231, 'en', 'name', 'Accessories', true, true, NOW(), NOW()),
('category', 232, 'en', 'name', 'Oils and Auto Chemicals', true, true, NOW(), NOW()),
('category', 233, 'en', 'name', 'Tools', true, true, NOW(), NOW()),
('category', 234, 'en', 'name', 'Roof Racks and Tow Bars', true, true, NOW(), NOW()),
('category', 235, 'en', 'name', 'Trailers', true, true, NOW(), NOW()),
('category', 236, 'en', 'name', 'Equipment', true, true, NOW(), NOW()),
('category', 237, 'en', 'name', 'Anti-theft Devices', true, true, NOW(), NOW()),
('category', 238, 'en', 'name', 'GPS Navigators', true, true, NOW(), NOW()),

-- Electronics
('category', 239, 'en', 'name', 'Phones', true, true, NOW(), NOW()),
('category', 240, 'en', 'name', 'Audio and Video', true, true, NOW(), NOW()),
('category', 241, 'en', 'name', 'Computer Equipment', true, true, NOW(), NOW()),
('category', 242, 'en', 'name', 'Games, Consoles and Software', true, true, NOW(), NOW()),
('category', 243, 'en', 'name', 'Laptops', true, true, NOW(), NOW()),
('category', 244, 'en', 'name', 'Photo Equipment', true, true, NOW(), NOW()),
('category', 245, 'en', 'name', 'Tablets and E-readers', true, true, NOW(), NOW()),
('category', 246, 'en', 'name', 'Office Equipment and Supplies', true, true, NOW(), NOW()),
('category', 247, 'en', 'name', 'Household Appliances', true, true, NOW(), NOW()),

-- Phone Subcategories
('category', 248, 'en', 'name', 'Mobile Phones', true, true, NOW(), NOW()),
('category', 249, 'en', 'name', 'Accessories', true, true, NOW(), NOW()),
('category', 250, 'en', 'name', 'Two-way Radios', true, true, NOW(), NOW()),
('category', 251, 'en', 'name', 'Landline Phones', true, true, NOW(), NOW()),

-- Phone Accessories
('category', 252, 'en', 'name', 'Batteries', true, true, NOW(), NOW()),
('category', 253, 'en', 'name', 'Headsets and Earphones', true, true, NOW(), NOW()),
('category', 254, 'en', 'name', 'Chargers', true, true, NOW(), NOW()),
('category', 255, 'en', 'name', 'Cables and Adapters', true, true, NOW(), NOW()),
('category', 256, 'en', 'name', 'Modems and Routers', true, true, NOW(), NOW()),
('category', 257, 'en', 'name', 'Cases and Screen Protectors', true, true, NOW(), NOW()),
('category', 258, 'en', 'name', 'Spare Parts', true, true, NOW(), NOW()),

-- Audio and Video
('category', 259, 'en', 'name', 'TVs and Projectors', true, true, NOW(), NOW()),
('category', 260, 'en', 'name', 'Headphones', true, true, NOW(), NOW()),
('category', 261, 'en', 'name', 'Speakers, Sound Systems, Subwoofers', true, true, NOW(), NOW()),
('category', 262, 'en', 'name', 'Accessories', true, true, NOW(), NOW()),
('category', 263, 'en', 'name', 'Music Centers, Radio Players', true, true, NOW(), NOW()),
('category', 264, 'en', 'name', 'Amplifiers and Receivers', true, true, NOW(), NOW()),
('category', 265, 'en', 'name', 'Video Cameras', true, true, NOW(), NOW()),
('category', 266, 'en', 'name', 'Video, DVD and Blu-ray Players', true, true, NOW(), NOW()),
('category', 267, 'en', 'name', 'Cables and Adapters', true, true, NOW(), NOW()),
('category', 268, 'en', 'name', 'Music and Movies', true, true, NOW(), NOW()),
('category', 269, 'en', 'name', 'Microphones', true, true, NOW(), NOW()),
('category', 270, 'en', 'name', 'MP3 Players', true, true, NOW(), NOW()),

-- Computer Equipment
('category', 271, 'en', 'name', 'Desktop PCs', true, true, NOW(), NOW()),
('category', 272, 'en', 'name', 'All-in-One PCs', true, true, NOW(), NOW()),
('category', 273, 'en', 'name', 'Components', true, true, NOW(), NOW()),
('category', 274, 'en', 'name', 'Monitors and Parts', true, true, NOW(), NOW()),
('category', 275, 'en', 'name', 'Network Equipment', true, true, NOW(), NOW()),
('category', 276, 'en', 'name', 'Keyboards and Mice', true, true, NOW(), NOW()),
('category', 277, 'en', 'name', 'Joysticks and Steering Wheels', true, true, NOW(), NOW()),
('category', 278, 'en', 'name', 'USB Drives and Memory Cards', true, true, NOW(), NOW()),
('category', 279, 'en', 'name', 'Webcams', true, true, NOW(), NOW()),
('category', 280, 'en', 'name', 'TV Tuners', true, true, NOW(), NOW()),

-- Computer Components
('category', 281, 'en', 'name', 'CD, DVD and Blu-ray Drives', true, true, NOW(), NOW()),
('category', 282, 'en', 'name', 'Power Supplies', true, true, NOW(), NOW()),
('category', 283, 'en', 'name', 'Graphics Cards', true, true, NOW(), NOW()),
('category', 284, 'en', 'name', 'Hard Drives', true, true, NOW(), NOW()),
('category', 285, 'en', 'name', 'Sound Cards', true, true, NOW(), NOW()),
('category', 286, 'en', 'name', 'Controllers', true, true, NOW(), NOW()),
('category', 287, 'en', 'name', 'Cases', true, true, NOW(), NOW()),
('category', 288, 'en', 'name', 'Motherboards', true, true, NOW(), NOW()),
('category', 289, 'en', 'name', 'RAM', true, true, NOW(), NOW()),
('category', 290, 'en', 'name', 'Processors', true, true, NOW(), NOW()),
('category', 291, 'en', 'name', 'Cooling Systems', true, true, NOW(), NOW()),

-- Games, Consoles and Software
('category', 292, 'en', 'name', 'Gaming Consoles', true, true, NOW(), NOW()),
('category', 293, 'en', 'name', 'Console Games', true, true, NOW(), NOW()),
('category', 294, 'en', 'name', 'PC Games', true, true, NOW(), NOW()),
('category', 295, 'en', 'name', 'Software', true, true, NOW(), NOW()),

-- Photo Equipment
('category', 296, 'en', 'name', 'Equipment and Accessories', true, true, NOW(), NOW()),
('category', 297, 'en', 'name', 'Lenses', true, true, NOW(), NOW()),
('category', 298, 'en', 'name', 'Compact Cameras', true, true, NOW(), NOW()),
('category', 299, 'en', 'name', 'Film Cameras', true, true, NOW(), NOW()),
('category', 300, 'en', 'name', 'DSLR Cameras', true, true, NOW(), NOW()),

-- Tablets and E-readers
('category', 301, 'en', 'name', 'Tablets', true, true, NOW(), NOW()),
('category', 302, 'en', 'name', 'E-readers', true, true, NOW(), NOW()),
('category', 303, 'en', 'name', 'Accessories', true, true, NOW(), NOW()),

-- Office Equipment and Supplies
('category', 304, 'en', 'name', 'MFPs, Copiers and Scanners', true, true, NOW(), NOW()),
('category', 305, 'en', 'name', 'Printers', true, true, NOW(), NOW()),
('category', 306, 'en', 'name', 'Stationery', true, true, NOW(), NOW()),
('category', 307, 'en', 'name', 'UPS, Surge Protectors', true, true, NOW(), NOW()),
('category', 308, 'en', 'name', 'Telephony', true, true, NOW(), NOW()),
('category', 309, 'en', 'name', 'Paper Shredders', true, true, NOW(), NOW()),
('category', 310, 'en', 'name', 'Consumables', true, true, NOW(), NOW()),

-- Household Appliances
('category', 311, 'en', 'name', 'Kitchen Appliances', true, true, NOW(), NOW()),
('category', 312, 'en', 'name', 'Home Appliances', true, true, NOW(), NOW()),
('category', 313, 'en', 'name', 'Climate Control Equipment', true, true, NOW(), NOW()),
('category', 314, 'en', 'name', 'Personal Care Appliances', true, true, NOW(), NOW()),

-- Kitchen Appliances
('category', 315, 'en', 'name', 'Hoods', true, true, NOW(), NOW()),
('category', 316, 'en', 'name', 'Small Kitchen Appliances', true, true, NOW(), NOW()),
('category', 317, 'en', 'name', 'Microwave Ovens', true, true, NOW(), NOW()),
('category', 318, 'en', 'name', 'Stoves and Ovens', true, true, NOW(), NOW()),
('category', 319, 'en', 'name', 'Dishwashers', true, true, NOW(), NOW()),
('category', 320, 'en', 'name', 'Refrigerators and Freezers', true, true, NOW(), NOW()),

-- Home Appliances
('category', 321, 'en', 'name', 'Vacuum Cleaners and Parts', true, true, NOW(), NOW()),
('category', 322, 'en', 'name', 'Washing and Drying Machines', true, true, NOW(), NOW()),
('category', 323, 'en', 'name', 'Irons', true, true, NOW(), NOW()),
('category', 324, 'en', 'name', 'Sewing Equipment', true, true, NOW(), NOW()),

-- Climate Control Equipment
('category', 325, 'en', 'name', 'Fans', true, true, NOW(), NOW()),
('category', 326, 'en', 'name', 'Air Conditioners and Parts', true, true, NOW(), NOW()),
('category', 327, 'en', 'name', 'Heaters', true, true, NOW(), NOW()),
('category', 328, 'en', 'name', 'Air Purifiers', true, true, NOW(), NOW()),
('category', 329, 'en', 'name', 'Thermometers and Weather Stations', true, true, NOW(), NOW()),

-- Personal Care Appliances
('category', 330, 'en', 'name', 'Shavers and Trimmers', true, true, NOW(), NOW()),
('category', 331, 'en', 'name', 'Hair Clippers', true, true, NOW(), NOW()),
('category', 332, 'en', 'name', 'Hair Dryers and Styling Tools', true, true, NOW(), NOW()),
('category', 333, 'en', 'name', 'Epilators', true, true, NOW(), NOW()),

-- For Home and Apartment
('category', 334, 'en', 'name', 'Repair and Construction', true, true, NOW(), NOW()),
('category', 335, 'en', 'name', 'Furniture and Interior', true, true, NOW(), NOW()),
('category', 336, 'en', 'name', 'Food Products', true, true, NOW(), NOW()),
('category', 337, 'en', 'name', 'Cookware and Kitchen Utensils', true, true, NOW(), NOW()),

-- Main Categories and Other Important for Demonstration
('category', 338, 'en', 'name', 'Doors', true, true, NOW(), NOW()),
('category', 339, 'en', 'name', 'Tools', true, true, NOW(), NOW()),
('category', 340, 'en', 'name', 'Fireplaces and Heaters', true, true, NOW(), NOW()),
('category', 341, 'en', 'name', 'Windows and Balconies', true, true, NOW(), NOW()),
('category', 342, 'en', 'name', 'Ceilings', true, true, NOW(), NOW()),
('category', 343, 'en', 'name', 'For Garden and Cottage', true, true, NOW(), NOW()),
('category', 344, 'en', 'name', 'Plumbing, Water Supply and Sauna', true, true, NOW(), NOW()),
('category', 345, 'en', 'name', 'Prefabricated Buildings and Log Houses', true, true, NOW(), NOW()),
('category', 346, 'en', 'name', 'Gates, Fences and Barriers', true, true, NOW(), NOW()),
('category', 347, 'en', 'name', 'Security and Alarms', true, true, NOW(), NOW()),
('category', 348, 'en', 'name', 'Beds, Sofas and Armchairs', true, true, NOW(), NOW()),
('category', 349, 'en', 'name', 'Textiles and Carpets', true, true, NOW(), NOW()),
('category', 350, 'en', 'name', 'Lighting', true, true, NOW(), NOW()),
('category', 351, 'en', 'name', 'Computer Desks and Chairs', true, true, NOW(), NOW()),
('category', 352, 'en', 'name', 'Wardrobes, Chests of Drawers and Shelving', true, true, NOW(), NOW()),
('category', 353, 'en', 'name', 'Kitchen Sets', true, true, NOW(), NOW()),
('category', 354, 'en', 'name', 'Tables and Chairs', true, true, NOW(), NOW()),
('category', 355, 'en', 'name', 'Indoor Plants', true, true, NOW(), NOW()),
('category', 356, 'en', 'name', 'Decorative Foliage Plants', true, true, NOW(), NOW()),
('category', 357, 'en', 'name', 'Flowering Plants', true, true, NOW(), NOW()),
('category', 358, 'en', 'name', 'Palms and Ficus Trees', true, true, NOW(), NOW()),
('category', 359, 'en', 'name', 'Cacti and Succulents', true, true, NOW(), NOW()),
('category', 360, 'en', 'name', 'Tea, Coffee, Cocoa', true, true, NOW(), NOW()),
('category', 361, 'en', 'name', 'Beverages', true, true, NOW(), NOW()),
('category', 362, 'en', 'name', 'Fish, Seafood, Caviar', true, true, NOW(), NOW()),
('category', 363, 'en', 'name', 'Meat, Poultry, Offal', true, true, NOW(), NOW()),
('category', 364, 'en', 'name', 'Confectionery', true, true, NOW(), NOW()),
('category', 366, 'en', 'name', 'Tableware', true, true, NOW(), NOW()),
('category', 367, 'en', 'name', 'Kitchen Utensils', true, true, NOW(), NOW()),
('category', 368, 'en', 'name', 'Table Setting', true, true, NOW(), NOW()),
('category', 369, 'en', 'name', 'Food Preparation', true, true, NOW(), NOW()),
('category', 370, 'en', 'name', 'Food Storage', true, true, NOW(), NOW()),
('category', 371, 'en', 'name', 'Beverage Preparation', true, true, NOW(), NOW()),
('category', 372, 'en', 'name', 'Household Goods', true, true, NOW(), NOW()),

-- Все для сада (родитель id=5)
('category', 373, 'en', 'name', 'Garden Furniture', true, true, NOW(), NOW()),
('category', 374, 'en', 'name', 'Lighting', true, true, NOW(), NOW()),
('category', 375, 'en', 'name', 'Interior Design', true, true, NOW(), NOW()),
('category', 376, 'en', 'name', 'Plants and Seeds', true, true, NOW(), NOW()),
('category', 377, 'en', 'name', 'Garden and Vegetable Garden', true, true, NOW(), NOW()),
('category', 378, 'en', 'name', 'Garden Plants', true, true, NOW(), NOW()),
('category', 379, 'en', 'name', 'Seeds, Bulbs, Tubers', true, true, NOW(), NOW()),
('category', 380, 'en', 'name', 'Plant Care Products', true, true, NOW(), NOW()),
('category', 381, 'en', 'name', 'Panels and Artificial Plants', true, true, NOW(), NOW()),

-- Подкатегории Садовые растения (родитель id=378)
('category', 382, 'en', 'name', 'Decorative Shrubs and Trees', true, true, NOW(), NOW()),
('category', 383, 'en', 'name', 'Coniferous Plants', true, true, NOW(), NOW()),
('category', 384, 'en', 'name', 'Perennial Plants', true, true, NOW(), NOW()),
('category', 385, 'en', 'name', 'Fruit Plants', true, true, NOW(), NOW()),
('category', 386, 'en', 'name', 'Lawn', true, true, NOW(), NOW()),
('category', 387, 'en', 'name', 'Greens and Herbs', true, true, NOW(), NOW()),

-- Подкатегории Товары для ухода за растениями (родитель id=380)
('category', 388, 'en', 'name', 'Soils and Substrates', true, true, NOW(), NOW()),
('category', 389, 'en', 'name', 'Fertilizers', true, true, NOW(), NOW()),
('category', 390, 'en', 'name', 'Pest and Weed Control', true, true, NOW(), NOW()),
('category', 391, 'en', 'name', 'Pots and Planters', true, true, NOW(), NOW()),
('category', 392, 'en', 'name', 'Grow Lights', true, true, NOW(), NOW()),
('category', 393, 'en', 'name', 'Moisture Meters', true, true, NOW(), NOW()),
('category', 394, 'en', 'name', 'Greenhouses, Beds, Flowerbeds', true, true, NOW(), NOW()),

-- Aquarium (subcategories)
('category', 401, 'en', 'name', 'Aquariums and Equipment', true, true, NOW(), NOW()),
('category', 402, 'en', 'name', 'Fish', true, true, NOW(), NOW()),
('category', 403, 'en', 'name', 'Aquarium Plants', true, true, NOW(), NOW()),
('category', 404, 'en', 'name', 'Fish Food', true, true, NOW(), NOW()),
('category', 405, 'en', 'name', 'Filters and Accessories', true, true, NOW(), NOW()),

-- Other Animals (subcategories)
('category', 406, 'en', 'name', 'Small Rodents', true, true, NOW(), NOW()),
('category', 407, 'en', 'name', 'Reptiles', true, true, NOW(), NOW()),
('category', 408, 'en', 'name', 'Exotic Animals', true, true, NOW(), NOW()),
('category', 409, 'en', 'name', 'Bees and Hives', true, true, NOW(), NOW()),
('category', 410, 'en', 'name', 'Farm Animals', true, true, NOW(), NOW()),
('category', 411, 'en', 'name', 'Horses', true, true, NOW(), NOW()),
('category', 412, 'en', 'name', 'Animal Food', true, true, NOW(), NOW()),
('category', 413, 'en', 'name', 'Animal Equipment', true, true, NOW(), NOW()),

-- Other Categories
('category', 473, 'en', 'name', 'Rakija and Wine', true, true, NOW(), NOW()),
('category', 474, 'en', 'name', 'Homemade Cheese', true, true, NOW(), NOW()),
('category', 475, 'en', 'name', 'Kaymak', true, true, NOW(), NOW()),
('category', 476, 'en', 'name', 'Ajvar', true, true, NOW(), NOW()),
('category', 477, 'en', 'name', 'Plum Rakija', true, true, NOW(), NOW()),
('category', 478, 'en', 'name', 'Grape Rakija', true, true, NOW(), NOW()),
('category', 479, 'en', 'name', 'Fruit Rakija', true, true, NOW(), NOW()),
('category', 480, 'en', 'name', 'Homemade Wine', true, true, NOW(), NOW()),
('category', 481, 'en', 'name', 'Folk Crafts', true, true, NOW(), NOW()),
('category', 482, 'en', 'name', 'Opanci', true, true, NOW(), NOW()),
('category', 483, 'en', 'name', 'Ceramics', true, true, NOW(), NOW()),
('category', 484, 'en', 'name', 'Embroidery', true, true, NOW(), NOW()),
('category', 485, 'en', 'name', 'Weaving', true, true, NOW(), NOW()),
('category', 486, 'en', 'name', 'Folk Instruments', true, true, NOW(), NOW()),

-- Agricultural Categories
('category', 487, 'en', 'name', 'Beekeeping', true, true, NOW(), NOW()),

-- Beekeeping Subcategories (parent id=487)
('category', 488, 'en', 'name', 'Honey', true, true, NOW(), NOW()),
('category', 489, 'en', 'name', 'Beeswax', true, true, NOW(), NOW()),
('category', 490, 'en', 'name', 'Propolis', true, true, NOW(), NOW()),
('category', 491, 'en', 'name', 'Beekeeping Equipment', true, true, NOW(), NOW()),

-- Categories for Seasonal Work
('category', 492, 'en', 'name', 'Seasonal Work', true, true, NOW(), NOW()),
-- Seasonal Work Subcategories (parent id=492)
('category', 493, 'en', 'name', 'Harvesting', true, true, NOW(), NOW()),
('category', 494, 'en', 'name', 'Vineyard Work', true, true, NOW(), NOW()),
('category', 495, 'en', 'name', 'Seasonal Construction Work', true, true, NOW(), NOW()),

-- Tourism Services
('category', 496, 'en', 'name', 'Rural Tourism', true, true, NOW(), NOW()),

-- Rural Tourism Subcategories (parent id=496)
('category', 497, 'en', 'name', 'Ethno Villages', true, true, NOW(), NOW()),
('category', 498, 'en', 'name', 'Wine Tours', true, true, NOW(), NOW()),
('category', 499, 'en', 'name', 'Agritourism', true, true, NOW(), NOW()),
('category', 500, 'en', 'name', 'Mountain Tourism', true, true, NOW(), NOW());


-- Update sequence for translations
SELECT setval('translations_id_seq', (SELECT MAX(id) FROM translations), true);