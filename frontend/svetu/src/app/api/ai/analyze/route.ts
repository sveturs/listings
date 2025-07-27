import { NextRequest, NextResponse } from 'next/server';

const CLAUDE_API_URL = 'https://api.anthropic.com/v1/messages';
const CLAUDE_API_KEY = process.env.NEXT_PUBLIC_CLAUDE_API_KEY || '';

export async function POST(request: NextRequest) {
  try {
    const { imageBase64, userLanguage = 'ru' } = await request.json();

    if (!imageBase64) {
      return NextResponse.json(
        { error: 'Image data is required' },
        { status: 400 }
      );
    }

    if (!CLAUDE_API_KEY) {
      return NextResponse.json(
        { error: 'Claude API key is not configured' },
        { status: 500 }
      );
    }

    console.log('Proxying request to Claude API...');
    console.log('User language:', userLanguage);
    console.log('Image data length:', imageBase64.length);

    const response = await fetch(CLAUDE_API_URL, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'X-API-Key': CLAUDE_API_KEY,
        'anthropic-version': '2023-06-01',
      },
      body: JSON.stringify({
        model: 'claude-3-5-sonnet-20241022',
        max_tokens: 4096,
        messages: [
          {
            role: 'user',
            content: [
              {
                type: 'image',
                source: {
                  type: 'base64',
                  media_type: 'image/jpeg',
                  data: imageBase64,
                },
              },
              {
                type: 'text',
                text: getAnalysisPrompt(userLanguage),
              },
            ],
          },
        ],
      }),
    });

    console.log('Claude API response status:', response.status);
    console.log(
      'Claude API response headers:',
      Object.fromEntries(response.headers)
    );

    if (!response.ok) {
      const errorData = await response.text();
      console.error('Claude API error:', response.status, errorData);
      return NextResponse.json(
        { error: `Claude API error: ${response.status}` },
        { status: response.status }
      );
    }

    const data = await response.json();
    console.log('Claude API response structure keys:', Object.keys(data));
    console.log(
      'Claude API response content array length:',
      data.content?.length
    );
    const content = data.content[0]?.text;

    if (!content) {
      return NextResponse.json(
        { error: 'No content in Claude response' },
        { status: 500 }
      );
    }

    // Parse and return the JSON response
    // Очищаем content от возможных проблемных символов, но сохраняем переносы строк (\n = \u000A)
    const cleanContent = content.replace(
      /[\u0000-\u0009\u000B-\u001F\u007F-\u009F]/g,
      ''
    );

    try {
      console.log('Parsing Claude response, length:', cleanContent.length);

      const analysis = JSON.parse(cleanContent);
      console.log('Successfully parsed Claude response');
      return NextResponse.json(analysis);
    } catch (parseError) {
      console.error('JSON parsing error:', parseError);
      console.error('Full content that failed to parse:', content);
      console.error('Clean content that failed to parse:', cleanContent);
      console.error('Content length:', content.length);
      console.error('Clean content length:', cleanContent.length);

      // Попробуем извлечь JSON из ответа если он обёрнут в ```json
      const jsonMatch = content.match(/```json\n([\s\S]*?)\n```/);
      if (jsonMatch) {
        try {
          const analysis = JSON.parse(jsonMatch[1]);
          console.log('Successfully parsed JSON from code block');
          return NextResponse.json(analysis);
        } catch (codeBlockError) {
          console.error(
            'Failed to parse JSON from code block:',
            codeBlockError
          );
        }
      }

      return NextResponse.json(
        { error: 'Failed to parse AI response as JSON' },
        { status: 500 }
      );
    }
  } catch (error) {
    console.error('API route error:', error);
    return NextResponse.json(
      { error: 'Failed to analyze image' },
      { status: 500 }
    );
  }
}

