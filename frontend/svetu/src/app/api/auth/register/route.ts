import { NextRequest, NextResponse } from 'next/server';

const BACKEND_URL = process.env.BACKEND_INTERNAL_URL || 'http://localhost:3000';

export async function POST(request: NextRequest) {
  try {
    const body = await request.json();

    // Прокси запрос к backend
    // Auth service expects: email, password, name, terms_accepted
    const requestBody = {
      email: body.email,
      password: body.password,
      name: body.name || '',
      terms_accepted: true,
    };

    const response = await fetch(`${BACKEND_URL}/api/v1/auth/register`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(requestBody),
    });

    const data = await response.json();

    if (response.ok && data.access_token) {
      // Создаем response
      const res = NextResponse.json({
        success: true,
        user: data.user,
      });

      // Set cookie with token
      res.cookies.set({
        name: 'access_token',
        value: data.access_token,
        httpOnly: true,
        secure: process.env.NODE_ENV === 'production',
        sameSite: 'lax',
        path: '/',
        maxAge: 60 * 60 * 48, // 48 hours (соответствует времени жизни токена)
      });

      if (data.refresh_token) {
        res.cookies.set({
          name: 'refresh_token',
          value: data.refresh_token,
          httpOnly: true,
          secure: process.env.NODE_ENV === 'production',
          sameSite: 'lax',
          path: '/',
          maxAge: 60 * 60 * 24 * 30, // 30 days
        });
      }

      return res;
    }

    return NextResponse.json(data, { status: response.status });
  } catch (error) {
    console.error('Registration error:', error);
    return NextResponse.json({ error: 'Failed to register' }, { status: 500 });
  }
}
