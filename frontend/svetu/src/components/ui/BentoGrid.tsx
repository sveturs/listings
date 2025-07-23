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
  ArrowRight,
} from 'lucide-react';

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
  };
}

export const BentoGrid: React.FC<BentoGridProps> = ({
  categories = [],
  featuredListing,
  stats,
}) => {
  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 p-4">
      {/* Hero Card - Большая карточка */}
      <div className="col-span-1 md:col-span-2 lg:col-span-2 row-span-2 group">
        <div className="card h-full bg-gradient-to-br from-primary/10 to-secondary/10 border-0 shadow-lg hover:shadow-xl transition-all duration-300 overflow-hidden">
          <div className="card-body p-8 flex flex-col justify-between h-full">
            <div>
              <div className="flex items-center gap-2 mb-4">
                <Sparkles className="w-6 h-6 text-primary" />
                <span className="badge badge-primary">Новое</span>
              </div>
              <h2 className="card-title text-3xl font-bold mb-4">
                Найдите то, что нужно именно вам
              </h2>
              <p className="text-base-content/70 mb-6">
                Тысячи товаров и услуг от проверенных продавцов в вашем районе
              </p>
            </div>
            <div className="flex justify-between items-end">
              <Link
                href="/search"
                className="btn btn-primary group-hover:scale-105 transition-transform"
              >
                Начать поиск
                <ArrowRight className="w-4 h-4 ml-2" />
              </Link>
              <div className="text-right">
                <p className="text-sm text-base-content/60">
                  Активных объявлений
                </p>
                <p className="text-2xl font-bold text-primary">
                  {stats?.totalListings?.toLocaleString() || '10,000+'}
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Популярные категории */}
      <div className="col-span-1 row-span-1 group">
        <div className="card h-full bg-base-100 shadow-lg hover:shadow-xl transition-all duration-300">
          <div className="card-body p-6">
            <div className="flex items-center justify-between mb-4">
              <TrendingUp className="w-6 h-6 text-warning" />
              <span className="text-xs text-base-content/60">Обновлено</span>
            </div>
            <h3 className="font-semibold mb-3">Популярные категории</h3>
            <div className="space-y-2">
              {categories.slice(0, 3).map((cat) => (
                <Link
                  key={cat.id}
                  href={`/search?category=${cat.id}`}
                  className="flex items-center justify-between p-2 rounded-lg hover:bg-base-200 transition-colors"
                >
                  <span className="text-sm">{cat.name}</span>
                  <span className="text-xs text-base-content/60">
                    {cat.count}
                  </span>
                </Link>
              ))}
            </div>
          </div>
        </div>
      </div>

      {/* Безопасные сделки */}
      <div className="col-span-1 row-span-1 group">
        <div className="card h-full bg-success/10 shadow-lg hover:shadow-xl transition-all duration-300">
          <div className="card-body p-6">
            <div className="flex items-center gap-3 mb-4">
              <div className="p-3 bg-success/20 rounded-full">
                <Shield className="w-6 h-6 text-success" />
              </div>
              <div className="flex-1">
                <h3 className="font-semibold">Безопасные сделки</h3>
                <p className="text-xs text-base-content/60">
                  Защита покупателей
                </p>
              </div>
            </div>
            <div className="stats stats-vertical shadow-none bg-transparent p-0">
              <div className="stat p-0">
                <div className="stat-value text-2xl text-success">98%</div>
                <div className="stat-desc">успешных сделок</div>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Избранное объявление */}
      {featuredListing && (
        <div className="col-span-1 md:col-span-2 lg:col-span-1 row-span-1 group">
          <Link href={`/listing/${featuredListing.id}`}>
            <div className="card h-full bg-base-100 shadow-lg hover:shadow-xl transition-all duration-300 overflow-hidden">
              <figure className="h-32 relative">
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
      )}

      {/* Статистика пользователей */}
      <div className="col-span-1 row-span-1 group">
        <div className="card h-full bg-info/10 shadow-lg hover:shadow-xl transition-all duration-300">
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
                  {stats?.activeUsers?.toLocaleString() || '1,500+'}
                </span>
              </div>
              <div className="flex justify-between items-center">
                <span className="text-sm text-base-content/70">
                  Новых сегодня
                </span>
                <span className="font-bold text-info">+23</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Рядом с вами */}
      <div className="col-span-1 row-span-1 group">
        <div className="card h-full bg-secondary/10 shadow-lg hover:shadow-xl transition-all duration-300">
          <div className="card-body p-6">
            <div className="flex items-center gap-3 mb-4">
              <MapPin className="w-6 h-6 text-secondary" />
              <div>
                <h3 className="font-semibold">Рядом с вами</h3>
                <p className="text-xs text-base-content/60">В радиусе 5 км</p>
              </div>
            </div>
            <div className="text-center pt-2">
              <p className="text-3xl font-bold text-secondary">342</p>
              <p className="text-sm text-base-content/60">объявления</p>
            </div>
          </div>
        </div>
      </div>

      {/* Последние обновления */}
      <div className="col-span-1 row-span-1 group">
        <div className="card h-full bg-base-100 shadow-lg hover:shadow-xl transition-all duration-300">
          <div className="card-body p-6">
            <div className="flex items-center gap-3 mb-4">
              <Clock className="w-6 h-6 text-accent" />
              <h3 className="font-semibold">Обновления</h3>
            </div>
            <div className="space-y-2">
              <div className="flex items-center gap-2 text-sm">
                <div className="w-2 h-2 bg-success rounded-full animate-pulse"></div>
                <span className="text-base-content/70">15 новых за час</span>
              </div>
              <div className="flex items-center gap-2 text-sm">
                <div className="w-2 h-2 bg-warning rounded-full"></div>
                <span className="text-base-content/70">8 снижений цен</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};
