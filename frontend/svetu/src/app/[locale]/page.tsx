import { getTranslations } from 'next-intl/server';
import { UnifiedSearchService } from '@/services/unifiedSearch';
import HomePage from '@/components/marketplace/HomePage';
import { Link } from '@/i18n/routing';
import configManager from '@/config';
import { SearchBar } from '@/components/SearchBar';

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
    <div className="min-h-screen">
      {/* Hero секция */}
      <div className="bg-gradient-to-b from-base-200/50 to-base-100 py-12 lg:py-16 mb-8">
        <div className="container mx-auto px-4">
          <h1 className="text-4xl lg:text-5xl font-bold text-center mb-4">
            {t('marketplace')}
          </h1>
          <p className="text-center text-base-content/70 text-lg max-w-2xl mx-auto mb-8">
            {t('description')}
          </p>

          {/* Search Bar с поддержкой fuzzy search */}
          <div className="max-w-3xl mx-auto">
            <SearchBar variant="hero" showTrending={true} />
          </div>
        </div>
      </div>

      <div className="container mx-auto px-4">
        <HomePage
          initialData={marketplaceData}
          locale={locale}
          error={error}
          paymentsEnabled={paymentsEnabled}
        />

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
