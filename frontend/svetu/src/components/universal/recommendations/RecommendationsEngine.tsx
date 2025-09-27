'use client';

import { FC, useState, useEffect, useMemo } from 'react';
import { useTranslations } from 'next-intl';
import UniversalListingCard from '../cards/UniversalListingCard';
import type { UniversalListingData } from '../cards/UniversalListingCard';
import {
  FaChartLine,
  FaHeart,
  FaHistory,
  FaLightbulb,
  FaFire,
  FaUserFriends,
} from 'react-icons/fa';

// Типы рекомендаций
export type RecommendationType =
  | 'similar' // Похожие товары
  | 'personal' // Персональные на основе истории
  | 'trending' // Популярные сейчас
  | 'new' // Новые поступления
  | 'discounted' // Со скидками
  | 'viewed_together' // Часто просматривают вместе
  | 'bought_together' // Часто покупают вместе
  | 'from_same_seller' // От того же продавца
  | 'in_same_category' // В той же категории
  | 'collaborative' // Коллаборативная фильтрация
  | 'content_based' // На основе контента
  | 'hybrid'; // Гибридные рекомендации

export interface RecommendationConfig {
  type: RecommendationType;
  title: string;
  icon?: React.ComponentType<any>;
  description?: string;
  maxItems?: number;
  algorithm?: 'simple' | 'ml' | 'collaborative' | 'content' | 'hybrid';
  params?: Record<string, any>;
  showReason?: boolean; // Показывать причину рекомендации
  refresh?: boolean; // Возможность обновить рекомендации
}

interface RecommendationsEngineProps {
  type: RecommendationType | RecommendationType[];
  category?: string;
  currentItemId?: number;
  userId?: number;
  limit?: number;
  layout?: 'grid' | 'carousel' | 'list';
  showTitle?: boolean;
  showDescription?: boolean;
  onItemClick?: (item: UniversalListingData) => void;
  className?: string;
  config?: Partial<RecommendationConfig>;
}

// Конфигурации по умолчанию для типов рекомендаций
const DEFAULT_CONFIGS: Record<RecommendationType, RecommendationConfig> = {
  similar: {
    type: 'similar',
    title: 'Similar Items',
    icon: FaLightbulb,
    description: 'Based on item characteristics',
    maxItems: 12,
    algorithm: 'content',
    showReason: true,
  },
  personal: {
    type: 'personal',
    title: 'Recommended for You',
    icon: FaHeart,
    description: 'Based on your browsing history',
    maxItems: 20,
    algorithm: 'hybrid',
    showReason: false,
    refresh: true,
  },
  trending: {
    type: 'trending',
    title: 'Trending Now',
    icon: FaFire,
    description: 'Popular items right now',
    maxItems: 10,
    algorithm: 'simple',
    refresh: true,
  },
  new: {
    type: 'new',
    title: 'New Arrivals',
    icon: FaChartLine,
    description: 'Recently added items',
    maxItems: 15,
    algorithm: 'simple',
  },
  discounted: {
    type: 'discounted',
    title: 'Special Offers',
    icon: FaChartLine,
    description: 'Items with discounts',
    maxItems: 12,
    algorithm: 'simple',
  },
  viewed_together: {
    type: 'viewed_together',
    title: 'Customers Also Viewed',
    icon: FaHistory,
    description: 'Often viewed together',
    maxItems: 8,
    algorithm: 'collaborative',
    showReason: true,
  },
  bought_together: {
    type: 'bought_together',
    title: 'Frequently Bought Together',
    icon: FaUserFriends,
    description: 'Items often purchased together',
    maxItems: 6,
    algorithm: 'collaborative',
  },
  from_same_seller: {
    type: 'from_same_seller',
    title: 'More from this Seller',
    icon: FaUserFriends,
    description: 'Other items from the same seller',
    maxItems: 10,
    algorithm: 'simple',
  },
  in_same_category: {
    type: 'in_same_category',
    title: 'More in this Category',
    icon: FaChartLine,
    description: 'Similar items in the same category',
    maxItems: 15,
    algorithm: 'simple',
  },
  collaborative: {
    type: 'collaborative',
    title: 'Users Like You Also Liked',
    icon: FaUserFriends,
    description: 'Based on similar users',
    maxItems: 12,
    algorithm: 'collaborative',
  },
  content_based: {
    type: 'content_based',
    title: 'Based on Your Interests',
    icon: FaLightbulb,
    description: 'Matching your preferences',
    maxItems: 15,
    algorithm: 'content',
  },
  hybrid: {
    type: 'hybrid',
    title: 'Top Picks for You',
    icon: FaHeart,
    description: 'Our best recommendations',
    maxItems: 20,
    algorithm: 'hybrid',
    refresh: true,
  },
};

