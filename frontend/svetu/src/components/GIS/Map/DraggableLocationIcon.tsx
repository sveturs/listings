'use client';

import React, { useState, useRef, useEffect } from 'react';
// import { MapMouseEvent } from 'react-map-gl';

interface DraggableLocationIconProps {
  onDropLocation: (lng: number, lat: number) => void;
  mapRef: React.RefObject<any>;
}

export default function DraggableLocationIcon({
  onDropLocation,
  mapRef,
}: DraggableLocationIconProps) {
  const [isDragging, setIsDragging] = useState(false);
  const [dragPosition, setDragPosition] = useState({ x: 0, y: 0 });
  const iconRef = useRef<HTMLDivElement>(null);
  const draggedIconRef = useRef<HTMLDivElement>(null);

  // Обработка начала перетаскивания
  const handleDragStart = (e: React.DragEvent<HTMLDivElement>) => {
    e.dataTransfer.effectAllowed = 'move';
    setIsDragging(true);

    // Создаем пустое изображение для drag preview
    const dragImage = new Image();
    dragImage.src =
      'data:image/gif;base64,R0lGODlhAQABAIAAAAUEBAAAACwAAAAAAQABAAACAkQBADs=';
    e.dataTransfer.setDragImage(dragImage, 0, 0);
  };

  // Обработка перетаскивания над картой
  const handleDragOver = (e: DragEvent) => {
    e.preventDefault();
    e.dataTransfer!.dropEffect = 'move';

    // Обновляем позицию перетаскиваемой иконки
    setDragPosition({
      x: e.clientX,
      y: e.clientY,
    });
  };

  // Обработка окончания перетаскивания
  const handleDrop = (e: DragEvent) => {
    e.preventDefault();
    setIsDragging(false);

    if (!mapRef.current) return;

    const map = mapRef.current;
    const rect = map.getContainer().getBoundingClientRect();

    // Вычисляем позицию относительно карты
    const x = e.clientX - rect.left;
    const y = e.clientY - rect.top;

    // Конвертируем пиксели в координаты
    const lngLat = map.unproject([x, y]);

    // Вызываем callback с новыми координатами
    onDropLocation(lngLat.lng, lngLat.lat);
  };

  // Обработка окончания перетаскивания
  const handleDragEnd = () => {
    setIsDragging(false);
  };

  // Добавляем обработчики событий на карту
  useEffect(() => {
    if (!mapRef.current) return;

    const mapContainer = mapRef.current.getContainer();

    const handleMapDragOver = (e: DragEvent) => {
      if (isDragging) {
        handleDragOver(e);
      }
    };

    const handleMapDrop = (e: DragEvent) => {
      if (isDragging) {
        handleDrop(e);
      }
    };

    mapContainer.addEventListener('dragover', handleMapDragOver);
    mapContainer.addEventListener('drop', handleMapDrop);

    return () => {
      mapContainer.removeEventListener('dragover', handleMapDragOver);
      mapContainer.removeEventListener('drop', handleMapDrop);
    };
  }, [isDragging, mapRef]);

  // Обработка для мобильных устройств
  const handleTouchStart = (e: React.TouchEvent) => {
    e.preventDefault(); // Предотвращаем прокрутку карты
    const touch = e.touches[0];
    setIsDragging(true);
    setDragPosition({
      x: touch.clientX,
      y: touch.clientY,
    });
  };

  const handleTouchMove = (e: React.TouchEvent) => {
    if (!isDragging) return;

    e.preventDefault(); // Предотвращаем прокрутку карты
    const touch = e.touches[0];
    setDragPosition({
      x: touch.clientX,
      y: touch.clientY,
    });
  };

  const handleTouchEnd = (e: React.TouchEvent) => {
    if (!isDragging || !mapRef.current) return;

    e.preventDefault(); // Предотвращаем прокрутку карты
    setIsDragging(false);

    const map = mapRef.current;
    const rect = map.getContainer().getBoundingClientRect();

    // Используем последнюю позицию из dragPosition
    const x = dragPosition.x - rect.left;
    const y = dragPosition.y - rect.top;

    // Проверяем, что координаты внутри карты
    if (x >= 0 && x <= rect.width && y >= 0 && y <= rect.height) {
      const lngLat = map.unproject([x, y]);
      onDropLocation(lngLat.lng, lngLat.lat);
    }
  };

  return (
    <>
      {/* Иконка для перетаскивания */}
      <div
        ref={iconRef}
        className={`absolute right-4 z-10 cursor-move transition-all duration-200 group ${
          isDragging ? 'opacity-50 scale-95' : 'opacity-100 hover:scale-105'
        } top-[calc(50%-60px)] md:top-[calc(50%+60px)] -translate-y-1/2`}
        draggable
        onDragStart={handleDragStart}
        onDragEnd={handleDragEnd}
        onTouchStart={handleTouchStart}
        onTouchMove={handleTouchMove}
        onTouchEnd={handleTouchEnd}
        style={{
          touchAction: 'none', // Отключаем touch-action для лучшего контроля
        }}
      >
        <div className="relative">
          <div className="bg-white rounded-lg shadow-lg p-3 hover:shadow-xl transition-shadow border-2 border-transparent hover:border-red-100">
            <svg
              width="24"
              height="24"
              viewBox="0 0 24 24"
              fill="none"
              xmlns="http://www.w3.org/2000/svg"
              className="text-red-500"
            >
              <circle cx="12" cy="12" r="3" fill="currentColor" />
              <path
                d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z"
                fill="currentColor"
              />
            </svg>
          </div>
          {/* Подсказка */}
          <div className="absolute top-full mt-2 right-0 bg-gray-800 text-white text-xs rounded-lg px-3 py-2 whitespace-nowrap pointer-events-none opacity-0 group-hover:opacity-100 transition-opacity">
            Перетащите на карту
          </div>
        </div>
      </div>

      {/* Перетаскиваемая иконка, следующая за курсором */}
      {isDragging && (
        <div
          ref={draggedIconRef}
          className="fixed pointer-events-none z-50"
          style={{
            left: dragPosition.x - 12,
            top: dragPosition.y - 24,
          }}
        >
          <svg
            width="24"
            height="24"
            viewBox="0 0 24 24"
            fill="none"
            xmlns="http://www.w3.org/2000/svg"
            className="text-red-500"
          >
            <circle cx="12" cy="12" r="3" fill="currentColor" />
            <path
              d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z"
              fill="currentColor"
            />
          </svg>
        </div>
      )}
    </>
  );
}
