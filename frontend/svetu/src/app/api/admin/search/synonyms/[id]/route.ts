import { NextRequest, NextResponse } from 'next/server';
import { apiClient } from '@/lib/api-client';

export async function PUT(
  request: NextRequest,
  { params }: { params: Promise<{ id: string }> }
) {
  try {
    // Получаем токен из заголовка Authorization
    const authHeader = request.headers.get('authorization');
    const accessToken = authHeader?.replace('Bearer ', '');

    if (!accessToken) {
      return NextResponse.json({ error: 'Unauthorized' }, { status: 401 });
    }

    const { id } = await params;
    const data = await request.json();
    const response = await apiClient.put(
      `/api/v1/admin/search/synonyms/${id}`,
      data,
      {
        headers: {
          Authorization: `Bearer ${accessToken}`,
        },
      }
    );
    return NextResponse.json(response.data);
  } catch (error) {
    console.error('Error updating synonym:', error);
    return NextResponse.json(
      { error: 'Failed to update synonym' },
      { status: 500 }
    );
  }
}

export async function DELETE(
  request: NextRequest,
  { params }: { params: Promise<{ id: string }> }
) {
  try {
    // Получаем токен из заголовка Authorization
    const authHeader = request.headers.get('authorization');
    const accessToken = authHeader?.replace('Bearer ', '');

    if (!accessToken) {
      return NextResponse.json({ error: 'Unauthorized' }, { status: 401 });
    }

    const { id } = await params;
    await apiClient.delete(`/api/v1/admin/search/synonyms/${id}`, {
      headers: {
        Authorization: `Bearer ${accessToken}`,
      },
    });
    return NextResponse.json({ success: true });
  } catch (error) {
    console.error('Error deleting synonym:', error);
    return NextResponse.json(
      { error: 'Failed to delete synonym' },
      { status: 500 }
    );
  }
}
