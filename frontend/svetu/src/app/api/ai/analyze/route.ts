import { NextRequest, NextResponse } from 'next/server';

const CLAUDE_API_URL = 'https://api.anthropic.com/v1/messages';
const CLAUDE_API_KEY = process.env.NEXT_PUBLIC_CLAUDE_API_KEY || '';

export async function POST(request: NextRequest) {
  try {
    const { imageBase64 } = await request.json();

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
                text: `Проанализируй это изображение товара и предоставь детальный анализ в формате JSON. Ответ должен включать:

1. title: Основное название товара на русском языке (краткое и описательное)
2. titleVariants: Массив из 3 альтернативных заголовков для A/B тестирования
3. description: Детальное описание товара на русском с эмодзи и форматированием
4. category: Основная категория (electronics, fashion, home, auto, или other)
5. categoryProbabilities: Массив объектов с {name, probability} для топ-3 категорий
6. price: Предлагаемая цена в РСД (сербских динарах) как строка
7. priceRange: Объект с {min, max} диапазоном цен
8. attributes: Объект с парами ключ-значение для атрибутов товара (бренд, модель, состояние и т.д.)
9. tags: Массив релевантных поисковых тегов на русском
10. suggestedPhotos: Массив предложений дополнительных ракурсов/типов фото
11. translations: Объект с переводами на en, sr (сербская латиница)
12. socialPosts: Объект с оптимизированными постами для whatsapp, telegram, instagram

Отвечай ТОЛЬКО валидным JSON, без дополнительного текста.`,
              },
            ],
          },
        ],
      }),
    });

    if (!response.ok) {
      const errorData = await response.text();
      console.error('Claude API error:', response.status, errorData);
      return NextResponse.json(
        { error: `Claude API error: ${response.status}` },
        { status: response.status }
      );
    }

    const data = await response.json();
    const content = data.content[0]?.text;

    if (!content) {
      return NextResponse.json(
        { error: 'No content in Claude response' },
        { status: 500 }
      );
    }

    // Parse and return the JSON response
    const analysis = JSON.parse(content);
    return NextResponse.json(analysis);
  } catch (error) {
    console.error('API route error:', error);
    return NextResponse.json(
      { error: 'Failed to analyze image' },
      { status: 500 }
    );
  }
}