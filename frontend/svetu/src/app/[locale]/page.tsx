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

  let marketplaceData;
  let error: Error | null = null;

  try {
    // MarketplaceService автоматически использует внутренний URL для SSR
    marketplaceData = await MarketplaceService.search({
      sort_by: 'date_desc',
      page: 0,
      size: 20,
    });

    // Отладка: проверим объявление 177
    const item177 = marketplaceData?.data?.find((item) => item.id === 177);
    if (item177) {
      console.log('SSR - Item 177 found:', {
        id: item177.id,
        title: item177.title,
        images: item177.images,
        images_count: item177.images?.length || 0,
      });
    } else {
      console.log('SSR - Item 177 not found in initial data');
    }
  } catch (err) {
    error = err as Error;
    console.error('Failed to fetch marketplace data:', error);
  }

  return (
    <div className="min-h-screen">
      {/* Hero секция */}
      <div className="bg-gradient-to-b from-base-200/50 to-base-100 py-12 lg:py-16 mb-8">
        <div className="container mx-auto px-4">
          <h1 className="text-4xl lg:text-5xl font-bold text-center mb-4">
            {t('marketplace')}
          </h1>
          <p className="text-center text-base-content/70 text-lg max-w-2xl mx-auto">
            Находите лучшие предложения от частных лиц и магазинов
          </p>
        </div>
      </div>

      <div className="container mx-auto px-4">
        {paymentsEnabled && (
          <div className="alert alert-info mb-8 shadow-lg">
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

        {error && (
          <div className="alert alert-error mb-8 shadow-lg">
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
            <div className="text-center py-16">
              <div className="inline-flex items-center justify-center w-24 h-24 rounded-full bg-base-200 mb-6">
                <svg
                  className="w-12 h-12 text-base-content/40"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={1.5}
                    d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4"
                  />
                </svg>
              </div>
              <p className="text-xl text-base-content/70 mb-4">
                {t('noItems')}
              </p>
              <p className="text-base-content/50">
                Начните с создания первого объявления
              </p>
            </div>
          )
        )}

        {/* Плавающая кнопка создания объявления */}
        <Link
          href="/create-listing"
          className="fixed bottom-6 right-6 btn btn-primary btn-circle btn-lg shadow-xl hover:shadow-2xl hover:scale-110 transition-all duration-200 z-50"
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
    </div>
  );
}
