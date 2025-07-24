# Примеры кода UI/UX для Sve Tu

## 1. Улучшенная карточка товара с DaisyUI и Tailwind v4

```tsx
// src/components/marketplace/EnhancedListingCard.tsx
import React from 'react';
import Link from 'next/link';
import Image from 'next/image';
import { Heart, MapPin, Shield, Clock, Eye } from 'lucide-react';
import { formatDistanceToNow } from 'date-fns';
import { ru } from 'date-fns/locale';

interface EnhancedListingCardProps {
  listing: {
    id: string;
    title: string;
    price: number;
    currency: string;
    images: string[];
    location: {
      city: string;
      distance?: number;
    };
    seller: {
      name: string;
      rating: number;
      verified: boolean;
    };
    createdAt: string;
    viewCount: number;
    isEscrowEnabled: boolean;
    isFavorite?: boolean;
    category: {
      id: string;
      name: string;
      icon?: string;
    };
    condition?: 'new' | 'like_new' | 'good' | 'fair';
    ecoScore?: number;
  };
  onToggleFavorite?: (id: string) => void;
}

export const EnhancedListingCard: React.FC<EnhancedListingCardProps> = ({ 
  listing, 
  onToggleFavorite 
}) => {
  const conditionBadge = {
    new: { text: 'Новое', class: 'badge-success' },
    like_new: { text: 'Как новое', class: 'badge-info' },
    good: { text: 'Хорошее', class: 'badge-primary' },
    fair: { text: 'Удовлетворительное', class: 'badge-warning' }
  };

  return (
    <div className="card card-compact bg-base-100 shadow-sm hover:shadow-lg transition-all duration-300 group">
      {/* Изображение с оверлеями */}
      <figure className="relative aspect-square overflow-hidden">
        <Image
          src={listing.images[0]}
          alt={listing.title}
          fill
          className="object-cover group-hover:scale-105 transition-transform duration-300"
          sizes="(max-width: 768px) 50vw, (max-width: 1200px) 33vw, 25vw"
        />
        
        {/* Badges слева сверху */}
        <div className="absolute top-2 left-2 flex flex-col gap-1">
          {listing.condition && (
            <div className={`badge badge-sm ${conditionBadge[listing.condition].class}`}>
              {conditionBadge[listing.condition].text}
            </div>
          )}
          {listing.ecoScore && listing.ecoScore > 7 && (
            <div className="badge badge-sm badge-success gap-1">
              <span className="text-xs">♻️</span>
              <span>Эко</span>
            </div>
          )}
        </div>

        {/* Действия справа сверху */}
        <div className="absolute top-2 right-2">
          <button
            onClick={(e) => {
              e.preventDefault();
              onToggleFavorite?.(listing.id);
            }}
            className="btn btn-circle btn-sm btn-ghost bg-base-100/80 backdrop-blur-sm"
          >
            <Heart 
              className={`w-4 h-4 ${listing.isFavorite ? 'fill-error text-error' : ''}`} 
            />
          </button>
        </div>

        {/* Расстояние */}
        {listing.location.distance && (
          <div className="absolute bottom-2 left-2">
            <div className="badge badge-neutral badge-sm gap-1">
              <MapPin className="w-3 h-3" />
              <span>{listing.location.distance} км</span>
            </div>
          </div>
        )}

        {/* Количество фото */}
        {listing.images.length > 1 && (
          <div className="absolute bottom-2 right-2">
            <div className="badge badge-neutral badge-sm">
              {listing.images.length} фото
            </div>
          </div>
        )}
      </figure>

      <div className="card-body p-3">
        {/* Категория */}
        <div className="text-xs text-base-content/60 flex items-center gap-1">
          {listing.category.icon && <span>{listing.category.icon}</span>}
          <span>{listing.category.name}</span>
        </div>

        {/* Заголовок */}
        <h3 className="card-title text-sm line-clamp-2 min-h-[2.5rem]">
          {listing.title}
        </h3>

        {/* Продавец */}
        <div className="flex items-center gap-2 text-xs">
          <div className="avatar placeholder">
            <div className="bg-neutral text-neutral-content rounded-full w-5">
              <span className="text-xs">{listing.seller.name[0]}</span>
            </div>
          </div>
          <span className="font-medium">{listing.seller.name}</span>
          {listing.seller.verified && (
            <Shield className="w-3 h-3 text-success" />
          )}
          <div className="rating rating-xs">
            <input 
              type="radio" 
              className="rating-hidden" 
              checked={listing.seller.rating === 0}
              readOnly
            />
            {[1, 2, 3, 4, 5].map((star) => (
              <input
                key={star}
                type="radio"
                className="mask mask-star-2 bg-warning"
                checked={listing.seller.rating >= star}
                readOnly
              />
            ))}
          </div>
        </div>

        {/* Статистика */}
        <div className="flex items-center gap-3 text-xs text-base-content/60">
          <span className="flex items-center gap-1">
            <Clock className="w-3 h-3" />
            {formatDistanceToNow(new Date(listing.createdAt), { 
              addSuffix: true,
              locale: ru 
            })}
          </span>
          <span className="flex items-center gap-1">
            <Eye className="w-3 h-3" />
            {listing.viewCount}
          </span>
        </div>

        {/* Цена и действия */}
        <div className="card-actions justify-between items-end mt-2">
          <div>
            <div className="text-lg font-bold">
              {listing.price.toLocaleString('ru-RU')} {listing.currency}
            </div>
            {listing.isEscrowEnabled && (
              <div className="text-xs text-success flex items-center gap-1">
                <Shield className="w-3 h-3" />
                Безопасная сделка
              </div>
            )}
          </div>
          <Link 
            href={`/marketplace/listing/${listing.id}`}
            className="btn btn-primary btn-sm"
          >
            Подробнее
          </Link>
        </div>
      </div>
    </div>
  );
};
```

