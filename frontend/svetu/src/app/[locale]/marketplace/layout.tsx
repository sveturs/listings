import { loadMessages } from '@/lib/i18n/loadMessages';
import { NextIntlClientProvider } from 'next-intl';

export default async function MarketplaceLayout({
  children,
  params,
}: {
  children: React.ReactNode;
  params: Promise<{ locale: string }>;
}) {
  const { locale } = await params;

  // Загружаем дополнительные модули для страниц marketplace
  const additionalMessages = await loadMessages(locale as any, [
    'marketplace',
    'reviews',
    'misc',
    'cart',
  ]);

  return (
    <NextIntlClientProvider messages={additionalMessages}>
      {children}
    </NextIntlClientProvider>
  );
}
