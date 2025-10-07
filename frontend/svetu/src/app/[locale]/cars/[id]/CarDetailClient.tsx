'use client';

import { useState, useMemo } from 'react';
import { useTranslations } from 'next-intl';
import Image from 'next/image';
import { useRouter } from 'next/navigation';
import type { components } from '@/types/generated/api';
import {
  Car,
  Calendar,
  Fuel,
  Settings,
  MapPin,
  Phone,
  Heart,
  Share2,
  Shield,
  Clock,
  Eye,
  ArrowLeft,
  ChevronLeft,
  ChevronRight,
} from 'lucide-react';
import UniversalCreditCalculator from '@/components/universal/calculators/UniversalCreditCalculator';
import RecommendationsEngine from '@/components/universal/recommendations/RecommendationsEngine';
import { useAppDispatch, useAppSelector } from '@/store/hooks';
import {
  addToFavorites,
  removeFromFavorites,
} from '@/store/slices/favoritesSlice';
import { addItem, removeItem } from '@/store/slices/universalCompareSlice';
import type { UniversalListingData } from '@/components/universal/cards/UniversalListingCard';
import {
  makeSelectCompareItemsByType,
  selectFavorites,
} from '@/store/selectors/compareSelectors';

type MarketplaceListing =
  components['schemas']['models.MarketplaceListing'];

interface CarDetailClientProps {
  car: MarketplaceListing;
  locale: string;
}

