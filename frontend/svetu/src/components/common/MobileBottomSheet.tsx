'use client';

import React, { useEffect, useRef, useState } from 'react';
import { X, GripHorizontal } from 'lucide-react';

interface MobileBottomSheetProps {
  isOpen: boolean;
  onClose: () => void;
  title?: string;
  children: React.ReactNode;
  height?: 'auto' | 'full' | '3/4' | '1/2';
  showHandle?: boolean;
  closeOnOverlayClick?: boolean;
}

/**
 * Bottom Sheet компонент для мобильных устройств
 * Поддерживает свайп для закрытия и различные высоты
 */
export const MobileBottomSheet: React.FC<MobileBottomSheetProps> = ({
  isOpen,
  onClose,
  title,
  children,
  height = '3/4',
  showHandle = true,
  closeOnOverlayClick = true,
}) => {
  const sheetRef = useRef<HTMLDivElement>(null);
  const [isDragging, setIsDragging] = useState(false);
  const [startY, setStartY] = useState(0);
  const [currentY, setCurrentY] = useState(0);
  const [translateY, setTranslateY] = useState(100);

  // Определяем высоту в зависимости от prop
  const getHeightClass = () => {
    switch (height) {
      case 'full':
        return 'h-full';
      case '3/4':
        return 'h-3/4';
      case '1/2':
        return 'h-1/2';
      case 'auto':
      default:
        return 'max-h-[90vh]';
    }
  };

  // Открытие/закрытие анимации
  useEffect(() => {
    if (isOpen) {
      setTranslateY(0);
      document.body.style.overflow = 'hidden';
    } else {
      setTranslateY(100);
      document.body.style.overflow = '';
    }

    return () => {
      document.body.style.overflow = '';
    };
  }, [isOpen]);

  // Обработка начала свайпа
  const handleTouchStart = (e: React.TouchEvent) => {
    setIsDragging(true);
    setStartY(e.touches[0].clientY);
    setCurrentY(e.touches[0].clientY);
  };

  // Обработка движения свайпа
  const handleTouchMove = (e: React.TouchEvent) => {
    if (!isDragging) return;

    const deltaY = e.touches[0].clientY - startY;

    // Разрешаем только свайп вниз
    if (deltaY > 0) {
      setCurrentY(e.touches[0].clientY);
      const percentage = (deltaY / window.innerHeight) * 100;
      setTranslateY(Math.min(percentage, 100));
    }
  };

  // Обработка окончания свайпа
  const handleTouchEnd = () => {
    if (!isDragging) return;

    const deltaY = currentY - startY;
    const threshold = window.innerHeight * 0.3; // 30% высоты экрана

    // Если свайп больше порога - закрываем
    if (deltaY > threshold) {
      onClose();
    } else {
      // Иначе возвращаем на место
      setTranslateY(0);
    }

    setIsDragging(false);
    setStartY(0);
    setCurrentY(0);
  };

  // Обработка клика по оверлею
  const handleOverlayClick = () => {
    if (closeOnOverlayClick) {
      onClose();
    }
  };

  // Предотвращение клика внутри контента
  const handleContentClick = (e: React.MouseEvent) => {
    e.stopPropagation();
  };

  if (!isOpen && translateY === 100) {
    return null;
  }

  return (
    <>
      {/* Overlay */}
      <div
        className={`fixed inset-0 bg-black/50 z-40 transition-opacity duration-300 ${
          isOpen ? 'opacity-100' : 'opacity-0 pointer-events-none'
        }`}
        onClick={handleOverlayClick}
      />

      {/* Bottom Sheet */}
      <div
        ref={sheetRef}
        className={`fixed bottom-0 left-0 right-0 z-50 bg-base-100 rounded-t-3xl shadow-2xl ${getHeightClass()} transition-transform duration-300 ease-out`}
        style={{
          transform: `translateY(${translateY}%)`,
        }}
        onClick={handleContentClick}
      >
        {/* Handle для свайпа */}
        {showHandle && (
          <div
            className="flex justify-center py-3 cursor-grab active:cursor-grabbing"
            onTouchStart={handleTouchStart}
            onTouchMove={handleTouchMove}
            onTouchEnd={handleTouchEnd}
          >
            <GripHorizontal className="w-10 h-1 bg-base-300 rounded-full" />
          </div>
        )}

        {/* Header */}
        {title && (
          <div className="flex items-center justify-between px-4 pb-3 border-b border-base-200">
            <h3 className="text-lg font-semibold">{title}</h3>
            <button
              onClick={onClose}
              className="btn btn-ghost btn-sm btn-circle"
              aria-label="Close"
            >
              <X className="w-5 h-5" />
            </button>
          </div>
        )}

        {/* Content */}
        <div className="flex-1 overflow-y-auto overscroll-contain p-4">
          {children}
        </div>
      </div>
    </>
  );
};

/**
 * Hook для управления состоянием Bottom Sheet
 */
export const useBottomSheet = (initialState = false) => {
  const [isOpen, setIsOpen] = useState(initialState);

  const open = () => setIsOpen(true);
  const close = () => setIsOpen(false);
  const toggle = () => setIsOpen((prev) => !prev);

  return {
    isOpen,
    open,
    close,
    toggle,
  };
};

export default MobileBottomSheet;