function getAnalysisPrompt(userLanguage: string): string {
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

4. category: Выбери из: electronics, fashion, automotive, real-estate, home-garden, agriculture, industrial, food-beverages, services, sports-recreation
5. categoryProbabilities: Топ-3 категории с вероятностями
6. price: Рыночная цена в РСД как строка
7. priceRange: {min, max} диапазон цен
8. attributes: Для авто ТОЛЬКО: brand, car_model, year, color (из: black, white, silver, gold, blue, red, green, yellow, purple, other), fuel_type (petrol, diesel, electric, hybrid, lpg, cng), transmission (manual, automatic, semi-automatic, cvt), mileage, engine_size
9. tags: 5-8 поисковых тегов на русском
10. suggestedPhotos: Какие фото добавить для лучшей продажи
11. translations: ОБЯЗАТЕЛЬНО создай ПОЛНЫЕ переводы title и description на ${targetLanguages.map((l) => languageNames[userLanguage][l]).join(' и ')} (${targetLanguages.join(', ')}). 
    Формат: {"${targetLanguages[0] || 'en'}": {"title": "полный переведенный заголовок", "description": "полное переведенное описание со ВСЕМИ эмодзи, деталями и форматированием"}, "${targetLanguages[1] || 'sr'}": {"title": "полный переведенный заголовок", "description": "полное переведенное описание со ВСЕМИ эмодзи, деталями и форматированием"}}
    ЗАПРЕЩЕНО использовать заглушки типа [...] или многоточия! Переводи КАЖДУЮ строку описания!
12. socialPosts: Короткие продающие посты для whatsapp, telegram, instagram
13. location: {city: "город", region: "регион", suggestedLocation: "район"}. Города Сербии: Белград, Нови-Сад, Ниш, Крагуевац, Суботица
14. condition: "new", "used" или "refurbished"
15. insights: {ru: {demand: "анализ спроса", audience: "кто покупает", recommendations: "как продать быстрее"}, en: {...}, sr: {...}}
16. originalLanguage: "${userLanguage}"

ВАЖНО: 
1. Создавай ПРОДАЮЩЕЕ ОБЪЯВЛЕНИЕ, а НЕ описание фотографии!
2. Отвечай ТОЛЬКО в формате JSON! Никакого дополнительного текста!
3. Не используй markdown блоки - только чистый JSON!

ПРИМЕР ПРАВИЛЬНОГО ФОРМАТА ПЕРЕВОДОВ:
"translations": {
  "en": {
    "title": "Volkswagen Touran 2.0 TDI • 7 seats • Excellent condition",
    "description": "🚗 RELIABLE FAMILY VEHICLE IN EXCELLENT CONDITION!\n\n✨ MAIN ADVANTAGES:\n- Spacious and comfortable family car\n- Economical 2.0 TDI engine\n- 7 seats with Isofix system\n- Large trunk space\n\n🔧 EQUIPMENT:\n- Automatic climate control..."
  },
  "sr": {
    "title": "Volkswagen Touran 2.0 TDI • 7 sedišta • Odlično stanje", 
    "description": "🚗 POUZDANO PORODIČNO VOZILO U ODLIČNOM STANJU!\n\n✨ GLAVNE PREDNOSTI:\n- Prostran i komforan porodični automobil..."
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

4. category: Choose from: electronics, fashion, automotive, real-estate, home-garden, agriculture, industrial, food-beverages, services, sports-recreation
5. categoryProbabilities: Top 3 categories with probabilities
6. price: Market price in RSD as string
7. priceRange: {min, max} price range
8. attributes: For cars ONLY: brand, car_model, year, color (from: black, white, silver, gold, blue, red, green, yellow, purple, other), fuel_type (petrol, diesel, electric, hybrid, lpg, cng), transmission (manual, automatic, semi-automatic, cvt), mileage, engine_size
9. tags: 5-8 search tags in English
10. suggestedPhotos: What photos to add for better sales
11. translations: MANDATORY create COMPLETE translations of title & description to ${targetLanguages.map((l) => languageNames[userLanguage][l]).join(' and ')} (${targetLanguages.join(', ')}). 
    Format: {"${targetLanguages[0] || 'ru'}": {"title": "complete translated title", "description": "complete translated description with ALL emojis, details and formatting"}, "${targetLanguages[1] || 'sr'}": {"title": "complete translated title", "description": "complete translated description with ALL emojis, details and formatting"}}
    FORBIDDEN to use placeholders like [...] or dots! Translate EVERY line of description!
12. socialPosts: Short selling posts for whatsapp, telegram, instagram
13. location: {city: "city", region: "region", suggestedLocation: "area"}. Serbia cities: Belgrade, Novi Sad, Nis, Kragujevac, Subotica
14. condition: "new", "used" or "refurbished"
15. insights: {ru: {demand: "demand analysis", audience: "who buys", recommendations: "how to sell faster"}, en: {...}, sr: {...}}
16. originalLanguage: "${userLanguage}"

IMPORTANT: 
1. Create a SELLING LISTING, NOT a photo description!
2. Reply ONLY in JSON format! No additional text!
3. Don't use markdown blocks - only clean JSON!

EXAMPLE OF CORRECT TRANSLATION FORMAT:
"translations": {
  "ru": {
    "title": "Volkswagen Touran 2.0 TDI • 7 мест • Отличное состояние",
    "description": "🚗 НАДЕЖНЫЙ СЕМЕЙНЫЙ АВТОМОБИЛЬ В ОТЛИЧНОМ СОСТОЯНИИ!\n\n✨ ОСНОВНЫЕ ПРЕИМУЩЕСТВА:\n- Просторный и комфортный семейный автомобиль\n- Экономичный двигатель 2.0 TDI\n- 7 мест с системой Isofix..."
  },
  "sr": {
    "title": "Volkswagen Touran 2.0 TDI • 7 sedišta • Odlično stanje",
    "description": "🚗 POUZDANO PORODIČNO VOZILO U ODLIČNOM STANJU!\n\n✨ GLAVNE PREDNOSTI:\n- Prostran i komforan porodični automobil..."
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

4. category: Izaberi iz: electronics, fashion, automotive, real-estate, home-garden, agriculture, industrial, food-beverages, services, sports-recreation
5. categoryProbabilities: Top 3 kategorije sa verovatnoćama
6. price: Tržišna cena u RSD kao string
7. priceRange: {min, max} raspon cena
8. attributes: Za automobile SAMO: brand, car_model, year, color (iz: black, white, silver, gold, blue, red, green, yellow, purple, other), fuel_type (petrol, diesel, electric, hybrid, lpg, cng), transmission (manual, automatic, semi-automatic, cvt), mileage, engine_size
9. tags: 5-8 tagova za pretragu na srpskom
10. suggestedPhotos: Koje fotografije dodati za bolju prodaju
11. translations: OBAVEZNO napravi KOMPLETNE prevode title i description na ${targetLanguages.map((l) => languageNames[userLanguage][l]).join(' i ')} (${targetLanguages.join(', ')}). 
    Format: {"${targetLanguages[0] || 'ru'}": {"title": "kompletan prevedeni naslov", "description": "kompletan prevedeni opis sa SVIM emoji, detaljima i formatiranjem"}, "${targetLanguages[1] || 'en'}": {"title": "kompletan prevedeni naslov", "description": "kompletan prevedeni opis sa SVIM emoji, detaljima i formatiranjem"}}
    ZABRANJENO koristiti placeholder-e kao [...] ili tri tačke! Prevedi SVAKI red opisa!
12. socialPosts: Kratke prodajne objave za whatsapp, telegram, instagram
13. location: {city: "grad", region: "region", suggestedLocation: "kvart"}. Gradovi Srbije: Beograd, Novi Sad, Niš, Kragujevac, Subotica
14. condition: "new", "used" ili "refurbished"
15. insights: {ru: {demand: "analiza potražnje", audience: "ko kupuje", recommendations: "kako prodati brže"}, en: {...}, sr: {...}}
16. originalLanguage: "${userLanguage}"

VAŽNO: 
1. Napravi PRODAJNI OGLAS, a NE opis fotografije!
2. Odgovori SAMO u JSON formatu! Nema dodatnog teksta!
3. Ne koristi markdown blokove - samo čist JSON!

PRIMER ISPRAVNOG FORMATA PREVODA:
"translations": {
  "ru": {
    "title": "Volkswagen Touran 2.0 TDI • 7 мест • Отличное состояние",
    "description": "🚗 НАДЕЖНЫЙ СЕМЕЙНЫЙ АВТОМОБИЛЬ В ОТЛИЧНОМ СОСТОЯНИИ!\n\n✨ ОСНОВНЫЕ ПРЕИМУЩЕСТВА:\n- Просторный и комфортный семейный автомобиль\n- Экономичный двигатель 2.0 TDI..."
  },
  "en": {
    "title": "Volkswagen Touran 2.0 TDI • 7 seats • Excellent condition",
    "description": "🚗 RELIABLE FAMILY VEHICLE IN EXCELLENT CONDITION!\n\n✨ MAIN ADVANTAGES:\n- Spacious and comfortable family car..."
  }
}`,
  };

  return prompts[userLanguage] || prompts.ru;
}