## 2. Интерактивная карта с кластеризацией

```tsx
// src/components/maps/InteractiveListingMap.tsx
import React, { useCallback, useState } from 'react';
import { GoogleMap, useJsApiLoader, MarkerClusterer, Marker, InfoWindow } from '@react-google-maps/api';
import { Listing } from '@/types/marketplace';

interface InteractiveListingMapProps {
  listings: Listing[];
  center?: { lat: number; lng: number };
  zoom?: number;
  onListingClick?: (listing: Listing) => void;
  height?: string;
}

const mapContainerStyle = {
  width: '100%',
  height: '100%'
};

const options = {
  styles: [
    {
      featureType: 'poi',
      elementType: 'labels',
      stylers: [{ visibility: 'off' }]
    }
  ],
  disableDefaultUI: false,
  zoomControl: true,
  streetViewControl: false,
  mapTypeControl: false,
  fullscreenControl: true
};

export const InteractiveListingMap: React.FC<InteractiveListingMapProps> = ({
  listings,
  center = { lat: 44.787197, lng: 20.457273 }, // Белград
  zoom = 12,
  onListingClick,
  height = '400px'
}) => {
  const { isLoaded, loadError } = useJsApiLoader({
    googleMapsApiKey: process.env.NEXT_PUBLIC_GOOGLE_MAPS_KEY!,
    libraries: ['places']
  });

  const [selectedListing, setSelectedListing] = useState<Listing | null>(null);
  const [map, setMap] = useState<google.maps.Map | null>(null);

  const onLoad = useCallback((map: google.maps.Map) => {
    setMap(map);
  }, []);

  const onUnmount = useCallback(() => {
    setMap(null);
  }, []);

  const handleMarkerClick = (listing: Listing) => {
    setSelectedListing(listing);
    onListingClick?.(listing);
  };

  if (loadError) {
    return (
      <div className="alert alert-error">
        <span>Ошибка загрузки карты</span>
      </div>
    );
  }

  if (!isLoaded) {
    return (
      <div className="flex items-center justify-center" style={{ height }}>
        <span className="loading loading-spinner loading-lg"></span>
      </div>
    );
  }

  return (
    <div className="relative rounded-lg overflow-hidden shadow-md" style={{ height }}>
      <GoogleMap
        mapContainerStyle={mapContainerStyle}
        center={center}
        zoom={zoom}
        options={options}
        onLoad={onLoad}
        onUnmount={onUnmount}
      >
        <MarkerClusterer
          options={{
            imagePath: 'https://developers.google.com/maps/documentation/javascript/examples/markerclusterer/m',
            gridSize: 50,
            maxZoom: 15
          }}
        >
          {(clusterer) =>
            listings.map((listing) => (
              <Marker
                key={listing.id}
                position={{
                  lat: listing.location.coordinates.lat,
                  lng: listing.location.coordinates.lng
                }}
                clusterer={clusterer}
                onClick={() => handleMarkerClick(listing)}
                icon={{
                  url: '/icons/marker.svg',
                  scaledSize: new google.maps.Size(30, 30)
                }}
              />
            ))
          }
        </MarkerClusterer>

        {selectedListing && (
          <InfoWindow
            position={{
              lat: selectedListing.location.coordinates.lat,
              lng: selectedListing.location.coordinates.lng
            }}
            onCloseClick={() => setSelectedListing(null)}
          >
            <div className="p-2 max-w-xs">
              <img 
                src={selectedListing.images[0]} 
                alt={selectedListing.title}
                className="w-full h-24 object-cover rounded mb-2"
              />
              <h3 className="font-bold text-sm">{selectedListing.title}</h3>
              <p className="text-lg font-bold text-primary">
                {selectedListing.price.toLocaleString('ru-RU')} RSD
              </p>
              <button 
                onClick={() => onListingClick?.(selectedListing)}
                className="btn btn-primary btn-xs mt-2 w-full"
              >
                Посмотреть
              </button>
            </div>
          </InfoWindow>
        )}
      </GoogleMap>

      {/* Легенда карты */}
      <div className="absolute bottom-4 left-4 bg-base-100/90 backdrop-blur-sm rounded-lg p-3 shadow-lg">
        <h4 className="font-bold text-sm mb-2">Товары рядом</h4>
        <div className="text-xs space-y-1">
          <div className="flex items-center gap-2">
            <div className="w-3 h-3 bg-primary rounded-full"></div>
            <span>Доступные товары</span>
          </div>
          <div className="flex items-center gap-2">
            <div className="w-6 h-6 bg-warning rounded-full flex items-center justify-center text-xs font-bold">
              5
            </div>
            <span>Группа товаров</span>
          </div>
        </div>
      </div>

      {/* Фильтр расстояния */}
      <div className="absolute top-4 right-4">
        <select className="select select-sm bg-base-100/90 backdrop-blur-sm">
          <option>В радиусе 5 км</option>
          <option>В радиусе 10 км</option>
          <option>В радиусе 25 км</option>
          <option>Весь город</option>
        </select>
      </div>
    </div>
  );
};
```

