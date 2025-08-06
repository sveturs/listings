import { loadMessages } from '@/lib/i18n/loadMessages';

export default async function AdminServerLayout({
  children,
  params,
}: {
  children: React.ReactNode;
  params: { locale: string };
}) {
  // Загружаем необходимые модули для админки
  const _messages = await loadMessages(params.locale as any, [
    'common',
    'admin',
    'auth',
  ]);

  return <>{children}</>;
}
