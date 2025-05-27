'use client';

import { useQuery, useInfiniteQuery } from '@tanstack/react-query';
import { useTranslations } from 'next-intl';
import { useEffect, useRef, useCallback } from 'react';
import { listingService } from '@/services/listing.service';
import ListingCard from './ListingCard';
import { 
  Alert, 
  Box, 
  Skeleton, 
  Button, 
  CircularProgress,
  List,
  ListItem 
} from '@mui/material';

interface ListingGridProps {
  viewMode?: 'grid' | 'list';
  searchQuery?: string;
  filters?: any;
}

export default function ListingGrid({ viewMode = 'grid', searchQuery = '', filters = {} }: ListingGridProps) {
  const t = useTranslations('marketplace.listings');
  const loadMoreRef = useRef<HTMLDivElement>(null);
  
  // Build query params from filters and search
  const queryParams = {
    q: searchQuery,
    ...filters,
    page: 1,
    size: 20
  };

  // Use infinite query for pagination
  const {
    data,
    error,
    fetchNextPage,
    hasNextPage,
    isFetchingNextPage,
    isLoading,
  } = useInfiniteQuery({
    queryKey: ['listings', queryParams],
    queryFn: ({ pageParam = 1 }) => 
      listingService.searchListings({
        ...queryParams,
        page: pageParam
      }),
    getNextPageParam: (lastPage, pages) => {
      if (lastPage.page < lastPage.totalPages) {
        return lastPage.page + 1;
      }
      return undefined;
    },
    initialPageParam: 1,
  });

  // Infinite scroll observer
  useEffect(() => {
    const observer = new IntersectionObserver(
      (entries) => {
        if (entries[0].isIntersecting && hasNextPage && !isFetchingNextPage) {
          fetchNextPage();
        }
      },
      { threshold: 0.1 }
    );

    const currentRef = loadMoreRef.current;
    if (currentRef) {
      observer.observe(currentRef);
    }

    return () => {
      if (currentRef) {
        observer.unobserve(currentRef);
      }
    };
  }, [hasNextPage, isFetchingNextPage, fetchNextPage]);

  if (isLoading) {
    return (
      <Box
        sx={{
          display: viewMode === 'grid' ? 'grid' : 'block',
          gridTemplateColumns: viewMode === 'grid' ? {
            xs: '1fr',
            sm: 'repeat(2, 1fr)',
            md: 'repeat(3, 1fr)',
          } : undefined,
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

  const allListings = data?.pages.flatMap(page => page.items) || [];

  if (allListings.length === 0) {
    return (
      <Box textAlign="center" py={8}>
        <Alert severity="info">{t('noListingsFound')}</Alert>
      </Box>
    );
  }

  const totalCount = data?.pages[0]?.total || 0;

  return (
    <>
      <Box mb={2}>
        <p className="text-gray-600">{t('found', { count: totalCount })}</p>
      </Box>
      
      {viewMode === 'grid' ? (
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
          {allListings.map((listing) => (
            <ListingCard key={listing.id} listing={listing} />
          ))}
        </Box>
      ) : (
        <List sx={{ width: '100%' }}>
          {allListings.map((listing) => (
            <ListItem key={listing.id} sx={{ px: 0 }}>
              <ListingCard listing={listing} viewMode="list" />
            </ListItem>
          ))}
        </List>
      )}

      {/* Load more trigger */}
      <Box ref={loadMoreRef} sx={{ mt: 4, textAlign: 'center' }}>
        {isFetchingNextPage && <CircularProgress />}
        {!hasNextPage && allListings.length > 0 && (
          <p className="text-gray-500">{t('noMoreListings')}</p>
        )}
      </Box>

      {/* Manual load more button for mobile */}
      {hasNextPage && !isFetchingNextPage && (
        <Box sx={{ mt: 4, textAlign: 'center', display: { md: 'none' } }}>
          <Button 
            variant="outlined" 
            onClick={() => fetchNextPage()}
            disabled={isFetchingNextPage}
          >
            {t('loadMore')}
          </Button>
        </Box>
      )}
    </>
  );
}