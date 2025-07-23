'use client';

import React, { useEffect, useState } from 'react';
import {
  X,
  Heart,
  Share2,
  ShoppingCart,
  Eye,
  Star,
  Shield,
  ChevronLeft,
  ChevronRight,
} from 'lucide-react';
import { DistanceBadge } from './DistanceBadge';

interface QuickViewProps {
  isOpen: boolean;
  onClose: () => void;
  product: {
    id: string;
    title: string;
    price: string;
    description: string;
    images: string[];
    category: string;
    seller: {
      name: string;
      rating: number;
      totalReviews: number;
      avatar?: string;
    };
    location?: {
      address: string;
      distance: number;
    };
    stats?: {
      views: number;
      favorites: number;
    };
    condition?: 'new' | 'used' | 'refurbished';
  };
}

export const QuickView: React.FC<QuickViewProps> = ({
  isOpen,
  onClose,
  product,
}) => {
  const [currentImageIndex, setCurrentImageIndex] = useState(0);
  const [isImageLoading, setIsImageLoading] = useState(false);

  useEffect(() => {
    if (isOpen) {
      document.body.style.overflow = 'hidden';
    } else {
      document.body.style.overflow = 'unset';
    }

    return () => {
      document.body.style.overflow = 'unset';
    };
  }, [isOpen]);

  const handlePrevImage = () => {
    setIsImageLoading(true);
    setCurrentImageIndex((prev) =>
      prev === 0 ? product.images.length - 1 : prev - 1
    );
  };

  const handleNextImage = () => {
    setIsImageLoading(true);
    setCurrentImageIndex((prev) =>
      prev === product.images.length - 1 ? 0 : prev + 1
    );
  };

  if (!isOpen) return null;

  return (
    <>
      {/* Backdrop */}
      <div
        className="fixed inset-0 bg-black/50 backdrop-blur-sm z-50 animate-fadeIn"
        onClick={onClose}
      />

      {/* Modal */}
      <div className="fixed inset-4 md:inset-8 lg:inset-12 z-50 flex items-center justify-center pointer-events-none">
        <div
          className="bg-base-100 rounded-2xl shadow-2xl w-full h-full max-w-6xl max-h-[90vh] overflow-hidden pointer-events-auto animate-slideUp"
          onClick={(e) => e.stopPropagation()}
        >
          {/* Header */}
          <div className="flex items-center justify-between p-4 border-b border-base-200">
            <h2 className="text-lg font-semibold">Быстрый просмотр</h2>
            <button
              onClick={onClose}
              className="btn btn-ghost btn-sm btn-circle"
            >
              <X className="w-5 h-5" />
            </button>
          </div>

          {/* Content */}
          <div className="flex flex-col lg:flex-row h-[calc(100%-4rem)]">
            {/* Images Section */}
            <div className="lg:w-1/2 h-1/2 lg:h-full bg-base-200 relative group">
              {/* Main Image */}
              <div className="relative h-full flex items-center justify-center p-4">
                <img
                  src={
                    product.images[currentImageIndex] ||
                    'https://via.placeholder.com/600x400'
                  }
                  alt={product.title}
                  className={`max-w-full max-h-full object-contain rounded-lg transition-opacity duration-300 ${
                    isImageLoading ? 'opacity-0' : 'opacity-100'
                  }`}
                  onLoad={() => setIsImageLoading(false)}
                />

                {/* Navigation Arrows */}
                {product.images.length > 1 && (
                  <>
                    <button
                      onClick={handlePrevImage}
                      className="absolute left-2 top-1/2 -translate-y-1/2 btn btn-circle btn-sm bg-base-100/80 hover:bg-base-100 opacity-0 group-hover:opacity-100 transition-opacity"
                    >
                      <ChevronLeft className="w-4 h-4" />
                    </button>
                    <button
                      onClick={handleNextImage}
                      className="absolute right-2 top-1/2 -translate-y-1/2 btn btn-circle btn-sm bg-base-100/80 hover:bg-base-100 opacity-0 group-hover:opacity-100 transition-opacity"
                    >
                      <ChevronRight className="w-4 h-4" />
                    </button>
                  </>
                )}

                {/* Image Counter */}
                {product.images.length > 1 && (
                  <div className="absolute bottom-4 left-1/2 -translate-x-1/2 bg-base-100/80 backdrop-blur-sm rounded-full px-3 py-1">
                    <span className="text-sm font-medium">
                      {currentImageIndex + 1} / {product.images.length}
                    </span>
                  </div>
                )}
              </div>

              {/* Thumbnails */}
              {product.images.length > 1 && (
                <div className="absolute bottom-0 left-0 right-0 p-4 flex gap-2 justify-center overflow-x-auto">
                  {product.images.map((image, index) => (
                    <button
                      key={index}
                      onClick={() => {
                        setIsImageLoading(true);
                        setCurrentImageIndex(index);
                      }}
                      className={`w-16 h-16 rounded-lg overflow-hidden border-2 transition-all ${
                        index === currentImageIndex
                          ? 'border-primary ring-2 ring-primary/30'
                          : 'border-transparent hover:border-base-300'
                      }`}
                    >
                      <img
                        src={image}
                        alt={`${product.title} ${index + 1}`}
                        className="w-full h-full object-cover"
                      />
                    </button>
                  ))}
                </div>
              )}
            </div>

            {/* Details Section */}
            <div className="lg:w-1/2 h-1/2 lg:h-full overflow-y-auto">
              <div className="p-6 space-y-4">
                {/* Category & Condition */}
                <div className="flex items-center gap-2">
                  <span className="badge badge-ghost">{product.category}</span>
                  {product.condition && product.condition !== 'new' && (
                    <span className="badge badge-success">
                      {product.condition === 'used' ? 'Б/У' : 'Восстановлено'}
                    </span>
                  )}
                  {product.location && (
                    <DistanceBadge
                      distance={product.location.distance}
                      variant="compact"
                    />
                  )}
                </div>

                {/* Title & Price */}
                <div>
                  <h1 className="text-2xl font-bold mb-2">{product.title}</h1>
                  <p className="text-3xl font-bold text-primary">
                    {product.price}
                  </p>
                </div>

                {/* Stats */}
                {product.stats && (
                  <div className="flex items-center gap-4 text-sm text-base-content/60">
                    <div className="flex items-center gap-1">
                      <Eye className="w-4 h-4" />
                      <span>{product.stats.views} просмотров</span>
                    </div>
                    <div className="flex items-center gap-1">
                      <Heart className="w-4 h-4" />
                      <span>{product.stats.favorites} в избранном</span>
                    </div>
                  </div>
                )}

                {/* Description */}
                <div className="py-4">
                  <h3 className="font-semibold mb-2">Описание</h3>
                  <p className="text-base-content/80 whitespace-pre-wrap">
                    {product.description}
                  </p>
                </div>

                {/* Seller */}
                <div className="card bg-base-200 p-4">
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-3">
                      <div className="avatar">
                        <div className="w-12 h-12 rounded-full bg-base-300">
                          {product.seller.avatar ? (
                            <img
                              src={product.seller.avatar}
                              alt={product.seller.name}
                            />
                          ) : (
                            <span className="text-lg font-semibold flex items-center justify-center h-full">
                              {product.seller.name[0]}
                            </span>
                          )}
                        </div>
                      </div>
                      <div>
                        <p className="font-medium">{product.seller.name}</p>
                        <div className="flex items-center gap-1">
                          <Star className="w-4 h-4 text-warning fill-warning" />
                          <span className="text-sm">
                            {product.seller.rating.toFixed(1)} (
                            {product.seller.totalReviews})
                          </span>
                        </div>
                      </div>
                    </div>
                    <button className="btn btn-sm btn-ghost">Профиль</button>
                  </div>
                </div>

                {/* Location */}
                {product.location && (
                  <div className="py-4">
                    <h3 className="font-semibold mb-2">Местоположение</h3>
                    <p className="text-base-content/80">
                      {product.location.address}
                    </p>
                  </div>
                )}

                {/* Actions */}
                <div className="flex gap-3 pt-4">
                  <button className="btn btn-primary flex-1">
                    <ShoppingCart className="w-5 h-5" />
                    Связаться
                  </button>
                  <button className="btn btn-ghost btn-square">
                    <Heart className="w-5 h-5" />
                  </button>
                  <button className="btn btn-ghost btn-square">
                    <Share2 className="w-5 h-5" />
                  </button>
                </div>

                {/* Safe Deal Badge */}
                <div className="flex items-center gap-2 p-3 bg-success/10 text-success rounded-lg">
                  <Shield className="w-5 h-5" />
                  <span className="text-sm font-medium">Безопасная сделка</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </>
  );
};
