import { getTranslations } from 'next-intl/server';
import { UnifiedSearchService } from '@/services/unifiedSearch';
import configManager from '@/config';
import HomePageClient from './HomePageClient';
import { getHomePageData } from './actions';

export default async function Home({
  params,
}: {
  params: Promise<{ locale: string }>;
}) {
  const { locale } = await params;
  const t = await getTranslations('home');

  // Проверяем feature flags
  const paymentsEnabled = configManager.isFeatureEnabled('enablePayments');

  let marketplaceData = null;
  let error: Error | null = null;
  let homePageData = null;

  // ВАЖНО: SSR загрузка отключена в development из-за проблем с сетевой конфигурацией
  // В production SSR будет работать нормально когда:
  // 1. Frontend и backend находятся в одной Docker сети
  // 2. Используется правильный INTERNAL_API_URL (например: http://backend:3000)
  // 3. Настроен правильный сетевой доступ между контейнерами
  const skipSSR = process.env.NODE_ENV === 'development';

  if (!skipSSR) {
    try {
      // SSR загрузка данных через унифицированный поиск с таймаутом и обработкой ошибок
      const [marketplaceResult, homePageResult] = await Promise.allSettled([
        UnifiedSearchService.search({
          query: '',
          product_types: ['marketplace', 'storefront'],
          sort_by: 'date',
          sort_order: 'desc',
          page: 1,
          limit: 20,
        }),
        getHomePageData(locale),
      ]);

      if (marketplaceResult.status === 'fulfilled') {
        marketplaceData = marketplaceResult.value;
      } else {
        console.error(
          'SSR marketplace search failed:',
          marketplaceResult.reason
        );
      }

      if (homePageResult.status === 'fulfilled') {
        homePageData = homePageResult.value;
      } else {
        console.error(
          'SSR home page data fetch failed:',
          homePageResult.reason
        );
      }
    } catch (err) {
      error = err as Error;
      console.error('SSR fetch failed:', error);
      // Не падаем, просто загрузим данные на клиенте
    }
  }

  return (
    <HomePageClient
      title={t('marketplace')}
      description={t('description')}
      createListingText={t('createListing')}
      homePageData={homePageData}
      locale={locale}
    />
  );
}
