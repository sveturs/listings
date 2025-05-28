'use client';

import { useInfiniteQuery } from '@tanstack/react-query';
import { useTranslations } from 'next-intl';
import { useEffect, useRef } from 'react';
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
  filters?: Record<string, unknown>;
}

export default function ListingGrid({ viewMode = 'grid', searchQuery = '', filters = {} }: ListingGridProps) {
  const t = useTranslations('marketplace.listings');
  const loadMoreRef = useRef<HTMLDivElement>(null);
  
  // Build query params from filters and search
  const queryParams = {
    q: searchQuery,
    ...filters,
    page: 1,
    size: 20, // Changed from 25 to match backend pagination
    sort_by: 'date_desc'
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
    getNextPageParam: (lastPage, allPages) => {
      console.log('Pagination debug:', {
        currentPage: lastPage.page,
        totalPages: lastPage.totalPages,
        total: lastPage.total,
        pageSize: lastPage.pageSize,
        itemsInPage: lastPage.items?.length || 0,
        allPagesCount: allPages.length,
        totalItemsSoFar: allPages.reduce((acc, page) => acc + page.items.length, 0)
      });
      
      // Check if there are more pages
      // Handle both 0-based and 1-based pagination
      const hasMorePages = lastPage.page < lastPage.totalPages;
      const hasMoreItems = lastPage.items.length === lastPage.pageSize;
      
      if (hasMorePages || (hasMoreItems && lastPage.totalPages === 0)) {
        return lastPage.page + 1;
      }
      return undefined;
    },
    initialPageParam: 1,
  });

  // Infinite scroll observer
  useEffect(() => {
    console.log('Infinite scroll state:', {
      hasNextPage,
      isFetchingNextPage,
      totalPages: data?.pages?.length || 0,
      lastPage: data?.pages?.[data.pages.length - 1]
    });
    
    const observer = new IntersectionObserver(
      (entries) => {
        if (entries[0].isIntersecting && hasNextPage && !isFetchingNextPage) {
          console.log('Triggering fetchNextPage');
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
  const displayedCount = allListings.length;

  return (
    <>
      <Box mb={2}>
        <p className="text-gray-600">
          {totalCount > 0 
            ? t('found', { count: totalCount })
            : `${displayedCount} ${t('listingsDisplayed')}`
          }
        </p>
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
          <p className="text-gray-500">
            {t('noMoreListings')} ({t('showingXofY', { displayed: displayedCount, total: totalCount })})
          </p>
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