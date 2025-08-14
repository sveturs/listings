'use client';

import React, { useEffect, useState } from 'react';
import L from 'leaflet';
import 'leaflet/dist/leaflet.css';
import { FaMapMarkerAlt, FaEye } from 'react-icons/fa';

// Fix for default markers in Leaflet
delete (L.Icon.Default.prototype as any)._getIconUrl;
L.Icon.Default.mergeOptions({
  iconRetinaUrl:
    'https://unpkg.com/leaflet@1.9.4/dist/images/marker-icon-2x.png',
  iconUrl: 'https://unpkg.com/leaflet@1.9.4/dist/images/marker-icon.png',
  shadowUrl: 'https://unpkg.com/leaflet@1.9.4/dist/images/marker-shadow.png',
});

interface MapItem {
  id: number;
  title: string;
  price: number;
  category: string;
  lat: number;
  lng: number;
  image: string;
  distance: string;
}

export default function MapSection() {
  const [map, setMap] = useState<L.Map | null>(null);
  const [userLocation, setUserLocation] = useState<[number, number]>([
    44.8125, 20.4612,
  ]); // Belgrade default

  // Примеры товаров рядом
  const nearbyItems: MapItem[] = [
    {
      id: 1,
      title: 'iPhone 14 Pro',
      price: 999,
      category: 'Электроника',
      lat: 44.815,
      lng: 20.465,
      image:
        'https://images.unsplash.com/photo-1678685888221-cda773a3dcdb?w=400',
      distance: '0.5 км',
    },
    {
      id: 2,
      title: 'Диван IKEA',
      price: 450,
      category: 'Мебель',
      lat: 44.81,
      lng: 20.458,
      image: 'https://images.unsplash.com/photo-1555041469-a586c61ea9bc?w=400',
      distance: '0.8 км',
    },
    {
      id: 3,
      title: 'PlayStation 5',
      price: 550,
      category: 'Игры',
      lat: 44.818,
      lng: 20.47,
      image:
        'https://images.unsplash.com/photo-1606813907291-d86efa9b94db?w=400',
      distance: '1.2 км',
    },
    {
      id: 4,
      title: 'MacBook Air M2',
      price: 1299,
      category: 'Компьютеры',
      lat: 44.808,
      lng: 20.455,
      image:
        'https://images.unsplash.com/photo-1517336714731-489689fd1ca4?w=400',
      distance: '1.5 км',
    },
    {
      id: 5,
      title: 'Электросамокат Xiaomi',
      price: 399,
      category: 'Транспорт',
      lat: 44.82,
      lng: 20.463,
      image: 'https://images.unsplash.com/photo-1558618666-fcd25c85cd64?w=400',
      distance: '2.0 км',
    },
  ];

  useEffect(() => {
    // Получаем локацию пользователя
    if (navigator.geolocation) {
      navigator.geolocation.getCurrentPosition(
        (position) => {
          setUserLocation([
            position.coords.latitude,
            position.coords.longitude,
          ]);
        },
        (error) => {
          console.log('Geolocation error:', error);
        }
      );
    }
  }, []);

  useEffect(() => {
    // Инициализация карты
    if (!map && typeof window !== 'undefined') {
      // Проверяем, что контейнер существует и не инициализирован
      const container = document.getElementById('map-container');
      if (!container || container.hasChildNodes()) {
        return;
      }

      const mapInstance = L.map('map-container').setView(userLocation, 13);

      // Добавляем слой OpenStreetMap
      L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
        attribution:
          '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors',
      }).addTo(mapInstance);

      // Добавляем маркер пользователя
      const userIcon = L.divIcon({
        html: `
          <div class="relative">
            <div class="absolute -inset-4 bg-blue-500 rounded-full animate-ping opacity-30"></div>
            <div class="relative w-4 h-4 bg-blue-500 rounded-full border-2 border-white shadow-lg"></div>
          </div>
        `,
        className: 'user-marker',
        iconSize: [20, 20],
        iconAnchor: [10, 10],
      });

      L.marker(userLocation, { icon: userIcon })
        .addTo(mapInstance)
        .bindPopup('<b>Вы здесь</b>');

      // Добавляем маркеры товаров
      nearbyItems.forEach((item) => {
        const itemIcon = L.divIcon({
          html: `
            <div class="bg-white rounded-lg shadow-lg p-1 border-2 border-primary">
              <div class="text-xs font-bold text-center">€${item.price}</div>
            </div>
          `,
          className: 'item-marker',
          iconSize: [60, 30],
          iconAnchor: [30, 30],
        });

        const popupContent = `
          <div class="w-48">
            <img src="${item.image}" alt="${item.title}" class="w-full h-24 object-cover rounded-lg mb-2">
            <h3 class="font-bold text-sm">${item.title}</h3>
            <p class="text-xs text-gray-600 mb-1">${item.category}</p>
            <div class="flex justify-between items-center">
              <span class="text-lg font-bold text-primary">€${item.price}</span>
              <span class="text-xs text-gray-500">${item.distance}</span>
            </div>
            <button class="btn btn-primary btn-xs btn-block mt-2">Посмотреть</button>
          </div>
        `;

        L.marker([item.lat, item.lng], { icon: itemIcon })
          .addTo(mapInstance)
          .bindPopup(popupContent);
      });

      // Добавляем круг радиуса поиска
      L.circle(userLocation, {
        color: 'rgba(59, 130, 246, 0.5)',
        fillColor: 'rgba(59, 130, 246, 0.1)',
        fillOpacity: 0.2,
        radius: 2000, // 2км радиус
      }).addTo(mapInstance);

      setMap(mapInstance);
    }

    // Cleanup
    return () => {
      if (map) {
        map.remove();
        setMap(null);
      }
    };
  }, [userLocation]); // eslint-disable-line react-hooks/exhaustive-deps

  return (
    <div className="relative w-full h-full">
      <div id="map-container" className="w-full h-full rounded-lg" />

      {/* Легенда карты */}
      <div className="absolute top-4 left-4 bg-white rounded-lg shadow-lg p-3 z-[1000]">
        <div className="flex items-center gap-2 mb-2">
          <div className="w-3 h-3 bg-blue-500 rounded-full"></div>
          <span className="text-xs">Ваше местоположение</span>
        </div>
        <div className="flex items-center gap-2">
          <div className="w-8 h-4 bg-primary rounded text-white text-[8px] flex items-center justify-center">
            €
          </div>
          <span className="text-xs">Товары поблизости</span>
        </div>
      </div>

      {/* Кнопки управления */}
      <div className="absolute bottom-4 right-4 flex flex-col gap-2 z-[1000]">
        <button
          className="btn btn-circle btn-sm bg-white shadow-lg hover:shadow-xl"
          onClick={() => {
            if (map) {
              map.setView(userLocation, 13);
            }
          }}
        >
          <FaMapMarkerAlt />
        </button>
        <button className="btn btn-circle btn-sm bg-white shadow-lg hover:shadow-xl">
          <FaEye />
        </button>
      </div>

      {/* Быстрые категории на карте */}
      <div className="absolute top-4 right-4 flex flex-col gap-2 z-[1000]">
        {['Все', 'Электроника', 'Мебель', 'Авто'].map((cat) => (
          <button
            key={cat}
            className="btn btn-xs btn-outline bg-white/90 backdrop-blur-sm"
          >
            {cat}
          </button>
        ))}
      </div>
    </div>
  );
}
