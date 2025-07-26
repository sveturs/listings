import { NextRequest, NextResponse } from 'next/server';

const CLAUDE_API_URL = 'https://api.anthropic.com/v1/messages';
const CLAUDE_API_KEY = process.env.NEXT_PUBLIC_CLAUDE_API_KEY || '';

export async function POST(request: NextRequest) {
  try {
    const { title, category, attributes } = await request.json();

    if (!CLAUDE_API_KEY) {
      return NextResponse.json(
        { error: 'Claude API key is not configured' },
        { status: 500 }
      );
    }

    const response = await fetch(CLAUDE_API_URL, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'X-API-Key': CLAUDE_API_KEY,
        'anthropic-version': '2023-06-01',
      },
      body: JSON.stringify({
        model: 'claude-3-5-sonnet-20241022',
        max_tokens: 1024,
        messages: [
          {
            role: 'user',
            content: `Создай оптимизированное описание товара на русском языке для:
Название: ${title}
Категория: ${category}
Атрибуты: ${JSON.stringify(attributes)}

Требования:
- Используй эмодзи для лучшей читаемости
- Включи ключевые преимущества
- Укажи состояние и важные детали
- Будь кратким, но информативным
- Оптимизируй для поисковых ключевых слов
- НЕ включай контактную информацию

Ответь только текстом описания.`,
          },
        ],
      }),
    });

    if (!response.ok) {
      const errorData = await response.text();
      return NextResponse.json(
        { error: `Claude API error: ${response.status}` },
        { status: response.status }
      );
    }

    const data = await response.json();
    const description = data.content[0]?.text || '';

    return NextResponse.json({ description });
  } catch (error) {
    console.error('API route error:', error);
    return NextResponse.json(
      { error: 'Failed to generate description' },
      { status: 500 }
    );
  }
}