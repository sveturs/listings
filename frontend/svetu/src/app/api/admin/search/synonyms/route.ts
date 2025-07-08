import { NextRequest, NextResponse } from 'next/server';
import { apiClient } from '@/lib/api-client';

export async function GET(request: NextRequest) {
  try {
    // Получаем токен из заголовка Authorization
    const authHeader = request.headers.get('authorization');
    const accessToken = authHeader?.replace('Bearer ', '');

    if (!accessToken) {
      return NextResponse.json({ error: 'Unauthorized' }, { status: 401 });
    }

    const { searchParams } = new URL(request.url);
    const lang = searchParams.get('lang') || 'ru';

    const response = await apiClient.get(
      `/api/v1/admin/search/synonyms?lang=${lang}`,
      {
        headers: {
          Authorization: `Bearer ${accessToken}`,
        },
      }
    );
    return NextResponse.json(response.data);
  } catch (error) {
    console.error('Error fetching synonyms:', error);
    return NextResponse.json(
      { error: 'Failed to fetch synonyms' },
      { status: 500 }
    );
  }
}

export async function POST(request: NextRequest) {
  try {
    // Получаем токен из заголовка Authorization
    const authHeader = request.headers.get('authorization');
    const accessToken = authHeader?.replace('Bearer ', '');

    if (!accessToken) {
      return NextResponse.json({ error: 'Unauthorized' }, { status: 401 });
    }

    const data = await request.json();
    const response = await apiClient.post(
      '/api/v1/admin/search/synonyms',
      data,
      {
        headers: {
          Authorization: `Bearer ${accessToken}`,
        },
      }
    );
    return NextResponse.json(response.data);
  } catch (error) {
    console.error('Error creating synonym:', error);
    return NextResponse.json(
      { error: 'Failed to create synonym' },
      { status: 500 }
    );
  }
}
