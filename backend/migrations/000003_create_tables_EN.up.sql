-- English translations of categories for the translations table
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text, is_machine_translated, is_verified, created_at, updated_at) VALUES

('category', 1000, 'en', 'name', 'Real Estate', true, true, NOW(), NOW()),
    -- Real Estate subcategories 1000
    ('category', 1100, 'en', 'name', 'Apartment', true, true, NOW(), NOW()),
    ('category', 1200, 'en', 'name', 'Room', true, true, NOW(), NOW()),
    ('category', 1300, 'en', 'name', 'House, Cottage, Villa', true, true, NOW(), NOW()),
    ('category', 1400, 'en', 'name', 'Land Plot', true, true, NOW(), NOW()),
    ('category', 1500, 'en', 'name', 'Garage and Parking Space', true, true, NOW(), NOW()),
    ('category', 1600, 'en', 'name', 'Commercial Real Estate', true, true, NOW(), NOW()),
    ('category', 1700, 'en', 'name', 'Overseas Real Estate', true, true, NOW(), NOW()),
    ('category', 1800, 'en', 'name', 'Hotel', true, true, NOW(), NOW()),
    ('category', 1900, 'en', 'name', 'Apartments', true, true, NOW(), NOW()),

('category', 2000, 'en', 'name', 'Vehicles', true, true, NOW(), NOW()),

    -- Auto categories 2000
    ('category', 2100, 'en', 'name', 'Passenger Cars', true, true, NOW(), NOW()),
    ('category', 2200, 'en', 'name', 'Commercial Vehicles', true, true, NOW(), NOW()),
        -- Commercial vehicles subcategories 2200
        ('category', 2210, 'en', 'name', 'Trucks', true, true, NOW(), NOW()),
        ('category', 2220, 'en', 'name', 'Semi-trailers', true, true, NOW(), NOW()),
        ('category', 2230, 'en', 'name', 'Light Commercial Vehicles', true, true, NOW(), NOW()),
        ('category', 2240, 'en', 'name', 'Buses', true, true, NOW(), NOW()),
        
    ('category', 2300, 'en', 'name', 'Special Equipment', true, true, NOW(), NOW()),
        -- Special equipment subcategories 2300
        ('category', 2310, 'en', 'name', 'Excavators', true, true, NOW(), NOW()),
        ('category', 2315, 'en', 'name', 'Loaders', true, true, NOW(), NOW()),
        ('category', 2320, 'en', 'name', 'Backhoe Loaders', true, true, NOW(), NOW()),
        ('category', 2325, 'en', 'name', 'Mobile Cranes', true, true, NOW(), NOW()),
        ('category', 2330, 'en', 'name', 'Concrete Mixers', true, true, NOW(), NOW()),
        ('category', 2335, 'en', 'name', 'Road Rollers', true, true, NOW(), NOW()),
        ('category', 2340, 'en', 'name', 'Street Sweepers', true, true, NOW(), NOW()),
        ('category', 2345, 'en', 'name', 'Garbage Trucks', true, true, NOW(), NOW()),
        ('category', 2350, 'en', 'name', 'Aerial Platforms', true, true, NOW(), NOW()),
        ('category', 2355, 'en', 'name', 'Bulldozers', true, true, NOW(), NOW()),
        ('category', 2360, 'en', 'name', 'Motor Graders', true, true, NOW(), NOW()),
        ('category', 2365, 'en', 'name', 'Drilling Rigs', true, true, NOW(), NOW()),

    ('category', 2400, 'en', 'name', 'Agricultural Machinery', true, true, NOW(), NOW()),
        -- Agricultural machinery subcategories 2400
        ('category', 2410, 'en', 'name', 'Tractors', true, true, NOW(), NOW()),
        ('category', 2415, 'en', 'name', 'Mini Tractors', true, true, NOW(), NOW()),
        ('category', 2420, 'en', 'name', 'Balers', true, true, NOW(), NOW()),
        ('category', 2425, 'en', 'name', 'Harrows', true, true, NOW(), NOW()),
        ('category', 2430, 'en', 'name', 'Mowers', true, true, NOW(), NOW()),
        ('category', 2435, 'en', 'name', 'Combines', true, true, NOW(), NOW()),
        ('category', 2440, 'en', 'name', 'Telescopic Loaders', true, true, NOW(), NOW()),
        ('category', 2445, 'en', 'name', 'Seeders', true, true, NOW(), NOW()),
        ('category', 2450, 'en', 'name', 'Cultivators', true, true, NOW(), NOW()),
        ('category', 2455, 'en', 'name', 'Plows', true, true, NOW(), NOW()),
        ('category', 2460, 'en', 'name', 'Sprayers', true, true, NOW(), NOW()),

    ('category', 2500, 'en', 'name', 'Vehicle and Equipment Rental', true, true, NOW(), NOW()),
            -- Vehicle and Equipment Rental (parent id=2500)
        ('category', 2510, 'en', 'name', 'Cars', true, true, NOW(), NOW()),
        ('category', 2515, 'en', 'name', 'Lifting Equipment', true, true, NOW(), NOW()),
        ('category', 2520, 'en', 'name', 'Earthmoving Equipment', true, true, NOW(), NOW()),
        ('category', 2525, 'en', 'name', 'Municipal Equipment', true, true, NOW(), NOW()),
        ('category', 2530, 'en', 'name', 'Road Construction Equipment', true, true, NOW(), NOW()),
        ('category', 2535, 'en', 'name', 'Freight Transport', true, true, NOW(), NOW()),
        ('category', 2540, 'en', 'name', 'Loading Equipment', true, true, NOW(), NOW()),
        ('category', 2545, 'en', 'name', 'Attachments', true, true, NOW(), NOW()),
        ('category', 2550, 'en', 'name', 'Trailers', true, true, NOW(), NOW()),
        ('category', 2555, 'en', 'name', 'Agricultural Equipment', true, true, NOW(), NOW()),
        ('category', 2560, 'en', 'name', 'Motorhomes', true, true, NOW(), NOW()),

    ('category', 2600, 'en', 'name', 'Motorcycles and Motor Vehicles', true, true, NOW(), NOW()),
        -- Motorcycles and Motor Vehicles (parent id=2600)
        ('category', 2610, 'en', 'name', 'All-terrain Vehicles', true, true, NOW(), NOW()),
        ('category', 2615, 'en', 'name', 'Karting', true, true, NOW(), NOW()),
        ('category', 2620, 'en', 'name', 'ATVs and Buggies', true, true, NOW(), NOW()),
        ('category', 2625, 'en', 'name', 'Mopeds', true, true, NOW(), NOW()),
        ('category', 2630, 'en', 'name', 'Scooters', true, true, NOW(), NOW()),
        ('category', 2635, 'en', 'name', 'Motorcycles', true, true, NOW(), NOW()),
        ('category', 2640, 'en', 'name', 'Snowmobiles', true, true, NOW(), NOW()),

    ('category', 2700, 'en', 'name', 'Watercraft', true, true, NOW(), NOW()),
        -- Watercraft (parent id=2700)
        ('category', 2710, 'en', 'name', 'Rowing Boats', true, true, NOW(), NOW()),
        ('category', 2720, 'en', 'name', 'Kayaks', true, true, NOW(), NOW()),
        ('category', 2730, 'en', 'name', 'Jet Skis', true, true, NOW(), NOW()),
        ('category', 2740, 'en', 'name', 'Boats and Yachts', true, true, NOW(), NOW()),
        ('category', 2750, 'en', 'name', 'Motorboats and Motors', true, true, NOW(), NOW()),

    ('category', 2800, 'en', 'name', 'Parts and Accessories', true, true, NOW(), NOW()),
        -- Parts and Accessories (parent id=2800)
        ('category', 2810, 'en', 'name', 'Spare Parts', true, true, NOW(), NOW()),
        ('category', 2815, 'en', 'name', 'Tires, Rims and Wheels', true, true, NOW(), NOW()),
        ('category', 2820, 'en', 'name', 'Audio and Video Equipment', true, true, NOW(), NOW()),
        ('category', 2825, 'en', 'name', 'Accessories', true, true, NOW(), NOW()),
        ('category', 2830, 'en', 'name', 'Oils and Auto Chemistry', true, true, NOW(), NOW()),
        ('category', 2835, 'en', 'name', 'Tools', true, true, NOW(), NOW()),
        ('category', 2840, 'en', 'name', 'Roof Racks and Towbars', true, true, NOW(), NOW()),
        ('category', 2845, 'en', 'name', 'Trailers', true, true, NOW(), NOW()),
        ('category', 2850, 'en', 'name', 'Equipment', true, true, NOW(), NOW()),
        ('category', 2855, 'en', 'name', 'Anti-theft Devices', true, true, NOW(), NOW()),
        ('category', 2860, 'en', 'name', 'GPS Navigation', true, true, NOW(), NOW()),

