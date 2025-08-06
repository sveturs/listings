import { getTranslations } from 'next-intl/server';
import configManager from '@/config';
import HomePageClient from './HomePageClient';
import { getHomePageData } from './actions';

export default async function Home({
  params,
}: {
  params: Promise<{ locale: string }>;
}) {
  const { locale } = await params;
  const t = await getTranslations('marketplace.home');

  // Проверяем feature flags
  const _paymentsEnabled = configManager.isFeatureEnabled('enablePayments');

  const _marketplaceData = null;
  let _error: Error | null = null;
  let homePageData = null;

  // ВАЖНО: SSR загрузка отключена в development из-за проблем с сетевой конфигурацией
  // В production SSR будет работать нормально когда:
  // 1. Frontend и backend находятся в одной Docker сети
  // 2. Используется правильный INTERNAL_API_URL (например: http://backend:3000)
  // 3. Настроен правильный сетевой доступ между контейнерами
  const skipSSR = process.env.NODE_ENV === 'development';

  if (!skipSSR) {
    try {
      homePageData = await getHomePageData(locale);
    } catch (e) {
      console.error('[SSR] Failed to load homepage data:', e);
      _error = e instanceof Error ? e : new Error('Failed to load data');
    }
  }

  return (
    <HomePageClient
      title={t('title')}
      description={t('description')}
      createListingText={t('createListing')}
      homePageData={homePageData}
      locale={locale}
    />
  );
}
