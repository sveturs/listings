-- Добавляем переводы для товаров витрин ID 225-242
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text, is_machine_translated)
VALUES 
    -- Product 225: Svež hleb
    ('storefront_product', 225, 'en', 'name', 'Fresh Bread', true),
    ('storefront_product', 225, 'en', 'description', 'Homemade white bread, baked daily', true),
    ('storefront_product', 225, 'ru', 'name', 'Свежий хлеб', true),
    ('storefront_product', 225, 'ru', 'description', 'Домашний белый хлеб, выпекается ежедневно', true),
    ('storefront_product', 225, 'sr', 'name', 'Svež hleb', true),
    ('storefront_product', 225, 'sr', 'description', 'Domaći beli hleb, peče se svakodnevno', true),
    
    -- Product 226: Integralni hleb
    ('storefront_product', 226, 'en', 'name', 'Whole Wheat Bread', true),
    ('storefront_product', 226, 'en', 'description', 'Healthy whole wheat bread with seeds', true),
    ('storefront_product', 226, 'ru', 'name', 'Цельнозерновой хлеб', true),
    ('storefront_product', 226, 'ru', 'description', 'Полезный цельнозерновой хлеб с семечками', true),
    ('storefront_product', 226, 'sr', 'name', 'Integralni hleb', true),
    ('storefront_product', 226, 'sr', 'description', 'Zdrav integralni hleb sa semenkama', true),
    
    -- Product 227: Kroasan
    ('storefront_product', 227, 'en', 'name', 'Croissant', true),
    ('storefront_product', 227, 'en', 'description', 'Butter croissant, fresh from France', true),
    ('storefront_product', 227, 'ru', 'name', 'Круассан', true),
    ('storefront_product', 227, 'ru', 'description', 'Масляный круассан, свежий из Франции', true),
    ('storefront_product', 227, 'sr', 'name', 'Kroasan', true),
    ('storefront_product', 227, 'sr', 'description', 'Masleni kroasan, svež iz Francuske', true),
    
    -- Product 228: Đevrek
    ('storefront_product', 228, 'en', 'name', 'Bagel', true),
    ('storefront_product', 228, 'en', 'description', 'Traditional bagel with sesame seeds', true),
    ('storefront_product', 228, 'ru', 'name', 'Бублик', true),
    ('storefront_product', 228, 'ru', 'description', 'Традиционный бублик с кунжутом', true),
    ('storefront_product', 228, 'sr', 'name', 'Đevrek', true),
    ('storefront_product', 228, 'sr', 'description', 'Tradicionalni đevrek sa susumom', true),
    
    -- Product 229: Gaming stolica Razer
    ('storefront_product', 229, 'en', 'name', 'Razer Gaming Chair', true),
    ('storefront_product', 229, 'en', 'description', 'Professional gaming chair with lumbar support', true),
    ('storefront_product', 229, 'ru', 'name', 'Игровое кресло Razer', true),
    ('storefront_product', 229, 'ru', 'description', 'Профессиональное игровое кресло с поддержкой поясницы', true),
    ('storefront_product', 229, 'sr', 'name', 'Gaming stolica Razer', true),
    ('storefront_product', 229, 'sr', 'description', 'Profesionalna gejming stolica sa lumbalnom podrškom', true),
    
    -- Product 230: Mehanička tastatura RGB
    ('storefront_product', 230, 'en', 'name', 'RGB Mechanical Keyboard', true),
    ('storefront_product', 230, 'en', 'description', 'Gaming keyboard with Cherry MX switches', true),
    ('storefront_product', 230, 'ru', 'name', 'Механическая клавиатура RGB', true),
    ('storefront_product', 230, 'ru', 'description', 'Игровая клавиатура с переключателями Cherry MX', true),
    ('storefront_product', 230, 'sr', 'name', 'Mehanička tastatura RGB', true),
    ('storefront_product', 230, 'sr', 'description', 'Gejming tastatura sa Cherry MX prekidačima', true),
    
    -- Product 231: Gaming miš Logitech G502
    ('storefront_product', 231, 'en', 'name', 'Logitech G502 Gaming Mouse', true),
    ('storefront_product', 231, 'en', 'description', 'Professional gaming mouse with adjustable DPI', true),
    ('storefront_product', 231, 'ru', 'name', 'Игровая мышь Logitech G502', true),
    ('storefront_product', 231, 'ru', 'description', 'Профессиональная игровая мышь с регулируемым DPI', true),
    ('storefront_product', 231, 'sr', 'name', 'Gaming miš Logitech G502', true),
    ('storefront_product', 231, 'sr', 'description', 'Profesionalni gejming miš sa podesivim DPI', true),
    
    -- Product 232: Monitor 4K 144Hz
    ('storefront_product', 232, 'en', 'name', '4K 144Hz Monitor', true),
    ('storefront_product', 232, 'en', 'description', '27 inch gaming monitor with HDR support', true),
    ('storefront_product', 232, 'ru', 'name', 'Монитор 4K 144Hz', true),
    ('storefront_product', 232, 'ru', 'description', '27-дюймовый игровой монитор с поддержкой HDR', true),
    ('storefront_product', 232, 'sr', 'name', 'Monitor 4K 144Hz', true),
    ('storefront_product', 232, 'sr', 'description', '27 inčni gejming monitor sa HDR podrškom', true),
    
    -- Product 233: Jabuke Crveni Delišes
    ('storefront_product', 233, 'en', 'name', 'Red Delicious Apples', true),
    ('storefront_product', 233, 'en', 'description', 'Fresh sweet apples from local orchard', true),
    ('storefront_product', 233, 'ru', 'name', 'Яблоки Ред Делишес', true),
    ('storefront_product', 233, 'ru', 'description', 'Свежие сладкие яблоки из местного сада', true),
    ('storefront_product', 233, 'sr', 'name', 'Jabuke Crveni Delišes', true),
    ('storefront_product', 233, 'sr', 'description', 'Sveže slatke jabuke iz lokalnog voćnjaka', true),
    
    -- Product 234: Paradajz domaći
    ('storefront_product', 234, 'en', 'name', 'Homegrown Tomatoes', true),
    ('storefront_product', 234, 'en', 'description', 'Organic tomatoes from greenhouse', true),
    ('storefront_product', 234, 'ru', 'name', 'Домашние помидоры', true),
    ('storefront_product', 234, 'ru', 'description', 'Органические помидоры из теплицы', true),
    ('storefront_product', 234, 'sr', 'name', 'Paradajz domaći', true),
    ('storefront_product', 234, 'sr', 'description', 'Organski paradajz iz plastenika', true),
    
    -- Product 235: Krompir
    ('storefront_product', 235, 'en', 'name', 'Potatoes', true),
    ('storefront_product', 235, 'en', 'description', 'Fresh local potatoes, perfect for frying', true),
    ('storefront_product', 235, 'ru', 'name', 'Картофель', true),
    ('storefront_product', 235, 'ru', 'description', 'Свежий местный картофель, идеальный для жарки', true),
    ('storefront_product', 235, 'sr', 'name', 'Krompir', true),
    ('storefront_product', 235, 'sr', 'description', 'Svež domaći krompir, savršen za prženje', true),
    
    -- Product 236: Luk crni
    ('storefront_product', 236, 'en', 'name', 'Red Onions', true),
    ('storefront_product', 236, 'en', 'description', 'Fresh red onions, great for salads', true),
    ('storefront_product', 236, 'ru', 'name', 'Красный лук', true),
    ('storefront_product', 236, 'ru', 'description', 'Свежий красный лук, отлично подходит для салатов', true),
    ('storefront_product', 236, 'sr', 'name', 'Luk crni', true),
    ('storefront_product', 236, 'sr', 'description', 'Svež crni luk, odličan za salate', true),
    
    -- Product 237: Nike Air Max 270
    ('storefront_product', 237, 'en', 'name', 'Nike Air Max 270', true),
    ('storefront_product', 237, 'en', 'description', 'Comfortable sports shoes with Air cushioning', true),
    ('storefront_product', 237, 'ru', 'name', 'Nike Air Max 270', true),
    ('storefront_product', 237, 'ru', 'description', 'Удобные спортивные кроссовки с воздушной подушкой', true),
    ('storefront_product', 237, 'sr', 'name', 'Nike Air Max 270', true),
    ('storefront_product', 237, 'sr', 'description', 'Udobne sportske patike sa Air jastučićem', true),
    
    -- Product 238: Adidas Ultraboost
    ('storefront_product', 238, 'en', 'name', 'Adidas Ultraboost', true),
    ('storefront_product', 238, 'en', 'description', 'Running shoes with Boost technology', true),
    ('storefront_product', 238, 'ru', 'name', 'Adidas Ultraboost', true),
    ('storefront_product', 238, 'ru', 'description', 'Беговые кроссовки с технологией Boost', true),
    ('storefront_product', 238, 'sr', 'name', 'Adidas Ultraboost', true),
    ('storefront_product', 238, 'sr', 'description', 'Patike za trčanje sa Boost tehnologijom', true),
    
    -- Product 239: Kožne čizme
    ('storefront_product', 239, 'en', 'name', 'Leather Boots', true),
    ('storefront_product', 239, 'en', 'description', 'Genuine leather winter boots', true),
    ('storefront_product', 239, 'ru', 'name', 'Кожаные сапоги', true),
    ('storefront_product', 239, 'ru', 'description', 'Зимние сапоги из натуральной кожи', true),
    ('storefront_product', 239, 'sr', 'name', 'Kožne čizme', true),
    ('storefront_product', 239, 'sr', 'description', 'Zimske čizme od prave kože', true),
    
    -- Product 240: Sportske sandale
    ('storefront_product', 240, 'en', 'name', 'Sports Sandals', true),
    ('storefront_product', 240, 'en', 'description', 'Comfortable sandals for summer', true),
    ('storefront_product', 240, 'ru', 'name', 'Спортивные сандалии', true),
    ('storefront_product', 240, 'ru', 'description', 'Удобные сандалии для лета', true),
    ('storefront_product', 240, 'sr', 'name', 'Sportske sandale', true),
    ('storefront_product', 240, 'sr', 'description', 'Udobne sandale za leto', true),
    
    -- Product 241: Muška majica polo
    ('storefront_product', 241, 'en', 'name', 'Men''s Polo Shirt', true),
    ('storefront_product', 241, 'en', 'description', 'Cotton polo shirt, various colors', true),
    ('storefront_product', 241, 'ru', 'name', 'Мужская рубашка поло', true),
    ('storefront_product', 241, 'ru', 'description', 'Хлопковая рубашка поло, разные цвета', true),
    ('storefront_product', 241, 'sr', 'name', 'Muška majica polo', true),
    ('storefront_product', 241, 'sr', 'description', 'Pamučna polo majica, razne boje', true),
    
    -- Product 242: Ženska haljina letnja
    ('storefront_product', 242, 'en', 'name', 'Women''s Summer Dress', true),
    ('storefront_product', 242, 'en', 'description', 'Light summer dress, floral pattern', true),
    ('storefront_product', 242, 'ru', 'name', 'Женское летнее платье', true),
    ('storefront_product', 242, 'ru', 'description', 'Легкое летнее платье с цветочным узором', true),
    ('storefront_product', 242, 'sr', 'name', 'Ženska haljina letnja', true),
    ('storefront_product', 242, 'sr', 'description', 'Lagana letnja haljina, cvetni dezen', true)
ON CONFLICT (entity_type, entity_id, language, field_name) DO NOTHING;