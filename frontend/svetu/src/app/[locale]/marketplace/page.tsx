import { Suspense } from 'react';
import { getTranslations } from 'next-intl/server';
import { Container, Typography } from '@mui/material';
import MarketplaceContent from '@/components/marketplace/MarketplaceContent';
import Loading from './loading';

export default async function MarketplacePage({
  params
}: {
  params: Promise<{ locale: string }>
}) {
  const { locale } = await params;
  const t = await getTranslations({ locale, namespace: 'marketplace' });

  return (
    <Container maxWidth="xl" sx={{ py: 2 }}>
      <Typography variant="h4" component="h1" gutterBottom>
        {t('listings.title')}
      </Typography>

      <Suspense fallback={<Loading />}>
        <MarketplaceContent />
      </Suspense>
    </Container>
  );
}