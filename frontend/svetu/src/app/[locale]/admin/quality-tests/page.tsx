import QualityTestsClient from './QualityTestsClient';
import type { Metadata } from 'next';

export const metadata: Metadata = {
  title: 'Quality Tests - Admin - SveTu',
  description: 'WCAG 2.1 AA compliance and quality testing dashboard',
};

export default async function QualityTestsPage({
  params,
}: {
  params: Promise<{ locale: string }>;
}) {
  const { locale } = await params;

  return <QualityTestsClient locale={locale} />;
}
