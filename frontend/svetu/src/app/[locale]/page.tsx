import { getTranslations } from 'next-intl/server';
import { MarketplaceService } from '@/services/marketplace';
import MarketplaceCard from '@/components/MarketplaceCard';

export default async function Home({
  params,
}: {
  params: Promise<{ locale: string }>;
}) {
  const { locale } = await params;
  const t = await getTranslations('home');

  let marketplaceData;
  let error: Error | null = null;

  try {
    marketplaceData = await MarketplaceService.search({
      sort_by: 'date_desc',
      page: 0,
      size: 25,
    });
  } catch (err) {
    error = err as Error;
    console.error('Failed to fetch marketplace data:', error);
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-8">{t('marketplace')}</h1>

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
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
          {marketplaceData.data.map((item) => (
            <MarketplaceCard key={item.id} item={item} locale={locale} />
          ))}
        </div>
      ) : (
        !error && (
          <div className="text-center py-12">
            <p className="text-lg text-base-content/70">{t('noItems')}</p>
          </div>
        )
      )}

      {marketplaceData && marketplaceData.meta.has_more && (
        <div className="text-center mt-8">
          <button className="btn btn-primary">{t('loadMore')}</button>
        </div>
      )}
    </div>
  );
}
