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

    // Получаем refresh_token cookie для проверки наличия
    // const refreshToken = request.cookies.get('refresh_token')?.value;

    // Отправляем запрос на backend с всеми cookies
    const cookieHeader = request.headers.get('cookie') || '';

    const response = await fetch(`${BACKEND_URL}/api/auth/logout`, {
      method: 'POST',
      headers: {
        Cookie: cookieHeader, // Передаем все cookies, включая refresh_token
      },
    });

    // Создаем ответ
    const res = NextResponse.json(
      { success: true },
      { status: response.status }
    );

    // Принудительно удаляем ВСЕ auth cookies
    res.cookies.set('refresh_token', '', {
      httpOnly: true,
      secure: process.env.NODE_ENV === 'production',
      sameSite: 'lax',
      path: '/',
      maxAge: 0, // Удаляем cookie
    });

    // Удаляем session_token (старая сессия)
    res.cookies.set('session_token', '', {
      httpOnly: true,
      secure: process.env.NODE_ENV === 'production',
      sameSite: 'lax',
      path: '/',
      maxAge: 0,
    });

    // Удаляем jwt_token (старый JWT в cookie)
    res.cookies.set('jwt_token', '', {
      httpOnly: true,
      secure: process.env.NODE_ENV === 'production',
      sameSite: 'lax',
      path: '/',
      maxAge: 0,
    });

    // Добавляем дополнительные заголовки для удаления cookies
    res.headers.append(
      'Set-Cookie',
      'refresh_token=; Path=/; Max-Age=0; HttpOnly; SameSite=Lax'
    );
    res.headers.append(
      'Set-Cookie',
      'session_token=; Path=/; Max-Age=0; HttpOnly; SameSite=Lax'
    );
    res.headers.append(
      'Set-Cookie',
      'jwt_token=; Path=/; Max-Age=0; HttpOnly; SameSite=Lax'
    );

    return res;
  } catch (error) {
    console.error('Logout API error:', error);
    return NextResponse.json(
      { error: 'Internal server error' },
      { status: 500 }
    );
  }
}
