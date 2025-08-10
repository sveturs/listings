import { loadMessages } from '@/lib/i18n/loadMessages';

export default async function DepositLayout({
  children,
  params,
}: {
  children: React.ReactNode;
  params: Promise<{ locale: string }>;
}) {
  const { locale } = await params;

  // Загружаем необходимые модули для страниц депозита
  const _messages = await loadMessages(locale as any, [
    'admin',
    'misc',
    'profile',
    'common',
  ]);

  return <>{children}</>;
}
