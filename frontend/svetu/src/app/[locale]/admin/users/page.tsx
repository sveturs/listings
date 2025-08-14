import { Metadata } from 'next';
import { notFound } from 'next/navigation';
import { getTranslations } from 'next-intl/server';
import { loadMessages } from '@/i18n/loadMessages';
import UsersPageClient from './UsersPageClient';
import AdminGuard from '@/components/AdminGuard';

export async function generateMetadata({
  params,
}: {
  params: Promise<{ locale: string }>;
}): Promise<Metadata> {
  const { locale } = await params;
  const t = await getTranslations({ locale, namespace: 'admin' });
  return {
    title: t('users.pageTitle'),
    description: t('users.pageDescription'),
  };
}

export default async function UsersPage({
  params,
}: {
  params: Promise<{ locale: string }>;
}) {
  const { locale } = await params;
  const messages = await loadMessages(locale as any, ['admin', 'common']);

  if (!messages) {
    notFound();
  }

  return (
    <AdminGuard>
      <UsersPageClient />
    </AdminGuard>
  );
}
