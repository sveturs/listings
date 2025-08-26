import { Geist, Geist_Mono } from 'next/font/google';
import { notFound } from 'next/navigation';
import { getTranslations } from 'next-intl/server';
import { PublicEnvScript } from 'next-runtime-env';
import { routing } from '@/i18n/routing';
import { ModularIntlProvider } from '@/providers/ModularIntlProvider';
import { loadMessages } from '@/i18n/loadMessages';
import HeaderWrapper from '@/components/HeaderWrapper';
import { AuthProvider } from '@/contexts/AuthContext';
import { ReduxProvider } from '@/components/ReduxProvider';
import WebSocketManager from '@/components/WebSocketManager';
import AuthStateManager from '@/components/AuthStateManager';
import { VisibleCitiesProvider } from '@/components/GIS/contexts/VisibleCitiesContext';
import { SmartMobileBottomNav } from '@/components/navigation/SmartMobileBottomNav';
import { CartSyncProvider } from '@/components/CartSyncProvider';
import { themeInitScript } from '@/scripts/theme-init';
import ErrorBoundaryClass from '@/components/ErrorBoundary';
import '../globals.css';

const geistSans = Geist({
  variable: '--font-geist-sans',
  subsets: ['latin'],
});

const geistMono = Geist_Mono({
  variable: '--font-geist-mono',
  subsets: ['latin'],
});

export function generateStaticParams() {
  return routing.locales.map((locale) => ({ locale }));
}

export const viewport = {
  width: 'device-width',
  initialScale: 1,
  maximumScale: 1,
  userScalable: false,
};

export async function generateMetadata({
  params,
}: {
  params: Promise<{ locale: string }>;
}) {
  const { locale } = await params;
  const t = await getTranslations({ locale, namespace: 'metadata' });

  // Определяем альтернативные языковые версии для SEO
  const languages = {
    sr: 'sr-RS',
    ru: 'ru-RU',
    en: 'en-US',
  };

  return {
    title: t('title'),
    description: t('description'),
    alternates: {
      languages: Object.fromEntries(
        routing.locales
          .filter((l) => l !== locale)
          .map((l) => [languages[l] || l, `/${l}`])
      ),
      canonical: `/${locale}`,
    },
    openGraph: {
      locale: languages[locale as keyof typeof languages] || locale,
      alternateLocale: routing.locales
        .filter((l) => l !== locale)
        .map((l) => languages[l as keyof typeof languages] || l),
    },
  };
}

export default async function RootLayout({
  children,
  params,
}: Readonly<{
  children: React.ReactNode;
  params: Promise<{ locale: string }>;
}>) {
  const { locale } = await params;

  // Ensure that the incoming locale is valid
  if (!routing.locales.includes(locale as 'en' | 'ru' | 'sr')) {
    notFound();
  }

  // Загружаем базовые модули для layout
  const messages = await loadMessages(locale as any, [
    'common',
    'auth',
    'auth-shared', // Добавляем для AdminGuard и других компонентов защиты
    'balance', // Добавляем для BalanceWidget
    'misc',
    'cart',
    'chat', // Добавляем chat модуль для страницы чата
    'map', // Добавляем map модуль для компонентов GIS
    'marketplace', // И marketplace, так как многие компоненты его используют
    'admin', // Добавляем admin для страниц админки
    'profile', // Добавляем profile для страниц профиля
    'cars', // Добавляем cars для car-selector
    'search', // Добавляем search для SearchBar
    'checkout', // Добавляем checkout для страницы оформления заказа
    'orders', // Добавляем orders для страницы успешного заказа
    'storefronts', // Добавляем storefronts для страниц витрин
    'userContacts', // Добавляем userContacts для страницы контактов
    'subscription', // Добавляем subscription для страниц подписки
  ]);

  return (
    <html lang={locale} suppressHydrationWarning>
      <head>
        <PublicEnvScript />
        <script dangerouslySetInnerHTML={{ __html: themeInitScript }} />
      </head>
      <body
        className={`${geistSans.variable} ${geistMono.variable} antialiased`}
        suppressHydrationWarning
      >
        <ModularIntlProvider locale={locale} messages={messages}>
          <ReduxProvider>
            <AuthProvider>
              <VisibleCitiesProvider>
                <CartSyncProvider>
                  <ErrorBoundaryClass name="RootLayout">
                    <AuthStateManager />
                    <WebSocketManager />
                    <HeaderWrapper />
                    <main className="min-h-screen pt-16 pb-16 md:pb-0">
                      {children}
                    </main>
                    <SmartMobileBottomNav />
                  </ErrorBoundaryClass>
                </CartSyncProvider>
              </VisibleCitiesProvider>
            </AuthProvider>
          </ReduxProvider>
        </ModularIntlProvider>
      </body>
    </html>
  );
}
