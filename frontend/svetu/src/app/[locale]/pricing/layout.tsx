import { loadMessages } from '@/lib/i18n/loadMessages';
import { NextIntlClientProvider } from 'next-intl';

export default async function PricingLayout({
  children,
  params,
}: {
  children: React.ReactNode;
  params: Promise<{ locale: string }>;
}) {
  const { locale } = await params;

  // Загружаем необходимые модули для pricing
  const messages = await loadMessages(locale as any, [
    'subscription',
    'misc',
    'common',
  ]);

  return (
    <NextIntlClientProvider locale={locale} messages={messages}>
      {children}
    </NextIntlClientProvider>
  );
}