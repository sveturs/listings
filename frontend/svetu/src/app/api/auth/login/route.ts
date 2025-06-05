import { NextRequest, NextResponse } from 'next/server';

const BACKEND_URL = process.env.BACKEND_URL || 'http://localhost:3000';

export async function POST(request: NextRequest) {
  try {
    const body = await request.json();

    // Отправляем запрос на backend
    const response = await fetch(`${BACKEND_URL}/api/auth/login`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(body),
    });

    const data = await response.json();

    if (!response.ok) {
      return NextResponse.json(data, { status: response.status });
    }

    // Создаем ответ
    const res = NextResponse.json(data);

    // Копируем refresh_token cookie из backend ответа
    // Headers.getSetCookie() возвращает массив всех Set-Cookie заголовков
    const setCookieHeaders = response.headers.getSetCookie();
    console.log(
      '[Login API Route] Set-Cookie headers from backend:',
      setCookieHeaders
    );

    // Обрабатываем каждый Set-Cookie header отдельно
    setCookieHeaders.forEach((cookieHeader: string) => {
      if (cookieHeader.includes('refresh_token=')) {
        // Убираем Domain и меняем SameSite для same-origin
        const modifiedCookie = cookieHeader
          .replace(/Domain=[^;]+;?\s*/i, '') // Убираем Domain
          .replace('SameSite=None', 'SameSite=Lax');

        console.log(
          '[Login API Route] Setting refresh_token cookie:',
          modifiedCookie
        );
        res.headers.append('Set-Cookie', modifiedCookie);
      }
    });

    return res;
  } catch (error) {
    console.error('Login API error:', error);
    return NextResponse.json(
      { error: 'Internal server error' },
      { status: 500 }
    );
  }
}
