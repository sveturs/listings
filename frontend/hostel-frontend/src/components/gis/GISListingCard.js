import React from 'react';
import { useTranslation } from 'react-i18next';
import { useNavigate } from 'react-router-dom';
import {
  Card,
  CardActionArea,
  CardMedia,
  CardContent,
  Typography,
  Box,
  Button,
  IconButton,
  Chip,
  Divider,
  Tooltip
} from '@mui/material';
import {
  FavoriteBorder as FavoriteBorderIcon,
  Favorite as FavoriteIcon,
  LocationOn as LocationIcon,
  Place as PlaceIcon,
  Phone as PhoneIcon,
  Message as MessageIcon
} from '@mui/icons-material';
import { styled } from '@mui/material/styles';
import { formatDistanceToNow } from 'date-fns';
import { enUS, ru, sr } from 'date-fns/locale';

const StyledCard = styled(Card)(({ theme }) => ({
  width: '100%',
  display: 'flex',
  marginBottom: theme.spacing(1.5),
  borderRadius: theme.spacing(1),
  overflow: 'visible',
  boxShadow: '0 3px 10px rgba(0, 0, 0, 0.08)',
  '&:hover': {
    transform: 'translateY(-2px)',
    boxShadow: '0 5px 15px rgba(0, 0, 0, 0.1)',
    transition: 'transform 0.2s ease-out, box-shadow 0.2s ease-out'
  },
  [theme.breakpoints.down('sm')]: {
    flexDirection: 'column',
    height: 'auto'
  }
}));

const CardImage = styled(CardMedia)(({ theme }) => ({
  width: 160,
  height: 160,
  borderRadius: `${theme.spacing(1)}px 0 0 ${theme.spacing(1)}px`,
  backgroundSize: 'cover',
  [theme.breakpoints.down('sm')]: {
    width: '100%',
    height: 180,
    borderRadius: `${theme.spacing(1)}px ${theme.spacing(1)}px 0 0`
  }
}));

const ContentContainer = styled(Box)(({ theme }) => ({
  flex: 1,
  display: 'flex',
  flexDirection: 'column',
  padding: theme.spacing(1.5),
  overflow: 'hidden'
}));

const PriceText = styled(Typography)(({ theme }) => ({
  fontWeight: 'bold',
  color: theme.palette.primary.main,
  marginRight: theme.spacing(1)
}));

const ConditionChip = styled(Chip)(({ theme, condition }) => ({
  backgroundColor: condition === 'new' 
    ? theme.palette.success.light 
    : theme.palette.grey[300],
  color: condition === 'new' 
    ? theme.palette.success.contrastText 
    : theme.palette.text.primary,
  fontSize: '0.75rem',
  height: 24
}));

const LocationText = styled(Typography)(({ theme }) => ({
  display: 'flex',
  alignItems: 'center',
  color: theme.palette.text.secondary,
  fontSize: '0.875rem',
  marginTop: theme.spacing(0.5),
  '& svg': {
    fontSize: '1rem',
    marginRight: theme.spacing(0.5)
  }
}));

const ActionButtons = styled(Box)(({ theme }) => ({
  display: 'flex',
  justifyContent: 'space-between',
  alignItems: 'center',
  marginTop: 'auto',
  paddingTop: theme.spacing(1),
  [theme.breakpoints.down('sm')]: {
    flexDirection: 'column',
    alignItems: 'stretch',
    gap: theme.spacing(1)
  }
}));

const DateText = styled(Typography)(({ theme }) => ({
  fontSize: '0.75rem',
  color: theme.palette.text.secondary,
  marginLeft: theme.spacing(1),
  [theme.breakpoints.down('sm')]: {
    marginLeft: 0
  }
}));

