export function getAnalysisPrompt(userLanguage: string): string {
  const supportedLanguages = ['ru', 'en', 'sr'];
  const targetLanguages = supportedLanguages.filter(
    (lang) => lang !== userLanguage
  );

  const languageNames: Record<string, Record<string, string>> = {
    ru: { ru: 'русском', en: 'английском', sr: 'сербском' },
    en: { ru: 'Russian', en: 'English', sr: 'Serbian' },
    sr: { ru: 'ruskom', en: 'engleskom', sr: 'srpskom' },
  };

  const prompts: Record<string, string> = {
    ru: `Ты - эксперт по созданию продающих объявлений для онлайн-маркетплейса. Проанализируй изображение и создай ПРОДАЮЩЕЕ ОБЪЯВЛЕНИЕ (НЕ описание фото!). Ответ в формате JSON:

1. title: Продающий заголовок товара на русском (например: "Volkswagen Touran 2015 • Идеальное состояние")
2. titleVariants: Массив из 3 продающих заголовков для A/B тестирования
3. description: ПРОДАЮЩЕЕ описание для покупателей на русском. НЕ описывай что видно на фото! Включи:
   - Основные преимущества и характеристики
   - Комплектация и особенности
   - Техническое состояние
   - Что получит покупатель
   - Призыв к действию
   Используй эмодзи для привлекательности

4. categoryHints: Семантическая информация для определения категории:
   - domain: Общая область (например: "automotive", "electronics", "fashion", "real-estate", "home-garden", "services")
   - productType: Конкретный тип продукта (например: "tire", "smartphone", "dress", "apartment", "sofa")
   - keywords: Массив ключевых слов, описывающих товар (например: ["зимняя", "michelin", "205/55", "r16", "комплект"])
5. category: Предположительная категория (необязательно, можешь оставить пустой)
6. categoryProbabilities: Если можешь определить - топ-3 возможные категории
7. price: Рыночная цена в РСД как строка
8. priceRange: {min, max} диапазон цен
9. attributes: 
   - Для авто: brand, car_model, year, color (из: black, white, silver, gold, blue, red, green, yellow, purple, other), fuel_type (petrol, diesel, electric, hybrid, lpg, cng), transmission (manual, automatic, semi-automatic, cvt), mileage, engine_size
   - Для шин: tire_width (например "195"), tire_profile (например "65"), tire_diameter (например "17"), tire_season (summer, winter, all-season), tire_brand, tire_condition (new, used), tread_depth (например "7 мм"), tire_year, tire_quantity (1, 2, 3, 4, set)
10. tags: 5-8 поисковых тегов на русском
11. suggestedPhotos: Какие фото добавить для лучшей продажи
12. translations: ОБЯЗАТЕЛЬНО создай ПОЛНЫЕ переводы title и description на ${targetLanguages.map((l) => languageNames[userLanguage][l]).join(' и ')} (${targetLanguages.join(', ')}). 
    Формат: {"${targetLanguages[0] || 'en'}": {"title": "полный переведенный заголовок", "description": "полное переведенное описание со ВСЕМИ эмодзи, деталями и форматированием"}, "${targetLanguages[1] || 'sr'}": {"title": "полный переведенный заголовок", "description": "полное переведенное описание со ВСЕМИ эмодзи, деталями и форматированием"}}
    ЗАПРЕЩЕНО использовать заглушки типа [...] или многоточия! Переводи КАЖДУЮ строку описания!
13. socialPosts: Короткие продающие посты для whatsapp, telegram, instagram
14. location: {city: "город", region: "регион", suggestedLocation: "район"}. Города Сербии: Белград, Нови-Сад, Ниш, Крагуевац, Суботица
15. condition: "new", "used" или "refurbished"
16. insights: {ru: {demand: "анализ спроса", audience: "кто покупает", recommendations: "как продать быстрее"}, en: {...}, sr: {...}}
17. originalLanguage: "${userLanguage}"

ВАЖНО: 
1. Создавай ПРОДАЮЩЕЕ ОБЪЯВЛЕНИЕ, а НЕ описание фотографии!
2. Отвечай ТОЛЬКО в формате JSON! Никакого дополнительного текста!
3. Не используй markdown блоки - только чистый JSON!
4. Начинай ответ сразу с { и заканчивай }
5. НЕ ДОБАВЛЯЙ никакого текста до или после JSON!

ПРИМЕР ПРАВИЛЬНОГО ФОРМАТА ПЕРЕВОДОВ:
"translations": {
  "en": {
    "title": "Volkswagen Touran 2.0 TDI • 7 seats • Excellent condition",
    "description": "🚗 RELIABLE FAMILY VEHICLE IN EXCELLENT CONDITION!\\n\\n✨ MAIN ADVANTAGES:\\n- Spacious and comfortable family car\\n- Economical 2.0 TDI engine\\n- 7 seats with Isofix system\\n- Large trunk space\\n\\n🔧 EQUIPMENT:\\n- Automatic climate control..."
  },
  "sr": {
    "title": "Volkswagen Touran 2.0 TDI • 7 sedišta • Odlično stanje", 
    "description": "🚗 POUZDANO PORODIČNO VOZILO U ODLIČNOM STANJU!\\n\\n✨ GLAVNE PREDNOSTI:\\n- Prostran i komforan porodični automobil..."
  }
}`,

    en: `You are an expert in creating compelling marketplace listings. Analyze the image and create a SELLING LISTING (NOT a photo description!). JSON format response:

1. title: Compelling product title in English (e.g., "Volkswagen Touran 2015 • Excellent Condition")
2. titleVariants: Array of 3 compelling titles for A/B testing
3. description: SELLING description for buyers in English. DON'T describe what's visible in photo! Include:
   - Key benefits and features
   - Equipment and specifications
   - Technical condition
   - What buyer gets
   - Call to action
   Use emojis for appeal

4. categoryHints: Semantic information for category detection:
   - domain: General area (e.g., "automotive", "electronics", "fashion", "real-estate", "home-garden", "services")
   - productType: Specific product type (e.g., "tire", "smartphone", "dress", "apartment", "sofa")
   - keywords: Array of keywords describing the item (e.g., ["winter", "michelin", "205/55", "r16", "set"])
5. category: Suggested category (optional, can leave empty)
6. categoryProbabilities: If you can determine - top 3 possible categories
7. price: Market price in RSD as string
8. priceRange: {min, max} price range
9. attributes: 
   - For cars: brand, car_model, year, color (from: black, white, silver, gold, blue, red, green, yellow, purple, other), fuel_type (petrol, diesel, electric, hybrid, lpg, cng), transmission (manual, automatic, semi-automatic, cvt), mileage, engine_size
   - For tires: tire_width (e.g. "195"), tire_profile (e.g. "65"), tire_diameter (e.g. "17"), tire_season (summer, winter, all-season), tire_brand, tire_condition (new, used), tread_depth (e.g. "7 mm"), tire_year, tire_quantity (1, 2, 3, 4, set)
10. tags: 5-8 search tags in English
11. suggestedPhotos: What photos to add for better sales
12. translations: MANDATORY create COMPLETE translations of title & description to ${targetLanguages.map((l) => languageNames[userLanguage][l]).join(' and ')} (${targetLanguages.join(', ')}). 
    Format: {"${targetLanguages[0] || 'ru'}": {"title": "complete translated title", "description": "complete translated description with ALL emojis, details and formatting"}, "${targetLanguages[1] || 'sr'}": {"title": "complete translated title", "description": "complete translated description with ALL emojis, details and formatting"}}
    FORBIDDEN to use placeholders like [...] or dots! Translate EVERY line of description!
13. socialPosts: Short selling posts for whatsapp, telegram, instagram
14. location: {city: "city", region: "region", suggestedLocation: "area"}. Serbia cities: Belgrade, Novi Sad, Nis, Kragujevac, Subotica
15. condition: "new", "used" or "refurbished"
16. insights: {ru: {demand: "demand analysis", audience: "who buys", recommendations: "how to sell faster"}, en: {...}, sr: {...}}
17. originalLanguage: "${userLanguage}"

IMPORTANT: 
1. Create a SELLING LISTING, NOT a photo description!
2. Reply ONLY in JSON format! No additional text!
3. Don't use markdown blocks - only clean JSON!
4. Start your response immediately with { and end with }
5. DO NOT ADD any text before or after the JSON!

EXAMPLE OF CORRECT TRANSLATION FORMAT:
"translations": {
  "ru": {
    "title": "Volkswagen Touran 2.0 TDI • 7 мест • Отличное состояние",
    "description": "🚗 НАДЕЖНЫЙ СЕМЕЙНЫЙ АВТОМОБИЛЬ В ОТЛИЧНОМ СОСТОЯНИИ!\\n\\n✨ ОСНОВНЫЕ ПРЕИМУЩЕСТВА:\\n- Просторный и комфортный семейный автомобиль\\n- Экономичный двигатель 2.0 TDI\\n- 7 мест с системой Isofix..."
  },
  "sr": {
    "title": "Volkswagen Touran 2.0 TDI • 7 sedišta • Odlično stanje",
    "description": "🚗 POUZDANO PORODIČNO VOZILO U ODLIČNOM STANJU!\\n\\n✨ GLAVNE PREDNOSTI:\\n- Prostran i komforan porodični automobil..."
  }
}`,

    sr: `Ti si ekspert za kreiranje prodajnih oglasa za online tržište. Analiziraj sliku i napravi PRODAJNI OGLAS (NE opis fotografije!). Odgovor u JSON formatu:

1. title: Prodajni naslov proizvoda na srpskom (npr. "Volkswagen Touran 2015 • Odlično stanje")
2. titleVariants: Niz od 3 prodajna naslova za A/B testiranje
3. description: PRODAJNI opis za kupce na srpskom. NE opisuj šta se vidi na slici! Uključi:
   - Glavne prednosti i karakteristike
   - Oprema i specifikacije
   - Tehničko stanje
   - Šta kupac dobija
   - Poziv na akciju
   Koristi emoji za privlačnost

4. categoryHints: Semantičke informacije za određivanje kategorije:
   - domain: Opšta oblast (npr. "automotive", "electronics", "fashion", "real-estate", "home-garden", "services")
   - productType: Konkretan tip proizvoda (npr. "tire", "smartphone", "dress", "apartment", "sofa")
   - keywords: Niz ključnih reči koje opisuju proizvod (npr. ["zimska", "michelin", "205/55", "r16", "komplet"])
5. category: Predložena kategorija (opciono, možeš ostaviti prazno)
6. categoryProbabilities: Ako možeš odrediti - top 3 moguće kategorije
7. price: Tržišna cena u RSD kao string
8. priceRange: {min, max} raspon cena
9. attributes: 
   - Za automobile: brand, car_model, year, color (iz: black, white, silver, gold, blue, red, green, yellow, purple, other), fuel_type (petrol, diesel, electric, hybrid, lpg, cng), transmission (manual, automatic, semi-automatic, cvt), mileage, engine_size
   - Za gume: tire_width (npr. "195"), tire_profile (npr. "65"), tire_diameter (npr. "17"), tire_season (summer, winter, all-season), tire_brand, tire_condition (new, used), tread_depth (npr. "7 mm"), tire_year, tire_quantity (1, 2, 3, 4, set)
10. tags: 5-8 tagova za pretragu na srpskom
11. suggestedPhotos: Koje fotografije dodati za bolju prodaju
12. translations: OBAVEZNO napravi KOMPLETNE prevode title i description na ${targetLanguages.map((l) => languageNames[userLanguage][l]).join(' i ')} (${targetLanguages.join(', ')}). 
    Format: {"${targetLanguages[0] || 'ru'}": {"title": "kompletan prevedeni naslov", "description": "kompletan prevedeni opis sa SVIM emoji, detaljima i formatiranjem"}, "${targetLanguages[1] || 'en'}": {"title": "kompletan prevedeni naslov", "description": "kompletan prevedeni opis sa SVIM emoji, detaljima i formatiranjem"}}
    ZABRANJENO koristiti placeholder-e kao [...] ili tri tačke! Prevedi SVAKI red opisa!
13. socialPosts: Kratke prodajne objave za whatsapp, telegram, instagram
14. location: {city: "grad", region: "region", suggestedLocation: "kvart"}. Gradovi Srbije: Beograd, Novi Sad, Niš, Kragujevac, Subotica
15. condition: "new", "used" ili "refurbished"
16. insights: {ru: {demand: "analiza potražnje", audience: "ko kupuje", recommendations: "kako prodati brže"}, en: {...}, sr: {...}}
17. originalLanguage: "${userLanguage}"

VAŽNO: 
1. Napravi PRODAJNI OGLAS, a NE opis fotografije!
2. Odgovori SAMO u JSON formatu! Nema dodatnog teksta!
3. Ne koristi markdown blokove - samo čist JSON!
4. Počni odgovor odmah sa { i završi sa }
5. NE DODAJ nikakav tekst pre ili posle JSON-a!

PRIMER ISPRAVNOG FORMATA PREVODA:
"translations": {
  "ru": {
    "title": "Volkswagen Touran 2.0 TDI • 7 мест • Отличное состояние",
    "description": "🚗 НАДЕЖНЫЙ СЕМЕЙНЫЙ АВТОМОБИЛЬ В ОТЛИЧНОМ СОСТОЯНИИ!\\n\\n✨ ОСНОВНЫЕ ПРЕИМУЩЕСТВА:\\n- Просторный и комфортный семейный автомобиль\\n- Экономичный двигатель 2.0 TDI..."
  },
  "en": {
    "title": "Volkswagen Touran 2.0 TDI • 7 seats • Excellent condition",
    "description": "🚗 RELIABLE FAMILY VEHICLE IN EXCELLENT CONDITION!\\n\\n✨ MAIN ADVANTAGES:\\n- Spacious and comfortable family car..."
  }
}`,
  };

  return prompts[userLanguage] || prompts.ru;
}