('category', 3000, 'en', 'name', 'Electronics', true, true, NOW(), NOW()),
    -- Electronics (parent id=3000)
    ('category', 3100, 'en', 'name', 'Phones', true, true, NOW(), NOW()),
            -- Phone subcategories (parent id=3100)
        ('category', 3110, 'en', 'name', 'Mobile Phones', true, true, NOW(), NOW()),
        ('category', 3120, 'en', 'name', 'Accessories', true, true, NOW(), NOW()),
                    -- Phone accessories subcategories (parent id=3120)
            ('category', 3121, 'en', 'name', 'Batteries', true, true, NOW(), NOW()),
            ('category', 3122, 'en', 'name', 'Headsets and Earphones', true, true, NOW(), NOW()),
            ('category', 3123, 'en', 'name', 'Chargers', true, true, NOW(), NOW()),
            ('category', 3124, 'en', 'name', 'Cables and Adapters', true, true, NOW(), NOW()),
            ('category', 3125, 'en', 'name', 'Modems and Routers', true, true, NOW(), NOW()),
            ('category', 3126, 'en', 'name', 'Cases and Screen Protectors', true, true, NOW(), NOW()),
            ('category', 3127, 'en', 'name', 'Spare Parts', true, true, NOW(), NOW()),

        ('category', 3130, 'en', 'name', 'Two-way Radios', true, true, NOW(), NOW()),
        ('category', 3140, 'en', 'name', 'Landline Phones', true, true, NOW(), NOW()),

    ('category', 3200, 'en', 'name', 'Audio and Video', true, true, NOW(), NOW()),
            -- Audio and Video subcategories (parent id=3200)
        ('category', 3210, 'en', 'name', 'TVs and Projectors', true, true, NOW(), NOW()),
        ('category', 3215, 'en', 'name', 'Headphones', true, true, NOW(), NOW()),
        ('category', 3220, 'en', 'name', 'Speakers, Sound Systems, Subwoofers', true, true, NOW(), NOW()),
        ('category', 3225, 'en', 'name', 'Accessories', true, true, NOW(), NOW()),
        ('category', 3230, 'en', 'name', 'Stereo Systems, Boomboxes', true, true, NOW(), NOW()),
        ('category', 3235, 'en', 'name', 'Amplifiers and Receivers', true, true, NOW(), NOW()),
        ('category', 3240, 'en', 'name', 'Video Cameras', true, true, NOW(), NOW()),
        ('category', 3245, 'en', 'name', 'Video, DVD and Blu-ray Players', true, true, NOW(), NOW()),
        ('category', 3250, 'en', 'name', 'Cables and Adapters', true, true, NOW(), NOW()),
        ('category', 3255, 'en', 'name', 'Music and Movies', true, true, NOW(), NOW()),
        ('category', 3260, 'en', 'name', 'Microphones', true, true, NOW(), NOW()),
        ('category', 3265, 'en', 'name', 'MP3 Players', true, true, NOW(), NOW()),

    ('category', 3300, 'en', 'name', 'Computer Equipment', true, true, NOW(), NOW()),
    -- Computer Equipment subcategories (parent id=3300)
        ('category', 3310, 'en', 'name', 'Desktop Computers', true, true, NOW(), NOW()),
        ('category', 3320, 'en', 'name', 'All-in-Ones', true, true, NOW(), NOW()),
        ('category', 3330, 'en', 'name', 'Components', true, true, NOW(), NOW()),
        -- Components subcategories (parent id=3330)
            ('category', 3331, 'en', 'name', 'CD, DVD and Blu-ray Drives', true, true, NOW(), NOW()),
            ('category', 3332, 'en', 'name', 'Power Supplies', true, true, NOW(), NOW()),
            ('category', 3333, 'en', 'name', 'Video Cards', true, true, NOW(), NOW()),
            ('category', 3334, 'en', 'name', 'Hard Drives', true, true, NOW(), NOW()),
            ('category', 3335, 'en', 'name', 'Sound Cards', true, true, NOW(), NOW()),
            ('category', 3336, 'en', 'name', 'Controllers', true, true, NOW(), NOW()),
            ('category', 3337, 'en', 'name', 'Cases', true, true, NOW(), NOW()),
            ('category', 3338, 'en', 'name', 'Motherboards', true, true, NOW(), NOW()),
            ('category', 3339, 'en', 'name', 'RAM', true, true, NOW(), NOW()),
            ('category', 3340, 'en', 'name', 'Processors', true, true, NOW(), NOW()),
            ('category', 3341, 'en', 'name', 'Cooling Systems', true, true, NOW(), NOW()),

        ('category', 3360, 'en', 'name', 'Monitors and Parts', true, true, NOW(), NOW()),
        ('category', 3365, 'en', 'name', 'Network Equipment', true, true, NOW(), NOW()),
        ('category', 3370, 'en', 'name', 'Keyboards and Mice', true, true, NOW(), NOW()),
        ('category', 3375, 'en', 'name', 'Joysticks and Steering Wheels', true, true, NOW(), NOW()),
        ('category', 3380, 'en', 'name', 'Flash Drives and Memory Cards', true, true, NOW(), NOW()),
        ('category', 3385, 'en', 'name', 'Webcams', true, true, NOW(), NOW()),
        ('category', 3390, 'en', 'name', 'TV Tuners', true, true, NOW(), NOW()),

    ('category', 3500, 'en', 'name', 'Games, Consoles and Software', true, true, NOW(), NOW()),
            -- Games, Consoles and Software subcategories (parent id=3500)
        ('category', 3510, 'en', 'name', 'Game Consoles', true, true, NOW(), NOW()),
        ('category', 3520, 'en', 'name', 'Console Games', true, true, NOW(), NOW()),
        ('category', 3530, 'en', 'name', 'Computer Games', true, true, NOW(), NOW()),
        ('category', 3540, 'en', 'name', 'Software', true, true, NOW(), NOW()),

    ('category', 3600, 'en', 'name', 'Laptops', true, true, NOW(), NOW()),
    ('category', 3700, 'en', 'name', 'Photography Equipment', true, true, NOW(), NOW()),
            -- Photography Equipment subcategories (parent id=3700)
        ('category', 3710, 'en', 'name', 'Equipment and Accessories', true, true, NOW(), NOW()),
        ('category', 3720, 'en', 'name', 'Lenses', true, true, NOW(), NOW()),
        ('category', 3730, 'en', 'name', 'Compact Cameras', true, true, NOW(), NOW()),
        ('category', 3740, 'en', 'name', 'Film Cameras', true, true, NOW(), NOW()),
        ('category', 3750, 'en', 'name', 'DSLR Cameras', true, true, NOW(), NOW()),

    ('category', 3800, 'en', 'name', 'Tablets and E-readers', true, true, NOW(), NOW()),
            -- Tablets and E-readers subcategories (parent id=3800)
        ('category', 3810, 'en', 'name', 'Tablets', true, true, NOW(), NOW()),
        ('category', 3820, 'en', 'name', 'E-readers', true, true, NOW(), NOW()),
        ('category', 3830, 'en', 'name', 'Accessories', true, true, NOW(), NOW()),

    ('category', 3900, 'en', 'name', 'Office Equipment and Supplies', true, true, NOW(), NOW()),
            -- Office Equipment and Supplies subcategories (parent id=3900)
        ('category', 3910, 'en', 'name', 'MFPs, Copiers and Scanners', true, true, NOW(), NOW()),
        ('category', 3920, 'en', 'name', 'Printers', true, true, NOW(), NOW()),
        ('category', 3930, 'en', 'name', 'Stationery', true, true, NOW(), NOW()),
        ('category', 3940, 'en', 'name', 'UPS, Surge Protectors', true, true, NOW(), NOW()),
        ('category', 3950, 'en', 'name', 'Telephony', true, true, NOW(), NOW()),
        ('category', 3960, 'en', 'name', 'Paper Shredders', true, true, NOW(), NOW()),
        ('category', 3970, 'en', 'name', 'Consumables', true, true, NOW(), NOW()),

    ('category', 4100, 'en', 'name', 'Home Appliances', true, true, NOW(), NOW()),
        -- Home Appliances subcategories (parent id=4100)
        ('category', 4110, 'en', 'name', 'Kitchen Appliances', true, true, NOW(), NOW()),
                    -- Kitchen Appliances subcategories (parent id=4110)
            ('category', 4111, 'en', 'name', 'Exhaust Hoods', true, true, NOW(), NOW()),
            ('category', 4112, 'en', 'name', 'Small Kitchen Appliances', true, true, NOW(), NOW()),
            ('category', 4113, 'en', 'name', 'Microwave Ovens', true, true, NOW(), NOW()),
            ('category', 4114, 'en', 'name', 'Stoves and Ovens', true, true, NOW(), NOW()),
            ('category', 4115, 'en', 'name', 'Dishwashers', true, true, NOW(), NOW()),
            ('category', 4116, 'en', 'name', 'Refrigerators and Freezers', true, true, NOW(), NOW()),

        ('category', 4120, 'en', 'name', 'Home Appliances', true, true, NOW(), NOW()),
                    -- Home Appliances subcategories (parent id=4120)
            ('category', 4121, 'en', 'name', 'Vacuum Cleaners and Parts', true, true, NOW(), NOW()),
            ('category', 4122, 'en', 'name', 'Washing and Drying Machines', true, true, NOW(), NOW()),
            ('category', 4123, 'en', 'name', 'Irons', true, true, NOW(), NOW()),
            ('category', 4124, 'en', 'name', 'Sewing Equipment', true, true, NOW(), NOW()),

        ('category', 4130, 'en', 'name', 'Climate Equipment', true, true, NOW(), NOW()),
                    -- Climate Equipment subcategories (parent id=4130)
            ('category', 4131, 'en', 'name', 'Fans', true, true, NOW(), NOW()),
            ('category', 4132, 'en', 'name', 'Air Conditioners and Parts', true, true, NOW(), NOW()),
            ('category', 4133, 'en', 'name', 'Heaters', true, true, NOW(), NOW()),
            ('category', 4134, 'en', 'name', 'Air Purifiers', true, true, NOW(), NOW()),
            ('category', 4135, 'en', 'name', 'Thermometers and Weather Stations', true, true, NOW(), NOW()),

        ('category', 4140, 'en', 'name', 'Personal Care', true, true, NOW(), NOW()),
            -- Personal Care subcategories (parent id=4140)
            ('category', 4141, 'en', 'name', 'Shavers and Trimmers', true, true, NOW(), NOW()),
            ('category', 4142, 'en', 'name', 'Hair Clippers', true, true, NOW(), NOW()),
            ('category', 4143, 'en', 'name', 'Hair Dryers and Styling Tools', true, true, NOW(), NOW()),
            ('category', 4144, 'en', 'name', 'Epilators', true, true, NOW(), NOW()),
            
