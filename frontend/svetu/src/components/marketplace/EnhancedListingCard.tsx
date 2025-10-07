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
import { MarketplaceItem, MarketplaceImage } from '@/types/marketplace';
import type { GridColumns } from '@/components/common/GridColumnsToggle';
import SafeImage from '@/components/SafeImage';
import { useAuth } from '@/contexts/AuthContext';
import { useRouter } from 'next/navigation';
import { useTranslations } from 'next-intl';
import configManager from '@/config';
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
import type { components } from '@/types/generated/api';

type ProductVariant =
  components['schemas']['models.StorefrontProductVariant'];

interface EnhancedListingCardProps {
  item: MarketplaceItem;
  locale: string;
  viewMode?: 'grid' | 'list';
  gridColumns?: GridColumns;
  onToggleFavorite?: (id: string) => Promise<void>;
}

export const EnhancedListingCard: React.FC<EnhancedListingCardProps> = ({
  item,
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

  // Попытаемся найти condition в атрибутах, если его нет в основном объекте
  let itemCondition = item.condition;
  if (!itemCondition && item.attributes) {
    const conditionAttr = item.attributes.find(
      (attr: any) =>
        attr.attribute_name === 'condition' || attr.name === 'condition'
    );
    if (conditionAttr) {
      itemCondition = String(
        conditionAttr.value || conditionAttr.text_value || ''
      );
    }
  }
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

  const conditionBadge: Record<string, { text: string; class: string }> = {
    new: { text: t('condition.new'), class: 'badge-success' },
    like_new: { text: t('condition.likeNew'), class: 'badge-info' },
    used: { text: t('condition.used'), class: 'badge-primary' },
    refurbished: { text: t('condition.refurbished'), class: 'badge-warning' },
  };

  const formatPrice = (price?: number, currency?: string) => {
    if (!price) return '';

    const formatter = new Intl.NumberFormat(locale, {
      style: 'currency',
      currency: currency || 'RSD',
      minimumFractionDigits: 0,
      maximumFractionDigits: 0,
    });

    return formatter.format(price);
  };

  const getImageUrl = (image?: MarketplaceImage) => {
    if (!image) return null;
    return configManager.buildImageUrl(image.public_url);
  };

  const mainImage = item.images?.find((img) => img.is_main) || item.images?.[0];
  const imageUrl = getImageUrl(mainImage);

  const handleChatClick = async (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();

    if (!isAuthenticated) {
      const redirectPath =
        typeof window !== 'undefined' ? window.location.pathname : '/';
      router.push(
        `/${locale}/auth/login?redirect=${encodeURIComponent(redirectPath)}`
      );
      return;
    }

    if (item.user_id === user?.id) {
      return;
    }

    router.push(
      `/${locale}/chat?listing_id=${item.id}&seller_id=${item.user_id}`
    );
  };

  const handleBuyClick = (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();

    if (!isAuthenticated) {
      const redirectPath =
        typeof window !== 'undefined' ? window.location.pathname : '/';
      router.push(
        `/${locale}/auth/login?redirect=${encodeURIComponent(redirectPath)}`
      );
      return;
    }

    if (item.user_id === user?.id) {
      return;
    }

    // Для товаров витрин добавляем в корзину
    if (isOnlinePurchaseAvailable()) {
      handleAddToCart();
    } else {
      // Для обычных товаров переходим на страницу покупки (можно убрать или оставить для legacy)
      router.push(`/${locale}/marketplace/${item.id}/buy`);
    }
  };

  const handleFavoriteClick = async (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();

    if (!isAuthenticated) {
      const redirectPath =
        typeof window !== 'undefined' ? window.location.pathname : '/';
      router.push(
        `/${locale}/auth/login?redirect=${encodeURIComponent(redirectPath)}`
      );
      return;
    }

    if (onToggleFavorite && !isProcessingFavorite) {
      setIsProcessingFavorite(true);
      try {
        await onToggleFavorite(item.id.toString());
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

  const isOnlinePurchaseAvailable = () => {
    // Онлайн покупка доступна если у товара есть витрина
    return item.storefront_id && item.storefront_id > 0;
  };

  const handleAddToCart = () => {
    if (!item.storefront_id) return;

    // Show variant selection modal
    setShowVariantModal(true);
  };

  const handleAddToCartWithVariant = (
    variant: ProductVariant | null,
    quantity: number
  ) => {
    if (!item.storefront_id) return;

    dispatch(
      addItem({
        productId: item.id,
        variantId: variant?.id,
        name: item.title,
        variantName: variant?.sku,
        price: variant?.price || item.price || 0,
        quantity,
        storefrontId: item.storefront_id,
        storefrontName: item.storefront?.name || 'Store',
        storefrontSlug: item.storefront?.slug || 'store',
        image: mainImage?.public_url || '',
        currency: 'RSD',
        stockQuantity: variant?.stock_quantity || 100, // TODO: get from API
      })
    );
  };

  const getItemUrl = () => {
    // Настоящие товары витрин (из таблицы storefront_products) имеют тип 'storefront'
    // и для них используется специальный URL с slug витрины
    if (item.product_type === 'storefront' && item.storefront_slug) {
      return `/${locale}/storefronts/${item.storefront_slug}/products/${item.id}`;
    }
    // Все объявления маркетплейса открываются через marketplace URL,
    // даже если они связаны с витриной
    return `/${locale}/marketplace/${item.id}`;
  };

  // Рассчитываем эко-скор для б/у товаров
  const ecoScore =
    item.condition === 'used' || item.condition === 'refurbished' ? 8 : 0;

  // TODO: Получать реальное расстояние из API
  // Временное решение - генерируем расстояние на основе ID для консистентности
  const getDistance = () => {
    const hash = item.id
      .toString()
      .split('')
      .reduce((acc, char) => acc + char.charCodeAt(0), 0);
    return (hash % 20) + 0.5; // От 0.5 до 20.5 км
  };

  if (viewMode === 'list') {
    return (
      <>
        <div
          className={`card bg-base-100 shadow-sm hover:shadow-md border ${
            item.storefront_id
              ? 'border-primary/50 ring-1 ring-primary/20'
              : 'border-base-300 dark:border-base-600'
          } transition-shadow duration-200 group`}
        >
          <Link href={getItemUrl()} className="block">
            <div className="card-body p-4">
              <div className="flex gap-4">
                {/* Изображение слева */}
                <figure className="relative w-32 h-32 flex-shrink-0 bg-base-200 rounded-lg overflow-hidden">
                  <SafeImage
                    src={imageUrl}
                    alt={item.title}
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
                    {itemCondition && conditionBadge[itemCondition] && (
                      <div
                        className={`badge badge-sm ${conditionBadge[itemCondition].class}`}
                      >
                        {conditionBadge[itemCondition].text}
                      </div>
                    )}
                    {ecoScore > 7 && (
                      <div className="badge badge-sm badge-success gap-1">
                        <span className="text-xs">♻️</span>
                        <span>{t('eco.badge')}</span>
                      </div>
                    )}
                  </div>

                  {/* Количество фото */}
                  {item.images && item.images.length > 1 && (
                    <div className="absolute bottom-1 right-1">
                      <div className="badge badge-neutral badge-sm">
                        {item.images.length} {t('photos')}
                      </div>
                    </div>
                  )}
                </figure>

                {/* Информация */}
                <div className="flex-grow">
                  <div className="flex justify-between items-start gap-4">
                    <div className="flex-grow">
                      {/* Категория */}
                      {item.category?.name && (
                        <div className="text-xs text-base-content/60 mb-1">
                          {item.category.translations &&
                          item.category.translations[locale]
                            ? item.category.translations[locale]
                            : item.category.name}
                        </div>
                      )}

                      <h2 className="text-lg font-semibold line-clamp-1">
                        {item.title}
                      </h2>

                      {item.description && (
                        <p className="text-sm text-base-content/70 line-clamp-2 mt-1">
                          {item.description}
                        </p>
                      )}

                      {/* Продавец и рейтинг */}
                      {item.user?.name && (
                        <div className="flex items-center gap-2 mt-2 text-sm">
                          <div className="avatar placeholder">
                            <div className="bg-neutral text-neutral-content rounded-full w-5">
                              <span className="text-xs">
                                {item.user.name[0]}
                              </span>
                            </div>
                          </div>
                          <span className="font-medium">{item.user.name}</span>
                          {/* TODO: добавить поле verified в API */}
                          {/* {item.seller?.verified && (
                          <Shield className="w-3 h-3 text-success" />
                        )} */}
                          {/* TODO: добавить рейтинг продавца */}
                          <div className="flex items-center gap-1">
                            <Star className="w-3 h-3 fill-warning text-warning" />
                            <span className="text-xs">4.8</span>
                          </div>
                        </div>
                      )}

                      {/* Статистика */}
                      <div className="flex items-center gap-3 mt-2 text-xs text-base-content/60">
                        {item.location && (
                          <div className="flex items-center gap-2">
                            <span className="flex items-center gap-1">
                              <MapPin className="w-3 h-3" />
                              {formatAddressWithPrivacy(
                                getFullLocalizedAddress(item, locale),
                                item.location_privacy
                              )}
                            </span>
                            <DistanceBadge
                              distance={getDistance()}
                              variant="compact"
                            />
                          </div>
                        )}
                        <span className="flex items-center gap-1">
                          <Clock className="w-3 h-3" />
                          {item.created_at &&
                          !isNaN(new Date(item.created_at).getTime())
                            ? formatDistanceToNow(new Date(item.created_at), {
                                addSuffix: true,
                                locale: dateLocale,
                              })
                            : t('dateUnknown')}
                        </span>
                        {item.views_count !== undefined && (
                          <span className="flex items-center gap-1">
                            <Eye className="w-3 h-3" />
                            {item.views_count}
                          </span>
                        )}
                      </div>
                    </div>

                    {/* Цена и действия */}
                    <div className="flex flex-col items-end gap-2">
                      {item.price && (
                        <div className="text-right">
                          <div className="flex items-center gap-2 justify-end">
                            {item.has_discount && item.old_price && (
                              <>
                                <p className="text-sm line-through text-base-content/50">
                                  {formatPrice(item.old_price, 'RSD')}
                                </p>
                                <DiscountBadge
                                  oldPrice={item.old_price}
                                  currentPrice={item.price}
                                  onClick={handlePriceHistoryClick}
                                  size="sm"
                                />
                              </>
                            )}
                          </div>
                          <p className="text-xl font-bold text-primary">
                            {formatPrice(item.price, 'RSD')}
                          </p>
                          {/* TODO: добавить поле is_escrow_enabled в API */}
                          {item.product_type === 'storefront' && (
                            <div className="text-xs text-success flex items-center gap-1 mt-1">
                              <Shield className="w-3 h-3" />
                              {t('secureDeal')}
                            </div>
                          )}
                        </div>
                      )}

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

                          {isAuthenticated && item.user_id !== user?.id && (
                            <>
                              {isOnlinePurchaseAvailable() && (
                                <button
                                  onClick={handleBuyClick}
                                  className="btn btn-primary btn-sm"
                                >
                                  <ShoppingCart className="w-4 h-4" />
                                  {t('addToCart')}
                                </button>
                              )}
                              <button
                                onClick={handleChatClick}
                                className="btn btn-outline btn-sm"
                              >
                                <MessageCircle className="w-4 h-4" />
                                {t('chat')}
                              </button>
                            </>
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

        {showQuickView && (
          <QuickView
            isOpen={showQuickView}
            onClose={() => setShowQuickView(false)}
            product={{
              id: item.id.toString(),
              title: item.title,
              price: formatPrice(item.price || 0, 'RSD'),
              description: item.description || '',
              images:
                item.images?.map((img) =>
                  configManager.buildImageUrl(img.public_url)
                ) || [],
              category: item.category?.name || '',
              seller: {
                name: item.user?.name || '',
                rating: 4.8, // TODO: получать из API
                totalReviews: 25, // TODO: получать из API
                avatar: '',
              },
              location: {
                address: formatAddressWithPrivacy(
                  getFullLocalizedAddress(item, locale),
                  item.location_privacy
                ),
                distance: getDistance(),
              },
              stats: {
                views: item.views_count || 0,
                favorites: 0, // TODO: получать из API
              },
              condition: item.condition as
                | 'new'
                | 'used'
                | 'refurbished'
                | undefined,
              storefrontId: item.storefront_id,
              storefrontName: item.storefront?.name,
              storefrontSlug: item.storefront?.slug,
              stockQuantity: 100, // TODO: get from API
            }}
            onAddToCart={handleAddToCart}
            onContact={() => handleChatClick({} as React.MouseEvent)}
          />
        )}

        {/* Price History Modal */}
        {showPriceHistory && (
          <PriceHistoryModal
            listingId={item.id}
            isOpen={showPriceHistory}
            onClose={() => setShowPriceHistory(false)}
          />
        )}
      </>
    );
  }

  // Grid view (default)
  return (
    <>
      <div
        className={`card card-compact bg-base-100 shadow-sm hover:shadow-md border ${
          item.storefront_id
            ? 'border-primary/50 ring-1 ring-primary/20'
            : 'border-base-300 dark:border-base-600'
        } transition-shadow duration-200 group`}
      >
        {/* Изображение с оверлеями */}
        <figure className="relative aspect-square overflow-hidden">
          <Link href={getItemUrl()}>
            <SafeImage
              src={imageUrl}
              alt={item.title}
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
            {itemCondition && conditionBadge[itemCondition] && (
              <div
                className={`badge badge-sm ${conditionBadge[itemCondition].class}`}
              >
                {conditionBadge[itemCondition].text}
              </div>
            )}
            {ecoScore > 7 && (
              <div className="badge badge-sm badge-success gap-1">
                <span className="text-xs">♻️</span>
                <span>{t('eco.badge')}</span>
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
          {item.images && item.images.length > 1 && (
            <div className="absolute bottom-2 right-2">
              <div className="badge badge-neutral badge-sm">
                {item.images.length} {t('photos')}
              </div>
            </div>
          )}
        </figure>

        <div className={`card-body ${gridColumns === 3 ? 'p-2' : 'p-3'}`}>
          {/* Категория */}
          {item.category?.name && (
            <div className="text-xs text-base-content/60">
              {item.category.translations && item.category.translations[locale]
                ? item.category.translations[locale]
                : item.category.name}
            </div>
          )}

          {/* Заголовок */}
          <h3
            className={`card-title ${
              gridColumns === 3 ? 'text-xs' : 'text-sm'
            } line-clamp-2 min-h-[2.5rem]`}
          >
            {item.title}
          </h3>

          {/* Продавец - скрываем при 3 столбцах */}
          {item.user?.name && gridColumns < 3 && (
            <div className="flex items-center gap-2 text-xs">
              <div className="avatar placeholder">
                <div className="bg-neutral text-neutral-content rounded-full w-5">
                  <span className="text-xs">{item.user.name[0]}</span>
                </div>
              </div>
              <span className="font-medium">{item.user.name}</span>
              {/* TODO: добавить рейтинг */}
              <div className="flex items-center gap-1">
                <Star className="w-3 h-3 fill-warning text-warning" />
                <span className="text-xs">4.8</span>
              </div>
            </div>
          )}

          {/* Локация и статистика - упрощаем для маленьких карточек */}
          <div className="flex flex-col gap-1 text-xs text-base-content/60">
            {item.location && gridColumns < 3 && (
              <div className="flex items-center gap-2">
                <span className="flex items-center gap-1">
                  <MapPin className="w-3 h-3" />
                  {formatAddressWithPrivacy(
                    getFullLocalizedAddress(item, locale),
                    item.location_privacy
                  )}
                </span>
                <DistanceBadge distance={getDistance()} variant="compact" />
              </div>
            )}
            {gridColumns === 1 && (
              <div className="flex items-center gap-3">
                <span className="flex items-center gap-1">
                  <Clock className="w-3 h-3" />
                  {(() => {
                    const date = new Date(item.created_at);
                    if (isNaN(date.getTime())) {
                      return '-';
                    }
                    return formatDistanceToNow(date, {
                      addSuffix: true,
                      locale: dateLocale,
                    });
                  })()}
                </span>
                {item.views_count !== undefined && (
                  <span className="flex items-center gap-1">
                    <Eye className="w-3 h-3" />
                    {item.views_count}
                  </span>
                )}
              </div>
            )}
          </div>

          {/* Цена и действия */}
          <div className="card-actions justify-between items-end mt-2">
            <div>
              {item.price && (
                <>
                  {item.has_discount && item.old_price && (
                    <div className="flex items-center gap-2 mb-1">
                      <p className="text-xs line-through text-base-content/50">
                        {formatPrice(item.old_price, 'RSD')}
                      </p>
                      <DiscountBadge
                        oldPrice={item.old_price}
                        currentPrice={item.price}
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
                    {formatPrice(item.price, 'RSD')}
                  </div>
                </>
              )}
              {/* TODO: добавить поле is_escrow_enabled в API */}
              {item.product_type === 'storefront' && (
                <div className="text-xs text-success flex items-center gap-1">
                  <Shield className="w-3 h-3" />
                  {t('secureDeal')}
                </div>
              )}
            </div>

            {mounted &&
              isAuthenticated &&
              item.user_id !== user?.id &&
              gridColumns < 3 && (
                <div className="flex gap-1">
                  {isOnlinePurchaseAvailable() && (
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
                  <button
                    onClick={handleChatClick}
                    className="btn btn-outline btn-xs"
                  >
                    <MessageCircle className="w-3 h-3" />
                  </button>
                </div>
              )}
          </div>
        </div>
      </div>

      {/* Quick View Modal */}
      {showQuickView && (
        <QuickView
          isOpen={showQuickView}
          onClose={() => setShowQuickView(false)}
          product={{
            id: item.id.toString(),
            title: item.title,
            price: formatPrice(item.price || 0, 'RSD'),
            description: item.description || '',
            images:
              item.images?.map((img) =>
                configManager.buildImageUrl(img.public_url)
              ) || [],
            category: item.category?.name || '',
            seller: {
              name: item.user?.name || '',
              rating: 4.8, // TODO: получать из API
              totalReviews: 25, // TODO: получать из API
              avatar: '',
            },
            location: {
              address: formatAddressWithPrivacy(
                getFullLocalizedAddress(item, locale),
                item.location_privacy
              ),
              distance: getDistance(),
            },
            stats: {
              views: item.views_count || 0,
              favorites: 0, // TODO: получать из API
            },
            condition: item.condition as
              | 'new'
              | 'used'
              | 'refurbished'
              | undefined,
            storefrontId: item.storefront_id,
            storefrontName: item.storefront?.name,
            storefrontSlug: item.storefront?.slug,
            stockQuantity: 100, // TODO: get from API
          }}
          onAddToCart={handleAddToCart}
          onContact={() => handleChatClick({} as React.MouseEvent)}
        />
      )}

      {/* Price History Modal */}
      {showPriceHistory && (
        <PriceHistoryModal
          listingId={item.id}
          isOpen={showPriceHistory}
          onClose={() => setShowPriceHistory(false)}
        />
      )}

      {/* Variant Selection Modal */}
      {showVariantModal && item.storefront?.slug && (
        <VariantSelectionModal
          isOpen={showVariantModal}
          onClose={() => setShowVariantModal(false)}
          productId={item.id}
          productName={item.title}
          productImage={mainImage?.public_url}
          storefrontSlug={item.storefront.slug}
          basePrice={item.price || 0}
          baseCurrency="RSD"
          onAddToCart={handleAddToCartWithVariant}
        />
      )}
    </>
  );
};
