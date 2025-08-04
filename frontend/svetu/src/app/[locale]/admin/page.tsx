import AdminPageWrapper from './AdminPageWrapper';

export default async function AdminPage({
  params,
}: {
  params: Promise<{ locale: string }>;
}) {
  const { locale } = await params;

  return <AdminPageWrapper locale={locale} />;
}
