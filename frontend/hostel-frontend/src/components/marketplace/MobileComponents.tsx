// frontend/hostel-frontend/src/components/marketplace/MobileComponents.tsx
import React, { useState, useCallback, useEffect, MouseEvent, TouchEvent } from 'react';
import { useTranslation } from 'react-i18next';
import AutocompleteInput from '../shared/AutocompleteInput';
import { Link } from 'react-router-dom';
import { Listing, ListingImage } from '../../types/listing';
import {
  Box, Button, IconButton, Typography, InputBase, Toolbar, TextField, Select, MenuItem,
  Paper, Grid, Drawer, Stack, List, ListItem, ListItemText, ListItemAvatar, Avatar, Divider,
  ToggleButton, ToggleButtonGroup, SelectChangeEvent
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
  Percent,
  Eye
} from 'lucide-react';
import { debounce } from 'lodash';

// Type definitions
interface Category {
  id: string | number;
  name: string;
  parent_id?: string | number | null;
  children?: Category[];
  translations?: {
    [language: string]: string;
  };
}

interface DiscountInfo {
  percent: number;
  oldPrice: number;
}

interface FilterOptions {
  query: string;
  category_id: string | number | null;
  min_price: string | number | null;
  max_price: string | number | null;
  condition: string | null;
  sort_by: string;
  distance?: string | number | null;
  latitude?: number | null;
  longitude?: number | null;
  [key: string]: any; // For other filter properties
}

// Component props interfaces
interface MobileHeaderProps {
  onOpenFilters: () => void;
  filtersCount: number;
  onSearch: (value: string) => void;
  searchValue: string;
  viewMode: 'grid' | 'list';
  onViewModeChange: (event: React.MouseEvent<HTMLElement>, newMode: 'grid' | 'list' | null) => void;
}

interface MobileListingCardProps {
  listing: Listing;
  viewMode?: 'grid' | 'list';
}

interface MobileListingGridProps {
  listings: Listing[];
  viewMode?: 'grid' | 'list';
}

interface MobileFiltersProps {
  open: boolean;
  onClose: () => void;
  filters: FilterOptions;
  onFilterChange: (filters: FilterOptions | Partial<FilterOptions>) => void;
  categories: Category[];
  onToggleMapView?: () => void;
}

