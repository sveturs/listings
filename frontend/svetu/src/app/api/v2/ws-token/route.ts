import { NextResponse } from 'next/server';
import { cookies } from 'next/headers';

/**
 * Endpoint для получения access token для WebSocket соединения
 *
 * WebSocket не может использовать httpOnly cookies напрямую,
 * поэтому frontend должен получить токен через этот endpoint
 * и передать его в query параметре при подключении к WebSocket
 */
export async function GET() {
  try {
    const cookieStore = await cookies();
    const accessToken = cookieStore.get('access_token')?.value;

    if (!accessToken) {
      return NextResponse.json(
        { error: 'No access token found' },
        { status: 401 }
      );
    }

    return NextResponse.json({
      token: accessToken,
    });
  } catch (error) {
    console.error('[WS Token] Failed to get token:', error);
    return NextResponse.json(
      { error: 'Failed to retrieve token' },
      { status: 500 }
    );
  }
}
