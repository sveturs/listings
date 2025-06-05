import { NextRequest, NextResponse } from 'next/server';

const BACKEND_URL = process.env.BACKEND_URL || 'http://localhost:3000';

// Валидация backend URL для предотвращения SSRF
// Можно расширить список через переменную окружения ALLOWED_BACKEND_HOSTS
const DEFAULT_ALLOWED_HOSTS = [
  'localhost:3000', // Development
  '127.0.0.1:3000', // Alternative localhost
  'backend:3000', // Docker internal network
  'svetu.rs', // Production (через nginx)
  'www.svetu.rs', // Production with www
];
const ALLOWED_HOSTS = process.env.ALLOWED_BACKEND_HOSTS
  ? [...DEFAULT_ALLOWED_HOSTS, ...process.env.ALLOWED_BACKEND_HOSTS.split(',')]
  : DEFAULT_ALLOWED_HOSTS;

function validateBackendUrl(url: string): boolean {
  try {
    const parsedUrl = new URL(url);
    return ALLOWED_HOSTS.includes(parsedUrl.host);
  } catch {
    return false;
  }
}

export async function POST(request: NextRequest) {
  try {
    // Проверяем валидность backend URL
    if (!validateBackendUrl(BACKEND_URL)) {
      console.error('Invalid backend URL:', BACKEND_URL);
      return NextResponse.json(
        { error: 'Configuration error' },
        { status: 500 }
      );
    }

    // Получаем refresh_token cookie из запроса
    const refreshToken = request.cookies.get('refresh_token')?.value;

    // Отправляем запрос на backend с cookie
    const response = await fetch(`${BACKEND_URL}/api/auth/refresh`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        // Передаем cookie в backend
        Cookie: refreshToken ? `refresh_token=${refreshToken}` : '',
      },
    });

    const data = await response.json();

    if (!response.ok) {
      return NextResponse.json(data, { status: response.status });
    }

    // Создаем ответ
    const res = NextResponse.json(data);

    // Копируем новый refresh_token cookie из backend ответа
    const setCookieHeaders = response.headers.getSetCookie();
    setCookieHeaders.forEach((cookieHeader: string) => {
      if (cookieHeader.includes('refresh_token=')) {
        // Убираем Domain и меняем SameSite для same-origin
        const modifiedCookie = cookieHeader
          .replace(/Domain=[^;]+;?\s*/i, '') // Убираем Domain
          .replace('SameSite=None', 'SameSite=Lax');
        res.headers.append('Set-Cookie', modifiedCookie);
      }
    });

    return res;
  } catch (error) {
    console.error('Refresh API error:', error);
    return NextResponse.json(
      { error: 'Internal server error' },
      { status: 500 }
    );
  }
}