## 3. Мобильная навигация с анимацией

```tsx
// src/components/navigation/MobileBottomNav.tsx
import React from 'react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { Home, Search, PlusCircle, MessageCircle, User } from 'lucide-react';
import { motion } from 'framer-motion';

interface NavItem {
  icon: React.ElementType;
  label: string;
  href: string;
  badge?: number;
}

export const MobileBottomNav: React.FC = () => {
  const pathname = usePathname();
  
  const navItems: NavItem[] = [
    { icon: Home, label: 'Главная', href: '/' },
    { icon: Search, label: 'Поиск', href: '/search' },
    { icon: PlusCircle, label: 'Создать', href: '/create' },
    { icon: MessageCircle, label: 'Чаты', href: '/chats', badge: 3 },
    { icon: User, label: 'Профиль', href: '/profile' }
  ];

  return (
    <nav className="btm-nav btm-nav-sm md:hidden bg-base-100 border-t border-base-200">
      {navItems.map((item) => {
        const isActive = pathname === item.href;
        const Icon = item.icon;
        
        return (
          <Link
            key={item.href}
            href={item.href}
            className={`relative ${isActive ? 'active' : ''}`}
          >
            <motion.div
              whileTap={{ scale: 0.9 }}
              className="flex flex-col items-center"
            >
              {/* Индикатор активности */}
              {isActive && (
                <motion.div
                  layoutId="activeTab"
                  className="absolute -top-0.5 w-12 h-1 bg-primary rounded-b-full"
                  transition={{ type: "spring", stiffness: 500, damping: 30 }}
                />
              )}
              
              {/* Иконка с бейджем */}
              <div className="relative">
                <Icon 
                  className={`w-5 h-5 ${isActive ? 'text-primary' : 'text-base-content/60'}`} 
                />
                {item.badge && item.badge > 0 && (
                  <span className="badge badge-error badge-xs absolute -top-1 -right-2">
                    {item.badge > 9 ? '9+' : item.badge}
                  </span>
                )}
              </div>
              
              {/* Подпись */}
              <span className={`text-xs mt-1 ${isActive ? 'text-primary' : 'text-base-content/60'}`}>
                {item.label}
              </span>
            </motion.div>
          </Link>
        );
      })}
    </nav>
  );
};
```

