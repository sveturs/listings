import { Metadata } from 'next';
import { getTranslations } from 'next-intl/server';
import VINDecoderClient from './VINDecoderClient';

interface Props {
  params: Promise<{
    locale: string;
  }>;
}

export async function generateMetadata(_: Props): Promise<Metadata> {
  const t = await getTranslations('cars');

  return {
    title: `${t('vinDecoder')} - Sve Tu`,
    description: t('vinDecoderDescription'),
  };
}

export default async function VINDecoderPage({ params }: Props) {
  const { locale } = await params;
  return <VINDecoderClient locale={locale} />;
}
