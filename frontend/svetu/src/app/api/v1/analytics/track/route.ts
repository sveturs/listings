import { NextRequest, NextResponse } from 'next/server';

export async function POST(request: NextRequest) {
  try {
    const body = await request.json();

    // В режиме разработки просто логируем события
    if (process.env.NODE_ENV === 'development') {
      console.log('🎯 Behavior tracking events:', {
        batch_id: body.batch_id,
        events_count: body.events?.length || 0,
        events: body.events,
      });

      return NextResponse.json({
        success: true,
        processed_count: body.events?.length || 0,
        message: 'Events logged in development mode',
      });
    }

    // В продакшене здесь должна быть интеграция с реальным сервисом аналитики
    // Например, отправка на backend или внешний сервис аналитики

    return NextResponse.json({
      success: true,
      processed_count: body.events?.length || 0,
      message: 'Events processed',
    });
  } catch (error) {
    console.error('Behavior tracking API error:', error);
    return NextResponse.json(
      { error: 'analytics.error.failed_to_record' },
      { status: 500 }
    );
  }
}
