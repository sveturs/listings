import { NextRequest, NextResponse } from 'next/server';
import { apiClient } from '@/lib/api-client';

export async function GET(request: NextRequest) {
  try {
    // Получаем токен из заголовка Authorization
    const authHeader = request.headers.get('authorization');
    const accessToken = authHeader?.replace('Bearer ', '');

    if (!accessToken) {
      console.error('No access token in Authorization header');
      return NextResponse.json(
        { error: 'Unauthorized - Please login again' },
        { status: 401 }
      );
    }

    const { searchParams } = new URL(request.url);
    const range = searchParams.get('range') || '7d';

    const response = await apiClient.get(
      `/api/v1/admin/search/analytics?range=${range}`,
      {
        headers: {
          Authorization: `Bearer ${accessToken}`,
        },
      }
    );
    return NextResponse.json(response.data);
  } catch (error: any) {
    console.error('Error fetching search analytics:', error);

    // Проверяем если это ошибка авторизации от backend
    if (error.response?.status === 401) {
      return NextResponse.json(
        { error: 'Unauthorized - Please login again' },
        { status: 401 }
      );
    }

    return NextResponse.json(
      { error: 'Failed to fetch search analytics' },
      { status: 500 }
    );
  }
}