('category', 5000, 'en', 'name', 'Home and Apartment', true, true, NOW(), NOW()),
    -- Home and Apartment (parent id=5000)
    -- Renovation and Construction
    ('category', 5100, 'en', 'name', 'Renovation and Construction', true, true, NOW(), NOW()),
        -- Renovation and Construction subcategories (parent id=5100)
        ('category', 5110, 'en', 'name', 'Doors', true, true, NOW(), NOW()),
        ('category', 5115, 'en', 'name', 'Tools', true, true, NOW(), NOW()),
        ('category', 5120, 'en', 'name', 'Fireplaces and Heaters', true, true, NOW(), NOW()),
        ('category', 5125, 'en', 'name', 'Windows and Balconies', true, true, NOW(), NOW()),
        ('category', 5130, 'en', 'name', 'Ceilings', true, true, NOW(), NOW()),
        ('category', 5135, 'en', 'name', 'Garden and Cottage', true, true, NOW(), NOW()),
        ('category', 5140, 'en', 'name', 'Plumbing, Water Supply and Sauna', true, true, NOW(), NOW()),
        ('category', 5145, 'en', 'name', 'Ready-made Buildings and Log Houses', true, true, NOW(), NOW()),
        ('category', 5150, 'en', 'name', 'Gates, Fences and Barriers', true, true, NOW(), NOW()),
        ('category', 5155, 'en', 'name', 'Security and Alarm Systems', true, true, NOW(), NOW()),

    ('category', 5200, 'en', 'name', 'Furniture and Interior', true, true, NOW(), NOW()),
        -- Furniture and Interior subcategories (parent id=5200)
        ('category', 5210, 'en', 'name', 'Beds, Sofas and Armchairs', true, true, NOW(), NOW()),
        ('category', 5215, 'en', 'name', 'Textiles and Rugs', true, true, NOW(), NOW()),
        ('category', 5220, 'en', 'name', 'Lighting', true, true, NOW(), NOW()),
        ('category', 5225, 'en', 'name', 'Computer Desks and Chairs', true, true, NOW(), NOW()),
        ('category', 5230, 'en', 'name', 'Wardrobes, Chests and Shelving', true, true, NOW(), NOW()),
        ('category', 5235, 'en', 'name', 'Kitchen Sets', true, true, NOW(), NOW()),
        ('category', 5240, 'en', 'name', 'Tables and Chairs', true, true, NOW(), NOW()),
        ('category', 5250, 'en', 'name', 'Indoor Plants', true, true, NOW(), NOW()),
            -- Indoor Plants subcategories (parent id=5250)
            ('category', 5251, 'en', 'name', 'Decorative Foliage Plants', true, true, NOW(), NOW()),
            ('category', 5252, 'en', 'name', 'Flowering Plants', true, true, NOW(), NOW()),
            ('category', 5253, 'en', 'name', 'Palms and Ficus', true, true, NOW(), NOW()),
            ('category', 5254, 'en', 'name', 'Cacti and Succulents', true, true, NOW(), NOW()),

    ('category', 5300, 'en', 'name', 'Food Products', true, true, NOW(), NOW()),
            -- Food Products subcategories (parent id=5300)
        ('category', 5310, 'en', 'name', 'Tea, Coffee, Cocoa', true, true, NOW(), NOW()),
        ('category', 5315, 'en', 'name', 'Beverages', true, true, NOW(), NOW()),
        ('category', 5320, 'en', 'name', 'Fish, Seafood, Caviar', true, true, NOW(), NOW()),
        ('category', 5325, 'en', 'name', 'Meat, Poultry, Offal', true, true, NOW(), NOW()),
        ('category', 5330, 'en', 'name', 'Confectionery', true, true, NOW(), NOW()),
        ('category', 5340, 'en', 'name', 'Rakia and Wine', true, true, NOW(), NOW()),
                    -- Rakia and Wine subcategories (parent id=5340)
            ('category', 5341, 'en', 'name', 'Plum Rakia', true, true, NOW(), NOW()),
            ('category', 5342, 'en', 'name', 'Grape Rakia', true, true, NOW(), NOW()),
            ('category', 5343, 'en', 'name', 'Fruit Rakia', true, true, NOW(), NOW()),
            ('category', 5344, 'en', 'name', 'Homemade Wine', true, true, NOW(), NOW()),

        ('category', 5350, 'en', 'name', 'Homemade Cheese', true, true, NOW(), NOW()),
        ('category', 5360, 'en', 'name', 'Kaymak', true, true, NOW(), NOW()),
        ('category', 5370, 'en', 'name', 'Ajvar', true, true, NOW(), NOW()),

    ('category', 5400, 'en', 'name', 'Kitchenware and Kitchen Goods', true, true, NOW(), NOW()),
        -- Kitchenware and Kitchen Goods subcategories (parent id=5400)
        ('category', 5405, 'en', 'name', 'Tableware', true, true, NOW(), NOW()),
        ('category', 5410, 'en', 'name', 'Kitchen Items', true, true, NOW(), NOW()),
        ('category', 5415, 'en', 'name', 'Table Setting', true, true, NOW(), NOW()),
        ('category', 5420, 'en', 'name', 'Cooking', true, true, NOW(), NOW()),
        ('category', 5425, 'en', 'name', 'Food Storage', true, true, NOW(), NOW()),
        ('category', 5430, 'en', 'name', 'Beverage Preparation', true, true, NOW(), NOW()),
        ('category', 5435, 'en', 'name', 'Household Goods', true, true, NOW(), NOW()),

