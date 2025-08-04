import { getTranslations } from 'next-intl/server';
import { NextIntlClientProvider } from 'next-intl';
import { UnifiedSearchService } from '@/services/unifiedSearch';
import configManager from '@/config';
import HomePageClient from './HomePageClient';
import { getHomePageData } from './actions';
import { loadMessages } from '@/lib/i18n/loadMessages';

export default async function Home({
  params,
}: {
  params: Promise<{ locale: string }>;
}) {
  const { locale } = await params;
  const t = await getTranslations('home');

  // Проверяем, используем ли модульную систему
  const useModular = process.env.USE_MODULAR_I18N === 'true';
  
  // Загружаем необходимые модули для главной страницы
  let additionalMessages = {};
  if (useModular) {
    additionalMessages = await loadMessages(locale as any, ['marketplace']);
  }

  // Проверяем feature flags
  const _paymentsEnabled = configManager.isFeatureEnabled('enablePayments');

  const _marketplaceData = null;
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
      homePageData = await getHomePageData(locale);
    } catch (e) {
      console.error('[SSR] Failed to load homepage data:', e);
      error = e instanceof Error ? e : new Error('Failed to load data');
    }
  }

  return useModular ? (
    <NextIntlClientProvider messages={additionalMessages}>
      <HomePageClient
        initialData={homePageData}
        error={error}
        locale={locale}
        skipSSR={skipSSR}
      />
    </NextIntlClientProvider>
  ) : (
    <HomePageClient
      initialData={homePageData}
      error={error}
      locale={locale}
      skipSSR={skipSSR}
    />
  );
}