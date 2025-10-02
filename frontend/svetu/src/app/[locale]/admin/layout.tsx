import { loadMessages } from '@/lib/i18n/loadMessages';
import { NextIntlClientProvider } from 'next-intl';
import AdminLayoutClient from './layout-client';

export default async function AdminLayout({
  children,
  params,
}: {
  children: React.ReactNode;
  params: Promise<{ locale: string }>;
}) {
  const { locale } = await params;

  // Загружаем необходимые модули для админ панели
  const messages = await loadMessages(
    locale as any,
    [
      'admin',
      'misc',
      'common',
      'auth',
      'auth-shared',
      'marketplace', // Для IconPicker компонента
    ] as any
  );

  return (
    <NextIntlClientProvider locale={locale} messages={messages}>
      <AdminLayoutClient>{children}</AdminLayoutClient>
    </NextIntlClientProvider>
  );
}
