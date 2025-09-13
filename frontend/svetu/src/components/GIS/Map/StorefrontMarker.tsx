import React from 'react';

interface StorefrontMarkerProps {
  title: string;
  productCount?: number;
  avgPrice?: number;
  onClick?: () => void;
  onMouseEnter?: () => void;
  onMouseLeave?: () => void;
}

const StorefrontMarker: React.FC<StorefrontMarkerProps> = ({
  title,
  productCount,
  avgPrice,
  onClick,
  onMouseEnter,
  onMouseLeave
}) => {
  return (
    <div
      className="relative cursor-pointer hover:scale-110 transition-all duration-300"
      onClick={onClick}
      onMouseEnter={onMouseEnter}
      onMouseLeave={onMouseLeave}
    >
      {/* Основная иконка витрины */}
      <div
        className="flex items-center justify-center w-12 h-12 rounded-full shadow-lg border-2 border-white"
        style={{
          background: 'linear-gradient(135deg, #FF6B6B 0%, #FF8E53 100%)',
        }}
      >
        <span className="text-2xl">🏪</span>
      </div>

      {/* Количество товаров */}
      {productCount && productCount > 0 && (
        <div
          className="absolute -top-1 -right-1 bg-primary text-primary-content text-xs rounded-full w-5 h-5 flex items-center justify-center font-bold shadow-md"
        >
          {productCount}
        </div>
      )}

      {/* Подсказка при наведении */}
      <div className="absolute bottom-full left-1/2 transform -translate-x-1/2 mb-2 opacity-0 hover:opacity-100 transition-opacity duration-200 pointer-events-none z-10">
        <div className="bg-base-100 text-base-content px-2 py-1 rounded shadow-xl whitespace-nowrap text-xs">
          <div className="font-semibold">{title}</div>
          {productCount && <div>{productCount} товаров</div>}
          {avgPrice && <div>Средняя цена: {avgPrice.toFixed(0)} ₽</div>}
        </div>
      </div>
    </div>
  );
};

export default StorefrontMarker;