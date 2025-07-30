'use client';

import React, { useState, useEffect } from 'react';
import Link from 'next/link';
import {
  Heart,
  MapPin,
  Shield,
  Clock,
  Eye,
  MessageCircle,
  ShoppingCart,
  Star,
  Expand,
} from 'lucide-react';
import { formatDistanceToNow } from 'date-fns';
import { ru, enUS, sr } from 'date-fns/locale';
import type { Locale } from 'date-fns';
import type { GridColumns } from '@/components/common/GridColumnsToggle';
import SafeImage from '@/components/SafeImage';
import { useAuth } from '@/contexts/AuthContext';
import { useRouter } from 'next/navigation';
import { useTranslations } from 'next-intl';
import {
  formatAddressWithPrivacy,
  getFullLocalizedAddress,
} from '@/utils/addressUtils';
import { QuickView } from '@/components/ui/QuickView';
import { DistanceBadge } from '@/components/ui/DistanceBadge';
import { DiscountBadge } from '@/components/ui/DiscountBadge';
import { PriceHistoryModal } from '@/components/marketplace/PriceHistoryModal';
import { useDispatch } from 'react-redux';
import { addItem } from '@/store/slices/localCartSlice';
import type { AppDispatch } from '@/store';
import VariantSelectionModal from '@/components/cart/VariantSelectionModal';
import type { UnifiedProduct } from '@/types/unified-product';
import {
  getProductUrl,
  canAddToCart,
  canStartChat,
  getMainImage,
  formatPrice,
  getConditionBadge,
  getStockStatusColor,
  getEcoScore,
  getMockDistance,
  showSecureDealBadge,
  getMinPrice,
  isInStock,
} from '@/utils/product-utils';

interface UnifiedProductCardProps {
  product: UnifiedProduct;
  locale: string;
  viewMode?: 'grid' | 'list';
  gridColumns?: GridColumns;
  onToggleFavorite?: (id: string) => Promise<void>;
}

