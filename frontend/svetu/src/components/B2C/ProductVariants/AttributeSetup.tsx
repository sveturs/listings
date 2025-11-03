'use client';

import React, { useState, useEffect, useCallback } from 'react';
import { useTranslations } from 'next-intl';
import { apiClient } from '@/services/api-client';

interface AttributeValue {
  value: string;
  display_name: string;
  color_hex?: string;
  image_url?: string;
  is_custom: boolean;
}

interface ProductVariantAttribute {
  id: number;
  name: string;
  display_name: string;
  type: string;
  is_required: boolean;
  sort_order: number;
}

interface ProductAttributeSetup {
  attribute_id: number;
  is_enabled: boolean;
  is_required: boolean;
  custom_values: AttributeValue[];
  selected_global_values: string[];
}

interface AttributeSetupProps {
  productId: number;
  categoryId: number;
  onSave: (attributes: ProductAttributeSetup[]) => void;
  onCancel: () => void;
}

export default function AttributeSetup({
  productId,
  categoryId,
  onSave,
  onCancel,
}: AttributeSetupProps) {
  const t = useTranslations('storefronts');
  const [availableAttributes, setAvailableAttributes] = useState<
    ProductVariantAttribute[]
  >([]);
  const [attributeSetups, setAttributeSetups] = useState<
    ProductAttributeSetup[]
  >([]);
  const [globalValues, setGlobalValues] = useState<
    Record<number, AttributeValue[]>
  >({});
  const [loading, setLoading] = useState(true);

  const loadAvailableAttributes = useCallback(async () => {
    try {
      const response = await apiClient.get(
        `/marketplace/storefronts/storefront/categories/${categoryId}/attributes`
      );
      if (response.data) {
        const attributes = response.data;
        setAvailableAttributes(attributes);

        // Load global values for each attribute
        for (const attr of attributes) {
          loadGlobalValues(attr.id);
        }
      } else if (response.error) {
        console.error('Failed to load available attributes:', response.error);
      }
    } catch (error) {
      console.error('Failed to load available attributes:', error);
    }
  }, [categoryId]);

  const loadCurrentSetup = useCallback(async () => {
    try {
      const response = await apiClient.get(
        `/marketplace/storefronts/storefront/products/${productId}/attributes`
      );
      if (response.data) {
        const currentAttributes = response.data;
        const setups = availableAttributes.map((attr) => {
          const existing = currentAttributes.find(
            (ca: any) => ca.attribute_id === attr.id
          );
          return existing
            ? {
                attribute_id: attr.id,
                is_enabled: existing.is_enabled,
                is_required: existing.is_required,
                custom_values: existing.custom_values || [],
                selected_global_values:
                  existing.custom_values
                    ?.filter((v: AttributeValue) => !v.is_custom)
                    ?.map((v: AttributeValue) => v.value) || [],
              }
            : {
                attribute_id: attr.id,
                is_enabled: false,
                is_required: false,
                custom_values: [],
                selected_global_values: [],
              };
        });
        setAttributeSetups(setups);
      } else if (response.error) {
        console.error('Failed to load current setup:', response.error);
      }
    } catch (error) {
      console.error('Failed to load current setup:', error);
    } finally {
      setLoading(false);
    }
  }, [productId, availableAttributes]);

  useEffect(() => {
    loadAvailableAttributes();
    loadCurrentSetup();
  }, [productId, categoryId, loadAvailableAttributes, loadCurrentSetup]);

  const loadGlobalValues = async (attributeId: number) => {
    try {
      const response = await apiClient.get(
        `/public/variants/attributes/${attributeId}/values`
      );
      if (response.data) {
        const values = response.data;
        setGlobalValues((prev) => ({
          ...prev,
          [attributeId]: values.map((v: any) => ({
            value: v.value,
            display_name: v.display_name,
            color_hex: v.color_hex,
            image_url: v.image_url,
            is_custom: false,
          })),
        }));
      } else if (response.error) {
        console.error(
          `Failed to load global values for attribute ${attributeId}:`,
          response.error
        );
      }
    } catch (error) {
      console.error(
        `Failed to load global values for attribute ${attributeId}:`,
        error
      );
    }
  };

  const updateAttributeSetup = (
    attributeId: number,
    updates: Partial<ProductAttributeSetup>
  ) => {
    setAttributeSetups((prev) =>
      prev.map((setup) =>
        setup.attribute_id === attributeId ? { ...setup, ...updates } : setup
      )
    );
  };

  const addCustomValue = (attributeId: number, value: AttributeValue) => {
    updateAttributeSetup(attributeId, {
      custom_values: [
        ...(attributeSetups.find((s) => s.attribute_id === attributeId)
          ?.custom_values || []),
        { ...value, is_custom: true },
      ],
    });
  };

  const removeCustomValue = (attributeId: number, valueIndex: number) => {
    const setup = attributeSetups.find((s) => s.attribute_id === attributeId);
    if (setup) {
      updateAttributeSetup(attributeId, {
        custom_values: setup.custom_values.filter(
          (_, index) => index !== valueIndex
        ),
      });
    }
  };

  const toggleGlobalValue = (attributeId: number, value: string) => {
    const setup = attributeSetups.find((s) => s.attribute_id === attributeId);
    if (setup) {
      const isSelected = setup.selected_global_values.includes(value);
      updateAttributeSetup(attributeId, {
        selected_global_values: isSelected
          ? setup.selected_global_values.filter((v) => v !== value)
          : [...setup.selected_global_values, value],
      });
    }
  };

  const handleSave = () => {
    const enabledSetups = attributeSetups.filter((setup) => setup.is_enabled);
    onSave(enabledSetups);
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center p-8">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h3 className="text-lg font-semibold">
          {t('setup_product_attributes')}
        </h3>
        <div className="space-x-2">
          <button
            onClick={onCancel}
            className="px-4 py-2 text-gray-600 border border-gray-300 rounded-md hover:bg-gray-50"
          >
            {t('cancel')}
          </button>
          <button
            onClick={handleSave}
            className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
          >
            {t('save')}
          </button>
        </div>
      </div>

      <div className="space-y-4">
        {availableAttributes.map((attribute) => {
          const setup = attributeSetups.find(
            (s) => s.attribute_id === attribute.id
          );
          const globalVals = globalValues[attribute.id] || [];

          return (
            <div
              key={attribute.id}
              className="border border-gray-200 rounded-lg p-4"
            >
              <div className="flex items-center justify-between mb-4">
                <div className="flex items-center space-x-3">
                  <input
                    type="checkbox"
                    checked={setup?.is_enabled || false}
                    onChange={(e) =>
                      updateAttributeSetup(attribute.id, {
                        is_enabled: e.target.checked,
                      })
                    }
                    className="h-4 w-4 text-blue-600 rounded"
                  />
                  <div>
                    <h4 className="font-medium">{attribute.display_name}</h4>
                    <p className="text-sm text-gray-500">
                      {attribute.name} • {attribute.type}
                    </p>
                  </div>
                </div>
                {setup?.is_enabled && (
                  <label className="flex items-center space-x-2">
                    <input
                      type="checkbox"
                      checked={setup.is_required}
                      onChange={(e) =>
                        updateAttributeSetup(attribute.id, {
                          is_required: e.target.checked,
                        })
                      }
                      className="h-4 w-4 text-red-600 rounded"
                    />
                    <span className="text-sm text-red-600">
                      {t('required')}
                    </span>
                  </label>
                )}
              </div>

              {setup?.is_enabled && (
                <div className="space-y-4">
                  {/* Global Values */}
                  {globalVals.length > 0 && (
                    <div>
                      <h5 className="font-medium mb-2">{t('global_values')}</h5>
                      <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-2">
                        {globalVals.map((value) => (
                          <label
                            key={value.value}
                            className="flex items-center space-x-2 p-2 border rounded cursor-pointer hover:bg-gray-50"
                          >
                            <input
                              type="checkbox"
                              checked={setup.selected_global_values.includes(
                                value.value
                              )}
                              onChange={() =>
                                toggleGlobalValue(attribute.id, value.value)
                              }
                              className="h-4 w-4 text-blue-600 rounded"
                            />
                            {attribute.type === 'color' && value.color_hex && (
                              <div
                                className="w-4 h-4 rounded border"
                                style={{ backgroundColor: value.color_hex }}
                              />
                            )}
                            <span className="text-sm">
                              {value.display_name}
                            </span>
                          </label>
                        ))}
                      </div>
                    </div>
                  )}

                  {/* Custom Values */}
                  <div>
                    <h5 className="font-medium mb-2">{t('custom_values')}</h5>
                    <div className="space-y-2">
                      {setup.custom_values.map((value, index) => (
                        <div
                          key={index}
                          className="flex items-center space-x-2 p-2 bg-blue-50 rounded"
                        >
                          {attribute.type === 'color' && value.color_hex && (
                            <div
                              className="w-4 h-4 rounded border"
                              style={{ backgroundColor: value.color_hex }}
                            />
                          )}
                          <span className="flex-1">{value.display_name}</span>
                          <button
                            onClick={() =>
                              removeCustomValue(attribute.id, index)
                            }
                            className="text-red-600 hover:text-red-800"
                          >
                            ×
                          </button>
                        </div>
                      ))}
                      <CustomValueInput
                        attributeType={attribute.type}
                        onAdd={(value) => addCustomValue(attribute.id, value)}
                      />
                    </div>
                  </div>
                </div>
              )}
            </div>
          );
        })}
      </div>
    </div>
  );
}

