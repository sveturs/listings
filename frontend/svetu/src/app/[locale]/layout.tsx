import { Geist, Geist_Mono } from 'next/font/google';
import { notFound } from 'next/navigation';
import { getTranslations } from 'next-intl/server';
import { routing } from '@/i18n/routing';
import { ModularIntlProvider } from '@/providers/ModularIntlProvider';
import { getMessages } from '@/i18n/messages';
import HeaderWrapper from '@/components/HeaderWrapper';
import { AuthProvider } from '@/contexts/AuthContext';
import { ReduxProvider } from '@/components/ReduxProvider';
import WebSocketManager from '@/components/WebSocketManager';
import AuthStateManager from '@/components/AuthStateManager';
import LocaleSync from '@/components/LocaleSync';
import { VisibleCitiesProvider } from '@/components/GIS/contexts/VisibleCitiesContext';
import { SmartMobileBottomNav } from '@/components/navigation/SmartMobileBottomNav';
import { CartSyncProvider } from '@/components/CartSyncProvider';
import { ThemeProvider } from '@/components/ThemeProvider';
import ErrorBoundaryClass from '@/components/ErrorBoundary';
import { Toaster } from 'react-hot-toast';
import { PublicEnvScript } from 'next-runtime-env';
import '@/utils/forceTokenCleanup'; // Автоматическая очистка старых HS256 токенов
import '../globals.css';

const geistSans = Geist({
  variable: '--font-geist-sans',
  subsets: ['latin'],
  display: 'swap',
});

const geistMono = Geist_Mono({
  variable: '--font-geist-mono',
  subsets: ['latin'],
  display: 'swap',
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
  const t = await getTranslations({ locale, namespace: 'misc' });

  // Определяем альтернативные языковые версии для SEO
  const languages = {
    sr: 'sr-RS',
    ru: 'ru-RU',
    en: 'en-US',
  };

  return {
    title: t('metadata.title'),
    description: t('metadata.description'),
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

  // Используем статические импорты для поддержки SSG
  const messages = getMessages(locale as 'en' | 'ru' | 'sr');

  return (
    <html
      lang={locale}
      suppressHydrationWarning
      className={`${geistSans.variable} ${geistMono.variable}`}
    >
      <head>
        <PublicEnvScript />
      </head>
      <body
        className="antialiased"
        suppressHydrationWarning
        style={
          {
            fontFamily:
              'var(--font-geist-sans), system-ui, -apple-system, sans-serif',
          } as React.CSSProperties
        }
      >
        <ThemeProvider>
          <ModularIntlProvider locale={locale} messages={messages}>
            <ReduxProvider>
              <AuthProvider>
                <VisibleCitiesProvider>
                  <CartSyncProvider>
                    <ErrorBoundaryClass name="RootLayout">
                      <AuthStateManager />
                      <WebSocketManager />
                      <LocaleSync />
                      <HeaderWrapper />
                      <Toaster position="top-right" />
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
        </ThemeProvider>
      </body>
    </html>
  );
}
