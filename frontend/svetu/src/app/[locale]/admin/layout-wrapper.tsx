import { loadMessages } from '@/lib/i18n/loadMessages';

export default async function AdminLayoutWrapper({
  children,
  params,
}: {
  children: React.ReactNode;
  params: Promise<{ locale: string }>;
}) {
  const { locale } = await params;

  // Загружаем необходимые модули для админ панели
  const _messages = await loadMessages(locale as any, [
    'admin',
    'misc',
    'common',
  ]);

  return <>{children}</>;
}
