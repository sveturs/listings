import React, { useState, useEffect } from 'react';
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
import {
  Edit,
  Eye,
  Tag,
  Clock,
  MapPin,
  Star
} from 'lucide-react';
import { Link } from 'react-router-dom';
import axios from '../../api/axios';
import { Listing } from '../marketplace/ListingCard';

interface Category {
  id: number | string;
  name: string;
  slug: string;
  translations?: {
    [language: string]: string;
  };
  [key: string]: any;
}

interface StorefrontListingsListProps {
  listings: Listing[];
  selectedItems: Array<number | string>;
  onSelectItem: (id: number | string) => void;
  onSelectAll: (selected: boolean) => void;
}

const StorefrontListingsList: React.FC<StorefrontListingsListProps> = ({ 
  listings, 
  selectedItems, 
  onSelectItem, 
  onSelectAll 
}) => {
  const { t, i18n } = useTranslation(['marketplace', 'common']);
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('sm'));
  const [allCategories, setAllCategories] = useState<Category[]>([]);
  const [categoriesMap, setCategoriesMap] = useState<Record<string | number, Category>>({});

  // Проверяем, выбраны ли все объявления
  const isAllSelected = listings.length > 0 && selectedItems.length === listings.length;

  // Загрузка всех категорий при монтировании компонента
  useEffect(() => {
    const fetchCategories = async () => {
      try {
        const response = await axios.get('/api/v1/marketplace/categories');
        if (response.data.success) {
          const categories = response.data.data as Category[];
          setAllCategories(categories);
          
          // Создаем мапу для быстрого доступа к категориям по id
          const map: Record<string | number, Category> = {};
          categories.forEach(cat => {
            map[cat.id] = cat;
          });
          setCategoriesMap(map);
          
          console.log('Загружено категорий:', categories.length);
          if (categories.length > 0) {
            console.log('Пример категории с переводами:', categories[0]);
          }
        }
      } catch (error) {
        console.error('Ошибка при загрузке категорий:', error);
      }
    };
    
    fetchCategories();
  }, []);

  // Форматирование цены
  const formatPrice = (price: number): string => {
    // Используем локаль из текущего языка интерфейса
    const locale = i18n.language === 'ru' ? 'ru-RU' :
      i18n.language === 'sr' ? 'sr-RS' : 'en-US';

    return new Intl.NumberFormat(locale, {
      style: 'currency',
      currency: 'RSD',
      maximumFractionDigits: 0
    }).format(price);
  };

  // Форматирование даты
  const formatDate = (dateString?: string): string => {
    if (!dateString) return '';
    
    try {
      const date = new Date(dateString);
      // Используем локаль из текущего языка интерфейса
      const locale = i18n.language === 'ru' ? 'ru-RU' :
        i18n.language === 'sr' ? 'sr-RS' : 'en-US';

      return new Intl.DateTimeFormat(locale, {
        day: '2-digit',
        month: '2-digit',
        year: 'numeric'
      }).format(date);
    } catch (e) {
      return dateString;
    }
  };

  // Функция для получения локализованного текста
  const getLocalizedText = (listing: Listing | null, field: string): string => {
    if (!listing) return '';

    // Получаем перевод для текущего языка интерфейса
    const currentLanguage = i18n.language;

    // Проверяем, совпадает ли текущий язык с оригинальным языком объявления
    if (currentLanguage === listing.original_language) {
      return listing[field] as string;
    }

    // Проверяем наличие перевода
    if (listing.translations &&
      listing.translations[currentLanguage] &&
      listing.translations[currentLanguage][field]) {
      return listing.translations[currentLanguage][field];
    }

    // Если перевода нет, возвращаем оригинальный текст
    return listing[field] as string;
  };

  // Функция для получения URL изображения - скопировано из ListingCard.js
  const getImageUrl = (image: string | any): string => {
    if (!image) {
      return '/placeholder.jpg';
    }

    // Используем переменную окружения из window.ENV вместо process.env
    const baseUrl = window.ENV?.REACT_APP_MINIO_URL || window.ENV?.REACT_APP_BACKEND_URL || '';

    // 1. Если image это строка, это простой путь к файлу
    if (typeof image === 'string') {
      // Относительный путь MinIO
      if (image.startsWith('/listings/')) {
        return `${baseUrl}${image}`;
      }

      // ID/filename.jpg (прямой путь MinIO)
      if (image.match(/^\d+\/[^\/]+$/)) {
        return `${baseUrl}/listings/${image}`;
      }

      // Локальное хранилище (обратная совместимость)
      return `${baseUrl}/uploads/${image}`;
    }

    // 2. Объект с информацией о файле
    if (typeof image === 'object' && image !== null) {
      // Приоритет 1: Используем PublicURL если он доступен
      if (image.public_url && typeof image.public_url === 'string' && image.public_url.trim() !== '') {
        const publicUrl = image.public_url;

        // Абсолютный URL
        if (publicUrl.startsWith('http')) {
          return publicUrl;
        }
        // Относительный URL с /listings/
        else if (publicUrl.startsWith('/listings/')) {
          return `${baseUrl}${publicUrl}`;
        }
        // Другой относительный URL
        else {
          return `${baseUrl}${publicUrl}`;
        }
      }

      // Приоритет 2: Формируем URL на основе типа хранилища и пути к файлу
      if (image.file_path) {
        if (image.storage_type === 'minio' || image.file_path.includes('listings/')) {
          // Учитываем возможность наличия префикса listings/ в пути
          const filePath = image.file_path.includes('listings/')
            ? image.file_path.replace('listings/', '')
            : image.file_path;

          return `${baseUrl}/listings/${filePath}`;
        }

        // Локальное хранилище
        return `${baseUrl}/uploads/${image.file_path}`;
      }
    }

    return '/placeholder.jpg';
  };

  // Функция для получения переведенного названия категории
  const getTranslatedCategoryName = (category: any): string => {
    if (!category) return '';
    
    // Проверяем, есть ли у категории id и существует ли она в нашей карте
    if (category.id && categoriesMap[category.id]) {
      // Ищем категорию в нашей карте категорий
      const fullCategory = categoriesMap[category.id];
      if (fullCategory && fullCategory.translations) {
        // Если найдена полная информация о категории, используем её переводы
        if (fullCategory.translations[i18n.language]) {
          return fullCategory.translations[i18n.language];
        }
        
        // Если нет прямого перевода, ищем по приоритету
        const langPriority = [i18n.language, 'ru', 'sr', 'en'];
        for (const lang of langPriority) {
          if (fullCategory.translations[lang]) {
            return fullCategory.translations[lang];
          }
        }
      }
    }
    
    // Если категория не найдена в нашей карте, пробуем найти по имени или slug
    if ((category.name || category.slug) && allCategories.length > 0) {
      const foundCategory = allCategories.find(fullCat => 
        fullCat.name === category.name || fullCat.slug === category.slug
      );
      
      if (foundCategory && foundCategory.translations) {
        if (foundCategory.translations[i18n.language]) {
          return foundCategory.translations[i18n.language];
        }
        
        // Если нет прямого перевода, ищем по приоритету
        const langPriority = [i18n.language, 'ru', 'sr', 'en'];
        for (const lang of langPriority) {
          if (foundCategory.translations[lang]) {
            return foundCategory.translations[lang];
          }
        }
      }
    }
    
    // Проверяем наличие переводов в самой категории
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
    
    // Хардкодные соответствия для некоторых известных категорий
    if (category.name === "Чехлы и плёнки" || category.slug === "phone-cases") {
      if (i18n.language === 'en') return 'Cases and Screen Protectors111';
      if (i18n.language === 'sr') return 'Maske i zaštitna stakla';
    }
    
    // Если все методы не сработали, возвращаем исходное имя
    return category.name || '';
  };

  // Локализация статусов через i18next
  const getLocalizedStatus = (status?: string): string => {
    if (!status) return '';

    // Приводим статус к нижнему регистру для надёжного сравнения
    const normalizedStatus = status.toLowerCase();
    
    if (normalizedStatus === 'active' || normalizedStatus === 'активно') {
      // Используем прямые значения в зависимости от языка интерфейса
      if (i18n.language === 'en') return 'Active';
      if (i18n.language === 'sr') return 'Aktivno';
      return 'Активно';
    }
    
    if (normalizedStatus === 'inactive' || normalizedStatus === 'неактивно') {
      if (i18n.language === 'en') return 'Inactive';
      if (i18n.language === 'sr') return 'Neaktivno';
      return 'Неактивно';
    }

    return status;
  };

  // Функция для получения локализованного текста UI компонентов
  const getUITranslation = (key: string, defaultValue: string, options?: any): string => {
    return t(`marketplace:${key}`, { defaultValue, ...options });
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
              {getUITranslation('store.listings.title', 'Наименование')}
            </Typography>
          </Grid>
          <Grid item xs={isMobile ? 3 : 2} md={2} sx={{ textAlign: isMobile ? 'center' : 'right' }}>
            <Typography variant="subtitle2" sx={{ fontWeight: 'bold' }}>
              {getUITranslation('store.listings.price', 'Цена')}
            </Typography>
          </Grid>
          <Grid item xs={2} md={1} sx={{ textAlign: 'center' }}>
            <Typography variant="subtitle2" sx={{ fontWeight: 'bold' }}>
              {getUITranslation('store.listings.status', 'Статус')}
            </Typography>
          </Grid>
          <Grid item xs={2} md={2} sx={{ textAlign: 'center' }}>
            <Typography variant="subtitle2" sx={{ fontWeight: 'bold' }}>
              {getUITranslation('store.listings.date', 'Дата')}
            </Typography>
          </Grid>
        </Grid>
      </Paper>

      {listings.length === 0 ? (
        <Box py={4} textAlign="center">
          <Typography color="text.secondary">
            {getUITranslation('store.listings.noListings', 'Нет объявлений')}
          </Typography>
        </Box>
      ) : (
        <Box>
          {listings.map((listing) => {
            // Получаем локализованные тексты
            const localizedTitle = getLocalizedText(listing, 'title');
            const localizedDescription = getLocalizedText(listing, 'description');

            return (
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
                          src={getImageUrl(listing.images[0])}
                          alt={localizedTitle}
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
                            '&:hover': { color: theme.palette.primary?.main }
                          }}>
                            {localizedTitle}
                          </Typography>
                        </Link>

                        <Box display="flex" alignItems="center" flexWrap="wrap" gap={1}>
                          {!isMobile && (
                            <>
                              <Chip
                                label={listing.condition === 'new' ? getUITranslation('listings.condition.new', 'Новое') : getUITranslation('listings.condition.used', 'Б/у')}
                                size="small"
                                variant="outlined"
                              />

                              {listing.category && (
                                <Chip
                                  icon={<Tag size={14} />}
                                  label={getTranslatedCategoryName(listing.category)}
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
                      label={getLocalizedStatus(listing.status)}
                      size="small"
                      color={listing.status === 'active' || listing.status === 'Активно' ? 'success' : 'default'}
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
                          {getUITranslation('listings.views', 'Просмотров: {{count}}', { count: listing.views_count || 0 })}
                        </Typography>
                      )}
                    </Box>
                  </Grid>
                </Grid>
              </Paper>
            );
          })}
        </Box>
      )}
    </Box>
  );
};

export default StorefrontListingsList;