'use client';

import React, { useState, useRef } from 'react';
import { useTranslations } from 'next-intl';
import {
  Calendar,
  User,
  Gauge,
  AlertCircle,
  Info,
  Fuel,
  Settings,
  DollarSign,
  Globe,
  Shield,
  Package,
} from 'lucide-react';

interface FloatingInputProps {
  label: string;
  value: string;
  onChange: (value: string) => void;
  type?: string;
  icon?: React.ReactNode;
  required?: boolean;
  manualOnly?: boolean;
  aiValue?: string;
  helpText?: string;
  options?: Array<{ value: string; label: string }>;
  placeholder?: string;
  unit?: string;
  step?: string;
}

const FloatingInput: React.FC<FloatingInputProps> = ({
  label,
  value,
  onChange,
  type = 'text',
  icon,
  required = false,
  manualOnly = false,
  aiValue,
  helpText,
  options,
  placeholder,
  unit,
  step,
}) => {
  const [isFocused, setIsFocused] = useState(false);
  const [showTooltip, setShowTooltip] = useState(false);
  const inputRef = useRef<HTMLInputElement | HTMLSelectElement>(null);
  const hasValue = value && value.length > 0;
  const shouldFloatLabel = isFocused || hasValue;

  // Определяем цвет рамки и фона
  const getBorderClass = () => {
    if (manualOnly && !hasValue && required) {
      return 'border-error border-2 bg-error/5';
    }
    if (isFocused) {
      return 'border-primary border-2';
    }
    if (aiValue && !value) {
      return 'border-success border-dashed bg-success/5';
    }
    return 'border-base-300';
  };

  return (
    <div className="relative group">
      <div
        className={`relative transition-all duration-200 ${
          manualOnly && !hasValue ? 'transform hover:scale-[1.01]' : ''
        }`}
        onMouseEnter={() => setShowTooltip(true)}
        onMouseLeave={() => setShowTooltip(false)}
      >
        {/* Иконка слева */}
        {icon && (
          <div
            className={`absolute left-3 top-1/2 -translate-y-1/2 z-10 transition-colors ${
              manualOnly && !hasValue
                ? 'text-error'
                : aiValue && !value
                  ? 'text-success'
                  : 'text-base-content/50'
            }`}
          >
            {icon}
          </div>
        )}

        {/* Индикатор обязательного ручного ввода */}
        {manualOnly && !hasValue && required && (
          <div className="absolute -top-2 -right-2 z-20">
            <div className="relative">
              <span className="absolute inset-0 animate-ping bg-error rounded-full opacity-75"></span>
              <AlertCircle className="w-5 h-5 text-error" />
            </div>
          </div>
        )}

        {/* Индикатор AI заполнения */}
        {aiValue && !value && (
          <div className="absolute -top-2 -left-2 z-20">
            <div className="badge badge-success badge-xs">AI</div>
          </div>
        )}

        {/* Input или Select */}
        {options ? (
          <select
            ref={inputRef as React.RefObject<HTMLSelectElement>}
            value={value || aiValue || ''}
            onChange={(e) => onChange(e.target.value)}
            onFocus={() => setIsFocused(true)}
            onBlur={() => setIsFocused(false)}
            className={`select select-bordered w-full ${icon ? 'pl-10' : 'pl-4'} ${
              unit ? 'pr-16' : 'pr-4'
            } pt-6 pb-2 ${getBorderClass()}`}
          >
            <option value="">
              {placeholder || `Выберите ${label.toLowerCase()}`}
            </option>
            {options.map((opt, idx) => (
              <option key={`${opt.value}-${idx}`} value={opt.value}>
                {opt.label}
              </option>
            ))}
          </select>
        ) : (
          <input
            ref={inputRef as React.RefObject<HTMLInputElement>}
            type={type}
            value={value || aiValue || ''}
            onChange={(e) => onChange(e.target.value)}
            onFocus={() => setIsFocused(true)}
            onBlur={() => setIsFocused(false)}
            placeholder={isFocused && !value ? placeholder : ''}
            step={step}
            className={`input input-bordered w-full ${icon ? 'pl-10' : 'pl-4'} ${
              unit ? 'pr-16' : 'pr-4'
            } pt-6 pb-2 ${getBorderClass()}`}
          />
        )}

        {/* Плавающий лейбл */}
        <label
          className={`absolute transition-all duration-200 pointer-events-none
            ${icon ? 'left-10' : 'left-4'}
            ${
              shouldFloatLabel
                ? 'top-1 text-xs ' +
                  (manualOnly && !hasValue
                    ? 'text-error font-semibold'
                    : aiValue && !value
                      ? 'text-success font-semibold'
                      : 'text-base-content/70')
                : 'top-1/2 -translate-y-1/2 text-base text-base-content/50'
            }
          `}
        >
          {label}
          {required && <span className="text-error ml-1">*</span>}
        </label>

        {/* Единица измерения */}
        {unit && (
          <span className="absolute right-4 top-1/2 -translate-y-1/2 text-base-content/50">
            {unit}
          </span>
        )}
      </div>

      {/* Tooltip с подсказкой */}
      {helpText && showTooltip && (
        <div className="absolute bottom-full left-0 mb-2 z-30">
          <div className="bg-base-100 border border-base-300 rounded-lg shadow-xl p-3 max-w-xs">
            <div className="flex items-start gap-2">
              <Info className="w-4 h-4 text-info flex-shrink-0 mt-0.5" />
              <p className="text-sm">{helpText}</p>
            </div>
          </div>
        </div>
      )}

      {/* Информация о статусе поля */}
      {manualOnly && !hasValue && required && (
        <p className="text-error text-xs mt-1 flex items-center gap-1">
          <AlertCircle className="w-3 h-3" />
          Требуется ручной ввод
        </p>
      )}
      {aiValue && !value && (
        <p className="text-success text-xs mt-1 flex items-center gap-1">
          <Info className="w-3 h-3" />
          AI предложил: {aiValue}
        </p>
      )}
    </div>
  );
};

