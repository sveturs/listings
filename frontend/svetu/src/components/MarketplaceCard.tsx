import { MarketplaceItem, MarketplaceImage } from '@/types/marketplace';
import Image from 'next/image';
import Link from 'next/link';
import configManager from '@/config';

interface MarketplaceCardProps {
  item: MarketplaceItem;
  locale: string;
}

export default function MarketplaceCard({
  item,
  locale,
}: MarketplaceCardProps) {
  const formatPrice = (price?: number, currency?: string) => {
    if (!price) return '';

    const formatter = new Intl.NumberFormat(locale, {
      style: 'currency',
      currency: currency || 'USD',
      minimumFractionDigits: 0,
      maximumFractionDigits: 0,
    });

    return formatter.format(price);
  };

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    const now = new Date();
    const diffInHours = Math.floor(
      (now.getTime() - date.getTime()) / (1000 * 60 * 60)
    );

    if (diffInHours < 1) return locale === 'ru' ? 'Только что' : 'Just now';
    if (diffInHours < 24)
      return locale === 'ru'
        ? `${diffInHours} ч. назад`
        : `${diffInHours}h ago`;

    const diffInDays = Math.floor(diffInHours / 24);
    if (diffInDays < 7)
      return locale === 'ru' ? `${diffInDays} д. назад` : `${diffInDays}d ago`;

    return date.toLocaleDateString(locale);
  };

  const getImageUrl = (image?: MarketplaceImage) => {
    if (!image) return null;
    return configManager.buildImageUrl(image.public_url);
  };

  const mainImage = item.images?.find((img) => img.is_main) || item.images?.[0];
  const imageUrl = getImageUrl(mainImage);

  return (
    <Link
      href={`/${locale}/marketplace/${item.id}`}
      className="card card-compact bg-base-100 shadow-xl hover:shadow-2xl transition-shadow"
    >
      <figure className="relative h-48 bg-base-200">
        {imageUrl ? (
          <Image
            src={imageUrl}
            alt={item.title}
            fill
            className="object-cover"
            sizes="(max-width: 768px) 100vw, (max-width: 1200px) 50vw, 33vw"
          />
        ) : (
          <div className="flex items-center justify-center w-full h-full text-base-content/50">
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
        )}
      </figure>
      <div className="card-body">
        <h2 className="card-title line-clamp-2">{item.title}</h2>
        {item.price && (
          <div>
            {item.has_discount && item.old_price && (
              <p className="text-sm line-through text-base-content/50">
                {formatPrice(item.old_price, 'USD')}
              </p>
            )}
            <p className="text-lg font-bold text-primary">
              {formatPrice(item.price, 'USD')}
            </p>
          </div>
        )}
        {item.location && (
          <p className="text-sm text-base-content/70 flex items-center gap-1">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              strokeWidth={1.5}
              stroke="currentColor"
              className="w-4 h-4"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                d="M15 10.5a3 3 0 1 1-6 0 3 3 0 0 1 6 0Z"
              />
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                d="M19.5 10.5c0 7.142-7.5 11.25-7.5 11.25S4.5 17.642 4.5 10.5a7.5 7.5 0 1 1 15 0Z"
              />
            </svg>
            {item.location}
          </p>
        )}
        <p className="text-xs text-base-content/50">
          {formatDate(item.created_at)}
        </p>
      </div>
    </Link>
  );
}
