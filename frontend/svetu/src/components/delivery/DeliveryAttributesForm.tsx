'use client';

import { useState, useEffect } from 'react';
import {
  CubeIcon,
  ScaleIcon,
  ExclamationTriangleIcon,
  InformationCircleIcon,
  CheckIcon,
  SparklesIcon,
  ArchiveBoxIcon,
} from '@heroicons/react/24/outline';
import { useTranslations } from 'next-intl';
import {
  DeliveryAttributes,
  ValidationErrors,
  CategoryDefaults,
} from '@/types/delivery';

interface Props {
  attributes: Partial<DeliveryAttributes>;
  onChange: (attributes: Partial<DeliveryAttributes>) => void;
  categoryId?: number;
  showDefaults?: boolean;
  compact?: boolean;
  errors?: ValidationErrors;
  className?: string;
}

const PACKAGING_TYPES = [
  {
    value: 'box',
    label: 'Коробка',
    icon: ArchiveBoxIcon,
    description: 'Стандартная упаковка для большинства товаров',
  },
  {
    value: 'envelope',
    label: 'Конверт',
    icon: ArchiveBoxIcon,
    description: 'Для плоских и легких предметов',
  },
  {
    value: 'pallet',
    label: 'Паллета',
    icon: ArchiveBoxIcon,
    description: 'Для крупногабаритных и тяжелых товаров',
  },
  {
    value: 'custom',
    label: 'Особая упаковка',
    icon: SparklesIcon,
    description: 'Нестандартная упаковка',
  },
] as const;