interface CarAttributesFormProps {
  attributes: Record<string, any>;
  onChange: (attributes: Record<string, any>) => void;
  aiSuggestions?: Record<string, any>;
  categoryAttributes?: any[]; // Атрибуты категории из БД
}

export const CarAttributesForm: React.FC<CarAttributesFormProps> = ({
  attributes,
  onChange,
  aiSuggestions = {},
  categoryAttributes = [],
}) => {
  const t = useTranslations('cars');

  const updateAttribute = (key: string, value: string) => {
    onChange({
      ...attributes,
      [key]: value,
    });
  };

  // Список атрибутов, которые НЕВОЗМОЖНО определить по фото
  const MANUAL_ONLY_ATTRIBUTES = [
    'year',
    'mileage',
    'transmission',
    'fuel_type',
    'engine_size',
    'power',
    'owner_count',
    'vin',
    'registration_date',
    'in_serbia_since',
    'imported_from',
    'import_date',
    'accident_history',
    'service_history',
    'price',
    'technical_inspection',
    'chassis_number',
    'engine_number',
  ];

  // Функция для определения конфигурации поля
  const getFieldConfig = (attrName: string) => {
    const categoryAttr = categoryAttributes.find(
      (attr) => attr.name === attrName
    );

    // Пытаемся получить перевод для атрибута
    let translatedLabel = '';
    const translationKey = `attributes.${attrName}`;
    try {
      // Проверяем, существует ли ключ перевода
      const translated = t(translationKey as any);
      // Если вернулось то же самое, что и ключ - значит перевода нет
      translatedLabel =
        translated === translationKey
          ? categoryAttr?.display_name || attrName
          : translated;
    } catch {
      // Если перевод не найден, используем display_name из БД или имя атрибута
      translatedLabel = categoryAttr?.display_name || attrName;
    }

    const config: any = {
      label: translatedLabel,
      helpText: categoryAttr?.description,
      isRequired: categoryAttr?.is_required || false,
      placeholder: categoryAttr?.placeholder,
      type: 'text',
      icon: <Settings className="w-4 h-4" />,
      options: undefined,
      unit: '',
      step: undefined,
    };

    // Настройка специфичных полей
    switch (attrName) {
      case 'year':
        config.type = 'number';
        config.icon = <Calendar className="w-4 h-4" />;
        config.placeholder = '2020';
        break;
      case 'mileage':
        config.type = 'number';
        config.unit = 'км';
        config.icon = <Gauge className="w-4 h-4" />;
        config.placeholder = '50000';
        break;
      case 'engine_size':
        config.type = 'number';
        config.step = '0.1';
        config.unit = 'л';
        config.placeholder = '2.0';
        break;
      case 'power':
        config.type = 'number';
        config.unit = 'л.с.';
        config.icon = <Gauge className="w-4 h-4" />;
        config.placeholder = '150';
        break;
      case 'price':
        config.type = 'number';
        config.unit = '€';
        config.icon = <DollarSign className="w-4 h-4" />;
        config.placeholder = '15000';
        break;
      case 'registration_date':
      case 'in_serbia_since':
      case 'technical_inspection':
        config.type = 'date';
        config.icon = <Calendar className="w-4 h-4" />;
        break;
      case 'fuel_type':
        config.icon = <Fuel className="w-4 h-4" />;
        config.options = [
          { value: 'petrol', label: t('options.fuelType.petrol') },
          { value: 'diesel', label: t('options.fuelType.diesel') },
          { value: 'electric', label: t('options.fuelType.electric') },
          { value: 'hybrid', label: t('options.fuelType.hybrid') },
          {
            value: 'plug-in-hybrid',
            label: t('options.fuelType.plugInHybrid'),
          },
          { value: 'lpg', label: t('options.fuelType.lpg') },
          { value: 'cng', label: t('options.fuelType.cng') },
        ];
        break;
      case 'transmission':
        config.options = [
          { value: 'manual', label: t('options.transmission.manual') },
          { value: 'automatic', label: t('options.transmission.automatic') },
          {
            value: 'semi-automatic',
            label: t('options.transmission.semiAutomatic'),
          },
          { value: 'cvt', label: t('options.transmission.cvt') },
          { value: 'dsg', label: t('options.transmission.dsg') },
        ];
        break;
      case 'body_type':
        config.options = [
          { value: 'sedan', label: t('options.bodyType.sedan') },
          { value: 'hatchback', label: t('options.bodyType.hatchback') },
          { value: 'suv', label: t('options.bodyType.suv') },
          { value: 'crossover', label: t('options.bodyType.crossover') },
          { value: 'wagon', label: t('options.bodyType.wagon') },
          { value: 'coupe', label: t('options.bodyType.coupe') },
          { value: 'minivan', label: t('options.bodyType.minivan') },
          { value: 'pickup', label: t('options.bodyType.pickup') },
          { value: 'convertible', label: t('options.bodyType.convertible') },
        ];
        break;
      case 'wheel_drive':
        config.options = [
          { value: 'fwd', label: t('options.wheelDrive.fwd') },
          { value: 'rwd', label: t('options.wheelDrive.rwd') },
          { value: 'awd', label: t('options.wheelDrive.awd') },
          { value: '4wd', label: t('options.wheelDrive.4wd') },
        ];
        break;
      case 'condition':
        config.options = [
          { value: 'new', label: t('options.condition.new') },
          { value: 'excellent', label: t('options.condition.excellent') },
          { value: 'good', label: t('options.condition.good') },
          { value: 'fair', label: t('options.condition.fair') },
          { value: 'damaged', label: t('options.condition.damaged') },
          { value: 'for_parts', label: t('options.condition.forParts') },
        ];
        break;
      case 'owner_count':
        config.type = 'number';
        config.icon = <User className="w-4 h-4" />;
        config.placeholder = '1';
        break;
      case 'imported_from':
        config.icon = <Globe className="w-4 h-4" />;
        config.placeholder = 'Германия';
        break;
      case 'doors':
      case 'seats':
        config.type = 'number';
        config.icon = <Package className="w-4 h-4" />;
        break;
      case 'vin':
      case 'chassis_number':
      case 'engine_number':
        config.icon = <Shield className="w-4 h-4" />;
        break;
      case 'negotiable':
        config.icon = <DollarSign className="w-4 h-4" />;
        config.options = [
          { value: 'yes', label: t('options.negotiable.yes') },
          { value: 'no', label: t('options.negotiable.no') },
          { value: 'slight', label: t('options.negotiable.slight') },
        ];
        break;
      case 'exchange_possible':
        config.options = [
          { value: 'yes', label: t('options.exchange.yes') },
          { value: 'no', label: t('options.exchange.no') },
          { value: 'maybe', label: t('options.exchange.maybe') },
        ];
        break;
    }

    // Используем опции из атрибута категории, если они есть
    if (categoryAttr?.options && categoryAttr.options.length > 0) {
      config.options = categoryAttr.options.map((opt: any) => {
        // Опции могут быть либо строками, либо объектами {value, label}
        if (typeof opt === 'string') {
          return { value: opt, label: opt };
        }
        // Если это объект, используем его поля value и label
        return {
          value: String(opt.value || opt.label || opt),
          label: String(opt.label || opt.value || opt),
        };
      });
    }

    return config;
  };

  // Группируем атрибуты по секциям
  const attributeSections = [
    {
      title: t('form.basicInfo'),
      icon: <Settings className="w-5 h-5" />,
      attributes: [
        'brand',
        'model',
        'year',
        'mileage',
        'body_type',
        'color',
        'condition',
      ],
    },
    {
      title: t('form.technicalSpecs'),
      icon: <Settings className="w-5 h-5" />,
      attributes: [
        'transmission',
        'fuel_type',
        'engine_size',
        'power',
        'wheel_drive',
        'emission_class',
        'euro_standard',
      ],
    },
    {
      title: t('form.equipment'),
      icon: <Package className="w-5 h-5" />,
      attributes: [
        'doors',
        'seats',
        'interior_material',
        'interior_color',
        'equipment',
      ],
    },
    {
      title: t('form.historyAndDocs'),
      icon: <User className="w-5 h-5" />,
      attributes: [
        'owner_count',
        'in_serbia_since',
        'imported_from',
        'technical_inspection',
        'accident_history',
        'service_history',
      ],
    },
    {
      title: t('form.identification'),
      icon: <Shield className="w-5 h-5" />,
      attributes: [
        'vin',
        'chassis_number',
        'engine_number',
        'registration_date',
        'registration_plate',
      ],
    },
    {
      title: t('form.priceAndTerms'),
      icon: <DollarSign className="w-5 h-5" />,
      attributes: [
        'price',
        'negotiable',
        'exchange_possible',
        'credit_available',
      ],
    },
  ];

  // Собираем все уникальные атрибуты
  const allAttributeKeys = new Set<string>();

  // Из категории
  categoryAttributes.forEach((attr) => allAttributeKeys.add(attr.name));
  // Из текущих данных
  Object.keys(attributes).forEach((key) => allAttributeKeys.add(key));
  // Из AI предложений
  Object.keys(aiSuggestions).forEach((key) => allAttributeKeys.add(key));

  // Атрибуты, которые не попали ни в одну секцию
  const usedAttributes = new Set(
    attributeSections.flatMap((s) => s.attributes)
  );
  const additionalAttributes = Array.from(allAttributeKeys).filter(
    (key) => !usedAttributes.has(key)
  );

  return (
    <div className="space-y-8">
      {/* Информационные блоки */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div className="alert alert-info">
          <Info className="w-5 h-5" />
          <div>
            <h3 className="font-bold">{t('form.aiAutoFill')}</h3>
            <p className="text-sm">{t('form.aiAutoFillDescription')}</p>
          </div>
        </div>

        <div className="alert alert-warning">
          <AlertCircle className="w-5 h-5" />
          <div>
            <h3 className="font-bold">{t('form.manualInput')}</h3>
            <p className="text-sm">{t('form.manualInputDescription')}</p>
          </div>
        </div>
      </div>

      {/* Отображаем секции с атрибутами */}
      {attributeSections.map((section) => {
        // Показываем секцию только если есть хотя бы один атрибут из этой секции
        const sectionAttributes = section.attributes.filter((attr) =>
          allAttributeKeys.has(attr)
        );
        if (sectionAttributes.length === 0) return null;

        return (
          <div key={section.title} className="space-y-2">
            <h3 className="text-lg font-semibold flex items-center gap-2">
              {section.icon}
              {section.title}
            </h3>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              {sectionAttributes.map((attrName) => {
                const config = getFieldConfig(attrName);
                const isManualOnly = MANUAL_ONLY_ATTRIBUTES.includes(attrName);
                const value = attributes[attrName] || '';
                const aiValue = aiSuggestions[attrName];

                return (
                  <FloatingInput
                    key={attrName}
                    label={config.label}
                    value={value}
                    onChange={(val) => updateAttribute(attrName, val)}
                    type={config.type}
                    icon={config.icon}
                    required={config.isRequired}
                    manualOnly={isManualOnly && !aiValue}
                    aiValue={aiValue}
                    helpText={config.helpText}
                    options={config.options}
                    unit={config.unit}
                    step={config.step}
                    placeholder={config.placeholder}
                  />
                );
              })}
            </div>
          </div>
        );
      })}

      {/* Дополнительные атрибуты */}
      {additionalAttributes.length > 0 && (
        <div className="space-y-2">
          <h3 className="text-lg font-semibold flex items-center gap-2">
            <Package className="w-5 h-5" />
            Дополнительные характеристики
          </h3>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            {additionalAttributes.map((attrName) => {
              const config = getFieldConfig(attrName);
              const isManualOnly = MANUAL_ONLY_ATTRIBUTES.includes(attrName);
              const value = attributes[attrName] || '';
              const aiValue = aiSuggestions[attrName];

              return (
                <FloatingInput
                  key={attrName}
                  label={config.label}
                  value={value}
                  onChange={(val) => updateAttribute(attrName, val)}
                  type={config.type}
                  icon={config.icon}
                  required={config.isRequired}
                  manualOnly={isManualOnly && !aiValue}
                  aiValue={aiValue}
                  helpText={config.helpText}
                  options={config.options}
                  unit={config.unit}
                  step={config.step}
                  placeholder={config.placeholder}
                />
              );
            })}
          </div>
        </div>
      )}
    </div>
  );
};
