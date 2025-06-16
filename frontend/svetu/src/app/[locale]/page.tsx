import { getTranslations } from 'next-intl/server';
import { MarketplaceService } from '@/services/marketplace';
import MarketplaceList from '@/components/marketplace/MarketplaceList';
import { Link } from '@/i18n/routing';
import configManager from '@/config';

export default async function Home({
  params,
}: {
  params: Promise<{ locale: string }>;
}) {
  const { locale } = await params;
  const t = await getTranslations('home');

  // Проверяем feature flags
  const paymentsEnabled = configManager.isFeatureEnabled('enablePayments');
  const config = configManager.getConfig();

  let marketplaceData;
  let error: Error | null = null;

  try {
    // MarketplaceService автоматически использует внутренний URL для SSR
    marketplaceData = await MarketplaceService.search({
      sort_by: 'date_desc',
      page: 0,
      size: 20,
    });
  } catch (err) {
    error = err as Error;
    console.error('Failed to fetch marketplace data:', error);
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-8">{t('marketplace')}</h1>

      {paymentsEnabled && (
        <div className="alert alert-info mb-4">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            className="stroke-current shrink-0 w-6 h-6"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth="2"
              d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
            ></path>
          </svg>
          <span>🎉 {t('paymentsNowAvailable')}</span>
        </div>
      )}

      {/* Debug info в development */}
      {config.env.isDevelopment && (
        <div className="text-xs text-base-content/50 mb-2">
          API: {config.api.url} | Internal:{' '}
          {config.api.internalUrl || 'not set'}
        </div>
      )}

      {error && (
        <div className="alert alert-error mb-8">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            className="stroke-current shrink-0 h-6 w-6"
            fill="none"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth="2"
              d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"
            />
          </svg>
          <span>{t('errorLoadingData')}</span>
        </div>
      )}

      {marketplaceData && marketplaceData.data.length > 0 ? (
        <MarketplaceList initialData={marketplaceData} locale={locale} />
      ) : (
        !error && (
          <div className="text-center py-12">
            <p className="text-lg text-base-content/70">{t('noItems')}</p>
          </div>
        )
      )}

      {/* Плавающая кнопка создания объявления */}
      <Link
        href="/create-listing"
        className="fixed bottom-6 right-6 btn btn-primary btn-circle btn-lg shadow-lg hover:shadow-xl transition-shadow z-50"
        title={t('createListing')}
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          className="h-6 w-6"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M12 4v16m8-8H4"
          />
        </svg>
      </Link>
    </div>
  );
}
