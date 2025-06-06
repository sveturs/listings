import { NextRequest, NextResponse } from 'next/server';

const BACKEND_URL = process.env.BACKEND_URL || 'http://localhost:3000';

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

    // Fallback на старый метод с session token
    const response = await fetch(`${BACKEND_URL}/auth/session`, {
      method: 'GET',
      headers,
    });

    const data = await response.json();
    console.log('[API Route] Session response from backend:', data);

    if (!response.ok) {
      return NextResponse.json(data, { status: response.status });
    }

    console.log('[API Route] Returning session to frontend:', data);
    return NextResponse.json(data);
  } catch (error) {
    console.error('Session API error:', error);
    return NextResponse.json(
      { error: 'Internal server error' },
      { status: 500 }
    );
  }
}
