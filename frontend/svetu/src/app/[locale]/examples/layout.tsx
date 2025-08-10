import { loadMessages } from '@/lib/i18n/loadMessages';

export default async function ExamplesLayout({
  children,
  params,
}: {
  children: React.ReactNode;
  params: Promise<{ locale: string }>;
}) {
  const { locale } = await params;

  // Загружаем необходимые модули для примеров
  const _messages = await loadMessages(locale as any, [
    'admin',
    'misc',
    'common',
    'storefronts',
    'marketplace',
    'auth',
  ]);

  return <>{children}</>;
}