// Component for adding custom values
interface CustomValueInputProps {
  attributeType: string;
  onAdd: (value: AttributeValue) => void;
}

function CustomValueInput({ attributeType, onAdd }: CustomValueInputProps) {
  const t = useTranslations('storefronts');
  const [value, setValue] = useState('');
  const [displayName, setDisplayName] = useState('');
  const [colorHex, setColorHex] = useState('#000000');

  const handleAdd = () => {
    if (!value.trim() || !displayName.trim()) return;

    onAdd({
      value: value.trim(),
      display_name: displayName.trim(),
      color_hex: attributeType === 'color' ? colorHex : undefined,
      is_custom: true,
    });

    setValue('');
    setDisplayName('');
    setColorHex('#000000');
  };

  return (
    <div className="flex items-center space-x-2 p-2 border border-dashed border-gray-300 rounded">
      <input
        type="text"
        placeholder={t('value')}
        value={value}
        onChange={(e) => setValue(e.target.value)}
        className="flex-1 px-2 py-1 border border-gray-300 rounded text-sm"
      />
      <input
        type="text"
        placeholder={t('display_name')}
        value={displayName}
        onChange={(e) => setDisplayName(e.target.value)}
        className="flex-1 px-2 py-1 border border-gray-300 rounded text-sm"
      />
      {attributeType === 'color' && (
        <input
          type="color"
          value={colorHex}
          onChange={(e) => setColorHex(e.target.value)}
          className="w-8 h-8 border border-gray-300 rounded"
        />
      )}
      <button
        onClick={handleAdd}
        disabled={!value.trim() || !displayName.trim()}
        className="px-3 py-1 bg-blue-600 text-white rounded text-sm hover:bg-blue-700 disabled:opacity-50"
      >
        {t('add')}
      </button>
    </div>
  );
}
