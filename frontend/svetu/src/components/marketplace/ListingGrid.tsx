'use client';

import { useQuery } from '@tanstack/react-query';
import { useTranslations } from 'next-intl';
import { listingService } from '@/services/listing.service';
import ListingCard from './ListingCard';
import { Alert, Box, Skeleton } from '@mui/material';

export default function ListingGrid() {
  const t = useTranslations('marketplace.listings');
  
  const { data, isLoading, error } = useQuery({
    queryKey: ['listings'],
    queryFn: () => listingService.getListings(),
  });

  if (isLoading) {
    return (
      <Box
        sx={{
          display: 'grid',
          gridTemplateColumns: {
            xs: '1fr',
            sm: 'repeat(2, 1fr)',
            md: 'repeat(3, 1fr)',
          },
          gap: 3,
        }}
      >
        {[...Array(6)].map((_, i) => (
          <Box key={i}>
            <Skeleton variant="rectangular" height={300} />
          </Box>
        ))}
      </Box>
    );
  }

  if (error) {
    return <Alert severity="error">{t('error')}</Alert>;
  }

  if (!data?.items || data.items.length === 0) {
    return (
      <Box textAlign="center" py={8}>
        <Alert severity="info">{t('noListingsFound')}</Alert>
      </Box>
    );
  }

  return (
    <>
      <Box mb={2}>
        <p className="text-gray-600">{t('found', { count: data.total })}</p>
      </Box>
      <Box
        sx={{
          display: 'grid',
          gridTemplateColumns: {
            xs: '1fr',
            sm: 'repeat(2, 1fr)',
            md: 'repeat(3, 1fr)',
          },
          gap: 3,
        }}
      >
        {data.items.map((listing) => (
          <ListingCard key={listing.id} listing={listing} />
        ))}
      </Box>
    </>
  );
}