('category', 6000, 'en', 'name', 'Garden Supplies', true, true, NOW(), NOW()),
    -- Garden Supplies (parent id=6000, subcategories 51-60)
    ('category', 6050, 'en', 'name', 'Garden Furniture', true, true, NOW(), NOW()),
    ('category', 6100, 'en', 'name', 'Garden Tools', true, true, NOW(), NOW()),
    ('category', 6200, 'en', 'name', 'Seeds and Seedlings', true, true, NOW(), NOW()),
    ('category', 6250, 'en', 'name', 'BBQ and Accessories', true, true, NOW(), NOW()),
    ('category', 6300, 'en', 'name', 'Swimming Pools and Equipment', true, true, NOW(), NOW()),
    ('category', 6350, 'en', 'name', 'Irrigation Systems', true, true, NOW(), NOW()),
('category', 6400, 'en', 'name', 'Composting', true, true, NOW(), NOW()),
    ('category', 6450, 'en', 'name', 'Greenhouses and Hotbeds', true, true, NOW(), NOW()),
    ('category', 6500, 'en', 'name', 'Fertilizers and Soils', true, true, NOW(), NOW()),
    ('category', 6550, 'en', 'name', 'Lighting', true, true, NOW(), NOW()),
    ('category', 6600, 'en', 'name', 'Interior Design', true, true, NOW(), NOW()),
    ('category', 6650, 'en', 'name', 'Plants and Seeds', true, true, NOW(), NOW()),
    ('category', 6700, 'en', 'name', 'Garden and Vegetable Garden', true, true, NOW(), NOW()),
    ('category', 6750, 'en', 'name', 'Garden Plants', true, true, NOW(), NOW()),
            -- Garden Plants subcategories (parent id=6750)
        ('category', 6751, 'en', 'name', 'Decorative Shrubs and Trees', true, true, NOW(), NOW()),
        ('category', 6752, 'en', 'name', 'Coniferous Plants', true, true, NOW(), NOW()),
        ('category', 6753, 'en', 'name', 'Perennial Plants', true, true, NOW(), NOW()),
        ('category', 6754, 'en', 'name', 'Fruit Plants', true, true, NOW(), NOW()),
        ('category', 6755, 'en', 'name', 'Lawn', true, true, NOW(), NOW()),
        ('category', 6756, 'en', 'name', 'Greens and Herbs', true, true, NOW(), NOW()),

    ('category', 6850, 'en', 'name', 'Seeds, Bulbs, Tubers', true, true, NOW(), NOW()),
    ('category', 6900, 'en', 'name', 'Plant Care Products', true, true, NOW(), NOW()),
            -- Plant Care Products subcategories (parent id=6900)
        ('category', 6901, 'en', 'name', 'Soils and Substrates', true, true, NOW(), NOW()),
        ('category', 6902, 'en', 'name', 'Fertilizers', true, true, NOW(), NOW()),
        ('category', 6903, 'en', 'name', 'Pest and Weed Control', true, true, NOW(), NOW()),
        ('category', 6904, 'en', 'name', 'Pots and Planters', true, true, NOW(), NOW()),
        ('category', 6905, 'en', 'name', 'Plant Grow Lights', true, true, NOW(), NOW()),
        ('category', 6906, 'en', 'name', 'Moisture Meters', true, true, NOW(), NOW()),

    ('category', 6950, 'en', 'name', 'Panels and Artificial Plants', true, true, NOW(), NOW()),

