import AdminPageWrapper from './AdminPageWrapper';
import type { Metadata } from 'next';

export const metadata: Metadata = {
  title: 'Admin Dashboard - SveTu',
  description: 'Administration panel for SveTu marketplace',
};

export default async function AdminPage({
  params,
}: {
  params: Promise<{ locale: string }>;
}) {
  const { locale } = await params;

  return <AdminPageWrapper locale={locale} />;
}