// Компонент MobileHeader
export const MobileHeader: React.FC<MobileHeaderProps> = ({ 
  onOpenFilters, 
  filtersCount, 
  onSearch, 
  searchValue, 
  viewMode, 
  onViewModeChange 
}) => {
  const { t } = useTranslation(['marketplace', 'common']);

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
const formatPrice = (price?: number): string => {
  return new Intl.NumberFormat('sr-RS', {
    style: 'currency',
    currency: 'RSD',
    maximumFractionDigits: 0
  }).format(price || 0);
};

// Получаем URL изображения
const getImageUrl = (listing?: Listing): string => {
  if (!listing || !listing.images || !Array.isArray(listing.images) || listing.images.length === 0) {
    return '/placeholder.jpg';
  }

  const baseUrl = (window as any).ENV?.REACT_APP_MINIO_URL || (window as any).ENV?.REACT_APP_BACKEND_URL || '';

  // Находим главное изображение или используем первое в списке
  const mainImage = listing.images.find(img => typeof img === 'object' && img.is_main) || listing.images[0];

  // Если это строка URL
  if (typeof mainImage === 'string') {
    // Проверяем, абсолютный или относительный URL
    if (mainImage.startsWith('http')) {
      return mainImage;
    } else {
      return `${baseUrl}/uploads/${mainImage}`;
    }
  }

  // Если это объект
  if (mainImage && typeof mainImage === 'object') {
    // Проверяем на наличие публичного URL
    if (mainImage.public_url) {
      // Проверяем, абсолютный или относительный URL
      if (mainImage.public_url.startsWith('http')) {
        return mainImage.public_url;
      } else {
        return `${baseUrl}${mainImage.public_url}`;
      }
    }

    // Проверяем на MinIO хранилище
    if (mainImage.storage_type === 'minio' ||
        (mainImage.file_path && mainImage.file_path.includes('listings/'))) {
      return `${baseUrl}${mainImage.public_url}`;
    }

    // Проверяем file_path
    if (mainImage.file_path) {
      return `${baseUrl}/uploads/${mainImage.file_path}`;
    }
  }

  return '/placeholder.jpg';
};

// Функция для определения браузера Firefox
const isFirefox = (): boolean => {
  return typeof window !== 'undefined' && window.navigator.userAgent.indexOf('Firefox') !== -1;
};

// Получаем информацию о скидке
const getDiscountInfo = (listing: Listing): DiscountInfo | null => {
  if (listing.metadata && listing.metadata.discount) {
    return {
      percent: listing.metadata.discount.discount_percent || listing.metadata.discount.percent || 0,
      oldPrice: listing.metadata.discount.previous_price || listing.metadata.discount.oldPrice || 0
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
const formatDate = (dateString?: string): string => {
  if (!dateString) return '';
  return new Date(dateString).toLocaleDateString();
};

// Компонент MobileListingCard
export const MobileListingCard: React.FC<MobileListingCardProps> = ({ listing, viewMode = 'grid' }) => {
  const { t, i18n } = useTranslation('marketplace');

  const getTranslatedText = (field: keyof Listing): string => {
    if (!listing) return '';
    if (i18n.language === listing.original_language) {
      return listing[field] as string;
    }
    return listing.translations?.[i18n.language]?.[field] || (listing[field] as string);
  };

  // Функция для определения, является ли объявление частью витрины
  // Временное решение, пока не будет исправлена серверная часть
  const isStoreItem = (): boolean => {
    // Единственное надежное условие - наличие storefront_id
    return listing.storefront_id !== undefined && listing.storefront_id !== null;
  };

  // Получаем ID витрины
  const getStoreId = (): number | string | null => {
    // Просто возвращаем storefront_id, если он есть
    return listing.storefront_id || null;
  };

  const handleShopButtonClick = (e: React.MouseEvent | React.TouchEvent): boolean => {
    e.preventDefault();
    e.stopPropagation();
    // Останавливаем всплытие события для всех родительских элементов
    if (e.nativeEvent) {
      e.nativeEvent.stopImmediatePropagation();
    }
    // Используем прямой переход, с нашей функцией определения ID витрины
    window.location.href = `/shop/${getStoreId()}`;
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
              <Box sx={{ display: 'flex', alignItems: 'center' }}>
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
                
                {/* View count indicator - always show even if zero */}
                <Box sx={{ display: 'flex', alignItems: 'center', ml: 1 }}>
                  <Eye size={12} color="#666" />
                  <Typography variant="caption" color="text.secondary" sx={{ ml: 0.25 }}>
                    {listing.views_count || 0}
                  </Typography>
                </Box>
              </Box>

              {/* Бейдж магазина */}
              {/* Проверяем бирку магазина для режима списка - используем функцию */}
              {isStoreItem() && (
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
        {/* Контейнер для бирок (магазин и скидка) */}
        <Box>
          {/* бейдж магазина - используем функцию определения товаров из витрины */}
          {isStoreItem() && (
            <Box
              sx={{
                position: 'absolute',
                top: 10,
                left: 10,
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
              onMouseDown={(e: React.MouseEvent) => {
                e.preventDefault();
                e.stopPropagation();
                if (e.nativeEvent) {
                  e.nativeEvent.stopImmediatePropagation();
                }
                window.location.href = `/shop/${getStoreId()}`;
                return false;
              }}
              // Дублируем на onClick для надежности
              onClick={(e: React.MouseEvent) => {
                e.preventDefault();
                e.stopPropagation();
                if (e.nativeEvent) {
                  e.nativeEvent.stopImmediatePropagation();
                }
                window.location.href = `/shop/${getStoreId()}`;
                return false;
              }}
              // Дублируем на onTouchStart для надежности на мобильных
              onTouchStart={(e: React.TouchEvent) => {
                e.preventDefault();
                e.stopPropagation();
                if (e.nativeEvent) {
                  e.nativeEvent.stopImmediatePropagation();
                }
                window.location.href = `/shop/${getStoreId()}`;
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
                right: 10,
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
              {`-${discount.percent}%`}
            </Box>
          )}
        </Box>

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

        <Box sx={{ display: 'flex', alignItems: 'center', mt: 0.5 }}>
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
          
          {/* View count indicator - always show even if zero */}
          <Box sx={{ display: 'flex', alignItems: 'center', ml: 1 }}>
            <Eye size={12} color="#666" />
            <Typography variant="caption" color="text.secondary" sx={{ ml: 0.25 }}>
              {listing.views_count || 0}
            </Typography>
          </Box>
        </Box>

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
export const MobileListingGrid: React.FC<MobileListingGridProps> = ({ listings, viewMode = 'grid' }) => {
  const { t } = useTranslation('marketplace');
  const [isFirefoxBrowser, setIsFirefoxBrowser] = useState<boolean>(false);

  useEffect(() => {
    setIsFirefoxBrowser(isFirefox());
  }, []);

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

  // Специальный рендеринг для Firefox
  if (isFirefoxBrowser) {
    return (
      <div style={{
        display: 'flex',
        flexWrap: 'wrap',
        width: '100%',
        margin: '0 auto', // Центрирование контейнера
        padding: 0
      }}>
        {listings.map((listing) => (
          <div
            key={listing.id}
            style={{
              width: '50%',
              boxSizing: 'border-box',
              padding: '4px',
              float: 'left',
              maxWidth: '50%' // Принудительная максимальная ширина 50%
            }}
          >
            <Link
              to={`/marketplace/listings/${listing.id}`}
              style={{
                textDecoration: 'none',
                color: 'inherit',
                display: 'block',
                width: '100%' // Полная ширина внутри контейнера
              }}
            >
              <MobileListingCard listing={listing} viewMode="grid" />
            </Link>
          </div>
        ))}

        {listings.length === 0 && (
          <div style={{ width: '100%', padding: '16px 0', textAlign: 'center' }}>
            <Typography variant="body2" color="text.secondary">
              {t('search.noresults', { defaultValue: 'По вашему запросу ничего не найдено' })}
            </Typography>
          </div>
        )}
      </div>
    );
  }

  // Стандартный рендеринг для других браузеров
  return (
    <Box className="MobileListingGrid" sx={{ width: '100%', display: 'flex', flexWrap: 'wrap' }}>
      {listings.map((listing) => (
        <Box
          key={listing.id}
          component={Link}
          to={`/marketplace/listings/${listing.id}`}
          sx={{
            width: '50%',
            boxSizing: 'border-box',
            textDecoration: 'none',
            color: 'inherit',
            padding: '4px'
          }}
        >
          <MobileListingCard listing={listing} viewMode="grid" />
        </Box>
      ))}

      {listings.length === 0 && (
        <Box sx={{ width: '100%', py: 4, textAlign: 'center' }}>
          <Typography variant="body2" color="text.secondary">
            {t('search.noresults', { defaultValue: 'По вашему запросу ничего не найдено' })}
          </Typography>
        </Box>
      )}
    </Box>
  );
};

export const MobileFilters: React.FC<MobileFiltersProps> = ({ open, onClose, filters, onFilterChange, categories, onToggleMapView }) => {
  const { t, i18n } = useTranslation('marketplace'); // Добавляем i18n

  // Исправленная функция getTranslatedName в CategoryItem
  const getTranslatedName = (category: Category | null): string => {
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

  const [tempFilters, setTempFilters] = useState<FilterOptions>(filters);
  const [currentCategory, setCurrentCategory] = useState<Category | null>(null);
  const [navigationHistory, setNavigationHistory] = useState<Array<Category | null>>([]);

  useEffect(() => {
    setTempFilters(filters);
  }, [filters, open]);

  const handleApply = (): void => {
    onFilterChange(tempFilters);
    onClose();
  };

  // Получаем текущие категории для отображения
  const getCurrentCategories = (): Category[] => {
    // Если нет текущей категории, показываем только корневые категории
    if (!currentCategory) {
      return categories.filter(category => !category.parent_id) || [];
    }

    // Иначе показываем только дочерние категории текущей категории
    return categories.filter(category =>
      category.parent_id && String(category.parent_id) === String(currentCategory.id)
    ) || [];
  };

  const handleCategoryClick = (category: Category): void => {
    const hasChildren = categories.some(c => 
      c.parent_id && String(c.parent_id) === String(category.id)
    );

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

  const handleBack = (): void => {
    if (navigationHistory.length > 0) {
      const newHistory = [...navigationHistory];
      const lastCategory = newHistory.pop();
      setNavigationHistory(newHistory);
      setCurrentCategory(lastCategory);
    }
  };

  const handleClearFilters = (): void => {
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
              onChange={(e: SelectChangeEvent<string>) => setTempFilters(prev => ({
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
              onChange={(e: SelectChangeEvent<string>) => setTempFilters(prev => ({
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
              const hasChildren = categories.some(c => 
                c.parent_id && String(c.parent_id) === String(category.id)
              );
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