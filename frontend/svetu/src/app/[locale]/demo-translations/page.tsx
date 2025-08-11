import { NextIntlClientProvider } from 'next-intl';
import { loadMessages } from '@/lib/i18n/loadMessages';
import AITranslationsDemo from '@/components/admin/translations/AITranslationsDemo';

export default async function DemoTranslationsPage({
  params,
}: {
  params: Promise<{ locale: string }>;
}) {
  const { locale } = await params;

  // Загружаем модуль admin для страницы переводов
  const messages = await loadMessages(locale as any, ['admin']);

  return (
    <NextIntlClientProvider messages={messages}>
      <div className="container mx-auto px-4 py-8">
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
          <span>Демонстрационная версия системы управления переводами</span>
        </div>
        <AITranslationsDemo />
      </div>
    </NextIntlClientProvider>
  );
}
