import { loadMessages } from '@/lib/i18n/loadMessages';

export default async function BalanceLayout({
  children,
  params,
}: {
  children: React.ReactNode;
  params: Promise<{ locale: string }>;
}) {
  const { locale } = await params;

  // Загружаем необходимые модули для страниц баланса
  const _messages = await loadMessages(locale as any, ['admin', 'misc']);

  return <>{children}</>;
}
