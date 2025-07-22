'use client';

import { use, useEffect, useState } from 'react';
import { useLocale, useTranslations } from 'next-intl';
import Link from 'next/link';
import { storefrontProductsService } from '@/services/storefrontProducts';
import { storefrontApi } from '@/services/storefrontApi';
import type { components } from '@/types/generated/api';
import SafeImage from '@/components/SafeImage';
import AddToCartButton from '@/components/cart/AddToCartButton';

type StorefrontProduct =
  components['schemas']['backend_internal_domain_models.StorefrontProduct'];
type Storefront =
  components['schemas']['backend_internal_domain_models.Storefront'];

type Props = {
  params: Promise<{ slug: string; id: string }>;
};

export default function StorefrontProductPage({ params }: Props) {
  const { slug, id } = use(params);
  const locale = useLocale();
  const t = useTranslations();
  const [product, setProduct] = useState<StorefrontProduct | null>(null);
  const [storefront, setStorefront] = useState<Storefront | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [selectedImageIndex, setSelectedImageIndex] = useState(0);

  // Форматирование адреса с учетом приватности
  const formatAddressWithPrivacy = (address: string, privacyLevel?: string): string => {
    if (!address) return '';

    if (privacyLevel === 'exact') {
      return address;
    }

    const parts = address.split(',').map(part => part.trim());

    switch (privacyLevel) {
      case 'approximate':
      case 'street':
        // Убираем номер дома
        if (parts.length > 2) {
          const streetPart = parts[0].replace(/\d+[а-яА-Яa-zA-Z]?(\s|$)/g, '').trim();
          return streetPart ? [streetPart, ...parts.slice(1)].join(', ') : parts.slice(1).join(', ');
        }
        return parts.slice(1).join(', ');

      case 'district':
        // Оставляем только район и город
        if (parts.length > 2) {
          return parts.slice(-2).join(', ');
        }
        return address;

      case 'city_only':
      case 'city':
        // Оставляем только город
        if (parts.length > 1) {
          return parts[parts.length - 1];
        }
        return address;

      case 'hidden':
        // Скрываем адрес полностью
        return 'Адрес скрыт';

      default:
        return address;
    }
  };

  useEffect(() => {
    const fetchData = async () => {
      try {
        setIsLoading(true);

        // Загружаем данные витрины и товара параллельно
        const [productData, storefrontData] = await Promise.all([
          storefrontProductsService.getProduct(slug, parseInt(id)),
          storefrontApi.getStorefrontBySlug(slug),
        ]);

        setProduct(productData);
        setStorefront(storefrontData);
      } catch (error) {
        console.error('Error fetching product:', error);
      } finally {
        setIsLoading(false);
      }
    };

    fetchData();
  }, [slug, id]);

  if (isLoading) {
    return (
      <div className="min-h-screen bg-base-100 pt-24">
        <div className="container mx-auto px-4 py-8">
          <div className="animate-pulse">
            {/* Breadcrumbs skeleton */}
            <div className="breadcrumbs text-sm mb-6">
              <div className="h-4 bg-base-300 rounded w-80"></div>
            </div>

            <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
              {/* Image skeleton */}
              <div className="space-y-4">
                <div className="aspect-square bg-base-300 rounded-lg"></div>
                <div className="flex gap-2">
                  {[...Array(4)].map((_, i) => (
                    <div
                      key={i}
                      className="w-16 h-16 bg-base-300 rounded"
                    ></div>
                  ))}
                </div>
              </div>

              {/* Info skeleton */}
              <div className="space-y-4">
                <div className="h-8 bg-base-300 rounded w-3/4"></div>
                <div className="h-6 bg-base-300 rounded w-1/2"></div>
                <div className="h-20 bg-base-300 rounded"></div>
                <div className="h-12 bg-base-300 rounded w-1/3"></div>
              </div>
            </div>
          </div>
        </div>
      </div>
    );
  }

  if (!product || !storefront) {
    return (
      <div className="min-h-screen bg-base-100 flex items-center justify-center pt-24">
        <div className="text-center">
          <div className="text-6xl mb-4">📦</div>
          <h1 className="text-2xl font-bold mb-2">{t('products.notFound')}</h1>
          <p className="text-base-content/60 mb-6">
            {t('products.notFoundDescription')}
          </p>
          <Link
            href={`/${locale}/storefronts/${slug}`}
            className="btn btn-primary"
          >
            {t('common.back')}
          </Link>
        </div>
      </div>
    );
  }

  const images = product.images?.filter((img) => img.image_url) || [];
  const mainImage =
    images[selectedImageIndex]?.image_url || images[0]?.image_url || '';

  return (
    <div className="min-h-screen bg-base-100 pt-24">
      <div className="container mx-auto px-4 py-8">
        {/* Breadcrumbs */}
        <div className="breadcrumbs text-sm mb-6">
          <ul>
            <li>
              <Link href={`/${locale}`}>{t('common.home')}</Link>
            </li>
            <li>
              <Link href={`/${locale}/storefronts/${slug}`}>
                {storefront.name}
              </Link>
            </li>
            <li>
              <Link href={`/${locale}/storefronts/${slug}/products`}>
                {t('storefronts.products.title')}
              </Link>
            </li>
            <li className="text-base-content/60">{product.name}</li>
          </ul>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
          {/* Product Images */}
          <div className="space-y-4">
            {/* Main Image */}
            <div className="aspect-square relative bg-base-200 rounded-lg overflow-hidden">
              <SafeImage
                src={mainImage}
                alt={product.name || ''}
                fill
                className="object-cover"
                fallback={
                  <div className="w-full h-full flex items-center justify-center">
                    <svg
                      className="w-24 h-24 text-base-content/20"
                      fill="none"
                      viewBox="0 0 24 24"
                      stroke="currentColor"
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={2}
                        d="m2.25 15.75 5.159-5.159a2.25 2.25 0 0 1 3.182 0l5.159 5.159m-1.5-1.5 1.409-1.409a2.25 2.25 0 0 1 3.182 0l2.909 2.909m-18 3.75h16.5a1.5 1.5 0 0 0 1.5-1.5V6a1.5 1.5 0 0 0-1.5-1.5H3.75A1.5 1.5 0 0 0 2.25 6v12a1.5 1.5 0 0 0 1.5 1.5Zm10.5-11.25h.008v.008h-.008V8.25Zm.375 0a.375.375 0 1 1-.75 0 .375.375 0 0 1 .75 0Z"
                      />
                    </svg>
                  </div>
                }
              />
              {product.stock_status === 'out_of_stock' && (
                <div className="absolute top-4 right-4 badge badge-warning">
                  {t('cart.outOfStock')}
                </div>
              )}
            </div>

            {/* Thumbnail Images */}
            {images.length > 1 && (
              <div className="flex gap-2 overflow-x-auto">
                {images.map((image, index) => (
                  <button
                    key={index}
                    onClick={() => setSelectedImageIndex(index)}
                    className={`w-16 h-16 flex-shrink-0 rounded overflow-hidden border-2 transition-colors ${
                      selectedImageIndex === index
                        ? 'border-primary'
                        : 'border-base-300'
                    }`}
                  >
                    <SafeImage
                      src={image.image_url || ''}
                      alt={`${product.name} ${index + 1}`}
                      width={64}
                      height={64}
                      className="object-cover w-full h-full"
                    />
                  </button>
                ))}
              </div>
            )}
          </div>

          {/* Product Info */}
          <div className="space-y-6">
            {/* Store Info */}
            <div className="flex items-center gap-3 p-3 bg-base-200 rounded-lg">
              <div className="avatar placeholder">
                <div className="bg-neutral text-neutral-content rounded-full w-10">
                  <span className="text-sm">
                    {storefront.name?.charAt(0).toUpperCase()}
                  </span>
                </div>
              </div>
              <div>
                <Link
                  href={`/${locale}/storefronts/${slug}`}
                  className="font-semibold hover:link"
                >
                  {storefront.name}
                </Link>
                <div className="text-sm text-base-content/60">
                  {storefront.is_verified && (
                    <span className="badge badge-success badge-sm mr-2">
                      {t('storefronts.verified')}
                    </span>
                  )}
                  {t('storefronts.store')}
                </div>
              </div>
            </div>

            {/* Product Name */}
            <h1 className="text-3xl font-bold">{product.name}</h1>

            {/* Price */}
            <div className="text-3xl font-bold text-primary">
              {product.price} {product.currency || 'RSD'}
            </div>

            {/* Stock Info */}
            <div className="flex items-center gap-4">
              <div
                className={`badge ${
                  product.stock_status === 'in_stock'
                    ? 'badge-success'
                    : product.stock_status === 'low_stock'
                      ? 'badge-warning'
                      : 'badge-error'
                }`}
              >
                {t(`products.status.${product.stock_status}`)}
              </div>
              {product.stock_quantity !== undefined &&
                product.stock_quantity > 0 && (
                  <span className="text-sm text-base-content/60">
                    {t('products.inStock')}: {product.stock_quantity}
                  </span>
                )}
              {product.sold_count && product.sold_count > 0 && (
                <span className="text-sm text-base-content/60">
                  {t('products.sold')}: {product.sold_count}
                </span>
              )}
            </div>

            {/* Category */}
            {product.category && (
              <div className="flex items-center gap-2">
                <span className="text-sm text-base-content/60">
                  {t('common.category')}:
                </span>
                <span className="badge badge-outline">
                  {product.category.name}
                </span>
              </div>
            )}

            {/* SKU */}
            {product.sku && (
              <div className="text-sm text-base-content/60">
                SKU: <span className="font-mono">{product.sku}</span>
              </div>
            )}

            {/* Location */}
            {product.individual_address && (
              <div className="flex items-center gap-2 text-sm text-base-content/70">
                <svg
                  className="w-5 h-5 flex-shrink-0"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z"
                  />
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M15 11a3 3 0 11-6 0 3 3 0 016 0z"
                  />
                </svg>
                {formatAddressWithPrivacy(product.individual_address, product.location_privacy)}
              </div>
            )}

            {/* Description */}
            {product.description && (
              <div className="prose max-w-none">
                <h3 className="text-lg font-semibold mb-2">
                  {t('products.description')}
                </h3>
                <p className="text-base-content/80 whitespace-pre-wrap">
                  {product.description}
                </p>
              </div>
            )}

            {/* Add to Cart */}
            <div className="pt-4">
              <AddToCartButton
                product={{
                  id: product.id!,
                  name: product.name!,
                  price: product.price!,
                  currency: product.currency || 'RSD',
                  image: mainImage,
                  storefrontId: product.storefront_id!,
                  stockQuantity: product.stock_quantity || 0,
                  stockStatus: product.stock_status || 'out_of_stock',
                }}
                className="btn btn-primary btn-lg w-full"
                disabled={product.stock_status === 'out_of_stock'}
              />
            </div>

            {/* Trust Badges */}
            <div className="alert alert-info">
              <svg
                className="w-5 h-5 flex-shrink-0"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
                />
              </svg>
              <div>
                <h4 className="font-semibold">
                  {t('storefronts.trustSafety')}
                </h4>
                <ul className="text-sm space-y-1 mt-1">
                  <li>• {t('storefronts.securePayments')}</li>
                  <li>• {t('storefronts.buyerProtection')}</li>
                  <li>• {t('orders.trackingAvailable')}</li>
                </ul>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
