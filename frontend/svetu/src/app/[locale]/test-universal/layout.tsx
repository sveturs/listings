import { loadMessages } from '@/i18n/loadMessages';
import { type Metadata } from 'next';

export const metadata: Metadata = {
  title: 'Test Universal Components',
  description: 'Testing page for universal marketplace components',
};

export default async function TestUniversalLayout({
  children,
  params,
}: {
  children: React.ReactNode;
  params: Promise<{ locale: string }>;
}) {
  const { locale } = await params;

  // Загружаем необходимые модули переводов для универсальных компонентов
  const _messages = await loadMessages(locale as any, [
    'misc',
    'filters',
    'calculator',
    'recommendations',
    'common',
    'nav',
  ]);

  return <>{children}</>;
}
