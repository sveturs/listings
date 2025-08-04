import { NextIntlClientProvider } from 'next-intl';
import { loadMessages } from '@/lib/i18n/loadMessages';
import AdminPageClient from './AdminPageClient';

export default async function AdminPageWrapper({ locale }: { locale: string }) {
  // Проверяем, используем ли модульную систему
  const useModular = process.env.USE_MODULAR_I18N === 'true';
  
  if (useModular) {
    // Загружаем модуль admin для этой страницы
    const messages = await loadMessages(locale as any, ['admin']);
    
    return (
      <NextIntlClientProvider messages={messages}>
        <AdminPageClient />
      </NextIntlClientProvider>
    );
  }
  
  // Если модульная система отключена, просто рендерим клиентский компонент
  return <AdminPageClient />;
}