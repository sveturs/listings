'use client';

import React, { useState, useEffect, useCallback, MouseEvent } from 'react';
import { useTranslations } from 'next-intl';
import { useRouter } from 'next/navigation';
import { 
  MapPin as LocationIcon, 
  Clock as AccessTime, 
  Camera, 
  Store, 
  Eye 
} from 'lucide-react';
import {
  Card,
  CardContent,
  CardMedia,
  Typography,
  Box,
  Chip,
  Button,
  Rating,
  Modal,
} from '@mui/material';
import { apiClient } from '@/lib/api-client';
import { Listing, ListingImage } from '@/types/listing';

interface ListingCardProps {
  listing: Listing;
  isMobile?: boolean;
  onClick?: (listing: Listing) => void;
  showStatus?: boolean;
  viewMode?: 'grid' | 'list';
}

const ListingCard: React.FC<ListingCardProps> = ({ 
  listing, 
  isMobile = false, 
  onClick, 
  showStatus = false,
  viewMode = 'grid'
}) => {
  const t = useTranslations('marketplace');
  const router = useRouter();
  const [, setStoreName] = useState<string>(t('listings.details.seller.title'));
  const [isPriceHistoryOpen, setIsPriceHistoryOpen] = useState<boolean>(false);

  // Check if listing is part of a storefront
  const isStoreItem = useCallback((): boolean => {
    return listing.storefront_id !== undefined && listing.storefront_id !== null;
  }, [listing.storefront_id]);

  // Get store ID
  const getStoreId = useCallback((): number | string | null => {
    return listing.storefront_id || null;
  }, [listing.storefront_id]);

  // Get seller name
  useEffect(() => {
    if (listing && (listing.storefront_name || listing.storefrontName)) {
      setStoreName(listing.storefront_name || listing.storefrontName || '');
      return;
    }

    if (isStoreItem()) {
      apiClient.get<{ data: { name?: string } }>(`/api/v1/public/storefronts/${getStoreId()}`)
        .then(response => {
          if (response.data && response.data.data) {
            const name = response.data.data.name || t('listings.details.seller.title');
            setStoreName(name);
          }
        })
        .catch(error => {
          console.error('Error fetching storefront:', error);
          setStoreName(t('listings.details.seller.title'));
        });
    }
  }, [listing, t, isStoreItem, getStoreId, setStoreName]);

  const handleCardClick = (e: MouseEvent<HTMLDivElement>) => {
    if ((e.target as HTMLElement).closest('button')) {
      return;
    }

    if (onClick) {
      onClick(listing);
    } else {
      router.push(`/marketplace/listings/${listing.id}`);
    }
  };

  const getDiscountData = () => {
    // Check for discount in multiple formats
    if (listing.has_discount && listing.old_price && listing.price &&
        listing.old_price > listing.price) {
      const percent = Math.round(((listing.old_price - listing.price) / listing.old_price) * 100);
      return {
        percent,
        oldPrice: listing.old_price,
        hasPriceHistory: false
      };
    } else if (listing.old_price && listing.price && listing.old_price > listing.price) {
      const percent = Math.round(((listing.old_price - listing.price) / listing.old_price) * 100);
      return {
        percent,
        oldPrice: listing.old_price,
        hasPriceHistory: false
      };
    } else if (listing.metadata?.discount) {
      const discount = listing.metadata.discount;
      const oldPrice = discount.previous_price || discount.oldPrice || 0;
      const percent = discount.discount_percent || discount.percent || 0;

      if (oldPrice > 0 && percent > 0) {
        return {
          percent,
          oldPrice,
          hasPriceHistory: Boolean(discount.has_price_history)
        };
      }

      if (oldPrice > 0 && listing.price && oldPrice > listing.price) {
        const calculatedPercent = Math.round(((oldPrice - listing.price) / oldPrice) * 100);
        return {
          percent: calculatedPercent,
          oldPrice,
          hasPriceHistory: Boolean(discount.has_price_history)
        };
      }
    }

    return null;
  };

  const discountData = getDiscountData();

  const formatDate = (dateString?: string) => {
    if (!dateString) return '';
    try {
      const date = new Date(dateString);
      const locale = typeof window !== 'undefined' ? window.location.pathname.split('/')[1] : 'en';
      return date.toLocaleDateString(locale);
    } catch {
      return dateString;
    }
  };

  // Get first image URL
  const getImageUrl = () => {
    if (!listing.images || !Array.isArray(listing.images) || listing.images.length === 0) {
      return '/placeholder-listing.jpg';
    }

    const firstImage = listing.images[0] as string | ListingImage;
    const baseUrl = process.env.NEXT_PUBLIC_MINIO_URL || process.env.NEXT_PUBLIC_BACKEND_URL || '';

    if (typeof firstImage === 'string') {
      if (firstImage.startsWith('http')) {
        return firstImage;
      } else {
        return `${baseUrl}/uploads/${firstImage}`;
      }
    }

    if (firstImage && typeof firstImage === 'object') {
      if (firstImage.public_url) {
        if (firstImage.public_url.startsWith('http')) {
          return firstImage.public_url;
        } else {
          return `${baseUrl}${firstImage.public_url}`;
        }
      }

      if (firstImage.storage_type === 'minio' ||
          (firstImage.file_path && firstImage.file_path.includes('listings/'))) {
        return `${baseUrl}${firstImage.public_url}`;
      }

      if (firstImage.file_path) {
        return `${baseUrl}/uploads/${firstImage.file_path}`;
      }

      if (firstImage.url) {
        return firstImage.url;
      }
    }

    return '/placeholder-listing.jpg';
  };

  const handleOpenPriceHistory = (e: MouseEvent<HTMLButtonElement>) => {
    e.stopPropagation();
    setIsPriceHistoryOpen(true);
  };

  const handleClosePriceHistory = () => {
    setIsPriceHistoryOpen(false);
  };

  // List view layout
  if (viewMode === 'list') {
    return (
      <>
        <Card
          sx={{
            width: '100%',
            display: 'flex',
            transition: 'box-shadow 0.2s',
            cursor: 'pointer',
            '&:hover': {
              boxShadow: 3
            },
            mb: 2
          }}
          onClick={handleCardClick}
        >
          <Box sx={{ position: 'relative', width: isMobile ? 100 : 200, flexShrink: 0 }}>
            <CardMedia
              component="img"
              sx={{
                width: '100%',
                height: isMobile ? 100 : 150,
                objectFit: 'cover'
              }}
              image={getImageUrl()}
              alt={listing.title}
            />
            {/* Store badge */}
            {isStoreItem() && (
              <Box
                sx={{
                  position: 'absolute',
                  top: 5,
                  left: 5,
                  bgcolor: 'primary.main',
                  color: 'white',
                  borderRadius: '4px',
                  padding: '2px 6px',
                  fontSize: '0.65rem',
                }}
              >
                <Store size={12} />
              </Box>
            )}
          </Box>
          
          <CardContent sx={{ flex: 1, py: 2 }}>
            <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start' }}>
              <Box sx={{ flex: 1 }}>
                <Typography variant="h6" component="h2" sx={{ mb: 1, fontSize: '1rem' }}>
                  {listing.title}
                </Typography>
                
                <Box sx={{ display: 'flex', gap: 2, mb: 1 }}>
                  <Box sx={{ display: 'flex', alignItems: 'center' }}>
                    <LocationIcon size={16} style={{ marginRight: 4 }} />
                    <Typography variant="body2" color="text.secondary">
                      {listing.city || listing.location || t('listings.locationNotSpecified')}
                    </Typography>
                  </Box>
                  
                  <Box sx={{ display: 'flex', alignItems: 'center' }}>
                    <AccessTime size={16} style={{ marginRight: 4 }} />
                    <Typography variant="body2" color="text.secondary">
                      {formatDate(listing.created_at || listing.createdAt)}
                    </Typography>
                  </Box>
                  
                  {(listing.views_count !== undefined || listing.viewCount !== undefined) && (
                    <Box sx={{ display: 'flex', alignItems: 'center' }}>
                      <Eye size={16} style={{ marginRight: 4 }} />
                      <Typography variant="body2" color="text.secondary">
                        {listing.views_count || listing.viewCount}
                      </Typography>
                    </Box>
                  )}
                </Box>
              </Box>
              
              <Box sx={{ textAlign: 'right' }}>
                <Typography variant="h6" component="div" color="primary" fontWeight="bold">
                  {new Intl.NumberFormat('sr-RS', {
                    style: 'currency',
                    currency: listing.currency || 'RSD',
                    maximumFractionDigits: 0
                  }).format(listing.price)}
                </Typography>
                
                {discountData && (
                  <>
                    <Typography
                      variant="body2"
                      color="text.secondary"
                      sx={{ textDecoration: 'line-through' }}
                    >
                      {new Intl.NumberFormat('sr-RS', {
                        style: 'currency',
                        currency: listing.currency || 'RSD',
                        maximumFractionDigits: 0
                      }).format(discountData.oldPrice)}
                    </Typography>
                    <Chip
                      label={`-${discountData.percent}%`}
                      size="small"
                      color="error"
                      sx={{ mt: 0.5 }}
                    />
                  </>
                )}
              </Box>
            </Box>
          </CardContent>
        </Card>
        
        {/* Price history modal */}
        <Modal
          open={isPriceHistoryOpen}
          onClose={handleClosePriceHistory}
          aria-labelledby="price-history-modal"
          aria-describedby="price-history-modal-description"
        >
          <Box sx={{
            position: 'absolute',
            top: '50%',
            left: '50%',
            transform: 'translate(-50%, -50%)',
            width: 400,
            bgcolor: 'background.paper',
            border: '2px solid #000',
            boxShadow: 24,
            p: 4,
          }}>
            <Typography id="price-history-modal" variant="h6" component="h2">
              {t('listings.priceHistory.title')}
            </Typography>
            <Typography id="price-history-modal-description" sx={{ mt: 2 }}>
              {/* Price history chart component would go here */}
              Price history chart placeholder
            </Typography>
          </Box>
        </Modal>
      </>
    );
  }

  // Grid view layout (original)
  return (
    <>
      <Card
        sx={{
          height: '100%',
          display: 'flex',
          flexDirection: 'column',
          transition: 'transform 0.2s, box-shadow 0.2s',
          cursor: 'pointer',
          '&:hover': {
            transform: 'translateY(-4px)',
            boxShadow: 3
          }
        }}
        onClick={handleCardClick}
      >
        <Box sx={{ position: 'relative' }}>
          <CardMedia
            component="img"
            sx={{
              height: isMobile ? 120 : 200,
              objectFit: 'cover'
            }}
            image={getImageUrl()}
            alt={listing.title}
          />

          <Box>
            {/* Store badge */}
            {isStoreItem() && (
              <Box
                sx={{
                  position: 'absolute',
                  top: 10,
                  left: 10,
                  bgcolor: 'primary.main',
                  color: 'white',
                  borderRadius: '4px',
                  padding: '4px 8px',
                  fontWeight: 'bold',
                  display: 'flex',
                  alignItems: 'center',
                  gap: 0.5,
                  fontSize: '0.75rem',
                  cursor: 'pointer',
                  zIndex: 5
                }}
                onClick={(e) => {
                  e.stopPropagation();
                  if (e.nativeEvent) {
                    e.nativeEvent.stopImmediatePropagation();
                  }
                  router.push(`/shop/${getStoreId()}`);
                  return false;
                }}
              >
                <Store size={16} />
                {t('listings.details.goToStore')}
              </Box>
            )}

            {/* Discount badge */}
            {discountData && discountData.percent > 0 && (
              <Box
                sx={{
                  position: 'absolute',
                  top: 10,
                  right: 10,
                  bgcolor: 'error.main',
                  color: 'white',
                  borderRadius: '4px',
                  padding: '4px 8px',
                  fontWeight: 'bold',
                  display: 'flex',
                  alignItems: 'center',
                  gap: 0.5,
                  zIndex: 4
                }}
              >
                <Typography variant="body2" fontWeight="bold">
                  -{discountData.percent}%
                </Typography>
              </Box>
            )}
          </Box>
          
          {/* Status tag */}
          {showStatus && listing.status && (
            <Chip
              label={
                listing.status === 'active'
                  ? t('listings.status.active')
                  : t('listings.status.inactive')
              }
              size="small"
              color={listing.status === 'active' ? 'success' : 'default'}
              sx={{
                position: 'absolute',
                top: 10,
                left: 10
              }}
            />
          )}
          
          {/* Photo count indicator */}
          {listing.images && listing.images.length > 1 && (
            <Box
              sx={{
                position: 'absolute',
                bottom: 10,
                right: 10,
                bgcolor: 'rgba(0,0,0,0.6)',
                color: 'white',
                borderRadius: '4px',
                padding: '2px 6px',
                display: 'flex',
                alignItems: 'center',
                gap: 0.5
              }}
            >
              <Camera size={16} />
              <Typography variant="body2" fontWeight="medium">
                {listing.images.length}
              </Typography>
            </Box>
          )}
        </Box>
        
        <CardContent sx={{ flexGrow: 1, display: 'flex', flexDirection: 'column' }}>
          <Typography
            variant="h6"
            component="h2"
            sx={{
              mb: 1,
              fontSize: isMobile ? '0.9rem' : '1rem',
              fontWeight: 'medium',
              overflow: 'hidden',
              textOverflow: 'ellipsis',
              display: '-webkit-box',
              WebkitLineClamp: 2,
              WebkitBoxOrient: 'vertical'
            }}
          >
            {listing.title}
          </Typography>
          
          <Box sx={{ display: 'flex', alignItems: 'center', mb: 1 }}>
            <Typography
              variant="h6"
              component="div"
              color="primary"
              fontWeight="bold"
              sx={{ fontSize: isMobile ? '1rem' : '1.25rem' }}
            >
              {new Intl.NumberFormat('sr-RS', {
                style: 'currency',
                currency: listing.currency || 'RSD',
                maximumFractionDigits: 0
              }).format(listing.price)}
            </Typography>
            
            {discountData && (
              <Typography
                variant="body2"
                color="text.secondary"
                sx={{
                  ml: 1,
                  textDecoration: 'line-through',
                  fontSize: isMobile ? '0.75rem' : '0.875rem'
                }}
              >
                {new Intl.NumberFormat('sr-RS', {
                  style: 'currency',
                  currency: listing.currency || 'RSD',
                  maximumFractionDigits: 0
                }).format(discountData.oldPrice)}
              </Typography>
            )}
            
            {discountData?.hasPriceHistory && (
              <Button
                size="small"
                variant="text"
                onClick={handleOpenPriceHistory}
                sx={{ minWidth: 'auto', ml: 'auto', px: 1 }}
              >
                <AccessTime size={16} />
              </Button>
            )}
          </Box>
          
          <Box
            sx={{
              display: 'flex',
              alignItems: 'center',
              mb: 0.5,
              fontSize: isMobile ? '0.75rem' : '0.875rem',
              color: 'text.secondary'
            }}
          >
            <LocationIcon size={isMobile ? 14 : 16} style={{ marginRight: 4 }} />
            <Typography
              variant="body2"
              sx={{
                overflow: 'hidden',
                textOverflow: 'ellipsis',
                whiteSpace: 'nowrap'
              }}
            >
              {listing.city || listing.location || t('listings.locationNotSpecified')}
            </Typography>
          </Box>
          
          <Box
            sx={{
              display: 'flex',
              alignItems: 'center',
              mb: 0.5,
              fontSize: isMobile ? '0.75rem' : '0.875rem',
              color: 'text.secondary',
              mt: 'auto'
            }}
          >
            <AccessTime size={isMobile ? 14 : 16} style={{ marginRight: 4 }} />
            <Typography variant="body2">
              {formatDate(listing.created_at || listing.createdAt)}
            </Typography>
            
            {(listing.views_count !== undefined || listing.viewCount !== undefined) && (
              <Box
                sx={{
                  display: 'flex',
                  alignItems: 'center',
                  ml: 'auto',
                  fontSize: isMobile ? '0.75rem' : '0.875rem',
                  color: 'text.secondary'
                }}
              >
                <Eye size={isMobile ? 14 : 16} style={{ marginRight: 4 }} />
                <Typography variant="body2">
                  {listing.views_count || listing.viewCount}
                </Typography>
              </Box>
            )}
          </Box>
          
          {(listing.rating !== undefined || listing.average_rating !== undefined) && 
           (listing.reviews_count !== undefined || listing.review_count !== undefined) && (
            <Box
              sx={{
                display: 'flex',
                alignItems: 'center',
                mt: 1
              }}
            >
              <Rating
                value={listing.rating || listing.average_rating}
                readOnly
                precision={0.5}
                size={isMobile ? 'small' : 'medium'}
              />
              <Typography variant="body2" color="text.secondary" sx={{ ml: 1 }}>
                ({listing.reviews_count || listing.review_count})
              </Typography>
            </Box>
          )}
        </CardContent>
      </Card>
      
      {/* Price history modal */}
      <Modal
        open={isPriceHistoryOpen}
        onClose={handleClosePriceHistory}
        aria-labelledby="price-history-modal"
        aria-describedby="price-history-modal-description"
      >
        <Box sx={{
          position: 'absolute',
          top: '50%',
          left: '50%',
          transform: 'translate(-50%, -50%)',
          width: 400,
          bgcolor: 'background.paper',
          border: '2px solid #000',
          boxShadow: 24,
          p: 4,
        }}>
          <Typography id="price-history-modal" variant="h6" component="h2">
            {t('listings.priceHistory.title')}
          </Typography>
          <Typography id="price-history-modal-description" sx={{ mt: 2 }}>
            {/* Price history chart component would go here */}
            Price history chart placeholder
          </Typography>
        </Box>
      </Modal>
    </>
  );
};

export default ListingCard;