import { NextRequest, NextResponse } from 'next/server';

export async function POST(request: NextRequest) {
  try {
    const body = await request.json();

    // В режиме разработки просто логируем событие
    if (process.env.NODE_ENV === 'development') {
      console.log('📊 Analytics event:', body);
      return NextResponse.json({ success: true, message: 'Event logged' });
    }

    // В продакшене здесь должна быть интеграция с реальным сервисом аналитики
    // Например, отправка на backend или внешний сервис

    return NextResponse.json({ success: true, message: 'Event processed' });
  } catch (error) {
    console.error('Analytics API error:', error);
    return NextResponse.json(
      { error: 'analytics.error.failed_to_record' },
      { status: 500 }
    );
  }
}
