import type { Metadata } from "next";
import { Geist, Geist_Mono } from "next/font/google";
import { NextIntlClientProvider } from 'next-intl';
import { getMessages } from 'next-intl/server';
import { getTranslations } from 'next-intl/server';
import Link from 'next/link';
import LanguageSwitcher from '@/components/LanguageSwitcher';
import { Providers } from '@/components/providers';
import "../globals.css";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title: "Svetu",
  description: "Find your perfect place",
};

export default async function LocaleLayout({
  children,
  params,
}: {
  children: React.ReactNode;
  params: Promise<{ locale: string }>;
}) {
  const { locale } = await params;
  const messages = await getMessages();
  const t = await getTranslations('nav');

  return (
    <html lang={locale}>
      <body className={`${geistSans.variable} ${geistMono.variable} antialiased`}>
        <Providers>
          <NextIntlClientProvider messages={messages}>
            <nav className="bg-gray-800 text-white p-4">
              <div className="container mx-auto flex justify-between items-center">
                <Link href={`/${locale}`} className="text-xl font-bold">Svetu</Link>
                <div className="flex gap-4 items-center">
                  <Link href={`/${locale}`} className="hover:text-gray-300">{t('home')}</Link>
                  <Link href={`/${locale}/marketplace`} className="hover:text-gray-300">{t('marketplace')}</Link>
                  <LanguageSwitcher currentLocale={locale} />
                </div>
              </div>
            </nav>
            <main className="container mx-auto p-4">
              {children}
            </main>
          </NextIntlClientProvider>
        </Providers>
      </body>
    </html>
  );
}