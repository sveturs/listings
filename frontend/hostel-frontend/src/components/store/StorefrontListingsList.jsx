// frontend/hostel-frontend/src/components/store/StorefrontListingsList.jsx
import React from 'react';
import { useTranslation } from 'react-i18next';
import {
  Box,
  Checkbox,
  Paper,
  Typography,
  Chip,
  Grid,
  useTheme,
  useMediaQuery,
  Divider,
  Avatar
} from '@mui/material';
import { Edit, Eye, Tag, Clock, MapPin, Star } from 'lucide-react';
import { Link } from 'react-router-dom';

const StorefrontListingsList = ({ listings, selectedItems, onSelectItem, onSelectAll }) => {
  const { t } = useTranslation(['marketplace', 'common']);
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('sm'));

  // Проверяем, выбраны ли все объявления
  const isAllSelected = listings.length > 0 && selectedItems.length === listings.length;
  
  // Форматирование цены
  const formatPrice = (price) => {
    return new Intl.NumberFormat('ru-RU', {
      style: 'currency',
      currency: 'RSD',
      maximumFractionDigits: 0
    }).format(price);
  };

  // Форматирование даты
  const formatDate = (dateString) => {
    try {
      const date = new Date(dateString);
      return new Intl.DateTimeFormat('ru-RU', {
        day: '2-digit',
        month: '2-digit',
        year: 'numeric'
      }).format(date);
    } catch (e) {
      return dateString;
    }
  };

  return (
    <Box>
      <Paper 
        elevation={0} 
        sx={{ 
          p: 2, 
          mb: 2, 
          display: 'flex', 
          alignItems: 'center', 
          bgcolor: 'grey.100',
          borderRadius: 1
        }}
      >
        <Checkbox
          checked={isAllSelected}
          indeterminate={selectedItems.length > 0 && selectedItems.length < listings.length}
          onChange={(e) => onSelectAll(e.target.checked)}
          sx={{ mr: 1 }}
        />
        
        <Grid container>
          <Grid item xs={isMobile ? 5 : 6} md={7}>
            <Typography variant="subtitle2" sx={{ fontWeight: 'bold' }}>
              {t('store.listings.title')}
            </Typography>
          </Grid>
          <Grid item xs={isMobile ? 3 : 2} md={2} sx={{ textAlign: isMobile ? 'center' : 'right' }}>
            <Typography variant="subtitle2" sx={{ fontWeight: 'bold' }}>
              {t('store.listings.price')}
            </Typography>
          </Grid>
          <Grid item xs={2} md={1} sx={{ textAlign: 'center' }}>
            <Typography variant="subtitle2" sx={{ fontWeight: 'bold' }}>
              {t('store.listings.status')}
            </Typography>
          </Grid>
          <Grid item xs={2} md={2} sx={{ textAlign: 'center' }}>
            <Typography variant="subtitle2" sx={{ fontWeight: 'bold' }}>
              {t('store.listings.date')}
            </Typography>
          </Grid>
        </Grid>
      </Paper>

      {listings.length === 0 ? (
        <Box py={4} textAlign="center">
          <Typography color="text.secondary">
            {t('store.listings.noListings')}
          </Typography>
        </Box>
      ) : (
        <Box>
          {listings.map((listing) => (
            <Paper 
              key={listing.id} 
              elevation={0}
              sx={{ 
                p: 2, 
                mb: 1, 
                display: 'flex', 
                alignItems: 'flex-start',
                '&:hover': { bgcolor: 'grey.50' },
                transition: 'background-color 0.2s',
              }}
            >
              <Checkbox
                checked={selectedItems.includes(listing.id)}
                onChange={() => onSelectItem(listing.id)}
                sx={{ mt: isMobile ? 0.5 : 1.5, mr: 1 }}
              />
              
              <Grid container alignItems="center" spacing={1}>
                <Grid item xs={isMobile ? 10 : 6} md={7}>
                  <Box display="flex" alignItems="center">
                    {listing.images && listing.images.length > 0 ? (
                      <Avatar 
                        variant="rounded"
                        src={`${process.env.REACT_APP_BACKEND_URL}/uploads/${listing.images[0].file_path}`}
                        alt={listing.title}
                        sx={{ width: 60, height: 60, mr: 2 }}
                      />
                    ) : (
                      <Avatar 
                        variant="rounded"
                        sx={{ width: 60, height: 60, mr: 2, bgcolor: 'grey.300' }}
                      >
                        <Tag size={24} />
                      </Avatar>
                    )}
                    
                    <Box>
                      <Link to={`/marketplace/listings/${listing.id}`} style={{ textDecoration: 'none', color: 'inherit' }}>
                        <Typography variant="subtitle1" sx={{ 
                          fontWeight: 'medium',
                          mb: 0.5,
                          display: '-webkit-box',
                          WebkitBoxOrient: 'vertical',
                          WebkitLineClamp: 2,
                          overflow: 'hidden',
                          textOverflow: 'ellipsis',
                          '&:hover': { color: theme.palette.primary.main }
                        }}>
                          {listing.title}
                        </Typography>
                      </Link>
                      
                      <Box display="flex" alignItems="center" flexWrap="wrap" gap={1}>
                        {!isMobile && (
                          <>
                            <Chip 
                              label={listing.condition === 'new' ? t('listings.condition.new') : t('listings.condition.used')} 
                              size="small" 
                              variant="outlined"
                            />
                            
                            {listing.category?.name && (
                              <Chip 
                                icon={<Tag size={14} />}
                                label={listing.category.name} 
                                size="small" 
                                variant="outlined"
                              />
                            )}
                          </>
                        )}
                        
                        {!isMobile && listing.location && (
                          <Box component="span" sx={{ display: 'flex', alignItems: 'center', color: 'text.secondary', fontSize: '0.75rem' }}>
                            <MapPin size={14} style={{ marginRight: 4 }} />
                            {listing.location}
                          </Box>
                        )}
                      </Box>
                    </Box>
                  </Box>
                </Grid>
                
                <Grid item xs={isMobile ? 12 : 2} md={2} sx={{ 
                  textAlign: isMobile ? 'left' : 'right',
                  pl: isMobile ? 9 : undefined,
                  mt: isMobile ? 1 : 0,
                }}>
                  <Typography variant="subtitle1" fontWeight="bold" color="primary.main">
                    {formatPrice(listing.price)}
                  </Typography>
                </Grid>
                
                <Grid item xs={isMobile ? 6 : 2} md={1} sx={{ 
                  textAlign: 'center',
                  mt: isMobile ? 1 : 0,
                }}>
                  <Chip 
                    label={listing.status === 'active' ? t('listings.status.active') : t('listings.status.inactive')} 
                    size="small"
                    color={listing.status === 'active' ? 'success' : 'default'} 
                    variant="outlined"
                  />
                </Grid>
                
                <Grid item xs={isMobile ? 6 : 2} md={2} sx={{ 
                  textAlign: 'center',
                  mt: isMobile ? 1 : 0,
                }}>
                  <Box display="flex" flexDirection="column" alignItems="center">
                    <Typography variant="body2" color="text.secondary" sx={{ display: 'flex', alignItems: 'center' }}>
                      <Clock size={14} style={{ marginRight: 4 }} />
                      {formatDate(listing.created_at)}
                    </Typography>
                    
                    {!isMobile && (
                      <Typography variant="body2" color="text.secondary" sx={{ display: 'flex', alignItems: 'center', mt: 0.5 }}>
                        <Eye size={14} style={{ marginRight: 4 }} />
                        {t('listings.views', { count: listing.views_count || 0 })}
                      </Typography>
                    )}
                  </Box>
                </Grid>
              </Grid>
            </Paper>
          ))}
        </Box>
      )}
    </Box>
  );
};

export default StorefrontListingsList;