export const UnifiedProductCard: React.FC<UnifiedProductCardProps> = ({
  product,
  locale,
  viewMode = 'grid',
  gridColumns = 1,
  onToggleFavorite,
}) => {
  const router = useRouter();
  const dispatch = useDispatch<AppDispatch>();
  const [mounted, setMounted] = useState(false);
  const [isFavorite, setIsFavorite] = useState(false);
  const [isProcessingFavorite, setIsProcessingFavorite] = useState(false);
  const [showQuickView, setShowQuickView] = useState(false);
  const [showPriceHistory, setShowPriceHistory] = useState(false);
  const [showVariantModal, setShowVariantModal] = useState(false);
  const t = useTranslations('common');
  const { user, isAuthenticated } = useAuth();

  useEffect(() => {
    setMounted(true);
  }, []);

  // Маппинг локалей для date-fns
  const localeMap: { [key: string]: Locale } = {
    ru: ru,
    en: enUS,
    sr: sr,
  };

  const dateLocale = localeMap[locale] || enUS;

  // Получаем данные из унифицированного продукта
  const mainImageUrl = getMainImage(product);
  const productUrl = getProductUrl(product, locale);
  const conditionBadge = product.condition
    ? getConditionBadge(product.condition)
    : null;
  const ecoScore = getEcoScore(product);
  const distance = getMockDistance(product.id);
  const stockColor = getStockStatusColor(product.stockStatus);
  const minPrice = getMinPrice(product);
  const canBuy = canAddToCart(product, user?.id);
  const canChat = canStartChat(product, user?.id, isAuthenticated);
  const showSecureDeal = showSecureDealBadge(product);
  const inStock = isInStock(product);

  const handleChatClick = async (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();

    if (!isAuthenticated) {
      router.push(
        `/${locale}/auth/login?redirect=${encodeURIComponent(window.location.pathname)}`
      );
      return;
    }

    if (product.seller?.id === user?.id) {
      return;
    }

    router.push(
      `/${locale}/chat?listing_id=${product.id}&seller_id=${product.seller?.id}`
    );
  };

  const handleBuyClick = (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();

    // Для товаров витрин сначала проверяем авторизацию
    if (product.type === 'storefront' && !isAuthenticated) {
      router.push(
        `/${locale}/auth/login?redirect=${encodeURIComponent(window.location.pathname)}`
      );
      return;
    }

    // Проверяем можно ли купить товар (для маркетплейса)
    if (product.type === 'marketplace' && !canBuy) {
      return;
    }

    // Для товаров витрин показываем модальное окно выбора вариантов
    if (product.type === 'storefront') {
      setShowVariantModal(true);
    } else {
      // Для обычных товаров переходим на страницу покупки
      router.push(`/${locale}/marketplace/${product.id}/buy`);
    }
  };

  const handleFavoriteClick = async (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();

    if (!isAuthenticated) {
      router.push(
        `/${locale}/auth/login?redirect=${encodeURIComponent(window.location.pathname)}`
      );
      return;
    }

    if (onToggleFavorite && !isProcessingFavorite) {
      setIsProcessingFavorite(true);
      try {
        await onToggleFavorite(product.id.toString());
        setIsFavorite(!isFavorite);
      } catch (error) {
        console.error('Error toggling favorite:', error);
      } finally {
        setIsProcessingFavorite(false);
      }
    }
  };

  const handleQuickView = (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setShowQuickView(true);
  };

  const handlePriceHistoryClick = (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setShowPriceHistory(true);
  };

  const handleAddToCartWithVariant = (
    variant: any | null,
    quantity: number
  ) => {
    if (!product.storefront) return;

    const mainImage =
      product.images.find((img) => img.isMain) || product.images[0];

    dispatch(
      addItem({
        productId: product.id,
        variantId: variant?.id,
        name: product.name,
        variantName: variant?.name,
        price: variant?.price || minPrice,
        quantity,
        storefrontId: product.storefront.id,
        storefrontName: product.storefront.name,
        storefrontSlug: product.storefront.slug,
        image: mainImage?.publicUrl || '',
        currency: product.currency,
        stockQuantity: variant?.stock_quantity || product.stockQuantity || 100,
      })
    );
  };

  // Формируем адрес для отображения
  const getLocationAddress = () => {
    if (!product.location) return null;

    const fullAddress = getFullLocalizedAddress(
      {
        location: product.location.address,
        city: product.location.city,
        country: product.location.country,
        translations: product.location.translations,
      },
      locale
    );

    return formatAddressWithPrivacy(fullAddress, product.location.privacy);
  };

  if (viewMode === 'list') {
    return (
      <>
        <div
          className={`card bg-base-100 shadow-sm hover:shadow-md transition-shadow duration-200 group ${
            product.type === 'storefront'
              ? 'border-2 border-primary/30 bg-gradient-to-r from-primary/5 to-primary/10'
              : 'border border-base-300 dark:border-base-600'
          }`}
        >
          <Link href={productUrl} className="block">
            <div className="card-body p-4">
              <div className="flex gap-4">
                {/* Изображение слева */}
                <figure className="relative w-32 h-32 flex-shrink-0 bg-base-200 rounded-lg overflow-hidden">
                  <SafeImage
                    src={mainImageUrl}
                    alt={product.name}
                    fill
                    className="object-cover"
                    sizes="128px"
                    fallback={
                      <div className="flex items-center justify-center w-full h-full text-base-content/50">
                        <svg
                          xmlns="http://www.w3.org/2000/svg"
                          fill="none"
                          viewBox="0 0 24 24"
                          strokeWidth={1.5}
                          stroke="currentColor"
                          className="w-12 h-12"
                        >
                          <path
                            strokeLinecap="round"
                            strokeLinejoin="round"
                            d="m2.25 15.75 5.159-5.159a2.25 2.25 0 0 1 3.182 0l5.159 5.159m-1.5-1.5 1.409-1.409a2.25 2.25 0 0 1 3.182 0l2.909 2.909m-18 3.75h16.5a1.5 1.5 0 0 0 1.5-1.5V6a1.5 1.5 0 0 0-1.5-1.5H3.75A1.5 1.5 0 0 0 2.25 6v12a1.5 1.5 0 0 0 1.5 1.5Zm10.5-11.25h.008v.008h-.008V8.25Zm.375 0a.375.375 0 1 1-.75 0 .375.375 0 0 1 .75 0Z"
                          />
                        </svg>
                      </div>
                    }
                  />

                  {/* Badges */}
                  <div className="absolute top-1 left-1 flex flex-col gap-1">
                    {/* Бейдж витрины - показываем первым для товаров витрин */}
                    {product.type === 'storefront' && product.storefront && (
                      <div className="badge badge-xs badge-primary text-primary-content font-medium">
                        🏪 {product.storefront.name}
                      </div>
                    )}
                    {conditionBadge && (
                      <div className={`badge badge-sm ${conditionBadge.class}`}>
                        {t(conditionBadge.text)}
                      </div>
                    )}
                    {ecoScore > 7 && (
                      <div className="badge badge-sm badge-success gap-1">
                        <span className="text-xs">♻️</span>
                        <span>{t('eco.badge')}</span>
                      </div>
                    )}
                    {product.stockStatus === 'out_of_stock' && (
                      <div className="badge badge-sm badge-error">
                        {t('stockStatus.out_of_stock')}
                      </div>
                    )}
                  </div>

                  {/* Количество фото */}
                  {product.images.length > 1 && (
                    <div className="absolute bottom-1 right-1">
                      <div className="badge badge-neutral badge-sm">
                        {product.images.length} {t('photos')}
                      </div>
                    </div>
                  )}
                </figure>

                {/* Информация */}
                <div className="flex-grow">
                  <div className="flex justify-between items-start gap-4">
                    <div className="flex-grow">
                      {/* Категория */}
                      {product.category?.name && (
                        <div className="text-xs text-base-content/60 mb-1">
                          {product.category.name}
                        </div>
                      )}

                      <h2 className="text-lg font-semibold line-clamp-1">
                        {product.name}
                      </h2>

                      {product.description && (
                        <p className="text-sm text-base-content/70 line-clamp-2 mt-1">
                          {product.description}
                        </p>
                      )}

                      {/* Продавец/Витрина и рейтинг */}
                      <div className="flex items-center gap-2 mt-2 text-sm">
                        {product.seller && (
                          <>
                            <div className="avatar placeholder">
                              <div className="bg-neutral text-neutral-content rounded-full w-5">
                                <span className="text-xs">
                                  {product.seller.name[0]}
                                </span>
                              </div>
                            </div>
                            <span className="font-medium">
                              {product.seller.name}
                            </span>
                            {product.seller.verified && (
                              <Shield className="w-3 h-3 text-success" />
                            )}
                            {product.seller.rating && (
                              <div className="flex items-center gap-1">
                                <Star className="w-3 h-3 fill-warning text-warning" />
                                <span className="text-xs">
                                  {product.seller.rating}
                                </span>
                              </div>
                            )}
                          </>
                        )}
                        {product.storefront && (
                          <div className="flex items-center gap-1">
                            <span className="text-xs">🏪</span>
                            <span className="text-primary font-medium">
                              {product.storefront.name}
                            </span>
                            <div className="badge badge-xs badge-primary">
                              {t('storefront')}
                            </div>
                          </div>
                        )}
                      </div>

                      {/* Статистика */}
                      <div className="flex items-center gap-3 mt-2 text-xs text-base-content/60">
                        {product.location && (
                          <div className="flex items-center gap-2">
                            <span className="flex items-center gap-1">
                              <MapPin className="w-3 h-3" />
                              {getLocationAddress()}
                            </span>
                            <DistanceBadge
                              distance={distance}
                              variant="compact"
                            />
                          </div>
                        )}
                        <span className="flex items-center gap-1">
                          <Clock className="w-3 h-3" />
                          {!isNaN(new Date(product.createdAt).getTime())
                            ? formatDistanceToNow(new Date(product.createdAt), {
                                addSuffix: true,
                                locale: dateLocale,
                              })
                            : t('dateUnknown')}
                        </span>
                        {product.viewsCount !== undefined && (
                          <span className="flex items-center gap-1">
                            <Eye className="w-3 h-3" />
                            {product.viewsCount}
                          </span>
                        )}
                      </div>

                      {/* Статус наличия для товаров витрин */}
                      {product.type === 'storefront' && product.stockStatus && (
                        <div className="flex items-center gap-2 mt-2">
                          <span className={`text-sm font-medium ${stockColor}`}>
                            {t(`stockStatus.${product.stockStatus}`)}
                          </span>
                          {product.stockQuantity !== undefined && (
                            <span className="text-sm text-base-content/70">
                              ({product.stockQuantity} {t('units')})
                            </span>
                          )}
                        </div>
                      )}
                    </div>

                    {/* Цена и действия */}
                    <div className="flex flex-col items-end gap-2">
                      <div className="text-right">
                        <div className="flex items-center gap-2 justify-end">
                          {product.hasDiscount && product.oldPrice && (
                            <>
                              <p className="text-sm line-through text-base-content/50">
                                {formatPrice(
                                  product.oldPrice,
                                  product.currency,
                                  locale
                                )}
                              </p>
                              <DiscountBadge
                                oldPrice={product.oldPrice}
                                currentPrice={minPrice}
                                onClick={handlePriceHistoryClick}
                                size="sm"
                              />
                            </>
                          )}
                        </div>
                        <p className="text-xl font-bold text-primary">
                          {product.hasVariants &&
                            product.variants &&
                            product.variants.length > 1 && (
                              <span className="text-sm font-normal">
                                {t('from')}{' '}
                              </span>
                            )}
                          {formatPrice(minPrice, product.currency, locale)}
                        </p>
                        {showSecureDeal && (
                          <div className="text-xs text-success flex items-center gap-1 mt-1">
                            <Shield className="w-3 h-3" />
                            {t('secureDeal')}
                          </div>
                        )}
                      </div>

                      {mounted && (
                        <div className="flex gap-2">
                          <button
                            onClick={handleFavoriteClick}
                            className={`btn btn-circle btn-sm ${isFavorite ? 'btn-error' : 'btn-ghost'}`}
                            disabled={isProcessingFavorite}
                          >
                            <Heart
                              className={`w-4 h-4 ${isFavorite ? 'fill-current' : ''}`}
                            />
                          </button>

                          {/* Кнопка добавления в корзину - обязательна для товаров витрин */}
                          {product.type === 'storefront' && inStock && (
                            <button
                              onClick={handleBuyClick}
                              className="btn btn-primary btn-sm"
                            >
                              <ShoppingCart className="w-4 h-4" />
                              {t('addToCart')}
                            </button>
                          )}
                          {/* Кнопка добавления в корзину для обычных товаров - только для авторизованных */}
                          {product.type === 'marketplace' &&
                            isAuthenticated &&
                            canBuy &&
                            inStock && (
                              <button
                                onClick={handleBuyClick}
                                className="btn btn-primary btn-sm"
                              >
                                <ShoppingCart className="w-4 h-4" />
                                {t('addToCart')}
                              </button>
                            )}
                          {isAuthenticated && canChat && (
                            <button
                              onClick={handleChatClick}
                              className="btn btn-outline btn-sm"
                            >
                              <MessageCircle className="w-4 h-4" />
                              {t('chat')}
                            </button>
                          )}
                        </div>
                      )}
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </Link>

          {/* Quick View button for list mode */}
          {mounted && (
            <div className="absolute top-4 right-4 z-10">
              <button
                onClick={handleQuickView}
                className="btn btn-circle btn-sm btn-ghost bg-base-100/80 backdrop-blur-sm"
                title={t('quickView')}
              >
                <Expand className="w-4 h-4" />
              </button>
            </div>
          )}
        </div>

        {/* Модальные окна */}
        {showQuickView && (
          <QuickView
            isOpen={showQuickView}
            onClose={() => setShowQuickView(false)}
            product={{
              id: product.id.toString(),
              title: product.name,
              price: formatPrice(minPrice, product.currency, locale),
              description: product.description || '',
              images: product.images.map((img) => img.url),
              category: product.category?.name || '',
              seller: product.seller
                ? {
                    name: product.seller.name,
                    rating: product.seller.rating || 4.8,
                    totalReviews: product.seller.totalReviews || 25,
                    avatar: product.seller.pictureUrl || '',
                  }
                : {
                    name: product.storefront?.name || '',
                    rating: 4.8,
                    totalReviews: 25,
                    avatar: '',
                  },
              location: product.location
                ? {
                    address: getLocationAddress() || '',
                    distance: distance,
                  }
                : undefined,
              stats: {
                views: product.viewsCount || 0,
                favorites: product.favoritesCount || 0,
              },
              condition:
                product.condition === 'like_new' ? 'used' : product.condition,
              storefrontId: product.storefront?.id,
              storefrontName: product.storefront?.name,
              storefrontSlug: product.storefront?.slug,
              stockQuantity: product.stockQuantity || 100,
            }}
            onAddToCart={() => setShowVariantModal(true)}
            onContact={() => handleChatClick({} as React.MouseEvent)}
          />
        )}

        {showPriceHistory && product.type === 'marketplace' && (
          <PriceHistoryModal
            listingId={product.id}
            isOpen={showPriceHistory}
            onClose={() => setShowPriceHistory(false)}
          />
        )}

        {showVariantModal && product.storefront && (
          <VariantSelectionModal
            isOpen={showVariantModal}
            onClose={() => setShowVariantModal(false)}
            productId={product.id}
            productName={product.name}
            productImage={product.images[0]?.publicUrl}
            storefrontSlug={product.storefront.slug}
            basePrice={minPrice}
            baseCurrency={product.currency}
            onAddToCart={handleAddToCartWithVariant}
          />
        )}
      </>
    );
  }

  // Grid view (default)
  return (
    <>
      <div
        className={`card card-compact bg-base-100 shadow-sm hover:shadow-md transition-shadow duration-200 group ${
          product.type === 'storefront'
            ? 'border-2 border-primary/30 bg-gradient-to-br from-primary/5 to-primary/10'
            : 'border border-base-300 dark:border-base-600'
        }`}
      >
        {/* Изображение с оверлеями */}
        <figure className="relative aspect-square overflow-hidden">
          <Link href={productUrl}>
            <SafeImage
              src={mainImageUrl}
              alt={product.name}
              fill
              className="object-cover group-hover:scale-105 transition-transform duration-300"
              sizes="(max-width: 768px) 50vw, (max-width: 1200px) 33vw, 25vw"
              fallback={
                <div className="flex items-center justify-center w-full h-full text-base-content/50 bg-base-200">
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    fill="none"
                    viewBox="0 0 24 24"
                    strokeWidth={1.5}
                    stroke="currentColor"
                    className="w-16 h-16"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      d="m2.25 15.75 5.159-5.159a2.25 2.25 0 0 1 3.182 0l5.159 5.159m-1.5-1.5 1.409-1.409a2.25 2.25 0 0 1 3.182 0l2.909 2.909m-18 3.75h16.5a1.5 1.5 0 0 0 1.5-1.5V6a1.5 1.5 0 0 0-1.5-1.5H3.75A1.5 1.5 0 0 0 2.25 6v12a1.5 1.5 0 0 0 1.5 1.5Zm10.5-11.25h.008v.008h-.008V8.25Zm.375 0a.375.375 0 1 1-.75 0 .375.375 0 0 1 .75 0Z"
                    />
                  </svg>
                </div>
              }
            />
          </Link>

          {/* Badges слева сверху */}
          <div className="absolute top-2 left-2 flex flex-col gap-1">
            {/* Бейдж витрины - показываем первым для товаров витрин */}
            {product.type === 'storefront' && product.storefront && (
              <div className="badge badge-sm badge-primary text-primary-content font-medium">
                🏪 {product.storefront.name}
              </div>
            )}
            {conditionBadge && (
              <div className={`badge badge-sm ${conditionBadge.class}`}>
                {t(conditionBadge.text)}
              </div>
            )}
            {ecoScore > 7 && (
              <div className="badge badge-sm badge-success gap-1">
                <span className="text-xs">♻️</span>
                <span>{t('eco.badge')}</span>
              </div>
            )}
            {product.stockStatus === 'out_of_stock' && (
              <div className="badge badge-sm badge-error">
                {t('stockStatus.out_of_stock')}
              </div>
            )}
          </div>

          {/* Действия справа сверху */}
          {mounted && (
            <div className="absolute top-2 right-2 flex gap-1">
              <button
                onClick={handleQuickView}
                className="btn btn-circle btn-sm btn-ghost bg-base-100/80 backdrop-blur-sm"
                title={t('quickView')}
              >
                <Expand className="w-4 h-4" />
              </button>
              <button
                onClick={handleFavoriteClick}
                className="btn btn-circle btn-sm btn-ghost bg-base-100/80 backdrop-blur-sm"
                disabled={isProcessingFavorite}
              >
                <Heart
                  className={`w-4 h-4 ${isFavorite ? 'fill-error text-error' : ''}`}
                />
              </button>
            </div>
          )}

          {/* Количество фото */}
          {product.images.length > 1 && (
            <div className="absolute bottom-2 right-2">
              <div className="badge badge-neutral badge-sm">
                {product.images.length} {t('photos')}
              </div>
            </div>
          )}

          {/* Out of stock overlay для товаров витрин */}
          {product.type === 'storefront' &&
            product.stockStatus === 'out_of_stock' && (
              <div className="absolute inset-0 bg-black/50 flex items-center justify-center">
                <span className="badge badge-error badge-lg">
                  {t('stockStatus.out_of_stock')}
                </span>
              </div>
            )}
        </figure>

        <div className={`card-body ${gridColumns === 3 ? 'p-2' : 'p-3'}`}>
          {/* Категория */}
          {product.category?.name && (
            <div className="text-xs text-base-content/60">
              {product.category.name}
            </div>
          )}

          {/* Заголовок */}
          <h3
            className={`card-title ${
              gridColumns === 3 ? 'text-xs' : 'text-sm'
            } line-clamp-2 min-h-[2.5rem]`}
          >
            {product.name}
          </h3>

          {/* Продавец/Витрина - скрываем при 3 столбцах */}
          {gridColumns < 3 && (
            <div className="flex items-center gap-2 text-xs">
              {product.seller && (
                <>
                  <div className="avatar placeholder">
                    <div className="bg-neutral text-neutral-content rounded-full w-5">
                      <span className="text-xs">{product.seller.name[0]}</span>
                    </div>
                  </div>
                  <span className="font-medium">{product.seller.name}</span>
                  {product.seller.rating && (
                    <div className="flex items-center gap-1">
                      <Star className="w-3 h-3 fill-warning text-warning" />
                      <span className="text-xs">{product.seller.rating}</span>
                    </div>
                  )}
                </>
              )}
              {product.storefront && (
                <div className="flex items-center gap-1">
                  <span className="text-xs">🏪</span>
                  <span className="text-primary font-medium text-xs">
                    {product.storefront.name}
                  </span>
                  <div className="badge badge-xs badge-primary">витрина</div>
                </div>
              )}
            </div>
          )}

          {/* Локация и статистика - упрощаем для маленьких карточек */}
          <div className="flex flex-col gap-1 text-xs text-base-content/60">
            {product.location && gridColumns < 3 && (
              <div className="flex items-center gap-2">
                <span className="flex items-center gap-1">
                  <MapPin className="w-3 h-3" />
                  {getLocationAddress()}
                </span>
                <DistanceBadge distance={distance} variant="compact" />
              </div>
            )}
            {gridColumns === 1 && (
              <div className="flex items-center gap-3">
                <span className="flex items-center gap-1">
                  <Clock className="w-3 h-3" />
                  {(() => {
                    const date = new Date(product.createdAt);
                    if (isNaN(date.getTime())) {
                      return '-';
                    }
                    return formatDistanceToNow(date, {
                      addSuffix: true,
                      locale: dateLocale,
                    });
                  })()}
                </span>
                {product.viewsCount !== undefined && (
                  <span className="flex items-center gap-1">
                    <Eye className="w-3 h-3" />
                    {product.viewsCount}
                  </span>
                )}
              </div>
            )}
          </div>

          {/* Статус наличия для товаров витрин */}
          {product.type === 'storefront' &&
            product.stockStatus &&
            gridColumns < 3 && (
              <div className="flex items-center gap-2">
                <span className={`text-xs font-medium ${stockColor}`}>
                  {t(`stockStatus.${product.stockStatus}`)}
                </span>
                {product.stockQuantity !== undefined && (
                  <span className="text-xs text-base-content/70">
                    ({product.stockQuantity})
                  </span>
                )}
              </div>
            )}

          {/* Цена и действия */}
          <div className="card-actions justify-between items-end mt-2">
            <div>
              {product.hasDiscount && product.oldPrice && (
                <div className="flex items-center gap-2 mb-1">
                  <p className="text-xs line-through text-base-content/50">
                    {formatPrice(product.oldPrice, product.currency, locale)}
                  </p>
                  <DiscountBadge
                    oldPrice={product.oldPrice}
                    currentPrice={minPrice}
                    onClick={handlePriceHistoryClick}
                    size="sm"
                  />
                </div>
              )}
              <div
                className={`${
                  gridColumns === 3 ? 'text-base' : 'text-lg'
                } font-bold`}
              >
                {product.hasVariants &&
                  product.variants &&
                  product.variants.length > 1 && (
                    <span className="text-xs font-normal">{t('from')} </span>
                  )}
                {formatPrice(minPrice, product.currency, locale)}
              </div>
              {showSecureDeal && (
                <div className="text-xs text-success flex items-center gap-1">
                  <Shield className="w-3 h-3" />
                  {t('secureDeal')}
                </div>
              )}
            </div>

            {mounted && gridColumns < 3 && (
              <div className="flex gap-1">
                {/* Кнопка добавления в корзину - обязательна для товаров витрин */}
                {product.type === 'storefront' && inStock && (
                  <button
                    onClick={handleBuyClick}
                    className="btn btn-primary btn-xs"
                    title={t('addToCart')}
                  >
                    <ShoppingCart className="w-3 h-3" />
                  </button>
                )}
                {/* Кнопка добавления в корзину для обычных товаров - только для авторизованных */}
                {product.type === 'marketplace' &&
                  isAuthenticated &&
                  canBuy &&
                  inStock && (
                    <button
                      onClick={handleBuyClick}
                      className="btn btn-primary btn-xs"
                      title={t('addToCart')}
                    >
                      <ShoppingCart className="w-3 h-3" />
                    </button>
                  )}
                <button
                  onClick={handleQuickView}
                  className="btn btn-outline btn-xs"
                  title={t('quickView')}
                >
                  <Expand className="w-3 h-3" />
                </button>
                {canChat && (
                  <button
                    onClick={handleChatClick}
                    className="btn btn-outline btn-xs"
                  >
                    <MessageCircle className="w-3 h-3" />
                  </button>
                )}
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Модальные окна */}
      {showQuickView && (
        <QuickView
          isOpen={showQuickView}
          onClose={() => setShowQuickView(false)}
          product={{
            id: product.id.toString(),
            title: product.name,
            price: formatPrice(minPrice, product.currency, locale),
            description: product.description || '',
            images: product.images.map((img) => img.url),
            category: product.category?.name || '',
            seller: product.seller
              ? {
                  name: product.seller.name,
                  rating: product.seller.rating || 4.8,
                  totalReviews: product.seller.totalReviews || 25,
                  avatar: product.seller.pictureUrl || '',
                }
              : {
                  name: product.storefront?.name || '',
                  rating: 4.8,
                  totalReviews: 25,
                  avatar: '',
                },
            location: product.location
              ? {
                  address: getLocationAddress() || '',
                  distance: distance,
                }
              : undefined,
            stats: {
              views: product.viewsCount || 0,
              favorites: product.favoritesCount || 0,
            },
            condition:
              product.condition === 'like_new' ? 'used' : product.condition,
            storefrontId: product.storefront?.id,
            storefrontName: product.storefront?.name,
            storefrontSlug: product.storefront?.slug,
            stockQuantity: product.stockQuantity || 100,
          }}
          onAddToCart={() => setShowVariantModal(true)}
          onContact={() => handleChatClick({} as React.MouseEvent)}
        />
      )}

      {showPriceHistory && product.type === 'marketplace' && (
        <PriceHistoryModal
          listingId={product.id}
          isOpen={showPriceHistory}
          onClose={() => setShowPriceHistory(false)}
        />
      )}

      {showVariantModal && product.storefront && (
        <VariantSelectionModal
          isOpen={showVariantModal}
          onClose={() => setShowVariantModal(false)}
          productId={product.id}
          productName={product.name}
          productImage={product.images[0]?.publicUrl}
          storefrontSlug={product.storefront.slug}
          basePrice={minPrice}
          baseCurrency={product.currency}
          onAddToCart={handleAddToCartWithVariant}
        />
      )}
    </>
  );
};
