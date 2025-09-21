import { NextIntlClientProvider } from 'next-intl';
import { loadMessages } from '@/lib/i18n/loadMessages';
import DeliveryAdminClient from './DeliveryAdminClient';

export default async function DeliveryAdminPage({
  params,
}: {
  params: Promise<{ locale: string }>;
}) {
  const { locale } = await params;
  const messages = await loadMessages(locale as any, ['admin']);

  return (
    <NextIntlClientProvider messages={messages}>
      <DeliveryAdminClient />
    </NextIntlClientProvider>
  );
}
