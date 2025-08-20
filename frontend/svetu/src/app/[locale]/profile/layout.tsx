import { NextIntlClientProvider } from 'next-intl';
import { loadMessages } from '@/i18n/loadMessages';

export default async function ProfileLayout({
  children,
  params,
}: {
  children: React.ReactNode;
  params: Promise<{ locale: string }>;
}) {
  const { locale } = await params;

  // Загружаем дополнительные модули для профиля
  const messages = await loadMessages(locale as any, [
    'profile',
    'storefronts',
    'admin',
    'orders',
    'balance', // Добавляем модуль balance для BalanceWidget
  ]);

  return (
    <NextIntlClientProvider locale={locale} messages={messages}>
      {children}
    </NextIntlClientProvider>
  );
}
