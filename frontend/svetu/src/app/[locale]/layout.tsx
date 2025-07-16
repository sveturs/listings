import { Geist, Geist_Mono } from 'next/font/google';
import { notFound } from 'next/navigation';
import { NextIntlClientProvider } from 'next-intl';
import { getMessages, getTranslations } from 'next-intl/server';
import { PublicEnvScript } from 'next-runtime-env';
import { routing } from '@/i18n/routing';
import Header from '@/components/Header';
import { AuthProvider } from '@/contexts/AuthContext';
import { ReduxProvider } from '@/components/ReduxProvider';
import WebSocketManager from '@/components/WebSocketManager';
import AuthStateManager from '@/components/AuthStateManager';
import { VisibleCitiesProvider } from '@/components/GIS/contexts/VisibleCitiesContext';
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

  // Providing all messages to the client
  // side is the easiest way to get started
  const messages = await getMessages();

  return (
    <html lang={locale} data-theme="cupcake">
      <head>
        <PublicEnvScript />
      </head>
      <body
        className={`${geistSans.variable} ${geistMono.variable} antialiased`}
      >
        <NextIntlClientProvider messages={messages}>
          <ReduxProvider>
            <AuthProvider>
              <VisibleCitiesProvider>
                <AuthStateManager />
                <WebSocketManager />
                <Header />
                <main className="min-h-screen pt-28 lg:pt-16">{children}</main>
              </VisibleCitiesProvider>
            </AuthProvider>
          </ReduxProvider>
        </NextIntlClientProvider>
      </body>
    </html>
  );
}
