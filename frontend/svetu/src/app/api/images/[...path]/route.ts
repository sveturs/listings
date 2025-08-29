import { NextRequest, NextResponse } from 'next/server';

/**
 * API роут для проксирования изображений из MinIO
 * Это позволяет избежать использования IP адресов в URL изображений
 */
export async function GET(
  request: NextRequest,
  { params }: { params: Promise<{ path: string[] }> }
) {
  try {
    // Получаем параметры из Promise
    const { path } = await params;
    // Собираем путь из сегментов
    const imagePath = path.join('/');

    // Определяем базовый URL MinIO на основе окружения
    const minioBaseUrl =
      process.env.NEXT_PUBLIC_MINIO_URL || 'http://localhost:9000';

    // Формируем полный URL изображения
    const imageUrl = `${minioBaseUrl}/${imagePath}`;

    // Делаем запрос к MinIO
    const response = await fetch(imageUrl, {
      headers: {
        // Прокидываем некоторые заголовки из оригинального запроса
        Accept: request.headers.get('accept') || 'image/*',
        Range: request.headers.get('range') || '',
      },
    });

    if (!response.ok) {
      // Если изображение не найдено, возвращаем 404
      if (response.status === 404) {
        return NextResponse.json({ error: 'Image not found' }, { status: 404 });
      }

      // Для других ошибок возвращаем общую ошибку
      return NextResponse.json(
        { error: 'Failed to fetch image' },
        { status: response.status }
      );
    }

    // Получаем тело ответа как ArrayBuffer
    const imageBuffer = await response.arrayBuffer();

    // Создаем ответ с изображением
    const headers = new Headers();

    // Копируем важные заголовки из ответа MinIO
    const contentType = response.headers.get('content-type');
    if (contentType) {
      headers.set('Content-Type', contentType);
    }

    const contentLength = response.headers.get('content-length');
    if (contentLength) {
      headers.set('Content-Length', contentLength);
    }

    // Добавляем заголовки кеширования
    headers.set('Cache-Control', 'public, max-age=31536000, immutable');

    // Добавляем заголовок для корректной обработки изображений
    headers.set('X-Content-Type-Options', 'nosniff');

    return new NextResponse(imageBuffer, {
      status: 200,
      headers,
    });
  } catch (error) {
    console.error('Error proxying image:', error);

    return NextResponse.json(
      { error: 'Internal server error' },
      { status: 500 }
    );
  }
}
