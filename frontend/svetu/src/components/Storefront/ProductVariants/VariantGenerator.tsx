'use client';

import React, { useState, useEffect, useCallback } from 'react';
import { useTranslations } from 'next-intl';

interface AttributeValue {
  value: string;
  display_name: string;
  color_hex?: string;
  image_url?: string;
  is_custom: boolean;
}

interface StorefrontProductAttribute {
  id: number;
  attribute_id: number;
  custom_values: AttributeValue[];
  attribute: {
    name: string;
    display_name: string;
    type: string;
  };
}

interface VariantGeneratorProps {
  productId: number;
  basePrice: number;
  onGenerate: (variants: any[]) => void;
  onCancel: () => void;
}

export default function VariantGenerator({
  productId,
  basePrice: _basePrice,
  onGenerate,
  onCancel,
}: VariantGeneratorProps) {
  const t = useTranslations('storefront');
  const [attributes, setAttributes] = useState<StorefrontProductAttribute[]>(
    []
  );
  const [selectedValues, setSelectedValues] = useState<
    Record<string, string[]>
  >({});
  const [priceModifiers, setPriceModifiers] = useState<Record<string, number>>(
    {}
  );
  const [stockQuantities, setStockQuantities] = useState<
    Record<string, number>
  >({});
  const [defaultAttributes, setDefaultAttributes] = useState<
    Record<string, string>
  >({});
  const [loading, setLoading] = useState(true);
  const [generating, setGenerating] = useState(false);

  const loadProductAttributes = useCallback(async () => {
    try {
      const response = await fetch(
        `/api/v1/storefront/products/${productId}/attributes`
      );
      if (response.ok) {
        const attrs = await response.json();
        setAttributes(attrs);

        // Initialize selected values
        const initialSelected: Record<string, string[]> = {};
        attrs.forEach((attr: StorefrontProductAttribute) => {
          initialSelected[attr.attribute.name] = [];
        });
        setSelectedValues(initialSelected);
      }
    } catch (error) {
      console.error('Failed to load product attributes:', error);
    } finally {
      setLoading(false);
    }
  }, [productId]);

  useEffect(() => {
    loadProductAttributes();
  }, [productId, loadProductAttributes]);

  const toggleValueSelection = (attributeName: string, value: string) => {
    setSelectedValues((prev) => ({
      ...prev,
      [attributeName]: prev[attributeName]?.includes(value)
        ? prev[attributeName].filter((v) => v !== value)
        : [...(prev[attributeName] || []), value],
    }));
  };

  const updatePriceModifier = (value: string, modifier: number) => {
    setPriceModifiers((prev) => ({
      ...prev,
      [value]: modifier,
    }));
  };

  const updateStockQuantity = (combination: string, quantity: number) => {
    setStockQuantities((prev) => ({
      ...prev,
      [combination]: quantity,
    }));
  };

  const setDefaultAttribute = (attributeName: string, value: string) => {
    setDefaultAttributes((prev) => ({
      ...prev,
      [attributeName]: value,
    }));
  };

  const generateCombinations = () => {
    const attributeNames = Object.keys(selectedValues);
    const combinations: Record<string, string>[] = [];

    const generate = (index: number, current: Record<string, string>) => {
      if (index === attributeNames.length) {
        combinations.push({ ...current });
        return;
      }

      const attrName = attributeNames[index];
      const values = selectedValues[attrName] || [];

      for (const value of values) {
        current[attrName] = value;
        generate(index + 1, current);
      }
    };

    generate(0, {});
    return combinations;
  };

  const calculateTotalVariants = () => {
    return Object.values(selectedValues).reduce((total, values) => {
      return total * Math.max(values.length, 1);
    }, 1);
  };

  const handleGenerate = async () => {
    setGenerating(true);

    try {
      const request = {
        product_id: productId,
        attribute_matrix: selectedValues,
        price_modifiers: priceModifiers,
        stock_quantities: stockQuantities,
        default_attributes: defaultAttributes,
      };

      const response = await fetch('/api/v1/storefront/variants/generate', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(request),
      });

      if (response.ok) {
        const result = await response.json();
        onGenerate(result.variants);
      } else {
        const error = await response.json();
        alert(`Failed to generate variants: ${error.error}`);
      }
    } catch (error) {
      console.error('Failed to generate variants:', error);
      alert('Failed to generate variants');
    } finally {
      setGenerating(false);
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center p-8">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  const totalVariants = calculateTotalVariants();
  const combinations = generateCombinations();

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h3 className="text-lg font-semibold">{t('generate_variants')}</h3>
          <p className="text-sm text-gray-600">
            {t('total_variants_will_be_generated', { count: totalVariants })}
          </p>
        </div>
        <div className="space-x-2">
          <button
            onClick={onCancel}
            className="px-4 py-2 text-gray-600 border border-gray-300 rounded-md hover:bg-gray-50"
          >
            {t('cancel')}
          </button>
          <button
            onClick={handleGenerate}
            disabled={totalVariants === 0 || generating}
            className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50"
          >
            {generating ? t('generating') : t('generate')}
          </button>
        </div>
      </div>

      {/* Attribute Selection */}
      <div className="space-y-4">
        <h4 className="font-medium">{t('select_attribute_values')}</h4>
        {attributes.map((attr) => (
          <div key={attr.id} className="border border-gray-200 rounded-lg p-4">
            <h5 className="font-medium mb-3">{attr.attribute.display_name}</h5>
            <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-2">
              {attr.custom_values.map((value) => (
                <label
                  key={value.value}
                  className="flex items-center space-x-2 p-2 border rounded cursor-pointer hover:bg-gray-50"
                >
                  <input
                    type="checkbox"
                    checked={
                      selectedValues[attr.attribute.name]?.includes(
                        value.value
                      ) || false
                    }
                    onChange={() =>
                      toggleValueSelection(attr.attribute.name, value.value)
                    }
                    className="h-4 w-4 text-blue-600 rounded"
                  />
                  {attr.attribute.type === 'color' && value.color_hex && (
                    <div
                      className="w-4 h-4 rounded border"
                      style={{ backgroundColor: value.color_hex }}
                    />
                  )}
                  <span className="text-sm">{value.display_name}</span>
                  {value.is_custom && (
                    <span className="text-xs text-blue-600 bg-blue-100 px-1 rounded">
                      {t('custom')}
                    </span>
                  )}
                </label>
              ))}
            </div>

            {/* Default Value Selection */}
            {selectedValues[attr.attribute.name]?.length > 0 && (
              <div className="mt-3">
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  {t('default_value')}
                </label>
                <select
                  value={defaultAttributes[attr.attribute.name] || ''}
                  onChange={(e) =>
                    setDefaultAttribute(attr.attribute.name, e.target.value)
                  }
                  className="px-3 py-2 border border-gray-300 rounded-md text-sm"
                >
                  <option value="">{t('no_default')}</option>
                  {selectedValues[attr.attribute.name].map((value) => {
                    const valueObj = attr.custom_values.find(
                      (v) => v.value === value
                    );
                    return (
                      <option key={value} value={value}>
                        {valueObj?.display_name || value}
                      </option>
                    );
                  })}
                </select>
              </div>
            )}
          </div>
        ))}
      </div>

      {/* Price Modifiers */}
      {Object.values(selectedValues).some((values) => values.length > 0) && (
        <div className="space-y-4">
          <h4 className="font-medium">{t('price_modifiers')}</h4>
          <p className="text-sm text-gray-600">
            {t('price_modifiers_description')}
          </p>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {Object.entries(selectedValues).flatMap(([attrName, values]) =>
              values.map((value) => {
                const attr = attributes.find(
                  (a) => a.attribute.name === attrName
                );
                const valueObj = attr?.custom_values.find(
                  (v) => v.value === value
                );
                return (
                  <div
                    key={`${attrName}-${value}`}
                    className="flex items-center space-x-2"
                  >
                    <span className="text-sm flex-1">
                      {valueObj?.display_name || value}
                    </span>
                    <div className="flex items-center space-x-1">
                      <span className="text-sm">+</span>
                      <input
                        type="number"
                        step="0.01"
                        value={priceModifiers[value] || 0}
                        onChange={(e) =>
                          updatePriceModifier(
                            value,
                            parseFloat(e.target.value) || 0
                          )
                        }
                        className="w-20 px-2 py-1 border border-gray-300 rounded text-sm"
                      />
                      <span className="text-sm">RSD</span>
                    </div>
                  </div>
                );
              })
            )}
          </div>
        </div>
      )}

      {/* Stock Quantities Preview */}
      {combinations.length > 0 && combinations.length <= 20 && (
        <div className="space-y-4">
          <h4 className="font-medium">{t('stock_quantities')}</h4>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            {combinations.map((combination, index) => {
              const key = Object.values(combination).join('-');
              const displayName = Object.entries(combination)
                .map(([attr, value]) => {
                  const attrObj = attributes.find(
                    (a) => a.attribute.name === attr
                  );
                  const valueObj = attrObj?.custom_values.find(
                    (v) => v.value === value
                  );
                  return valueObj?.display_name || value;
                })
                .join(' • ');

              return (
                <div key={index} className="flex items-center space-x-2">
                  <span className="text-sm flex-1">{displayName}</span>
                  <input
                    type="number"
                    min="0"
                    value={stockQuantities[key] || 10}
                    onChange={(e) =>
                      updateStockQuantity(key, parseInt(e.target.value) || 0)
                    }
                    className="w-20 px-2 py-1 border border-gray-300 rounded text-sm"
                  />
                </div>
              );
            })}
          </div>
        </div>
      )}

      {combinations.length > 20 && (
        <div className="bg-yellow-50 border border-yellow-200 rounded-lg p-4">
          <p className="text-sm text-yellow-800">
            {t('too_many_combinations_warning', { count: combinations.length })}
          </p>
        </div>
      )}
    </div>
  );
}