## 4. Компонент эскроу процесса

```tsx
// src/components/escrow/EscrowFlowVisualizer.tsx
import React, { useState } from 'react';
import { CreditCard, Lock, Package, CheckCircle, AlertCircle, X } from 'lucide-react';
import { motion, AnimatePresence } from 'framer-motion';

interface EscrowStep {
  id: string;
  title: string;
  description: string;
  icon: React.ElementType;
  status: 'pending' | 'active' | 'completed' | 'error';
}

interface EscrowFlowVisualizerProps {
  currentStep: string;
  amount: number;
  onCancel?: () => void;
  onRefund?: () => void;
}

export const EscrowFlowVisualizer: React.FC<EscrowFlowVisualizerProps> = ({
  currentStep,
  amount,
  onCancel,
  onRefund
}) => {
  const [showDetails, setShowDetails] = useState(false);

  const steps: EscrowStep[] = [
    {
      id: 'payment',
      title: 'Оплата',
      description: 'Покупатель оплачивает товар',
      icon: CreditCard,
      status: currentStep === 'payment' ? 'active' : 
              ['escrow', 'delivery', 'complete'].includes(currentStep) ? 'completed' : 'pending'
    },
    {
      id: 'escrow',
      title: 'Блокировка средств',
      description: 'Деньги заблокированы до получения товара',
      icon: Lock,
      status: currentStep === 'escrow' ? 'active' : 
              ['delivery', 'complete'].includes(currentStep) ? 'completed' : 'pending'
    },
    {
      id: 'delivery',
      title: 'Доставка',
      description: 'Товар отправлен покупателю',
      icon: Package,
      status: currentStep === 'delivery' ? 'active' : 
              currentStep === 'complete' ? 'completed' : 'pending'
    },
    {
      id: 'complete',
      title: 'Завершение',
      description: 'Средства переведены продавцу',
      icon: CheckCircle,
      status: currentStep === 'complete' ? 'completed' : 'pending'
    }
  ];

  const getStepColor = (status: string) => {
    switch (status) {
      case 'completed': return 'text-success bg-success/20';
      case 'active': return 'text-primary bg-primary/20 animate-pulse';
      case 'error': return 'text-error bg-error/20';
      default: return 'text-base-content/40 bg-base-200';
    }
  };

  return (
    <div className="card bg-base-100 shadow-lg">
      <div className="card-body">
        <div className="flex justify-between items-start mb-4">
          <div>
            <h3 className="card-title">Безопасная сделка</h3>
            <p className="text-sm text-base-content/60">
              Сумма: <span className="font-bold text-primary">{amount.toLocaleString('ru-RU')} RSD</span>
            </p>
          </div>
          <button 
            onClick={() => setShowDetails(!showDetails)}
            className="btn btn-ghost btn-sm"
          >
            {showDetails ? <X className="w-4 h-4" /> : <AlertCircle className="w-4 h-4" />}
          </button>
        </div>

        {/* Визуализация шагов */}
        <div className="relative">
          {/* Линия прогресса */}
          <div className="absolute top-8 left-8 right-8 h-0.5 bg-base-300">
            <motion.div 
              className="h-full bg-primary"
              initial={{ width: '0%' }}
              animate={{ 
                width: `${(steps.findIndex(s => s.id === currentStep) / (steps.length - 1)) * 100}%` 
              }}
              transition={{ duration: 0.5 }}
            />
          </div>

          {/* Шаги */}
          <div className="relative grid grid-cols-4 gap-2">
            {steps.map((step, index) => {
              const Icon = step.icon;
              return (
                <motion.div
                  key={step.id}
                  className="flex flex-col items-center"
                  initial={{ opacity: 0, y: 20 }}
                  animate={{ opacity: 1, y: 0 }}
                  transition={{ delay: index * 0.1 }}
                >
                  <div className={`
                    w-16 h-16 rounded-full flex items-center justify-center
                    ${getStepColor(step.status)}
                    transition-all duration-300
                  `}>
                    <Icon className="w-6 h-6" />
                  </div>
                  <p className="text-xs font-medium mt-2 text-center">
                    {step.title}
                  </p>
                </motion.div>
              );
            })}
          </div>
        </div>

        {/* Детали */}
        <AnimatePresence>
          {showDetails && (
            <motion.div
              initial={{ height: 0, opacity: 0 }}
              animate={{ height: 'auto', opacity: 1 }}
              exit={{ height: 0, opacity: 0 }}
              transition={{ duration: 0.3 }}
              className="overflow-hidden"
            >
              <div className="divider"></div>
              
              <div className="space-y-3">
                {steps.map((step) => (
                  <div 
                    key={step.id}
                    className={`flex items-start gap-3 p-3 rounded-lg ${
                      step.status === 'active' ? 'bg-primary/10' : ''
                    }`}
                  >
                    <step.icon className={`w-5 h-5 mt-0.5 ${
                      step.status === 'completed' ? 'text-success' :
                      step.status === 'active' ? 'text-primary' :
                      'text-base-content/40'
                    }`} />
                    <div className="flex-1">
                      <p className="font-medium text-sm">{step.title}</p>
                      <p className="text-xs text-base-content/60">{step.description}</p>
                    </div>
                    {step.status === 'completed' && (
                      <CheckCircle className="w-4 h-4 text-success" />
                    )}
                  </div>
                ))}
              </div>

              <div className="alert alert-info mt-4">
                <AlertCircle className="w-4 h-4" />
                <span className="text-sm">
                  Ваши средства защищены. В случае проблем с товаром вы можете открыть спор 
                  и вернуть деньги.
                </span>
              </div>

              {/* Действия */}
              {currentStep !== 'complete' && (
                <div className="flex gap-2 mt-4">
                  <button 
                    onClick={onCancel}
                    className="btn btn-error btn-sm"
                  >
                    Отменить сделку
                  </button>
                  <button 
                    onClick={onRefund}
                    className="btn btn-warning btn-sm"
                  >
                    Открыть спор
                  </button>
                </div>
              )}
            </motion.div>
          )}
        </AnimatePresence>
      </div>
    </div>
  );
};
```

