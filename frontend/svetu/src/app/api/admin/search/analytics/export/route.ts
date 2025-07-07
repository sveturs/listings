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
    const range = searchParams.get('range') || '7d';

    const response = await fetch(
      `${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:3000'}/api/v1/admin/search/analytics/export?range=${range}`,
      {
        headers: {
          Authorization: `Bearer ${accessToken}`,
        },
      }
    );

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const data = await response.text();

    return new NextResponse(data, {
      headers: {
        'Content-Type': 'text/csv',
        'Content-Disposition': `attachment; filename="search-analytics-${range}.csv"`,
      },
    });
  } catch (error) {
    console.error('Error exporting search analytics:', error);
    return NextResponse.json(
      { error: 'Failed to export search analytics' },
      { status: 500 }
    );
  }
}
