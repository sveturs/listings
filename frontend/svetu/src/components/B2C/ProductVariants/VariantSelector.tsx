'use client';

import React, { useState, useEffect, useCallback } from 'react';
import { useTranslations } from 'next-intl';
import SafeImage from '../../SafeImage';
import type { components } from '@/types/generated/api';

// Use generated types from API
type ProductVariant = components['schemas']['models.StorefrontProductVariant'] & {
  // Add additional fields that might be needed
  images?: ProductVariantImage[];
  stock_status?: string;
  reserved_quantity?: number;
  available_quantity?: number;
  is_default?: boolean;
  variant_attributes?: Record<string, any>;
  attributes?: Record<string, any>;
};

interface ProductVariantImage {
  id: number;
  image_url: string;
  thumbnail_url?: string;
  alt_text?: string;
  is_main: boolean;
}

interface VariantSelectorProps {
  productId: number;
  storefrontSlug: string;
  basePrice: number;
  baseCurrency: string;
  onVariantChange: (variant: ProductVariant | null) => void;
  className?: string;
}

export default function VariantSelector({
  productId,
  storefrontSlug,
  basePrice,
  baseCurrency,
  onVariantChange,
  className = '',
}: VariantSelectorProps) {
  const t = useTranslations('storefronts');
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

  const loadVariants = useCallback(async () => {
    console.log(
      'VariantSelector: Loading variants for productId:',
      productId,
      'storefrontSlug:',
      storefrontSlug
    );
    try {
      const response = await fetch(
        `/api/v1/public/b2c/${storefrontSlug}/products/${productId}/variants`
      );
      console.log('VariantSelector: Response status:', response.status);
      if (response.ok) {
        const variantData = await response.json();
        console.log('VariantSelector: Loaded variants:', variantData);
        setVariants(variantData);

        // Extract available attributes from variant_attributes field
        const attributes: Record<string, Set<string>> = {};
        variantData.forEach((variant: ProductVariant) => {
          const variantAttrs =
            variant.variant_attributes || variant.variant_attributes;
          if (variantAttrs) {
            Object.entries(variantAttrs).forEach(([key, value]) => {
              if (!attributes[key]) {
                attributes[key] = new Set();
              }
              attributes[key].add(String(value));
            });
          }
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
          const defaultAttrs =
            defaultVariant.variant_attributes || defaultVariant.attributes;
          if (defaultAttrs) {
            // Convert attributes to string-based format for compatibility
            const stringAttributes: Record<string, string> = {};
            Object.entries(defaultAttrs).forEach(([key, value]) => {
              stringAttributes[key] = String(value);
            });
            setSelectedAttributes(stringAttributes);
          }
        }
      } else {
        console.error(
          'VariantSelector: Failed to load variants, response not ok:',
          response.status
        );
      }
    } catch (error) {
      console.error('VariantSelector: Failed to load variants:', error);
    } finally {
      setLoading(false);
    }
  }, [productId, storefrontSlug]);

  const findMatchingVariant = useCallback(
    (attributes: Record<string, string>): ProductVariant | null => {
      return (
        variants.find((variant) => {
          const variantAttrs =
            variant.variant_attributes || variant.variant_attributes;
          if (!variantAttrs) return false;
          return Object.entries(attributes).every(
            ([key, value]) => String(variantAttrs[key]) === value
          );
        }) || null
      );
    },
    [variants]
  );

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
    // Use available_quantity if present, otherwise fall back to stock_quantity
    return (
      selectedVariant?.available_quantity ??
      selectedVariant?.stock_quantity ??
      0
    );
  };

  const getStockStatus = (): string => {
    if (selectedVariant?.stock_status) {
      return selectedVariant.stock_status;
    }

    // Calculate status based on stock
    const stock = getCurrentStock();
    if (stock === 0) return 'out_of_stock';
    if (stock < 5) return 'low_stock';
    return 'in_stock';
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

  useEffect(() => {
    loadVariants();
  }, [loadVariants]);

  useEffect(() => {
    // Find matching variant when attributes change
    const matchingVariant = findMatchingVariant(selectedAttributes);
    setSelectedVariant(matchingVariant);
    onVariantChange(matchingVariant);
  }, [selectedAttributes, variants, findMatchingVariant, onVariantChange]);

  console.log('VariantSelector: Rendering with state:', {
    loading,
    variantsLength: variants.length,
    selectedVariant,
    availableAttributes: Object.keys(availableAttributes),
  });

  if (loading) {
    console.log('VariantSelector: Showing loading state');
    return (
      <div className={`animate-pulse space-y-4 ${className}`}>
        <div className="h-4 bg-gray-200 rounded w-1/4"></div>
        <div className="h-10 bg-gray-200 rounded"></div>
        <div className="h-4 bg-gray-200 rounded w-1/3"></div>
      </div>
    );
  }

  if (variants.length === 0) {
    console.log('VariantSelector: No variants available, returning null');
    return null; // No variants available
  }

  console.log('VariantSelector: Rendering variant selector UI');

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

            {attributeName.toLowerCase() === 'color' ? (
              // Color picker for color attributes
              <div className="flex flex-wrap gap-3">
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
                        group relative flex flex-col items-center gap-1 p-1 rounded-lg transition-all
                        ${!isAvailable ? 'opacity-50 cursor-not-allowed' : 'hover:bg-gray-50 cursor-pointer'}
                      `}
                      title={value}
                    >
                      <div
                        className={`
                          w-12 h-12 rounded-lg border-2 relative transition-all shadow-sm
                          ${isSelected ? 'border-blue-500 ring-2 ring-blue-200 scale-105' : 'border-gray-200'}
                          ${!isAvailable ? '' : 'group-hover:border-gray-300 group-hover:shadow-md'}
                        `}
                        style={{ backgroundColor: getColorHex(value) }}
                      >
                        {/* Checkmark for selected state */}
                        {isSelected && (
                          <div className="absolute inset-0 flex items-center justify-center">
                            <div className="w-5 h-5 bg-white rounded-full flex items-center justify-center shadow-sm">
                              <svg
                                className="w-3 h-3 text-blue-500"
                                fill="currentColor"
                                viewBox="0 0 20 20"
                              >
                                <path
                                  fillRule="evenodd"
                                  d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
                                  clipRule="evenodd"
                                />
                              </svg>
                            </div>
                          </div>
                        )}

                        {/* Cross for unavailable colors */}
                        {!isAvailable && (
                          <div className="absolute inset-0 flex items-center justify-center">
                            <div className="w-8 h-0.5 bg-red-500 rotate-45 shadow-sm"></div>
                            <div className="absolute w-8 h-0.5 bg-red-500 -rotate-45 shadow-sm"></div>
                          </div>
                        )}
                      </div>

                      {/* Color name label */}
                      <span
                        className={`
                        text-xs text-center capitalize leading-tight min-w-0 max-w-16
                        ${isSelected ? 'text-blue-600 font-medium' : 'text-gray-600'}
                        ${!isAvailable ? 'text-gray-400' : ''}
                      `}
                      >
                        {value}
                      </span>
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
                <option value="">
                  {t('variants.selectOption', {
                    attribute: getAttributeDisplayName(attributeName),
                  })}
                </option>
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
          <span className="text-sm text-gray-600">
            {t('variants.availability')}
          </span>
          <span
            className={`text-sm font-medium ${
              getStockStatus() === 'in_stock'
                ? 'text-green-600'
                : getStockStatus() === 'low_stock'
                  ? 'text-yellow-600'
                  : 'text-red-600'
            }`}
          >
            {getStockStatus() === 'in_stock' && t('variants.inStock')}
            {getStockStatus() === 'low_stock' && t('variants.lowStock')}
            {getStockStatus() === 'out_of_stock' && t('variants.outOfStock')}
          </span>
        </div>

        {getCurrentStock() > 0 && (
          <div className="space-y-1">
            <div className="text-sm text-gray-600">
              {t('variants.itemsAvailable', { count: getCurrentStock() })}
            </div>
            {selectedVariant?.reserved_quantity &&
              selectedVariant.reserved_quantity > 0 && (
                <div className="text-xs text-amber-600">
                  {t('variants.reserved', {
                    count: selectedVariant.reserved_quantity,
                  })}
                </div>
              )}
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
            {t('variants.moreImages')}
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

// Helper function to get color hex value with enhanced color mapping
function getColorHex(colorName: string): string {
  const colorMap: Record<string, string> = {
    // Basic colors
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

    // Extended colors from the backend data
    navy: '#000080',
    teal: '#008080',
    lime: '#00FF00',
    maroon: '#800000',
    olive: '#808000',
    aqua: '#00FFFF',
    fuchsia: '#FF00FF',
    beige: '#F5F5DC',
    khaki: '#F0E68C',
    coral: '#FF7F50',
    salmon: '#FA8072',
    crimson: '#DC143C',
    indigo: '#4B0082',
    violet: '#EE82EE',
    turquoise: '#40E0D0',

    // Russian color names (if needed)
    чёрный: '#000000',
    белый: '#FFFFFF',
    красный: '#FF0000',
    синий: '#0000FF',
    зелёный: '#00FF00',
    жёлтый: '#FFFF00',
    фиолетовый: '#800080',
    оранжевый: '#FFA500',
    розовый: '#FFC0CB',
    серый: '#808080',
    коричневый: '#A52A2A',
    золотой: '#FFD700',
    серебряный: '#C0C0C0',
  };

  return colorMap[colorName.toLowerCase()] || '#CCCCCC';
}
