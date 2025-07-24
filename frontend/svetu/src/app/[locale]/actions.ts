'use server';

import { apiClientServer } from '@/lib/api-client-server';

interface CategoryWithCount {
  id: string;
  name: string;
  count: number;
}

interface FeaturedListing {
  id: string;
  title: string;
  price: string;
  image: string;
  category: string;
}

interface MarketplaceStats {
  totalListings: number;
  activeUsers: number;
  successfulDeals: number;
  newUsersToday?: number;
  nearbyListings?: number;
  newListingsLastHour?: number;
  priceDropsToday?: number;
}

export async function getHomePageData(locale: string) {
  console.log('[getHomePageData] Starting to fetch data for locale:', locale);

  try {
    // Получаем категории
    const categoriesResponse = await apiClientServer.get(
      `/api/v1/marketplace/categories?lang=${locale}`
    );
    console.log(
      '[getHomePageData] Categories response:',
      categoriesResponse.data
    );

    // Получаем популярные поисковые запросы
    const popularSearchesResponse = await apiClientServer.get(
      '/api/v1/search/statistics/popular?limit=10'
    );
    console.log(
      '[getHomePageData] Popular searches response:',
      popularSearchesResponse.data
    );

    // Получаем новые объявления для featured секции
    const featuredResponse = await apiClientServer.post(
      '/api/v1/marketplace/search',
      {
        sort: 'created_at',
        sortDirection: 'desc',
        limit: 1,
        offset: 0,
      }
    );
    console.log(
      '[getHomePageData] Featured listing response:',
      featuredResponse.data
    );

    // Получаем топ просматриваемые товары (для будущего использования)
    // const topItemsResponse = await apiClientServer.get(
    //   '/api/v1/analytics/metrics/items?item_type=marketplace&sort_by=views&limit=5'
    // );

    // Обрабатываем категории и добавляем количество (пока моковые данные)
    const categoriesWithCount: CategoryWithCount[] =
      categoriesResponse.data?.data?.slice(0, 5).map((cat: any) => ({
        id: cat.id,
        name: cat.name,
        count: Math.floor(Math.random() * 1000) + 100, // Временно, пока нет реального подсчета
      })) || [];

    // Обрабатываем featured listing
    let featuredListing: FeaturedListing | undefined;

    // Проверяем разные возможные структуры ответа
    const listings = featuredResponse.data?.data || featuredResponse.data || [];
    console.log('[getHomePageData] Listings array:', listings);

    if (Array.isArray(listings) && listings.length > 0) {
      const listing = listings[0];
      const currencySymbol = locale === 'ru' ? '₽' : 'RSD';

      featuredListing = {
        id: listing.id,
        title: listing.title || 'Без названия',
        price: listing.price
          ? `${listing.price.toLocaleString()} ${currencySymbol}`
          : 'Цена не указана',
        image:
          listing.images?.[0]?.url ||
          listing.images?.[0]?.public_url ||
          (listing.images?.[0]?.public_url
            ? `${process.env.NEXT_PUBLIC_MINIO_URL || 'http://localhost:9000'}${listing.images[0].public_url}`
            : '/api/placeholder/300/200'),
        category:
          listing.category_name || listing.category?.name || 'Без категории',
      };
    }

    // Получаем все объявления для анализа статистики
    const allListingsResponse = await apiClientServer.post(
      '/api/v1/marketplace/search',
      {
        limit: 100,
        offset: 0,
      }
    );

    const allListings = allListingsResponse.data?.data || [];
    const totalListings = allListingsResponse.data?.meta?.total || 0;

    // Подсчитываем уникальных пользователей
    const uniqueUserIds = new Set(
      allListings.map((listing: any) => listing.user_id)
    );
    const activeUsersCount = uniqueUserIds.size;

    // Подсчитываем объявления за последний час и день
    const now = new Date();
    const hourAgo = new Date(now.getTime() - 60 * 60 * 1000);
    const dayAgo = new Date(now.getTime() - 24 * 60 * 60 * 1000);

    let newListingsLastHour = 0;
    let newListingsToday = 0;
    let priceChangesToday = 0;

    allListings.forEach((listing: any) => {
      const createdAt = new Date(listing.created_at);
      const updatedAt = new Date(listing.updated_at);

      if (createdAt > hourAgo) {
        newListingsLastHour++;
      }

      if (createdAt > dayAgo) {
        newListingsToday++;
      }

      // Если updated_at отличается от created_at и обновлено сегодня
      if (updatedAt > dayAgo && updatedAt.getTime() !== createdAt.getTime()) {
        priceChangesToday++;
      }
    });

    // Получаем количество объявлений рядом (все, так как нет геолокации)
    // В будущем можно будет фильтровать по расстоянию
    const nearbyListings = totalListings;

    // Подсчитываем новых пользователей сегодня (тех, кто создал объявления сегодня)
    const todayUserIds = new Set(
      allListings
        .filter((listing: any) => new Date(listing.created_at) > dayAgo)
        .map((listing: any) => listing.user_id)
    );
    const newUsersToday = todayUserIds.size;

    // Считаем успешные сделки как все объявления, кроме созданных сегодня
    // (предполагаем, что большинство старых объявлений = успешные сделки)
    const successfulDeals = totalListings - newListingsToday;

    // Статистика на основе реальных данных
    const stats: MarketplaceStats = {
      totalListings,
      activeUsers: activeUsersCount,
      successfulDeals,
      newUsersToday,
      nearbyListings,
      newListingsLastHour,
      priceDropsToday: priceChangesToday,
    };

    const result = {
      categories: categoriesWithCount,
      featuredListing,
      stats,
      popularSearches: popularSearchesResponse.data?.data || [],
      error: null,
    };

    console.log('[getHomePageData] Final result:', result);
    return result;
  } catch (error) {
    console.error('Error fetching home page data:', error);
    return {
      categories: [],
      featuredListing: undefined,
      stats: {
        totalListings: 0,
        activeUsers: 0,
        successfulDeals: 0,
      },
      popularSearches: [],
      error: error as Error,
    };
  }
}
