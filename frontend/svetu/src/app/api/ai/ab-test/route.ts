import { NextRequest, NextResponse } from 'next/server';

const CLAUDE_API_URL = 'https://api.anthropic.com/v1/messages';
const CLAUDE_API_KEY = process.env.NEXT_PUBLIC_CLAUDE_API_KEY || '';

// Кэш для хранения результатов A/B тестов
const abTestCache = new Map<string, { result: any; timestamp: number }>();
const CACHE_TTL = 60 * 60 * 1000; // 1 час

// Эвристический анализ заголовков (быстрая альтернатива)
function analyzeHeuristically(titleVariants: string[]): any {
  const scores = titleVariants.map((title) => {
    let score = 50; // базовый балл

    // Факторы повышающие привлекательность
    if (title.length >= 10 && title.length <= 60) score += 15; // оптимальная длина
    if (/\d+/.test(title)) score += 10; // содержит числа (модель, размер)
    if (/[A-Z]/.test(title)) score += 5; // есть заглавные буквы (бренды)
    if (/новый|новая|новое|отличн|идеальн|гарант/i.test(title)) score += 10; // позитивные слова
    if (/\b(GB|ГБ|TB|ТБ|MHz|МГц|MP|МП)\b/i.test(title)) score += 8; // технические характеристики

    // Факторы понижающие привлекательность
    if (title.length < 5 || title.length > 100) score -= 20; // слишком короткий/длинный
    if (/[!?]{2,}/.test(title)) score -= 10; // избыток восклицаний
    if (title === title.toUpperCase()) score -= 15; // весь КАПСОМ

    // Нормализация (0-100)
    score = Math.max(0, Math.min(100, score));

    return { title, score };
  });

  // Сортировка по баллам
  scores.sort((a, b) => b.score - a.score);

  return {
    bestVariant: scores[0]?.title || titleVariants[0],
    scores,
  };
}

export async function POST(request: NextRequest) {
  let titleVariants: string[] = [];

  try {
    const body = await request.json();
    titleVariants = body.titleVariants;

    if (!titleVariants || titleVariants.length === 0) {
      return NextResponse.json(
        { error: 'No title variants provided' },
        { status: 400 }
      );
    }

    // Создаем ключ для кэша
    const cacheKey = JSON.stringify(titleVariants);

    // Проверяем кэш
    const cached = abTestCache.get(cacheKey);
    if (cached && Date.now() - cached.timestamp < CACHE_TTL) {
      console.log('A/B test result from cache');
      return NextResponse.json(cached.result);
    }

    // Если нет API ключа или слишком много вариантов - используем эвристику
    if (!CLAUDE_API_KEY || titleVariants.length > 5) {
      console.log('Using heuristic analysis for A/B test');
      const result = analyzeHeuristically(titleVariants);

      // Кэшируем результат
      abTestCache.set(cacheKey, { result, timestamp: Date.now() });

      return NextResponse.json(result);
    }

    // Используем более быструю модель Claude Haiku для простых задач
    const response = await fetch(CLAUDE_API_URL, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'X-API-Key': CLAUDE_API_KEY,
        'anthropic-version': '2023-06-01',
      },
      body: JSON.stringify({
        model: 'claude-3-haiku-20240307', // Быстрая модель для простых задач
        max_tokens: 256, // Уменьшаем лимит токенов
        temperature: 0.3, // Более детерминированный ответ
        messages: [
          {
            role: 'user',
            content: `Быстро оцени заголовки (1-100):
${titleVariants.map((t: string, i: number) => `${i + 1}. ${t}`).join('\n')}

JSON: {bestVariant: "заголовок", scores: [{title: "...", score: число}]}`,
          },
        ],
      }),
    });

    if (!response.ok) {
      console.error('Claude API error, falling back to heuristics');
      const result = analyzeHeuristically(titleVariants);
      return NextResponse.json(result);
    }

    const data = await response.json();
    const content = data.content[0]?.text;

    try {
      // Извлекаем JSON из ответа
      const jsonMatch = content.match(/\{[\s\S]*\}/);
      if (!jsonMatch) throw new Error('No JSON in response');

      const result = JSON.parse(jsonMatch[0]);

      // Кэшируем результат
      abTestCache.set(cacheKey, { result, timestamp: Date.now() });

      return NextResponse.json(result);
    } catch (parseError) {
      console.error('Failed to parse Claude response, using heuristics');
      const result = analyzeHeuristically(titleVariants);
      return NextResponse.json(result);
    }
  } catch (error) {
    console.error('API route error:', error);
    // Fallback response с эвристическим анализом
    const result = analyzeHeuristically(titleVariants);
    return NextResponse.json(result);
  }
}