('category', 7000, 'en', 'name', 'Hobbies and Leisure', true, true, NOW(), NOW()),
        -- Hobbies and Leisure (parent id=7000)
    ('category', 7050, 'en', 'name', 'Musical Instruments', true, true, NOW(), NOW()),
        -- Musical Instruments (parent id=7050)
        ('category', 7055, 'en', 'name', 'String Instruments', true, true, NOW(), NOW()),
        ('category', 7060, 'en', 'name', 'Pianos and Keyboards', true, true, NOW(), NOW()),
        ('category', 7065, 'en', 'name', 'Percussion Instruments', true, true, NOW(), NOW()),
        ('category', 7070, 'en', 'name', 'Wind Instruments', true, true, NOW(), NOW()),
        ('category', 7075, 'en', 'name', 'Accordions and Harmonicas', true, true, NOW(), NOW()),
        ('category', 7080, 'en', 'name', 'Audio Equipment', true, true, NOW(), NOW()),
        ('category', 7085, 'en', 'name', 'Instrument Accessories', true, true, NOW(), NOW()),

    ('category', 7100, 'en', 'name', 'Books and Magazines', true, true, NOW(), NOW()),
            -- Books and Magazines subcategories (parent id=7100)
        ('category', 7105, 'en', 'name', 'Magazines, Newspapers, Brochures', true, true, NOW(), NOW()),
        ('category', 7115, 'en', 'name', 'Books', true, true, NOW(), NOW()),
        ('category', 7130, 'en', 'name', 'Educational Literature', true, true, NOW(), NOW()),
    ('category', 7150, 'en', 'name', 'Sports Equipment', true, true, NOW(), NOW()),
    ('category', 7250, 'en', 'name', 'Collecting', true, true, NOW(), NOW()),
            -- Collecting subcategories (parent id=7250)
        ('category', 7251, 'en', 'name', 'Banknotes', true, true, NOW(), NOW()),
        ('category', 7252, 'en', 'name', 'Tickets', true, true, NOW(), NOW()),
        ('category', 7253, 'en', 'name', 'Celebrity Items, Autographs', true, true, NOW(), NOW()),
        ('category', 7254, 'en', 'name', 'Military Items', true, true, NOW(), NOW()),
        ('category', 7255, 'en', 'name', 'Vinyl Records', true, true, NOW(), NOW()),
        ('category', 7256, 'en', 'name', 'Documents', true, true, NOW(), NOW()),
        ('category', 7257, 'en', 'name', 'Tokens, Medals, Badges', true, true, NOW(), NOW()),
        ('category', 7258, 'en', 'name', 'Games', true, true, NOW(), NOW()),
        ('category', 7259, 'en', 'name', 'Calendars', true, true, NOW(), NOW()),
        ('category', 7261, 'en', 'name', 'Paintings', true, true, NOW(), NOW()),
        ('category', 7262, 'en', 'name', 'Stamps', true, true, NOW(), NOW()),
        ('category', 7263, 'en', 'name', 'Models', true, true, NOW(), NOW()),
        ('category', 7264, 'en', 'name', 'Coins', true, true, NOW(), NOW()),

    ('category', 7300, 'en', 'name', 'Art Objects', true, true, NOW(), NOW()),
    ('category', 7400, 'en', 'name', 'Bicycles', true, true, NOW(), NOW()),
            -- Bicycles subcategories (parent id=7400)
        ('category', 7410, 'en', 'name', 'BMX', true, true, NOW(), NOW()),
        ('category', 7415, 'en', 'name', 'City Bikes', true, true, NOW(), NOW()),
        ('category', 7420, 'en', 'name', 'Road Bikes', true, true, NOW(), NOW()),
        ('category', 7425, 'en', 'name', 'Children Bikes', true, true, NOW(), NOW()),
        ('category', 7430, 'en', 'name', 'Mountain Bikes', true, true, NOW(), NOW()),
        ('category', 7435, 'en', 'name', 'Parts and Accessories', true, true, NOW(), NOW()),

    ('category', 7500, 'en', 'name', 'Hunting and Fishing', true, true, NOW(), NOW()),
        -- Hunting and Fishing subcategories (parent id=7500)
        ('category', 7510, 'en', 'name', 'Knives, Multi-tools, Axes', true, true, NOW(), NOW()),
        ('category', 7520, 'en', 'name', 'Hunting', true, true, NOW(), NOW()),
                    -- Hunting subcategories (parent id=7520)
            ('category', 7521, 'en', 'name', 'Scopes', true, true, NOW(), NOW()),
            ('category', 7522, 'en', 'name', 'Scope Accessories', true, true, NOW(), NOW()),
            ('category', 7523, 'en', 'name', 'Monoculars, Binoculars, Rangefinders', true, true, NOW(), NOW()),

        ('category', 7530, 'en', 'name', 'Fishing', true, true, NOW(), NOW()),
            -- Fishing subcategories (parent id=7530)
            ('category', 7531, 'en', 'name', 'Fishing Rods, Spinning Rods and Reels', true, true, NOW(), NOW()),
            ('category', 7551, 'en', 'name', 'Lures and Tackle', true, true, NOW(), NOW()),
            ('category', 7571, 'en', 'name', 'Fish Finders and Equipment', true, true, NOW(), NOW()),

    ('category', 7650, 'en', 'name', 'Camping', true, true, NOW(), NOW()),
    ('category', 7700, 'en', 'name', 'Antiques', true, true, NOW(), NOW()),
    ('category', 7750, 'en', 'name', 'Tickets, Events and Travel', true, true, NOW(), NOW()),
        -- Tickets, Events and Travel subcategories (parent id=7750)
        ('category', 7751, 'en', 'name', 'Cards, Coupons', true, true, NOW(), NOW()),
        ('category', 7752, 'en', 'name', 'Concerts', true, true, NOW(), NOW()),
        ('category', 7753, 'en', 'name', 'Travel', true, true, NOW(), NOW()),
        ('category', 7754, 'en', 'name', 'Sports', true, true, NOW(), NOW()),
        ('category', 7755, 'en', 'name', 'Theater, Opera, Ballet', true, true, NOW(), NOW()),
        ('category', 7756, 'en', 'name', 'Circus, Cinema', true, true, NOW(), NOW()),
        ('category', 7758, 'en', 'name', 'Shows, Musicals', true, true, NOW(), NOW()),

    ('category', 7800, 'en', 'name', 'Sports', true, true, NOW(), NOW()),
            -- Sports subcategories (parent id=7800)
        ('category', 7805, 'en', 'name', 'Billiards and Bowling', true, true, NOW(), NOW()),
        ('category', 7810, 'en', 'name', 'Diving and Water Sports', true, true, NOW(), NOW()),
        ('category', 7815, 'en', 'name', 'Martial Arts', true, true, NOW(), NOW()),
        ('category', 7820, 'en', 'name', 'Winter Sports', true, true, NOW(), NOW()),
        ('category', 7825, 'en', 'name', 'Ball Games', true, true, NOW(), NOW()),
        ('category', 7830, 'en', 'name', 'Board Games', true, true, NOW(), NOW()),
        ('category', 7835, 'en', 'name', 'Paintball and Airsoft', true, true, NOW(), NOW()),
        ('category', 7840, 'en', 'name', 'Roller Skating and Skateboarding', true, true, NOW(), NOW()),
        ('category', 7845, 'en', 'name', 'Tennis, Badminton, Ping-Pong', true, true, NOW(), NOW()),
        ('category', 7850, 'en', 'name', 'Tourism and Outdoor Recreation', true, true, NOW(), NOW()),
        ('category', 7855, 'en', 'name', 'Fitness and Exercise Equipment', true, true, NOW(), NOW()),
        ('category', 7860, 'en', 'name', 'Sports Nutrition', true, true, NOW(), NOW()),

    -- Traditional Crafts and Souvenirs
    ('category', 7865, 'en', 'name', 'Folk Crafts and Handicrafts', true, true, NOW(), NOW()),
        -- Folk Crafts subcategories (parent id=7865)
        ('category', 7866, 'en', 'name', 'Opanci', true, true, NOW(), NOW()),
        ('category', 7867, 'en', 'name', 'Ceramics', true, true, NOW(), NOW()),
        ('category', 7868, 'en', 'name', 'Embroidery', true, true, NOW(), NOW()),
        ('category', 7869, 'en', 'name', 'Weaving', true, true, NOW(), NOW()),
        ('category', 7870, 'en', 'name', 'Folk Instruments', true, true, NOW(), NOW()),
        ('category', 7871, 'en', 'name', 'Woodworking', true, true, NOW(), NOW()),

    -- Agricultural Categories
    ('category', 7900, 'en', 'name', 'Beekeeping', true, true, NOW(), NOW()),
        -- Beekeeping subcategories (parent id=7900)
        ('category', 7910, 'en', 'name', 'Honey', true, true, NOW(), NOW()),
        ('category', 7920, 'en', 'name', 'Beeswax', true, true, NOW(), NOW()),
        ('category', 7930, 'en', 'name', 'Propolis', true, true, NOW(), NOW()),
        ('category', 7935, 'en', 'name', 'Beekeeping Equipment', true, true, NOW(), NOW()),
        ('category', 7945, 'en', 'name', 'Bees', true, true, NOW(), NOW()),

    -- Tourism Services
    ('category', 7950, 'en', 'name', 'Rural Tourism', true, true, NOW(), NOW()),
            -- Rural Tourism subcategories (parent id=7950)
        ('category', 7951, 'en', 'name', 'Ethno-villages', true, true, NOW(), NOW()),
        ('category', 7952, 'en', 'name', 'Wine Tours', true, true, NOW(), NOW()),
        ('category', 7953, 'en', 'name', 'Agritourism', true, true, NOW(), NOW()),
        ('category', 7954, 'en', 'name', 'Mountain Tourism', true, true, NOW(), NOW()),