export default function DeliveryAttributesForm({
  attributes,
  onChange,
  categoryId,
  showDefaults = true,
  compact = false,
  errors,
  className = '',
}: Props) {
  const t = useTranslations('delivery');
  const [categoryDefaults, setCategoryDefaults] =
    useState<CategoryDefaults | null>(null);
  const [loading, setLoading] = useState(false);
  const [showVolumeCalculation, setShowVolumeCalculation] = useState(false);

  // Загрузка дефолтных значений для категории
  useEffect(() => {
    if (categoryId && showDefaults) {
      loadCategoryDefaults();
    }
  }, [categoryId, showDefaults]);

  const loadCategoryDefaults = async () => {
    if (!categoryId) return;

    setLoading(true);
    try {
      const response = await fetch(
        `/api/v1/categories/${categoryId}/delivery-defaults`
      );
      const data = await response.json();

      if (data.success && data.data) {
        setCategoryDefaults(data.data);
      }
    } catch (error) {
      console.error('Failed to load category defaults:', error);
    } finally {
      setLoading(false);
    }
  };

  const applyDefaults = () => {
    if (!categoryDefaults) return;

    const defaultAttributes: Partial<DeliveryAttributes> = {
      weight_kg: categoryDefaults.default_weight_kg || 1,
      dimensions: {
        length_cm: categoryDefaults.default_length_cm || 30,
        width_cm: categoryDefaults.default_width_cm || 20,
        height_cm: categoryDefaults.default_height_cm || 10,
      },
      packaging_type: (categoryDefaults.default_packaging_type as any) || 'box',
      is_fragile: categoryDefaults.is_typically_fragile || false,
      requires_special_handling: false,
      stackable: !categoryDefaults.is_typically_fragile,
    };

    onChange({ ...attributes, ...defaultAttributes });
  };

  const calculateVolume = () => {
    const { dimensions } = attributes;
    if (
      !dimensions?.length_cm ||
      !dimensions?.width_cm ||
      !dimensions?.height_cm
    ) {
      return 0;
    }
    return (
      (dimensions.length_cm * dimensions.width_cm * dimensions.height_cm) /
      1000000
    ); // Convert to m³
  };

  const calculateVolumetricWeight = () => {
    const volume = calculateVolume();
    return volume * 250; // Standard volumetric weight calculation (250 kg/m³)
  };

  const updateAttribute = (field: keyof DeliveryAttributes, value: any) => {
    const updated = { ...attributes, [field]: value };

    // Auto-calculate volume if dimensions change
    if (field === 'dimensions') {
      const volume = calculateVolume();
      updated.volume_m3 = volume;
    }

    onChange(updated);
  };

  const updateDimension = (
    dimension: 'length_cm' | 'width_cm' | 'height_cm',
    value: number
  ) => {
    const dimensions = {
      ...attributes.dimensions,
      [dimension]: value,
    };
    updateAttribute('dimensions', dimensions);
  };

  const getFieldError = (field: string) => {
    if (!errors) return undefined;
    return field.includes('.')
      ? field.split('.').reduce((obj, key) => obj?.[key], errors as any)
      : errors[field];
  };

  const isVolumetricWeightHigher = () => {
    const actualWeight = attributes.weight_kg || 0;
    const volumetricWeight = calculateVolumetricWeight();
    return volumetricWeight > actualWeight;
  };

  if (compact) {
    return (
      <div className={`space-y-4 ${className}`}>
        <div className="grid grid-cols-2 gap-3">
          {/* Weight */}
          <div className="form-control">
            <label className="label py-1">
              <span className="label-text text-sm">Вес (кг)</span>
            </label>
            <input
              type="number"
              min="0.1"
              max="100"
              step="0.1"
              className={`input input-sm input-bordered ${getFieldError('weight_kg') ? 'input-error' : ''}`}
              value={attributes.weight_kg || ''}
              onChange={(e) =>
                updateAttribute('weight_kg', parseFloat(e.target.value) || 0)
              }
              placeholder="0.5"
            />
          </div>

          {/* Packaging */}
          <div className="form-control">
            <label className="label py-1">
              <span className="label-text text-sm">Упаковка</span>
            </label>
            <select
              className="select select-sm select-bordered"
              value={attributes.packaging_type || 'box'}
              onChange={(e) =>
                updateAttribute('packaging_type', e.target.value)
              }
            >
              {PACKAGING_TYPES.map((type) => (
                <option key={type.value} value={type.value}>
                  {type.label}
                </option>
              ))}
            </select>
          </div>
        </div>

        {/* Dimensions */}
        <div className="form-control">
          <label className="label py-1">
            <span className="label-text text-sm">Размеры (см)</span>
          </label>
          <div className="grid grid-cols-3 gap-2">
            <input
              type="number"
              min="1"
              max="300"
              className={`input input-sm input-bordered ${getFieldError('dimensions.length_cm') ? 'input-error' : ''}`}
              placeholder="Длина"
              value={attributes.dimensions?.length_cm || ''}
              onChange={(e) =>
                updateDimension('length_cm', parseFloat(e.target.value) || 0)
              }
            />
            <input
              type="number"
              min="1"
              max="300"
              className={`input input-sm input-bordered ${getFieldError('dimensions.width_cm') ? 'input-error' : ''}`}
              placeholder="Ширина"
              value={attributes.dimensions?.width_cm || ''}
              onChange={(e) =>
                updateDimension('width_cm', parseFloat(e.target.value) || 0)
              }
            />
            <input
              type="number"
              min="1"
              max="300"
              className={`input input-sm input-bordered ${getFieldError('dimensions.height_cm') ? 'input-error' : ''}`}
              placeholder="Высота"
              value={attributes.dimensions?.height_cm || ''}
              onChange={(e) =>
                updateDimension('height_cm', parseFloat(e.target.value) || 0)
              }
            />
          </div>
        </div>

        {/* Fragile toggle */}
        <label className="label cursor-pointer justify-start gap-2">
          <input
            type="checkbox"
            className="checkbox checkbox-sm checkbox-primary"
            checked={attributes.is_fragile || false}
            onChange={(e) => updateAttribute('is_fragile', e.target.checked)}
          />
          <span className="label-text text-sm">Хрупкий товар</span>
        </label>
      </div>
    );
  }

  return (
    <div className={`space-y-6 ${className}`}>
      {/* Header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div className="p-2 bg-primary/10 rounded-lg">
            <CubeIcon className="w-6 h-6 text-primary" />
          </div>
          <div>
            <h3 className="text-lg font-semibold">Параметры доставки</h3>
            <p className="text-sm text-base-content/60">
              Укажите вес, размеры и особенности товара для точного расчета
              стоимости
            </p>
          </div>
        </div>

        {/* Apply defaults button */}
        {categoryDefaults && (
          <button
            type="button"
            className="btn btn-sm btn-outline btn-primary"
            onClick={applyDefaults}
            disabled={loading}
          >
            <SparklesIcon className="w-4 h-4" />
            Применить стандартные
          </button>
        )}
      </div>

      <div className="grid lg:grid-cols-2 gap-6">
        {/* Basic Parameters */}
        <div className="space-y-4">
          <h4 className="font-medium text-base-content/80 flex items-center gap-2">
            <ScaleIcon className="w-4 h-4" />
            Основные параметры
          </h4>

          {/* Weight */}
          <div className="form-control">
            <label className="label">
              <span className="label-text font-medium">Вес (кг) *</span>
              <span className="label-text-alt text-xs">до 100 кг</span>
            </label>
            <input
              type="number"
              min="0.1"
              max="100"
              step="0.1"
              className={`input input-bordered ${getFieldError('weight_kg') ? 'input-error' : 'focus:input-primary'}`}
              value={attributes.weight_kg || ''}
              onChange={(e) =>
                updateAttribute('weight_kg', parseFloat(e.target.value) || 0)
              }
              placeholder="0.5"
            />
            {getFieldError('weight_kg') && (
              <label className="label">
                <span className="label-text-alt text-error">
                  {getFieldError('weight_kg')}
                </span>
              </label>
            )}
          </div>

          {/* Dimensions */}
          <div className="form-control">
            <label className="label">
              <span className="label-text font-medium">Размеры (см) *</span>
              <button
                type="button"
                className="label-text-alt text-primary cursor-pointer hover:underline"
                onClick={() => setShowVolumeCalculation(!showVolumeCalculation)}
              >
                Показать расчеты
              </button>
            </label>
            <div className="grid grid-cols-3 gap-3">
              <div>
                <input
                  type="number"
                  min="1"
                  max="300"
                  className={`input input-bordered ${getFieldError('dimensions.length_cm') ? 'input-error' : 'focus:input-primary'}`}
                  placeholder="Длина"
                  value={attributes.dimensions?.length_cm || ''}
                  onChange={(e) =>
                    updateDimension(
                      'length_cm',
                      parseFloat(e.target.value) || 0
                    )
                  }
                />
                <label className="label">
                  <span className="label-text-alt text-xs">Длина</span>
                </label>
              </div>
              <div>
                <input
                  type="number"
                  min="1"
                  max="300"
                  className={`input input-bordered ${getFieldError('dimensions.width_cm') ? 'input-error' : 'focus:input-primary'}`}
                  placeholder="Ширина"
                  value={attributes.dimensions?.width_cm || ''}
                  onChange={(e) =>
                    updateDimension('width_cm', parseFloat(e.target.value) || 0)
                  }
                />
                <label className="label">
                  <span className="label-text-alt text-xs">Ширина</span>
                </label>
              </div>
              <div>
                <input
                  type="number"
                  min="1"
                  max="300"
                  className={`input input-bordered ${getFieldError('dimensions.height_cm') ? 'input-error' : 'focus:input-primary'}`}
                  placeholder="Высота"
                  value={attributes.dimensions?.height_cm || ''}
                  onChange={(e) =>
                    updateDimension(
                      'height_cm',
                      parseFloat(e.target.value) || 0
                    )
                  }
                />
                <label className="label">
                  <span className="label-text-alt text-xs">Высота</span>
                </label>
              </div>
            </div>

            {/* Volume calculations */}
            {showVolumeCalculation && calculateVolume() > 0 && (
              <div className="mt-3 p-3 bg-base-200 rounded-lg text-sm space-y-2">
                <div className="flex justify-between">
                  <span>Объем:</span>
                  <span className="font-medium">
                    {calculateVolume().toFixed(3)} м³
                  </span>
                </div>
                <div className="flex justify-between">
                  <span>Объемный вес:</span>
                  <span
                    className={`font-medium ${isVolumetricWeightHigher() ? 'text-warning' : ''}`}
                  >
                    {calculateVolumetricWeight().toFixed(2)} кг
                  </span>
                </div>
                {isVolumetricWeightHigher() && (
                  <div className="alert alert-warning alert-sm">
                    <ExclamationTriangleIcon className="w-4 h-4" />
                    <span className="text-xs">
                      Объемный вес больше фактического. Будет использован для
                      расчета.
                    </span>
                  </div>
                )}
              </div>
            )}
          </div>

          {/* Packaging Type */}
          <div className="form-control">
            <label className="label">
              <span className="label-text font-medium">Тип упаковки</span>
            </label>
            <div className="grid grid-cols-2 gap-2">
              {PACKAGING_TYPES.map((type) => {
                const Icon = type.icon;
                const isSelected = attributes.packaging_type === type.value;

                return (
                  <div
                    key={type.value}
                    className={`
                      card cursor-pointer transition-all border-2 hover:border-primary/30
                      ${isSelected ? 'border-primary bg-primary/5' : 'border-transparent'}
                    `}
                    onClick={() =>
                      updateAttribute('packaging_type', type.value)
                    }
                  >
                    <div className="card-body p-3">
                      <div className="flex items-center gap-2">
                        <Icon className="w-5 h-5 text-primary" />
                        <div className="flex-1">
                          <div className="font-medium text-sm">
                            {type.label}
                          </div>
                          <div className="text-xs text-base-content/60">
                            {type.description}
                          </div>
                        </div>
                        {isSelected && (
                          <CheckIcon className="w-4 h-4 text-success" />
                        )}
                      </div>
                    </div>
                  </div>
                );
              })}
            </div>
          </div>
        </div>

        {/* Special Properties */}
        <div className="space-y-4">
          <h4 className="font-medium text-base-content/80 flex items-center gap-2">
            <ExclamationTriangleIcon className="w-4 h-4" />
            Особые свойства
          </h4>

          <div className="space-y-4">
            {/* Fragile */}
            <label className="flex items-start gap-3 cursor-pointer p-3 rounded-lg hover:bg-base-200 transition-colors">
              <input
                type="checkbox"
                className="checkbox checkbox-primary mt-0.5"
                checked={attributes.is_fragile || false}
                onChange={(e) =>
                  updateAttribute('is_fragile', e.target.checked)
                }
              />
              <div className="flex-1">
                <div className="font-medium">Хрупкий товар</div>
                <div className="text-sm text-base-content/60">
                  Требует особой осторожности при транспортировке. Может
                  увеличить стоимость доставки.
                </div>
              </div>
            </label>

            {/* Special Handling */}
            <label className="flex items-start gap-3 cursor-pointer p-3 rounded-lg hover:bg-base-200 transition-colors">
              <input
                type="checkbox"
                className="checkbox checkbox-primary mt-0.5"
                checked={attributes.requires_special_handling || false}
                onChange={(e) =>
                  updateAttribute('requires_special_handling', e.target.checked)
                }
              />
              <div className="flex-1">
                <div className="font-medium">Требует особого обращения</div>
                <div className="text-sm text-base-content/60">
                  Для товаров с нестандартными требованиями к перевозке.
                </div>
              </div>
            </label>

            {/* Stackable */}
            <label className="flex items-start gap-3 cursor-pointer p-3 rounded-lg hover:bg-base-200 transition-colors">
              <input
                type="checkbox"
                className="checkbox checkbox-primary mt-0.5"
                checked={attributes.stackable !== false} // Default true
                onChange={(e) => updateAttribute('stackable', e.target.checked)}
              />
              <div className="flex-1">
                <div className="font-medium">Можно штабелировать</div>
                <div className="text-sm text-base-content/60">
                  На товар можно ставить другие посылки при транспортировке.
                </div>
              </div>
            </label>

            {/* Max stack weight (if stackable) */}
            {attributes.stackable !== false && (
              <div className="form-control ml-6">
                <label className="label">
                  <span className="label-text">
                    Максимальный вес сверху (кг)
                  </span>
                </label>
                <input
                  type="number"
                  min="0"
                  max="100"
                  className="input input-bordered input-sm focus:input-primary"
                  value={attributes.max_stack_weight_kg || ''}
                  onChange={(e) =>
                    updateAttribute(
                      'max_stack_weight_kg',
                      parseFloat(e.target.value) || undefined
                    )
                  }
                  placeholder="50"
                />
                <label className="label">
                  <span className="label-text-alt">
                    Оставьте пустым для стандартного лимита
                  </span>
                </label>
              </div>
            )}
          </div>

          {/* Info panel */}
          <div className="alert alert-info">
            <InformationCircleIcon className="w-5 h-5" />
            <div className="text-sm">
              <div className="font-semibold mb-1">Советы по заполнению:</div>
              <ul className="list-disc list-inside space-y-1 text-xs">
                <li>Указывайте реальные размеры упакованного товара</li>
                <li>Хрупкие товары требуют дополнительной упаковки</li>
                <li>Точные параметры помогают избежать доплат при доставке</li>
              </ul>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
