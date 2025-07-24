'use client';

import React from 'react';
import Link from 'next/link';
import {
  TrendingUp,
  Star,
  Shield,
  MapPin,
  Clock,
  Users,
  Sparkles,
} from 'lucide-react';
import { useCountAnimation } from '@/hooks/useCountAnimation';
import SearchBar from '@/components/SearchBar/SearchBar';
import { BentoGridMap } from './BentoGridMap';

interface BentoGridProps {
  categories?: Array<{
    id: string;
    name: string;
    count: number;
    icon?: React.ReactNode;
  }>;
  featuredListing?: {
    id: string;
    title: string;
    price: string;
    image: string;
    category: string;
  };
  stats?: {
    totalListings: number;
    activeUsers: number;
    successfulDeals: number;
    newUsersToday?: number;
    nearbyListings?: number;
    newListingsLastHour?: number;
    priceDropsToday?: number;
  };
  nearbyListings?: Array<{
    id: string;
    latitude: number;
    longitude: number;
    price: number;
  }>;
  userLocation?: {
    latitude: number;
    longitude: number;
  };
}

export const BentoGrid: React.FC<BentoGridProps> = ({
  categories = [],
  featuredListing,
  stats,
  nearbyListings = [],
  userLocation,
}) => {
  const animatedListings = useCountAnimation(stats?.totalListings || 0, 2000);
  const animatedUsers = useCountAnimation(stats?.activeUsers || 0, 2000);
  const animatedDeals = useCountAnimation(stats?.successfulDeals || 0, 2000);

  return (
    <>
      {/* Мобильная версия - только поиск */}
      <div className="lg:hidden p-4">
        <div className="card bg-gradient-to-br from-primary/20 via-primary/10 to-transparent border border-primary/20 shadow-lg">
          <div className="card-body">
            <h2 className="text-xl font-bold mb-4">
              Найдите то, что нужно именно вам
            </h2>
            <SearchBar variant="minimal" showTrending={true} />

            <div className="mt-4 grid grid-cols-2 gap-4 text-center">
              <div>
                <p className="text-sm text-base-content/60">Объявлений</p>
                <p className="text-lg font-bold text-primary">
                  {animatedListings.toLocaleString()}
                </p>
              </div>
              <div>
                <p className="text-sm text-base-content/60">Продавцов</p>
                <p className="text-lg font-bold text-secondary">
                  {animatedUsers.toLocaleString()}
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Десктопная версия - полный BentoGrid */}
      <div className="hidden lg:grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 p-4 max-w-7xl mx-auto">
        {/* Hero Card - Большая карточка с поиском */}
        <div className="col-span-1 md:col-span-2 lg:col-span-2 row-span-1 group">
          <div className="card h-full bg-gradient-to-br from-primary/20 via-primary/10 to-transparent border border-primary/20 shadow-lg hover:shadow-2xl hover:scale-[1.02] transition-all duration-500 overflow-hidden">
            <div className="card-body p-6 lg:p-8 flex flex-col justify-between h-full">
              <div>
                <div className="flex items-center gap-2 mb-4">
                  <Sparkles className="w-6 h-6 text-primary" />
                  <span className="badge badge-primary">Новое</span>
                </div>
                <h2 className="card-title text-2xl lg:text-3xl font-bold mb-6">
                  Найдите то, что нужно именно вам
                </h2>

                {/* Интегрированный поиск */}
                <div className="mb-6">
                  <SearchBar variant="minimal" showTrending={true} />
                </div>
              </div>

              <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
                <div>
                  <p className="text-sm text-base-content/60">
                    Активных объявлений
                  </p>
                  <p className="text-2xl font-bold text-primary">
                    {animatedListings.toLocaleString()}
                  </p>
                </div>
                <div className="text-base-content/70 text-sm text-right">
                  Тысячи товаров и услуг
                  <br className="hidden sm:inline" /> от проверенных продавцов
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* Популярные категории */}
        <div className="col-span-1 md:col-span-2 lg:col-span-2 row-span-1 group">
          <div className="card h-full bg-base-100 border border-base-200 shadow-lg hover:shadow-xl hover:scale-[1.02] transition-all duration-300">
            <div className="card-body p-6">
              <div className="flex items-center justify-between mb-4">
                <TrendingUp className="w-6 h-6 text-warning" />
                <span className="text-xs text-base-content/60">Обновлено</span>
              </div>
              <h3 className="font-semibold mb-3">Популярные категории</h3>
              <div className="space-y-2">
                {categories.length > 0 ? (
                  categories.slice(0, 4).map((cat) => (
                    <Link
                      key={cat.id}
                      href={`/search?category=${cat.id}`}
                      className="flex items-center justify-between p-2 rounded-lg hover:bg-base-200 transition-colors"
                    >
                      <span className="text-sm">{cat.name}</span>
                      <span className="text-xs text-base-content/60">
                        {cat.count.toLocaleString()}
                      </span>
                    </Link>
                  ))
                ) : (
                  <p className="text-sm text-base-content/60">
                    Загрузка категорий...
                  </p>
                )}
              </div>
            </div>
          </div>
        </div>

        {/* Избранное объявление */}
        {featuredListing ? (
          <div className="col-span-1 md:col-span-2 lg:col-span-1 row-span-1 group">
            <Link href={`/listing/${featuredListing.id}`}>
              <div className="card h-full bg-base-100 shadow-lg hover:shadow-xl transition-all duration-300 overflow-hidden">
                <figure className="h-32 relative">
                  {/* eslint-disable-next-line @next/next/no-img-element */}
                  <img
                    src={featuredListing.image}
                    alt={featuredListing.title}
                    className="w-full h-full object-cover group-hover:scale-110 transition-transform duration-300"
                  />
                  <div className="absolute top-2 left-2">
                    <span className="badge badge-warning gap-1">
                      <Star className="w-3 h-3" />
                      Рекомендуем
                    </span>
                  </div>
                </figure>
                <div className="card-body p-4">
                  <p className="text-xs text-base-content/60">
                    {featuredListing.category}
                  </p>
                  <h3 className="font-semibold line-clamp-1">
                    {featuredListing.title}
                  </h3>
                  <p className="text-lg font-bold text-primary">
                    {featuredListing.price}
                  </p>
                </div>
              </div>
            </Link>
          </div>
        ) : (
          <div className="col-span-1 md:col-span-2 lg:col-span-1 row-span-1">
            <div className="card h-full bg-base-100 shadow-lg">
              <div className="skeleton h-32 w-full"></div>
              <div className="card-body p-4">
                <div className="skeleton h-4 w-20 mb-2"></div>
                <div className="skeleton h-4 w-full mb-2"></div>
                <div className="skeleton h-6 w-24"></div>
              </div>
            </div>
          </div>
        )}

        {/* Статистика пользователей */}
        <div className="col-span-1 row-span-1 group">
          <div className="card h-full bg-gradient-to-br from-info/20 to-info/5 border border-info/20 shadow-lg hover:shadow-xl hover:scale-[1.02] transition-all duration-300">
            <div className="card-body p-6">
              <div className="flex items-center gap-3 mb-4">
                <Users className="w-6 h-6 text-info" />
                <h3 className="font-semibold">Сообщество</h3>
              </div>
              <div className="space-y-3">
                <div className="flex justify-between items-center">
                  <span className="text-sm text-base-content/70">
                    Активных продавцов
                  </span>
                  <span className="font-bold">
                    {animatedUsers.toLocaleString()}
                  </span>
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-sm text-base-content/70">
                    Новых сегодня
                  </span>
                  <span className="font-bold text-info">
                    +{stats?.newUsersToday || 12}
                  </span>
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* Рядом с вами - Карта */}
        <div className="col-span-1 md:col-span-2 lg:col-span-2 row-span-2 group">
          <div className="card h-full bg-gradient-to-br from-secondary/20 to-secondary/5 border border-secondary/20 shadow-lg hover:shadow-xl hover:scale-[1.02] transition-all duration-300 overflow-hidden">
            <div className="card-body p-4 flex flex-col h-full">
              <div className="flex items-center gap-3 mb-3">
                <MapPin className="w-6 h-6 text-secondary" />
                <div>
                  <h3 className="font-semibold">Рядом с вами</h3>
                  <p className="text-xs text-base-content/60">
                    {nearbyListings.length > 0
                      ? `${nearbyListings.length} объявлений в радиусе 5 км`
                      : 'Загрузка карты...'}
                  </p>
                </div>
              </div>
              <div className="flex-1 relative rounded-lg overflow-hidden">
                <BentoGridMap
                  listings={nearbyListings}
                  userLocation={userLocation}
                />
              </div>
            </div>
          </div>
        </div>

        {/* Безопасные сделки */}
        <div className="col-span-1 row-span-1 group">
          <div className="card h-full bg-gradient-to-br from-success/20 to-success/5 border border-success/20 shadow-lg hover:shadow-xl hover:scale-[1.02] transition-all duration-300">
            <div className="card-body p-6">
              <div className="flex items-center gap-3 mb-4">
                <Shield className="w-6 h-6 text-success" />
                <h3 className="font-semibold">Безопасные сделки</h3>
              </div>
              <div className="text-center">
                <p className="text-3xl font-bold text-success">
                  {animatedListings > 0
                    ? Math.min(
                        Math.floor((animatedDeals / animatedListings) * 100),
                        99
                      )
                    : 98}
                  %
                </p>
                <p className="text-sm text-base-content/60 mt-1">
                  успешных сделок
                </p>
              </div>
            </div>
          </div>
        </div>

        {/* Последние обновления */}
        <div className="col-span-1 row-span-1 group">
          <div className="card h-full bg-base-100 border border-base-200 shadow-lg hover:shadow-xl hover:scale-[1.02] transition-all duration-300">
            <div className="card-body p-6">
              <div className="flex items-center gap-3 mb-4">
                <Clock className="w-6 h-6 text-accent" />
                <h3 className="font-semibold">Обновления</h3>
              </div>
              <div className="space-y-2">
                <div className="flex items-center gap-2 text-sm">
                  <div className="w-2 h-2 bg-success rounded-full animate-pulse"></div>
                  <span className="text-base-content/70">
                    {stats?.newListingsLastHour || 10} новых за час
                  </span>
                </div>
                <div className="flex items-center gap-2 text-sm">
                  <div className="w-2 h-2 bg-warning rounded-full"></div>
                  <span className="text-base-content/70">
                    {stats?.priceDropsToday || 8} снижений цен
                  </span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </>
  );
};
