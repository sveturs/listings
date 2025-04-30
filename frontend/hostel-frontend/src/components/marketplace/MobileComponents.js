// frontend/hostel-frontend/src/components/marketplace/MobileComponents.js
import React, { useState, useCallback, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import AutocompleteInput from '../shared/AutocompleteInput';
import { Link } from 'react-router-dom';
import {
  Box, Button, IconButton, Typography, InputBase, Toolbar, TextField, Select, MenuItem,
  Paper, Grid, Drawer, Stack, List, ListItem, ListItemText, ListItemAvatar, Avatar, Divider,
  ToggleButton, ToggleButtonGroup
} from '@mui/material';
import {
  Search as SearchIcon,
  Filter,
  X,
  Check,
  ArrowLeft,
  ChevronRight,
  Plus,
  Store,
  List as ListIcon,
  Grid as GridIcon,
  MapPin,
  Calendar,
  Percent
} from 'lucide-react';
import { debounce } from 'lodash';

// Компонент MobileHeader
export const MobileHeader = ({ onOpenFilters, filtersCount, onSearch, searchValue, viewMode, onViewModeChange }) => {
  const { t } = useTranslation('marketplace', 'common');

  return (
    <Box sx={{
      borderBottom: 1,
      borderColor: 'divider',
      position: 'sticky',
      top: 0,
      zIndex: 1100,
      bgcolor: 'background.paper'
    }}>
      <Toolbar sx={{
        minHeight: '56px !important',
        px: 2,
        display: 'flex',
        justifyContent: 'space-between',
        gap: 2
      }}>
        <Box sx={{ display: 'flex', gap: 1, alignItems: 'center' }}>
          <IconButton
            onClick={onOpenFilters}
            sx={{
              position: 'relative',
              bgcolor: filtersCount > 0 ? 'action.selected' : 'transparent'
            }}
          >
            <Filter size={20} />
            {filtersCount > 0 && (
              <Box
                sx={{
                  position: 'absolute',
                  top: 4,
                  right: 4,
                  width: 16,
                  height: 16,
                  borderRadius: '50%',
                  bgcolor: 'primary.main',
                  color: 'primary.contrastText',
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                  fontSize: '0.75rem'
                }}
              >
                {filtersCount}
              </Box>
            )}
          </IconButton>

          {/* Переключатель режимов отображения */}
          <ToggleButtonGroup
            value={viewMode}
            exclusive
            onChange={onViewModeChange}
            size="small"
            aria-label="view mode"
          >
            <ToggleButton value="grid" aria-label="grid view">
              <GridIcon size={18} />
            </ToggleButton>
            <ToggleButton value="list" aria-label="list view">
              <ListIcon size={18} />
            </ToggleButton>
          </ToggleButtonGroup>
        </Box>

        <Button
          component={Link}
          to="/marketplace/create"
          variant="contained"
          size="small"
          startIcon={<Plus size={16} />}
          sx={{
            textTransform: 'none',
            fontWeight: 500,
            height: 32
          }}
        >
          {t('buttons.create', { ns: 'common' })}
        </Button>
      </Toolbar>
      <Box sx={{ px: 2, pb: 2 }}>
        <AutocompleteInput
          value={searchValue}
          onChange={onSearch}
          onSearch={onSearch}
          placeholder={t('buttons.search', { ns: 'common' })}
        />
      </Box>
    </Box>
  );
};

// Общая функция форматирования цены
const formatPrice = (price) => {
  return new Intl.NumberFormat('sr-RS', {
    style: 'currency',
    currency: 'RSD',
    maximumFractionDigits: 0
  }).format(price || 0);
};

// Получаем URL изображения
const getImageUrl = (listing) => {
  if (!listing || !listing.images || !listing.images.length) {
    return '/placeholder.jpg';
  }

  const baseUrl = process.env.REACT_APP_BACKEND_URL || '';

  // Находим главное изображение или используем первое в списке
  const mainImage = listing.images.find(img => img.is_main) || listing.images[0];

  if (mainImage && typeof mainImage === 'object') {
    // Если есть публичный URL, используем его напрямую
    if (mainImage.public_url && mainImage.public_url !== '') {
      // Проверяем, абсолютный или относительный URL
      if (mainImage.public_url.startsWith('http')) {
        return mainImage.public_url;
      } else {
        return `${baseUrl}${mainImage.public_url}`;
      }
    }

    // Для MinIO-объектов формируем URL на основе storage_type
    if (mainImage.storage_type === 'minio' ||
        (mainImage.file_path && mainImage.file_path.includes('listings/'))) {
      return `${baseUrl}${mainImage.public_url}`;
    }

    // Обычный файл
    if (mainImage.file_path) {
      return `${baseUrl}/uploads/${mainImage.file_path}`;
    }
  }

  // Для строк (обратная совместимость)
  if (mainImage && typeof mainImage === 'string') {
    return `${baseUrl}/uploads/${mainImage}`;
  }

  return '/placeholder.jpg';
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

// Форматируем дату
const formatDate = (dateString) => {
  if (!dateString) return '';
  return new Date(dateString).toLocaleDateString();
};

// Компонент MobileListingCard
export const MobileListingCard = ({ listing, viewMode = 'grid' }) => {
  const { t, i18n } = useTranslation('marketplace');

  const getTranslatedText = (field) => {
    if (!listing) return '';
    if (i18n.language === listing.original_language) {
      return listing[field];
    }
    return listing.translations?.[i18n.language]?.[field] || listing[field];
  };

  const handleShopButtonClick = (e) => {
    e.preventDefault();
    e.stopPropagation();
    // Останавливаем всплытие события для всех родительских элементов
    if (e.nativeEvent) {
      e.nativeEvent.stopImmediatePropagation();
    }
    // Используем прямой переход
    window.location.href = `/shop/${listing.storefront_id}`;
    // Предотвращаем дальнейшую обработку события
    return false;
  };

  // Получаем информацию о скидке
  const discount = getDiscountInfo(listing);

  // Отображение в режиме списка
  if (viewMode === 'list') {
    return (
      <Paper
        elevation={0}
        variant="outlined"
        sx={{
          mb: 1,
          borderRadius: 1,
          overflow: 'hidden'
        }}
      >
        <Box sx={{ display: 'flex', p: 1 }}>
          {/* Изображение */}
          <Box
            sx={{
              width: 80,
              height: 80,
              borderRadius: 1,
              overflow: 'hidden',
              flexShrink: 0,
              bgcolor: 'grey.100',
              position: 'relative'
            }}
          >
            <img
              src={getImageUrl(listing)}
              alt={getTranslatedText('title')}
              style={{
                width: '100%',
                height: '100%',
                objectFit: 'contain'
              }}
            />

            {/* Бейдж скидки */}
            {discount && (
              <Box
                sx={{
                  position: 'absolute',
                  top: 4,
                  left: 4,
                  bgcolor: 'warning.main',
                  color: 'warning.contrastText',
                  borderRadius: '4px',
                  px: 0.5,
                  py: 0.25,
                  display: 'flex',
                  alignItems: 'center',
                  gap: 0.25,
                  fontSize: '0.625rem',
                  fontWeight: 'bold'
                }}
              >
                <Percent size={10} />
                {`-${discount.percent}%`}
              </Box>
            )}
          </Box>

          {/* Информация о листинге */}
          <Box sx={{ ml: 1.5, flex: 1, minWidth: 0 }}>
            <Typography
              variant="subtitle2"
              sx={{
                fontSize: '0.875rem',
                fontWeight: 500,
                overflow: 'hidden',
                textOverflow: 'ellipsis',
                whiteSpace: 'nowrap'
              }}
            >
              {getTranslatedText('title')}
            </Typography>

            <Box sx={{ mt: 0.5, display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start' }}>
              <Typography
                variant="subtitle1"
                sx={{
                  fontSize: '1rem',
                  fontWeight: 600,
                  color: 'primary.main'
                }}
              >
                {formatPrice(listing.price)}
              </Typography>

              {/* Бейдж магазина */}
              {listing.storefront_id && (
                <Box
                  sx={{
                    bgcolor: 'primary.main',
                    color: 'white',
                    borderRadius: '4px',
                    px: 0.5,
                    py: 0.25,
                    display: 'flex',
                    alignItems: 'center',
                    gap: 0.25,
                    fontSize: '0.625rem',
                    fontWeight: 'bold',
                    cursor: 'pointer',
                    zIndex: 100
                  }}
                  onClick={handleShopButtonClick}
                  data-shop-button="true"
                >
                  <Store size={10} />
                  {t('listings.details.goToStore')}
                </Box>
              )}
            </Box>

            {/* Старая цена со скидкой */}
            {discount && (
              <Typography
                variant="caption"
                color="text.secondary"
                sx={{ textDecoration: 'line-through', display: 'block', mt: 0.25 }}
              >
                {formatPrice(discount.oldPrice)}
              </Typography>
            )}

            {/* Местоположение и дата */}
            <Box sx={{ display: 'flex', alignItems: 'center', mt: 0.5 }}>
              {listing.city && (
                <Box sx={{ display: 'flex', alignItems: 'center', mr: 1 }}>
                  <MapPin size={12} color="#666" />
                  <Typography variant="caption" color="text.secondary" sx={{ ml: 0.25 }}>
                    {listing.city}
                  </Typography>
                </Box>
              )}

              <Box sx={{ display: 'flex', alignItems: 'center' }}>
                <Calendar size={12} color="#666" />
                <Typography variant="caption" color="text.secondary" sx={{ ml: 0.25 }}>
                  {formatDate(listing.created_at)}
                </Typography>
              </Box>
            </Box>
          </Box>
        </Box>
      </Paper>
    );
  }

  // Отображение в режиме сетки (стандартный вариант)
  return (
    <Box sx={{ p: 1 }}>
      <Box
        sx={{
          position: 'relative',
          paddingTop: '100%',
          borderRadius: 1,
          overflow: 'hidden',
          bgcolor: 'grey.100'
        }}
      >
        {/* бейдж магазина */}
        {listing.storefront_id && (
          <Box
            sx={{
              position: 'absolute',
              top: 10,
              right: 10,
              zIndex: 100, // Очень высокий z-index
              bgcolor: 'primary.main', // MUI тема для цвета
              color: 'white',
              borderRadius: '4px',
              px: 1,
              py: 0.5,
              display: 'flex',
              alignItems: 'center',
              gap: 0.5,
              fontSize: '0.75rem',
              fontWeight: 'bold',
              cursor: 'pointer',
              isolation: 'isolate' // Изоляция событий
            }}
            // Используем onMouseDown для более раннего перехвата события, чем onClick
            onMouseDown={(e) => {
              e.preventDefault();
              e.stopPropagation();
              if (e.nativeEvent) {
                e.nativeEvent.stopImmediatePropagation();
              }
              window.location.href = `/shop/${listing.storefront_id}`;
              return false;
            }}
            // Дублируем на onClick для надежности
            onClick={(e) => {
              e.preventDefault();
              e.stopPropagation();
              if (e.nativeEvent) {
                e.nativeEvent.stopImmediatePropagation();
              }
              window.location.href = `/shop/${listing.storefront_id}`;
              return false;
            }}
            // Дублируем на onTouchStart для надежности на мобильных
            onTouchStart={(e) => {
              e.preventDefault();
              e.stopPropagation();
              if (e.nativeEvent) {
                e.nativeEvent.stopImmediatePropagation();
              }
              window.location.href = `/shop/${listing.storefront_id}`;
              return false;
            }}
            data-shop-button="true"
          >
            <Store size={14} />
            {t('listings.details.goToStore')}
          </Box>
        )}

        {/* Бейдж скидки */}
        {discount && (
          <Box
            sx={{
              position: 'absolute',
              top: 10,
              left: 10,
              zIndex: 90,
              bgcolor: 'warning.main',
              color: 'warning.contrastText',
              borderRadius: '4px',
              px: 1,
              py: 0.5,
              display: 'flex',
              alignItems: 'center',
              gap: 0.5,
              fontSize: '0.75rem',
              fontWeight: 'bold'
            }}
          >
            <Percent size={14} />
            {`-${discount.percent}%`}
          </Box>
        )}

        {/* Изображение */}
        <img
          src={getImageUrl(listing)}
          alt={getTranslatedText('title')}
          style={{
            position: 'absolute',
            top: 0,
            left: 0,
            width: '100%',
            height: '100%',
            objectFit: 'cover'
          }}
        />
      </Box>

      <Box sx={{ mt: 1 }}>
        <Typography
          variant="subtitle2"
          sx={{
            fontSize: '0.875rem',
            fontWeight: 500,
            overflow: 'hidden',
            textOverflow: 'ellipsis',
            whiteSpace: 'nowrap'
          }}
        >
          {getTranslatedText('title')}
        </Typography>

        <Typography
          variant="subtitle1"
          sx={{
            fontSize: '1rem',
            fontWeight: 600,
            color: 'primary.main',
            mt: 0.5
          }}
        >
          {formatPrice(listing.price)}
        </Typography>

        {/* Старая цена со скидкой */}
        {discount && (
          <Typography
            variant="caption"
            color="text.secondary"
            sx={{ textDecoration: 'line-through', display: 'block' }}
          >
            {formatPrice(discount.oldPrice)}
          </Typography>
        )}
      </Box>
    </Box>
  );
};

// Компонент MobileGrid - обертка для отображения сетки объявлений
export const MobileListingGrid = ({ listings, viewMode = 'grid' }) => {
  const { t } = useTranslation('marketplace');

  if (viewMode === 'list') {
    return (
      <Box sx={{ px: 1, pt: 1 }}>
        {listings.map((listing) => (
          <Box key={listing.id} component={Link} to={`/marketplace/listings/${listing.id}`} sx={{ textDecoration: 'none', color: 'inherit', display: 'block', mb: 1 }}>
            <MobileListingCard listing={listing} viewMode="list" />
          </Box>
        ))}

        {listings.length === 0 && (
          <Box sx={{ py: 4, textAlign: 'center' }}>
            <Typography variant="body2" color="text.secondary">
              {t('search.noresults', { defaultValue: 'По вашему запросу ничего не найдено' })}
            </Typography>
          </Box>
        )}
      </Box>
    );
  }

  return (
    <Grid container spacing={0}>
      {listings.map((listing) => (
        <Grid item xs={6} key={listing.id} component={Link} to={`/marketplace/listings/${listing.id}`} sx={{ textDecoration: 'none', color: 'inherit' }}>
          <MobileListingCard listing={listing} viewMode="grid" />
        </Grid>
      ))}

      {listings.length === 0 && (
        <Grid item xs={12}>
          <Box sx={{ py: 4, textAlign: 'center' }}>
            <Typography variant="body2" color="text.secondary">
              {t('search.noresults', { defaultValue: 'По вашему запросу ничего не найдено' })}
            </Typography>
          </Box>
        </Grid>
      )}
    </Grid>
  );
};

export const MobileFilters = ({ open, onClose, filters, onFilterChange, categories }) => {
  const { t, i18n } = useTranslation('marketplace'); // Добавляем i18n

  // Исправленная функция getTranslatedName в CategoryItem
  const getTranslatedName = (category) => {
    if (!category) return '';

    // Проверяем наличие переводов
    if (category.translations && typeof category.translations === 'object') {
      // Если есть прямой перевод на текущий язык
      if (category.translations[i18n.language]) {
        return category.translations[i18n.language];
      }

      // Если прямого перевода нет, пробуем найти по приоритету
      const langPriority = [i18n.language, 'ru', 'sr', 'en'];
      for (const lang of langPriority) {
        if (category.translations[lang]) {
          return category.translations[lang];
        }
      }
    }

    // Если переводов нет или они не подходят, возвращаем исходное имя
    return category.name;
  };

  const [tempFilters, setTempFilters] = useState(filters);
  const [currentCategory, setCurrentCategory] = useState(null);
  const [navigationHistory, setNavigationHistory] = useState([]);

  useEffect(() => {
    setTempFilters(filters);
  }, [filters, open]);

  const handleApply = () => {
    onFilterChange(tempFilters);
    onClose();
  };

  // Получаем текущие категории для отображения
  const getCurrentCategories = () => {
    // Если нет текущей категории, показываем только корневые категории
    if (!currentCategory) {
      return categories.filter(category => !category.parent_id) || [];
    }

    // Иначе показываем только дочерние категории текущей категории
    return categories.filter(category =>
      category.parent_id && String(category.parent_id) === String(currentCategory.id)
    ) || [];
  };

  const handleCategoryClick = (category) => {
    const hasChildren = category.children && category.children.length > 0;

    // Устанавливаем выбранную категорию в фильтрах, даже если у неё есть дочерние элементы
    setTempFilters(prev => ({
      ...prev,
      category_id: category.id
    }));

    // Если есть дочерние категории, переходим в них
    if (hasChildren) {
      setNavigationHistory(prev => [...prev, currentCategory]);
      setCurrentCategory(category);
    }

    // Не закрываем фильтры автоматически
  };

  const handleBack = () => {
    if (navigationHistory.length > 0) {
      const newHistory = [...navigationHistory];
      const lastCategory = newHistory.pop();
      setNavigationHistory(newHistory);
      setCurrentCategory(lastCategory);
    }
  };

  const handleClearFilters = () => {
    setTempFilters({
      query: "",
      category_id: "",
      min_price: "",
      max_price: "",
      condition: "",
      sort_by: "date_desc",
      distance: "", // Важно: сбрасываем distance
      latitude: null, // Сбрасываем координаты
      longitude: null
    });
  };

  return (
    <Drawer
      anchor="right"
      open={open}
      onClose={onClose}
      PaperProps={{
        sx: { width: '100%', maxWidth: 400 }
      }}
    >
      <Box sx={{ display: 'flex', flexDirection: 'column', height: '100%' }}>
        {/* Шапка с навигацией */}
        <Box sx={{
          p: 2,
          borderBottom: '1px solid',
          borderColor: 'divider',
          display: 'flex',
          alignItems: 'center',
          gap: 1
        }}>
          {navigationHistory.length > 0 && (
            <IconButton
              onClick={handleBack}
              sx={{
                color: 'text.secondary',
                '&:hover': { color: 'primary.main' }
              }}
            >
              <ArrowLeft size={20} />
            </IconButton>
          )}

          <Typography
            variant="subtitle1"
            sx={{
              flex: 1,
              fontWeight: 600,
              color: 'text.primary'
            }}
          >
            {currentCategory ? getTranslatedName(currentCategory) : t('listings.filters.title')}
          </Typography>

          <Button
            variant="text"
            size="small"
            onClick={handleClearFilters}
            sx={{ color: 'text.secondary' }}
          >
            {t('listings.filters.reset')}
          </Button>
        </Box>

        {/* Основное содержимое с фильтрами */}
        <Box sx={{ flex: 1, overflow: 'auto' }}>
          {/* Фильтр цены */}
          <Box sx={{ p: 2, borderBottom: '1px solid', borderColor: 'divider' }}>
            <Typography variant="subtitle2" gutterBottom>
              {t('listings.filters.price.label')}
            </Typography>

            <Stack direction="row" spacing={1}>
              <TextField
                fullWidth
                size="small"
                placeholder={t('listings.filters.price.min')}
                type="number"
                value={tempFilters.min_price || ''}
                onChange={(e) => setTempFilters(prev => ({
                  ...prev,
                  min_price: e.target.value
                }))}
              />

              <TextField
                fullWidth
                size="small"
                placeholder={t('listings.filters.price.max')}
                type="number"
                value={tempFilters.max_price || ''}
                onChange={(e) => setTempFilters(prev => ({
                  ...prev,
                  max_price: e.target.value
                }))}
              />
            </Stack>
          </Box>

          {/* Фильтр состояния */}
          <Box sx={{ p: 2, borderBottom: '1px solid', borderColor: 'divider' }}>
            <Typography variant="subtitle2" gutterBottom>
              {t('listings.details.condition')}
            </Typography>

            <Select
              fullWidth
              size="small"
              value={tempFilters.condition || ''}
              onChange={(e) => setTempFilters(prev => ({
                ...prev,
                condition: e.target.value
              }))}
            >
              <MenuItem value="">{t('listings.create.condition.any')}</MenuItem>
              <MenuItem value="new">{t('listings.create.condition.new')}</MenuItem>
              <MenuItem value="used">{t('listings.create.condition.used')}</MenuItem>
            </Select>
          </Box>

          {/* Сортировка */}
          <Box sx={{ p: 2, borderBottom: '1px solid', borderColor: 'divider' }}>
            <Typography variant="subtitle2" gutterBottom>
              {t('listings.filters.sort.label')}
            </Typography>

            <Select
              fullWidth
              size="small"
              value={tempFilters.sort_by || 'date_desc'}
              onChange={(e) => setTempFilters(prev => ({
                ...prev,
                sort_by: e.target.value
              }))}
            >
              <MenuItem value="date_desc">{t('listings.filters.sort.newest')}</MenuItem>
              <MenuItem value="price_asc">{t('listings.filters.sort.cheapest')}</MenuItem>
              <MenuItem value="price_desc">{t('listings.filters.sort.expensive')}</MenuItem>
            </Select>
          </Box>

          {/* Категории */}
          <Box sx={{ borderBottom: '1px solid', borderColor: 'divider' }}>
            {getCurrentCategories().map((category) => {
              const hasChildren = category.children && category.children.length > 0;
              const isSelected = tempFilters.category_id === category.id;

              return (
                <Box
                  key={category.id}
                  sx={{
                    borderBottom: '1px solid',
                    borderColor: 'divider',
                    '&:last-child': {
                      borderBottom: 'none'
                    }
                  }}
                >
                  <Button
                    onClick={() => handleCategoryClick(category)}
                    sx={{
                      width: '100%',
                      justifyContent: 'flex-start',
                      textTransform: 'none',
                      py: 2,
                      px: 2,
                      color: isSelected ? 'primary.main' : 'text.primary',
                      backgroundColor: isSelected ? 'action.selected' : 'transparent',
                      borderRadius: 0,
                      '&:hover': {
                        backgroundColor: isSelected ? 'action.selected' : 'action.hover'
                      }
                    }}
                  >
                    <Typography
                      sx={{
                        flex: 1,
                        textAlign: 'left',
                        fontWeight: isSelected ? 600 : 400,
                        color: isSelected ? 'primary.main' : 'text.primary'
                      }}
                    >
                      {getTranslatedName(category)}
                      {isSelected && (
                        <Box component="span" sx={{ ml: 1, color: 'primary.main' }}>
                          ✓
                        </Box>
                      )}
                    </Typography>

                    {hasChildren && (
                      <ChevronRight
                        size={20}
                        style={{
                          opacity: 0.5,
                          marginLeft: 8
                        }}
                      />
                    )}
                  </Button>
                </Box>
              );
            })}
          </Box>
        </Box>

        {/* Кнопки действий */}
        <Box sx={{ p: 2, borderTop: '1px solid', borderColor: 'divider' }}>
          <Button
            variant="contained"
            fullWidth
            onClick={handleApply}
            startIcon={<Check size={20} />}
            sx={{
              py: 1.5,
              textTransform: 'none',
              fontWeight: 500,
            }}
          >
            {t('listings.filters.apply')}
          </Button>
        </Box>
      </Box>
    </Drawer>
  );
};