('category', 8000, 'en', 'name', 'Animals', true, true, NOW(), NOW()),
    -- Animals categories (parent id=8000)
    ('category', 8050, 'en', 'name', 'Dogs', true, true, NOW(), NOW()),
    ('category', 8100, 'en', 'name', 'Cats', true, true, NOW(), NOW()),
    ('category', 8150, 'en', 'name', 'Birds', true, true, NOW(), NOW()),
    ('category', 8200, 'en', 'name', 'Aquarium', true, true, NOW(), NOW()),
            -- Aquarium subcategories (parent id=8200)
        ('category', 8205, 'en', 'name', 'Aquariums', true, true, NOW(), NOW()),
        ('category', 8210, 'en', 'name', 'Fish', true, true, NOW(), NOW()),
        ('category', 8215, 'en', 'name', 'Other Aquarium Animals', true, true, NOW(), NOW()),
        ('category', 8220, 'en', 'name', 'Equipment', true, true, NOW(), NOW()),
        ('category', 8225, 'en', 'name', 'Plants', true, true, NOW(), NOW()),
        ('category', 8230, 'en', 'name', 'Aquarium Furniture', true, true, NOW(), NOW()),
        ('category', 8235, 'en', 'name', 'Marine Aquariums', true, true, NOW(), NOW()),

    ('category', 8250, 'en', 'name', 'Other Animals', true, true, NOW(), NOW()),
        -- Other Animals subcategories (parent id=8250)
        ('category', 8251, 'en', 'name', 'Amphibians', true, true, NOW(), NOW()),
        ('category', 8252, 'en', 'name', 'Rodents', true, true, NOW(), NOW()),
        ('category', 8253, 'en', 'name', 'Rabbits', true, true, NOW(), NOW()),
        ('category', 8254, 'en', 'name', 'Horses', true, true, NOW(), NOW()),
        ('category', 8255, 'en', 'name', 'Reptiles', true, true, NOW(), NOW()),
        ('category', 8256, 'en', 'name', 'Farm Animals', true, true, NOW(), NOW()),
        ('category', 8257, 'en', 'name', 'Poultry', true, true, NOW(), NOW()),
        ('category', 8258, 'en', 'name', 'Pet Supplies', true, true, NOW(), NOW()),

