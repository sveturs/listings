import React, { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { useNavigate } from 'react-router-dom';
import {
  Box,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  Checkbox,
  Typography,
  IconButton,
  Chip,
  useTheme,
  useMediaQuery,
  TableSortLabel
} from '@mui/material';
import { Edit, Trash2, Eye, MapPin, Calendar, Percent, ArrowUpDown } from 'lucide-react';

// Вспомогательные функции для форматирования
const formatPrice = (price) => {
  return new Intl.NumberFormat('sr-RS', {
    style: 'currency',
    currency: 'RSD',
    maximumFractionDigits: 0
  }).format(price || 0);
};

const formatDate = (dateString) => {
  if (!dateString) return '';
  return new Date(dateString).toLocaleDateString();
};

const getImageUrl = (listing) => {
  if (!listing || !listing.images || listing.images.length === 0) {
    return '/placeholder.jpg';
  }
  
  const baseUrl = process.env.REACT_APP_BACKEND_URL || '';
  
  // Найдем главное изображение или используем первое
  const mainImage = listing.images.find(img => img.is_main) || listing.images[0];
  
  if (typeof mainImage === 'string') {
    return `${baseUrl}/uploads/${mainImage}`;
  }
  
  if (mainImage && mainImage.file_path) {
    return `${baseUrl}/uploads/${mainImage.file_path}`;
  }
  
  return '/placeholder.jpg';
};

const MarketplaceListingsList = ({ 
  listings, 
  selectedItems = [], 
  onSelectItem = null, 
  onSelectAll = null,
  showSelection = false,
  onSortChange = null, // пропс для передачи информации о сортировке родительскому компоненту
  filters = {},
  initialSortField = 'created_at',
  initialSortOrder = 'desc' 
}) => {
  const { t } = useTranslation('marketplace');
  const navigate = useNavigate();
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('md'));
  const [order, setOrder] = useState(initialSortOrder);
  const [orderBy, setOrderBy] = useState(initialSortField);
  
 
  // Синхронизация состояния сортировки с фильтрами
  // Извлекаем текущее направление сортировки из значения sort_by
useEffect(() => {
    // Если есть sort_by в формате field_direction
    if (filters && filters.sort_by) {
        const sortParts = filters.sort_by.split('_');
        if (sortParts.length >= 2) {
            // Получаем поле и направление
            let field = sortParts[0];
            let direction = sortParts.pop(); // Последний элемент - направление
            
            // Преобразуем поле API обратно в поле UI
            switch (field) {
                case 'date':
                    setOrderBy('created_at');
                    break;
                case 'price':
                case 'title':
                case 'location':
                    setOrderBy(field);
                    break;
                default:
                    setOrderBy('created_at');
            }
            
            if (direction === 'asc' || direction === 'desc') {
                setOrder(direction);
            }
        }
    }
}, [filters]);
  
  // Адаптация для мобильных: скрываем некоторые колонки
  const columns = isMobile 
    ? ['image', 'title', 'price'] 
    : ['image', 'title', 'price', 'location', 'date'];

  // Проверяем, выбраны ли все элементы
  const isAllSelected = listings.length > 0 && selectedItems.length === listings.length;
  
  // Обработчик клика по строке - переход к объявлению
  const handleRowClick = (id) => {
    navigate(`/marketplace/listings/${id}`);
  };
  
  // Обработчик изменения сортировки
  const handleRequestSort = (property) => {
    const isAsc = orderBy === property && order === 'asc';
    const newOrder = isAsc ? 'desc' : 'asc';
    setOrder(newOrder);
    setOrderBy(property);
    
    // Если предоставлен колбэк для внешней сортировки, вызываем его
    if (onSortChange) {
      onSortChange(property, newOrder);
    }
  };
  
  // Функция для создания props для заголовка с сортировкой
  const createSortHandler = (property) => () => {
    handleRequestSort(property);
  };
  
  // Получаем информацию о скидке
  const getDiscountInfo = (listing) => {
    if (listing.metadata && listing.metadata.discount) {
      return {
        percent: listing.metadata.discount.discount_percent,
        oldPrice: listing.metadata.discount.previous_price
      };
    }
    if (listing.has_discount && listing.old_price) {
      const percent = Math.round((1 - listing.price / Number(listing.old_price)) * 100);
      return {
        percent: percent,
        oldPrice: listing.old_price
      };
    }
    return null;
  };

  return (
    <TableContainer component={Paper} elevation={0} variant="outlined">
      <Table sx={{ minWidth: 650 }}>
        <TableHead>
          <TableRow>
            {showSelection && onSelectAll && (
              <TableCell padding="checkbox">
                <Checkbox
                  indeterminate={selectedItems.length > 0 && selectedItems.length < listings.length}
                  checked={isAllSelected}
                  onChange={(e) => onSelectAll(e.target.checked)}
                />
              </TableCell>
            )}
            
            {columns.includes('image') && (
              <TableCell width={80}></TableCell>  
            )}
            
            {columns.includes('title') && (
              <TableCell>
                <TableSortLabel
                  active={orderBy === 'title'}
                  direction={orderBy === 'title' ? order : 'asc'}
                  onClick={createSortHandler('title')}
                >
                  {t('listings.table.title')}
                </TableSortLabel>
              </TableCell>
            )}
            
            {columns.includes('price') && (
              <TableCell align="right">
                <TableSortLabel
                  active={orderBy === 'price'}
                  direction={orderBy === 'price' ? order : 'asc'}
                  onClick={createSortHandler('price')}
                >
                  {t('listings.table.price')}
                </TableSortLabel>
              </TableCell>
            )}
            
            {columns.includes('location') && (
              <TableCell>
                <TableSortLabel
                  active={orderBy === 'location'}
                  direction={orderBy === 'location' ? order : 'asc'}
                  onClick={createSortHandler('location')}
                >
                  {t('listings.table.location')}
                </TableSortLabel>
              </TableCell>
            )}
            
            {columns.includes('date') && (
              <TableCell>
                <TableSortLabel
                  active={orderBy === 'created_at'}
                  direction={orderBy === 'created_at' ? order : 'asc'}
                  onClick={createSortHandler('created_at')}
                >
                  {t('listings.table.date')}
                </TableSortLabel>
              </TableCell>
            )}
          </TableRow>
        </TableHead>
        
        <TableBody>
          {listings.map((listing) => {
            const isSelected = selectedItems.includes(listing.id);
            const discount = getDiscountInfo(listing);
            
            return (
              <TableRow
                key={listing.id}
                hover
                onClick={() => handleRowClick(listing.id)}
                selected={isSelected}
                sx={{ 
                  cursor: 'pointer',
                  '&:last-child td, &:last-child th': { border: 0 }
                }}
              >
                {showSelection && onSelectItem && (
                  <TableCell padding="checkbox" onClick={(e) => {
                    e.stopPropagation();
                    onSelectItem(listing.id);
                  }}>
                    <Checkbox checked={isSelected} />
                  </TableCell>
                )}
                
                {columns.includes('image') && (
                  <TableCell width={80}>
                    <Box
                      component="img"
                      src={getImageUrl(listing)}
                      alt={listing.title}
                      sx={{
                        width: 60,
                        height: 60,
                        objectFit: 'contain',
                        backgroundColor: '#f5f5f5',
                        borderRadius: 1
                      }}
                    />
                  </TableCell>
                )}
                
                {columns.includes('title') && (
                  <TableCell>
                    <Typography variant="subtitle2">{listing.title}</Typography>
                    {discount && (
                      <Chip
                        icon={<Percent size={12} />}
                        label={`-${discount.percent}%`}
                        color="warning"
                        size="small"
                        sx={{ mt: 0.5, height: 20, fontSize: '0.7rem' }}
                      />
                    )}
                  </TableCell>
                )}
                
                {columns.includes('price') && (
                  <TableCell align="right">
                    <Typography variant="subtitle2" color="primary.main">
                      {formatPrice(listing.price)}
                    </Typography>
                    {discount && (
                      <Typography
                        variant="caption"
                        color="text.secondary"
                        sx={{ textDecoration: 'line-through', display: 'block' }}
                      >
                        {formatPrice(discount.oldPrice)}
                      </Typography>
                    )}
                  </TableCell>
                )}
                
                {columns.includes('location') && (
                  <TableCell>
                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
                      <MapPin size={16} />
                      <Typography variant="body2">
                        {listing.city || listing.location || t('listings.table.noLocation')}
                      </Typography>
                    </Box>
                  </TableCell>
                )}
                
                {columns.includes('date') && (
                  <TableCell>
                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
                      <Calendar size={16} />
                      <Typography variant="body2">
                        {formatDate(listing.created_at)}
                      </Typography>
                    </Box>
                  </TableCell>
                )}
              </TableRow>
            );
          })}
          
          {listings.length === 0 && (
            <TableRow>
              <TableCell colSpan={showSelection ? columns.length + 1 : columns.length} align="center">
                <Typography variant="body2" color="text.secondary" sx={{ py: 3 }}>
                  {t('listings.table.noListings')}
                </Typography>
              </TableCell>
            </TableRow>
          )}
        </TableBody>
      </Table>
    </TableContainer>
  );
};

export default MarketplaceListingsList;