const RecommendationsEngine: FC<RecommendationsEngineProps> = ({
  type,
  category,
  currentItemId,
  userId,
  limit = 12,
  layout = 'grid',
  showTitle = true,
  showDescription = false,
  onItemClick,
  className = '',
  config: customConfig,
}) => {
  const t = useTranslations('recommendations');
  const [recommendations, setRecommendations] = useState<
    UniversalListingData[]
  >([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [refreshKey, setRefreshKey] = useState(0);

  // Получаем конфигурацию
  const configs = useMemo(() => {
    const types = Array.isArray(type) ? type : [type];
    return types.map((t) => ({
      ...DEFAULT_CONFIGS[t],
      ...customConfig,
    }));
  }, [type, customConfig]);

  // Загрузка рекомендаций
  useEffect(() => {
    const fetchRecommendations = async () => {
      setLoading(true);
      setError(null);

      try {
        // Здесь будет вызов API для получения рекомендаций
        // Пока используем моковые данные
        const mockData = generateMockRecommendations(configs[0].type, limit);
        setRecommendations(mockData);
      } catch (err) {
        setError('Failed to load recommendations');
        console.error('Recommendations error:', err);
      } finally {
        setLoading(false);
      }
    };

    fetchRecommendations();
  }, [configs, limit, category, currentItemId, userId, refreshKey]);

  // Генерация моковых данных (заменить на реальный API)
  const generateMockRecommendations = (
    type: RecommendationType,
    count: number
  ): UniversalListingData[] => {
    const items: UniversalListingData[] = [];

    for (let i = 0; i < count; i++) {
      items.push({
        id: Math.floor(Math.random() * 10000),
        title: `${type} Item ${i + 1}`,
        price: Math.floor(Math.random() * 50000) + 5000,
        currency: '€',
        images: [`https://picsum.photos/400/300?random=${i}`],
        location: {
          city: ['Belgrade', 'Novi Sad', 'Niš', 'Kragujevac'][
            Math.floor(Math.random() * 4)
          ],
        },
        category: category || 'marketplace',
        createdAt: new Date(
          Date.now() - Math.random() * 30 * 24 * 60 * 60 * 1000
        ).toISOString(),
        badges:
          Math.random() > 0.7
            ? [
                {
                  type: Math.random() > 0.5 ? 'new' : 'discount',
                  label: Math.random() > 0.5 ? 'New' : '-15%',
                },
              ]
            : [],
        stats: {
          views: Math.floor(Math.random() * 1000),
          favorites: Math.floor(Math.random() * 50),
        },
      });
    }

    return items;
  };

  const handleRefresh = () => {
    setRefreshKey((prev) => prev + 1);
  };

  const renderRecommendations = () => {
    if (loading) {
      return (
        <div className="flex justify-center items-center h-48">
          <span className="loading loading-spinner loading-lg"></span>
        </div>
      );
    }

    if (error) {
      return (
        <div className="alert alert-error">
          <span>{error}</span>
        </div>
      );
    }

    if (recommendations.length === 0) {
      return (
        <div className="text-center py-8 text-base-content/60">
          {t('noRecommendations')}
        </div>
      );
    }

    if (layout === 'carousel') {
      return (
        <div className="carousel carousel-center space-x-4 p-4">
          {recommendations.map((item, index) => (
            <div key={item.id} className="carousel-item">
              <UniversalListingCard
                data={item}
                type={category as any}
                layout="grid"
                showBadges={true}
                showFavorite={true}
                showCompare={true}
                className="w-72"
              />
            </div>
          ))}
        </div>
      );
    }

    if (layout === 'list') {
      return (
        <div className="space-y-4">
          {recommendations.map((item) => (
            <UniversalListingCard
              key={item.id}
              data={item}
              type={category as any}
              layout="list"
              showBadges={true}
              showFavorite={true}
              showCompare={true}
            />
          ))}
        </div>
      );
    }

    // Grid layout
    return (
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
        {recommendations.map((item) => (
          <UniversalListingCard
            key={item.id}
            data={item}
            type={category as any}
            layout="grid"
            showBadges={true}
            showFavorite={true}
            showCompare={true}
          />
        ))}
      </div>
    );
  };

  // Рендер для множественных типов рекомендаций
  if (Array.isArray(type) && type.length > 1) {
    return (
      <div className={`space-y-8 ${className}`}>
        {configs.map((config, index) => (
          <div key={config.type} className="space-y-4">
            {showTitle && (
              <div className="flex items-center justify-between">
                <h3 className="text-xl font-semibold flex items-center gap-2">
                  {config.icon && <config.icon className="w-5 h-5" />}
                  {t(config.title) || config.title}
                </h3>
                {config.refresh && (
                  <button
                    className="btn btn-ghost btn-sm"
                    onClick={handleRefresh}
                  >
                    {t('refresh')}
                  </button>
                )}
              </div>
            )}

            {showDescription && config.description && (
              <p className="text-base-content/70">
                {t(config.description) || config.description}
              </p>
            )}

            {index === 0 ? (
              renderRecommendations()
            ) : (
              <div className="text-center py-8 text-base-content/60">
                {/* Placeholder для других типов рекомендаций */}
                {t('loadingMore')}...
              </div>
            )}
          </div>
        ))}
      </div>
    );
  }

  // Рендер для одного типа рекомендаций
  const config = configs[0];

  return (
    <div className={`space-y-4 ${className}`}>
      {showTitle && (
        <div className="flex items-center justify-between">
          <h3 className="text-xl font-semibold flex items-center gap-2">
            {config.icon && <config.icon className="w-5 h-5" />}
            {t(config.title) || config.title}
          </h3>

          <div className="flex items-center gap-2">
            {config.refresh && (
              <button className="btn btn-ghost btn-sm" onClick={handleRefresh}>
                {t('refresh')}
              </button>
            )}

            {recommendations.length > 0 && (
              <span className="text-sm text-base-content/60">
                {recommendations.length} {t('items')}
              </span>
            )}
          </div>
        </div>
      )}

      {showDescription && config.description && (
        <p className="text-base-content/70">
          {t(config.description) || config.description}
        </p>
      )}

      {renderRecommendations()}
    </div>
  );
};

export default RecommendationsEngine;