('category', 8500, 'en', 'name', 'Ready Business and Equipment', true, true, NOW(), NOW()),


('category', 9000, 'en', 'name', 'Jobs', true, true, NOW(), NOW()),
    ('category', 9050, 'en', 'name', 'Vacancies', true, true, NOW(), NOW()),
    ('category', 9100, 'en', 'name', 'Resumes', true, true, NOW(), NOW()),
    ('category', 9150, 'en', 'name', 'Remote Work', true, true, NOW(), NOW()),
    ('category', 9200, 'en', 'name', 'Partnership and Collaboration', true, true, NOW(), NOW()),
    ('category', 9250, 'en', 'name', 'Training and Internship', true, true, NOW(), NOW()),
    ('category', 9300, 'en', 'name', 'Seasonal Work', true, true, NOW(), NOW()),
        -- Seasonal Work subcategories (parent id=9300)
        ('category', 9310, 'en', 'name', 'Harvesting', true, true, NOW(), NOW()),
        ('category', 9315, 'en', 'name', 'Vineyard Work', true, true, NOW(), NOW()),
        ('category', 9320, 'en', 'name', 'Seasonal Construction Work', true, true, NOW(), NOW()),

('category', 9500, 'en', 'name', 'Clothing, Footwear, Accessories', true, true, NOW(), NOW()),