## 5. Компонент устойчивого развития

```tsx
// src/components/sustainability/EcoImpactDashboard.tsx
import React from 'react';
import { Leaf, Droplets, Recycle, TrendingUp } from 'lucide-react';
import { motion } from 'framer-motion';

interface EcoStat {
  label: string;
  value: number;
  unit: string;
  icon: React.ElementType;
  color: string;
  comparison?: string;
}

interface EcoImpactDashboardProps {
  userStats: {
    co2Saved: number;
    waterSaved: number;
    itemsRecycled: number;
    moneyEarned: number;
  };
  globalStats?: {
    totalCo2Saved: number;
    totalUsers: number;
  };
}

export const EcoImpactDashboard: React.FC<EcoImpactDashboardProps> = ({
  userStats,
  globalStats
}) => {
  const stats: EcoStat[] = [
    {
      label: 'CO₂ сэкономлено',
      value: userStats.co2Saved,
      unit: 'кг',
      icon: Leaf,
      color: 'text-success',
      comparison: 'Как посадить 5 деревьев'
    },
    {
      label: 'Воды сохранено',
      value: userStats.waterSaved,
      unit: 'л',
      icon: Droplets,
      color: 'text-info',
      comparison: '10 полных ванн'
    },
    {
      label: 'Вещей переиспользовано',
      value: userStats.itemsRecycled,
      unit: 'шт',
      icon: Recycle,
      color: 'text-warning',
      comparison: undefined
    },
    {
      label: 'Заработано',
      value: userStats.moneyEarned,
      unit: 'RSD',
      icon: TrendingUp,
      color: 'text-primary',
      comparison: undefined
    }
  ];

  return (
    <div className="space-y-6">
      {/* Личная статистика */}
      <div>
        <h2 className="text-2xl font-bold mb-4">Ваш вклад в экологию</h2>
        <div className="grid grid-cols-2 lg:grid-cols-4 gap-4">
          {stats.map((stat, index) => {
            const Icon = stat.icon;
            return (
              <motion.div
                key={stat.label}
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: index * 0.1 }}
                className="card bg-base-100 shadow-sm hover:shadow-md transition-shadow"
              >
                <div className="card-body p-4">
                  <div className="flex items-start justify-between">
                    <Icon className={`w-8 h-8 ${stat.color}`} />
                    <div className="text-right">
                      <p className="text-2xl font-bold">
                        {stat.value.toLocaleString('ru-RU')}
                      </p>
                      <p className="text-sm text-base-content/60">{stat.unit}</p>
                    </div>
                  </div>
                  <p className="text-xs font-medium mt-2">{stat.label}</p>
                  {stat.comparison && (
                    <p className="text-xs text-base-content/60">{stat.comparison}</p>
                  )}
                </div>
              </motion.div>
            );
          })}
        </div>
      </div>

      {/* Достижения */}
      <div className="card bg-gradient-to-r from-success/10 to-primary/10">
        <div className="card-body">
          <h3 className="card-title">Ваши эко-достижения</h3>
          <div className="flex flex-wrap gap-2 mt-4">
            {[
              { emoji: '🌱', name: 'Новичок', unlocked: true },
              { emoji: '🌿', name: 'Эко-воин', unlocked: true },
              { emoji: '🌳', name: 'Защитник природы', unlocked: true },
              { emoji: '🌍', name: 'Планетарный герой', unlocked: false },
              { emoji: '♻️', name: 'Мастер переработки', unlocked: false },
            ].map((achievement) => (
              <div
                key={achievement.name}
                className={`
                  badge badge-lg p-4 gap-2
                  ${achievement.unlocked ? 'badge-primary' : 'badge-ghost opacity-50'}
                `}
              >
                <span className="text-2xl">{achievement.emoji}</span>
                <span>{achievement.name}</span>
              </div>
            ))}
          </div>
        </div>
      </div>

      {/* Глобальная статистика */}
      {globalStats && (
        <div className="card bg-base-200">
          <div className="card-body text-center">
            <h3 className="card-title justify-center">Вместе мы сохранили</h3>
            <div className="stats stats-vertical lg:stats-horizontal shadow">
              <div className="stat">
                <div className="stat-figure text-success">
                  <Leaf className="w-8 h-8" />
                </div>
                <div className="stat-title">Общий CO₂</div>
                <div className="stat-value text-success">
                  {(globalStats.totalCo2Saved / 1000).toFixed(1)}т
                </div>
                <div className="stat-desc">тонн углекислого газа</div>
              </div>
              
              <div className="stat">
                <div className="stat-figure text-primary">
                  <svg className="w-8 h-8 fill-current" viewBox="0 0 24 24">
                    <path d="M12,2A10,10 0 0,0 2,12A10,10 0 0,0 12,22A10,10 0 0,0 22,12A10,10 0 0,0 12,2Z" />
                  </svg>
                </div>
                <div className="stat-title">Участников</div>
                <div className="stat-value text-primary">
                  {globalStats.totalUsers.toLocaleString('ru-RU')}
                </div>
                <div className="stat-desc">активных пользователей</div>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};
```

