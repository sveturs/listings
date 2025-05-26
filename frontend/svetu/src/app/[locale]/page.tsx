'use client';

import { useEffect } from 'react';
import { useTranslations } from 'next-intl';
import Link from 'next/link';
import { useDispatch, useSelector } from 'react-redux';
import { Container, Typography, Button, Box, CircularProgress, Alert } from '@mui/material';
import { AppDispatch, RootState } from '@/store/store';
import { fetchListings } from '@/store/slices/listingsSlice';
import ListingGrid from '@/components/listings/ListingGrid';

export default function HomePage() {
  const t = useTranslations('home');
  const dispatch = useDispatch<AppDispatch>();
  
  const { listings, loading, error } = useSelector((state: RootState) => state.listings);

  useEffect(() => {
    // Fetch featured listings for homepage
    dispatch(fetchListings({
      filters: {},
      pagination: { page: 0, size: 6 }
    }));
  }, [dispatch]);
  
  return (
    <Container maxWidth="lg" sx={{ py: 4 }}>
      {/* Hero Section */}
      <Box sx={{ textAlign: 'center', py: 8 }}>
        <Typography variant="h2" component="h1" gutterBottom>
          {t('title')}
        </Typography>
        <Typography variant="h5" color="text.secondary" paragraph>
          {t('description')}
        </Typography>
        <Button
          component={Link}
          href="/marketplace/listings"
          variant="contained"
          size="large"
          sx={{ mt: 2 }}
        >
          {t('viewListings')}
        </Button>
      </Box>
      
      {/* Featured Listings Section */}
      <Box sx={{ mt: 8 }}>
        <Typography variant="h4" component="h2" gutterBottom sx={{ mb: 4 }}>
          {t('featuredListings')}
        </Typography>
        
        {loading && (
          <Box sx={{ display: 'flex', justifyContent: 'center', py: 8 }}>
            <CircularProgress />
          </Box>
        )}
        
        {error && (
          <Alert severity="error" sx={{ mb: 2 }}>
            {error}
          </Alert>
        )}
        
        {!loading && !error && (
          <>
            <ListingGrid 
              listings={listings.slice(0, 6)} 
              emptyMessage={t('noListingsAvailable')}
            />
            
            {listings.length > 0 && (
              <Box sx={{ textAlign: 'center', mt: 4 }}>
                <Button
                  component={Link}
                  href="/marketplace/listings"
                  variant="outlined"
                  size="large"
                >
                  {t('viewAllListings')}
                </Button>
              </Box>
            )}
          </>
        )}
      </Box>
    </Container>
  );
}