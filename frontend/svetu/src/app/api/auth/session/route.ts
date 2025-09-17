import { NextRequest, NextResponse } from 'next/server';
import configManager from '@/config';

const BACKEND_URL =
  process.env.BACKEND_URL || configManager.getApiUrl({ internal: true });

// Auth service URL - для auth endpoints
const AUTH_SERVICE_URL =
  process.env.NEXT_PUBLIC_AUTH_SERVICE_URL || configManager.getAuthServiceUrl();

export async function GET(request: NextRequest) {
  try {
    // Получаем refresh_token cookie из запроса
    const refreshToken = request.cookies.get('refresh_token')?.value;
    const sessionToken = request.cookies.get('session_token')?.value;

    // Получаем Authorization header для JWT токена
    const authHeader = request.headers.get('authorization');

    // Формируем headers для запроса к backend
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
    };

    // Передаем Authorization header если есть
    if (authHeader) {
      headers['Authorization'] = authHeader;
    }

    // Формируем cookie string для backend
    const cookies: string[] = [];
    if (refreshToken) cookies.push(`refresh_token=${refreshToken}`);
    if (sessionToken) cookies.push(`session_token=${sessionToken}`);
    if (cookies.length > 0) {
      headers['Cookie'] = cookies.join('; ');
    }

    // Для JWT-based auth сначала проверяем JWT токен
    const jwtToken = request.cookies.get('jwt_token')?.value;

    if (jwtToken) {
      // Если есть JWT токен, получаем информацию о пользователе
      const userResponse = await fetch(`${BACKEND_URL}/api/v1/users/me`, {
        method: 'GET',
        headers: {
          ...headers,
          Authorization: `Bearer ${jwtToken}`,
        },
      });

      if (userResponse.ok) {
        const userData = await userResponse.json();
        console.log('[API Route] User data from JWT:', userData);

        // Возвращаем данные в формате сессии
        return NextResponse.json({
          authenticated: true,
          user: userData.data || userData,
        });
      } else if (userResponse.status === 401) {
        // Токен невалидный или истек
        return NextResponse.json({
          authenticated: false,
        });
      }
    }

    // Если нет JWT токена, проверяем access token из localStorage через header
    if (authHeader && authHeader.startsWith('Bearer ')) {
      const accessToken = authHeader.substring(7);

      // Пробуем получить данные пользователя через auth service
      const response = await fetch(`${AUTH_SERVICE_URL}/api/v1/auth/me`, {
        method: 'GET',
        headers: {
          Authorization: `Bearer ${accessToken}`,
          'Content-Type': 'application/json',
        },
      });

      if (response.ok) {
        const data = await response.json();
        console.log('[API Route] User data from auth service:', data);

        return NextResponse.json({
          authenticated: true,
          user: data.user || data,
        });
      }
    }

    // Если ничего не сработало, возвращаем неаутентифицированный статус
    return NextResponse.json({
      authenticated: false,
    });
  } catch (error) {
    console.error('Session API error:', error);
    return NextResponse.json(
      { error: 'Internal server error' },
      { status: 500 }
    );
  }
}
