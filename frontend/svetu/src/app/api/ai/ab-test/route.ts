import { NextRequest, NextResponse } from 'next/server';

const CLAUDE_API_URL = 'https://api.anthropic.com/v1/messages';
const CLAUDE_API_KEY = process.env.NEXT_PUBLIC_CLAUDE_API_KEY || '';

export async function POST(request: NextRequest) {
  let titleVariants: string[] = [];

  try {
    const body = await request.json();
    titleVariants = body.titleVariants;

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
            content: `Проанализируй эти заголовки товаров на эффективность и отранжируй их:
${titleVariants.map((t: string, i: number) => `${i + 1}. ${t}`).join('\n')}

Учитывай:
- Потенциал кликабельности
- SEO ключевые слова
- Ясность и привлекательность
- Предпочтения местного рынка (сербско/русскоговорящие)

Ответь в формате JSON: {bestVariant: "заголовок", scores: [{title: "...", score: 0-100}]}`,
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
    const content = data.content[0]?.text;
    const result = JSON.parse(content);

    return NextResponse.json(result);
  } catch (error) {
    console.error('API route error:', error);
    // Fallback response
    return NextResponse.json({
      bestVariant: titleVariants[0] || '',
      scores: titleVariants.map((title: string, i: number) => ({
        title,
        score: 90 - i * 10,
      })),
    });
  }
}