export default function CarDetailClient({ car, locale }: CarDetailClientProps) {
  const _t = useTranslations('cars');
  const router = useRouter();
  const dispatch = useAppDispatch();
  const [currentImageIndex, setCurrentImageIndex] = useState(0);
  const [showCreditCalculator, setShowCreditCalculator] = useState(false);

  const selectCompareItemsByType = useMemo(makeSelectCompareItemsByType, []);
  const favorites = useAppSelector(selectFavorites);
  const compareItems = useAppSelector((state) =>
    selectCompareItemsByType(state, 'cars')
  );

  const isFavorite = useMemo(
    () => favorites.some((item) => item.id === car.id),
    [favorites, car.id]
  );
  const isComparing = useMemo(
    () => compareItems.some((item) => item.id === car.id),
    [compareItems, car.id]
  );

  // Parse attributes array into object for easy access
  const carAttrs = useMemo(() => {
    const attrs: Record<string, any> = {};
    const translations: Record<string, any> = {};

    if (Array.isArray(car.attributes)) {
      car.attributes.forEach((attr: any) => {
        // Use attribute_name as key and appropriate value based on type
        if (attr.attribute_name) {
          attrs[attr.attribute_name] =
            attr.text_value || attr.numeric_value || attr.display_value || null;

          // Store translations for select/multiselect attributes
          if (attr.option_translations && attr.text_value) {
            const localeTranslations =
              attr.option_translations[locale] ||
              attr.option_translations.ru ||
              {};
            translations[attr.attribute_name] =
              localeTranslations[attr.text_value] || attr.text_value;
          }
        }
      });
    }

    return { raw: attrs, translated: translations };
  }, [car.attributes, locale]);

  const images = car.images || [];

  // Convert to universal format for compare
  const universalData: UniversalListingData = {
    id: car.id || 0,
    title: car.title || '',
    price: car.price || 0,
    currency: (car as any).currency || 'USD',
    images: images
      .map((img) => (img as any).url || img.thumbnail_url || '')
      .filter(Boolean),
    location: {
      city: car.city,
      district: (car as any).district,
      address: (car as any).address,
    },
    category: 'cars',
    categorySlug: 'cars',
    createdAt: car.created_at,
    updatedAt: car.updated_at,
    attributes: carAttrs.raw,
    customFields: [],
    badges: [],
    stats: {
      views: (car as any).views || 0,
      favorites: 0,
    },
  };

  const handleFavorite = () => {
    if (isFavorite) {
      dispatch(removeFromFavorites({ id: car.id || 0 }));
    } else {
      dispatch(addToFavorites(universalData));
    }
  };

  const handleCompare = () => {
    if (isComparing) {
      dispatch(removeItem(car.id || 0));
    } else {
      dispatch(addItem({ ...universalData, category: 'cars' } as any));
    }
  };

  const handleShare = () => {
    if (navigator.share) {
      navigator.share({
        title: car.title,
        url: window.location.href,
      });
    } else {
      navigator.clipboard.writeText(window.location.href);
    }
  };

  const nextImage = () => {
    setCurrentImageIndex((prev) => (prev + 1) % images.length);
  };

  const prevImage = () => {
    setCurrentImageIndex((prev) => (prev - 1 + images.length) % images.length);
  };

  return (
    <div className="min-h-screen bg-base-100">
      {/* Header */}
      <div className="bg-base-200 py-4">
        <div className="container mx-auto px-4">
          <button
            onClick={() => router.back()}
            className="btn btn-ghost btn-sm gap-2"
          >
            <ArrowLeft className="w-4 h-4" />
            Назад
          </button>
        </div>
      </div>

      {/* Main Content */}
      <div className="container mx-auto px-4 py-8">
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          {/* Left Column - Images */}
          <div className="lg:col-span-2">
            {/* Main Image */}
            {images.length > 0 ? (
              <div className="relative aspect-video bg-base-200 rounded-lg overflow-hidden">
                <Image
                  src={
                    (images[currentImageIndex] as any)?.url ||
                    images[currentImageIndex]?.thumbnail_url ||
                    '/images/car-placeholder.jpg'
                  }
                  alt={car.title || ''}
                  fill
                  sizes="(max-width: 768px) 100vw, (max-width: 1200px) 66vw, 66vw"
                  className="object-cover"
                />
                {images.length > 1 && (
                  <>
                    <button
                      onClick={prevImage}
                      className="absolute left-2 top-1/2 -translate-y-1/2 btn btn-circle btn-sm"
                    >
                      <ChevronLeft className="w-4 h-4" />
                    </button>
                    <button
                      onClick={nextImage}
                      className="absolute right-2 top-1/2 -translate-y-1/2 btn btn-circle btn-sm"
                    >
                      <ChevronRight className="w-4 h-4" />
                    </button>
                  </>
                )}
                <div className="absolute bottom-4 right-4 badge badge-neutral">
                  {currentImageIndex + 1} / {images.length}
                </div>
              </div>
            ) : (
              <div className="aspect-video bg-base-200 rounded-lg flex items-center justify-center">
                <Car className="w-24 h-24 text-base-300" />
              </div>
            )}

            {/* Thumbnails */}
            {images.length > 1 && (
              <div className="grid grid-cols-6 gap-2 mt-4">
                {images.slice(0, 6).map((img, idx) => (
                  <button
                    key={idx}
                    onClick={() => setCurrentImageIndex(idx)}
                    className={`aspect-video bg-base-200 rounded overflow-hidden border-2 ${
                      idx === currentImageIndex
                        ? 'border-primary'
                        : 'border-transparent'
                    }`}
                  >
                    <Image
                      src={img.thumbnail_url || (img as any).url || ''}
                      alt=""
                      width={120}
                      height={80}
                      className="object-cover w-full h-full"
                    />
                  </button>
                ))}
              </div>
            )}

            {/* Details */}
            <div className="mt-8 space-y-6">
              <div>
                <h1 className="text-3xl font-bold">{car.title}</h1>
                <div className="flex items-center gap-4 mt-2 text-base-content/60">
                  <span className="flex items-center gap-1">
                    <Eye className="w-4 h-4" />
                    {(car as any).views || 0} просмотров
                  </span>
                  <span className="flex items-center gap-1">
                    <Clock className="w-4 h-4" />
                    {new Date(car.created_at || '').toLocaleDateString('ru')}
                  </span>
                </div>
              </div>

              {/* Specifications */}
              <div className="card bg-base-100 shadow-xl">
                <div className="card-body">
                  <h2 className="card-title">Характеристики</h2>
                  <div className="grid grid-cols-2 gap-4">
                    <div className="flex items-center gap-3">
                      <Calendar className="w-5 h-5 text-primary" />
                      <div>
                        <div className="text-xs text-base-content/60">Год</div>
                        <div className="font-medium">
                          {carAttrs.raw.year || 'Не указан'}
                        </div>
                      </div>
                    </div>
                    <div className="flex items-center gap-3">
                      <Car className="w-5 h-5 text-primary" />
                      <div>
                        <div className="text-xs text-base-content/60">
                          Пробег
                        </div>
                        <div className="font-medium">
                          {carAttrs.raw.mileage
                            ? `${carAttrs.raw.mileage.toLocaleString()} км`
                            : 'Не указан'}
                        </div>
                      </div>
                    </div>
                    <div className="flex items-center gap-3">
                      <Fuel className="w-5 h-5 text-primary" />
                      <div>
                        <div className="text-xs text-base-content/60">
                          Топливо
                        </div>
                        <div className="font-medium">
                          {carAttrs.translated.fuel_type ||
                            carAttrs.raw.fuel_type ||
                            'Не указано'}
                        </div>
                      </div>
                    </div>
                    <div className="flex items-center gap-3">
                      <Settings className="w-5 h-5 text-primary" />
                      <div>
                        <div className="text-xs text-base-content/60">КПП</div>
                        <div className="font-medium">
                          {carAttrs.translated.transmission ||
                            carAttrs.raw.transmission ||
                            'Не указана'}
                        </div>
                      </div>
                    </div>
                    <div className="flex items-center gap-3">
                      <Shield className="w-5 h-5 text-primary" />
                      <div>
                        <div className="text-xs text-base-content/60">
                          Мощность
                        </div>
                        <div className="font-medium">
                          {carAttrs.raw.power
                            ? `${carAttrs.raw.power} л.с.`
                            : 'Не указана'}
                        </div>
                      </div>
                    </div>
                    <div className="flex items-center gap-3">
                      <Settings className="w-5 h-5 text-primary" />
                      <div>
                        <div className="text-xs text-base-content/60">
                          Объем двигателя
                        </div>
                        <div className="font-medium">
                          {carAttrs.raw.engine_size
                            ? `${carAttrs.raw.engine_size} л`
                            : 'Не указан'}
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>

              {/* Description */}
              {car.description && (
                <div className="card bg-base-100 shadow-xl">
                  <div className="card-body">
                    <h2 className="card-title">Описание</h2>
                    <p className="whitespace-pre-wrap">{car.description}</p>
                  </div>
                </div>
              )}
            </div>
          </div>

          {/* Right Column - Price and Actions */}
          <div className="lg:col-span-1">
            <div className="sticky top-4 space-y-4">
              {/* Price Card */}
              <div className="card bg-base-100 shadow-xl">
                <div className="card-body">
                  <div className="text-3xl font-bold text-primary">
                    {car.price
                      ? `${car.price.toLocaleString()} ${(car as any).currency || '€'}`
                      : 'Цена договорная'}
                  </div>

                  {/* Action Buttons */}
                  <div className="space-y-2 mt-4">
                    <a
                      href={`tel:${(car as any).phone || ''}`}
                      className="btn btn-primary btn-block"
                    >
                      <Phone className="w-4 h-4" />
                      Позвонить
                    </a>
                    <button
                      onClick={() =>
                        setShowCreditCalculator(!showCreditCalculator)
                      }
                      className="btn btn-outline btn-block"
                    >
                      Рассчитать кредит
                    </button>
                    <div className="flex gap-2">
                      <button
                        onClick={handleFavorite}
                        className={`btn btn-outline flex-1 ${isFavorite ? 'btn-error' : ''}`}
                      >
                        <Heart
                          className={`w-4 h-4 ${isFavorite ? 'fill-current' : ''}`}
                        />
                        {isFavorite ? 'В избранном' : 'В избранное'}
                      </button>
                      <button
                        onClick={handleCompare}
                        className={`btn btn-outline flex-1 ${isComparing ? 'btn-success' : ''}`}
                      >
                        Сравнить
                      </button>
                      <button
                        onClick={handleShare}
                        className="btn btn-outline btn-square"
                      >
                        <Share2 className="w-4 h-4" />
                      </button>
                    </div>
                  </div>

                  {/* Seller Info */}
                  <div className="divider"></div>
                  <div className="space-y-2">
                    <div className="font-medium">Продавец</div>
                    <div className="text-sm text-base-content/60">
                      {(car as any).user_name || 'Частное лицо'}
                    </div>
                    {car.city && (
                      <div className="flex items-center gap-2 text-sm">
                        <MapPin className="w-4 h-4" />
                        {car.city}
                      </div>
                    )}
                  </div>
                </div>
              </div>

              {/* Credit Calculator */}
              {showCreditCalculator && car.price && (
                <div className="card bg-base-100 shadow-xl">
                  <div className="card-body">
                    <UniversalCreditCalculator
                      price={car.price}
                      category="cars"
                      config={
                        {
                          minDownPayment: 10,
                          maxDownPayment: 90,
                          defaultDownPayment: 20,
                          minTerm: 12,
                          maxTerm: 84,
                          defaultTerm: 60,
                          defaultRate: 8.5,
                        } as any
                      }
                    />
                  </div>
                </div>
              )}
            </div>
          </div>
        </div>

        {/* Recommendations */}
        <div className="mt-12">
          <h2 className="text-2xl font-bold mb-6">Похожие автомобили</h2>
          <RecommendationsEngine
            category="cars"
            currentItemId={car.id}
            type="similar"
            limit={4}
          />
        </div>

        <div className="mt-12">
          <h2 className="text-2xl font-bold mb-6">Рекомендуем посмотреть</h2>
          <RecommendationsEngine
            category="cars"
            currentItemId={car.id}
            type="personal"
            limit={4}
          />
        </div>
      </div>
    </div>
  );
}