const FavoriteButton = styled(IconButton)(({ theme }) => ({
  position: 'absolute',
  top: theme.spacing(1),
  right: theme.spacing(1),
  backgroundColor: 'rgba(255, 255, 255, 0.9)',
  padding: theme.spacing(0.5),
  '&:hover': {
    backgroundColor: 'rgba(255, 255, 255, 1)'
  }
}));

const TruncatedTitle = styled(Typography)(({ theme }) => ({
  overflow: 'hidden',
  textOverflow: 'ellipsis',
  display: '-webkit-box',
  WebkitLineClamp: 2,
  WebkitBoxOrient: 'vertical'
}));

const GISListingCard = ({ 
  listing, 
  compact = false,
  isFavorite = false,
  showMap = true, 
  onFavoriteToggle,
  onShowOnMap,
  onContactClick
}) => {
  const { t, i18n } = useTranslation(['gis', 'marketplace']);
  const navigate = useNavigate();

  // Handle different date-fns locales
  const getLocale = () => {
    switch(i18n.language) {
      case 'ru': return ru;
      case 'en': return enUS;
      case 'sr': return sr;
      default: return sr;
    }
  };

  const formatDate = (dateString) => {
    try {
      if (!dateString) return '';
      
      // Ensure proper date format
      const date = new Date(dateString);
      
      // Check if date is valid
      if (isNaN(date.getTime())) {
        return dateString;
      }
      
      // Use simple formatting as fallback
      return date.toLocaleDateString(i18n.language, {
        year: 'numeric',
        month: 'short',
        day: 'numeric'
      });
    } catch (error) {
      console.error('Error formatting date:', error);
      return dateString;
    }
  };

  const handleCardClick = () => {
    navigate(`/marketplace/listings/${listing.id}`);
  };

  const handleFavoriteClick = (e) => {
    e.stopPropagation();
    if (onFavoriteToggle) {
      onFavoriteToggle(listing.id, !isFavorite);
    }
  };

  const handleShowOnMap = (e) => {
    e.stopPropagation();
    if (onShowOnMap) {
      onShowOnMap(listing);
    }
  };

  const handleContactClick = (e) => {
    e.stopPropagation();
    if (onContactClick) {
      onContactClick(listing);
    }
  };

  const placeholderImage = '/placeholder-listing.jpg';
  
  // Оптимизированная функция получения URL изображения, полностью совместимая с ListingCard
  const getImageUrl = (imageData) => {
    if (!imageData) {
      return placeholderImage;
    }

    // Используем переменную окружения из window.ENV вместо process.env
    const baseUrl = window.ENV?.REACT_APP_MINIO_URL || window.ENV?.REACT_APP_BACKEND_URL || '';
    console.log('GISListingCard: Using baseUrl from env:', baseUrl);

    // 1. Строковые пути (для обратной совместимости)
    if (typeof imageData === 'string') {
      console.log('GISListingCard: Processing string image path:', imageData);

      // Относительный путь MinIO
      if (imageData.startsWith('/listings/')) {
        const url = `${baseUrl}${imageData}`;
        console.log('GISListingCard: Using MinIO relative path:', url);
        return url;
      }

      // ID/filename.jpg (прямой путь MinIO)
      if (imageData.match(/^\d+\/[^\/]+$/)) {
        const url = `${baseUrl}/listings/${imageData}`;
        console.log('GISListingCard: Using direct MinIO path pattern:', url);
        return url;
      }

      // Локальное хранилище (обратная совместимость)
      const url = `${baseUrl}/uploads/${imageData}`;
      console.log('GISListingCard: Using local storage path:', url);
      return url;
    }

    // 2. Объекты с информацией о файле
    if (typeof imageData === 'object' && imageData !== null) {
      console.log('GISListingCard: Processing image object:', imageData);

      // Приоритет 1: Используем PublicURL если он доступен
      if (imageData.public_url && typeof imageData.public_url === 'string' && imageData.public_url.trim() !== '') {
        const publicUrl = imageData.public_url;
        console.log('GISListingCard: Found public_url string:', publicUrl);

        // Абсолютный URL
        if (publicUrl.startsWith('http')) {
          console.log('GISListingCard: Using absolute URL:', publicUrl);
          return publicUrl;
        }
        // Относительный URL с /listings/
        else if (publicUrl.startsWith('/listings/')) {
          const url = `${baseUrl}${publicUrl}`;
          console.log('GISListingCard: Using public_url with listings path:', url);
          return url;
        }
        // Другой относительный URL
        else {
          const url = `${baseUrl}${publicUrl}`;
          console.log('GISListingCard: Using general relative public_url:', url);
          return url;
        }
      }

      // Приоритет 2: Формируем URL на основе типа хранилища и пути к файлу
      if (imageData.file_path) {
        if (imageData.storage_type === 'minio' || imageData.file_path.includes('listings/')) {
          // Учитываем возможность наличия префикса listings/ в пути
          const filePath = imageData.file_path.includes('listings/')
            ? imageData.file_path.replace('listings/', '')
            : imageData.file_path;

          const url = `${baseUrl}/listings/${filePath}`;
          console.log('GISListingCard: Constructed MinIO URL from path:', url);
          return url;
        }

        // Локальное хранилище
        const url = `${baseUrl}/uploads/${imageData.file_path}`;
        console.log('GISListingCard: Using local storage path from object:', url);
        return url;
      }
    }

    console.log('GISListingCard: Could not determine image URL, using placeholder');
    return placeholderImage;
  };

  return (
    <StyledCard>
      <Box position="relative">
        <CardActionArea onClick={handleCardClick}>
          <CardImage
            image={listing.images && listing.images.length > 0 ? getImageUrl(listing.images[0]) : placeholderImage}
            title={listing.title}
          />
        </CardActionArea>
        <FavoriteButton 
          size="small" 
          onClick={handleFavoriteClick}
          aria-label={isFavorite ? t('card.unfavorite') : t('card.favorite')}
        >
          {isFavorite ? 
            <FavoriteIcon fontSize="small" color="error" /> : 
            <FavoriteBorderIcon fontSize="small" />
          }
        </FavoriteButton>
      </Box>

      <ContentContainer>
        <Box mb={1}>
          <Box display="flex" justifyContent="space-between" alignItems="flex-start">
            <TruncatedTitle variant="subtitle1" component="h2">
              {listing.title}
            </TruncatedTitle>
          </Box>
          
          <Box display="flex" alignItems="center" mt={0.5}>
            <PriceText variant="h6">
              {typeof listing.price === 'number' 
                ? listing.price.toLocaleString() 
                : listing.price} RSD
            </PriceText>
            <ConditionChip 
              size="small" 
              label={t(`marketplace:condition.${listing.condition}`, { ns: 'marketplace' })}
              condition={listing.condition} 
            />
          </Box>
          
          {listing.location && (
            <LocationText variant="body2">
              <LocationIcon fontSize="small" />
              {listing.location}
            </LocationText>
          )}
        </Box>

        {!compact && (
          <>
            <Divider sx={{ my: 1 }} />
            
            <ActionButtons>
              <Box display="flex" alignItems="center">
                <Button 
                  size="small" 
                  startIcon={<MessageIcon />}
                  variant="outlined"
                  onClick={handleContactClick}
                >
                  {t('card.contact')}
                </Button>
                <DateText>{formatDate(listing.createdAt)}</DateText>
              </Box>
              
              {showMap && listing.latitude && listing.longitude && (
                <Tooltip title={t('card.showOnMap')}>
                  <IconButton
                    size="small"
                    onClick={handleShowOnMap}
                    sx={{ ml: 1 }}
                  >
                    <PlaceIcon fontSize="small" />
                  </IconButton>
                </Tooltip>
              )}
            </ActionButtons>
          </>
        )}
      </ContentContainer>
    </StyledCard>
  );
};

export default GISListingCard;