'use client';

import { MarketplaceItem, MarketplaceImage } from '@/types/marketplace';
import SafeImage from '@/components/SafeImage';
import Link from 'next/link';
import { useAuth } from '@/contexts/AuthContext';
import { useRouter } from 'next/navigation';
import { useState, useEffect } from 'react';
import { useTranslations } from 'next-intl';

interface MarketplaceCardProps {
  item: MarketplaceItem;
  locale: string;
  viewMode?: 'grid' | 'list';
}

export default function MarketplaceCard({
  item,
  locale,
  viewMode = 'grid',
}: MarketplaceCardProps) {
  const router = useRouter();
  const [mounted, setMounted] = useState(false);
  const t = useTranslations('common');

  // Всегда вызываем хук, но используем данные только после монтирования
  const { user, isAuthenticated } = useAuth();

  useEffect(() => {
    setMounted(true);
  }, []);

  const formatPrice = (price?: number, currency?: string) => {
    if (!price) return '';

    const formatter = new Intl.NumberFormat(locale, {
      style: 'currency',
      currency: currency || 'USD',
      minimumFractionDigits: 0,
      maximumFractionDigits: 0,
    });

    return formatter.format(price);
  };

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    const now = new Date();
    const diffInMs = now.getTime() - date.getTime();

    // Если дата в будущем, просто показываем её
    if (diffInMs < 0) {
      return date.toLocaleDateString(locale);
    }

    const diffInHours = Math.floor(diffInMs / (1000 * 60 * 60));

    if (diffInHours < 1) return t('justNow');
    if (diffInHours < 24) return `${diffInHours} ${t('hoursAgo')}`;

    const diffInDays = Math.floor(diffInHours / 24);
    if (diffInDays < 7) {
      return t('daysAgoWithCount', { count: diffInDays });
    }

    return date.toLocaleDateString(locale);
  };

  const getImageUrl = (image?: MarketplaceImage) => {
    if (!image) return null;
    // Возвращаем public_url как есть, SafeImage сам обработает путь
    return image.public_url;
  };

  const mainImage = item.images?.find((img) => img.is_main) || item.images?.[0];
  const imageUrl = getImageUrl(mainImage);

  // Отладка для объявления 177
  if (item.id === 177) {
    console.log('MarketplaceCard - Item 177 debug:', {
      item,
      mainImage,
      imageUrl,
    });
  }

  const handleChatClick = (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();

    if (!isAuthenticated) {
      // Если не авторизован, можно перенаправить на авторизацию
      // Или показать сообщение
      return;
    }

    // Проверяем, что это не собственное объявление
    if (item.user_id === user?.id) {
      return;
    }

    // Переходим на страницу чата с параметрами для создания чата
    router.push(
      `/${locale}/chat?listing_id=${item.id}&seller_id=${item.user_id}`
    );
  };

  if (viewMode === 'list') {
    return (
      <div className="card bg-base-100 shadow-xl hover:shadow-2xl transition-shadow">
        <Link href={`/${locale}/marketplace/${item.id}`} className="block">
          <div className="card-body p-4">
            <div className="flex gap-4">
              {/* Изображение слева */}
              <figure className="relative w-20 h-20 sm:w-32 sm:h-32 flex-shrink-0 bg-base-200 rounded-lg overflow-hidden">
                <SafeImage
                  src={imageUrl}
                  alt={item.title}
                  fill
                  className="object-cover"
                  sizes="128px"
                  fallback={
                    <div className="flex items-center justify-center w-full h-full text-base-content/50">
                      <svg
                        xmlns="http://www.w3.org/2000/svg"
                        fill="none"
                        viewBox="0 0 24 24"
                        strokeWidth={1.5}
                        stroke="currentColor"
                        className="w-12 h-12"
                      >
                        <path
                          strokeLinecap="round"
                          strokeLinejoin="round"
                          d="m2.25 15.75 5.159-5.159a2.25 2.25 0 0 1 3.182 0l5.159 5.159m-1.5-1.5 1.409-1.409a2.25 2.25 0 0 1 3.182 0l2.909 2.909m-18 3.75h16.5a1.5 1.5 0 0 0 1.5-1.5V6a1.5 1.5 0 0 0-1.5-1.5H3.75A1.5 1.5 0 0 0 2.25 6v12a1.5 1.5 0 0 0 1.5 1.5Zm10.5-11.25h.008v.008h-.008V8.25Zm.375 0a.375.375 0 1 1-.75 0 .375.375 0 0 1 .75 0Z"
                        />
                      </svg>
                    </div>
                  }
                />
              </figure>

              {/* Информация */}
              <div className="flex-grow">
                <div className="flex justify-between items-start gap-4">
                  <div className="flex-grow">
                    <h2 className="text-base sm:text-lg font-semibold line-clamp-1">
                      {item.title}
                    </h2>
                    {item.description && (
                      <p className="text-sm text-base-content/70 line-clamp-2 mt-1">
                        {item.description}
                      </p>
                    )}
                    <div className="flex items-center gap-4 mt-2 text-sm text-base-content/70">
                      {item.location && (
                        <span className="flex items-center gap-1">
                          <svg
                            xmlns="http://www.w3.org/2000/svg"
                            fill="none"
                            viewBox="0 0 24 24"
                            strokeWidth={1.5}
                            stroke="currentColor"
                            className="w-4 h-4"
                          >
                            <path
                              strokeLinecap="round"
                              strokeLinejoin="round"
                              d="M15 10.5a3 3 0 1 1-6 0 3 3 0 0 1 6 0Z"
                            />
                            <path
                              strokeLinecap="round"
                              strokeLinejoin="round"
                              d="M19.5 10.5c0 7.142-7.5 11.25-7.5 11.25S4.5 17.642 4.5 10.5a7.5 7.5 0 1 1 15 0Z"
                            />
                          </svg>
                          {item.location}
                        </span>
                      )}
                      <span className="hidden sm:inline">
                        {formatDate(item.created_at)}
                      </span>
                    </div>
                  </div>

                  {/* Цена и кнопка чата */}
                  <div className="flex flex-col items-end gap-2">
                    {item.price && (
                      <div className="text-right">
                        {item.has_discount && item.old_price && (
                          <p className="text-sm line-through text-base-content/50">
                            {formatPrice(item.old_price, 'RSD')}
                          </p>
                        )}
                        <p className="text-lg sm:text-xl font-bold text-primary">
                          {formatPrice(item.price, 'RSD')}
                        </p>
                      </div>
                    )}
                    {mounted &&
                      isAuthenticated &&
                      item.user_id !== user?.id && (
                        <button
                          onClick={handleChatClick}
                          className="btn btn-primary btn-sm"
                          title={t('sendMessage')}
                        >
                          <svg
                            xmlns="http://www.w3.org/2000/svg"
                            fill="none"
                            viewBox="0 0 24 24"
                            strokeWidth={1.5}
                            stroke="currentColor"
                            className="w-4 h-4"
                          >
                            <path
                              strokeLinecap="round"
                              strokeLinejoin="round"
                              d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"
                            />
                          </svg>
                          {t('chat')}
                        </button>
                      )}
                  </div>
                </div>
              </div>
            </div>
          </div>
        </Link>
      </div>
    );
  }

  // Grid view (default)
  return (
    <div className="card card-compact bg-base-100 shadow-xl hover:shadow-2xl transition-shadow relative">
      {/* Кнопка чата - показываем только после монтирования */}
      {mounted && isAuthenticated && item.user_id !== user?.id && (
        <button
          onClick={handleChatClick}
          className="btn btn-primary btn-circle btn-sm absolute top-3 right-3 shadow-lg z-10"
          title={t('sendMessage')}
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            strokeWidth={1.5}
            stroke="currentColor"
            className="w-4 h-4"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"
            />
          </svg>
        </button>
      )}

      <Link href={`/${locale}/marketplace/${item.id}`} className="block">
        <figure className="relative h-48 bg-base-200">
          <SafeImage
            src={imageUrl}
            alt={item.title}
            fill
            className="object-cover"
            sizes="(max-width: 768px) 100vw, (max-width: 1200px) 50vw, 33vw"
            fallback={
              <div className="flex items-center justify-center w-full h-full text-base-content/50">
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  fill="none"
                  viewBox="0 0 24 24"
                  strokeWidth={1.5}
                  stroke="currentColor"
                  className="w-16 h-16"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    d="m2.25 15.75 5.159-5.159a2.25 2.25 0 0 1 3.182 0l5.159 5.159m-1.5-1.5 1.409-1.409a2.25 2.25 0 0 1 3.182 0l2.909 2.909m-18 3.75h16.5a1.5 1.5 0 0 0 1.5-1.5V6a1.5 1.5 0 0 0-1.5-1.5H3.75A1.5 1.5 0 0 0 2.25 6v12a1.5 1.5 0 0 0 1.5 1.5Zm10.5-11.25h.008v.008h-.008V8.25Zm.375 0a.375.375 0 1 1-.75 0 .375.375 0 0 1 .75 0Z"
                  />
                </svg>
              </div>
            }
          />
        </figure>

        <div className="card-body">
          <h2 className="card-title line-clamp-2">{item.title}</h2>
          {item.price && (
            <div>
              {item.has_discount && item.old_price && (
                <p className="text-sm line-through text-base-content/50">
                  {formatPrice(item.old_price, 'RSD')}
                </p>
              )}
              <p className="text-lg font-bold text-primary">
                {formatPrice(item.price, 'RSD')}
              </p>
            </div>
          )}
          {item.location && (
            <p className="text-sm text-base-content/70 flex items-center gap-1">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
                strokeWidth={1.5}
                stroke="currentColor"
                className="w-4 h-4"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  d="M15 10.5a3 3 0 1 1-6 0 3 3 0 0 1 6 0Z"
                />
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  d="M19.5 10.5c0 7.142-7.5 11.25-7.5 11.25S4.5 17.642 4.5 10.5a7.5 7.5 0 1 1 15 0Z"
                />
              </svg>
              {item.location}
            </p>
          )}
          <p className="text-xs text-base-content/50">
            {formatDate(item.created_at)}
          </p>
        </div>
      </Link>
    </div>
  );
}
