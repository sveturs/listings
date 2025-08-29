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

  // SSR загрузка данных для главной страницы
  // Использует INTERNAL_API_URL для связи между контейнерами Docker
  try {
    homePageData = await getHomePageData(locale);
  } catch (e) {
    console.error('[SSR] Failed to load homepage data:', e);
    _error = e instanceof Error ? e : new Error('Failed to load data');
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
