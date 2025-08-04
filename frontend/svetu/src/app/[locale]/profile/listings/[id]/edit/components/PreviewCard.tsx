'use client';

import { useTranslations } from 'next-intl';
import Image from 'next/image';

interface PreviewCardProps {
  data: {
    title: string;
    description: string;
    price: number;
    condition: string;
    city: string;
    category: {
      name: string;
    };
    images: Array<{
      id: number;
      public_url: string;
      is_main: boolean;
    }>;
  };
  viewMode: 'card' | 'list' | 'mobile';
}

export function PreviewCard({ data, viewMode }: PreviewCardProps) {
  const t = useTranslations('profile');

  const mainImage = data.images.find((img) => img.is_main) || data.images[0];

  const formatPrice = (price: number) => {
    return new Intl.NumberFormat('sr-RS', {
      style: 'currency',
      currency: 'RSD',
      minimumFractionDigits: 0,
      maximumFractionDigits: 0,
    }).format(price);
  };

  const getConditionLabel = (condition: string) => {
    switch (condition) {
      case 'new':
        return t('condition.new');
      case 'used':
        return t('condition.used');
      case 'refurbished':
        return t('condition.refurbished');
      default:
        return condition;
    }
  };

  if (viewMode === 'list') {
    return (
      <div className="card card-side bg-base-100 shadow-xl">
        <figure className="w-48 h-48">
          {mainImage ? (
            <Image
              src={mainImage.public_url}
              alt={data.title}
              fill
              className="object-cover"
            />
          ) : (
            <div className="w-full h-full bg-base-200 flex items-center justify-center">
              <svg
                className="w-12 h-12 text-base-content/20"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"
                />
              </svg>
            </div>
          )}
        </figure>
        <div className="card-body">
          <div className="text-sm text-base-content/70">
            {data.category.name}
          </div>
          <h2 className="card-title">{data.title || t('preview.untitled')}</h2>
          <p className="text-sm line-clamp-2">
            {data.description || t('preview.noDescription')}
          </p>
          <div className="flex items-center gap-4 mt-2">
            <span className="text-2xl font-bold text-primary">
              {formatPrice(data.price)}
            </span>
            <div className="badge badge-outline">
              {getConditionLabel(data.condition)}
            </div>
            {data.city && (
              <span className="text-sm text-base-content/70">{data.city}</span>
            )}
          </div>
        </div>
      </div>
    );
  }

  if (viewMode === 'mobile') {
    return (
      <div className="w-full max-w-sm mx-auto">
        <div className="card bg-base-100 shadow-xl">
          <figure className="aspect-square">
            {mainImage ? (
              <Image
                src={mainImage.public_url}
                alt={data.title}
                fill
                className="object-cover"
              />
            ) : (
              <div className="w-full h-full bg-base-200 flex items-center justify-center">
                <svg
                  className="w-16 h-16 text-base-content/20"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"
                  />
                </svg>
              </div>
            )}
            {data.images.length > 1 && (
              <div className="absolute bottom-2 right-2 badge badge-neutral badge-sm">
                {data.images.length} {t('preview.photos')}
              </div>
            )}
          </figure>
          <div className="card-body p-4">
            <div className="text-xs text-base-content/70">
              {data.category.name}
            </div>
            <h2 className="card-title text-lg">
              {data.title || t('preview.untitled')}
            </h2>
            <p className="text-sm line-clamp-2">
              {data.description || t('preview.noDescription')}
            </p>
            <div className="flex flex-col gap-2 mt-2">
              <span className="text-xl font-bold text-primary">
                {formatPrice(data.price)}
              </span>
              <div className="flex items-center gap-2">
                <div className="badge badge-outline badge-sm">
                  {getConditionLabel(data.condition)}
                </div>
                {data.city && (
                  <span className="text-xs text-base-content/70">
                    {data.city}
                  </span>
                )}
              </div>
            </div>
          </div>
        </div>
      </div>
    );
  }

  // Default card view
  return (
    <div className="card bg-base-100 shadow-xl">
      <figure className="aspect-[4/3]">
        {mainImage ? (
          <Image
            src={mainImage.public_url}
            alt={data.title}
            fill
            className="object-cover"
          />
        ) : (
          <div className="w-full h-full bg-base-200 flex items-center justify-center">
            <svg
              className="w-16 h-16 text-base-content/20"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"
              />
            </svg>
          </div>
        )}
        {data.images.length > 1 && (
          <div className="absolute bottom-2 right-2 badge badge-neutral badge-sm">
            {data.images.length} {t('preview.photos')}
          </div>
        )}
      </figure>
      <div className="card-body">
        <div className="text-sm text-base-content/70">{data.category.name}</div>
        <h2 className="card-title">{data.title || t('preview.untitled')}</h2>
        <p className="line-clamp-2">
          {data.description || t('preview.noDescription')}
        </p>
        <div className="flex justify-between items-center mt-4">
          <span className="text-2xl font-bold text-primary">
            {formatPrice(data.price)}
          </span>
          <div className="flex items-center gap-2">
            <div className="badge badge-outline">
              {getConditionLabel(data.condition)}
            </div>
            {data.city && (
              <span className="text-sm text-base-content/70">{data.city}</span>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