## 6. PWA конфигурация

```javascript
// public/sw.js
const CACHE_NAME = 'svetu-v1';
const urlsToCache = [
  '/',
  '/manifest.json',
  '/icons/icon-192x192.png',
  '/icons/icon-512x512.png',
  '/offline.html'
];

// Установка Service Worker
self.addEventListener('install', event => {
  event.waitUntil(
    caches.open(CACHE_NAME)
      .then(cache => cache.addAll(urlsToCache))
  );
});

// Активация и очистка старых кэшей
self.addEventListener('activate', event => {
  event.waitUntil(
    caches.keys().then(cacheNames => {
      return Promise.all(
        cacheNames
          .filter(cacheName => cacheName !== CACHE_NAME)
          .map(cacheName => caches.delete(cacheName))
      );
    })
  );
});

// Стратегия Network First для API
self.addEventListener('fetch', event => {
  if (event.request.url.includes('/api/')) {
    event.respondWith(
      fetch(event.request)
        .then(response => {
          const responseClone = response.clone();
          caches.open(CACHE_NAME).then(cache => {
            cache.put(event.request, responseClone);
          });
          return response;
        })
        .catch(() => caches.match(event.request))
    );
  } else {
    // Cache First для статики
    event.respondWith(
      caches.match(event.request)
        .then(response => response || fetch(event.request))
        .catch(() => caches.match('/offline.html'))
    );
  }
});

// Background Sync для офлайн действий
self.addEventListener('sync', event => {
  if (event.tag === 'sync-favorites') {
    event.waitUntil(syncFavorites());
  }
});

async function syncFavorites() {
  const db = await openDB();
  const favorites = await db.getAll('pending-favorites');
  
  for (const favorite of favorites) {
    try {
      await fetch('/api/v1/favorites', {
        method: 'POST',
        body: JSON.stringify(favorite),
        headers: { 'Content-Type': 'application/json' }
      });
      await db.delete('pending-favorites', favorite.id);
    } catch (error) {
      console.error('Sync failed:', error);
    }
  }
}
```

