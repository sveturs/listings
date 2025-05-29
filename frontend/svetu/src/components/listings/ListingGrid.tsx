'use client';

import React from 'react';
import { Box, CircularProgress, Typography, Alert } from '@mui/material';
import { useTranslations } from 'next-intl';
import ListingCard from './ListingCard';
import { Listing } from '@/types/listing';

interface ListingGridProps {
  listings: Listing[];
  loading?: boolean;
  error?: string | null;
  emptyMessage?: string;
}

const ListingGrid: React.FC<ListingGridProps> = ({
  listings,
  loading = false,
  error = null,
  emptyMessage,
}) => {
  const t = useTranslations('marketplace');

  if (loading) {
    return (
      <Box sx={{ display: 'flex', justifyContent: 'center', py: 8 }}>
        <CircularProgress />
      </Box>
    );
  }

  if (error) {
    return (
      <Alert severity="error" sx={{ mb: 2 }}>
        {error}
      </Alert>
    );
  }

  if (!listings || listings.length === 0) {
    return (
      <Box sx={{ textAlign: 'center', py: 8 }}>
        <Typography variant="h6" color="text.secondary">
          {emptyMessage || t('listings.noListingsFound')}
        </Typography>
      </Box>
    );
  }

  return (
    <Box
      sx={{
        display: 'grid',
        gridTemplateColumns: {
          xs: '1fr',
          sm: 'repeat(2, 1fr)',
          md: 'repeat(3, 1fr)',
          lg: 'repeat(4, 1fr)',
        },
        gap: 3,
      }}
    >
      {listings.map((listing) => (
        <ListingCard key={listing.id} listing={listing} />
      ))}
    </Box>
  );
};

export default ListingGrid;