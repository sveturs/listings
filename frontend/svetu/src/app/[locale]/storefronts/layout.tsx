import { loadMessages } from '@/lib/i18n/loadMessages';
import { NextIntlClientProvider } from 'next-intl';

export default async function StorefrontsLayout({
  children,
  params,
}: {
  children: React.ReactNode;
  params: Promise<{ locale: string }>;
}) {
  const { locale } = await params;

  // Загружаем дополнительные модули для витрин
  const additionalMessages = await loadMessages(locale as any, [
    'storefronts',
    'products',
    'reviews',
    'admin',
    'notifications',
    'cart',
    'orders',
  ]);

  return (
    <NextIntlClientProvider messages={additionalMessages}>
      {children}
    </NextIntlClientProvider>
  );
}
