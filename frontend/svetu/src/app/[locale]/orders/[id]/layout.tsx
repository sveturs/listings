import { ReactNode } from 'react';
import { loadMessages } from '@/lib/i18n/loadMessages';
import { NextIntlClientProvider } from 'next-intl';

interface Props {
  children: ReactNode;
  params: Promise<{ locale: string }>;
}

export default async function OrderDetailsLayout({ children, params }: Props) {
  const { locale } = await params;

  // Load necessary translation modules for order details page
  const messages = await loadMessages(
    locale as any,
    ['orders', 'common', 'marketplace'] as any
  );

  return (
    <NextIntlClientProvider locale={locale} messages={messages}>
      {children}
    </NextIntlClientProvider>
  );
}
