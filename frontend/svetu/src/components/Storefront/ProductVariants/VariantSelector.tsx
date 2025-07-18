'use client';

import React, { useState, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import SafeImage from '../../SafeImage';

interface ProductVariant {
  id: number;
  sku?: string;
  price?: number;
  stock_quantity: number;
  stock_status: string;
  variant_attributes: Record<string, string>;
  is_default: boolean;
  images: ProductVariantImage[];
}

interface ProductVariantImage {
  id: number;
  image_url: string;
  thumbnail_url?: string;
  alt_text?: string;
  is_main: boolean;
}

interface VariantSelectorProps {
  productId: number;
  basePrice: number;
  baseCurrency: string;
  onVariantChange: (variant: ProductVariant | null) => void;
  className?: string;
}

export default function VariantSelector({
  productId,
  basePrice,
  baseCurrency,
  onVariantChange,
  className = '',
}: VariantSelectorProps) {
  const t = useTranslations('storefront');
  const [variants, setVariants] = useState<ProductVariant[]>([]);
  const [selectedVariant, setSelectedVariant] = useState<ProductVariant | null>(
    null
  );
  const [selectedAttributes, setSelectedAttributes] = useState<
    Record<string, string>
  >({});
  const [availableAttributes, setAvailableAttributes] = useState<
    Record<string, string[]>
  >({});
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadVariants();
  }, [productId]);

  useEffect(() => {
    // Find matching variant when attributes change
    const matchingVariant = findMatchingVariant(selectedAttributes);
    setSelectedVariant(matchingVariant);
    onVariantChange(matchingVariant);
  }, [selectedAttributes, variants]);

  const loadVariants = async () => {
    try {
      const response = await fetch(
        `/api/v1/storefront/products/${productId}/variants`
      );
      if (response.ok) {
        const variantData = await response.json();
        setVariants(variantData);

        // Extract available attributes
        const attributes: Record<string, Set<string>> = {};
        variantData.forEach((variant: ProductVariant) => {
          Object.entries(variant.variant_attributes).forEach(([key, value]) => {
            if (!attributes[key]) {
              attributes[key] = new Set();
            }
            attributes[key].add(value);
          });
        });

        // Convert Sets to arrays
        const availableAttrs: Record<string, string[]> = {};
        Object.entries(attributes).forEach(([key, valueSet]) => {
          availableAttrs[key] = Array.from(valueSet);
        });
        setAvailableAttributes(availableAttrs);

        // Set default variant
        const defaultVariant =
          variantData.find((v: ProductVariant) => v.is_default) ||
          variantData[0];
        if (defaultVariant) {
          setSelectedAttributes(defaultVariant.variant_attributes);
        }
      }
    } catch (error) {
      console.error('Failed to load variants:', error);
    } finally {
      setLoading(false);
    }
  };

  const findMatchingVariant = (
    attributes: Record<string, string>
  ): ProductVariant | null => {
    return (
      variants.find((variant) => {
        return Object.entries(attributes).every(
          ([key, value]) => variant.variant_attributes[key] === value
        );
      }) || null
    );
  };

  const updateAttribute = (attributeName: string, value: string) => {
    setSelectedAttributes((prev) => ({
      ...prev,
      [attributeName]: value,
    }));
  };

  const getAttributeDisplayName = (attributeName: string): string => {
    // Convert snake_case to Title Case
    return attributeName
      .split('_')
      .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
      .join(' ');
  };

  const isAttributeValueAvailable = (
    attributeName: string,
    value: string
  ): boolean => {
    // Check if selecting this value would result in a valid variant
    const testAttributes = { ...selectedAttributes, [attributeName]: value };
    return findMatchingVariant(testAttributes) !== null;
  };

  const getCurrentPrice = (): number => {
    return selectedVariant?.price || basePrice;
  };

  const getCurrentStock = (): number => {
    return selectedVariant?.stock_quantity || 0;
  };

  const getStockStatus = (): string => {
    return selectedVariant?.stock_status || 'out_of_stock';
  };

  const getMainImage = (): string | null => {
    if (selectedVariant?.images && selectedVariant.images.length > 0) {
      const mainImage =
        selectedVariant.images.find((img) => img.is_main) ||
        selectedVariant.images[0];
      return mainImage.image_url;
    }
    return null;
  };

  if (loading) {
    return (
      <div className={`animate-pulse space-y-4 ${className}`}>
        <div className="h-4 bg-gray-200 rounded w-1/4"></div>
        <div className="h-10 bg-gray-200 rounded"></div>
        <div className="h-4 bg-gray-200 rounded w-1/3"></div>
      </div>
    );
  }

  if (variants.length === 0) {
    return null; // No variants available
  }

  return (
    <div className={`space-y-6 ${className}`}>
      {/* Variant Image */}
      {getMainImage() && (
        <div className="relative w-full h-64 bg-gray-100 rounded-lg overflow-hidden">
          <SafeImage
            src={getMainImage()!}
            alt={selectedVariant?.sku || 'Product variant'}
            fill
            className="object-cover"
            sizes="(max-width: 768px) 100vw, 50vw"
          />
        </div>
      )}

      {/* Price Display */}
      <div className="space-y-2">
        <div className="text-2xl font-bold text-gray-900">
          {getCurrentPrice().toLocaleString()} {baseCurrency}
        </div>
        {selectedVariant?.price && selectedVariant.price !== basePrice && (
          <div className="text-sm text-gray-500 line-through">
            {basePrice.toLocaleString()} {baseCurrency}
          </div>
        )}
      </div>

      {/* Attribute Selectors */}
      <div className="space-y-4">
        {Object.entries(availableAttributes).map(([attributeName, values]) => (
          <div key={attributeName} className="space-y-2">
            <label className="block text-sm font-medium text-gray-700">
              {getAttributeDisplayName(attributeName)}
            </label>

            {attributeName === 'color' ? (
              // Color picker for color attributes
              <div className="flex flex-wrap gap-2">
                {values.map((value) => {
                  const isSelected =
                    selectedAttributes[attributeName] === value;
                  const isAvailable = isAttributeValueAvailable(
                    attributeName,
                    value
                  );

                  return (
                    <button
                      key={value}
                      onClick={() => updateAttribute(attributeName, value)}
                      disabled={!isAvailable}
                      className={`
                        w-8 h-8 rounded-full border-2 relative
                        ${isSelected ? 'border-blue-600 ring-2 ring-blue-200' : 'border-gray-300'}
                        ${!isAvailable ? 'opacity-50 cursor-not-allowed' : 'hover:border-gray-400'}
                      `}
                      style={{ backgroundColor: getColorHex(value) }}
                      title={value}
                    >
                      {!isAvailable && (
                        <div className="absolute inset-0 flex items-center justify-center">
                          <div className="w-6 h-0.5 bg-red-500 rotate-45"></div>
                        </div>
                      )}
                    </button>
                  );
                })}
              </div>
            ) : (
              // Dropdown for other attributes
              <select
                value={selectedAttributes[attributeName] || ''}
                onChange={(e) => updateAttribute(attributeName, e.target.value)}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                <option value="">{t('select_option')}</option>
                {values.map((value) => (
                  <option
                    key={value}
                    value={value}
                    disabled={!isAttributeValueAvailable(attributeName, value)}
                  >
                    {value}
                  </option>
                ))}
              </select>
            )}
          </div>
        ))}
      </div>

      {/* Stock Status */}
      <div className="space-y-2">
        <div className="flex items-center justify-between">
          <span className="text-sm text-gray-600">{t('availability')}</span>
          <span
            className={`text-sm font-medium ${
              getStockStatus() === 'in_stock'
                ? 'text-green-600'
                : getStockStatus() === 'low_stock'
                  ? 'text-yellow-600'
                  : 'text-red-600'
            }`}
          >
            {getStockStatus() === 'in_stock' && t('in_stock')}
            {getStockStatus() === 'low_stock' && t('low_stock')}
            {getStockStatus() === 'out_of_stock' && t('out_of_stock')}
          </span>
        </div>

        {getCurrentStock() > 0 && (
          <div className="text-sm text-gray-600">
            {t('items_available', { count: getCurrentStock() })}
          </div>
        )}

        {selectedVariant?.sku && (
          <div className="text-xs text-gray-500">
            SKU: {selectedVariant.sku}
          </div>
        )}
      </div>

      {/* Variant Images Gallery */}
      {selectedVariant?.images && selectedVariant.images.length > 1 && (
        <div className="space-y-2">
          <h4 className="text-sm font-medium text-gray-700">
            {t('more_images')}
          </h4>
          <div className="grid grid-cols-4 gap-2">
            {selectedVariant.images.map((image) => (
              <div
                key={image.id}
                className="relative aspect-square bg-gray-100 rounded overflow-hidden"
              >
                <SafeImage
                  src={image.thumbnail_url || image.image_url}
                  alt={image.alt_text || 'Product image'}
                  fill
                  className="object-cover"
                  sizes="100px"
                />
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}

// Helper function to get color hex value (you might want to store this in your data)
function getColorHex(colorName: string): string {
  const colorMap: Record<string, string> = {
    black: '#000000',
    white: '#FFFFFF',
    red: '#FF0000',
    blue: '#0000FF',
    green: '#00FF00',
    yellow: '#FFFF00',
    purple: '#800080',
    orange: '#FFA500',
    pink: '#FFC0CB',
    gray: '#808080',
    brown: '#A52A2A',
    gold: '#FFD700',
    silver: '#C0C0C0',
  };

  return colorMap[colorName.toLowerCase()] || '#CCCCCC';
}
