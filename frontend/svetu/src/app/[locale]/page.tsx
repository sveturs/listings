import { getTranslations } from 'next-intl/server';
import { UnifiedSearchService } from '@/services/unifiedSearch';
import configManager from '@/config';
import HomePageClient from './HomePageClient';

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

  // ВАЖНО: SSR загрузка отключена в development из-за проблем с сетевой конфигурацией
  // В production SSR будет работать нормально когда:
  // 1. Frontend и backend находятся в одной Docker сети
  // 2. Используется правильный INTERNAL_API_URL (например: http://backend:3000)
  // 3. Настроен правильный сетевой доступ между контейнерами
  const skipSSR = process.env.NODE_ENV === 'development';

  if (!skipSSR) {
    try {
      // SSR загрузка данных через унифицированный поиск с таймаутом и обработкой ошибок
      marketplaceData = await UnifiedSearchService.search({
        query: '',
        product_types: ['marketplace', 'storefront'],
        sort_by: 'date',
        sort_order: 'desc',
        page: 1,
        limit: 20,
      });
    } catch (err) {
      error = err as Error;
      console.error('SSR unified search fetch failed:', error);
      // Не падаем, просто загрузим данные на клиенте
    }
  }

  return (
    <HomePageClient
      title={t('marketplace')}
      description={t('description')}
      createListingText={t('createListing')}
      initialData={marketplaceData}
      locale={locale}
      error={error}
      paymentsEnabled={paymentsEnabled}
    />
  );
}
