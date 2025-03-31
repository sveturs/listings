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

const FavoriteListingsPage = () => {
  const [listings, setListings] = useState(null); // Изменили на null
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const { user } = useAuth();
  const { t } = useTranslation('marketplace');
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('md'));
  
  // Состояния для пагинации и бесконечной прокрутки
  const [page, setPage] = useState(1);
  const [hasMoreListings, setHasMoreListings] = useState(true);
  const [loadingMore, setLoadingMore] = useState(false);
  const [totalListings, setTotalListings] = useState(0);
  
  // Состояния для сортировки
  const [sortField, setSortField] = useState('created_at');
  const [sortDirection, setSortDirection] = useState('desc');
  
  // Состояние для режима отображения (grid/list)
  const [viewMode, setViewMode] = useState(() => {
    // Загружаем сохраненный режим просмотра из localStorage
    const savedViewMode = localStorage.getItem('favorites-view-mode');
    return savedViewMode === 'list' ? 'list' : 'grid';
  });
  
  // Функция для загрузки начальных данных
  const fetchInitialData = async (sortParams = {}) => {
    try {
      setLoading(true);
      setError(null);
      
      // Формируем параметры API сортировки
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
      
      const response = await axios.get('/api/v1/marketplace/favorites', {
        params: {
          page: 1,
          size: 20,
          sort_by: apiSort
        }
      });
      
      const data = response.data.data || [];
      setListings(data);
      
      // Обновляем информацию о пагинации
      if (response.data.meta) {
        setTotalListings(response.data.meta.total || 0);
        setHasMoreListings(data.length < response.data.meta.total);
      } else {
        setHasMoreListings(data.length >= 20);
      }
      
      // Обновляем состояние сортировки если переданы новые параметры
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
  
  // Функция для загрузки дополнительных объявлений
  const fetchMoreListings = async () => {
    if (!hasMoreListings || loadingMore) return;
    
    try {
      setLoadingMore(true);
      const nextPage = page + 1;
      
      // Формируем параметры API сортировки
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
      
      const response = await axios.get('/api/v1/marketplace/favorites', {
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
        
        // Проверяем, есть ли еще объявления для загрузки
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
  
  // Обработчик для бесконечной прокрутки
  const handleLoadMore = () => {
    fetchMoreListings();
  };
  
  // Обработчик переключения режима отображения (сетка/список)
  const handleViewModeChange = (event, newMode) => {
    if (newMode !== null) {
      setViewMode(newMode);
      // Сохраняем предпочтение пользователя в localStorage
      localStorage.setItem('favorites-view-mode', newMode);
    }
  };
  
  // Обработчик изменения сортировки
  const handleSortChange = (field, direction) => {
    console.log(`Обработчик сортировки получил: поле=${field}, направление=${direction}`);
    
    // Сбрасываем предыдущие объявления и состояние пагинации
    setListings([]);
    setPage(1);
    setHasMoreListings(true);
    
    // Запрашиваем данные с новыми параметрами сортировки
    fetchInitialData({ field, direction });
  };
  
  useEffect(() => {
    if (user) {
      fetchInitialData();
    } else {
      setListings([]);
      setLoading(false);
    }
  }, [user, t]);

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
        
        {/* Переключатель режима отображения */}
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
          autoLoad={!isMobile} // Автозагрузка только на десктопе
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
              filters={{ // Добавляем параметр фильтров для корректной работы сортировки
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