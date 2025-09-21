'use client';

import {
  CubeIcon,
  ScaleIcon,
  ArchiveBoxIcon,
  ExclamationTriangleIcon,
  InformationCircleIcon,
  TruckIcon,
} from '@heroicons/react/24/outline';
import { DeliveryAttributes } from '@/types/delivery';

interface Props {
  attributes: DeliveryAttributes;
  showCalculatedValues?: boolean;
  compact?: boolean;
  className?: string;
}

const PACKAGING_LABELS: Record<string, string> = {
  box: 'Коробка',
  envelope: 'Конверт',
  pallet: 'Паллета',
  custom: 'Особая упаковка',
};

export default function DeliveryAttributesDisplay({
  attributes,
  showCalculatedValues = true,
  compact = false,
  className = '',
}: Props) {
  const calculateVolume = () => {
    const { dimensions } = attributes;
    return (
      (dimensions.length_cm * dimensions.width_cm * dimensions.height_cm) /
      1000000
    ); // Convert to m³
  };

  const calculateVolumetricWeight = () => {
    const volume = calculateVolume();
    return volume * 250; // Standard volumetric weight calculation
  };

  const isVolumetricWeightHigher = () => {
    return calculateVolumetricWeight() > attributes.weight_kg;
  };

  const getBillingWeight = () => {
    return Math.max(attributes.weight_kg, calculateVolumetricWeight());
  };

  if (compact) {
    return (
      <div className={`space-y-2 ${className}`}>
        <div className="flex items-center gap-4 text-sm">
          <div className="flex items-center gap-1">
            <ScaleIcon className="w-4 h-4 text-base-content/60" />
            <span>{attributes.weight_kg} кг</span>
          </div>
          <div className="flex items-center gap-1">
            <CubeIcon className="w-4 h-4 text-base-content/60" />
            <span>
              {attributes.dimensions.length_cm}×{attributes.dimensions.width_cm}
              ×{attributes.dimensions.height_cm} см
            </span>
          </div>
        </div>

        {(attributes.is_fragile || attributes.requires_special_handling) && (
          <div className="flex gap-2">
            {attributes.is_fragile && (
              <span className="badge badge-warning badge-xs">Хрупкое</span>
            )}
            {attributes.requires_special_handling && (
              <span className="badge badge-info badge-xs">
                Особое обращение
              </span>
            )}
          </div>
        )}
      </div>
    );
  }

  return (
    <div className={`space-y-4 ${className}`}>
      {/* Header */}
      <div className="flex items-center gap-3">
        <div className="p-2 bg-primary/10 rounded-lg">
          <TruckIcon className="w-5 h-5 text-primary" />
        </div>
        <div>
          <h3 className="font-semibold">Параметры доставки</h3>
          <p className="text-sm text-base-content/60">
            Характеристики товара для расчета стоимости
          </p>
        </div>
      </div>

      <div className="grid md:grid-cols-2 gap-6">
        {/* Physical Properties */}
        <div className="space-y-4">
          <h4 className="font-medium text-base-content/80 flex items-center gap-2">
            <ScaleIcon className="w-4 h-4" />
            Физические характеристики
          </h4>

          <div className="space-y-3">
            {/* Weight */}
            <div className="flex justify-between items-center">
              <span className="text-sm text-base-content/70">Вес:</span>
              <span className="font-medium">{attributes.weight_kg} кг</span>
            </div>

            {/* Dimensions */}
            <div className="flex justify-between items-center">
              <span className="text-sm text-base-content/70">Размеры:</span>
              <span className="font-medium">
                {attributes.dimensions.length_cm} ×{' '}
                {attributes.dimensions.width_cm} ×{' '}
                {attributes.dimensions.height_cm} см
              </span>
            </div>

            {/* Volume */}
            <div className="flex justify-between items-center">
              <span className="text-sm text-base-content/70">Объем:</span>
              <span className="font-medium">
                {calculateVolume().toFixed(3)} м³
              </span>
            </div>

            {/* Packaging */}
            <div className="flex justify-between items-center">
              <span className="text-sm text-base-content/70">Упаковка:</span>
              <div className="flex items-center gap-2">
                <ArchiveBoxIcon className="w-4 h-4 text-base-content/60" />
                <span className="font-medium">
                  {PACKAGING_LABELS[attributes.packaging_type] ||
                    attributes.packaging_type}
                </span>
              </div>
            </div>
          </div>
        </div>

        {/* Special Properties */}
        <div className="space-y-4">
          <h4 className="font-medium text-base-content/80 flex items-center gap-2">
            <ExclamationTriangleIcon className="w-4 h-4" />
            Особые свойства
          </h4>

          <div className="space-y-3">
            {/* Fragile */}
            <div className="flex justify-between items-center">
              <span className="text-sm text-base-content/70">
                Хрупкий товар:
              </span>
              <span
                className={`badge badge-sm ${attributes.is_fragile ? 'badge-warning' : 'badge-ghost'}`}
              >
                {attributes.is_fragile ? 'Да' : 'Нет'}
              </span>
            </div>

            {/* Special Handling */}
            <div className="flex justify-between items-center">
              <span className="text-sm text-base-content/70">
                Особое обращение:
              </span>
              <span
                className={`badge badge-sm ${attributes.requires_special_handling ? 'badge-info' : 'badge-ghost'}`}
              >
                {attributes.requires_special_handling ? 'Да' : 'Нет'}
              </span>
            </div>

            {/* Stackable */}
            <div className="flex justify-between items-center">
              <span className="text-sm text-base-content/70">
                Штабелируемый:
              </span>
              <span
                className={`badge badge-sm ${attributes.stackable ? 'badge-success' : 'badge-ghost'}`}
              >
                {attributes.stackable ? 'Да' : 'Нет'}
              </span>
            </div>

            {/* Max Stack Weight */}
            {attributes.stackable && attributes.max_stack_weight_kg && (
              <div className="flex justify-between items-center">
                <span className="text-sm text-base-content/70">
                  Макс. вес сверху:
                </span>
                <span className="font-medium">
                  {attributes.max_stack_weight_kg} кг
                </span>
              </div>
            )}

            {/* Hazmat */}
            {attributes.hazmat_class && (
              <div className="flex justify-between items-center">
                <span className="text-sm text-base-content/70">
                  Класс опасности:
                </span>
                <span className="badge badge-error badge-sm">
                  {attributes.hazmat_class}
                </span>
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Calculated Values */}
      {showCalculatedValues && (
        <div className="card bg-base-200">
          <div className="card-body p-4">
            <h4 className="font-medium text-base-content/80 mb-3">
              Расчетные значения
            </h4>

            <div className="grid md:grid-cols-2 gap-4">
              <div className="space-y-2">
                <div className="flex justify-between items-center">
                  <span className="text-sm text-base-content/70">
                    Объемный вес:
                  </span>
                  <span
                    className={`font-medium ${isVolumetricWeightHigher() ? 'text-warning' : ''}`}
                  >
                    {calculateVolumetricWeight().toFixed(2)} кг
                  </span>
                </div>

                <div className="flex justify-between items-center">
                  <span className="text-sm text-base-content/70">
                    Расчетный вес:
                  </span>
                  <span className="font-bold text-primary">
                    {getBillingWeight().toFixed(2)} кг
                  </span>
                </div>
              </div>

              <div className="space-y-2">
                <div className="flex justify-between items-center">
                  <span className="text-sm text-base-content/70">
                    Тип груза:
                  </span>
                  <span className="font-medium">
                    {attributes.is_fragile ? 'Хрупкий' : 'Обычный'}
                  </span>
                </div>

                <div className="flex justify-between items-center">
                  <span className="text-sm text-base-content/70">
                    Категория:
                  </span>
                  <span className="font-medium">
                    {getBillingWeight() > 30
                      ? 'Тяжелый'
                      : getBillingWeight() > 10
                        ? 'Средний'
                        : 'Легкий'}
                  </span>
                </div>
              </div>
            </div>

            {isVolumetricWeightHigher() && (
              <div className="alert alert-warning alert-sm mt-3">
                <ExclamationTriangleIcon className="w-4 h-4" />
                <span className="text-xs">
                  Для расчета стоимости будет использован объемный вес как
                  больший
                </span>
              </div>
            )}
          </div>
        </div>
      )}

      {/* Info Panel */}
      <div className="card bg-gradient-to-r from-info/5 to-info/10">
        <div className="card-body p-4">
          <div className="flex items-start gap-3">
            <InformationCircleIcon className="w-5 h-5 text-info flex-shrink-0 mt-0.5" />
            <div className="space-y-2 text-sm">
              <div className="font-semibold">
                Влияние на стоимость доставки:
              </div>
              <ul className="list-disc list-inside space-y-1 text-xs">
                <li>
                  Стоимость рассчитывается по большему из физического или
                  объемного веса
                </li>
                {attributes.is_fragile && (
                  <li className="text-warning">
                    Хрупкие товары требуют дополнительной упаковки (+10-20%)
                  </li>
                )}
                {attributes.requires_special_handling && (
                  <li className="text-info">
                    Особое обращение может увеличить стоимость
                  </li>
                )}
                <li>Точные размеры помогают избежать доплат при доставке</li>
              </ul>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
