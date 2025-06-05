import { NextRequest, NextResponse } from 'next/server';

const BACKEND_URL = process.env.BACKEND_URL || 'http://localhost:3000';

export async function POST(request: NextRequest) {
  try {
    const body = await request.json();

    // Отправляем запрос на backend
    const response = await fetch(`${BACKEND_URL}/api/auth/register`, {
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
    const setCookieHeader = response.headers.get('set-cookie');
    if (setCookieHeader) {
      const cookies = setCookieHeader.split(', ');
      cookies.forEach((cookie) => {
        if (cookie.includes('refresh_token=')) {
          // Убираем SameSite=None для same-origin
          const modifiedCookie = cookie.replace(
            'SameSite=None',
            'SameSite=Lax'
          );
          res.headers.append('Set-Cookie', modifiedCookie);
        }
      });
    }

    return res;
  } catch (error) {
    console.error('Register API error:', error);
    return NextResponse.json(
      { error: 'Internal server error' },
      { status: 500 }
    );
  }
}