-- Children's Goods and Toys (new category)
('category', 9700, 'en', 'name', 'Childrens Goods and Toys', true, true, NOW(), NOW()),
    -- Children's Goods and Toys subcategories (parent id=9700)
    ('category', 9705, 'en', 'name', 'Baby Strollers', true, true, NOW(), NOW()),
    ('category', 9710, 'en', 'name', 'Childrens Furniture', true, true, NOW(), NOW()),
    ('category', 9715, 'en', 'name', 'Bicycles and Scooters', true, true, NOW(), NOW()),
    ('category', 9720, 'en', 'name', 'Feeding Products', true, true, NOW(), NOW()),
    ('category', 9725, 'en', 'name', 'Car Seats', true, true, NOW(), NOW()),
    ('category', 9730, 'en', 'name', 'Toys', true, true, NOW(), NOW()),
    ('category', 9735, 'en', 'name', 'Bedding', true, true, NOW(), NOW()),
    ('category', 9740, 'en', 'name', 'Bathing Products', true, true, NOW(), NOW()),
    ('category', 9745, 'en', 'name', 'School Supplies', true, true, NOW(), NOW()),
    ('category', 9750, 'en', 'name', 'Childrens Clothing and Footwear, Accessories', true, true, NOW(), NOW()),

('category', 9999, 'en', 'name', 'Other', true, true, NOW(), NOW()),
('category', 10000, 'en', 'name', 'Security', true, true, NOW(), NOW()),
    -- Children's Goods and Toys subcategories (parent id=10000)
    ('category', 10100, 'en', 'name', 'Video Surveillance', true, true, NOW(), NOW()),
        -- Video Surveillance subcategories (parent id=10100)
         ('category', 10410, 'en', 'name', 'Connectors and Video Baluns', true, true, NOW(), NOW());
 



-- Update sequence for translations
SELECT setval('translations_id_seq', (SELECT MAX(id) FROM translations), true);