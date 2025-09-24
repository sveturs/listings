import { NextRequest, NextResponse } from 'next/server';
import { cookies } from 'next/headers';

const BACKEND_URL = process.env.BACKEND_INTERNAL_URL || 'http://localhost:3000';

export async function POST(_request: NextRequest) {
  try {
    const cookieStore = await cookies();
    const token = cookieStore.get('access_token')?.value;

    if (token) {
      // Отправляем запрос на logout к backend
      await fetch(`${BACKEND_URL}/api/v1/auth/logout`, {
        method: 'POST',
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });
    }

    // Создаем response
    const res = NextResponse.json({ success: true });

    // Удаляем cookies
    res.cookies.delete('access_token');
    res.cookies.delete('refresh_token');

    return res;
  } catch (error) {
    console.error('Logout error:', error);
    // Даже при ошибке удаляем токены
    const res = NextResponse.json({ success: true });
    res.cookies.delete('access_token');
    res.cookies.delete('refresh_token');
    return res;
  }
}
