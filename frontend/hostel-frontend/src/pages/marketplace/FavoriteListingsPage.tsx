// frontend/hostel-frontend/src/pages/marketplace/FavoriteListingsPage.tsx
import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import {
  Container,
  Typography,
  Grid,
  Box,
  CircularProgress,
  Alert,
  ToggleButtonGroup,
  ToggleButton,
  useTheme,
  useMediaQuery
} from '@mui/material';
import { Grid as GridIcon, List } from 'lucide-react';
import ListingCard from '../../components/marketplace/ListingCard';
import MarketplaceListingsList from '../../components/marketplace/MarketplaceListingsList';
import InfiniteScroll from '../../components/marketplace/InfiniteScroll';
import axios from '../../api/axios';
import { useAuth } from '../../contexts/AuthContext';
import { Listing } from '../../components/marketplace/ListingCard';

// Define the view mode type
type ViewMode = 'grid' | 'list';

// Define the sort field type
type SortField = 'created_at' | 'title' | 'price' | 'reviews';

// Define the sort direction type
type SortDirection = 'asc' | 'desc';

// Define the API response interface
interface ListingsResponse {
  data: Listing[];
  meta?: {
    total: number;
    page: number;
    size: number;
  };
}

// Define the filter state interface
interface FilterState {
  sort_by?: string;
}

