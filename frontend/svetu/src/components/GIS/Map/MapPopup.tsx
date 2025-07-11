import React from 'react';
import Image from 'next/image';
import { Popup } from 'react-map-gl';
import { MapPopupData } from '../types/gis';

interface MapPopupProps {
  popup: MapPopupData;
  onClose?: () => void;
  className?: string;
}

const MapPopup: React.FC<MapPopupProps> = ({
  popup,
  onClose,
  className = '',
}) => {
  const handleClose = () => {
    if (popup.onClose) {
      popup.onClose();
    }
    if (onClose) {
      onClose();
    }
  };

  return (
    <Popup
      longitude={popup.position[0]}
      latitude={popup.position[1]}
      anchor="bottom"
      onClose={handleClose}
      closeButton={true}
      closeOnClick={false}
      maxWidth="320px"
      className={`map-popup ${className}`}
    >
      <div className="p-4 bg-white rounded-lg shadow-lg">
        {/* Заголовок */}
        <div className="flex items-center justify-between mb-2">
          <h3 className="font-semibold text-gray-900 text-lg leading-tight">
            {popup.title}
          </h3>
        </div>

        {/* Описание */}
        {popup.description && (
          <p className="text-gray-600 text-sm mb-3 leading-relaxed">
            {popup.description}
          </p>
        )}

        {/* Кастомный контент */}
        {popup.content && <div className="border-t pt-3">{popup.content}</div>}
      </div>
    </Popup>
  );
};

// Специализированный компонент для объявлений
interface ListingPopupProps {
  listing: {
    id: string;
    title: string;
    price: number;
    currency: string;
    rating?: number;
    reviews?: number;
    imageUrl?: string;
    distance?: number;
    category?: string;
  };
  position: [number, number];
  onClose?: () => void;
  onViewDetails?: (id: string) => void;
}

export const ListingPopup: React.FC<ListingPopupProps> = ({
  listing,
  position,
  onClose,
  onViewDetails,
}) => {
  const handleViewDetails = () => {
    if (onViewDetails) {
      onViewDetails(listing.id);
    }
  };

  return (
    <Popup
      longitude={position[0]}
      latitude={position[1]}
      anchor="bottom"
      onClose={onClose}
      closeButton={true}
      closeOnClick={false}
      maxWidth="300px"
      className="listing-popup"
    >
      <div className="bg-white rounded-lg shadow-lg overflow-hidden">
        {/* Изображение */}
        {listing.imageUrl && (
          <div className="h-32 bg-gray-200 overflow-hidden relative">
            <Image
              src={listing.imageUrl}
              alt={listing.title}
              fill
              className="object-cover"
              sizes="300px"
            />
          </div>
        )}

        {/* Контент */}
        <div className="p-4">
          {/* Заголовок и цена */}
          <div className="flex justify-between items-start mb-2">
            <h3 className="font-semibold text-gray-900 text-sm leading-tight flex-1 mr-2">
              {listing.title}
            </h3>
            <div className="text-right">
              <span className="font-bold text-primary text-lg">
                {listing.price}
                {listing.currency}
              </span>
              <div className="text-xs text-gray-500">за ночь</div>
            </div>
          </div>

          {/* Категория */}
          {listing.category && (
            <div className="mb-2">
              <span className="inline-block bg-gray-100 text-gray-700 text-xs px-2 py-1 rounded">
                {listing.category}
              </span>
            </div>
          )}

          {/* Рейтинг и отзывы */}
          {listing.rating && (
            <div className="flex items-center mb-2">
              <div className="flex items-center">
                <span className="text-yellow-400 text-sm">★</span>
                <span className="ml-1 text-sm font-medium text-gray-700">
                  {listing.rating.toFixed(1)}
                </span>
              </div>
              {listing.reviews && (
                <span className="ml-2 text-sm text-gray-500">
                  ({listing.reviews} отзывов)
                </span>
              )}
            </div>
          )}

          {/* Расстояние */}
          {listing.distance && (
            <div className="mb-3">
              <span className="text-sm text-gray-500">
                📍 {listing.distance.toFixed(1)} км от вас
              </span>
            </div>
          )}

          {/* Кнопки действий */}
          <div className="flex space-x-2">
            <button
              onClick={handleViewDetails}
              className="flex-1 bg-primary text-white text-sm py-2 px-3 rounded-md hover:bg-primary-dark transition-colors"
            >
              Подробнее
            </button>
            <button
              onClick={onClose}
              className="bg-gray-200 text-gray-700 text-sm py-2 px-3 rounded-md hover:bg-gray-300 transition-colors"
            >
              Закрыть
            </button>
          </div>
        </div>
      </div>
    </Popup>
  );
};

// Компонент для отображения информации о пользователе
interface UserPopupProps {
  user: {
    id: string;
    name: string;
    avatar?: string;
    status?: string;
    lastSeen?: Date;
  };
  position: [number, number];
  onClose?: () => void;
  onViewProfile?: (id: string) => void;
}

export const UserPopup: React.FC<UserPopupProps> = ({
  user,
  position,
  onClose,
  onViewProfile,
}) => {
  const handleViewProfile = () => {
    if (onViewProfile) {
      onViewProfile(user.id);
    }
  };

  return (
    <Popup
      longitude={position[0]}
      latitude={position[1]}
      anchor="bottom"
      onClose={onClose}
      closeButton={true}
      closeOnClick={false}
      maxWidth="250px"
      className="user-popup"
    >
      <div className="bg-white rounded-lg shadow-lg overflow-hidden">
        <div className="p-4">
          {/* Аватар и имя */}
          <div className="flex items-center mb-3">
            {user.avatar ? (
              <div className="w-12 h-12 rounded-full mr-3 relative overflow-hidden">
                <Image
                  src={user.avatar}
                  alt={user.name}
                  fill
                  className="object-cover"
                  sizes="48px"
                />
              </div>
            ) : (
              <div className="w-12 h-12 rounded-full bg-gray-200 flex items-center justify-center mr-3">
                <span className="text-gray-500 text-xl">👤</span>
              </div>
            )}
            <div>
              <h3 className="font-semibold text-gray-900 text-sm">
                {user.name}
              </h3>
              {user.status && (
                <p className="text-xs text-gray-500">{user.status}</p>
              )}
            </div>
          </div>

          {/* Последняя активность */}
          {user.lastSeen && (
            <div className="mb-3">
              <span className="text-xs text-gray-500">
                Был в сети: {user.lastSeen.toLocaleDateString()}
              </span>
            </div>
          )}

          {/* Кнопки действий */}
          <div className="flex space-x-2">
            <button
              onClick={handleViewProfile}
              className="flex-1 bg-primary text-white text-sm py-2 px-3 rounded-md hover:bg-primary-dark transition-colors"
            >
              Профиль
            </button>
            <button
              onClick={onClose}
              className="bg-gray-200 text-gray-700 text-sm py-2 px-3 rounded-md hover:bg-gray-300 transition-colors"
            >
              Закрыть
            </button>
          </div>
        </div>
      </div>
    </Popup>
  );
};

export default MapPopup;
