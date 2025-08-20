import { Metadata } from 'next';
import { getTranslations } from 'next-intl/server';
import UsersPageClient from './UsersPageClient';
import AdminGuard from '@/components/AdminGuard';

// Отключаем статическую генерацию для admin страниц
export const dynamic = 'force-dynamic';

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

export default async function UsersPage() {
  return (
    <AdminGuard>
      <UsersPageClient />
    </AdminGuard>
  );
}
