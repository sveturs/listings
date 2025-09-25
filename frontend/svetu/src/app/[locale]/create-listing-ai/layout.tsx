import { loadMessages } from '@/lib/i18n/loadMessages';
import { NextIntlClientProvider } from 'next-intl';

export default async function CreateListingAILayout({
  children,
  params,
}: {
  children: React.ReactNode;
  params: Promise<{ locale: string }>;
}) {
  const { locale } = await params;

  // Загружаем дополнительные модули для создания объявления и работы с автомобилями
  const messages = await loadMessages(locale as any, ['create_listing', 'cars']);

  return (
    <NextIntlClientProvider locale={locale} messages={messages}>
      {children}
    </NextIntlClientProvider>
  );
}