```json
// public/manifest.json
{
  "name": "Sve Tu - Локальный маркетплейс",
  "short_name": "Sve Tu",
  "description": "Покупайте и продавайте товары в вашем районе",
  "start_url": "/",
  "display": "standalone",
  "background_color": "#ffffff",
  "theme_color": "#570df8",
  "orientation": "portrait",
  "icons": [
    {
      "src": "/icons/icon-192x192.png",
      "sizes": "192x192",
      "type": "image/png",
      "purpose": "any maskable"
    },
    {
      "src": "/icons/icon-512x512.png",
      "sizes": "512x512",
      "type": "image/png"
    }
  ],
  "categories": ["shopping", "lifestyle"],
  "screenshots": [
    {
      "src": "/screenshots/home.png",
      "sizes": "412x915",
      "type": "image/png"
    },
    {
      "src": "/screenshots/map.png",
      "sizes": "412x915",
      "type": "image/png"
    }
  ],
  "related_applications": [],
  "prefer_related_applications": false
}
```

## 7. Оптимизированный поиск с AI

```tsx
// src/components/search/AIEnhancedSearch.tsx
import React, { useState, useCallback, useEffect } from 'react';
import { Search, Sparkles, MapPin, Filter, X } from 'lucide-react';
import { useDebounce } from '@/hooks/useDebounce';
import { motion, AnimatePresence } from 'framer-motion';

interface AIEnhancedSearchProps {
  onSearch: (query: string, filters: any) => void;
  suggestions?: string[];
  categories: Array<{ id: string; name: string; icon: string }>;
}

export const AIEnhancedSearch: React.FC<AIEnhancedSearchProps> = ({
  onSearch,
  suggestions = [],
  categories
}) => {
  const [query, setQuery] = useState('');
  const [showFilters, setShowFilters] = useState(false);
  const [showSuggestions, setShowSuggestions] = useState(false);
  const [filters, setFilters] = useState({
    category: '',
    priceMin: '',
    priceMax: '',
    radius: 10,
    condition: '',
    escrowOnly: false
  });
  
  const debouncedQuery = useDebounce(query, 300);

  useEffect(() => {
    if (debouncedQuery) {
      onSearch(debouncedQuery, filters);
    }
  }, [debouncedQuery, filters, onSearch]);

  const handleSuggestionClick = (suggestion: string) => {
    setQuery(suggestion);
    setShowSuggestions(false);
  };

  const clearFilters = () => {
    setFilters({
      category: '',
      priceMin: '',
      priceMax: '',
      radius: 10,
      condition: '',
      escrowOnly: false
    });
  };

  const activeFiltersCount = Object.values(filters).filter(
    v => v !== '' && v !== false && v !== 10
  ).length;

  return (
    <div className="relative w-full">
      {/* Основная строка поиска */}
      <div className="relative">
        <input
          type="text"
          value={query}
          onChange={(e) => {
            setQuery(e.target.value);
            setShowSuggestions(true);
          }}
          onFocus={() => setShowSuggestions(true)}
          placeholder="Поиск товаров с AI подсказками..."
          className="input input-bordered w-full pl-10 pr-20"
        />
        <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-base-content/60" />
        
        <div className="absolute right-2 top-1/2 -translate-y-1/2 flex items-center gap-1">
          <div className="badge badge-primary badge-sm gap-1">
            <Sparkles className="w-3 h-3" />
            AI
          </div>
          <button
            onClick={() => setShowFilters(!showFilters)}
            className={`btn btn-ghost btn-sm btn-circle ${activeFiltersCount > 0 ? 'text-primary' : ''}`}
          >
            <Filter className="w-4 h-4" />
            {activeFiltersCount > 0 && (
              <span className="badge badge-error badge-xs absolute -top-1 -right-1">
                {activeFiltersCount}
              </span>
            )}
          </button>
        </div>
      </div>

      {/* AI предложения */}
      <AnimatePresence>
        {showSuggestions && suggestions.length > 0 && (
          <motion.div
            initial={{ opacity: 0, y: -10 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -10 }}
            className="absolute top-full left-0 right-0 mt-1 bg-base-100 rounded-lg shadow-lg border border-base-200 z-50"
          >
            <div className="p-2">
              <div className="text-xs text-base-content/60 px-2 py-1 flex items-center gap-1">
                <Sparkles className="w-3 h-3" />
                AI предложения
              </div>
              {suggestions.map((suggestion, idx) => (
                <button
                  key={idx}
                  onClick={() => handleSuggestionClick(suggestion)}
                  className="w-full text-left px-3 py-2 hover:bg-base-200 rounded-md text-sm transition-colors"
                >
                  {suggestion}
                </button>
              ))}
            </div>
          </motion.div>
        )}
      </AnimatePresence>

      {/* Панель фильтров */}
      <AnimatePresence>
        {showFilters && (
          <motion.div
            initial={{ opacity: 0, height: 0 }}
            animate={{ opacity: 1, height: 'auto' }}
            exit={{ opacity: 0, height: 0 }}
            className="mt-4 p-4 bg-base-200 rounded-lg"
          >
            <div className="flex justify-between items-center mb-4">
              <h3 className="font-bold">Фильтры</h3>
              {activeFiltersCount > 0 && (
                <button 
                  onClick={clearFilters}
                  className="btn btn-ghost btn-xs"
                >
                  Очистить все
                </button>
              )}
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              {/* Категория */}
              <div className="form-control">
                <label className="label">
                  <span className="label-text">Категория</span>
                </label>
                <select 
                  className="select select-bordered select-sm w-full"
                  value={filters.category}
                  onChange={(e) => setFilters({ ...filters, category: e.target.value })}
                >
                  <option value="">Все категории</option>
                  {categories.map(cat => (
                    <option key={cat.id} value={cat.id}>
                      {cat.icon} {cat.name}
                    </option>
                  ))}
                </select>
              </div>

              {/* Ценовой диапазон */}
              <div className="form-control">
                <label className="label">
                  <span className="label-text">Цена (RSD)</span>
                </label>
                <div className="flex gap-2">
                  <input
                    type="number"
                    placeholder="От"
                    className="input input-bordered input-sm w-full"
                    value={filters.priceMin}
                    onChange={(e) => setFilters({ ...filters, priceMin: e.target.value })}
                  />
                  <input
                    type="number"
                    placeholder="До"
                    className="input input-bordered input-sm w-full"
                    value={filters.priceMax}
                    onChange={(e) => setFilters({ ...filters, priceMax: e.target.value })}
                  />
                </div>
              </div>

              {/* Радиус поиска */}
              <div className="form-control">
                <label className="label">
                  <span className="label-text">Радиус поиска</span>
                  <span className="label-text-alt">{filters.radius} км</span>
                </label>
                <input
                  type="range"
                  min="1"
                  max="50"
                  value={filters.radius}
                  onChange={(e) => setFilters({ ...filters, radius: parseInt(e.target.value) })}
                  className="range range-primary range-sm"
                />
                <div className="w-full flex justify-between text-xs px-2">
                  <span>1</span>
                  <span>10</span>
                  <span>25</span>
                  <span>50</span>
                </div>
              </div>

              {/* Состояние */}
              <div className="form-control">
                <label className="label">
                  <span className="label-text">Состояние</span>
                </label>
                <select 
                  className="select select-bordered select-sm w-full"
                  value={filters.condition}
                  onChange={(e) => setFilters({ ...filters, condition: e.target.value })}
                >
                  <option value="">Любое</option>
                  <option value="new">Новое</option>
                  <option value="like_new">Как новое</option>
                  <option value="good">Хорошее</option>
                  <option value="fair">Удовлетворительное</option>
                </select>
              </div>

              {/* Безопасная сделка */}
              <div className="form-control">
                <label className="label cursor-pointer">
                  <span className="label-text">Только с эскроу</span>
                  <input
                    type="checkbox"
                    className="toggle toggle-primary toggle-sm"
                    checked={filters.escrowOnly}
                    onChange={(e) => setFilters({ ...filters, escrowOnly: e.target.checked })}
                  />
                </label>
              </div>
            </div>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
};
```

Эти примеры кода демонстрируют современные подходы к созданию UI/UX для маркетплейса с учетом всех рассмотренных трендов и лучших практик.