const FavoriteListingsPage: React.FC = () => {
  const [listings, setListings] = useState<Listing[] | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const { user } = useAuth();
  const { t } = useTranslation('marketplace');
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('md'));
  
  // Pagination and infinite scroll states
  const [page, setPage] = useState<number>(1);
  const [hasMoreListings, setHasMoreListings] = useState<boolean>(true);
  const [loadingMore, setLoadingMore] = useState<boolean>(false);
  const [totalListings, setTotalListings] = useState<number>(0);
  
  // Sorting states
  const [sortField, setSortField] = useState<SortField>('created_at');
  const [sortDirection, setSortDirection] = useState<SortDirection>('desc');
  
  // View mode state (grid/list)
  const [viewMode, setViewMode] = useState<ViewMode>(() => {
    // Load saved view mode from localStorage
    const savedViewMode = localStorage.getItem('favorites-view-mode');
    return (savedViewMode === 'list' ? 'list' : 'grid') as ViewMode;
  });
  
  // Function to fetch initial data
  const fetchInitialData = async (sortParams: { field?: SortField, direction?: SortDirection } = {}): Promise<void> => {
    try {
      setLoading(true);
      setError(null);
      
      // Format API sort parameters
      let apiSort = '';
      const currentSortField = sortParams.field || sortField;
      const currentSortDirection = sortParams.direction || sortDirection;
      
      switch (currentSortField) {
        case 'created_at':
          apiSort = `date_${currentSortDirection}`;
          break;
        case 'title':
          apiSort = `title_${currentSortDirection}`;
          break;
        case 'price':
          apiSort = `price_${currentSortDirection}`;
          break;
        case 'reviews':
          apiSort = `rating_${currentSortDirection}`;
          break;
        default:
          apiSort = `date_${currentSortDirection}`;
      }
      
      console.log(`Запрос избранного с сортировкой: ${apiSort}`);
      
      const response = await axios.get<ListingsResponse>('/api/v1/marketplace/favorites', {
        params: {
          page: 1,
          size: 20,
          sort_by: apiSort
        }
      });
      
      const data = response.data.data || [];
      setListings(data);
      
      // Update pagination information
      if (response.data.meta) {
        setTotalListings(response.data.meta.total || 0);
        setHasMoreListings(data.length < response.data.meta.total);
      } else {
        setHasMoreListings(data.length >= 20);
      }
      
      // Update sort state if new parameters are provided
      if (sortParams.field) setSortField(sortParams.field);
      if (sortParams.direction) setSortDirection(sortParams.direction);
      
      setPage(1);
    } catch (err) {
      console.error('Error fetching favorites:', err);
      setError(t('favorites.errors.loadFailed'));
      setListings([]);
    } finally {
      setLoading(false);
    }
  };
  
  // Function to fetch additional listings
  const fetchMoreListings = async (): Promise<void> => {
    if (!hasMoreListings || loadingMore) return;
    
    try {
      setLoadingMore(true);
      const nextPage = page + 1;
      
      // Format API sort parameters
      let apiSort = '';
      switch (sortField) {
        case 'created_at':
          apiSort = `date_${sortDirection}`;
          break;
        case 'title':
          apiSort = `title_${sortDirection}`;
          break;
        case 'price':
          apiSort = `price_${sortDirection}`;
          break;
        case 'reviews':
          apiSort = `rating_${sortDirection}`;
          break;
        default:
          apiSort = `date_${sortDirection}`;
      }
      
      console.log(`Загрузка дополнительных избранных с сортировкой: ${apiSort}`);
      
      const response = await axios.get<ListingsResponse>('/api/v1/marketplace/favorites', {
        params: {
          page: nextPage,
          size: 20,
          sort_by: apiSort
        }
      });
      
      const newListings = response.data.data || [];
      
      if (newListings.length > 0) {
        setListings(prev => [...(prev || []), ...newListings]);
        setPage(nextPage);
        
        // Check if there are more listings to load
        if (response.data.meta) {
          const total = response.data.meta.total || 0;
          setHasMoreListings((listings?.length || 0) + newListings.length < total);
        } else {
          setHasMoreListings(newListings.length >= 20);
        }
      } else {
        setHasMoreListings(false);
      }
    } catch (err) {
      console.error('Error fetching more favorites:', err);
    } finally {
      setLoadingMore(false);
    }
  };
  
  // Handler for infinite scroll
  const handleLoadMore = (): void => {
    fetchMoreListings();
  };
  
  // Handler for view mode toggle (grid/list)
  const handleViewModeChange = (event: React.MouseEvent<HTMLElement>, newMode: ViewMode | null): void => {
    if (newMode !== null) {
      setViewMode(newMode);
      // Save user preference to localStorage
      localStorage.setItem('favorites-view-mode', newMode);
    }
  };
  
  // Handler for sort change
  const handleSortChange = (field: SortField, direction: SortDirection): void => {
    console.log(`Обработчик сортировки получил: поле=${field}, направление=${direction}`);
    
    // Reset previous listings and pagination state
    setListings([]);
    setPage(1);
    setHasMoreListings(true);
    
    // Request data with new sort parameters
    fetchInitialData({ field, direction });
  };
  
  useEffect(() => {
    if (user) {
      fetchInitialData();
    } else {
      setListings([]);
      setLoading(false);
    }
  }, [user]);

  if (loading) {
    return (
      <Box display="flex" justifyContent="center" p={4}>
        <CircularProgress />
      </Box>
    );
  }

  if (!user) {
    return (
      <Container sx={{ py: 4 }}>
        <Alert severity="info">
          {t('favorites.authRequired')}
        </Alert>
      </Container>
    );
  }

  return (
    <Container sx={{ py: 4 }}>
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
        <Typography variant="h4">
          {t('favorites.title')}
        </Typography>
        
        {/* View mode toggle */}
        <ToggleButtonGroup
          value={viewMode}
          exclusive
          onChange={handleViewModeChange}
          aria-label="view mode"
          size="small"
        >
          <ToggleButton value="grid" aria-label="grid view">
            <GridIcon size={18} />
          </ToggleButton>
          <ToggleButton value="list" aria-label="list view">
            <List size={18} />
          </ToggleButton>
        </ToggleButtonGroup>
      </Box>

      {error && (
        <Alert severity="error" sx={{ mb: 2 }}>
          {error}
        </Alert>
      )}

      {listings && listings.length === 0 ? (
        <Alert severity="info">
          {t('favorites.empty')}
        </Alert>
      ) : (
        <InfiniteScroll
          hasMore={hasMoreListings}
          loading={loadingMore}
          onLoadMore={handleLoadMore}
          autoLoad={!isMobile} // Auto-load only on desktop
          loadingMessage={t('marketplace:listings.loading', { defaultValue: 'Загрузка...' })}
          loadMoreButtonText={t('marketplace:listings.loadMore', { defaultValue: 'Показать ещё' })}
          noMoreItemsText={t('favorites.noMoreItems', { defaultValue: 'Больше нет избранных товаров' })}
        >
          {viewMode === 'grid' ? (
            <Grid container spacing={3}>
              {listings && listings.map((listing) => (
                <Grid item xs={12} sm={6} md={4} key={listing.id}>
                  <Link
                    to={`/marketplace/listings/${listing.id}`}
                    style={{ textDecoration: 'none' }}
                  >
                    <ListingCard listing={listing} />
                  </Link>
                </Grid>
              ))}
            </Grid>
          ) : (
            <MarketplaceListingsList
              listings={listings || []}
              showSelection={false}
              initialSortField={sortField}
              initialSortOrder={sortDirection}
              onSortChange={handleSortChange}
              filters={{ // Add filter parameters for proper sorting
                sort_by: `${sortField === 'created_at' ? 'date' : sortField}_${sortDirection}`
              }}
            />
          )}
        </InfiniteScroll>
      )}
    </Container>
  );
};

export default FavoriteListingsPage;