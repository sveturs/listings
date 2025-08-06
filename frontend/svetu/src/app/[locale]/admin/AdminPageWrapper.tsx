import { NextIntlClientProvider } from 'next-intl';
import { loadMessages } from '@/lib/i18n/loadMessages';
import AdminPageClient from './AdminPageClient';

export default async function AdminPageWrapper({ locale }: { locale: string }) {
  // Загружаем модуль admin для этой страницы
  const messages = await loadMessages(locale as any, ['admin']);

  return (
    <NextIntlClientProvider messages={messages}>
      <AdminPageClient />
    </NextIntlClientProvider>
  );
}
