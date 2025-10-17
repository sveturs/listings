import QualityTestsClient from './QualityTestsClient';

export default async function QualityTestsPage({
  params,
}: {
  params: Promise<{ locale: string }>;
}) {
  const { locale } = await params;

  return <QualityTestsClient locale={locale} />;
}
