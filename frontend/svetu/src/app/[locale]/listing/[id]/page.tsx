'use client';

import { useParams } from 'next/navigation';
import { useQuery } from '@tanstack/react-query';
import { useTranslations } from 'next-intl';
import Link from 'next/link';
import { 
  Box, 
  Container, 
  Typography, 
  Paper,
  Avatar,
  Button,
  Skeleton
} from '@mui/material';
import { ArrowBack, LocationOn, Person } from '@mui/icons-material';
import { listingService } from '@/services/listing.service';

export default function ListingPage() {
  const params = useParams();
  const id = params.id as string;
  const t = useTranslations('listing');
  
  const { data: listing, isLoading, error } = useQuery({
    queryKey: ['listing', id],
    queryFn: () => listingService.getListingById(id),
  });

  if (isLoading) {
    return (
      <Container maxWidth="lg">
        <Skeleton variant="text" width={100} height={40} />
        <Skeleton variant="text" width={300} height={60} sx={{ mb: 2 }} />
        <Skeleton variant="rectangular" height={400} sx={{ mb: 3 }} />
        <Box display="flex" gap={3} flexDirection={{ xs: 'column', md: 'row' }}>
          <Box flex={{ xs: 1, md: 8 }}>
            <Skeleton variant="rectangular" height={200} />
          </Box>
          <Box flex={{ xs: 1, md: 4 }}>
            <Skeleton variant="rectangular" height={200} />
          </Box>
        </Box>
      </Container>
    );
  }

  if (error || !listing) {
    return (
      <Container maxWidth="lg">
        <Typography variant="h6" color="error">
          {t('notFound')}
        </Typography>
      </Container>
    );
  }
  
  return (
    <Container maxWidth="lg">
      <Link 
        href={`/${params.locale}/marketplace`}
        style={{ textDecoration: 'none', display: 'inline-flex', alignItems: 'center', marginBottom: '1rem' }}
      >
        <ArrowBack sx={{ mr: 1 }} />
        <Typography variant="body1" color="primary">
          {t('backToHome')}
        </Typography>
      </Link>
      
      <Typography variant="h4" component="h1" gutterBottom>
        {listing.title}
      </Typography>
      
      <Box sx={{ mb: 3, display: 'flex', alignItems: 'center', gap: 1 }}>
        <LocationOn color="action" />
        <Typography variant="body1" color="text.secondary">
          {listing.location || 'Location not specified'}
        </Typography>
      </Box>
      
      <Box sx={{ position: 'relative', height: 400, bgcolor: 'grey.200', borderRadius: 2, mb: 4 }}>
        {listing.images?.[0] ? (
          <>
            {/* eslint-disable-next-line @next/next/no-img-element */}
            <img 
              src={listing.images[0].url} 
              alt={listing.title}
              style={{ width: '100%', height: '100%', objectFit: 'cover', borderRadius: 8 }}
            />
          </>
        ) : (
          <Box display="flex" alignItems="center" justifyContent="center" height="100%">
            <Typography color="text.secondary">No image available</Typography>
          </Box>
        )}
      </Box>
      
      <Box display="flex" gap={3} flexDirection={{ xs: 'column', md: 'row' }}>
        <Box flex={{ xs: 1, md: 8 }}>
          <Paper sx={{ p: 3, mb: 3 }}>
            <Typography variant="h6" gutterBottom>
              {t('description')}
            </Typography>
            <Typography variant="body1" color="text.secondary">
              {listing.description}
            </Typography>
          </Paper>
        </Box>
        
        <Box flex={{ xs: 1, md: 4 }}>
          <Paper sx={{ p: 3, mb: 3 }}>
            <Typography variant="h6" gutterBottom>
              {t('price')}
            </Typography>
            <Typography variant="h4" color="primary">
              {listing.currency === 'EUR' ? '€' : '$'}{listing.price.toLocaleString()}
            </Typography>
          </Paper>
          
          {listing.user && (
            <Paper sx={{ p: 3 }}>
              <Typography variant="h6" gutterBottom>
                Seller Information
              </Typography>
              <Box display="flex" alignItems="center" gap={2}>
                <Avatar>
                  <Person />
                </Avatar>
                <Box>
                  <Typography variant="body1">
                    {listing.user.name}
                  </Typography>
                  {listing.user.rating && (
                    <Typography variant="body2" color="text.secondary">
                      Rating: {listing.user.rating}/5
                    </Typography>
                  )}
                </Box>
              </Box>
              <Button 
                variant="contained" 
                fullWidth 
                sx={{ mt: 2 }}
              >
                Contact Seller
              </Button>
            </Paper>
          )}
        </Box>
      </Box>
    </Container>
  );
}