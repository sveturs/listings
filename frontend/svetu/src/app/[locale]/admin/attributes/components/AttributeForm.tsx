'use client';

import { useState, useEffect } from 'react';
import { Attribute, adminApi } from '@/services/admin';
import { useTranslations } from 'next-intl';
import { toast } from '@/utils/toast';
import IconPicker from '@/components/IconPicker';
import ValidationRulesEditor from '@/components/attributes/ValidationRulesEditor';
import CustomComponentSelector from '@/components/attributes/CustomComponentSelector';

interface AttributeFormProps {
  attribute?: Attribute | null;
  onSave: (data: Partial<Attribute>) => void;
  onCancel: () => void;
}

interface SelectOption {
  value: string;
  label?: string;
  icon?: string;
  color?: string;
}

export default function AttributeForm({
  attribute,
  onSave,
  onCancel,
}: AttributeFormProps) {
  const t = useTranslations('admin');
  const tCommon = useTranslations('admin');

  const [formData, setFormData] = useState<
    Partial<
      Attribute & {
        variant_type?:
          | 'text'
          | 'number'
          | 'select'
          | 'multiselect'
          | 'boolean'
          | 'date'
          | 'range';
        variant_is_required?: boolean;
        variant_sort_order?: number;
        variant_affects_stock?: boolean;
      }
    >
  >({
    name: '',
    display_name: '',
    attribute_type: 'text',
    icon: '',
    is_searchable: true,
    is_filterable: true,
    is_required: false,
    sort_order: 0,
    show_in_card: false,
    show_in_list: false,
    is_variant_compatible: false,
    affects_stock: false,
    variant_type: 'multiselect',
    variant_is_required: false,
    variant_sort_order: 0,
    variant_affects_stock: false,
    unit: '',
    min_value: undefined,
    max_value: undefined,
    min_length: undefined,
    max_length: undefined,
    pattern: '',
    default_value: '',
    validation_rules: {},
    custom_component: '',
  });

  const [options, setOptions] = useState<SelectOption[]>([]);
  const [isTranslating, setIsTranslating] = useState(false);
  const [translations, setTranslations] = useState<{
    display_name?: Record<string, string>;
    options?: Record<string, Record<string, string>>;
  }>({});
  const [isEditingTranslations, setIsEditingTranslations] = useState(false);

  useEffect(() => {
    if (attribute) {
      setFormData({
        name: attribute.name || '',
        display_name: attribute.display_name || '',
        attribute_type: attribute.attribute_type || 'text',
        icon: attribute.icon || '',
        is_searchable: attribute.is_searchable !== false,
        is_filterable: attribute.is_filterable !== false,
        is_required: attribute.is_required || false,
        sort_order: attribute.sort_order || 0,
        show_in_card: attribute.show_in_card || false,
        show_in_list: attribute.show_in_list || false,
        is_variant_compatible: attribute.is_variant_compatible || false,
        affects_stock: attribute.affects_stock || false,
        variant_type: (attribute as any).variant_type || 'multiselect',
        variant_is_required: (attribute as any).variant_is_required || false,
        variant_sort_order: (attribute as any).variant_sort_order || 0,
        variant_affects_stock:
          (attribute as any).variant_affects_stock || false,
        unit: attribute.unit || '',
        min_value: attribute.min_value,
        max_value: attribute.max_value,
        min_length: attribute.min_length,
        max_length: attribute.max_length,
        pattern: attribute.pattern || '',
        default_value: attribute.default_value || '',
        validation_rules: attribute.validation_rules || {},
        custom_component: attribute.custom_component || '',
      });

      // Parse options if select or multiselect type
      if (
        (attribute.attribute_type === 'select' ||
          attribute.attribute_type === 'multiselect') &&
        attribute.options
      ) {
        try {
          const parsedOptions =
            typeof attribute.options === 'string'
              ? JSON.parse(attribute.options)
              : attribute.options;

          if (Array.isArray(parsedOptions)) {
            setOptions(
              parsedOptions.map((opt) =>
                typeof opt === 'string' ? { value: opt } : opt
              )
            );
          }
        } catch (e) {
          console.error('Failed to parse options:', e);
        }
      }

      // Загружаем переводы если они есть
      if (attribute.translations || attribute.option_translations) {
        setTranslations({
          display_name: attribute.translations || {},
          options: attribute.option_translations || {},
        });
      }
    }
  }, [attribute]);

  const handleChange = (
    e: React.ChangeEvent<
      HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement
    >
  ) => {
    const { name, value, type } = e.target;

    setFormData((prev) => ({
      ...prev,
      [name]:
        type === 'checkbox'
          ? (e.target as HTMLInputElement).checked
          : type === 'number'
            ? value
              ? Number(value)
              : undefined
            : value,
    }));

    // Auto-generate system name from display name
    if (name === 'display_name' && !attribute) {
      const systemName = value
        .toLowerCase()
        .replace(/[^a-z0-9_]/g, '_')
        .replace(/_+/g, '_')
        .replace(/^_|_$/g, '');
      setFormData((prev) => ({ ...prev, name: systemName }));
    }
  };

  const handleAddOption = () => {
    setOptions([...options, { value: '', label: '' }]);
  };

  const handleRemoveOption = (index: number) => {
    setOptions(options.filter((_, i) => i !== index));
  };

  const handleOptionChange = (
    index: number,
    field: keyof SelectOption,
    value: string
  ) => {
    const newOptions = [...options];
    newOptions[index] = { ...newOptions[index], [field]: value };
    setOptions(newOptions);
  };

  const handleTranslationChange = (
    type: 'display_name' | 'options',
    lang: string,
    value: string,
    optionKey?: string
  ) => {
    setTranslations((prev) => {
      if (type === 'display_name') {
        return {
          ...prev,
          display_name: {
            ...prev.display_name,
            [lang]: value,
          },
        };
      } else if (type === 'options' && optionKey) {
        return {
          ...prev,
          options: {
            ...prev.options,
            [optionKey]: {
              ...(prev.options?.[optionKey] || {}),
              [lang]: value,
            },
          },
        };
      }
      return prev;
    });
  };

  const handleTranslate = async () => {
    if (!formData.display_name) {
      toast.error('Введите отображаемое имя для перевода');
      return;
    }

    // Если атрибут уже существует, используем API для перевода атрибута
    if (attribute?.id) {
      setIsTranslating(true);
      try {
        const result = await adminApi.translateAttribute(attribute.id);

        if (result.errors && result.errors.length > 0) {
          toast.error(
            `Переводы получены с ошибками: ${result.errors.join(', ')}`
          );
        } else {
          toast.success('Переводы успешно получены и сохранены');
          // Обновляем состояние с переводами
          if (result.translations) {
            setTranslations(result.translations);
          }
        }
      } catch (error) {
        toast.error('Ошибка при переводе атрибута');
        console.error('Translation error:', error);
      } finally {
        setIsTranslating(false);
      }
    } else {
      // Для нового атрибута переводим только текст
      setIsTranslating(true);
      try {
        const translations = await adminApi.translate(formData.display_name);

        // Сохраняем переводы в formData для последующего сохранения
        setFormData((prev) => ({
          ...prev,
          translations: translations,
        }));

        // Обновляем состояние переводов для отображения
        setTranslations((prev) => ({
          ...prev,
          display_name: translations,
        }));

        // Translate options if select or multiselect type
        if (
          (formData.attribute_type === 'select' ||
            formData.attribute_type === 'multiselect') &&
          options.length > 0
        ) {
          const optionTranslations: Record<string, Record<string, string>> = {};

          for (const option of options) {
            if (option.value) {
              const optTranslations = await adminApi.translate(option.value);
              optionTranslations[option.value] = optTranslations;
            }
          }

          setFormData((prev) => ({
            ...prev,
            option_translations: optionTranslations,
          }));

          // Обновляем состояние переводов для отображения
          setTranslations((prev) => ({
            ...prev,
            options: optionTranslations,
          }));
        }

        toast.success(
          'Переводы получены. Сохраните атрибут для применения переводов.'
        );
      } catch (error) {
        toast.error('Ошибка при переводе');
        console.error('Translation error:', error);
      } finally {
        setIsTranslating(false);
      }
    }
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    if (!formData.name || !formData.display_name || !formData.attribute_type) {
      toast.error('Заполните обязательные поля');
      return;
    }

    // Prepare data
    const dataToSave = { ...formData };

    // Add options for select/multiselect type
    if (
      formData.attribute_type === 'select' ||
      formData.attribute_type === 'multiselect'
    ) {
      dataToSave.options = options
        .filter((opt) => opt.value)
        .map((opt) => opt.value);
    }

    // Добавляем переводы, если они есть
    if (
      translations.display_name &&
      Object.keys(translations.display_name).length > 0
    ) {
      dataToSave.translations = translations.display_name;
    }

    if (translations.options && Object.keys(translations.options).length > 0) {
      dataToSave.option_translations = translations.options;
    }

    onSave(dataToSave);
  };

  const showValidationFields = ['text', 'number', 'range'].includes(
    formData.attribute_type || ''
  );
  const showNumberFields = ['number', 'range'].includes(
    formData.attribute_type || ''
  );
  const showSelectOptions =
    formData.attribute_type === 'select' ||
    formData.attribute_type === 'multiselect';
  const showRangeFields = formData.attribute_type === 'range';

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div className="form-control">
        <label className="label">
          <span className="label-text">{t('attributes.displayName')} *</span>
        </label>
        <input
          type="text"
          name="display_name"
          value={formData.display_name}
          onChange={handleChange}
          className="input input-bordered"
          required
        />
      </div>

      <div className="form-control">
        <label className="label">
          <span className="label-text">{t('attributes.systemName')} *</span>
        </label>
        <input
          type="text"
          name="name"
          value={formData.name}
          onChange={handleChange}
          className="input input-bordered"
          pattern="[a-z0-9_]+"
          required
        />
      </div>

      <div className="form-control">
        <label className="label">
          <span className="label-text">{t('attributes.type')} *</span>
        </label>
        <select
          name="attribute_type"
          value={formData.attribute_type}
          onChange={handleChange}
          className="select select-bordered"
          disabled={!!attribute}
        >
          <option value="text">{t('types.text')}</option>
          <option value="number">{t('types.number')}</option>
          <option value="select">{t('types.select')}</option>
          <option value="multiselect">{t('types.multiselect')}</option>
          <option value="boolean">{t('types.boolean')}</option>
          <option value="date">{t('types.date')}</option>
          <option value="range">{t('types.range')}</option>
          <option value="location">{t('types.location')}</option>
          <option value="file">{t('types.file')}</option>
          <option value="gallery">{t('types.gallery')}</option>
        </select>
      </div>

      <div className="form-control">
        <label className="label">
          <span className="label-text">{t('attributes.icon')}</span>
        </label>
        <IconPicker
          value={formData.icon || ''}
          onChange={(icon) => setFormData((prev) => ({ ...prev, icon }))}
          placeholder="Выберите иконку"
        />
      </div>

      {/* Select Options */}
      {showSelectOptions && (
        <div className="form-control">
          <label className="label">
            <span className="label-text">{t('attributes.options')}</span>
          </label>
          <div className="space-y-2">
            {options.map((option, index) => (
              <div key={index} className="flex gap-2">
                <input
                  type="text"
                  value={option.value}
                  onChange={(e) =>
                    handleOptionChange(index, 'value', e.target.value)
                  }
                  className="input input-bordered input-sm flex-1"
                  placeholder="Значение"
                />
                <button
                  type="button"
                  onClick={() => handleRemoveOption(index)}
                  className="btn btn-ghost btn-sm"
                >
                  ✕
                </button>
              </div>
            ))}
            <button
              type="button"
              onClick={handleAddOption}
              className="btn btn-outline btn-sm"
            >
              {t('attributes.addOption')}
            </button>
          </div>
        </div>
      )}

      {/* Unit field for numbers and range */}
      {showNumberFields && (
        <div className="form-control">
          <label className="label">
            <span className="label-text">{t('attributes.unit')}</span>
          </label>
          <input
            type="text"
            name="unit"
            value={formData.unit}
            onChange={handleChange}
            className="input input-bordered"
            placeholder="км, м², кг..."
          />
        </div>
      )}

      {/* Range specific fields */}
      {showRangeFields && (
        <>
          <div className="divider">{t('attributes.rangeSettings')}</div>
          <div className="alert alert-info">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              className="stroke-current shrink-0 w-6 h-6"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth="2"
                d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
              ></path>
            </svg>
            <span>{t('attributes.rangeDescription')}</span>
          </div>
        </>
      )}

      {/* Validation Rules Editor */}
      {showValidationFields && (
        <>
          <div className="divider">{t('attributes.validation')}</div>
          <ValidationRulesEditor
            value={formData.validation_rules as any}
            onChange={(rules) =>
              setFormData((prev) => ({ ...prev, validation_rules: rules }))
            }
            attributeType={formData.attribute_type}
          />
        </>
      )}

      {/* Custom Component Selector */}
      <CustomComponentSelector
        value={formData.custom_component}
        onChange={(component) =>
          setFormData((prev) => ({ ...prev, custom_component: component }))
        }
        attributeType={formData.attribute_type}
      />

      <div className="form-control">
        <label className="label">
          <span className="label-text">{t('attributes.defaultValue')}</span>
        </label>
        <input
          type="text"
          name="default_value"
          value={String(formData.default_value || '')}
          onChange={handleChange}
          className="input input-bordered"
        />
      </div>

      <div className="divider">Настройки отображения</div>

      <div className="space-y-2">
        <div className="form-control">
          <label className="label cursor-pointer">
            <span className="label-text">{t('attributes.isRequired')}</span>
            <input
              type="checkbox"
              name="is_required"
              checked={formData.is_required}
              onChange={handleChange}
              className="checkbox checkbox-primary"
            />
          </label>
        </div>

        <div className="form-control">
          <label className="label cursor-pointer">
            <span className="label-text">{t('attributes.isSearchable')}</span>
            <input
              type="checkbox"
              name="is_searchable"
              checked={formData.is_searchable}
              onChange={handleChange}
              className="checkbox checkbox-primary"
            />
          </label>
        </div>

        <div className="form-control">
          <label className="label cursor-pointer">
            <span className="label-text">{t('attributes.isFilterable')}</span>
            <input
              type="checkbox"
              name="is_filterable"
              checked={formData.is_filterable}
              onChange={handleChange}
              className="checkbox checkbox-primary"
            />
          </label>
        </div>

        <div className="form-control">
          <label className="label cursor-pointer">
            <span className="label-text">{t('attributes.showInCard')}</span>
            <input
              type="checkbox"
              name="show_in_card"
              checked={formData.show_in_card}
              onChange={handleChange}
              className="checkbox checkbox-primary"
            />
          </label>
        </div>

        <div className="form-control">
          <label className="label cursor-pointer">
            <span className="label-text">{t('attributes.showInList')}</span>
            <input
              type="checkbox"
              name="show_in_list"
              checked={formData.show_in_list}
              onChange={handleChange}
              className="checkbox checkbox-primary"
            />
          </label>
        </div>

        <div className="form-control">
          <label className="label cursor-pointer">
            <span className="label-text flex items-center gap-2">
              🔄 {t('attributes.isVariantCompatible')}
            </span>
            <input
              type="checkbox"
              name="is_variant_compatible"
              checked={formData.is_variant_compatible}
              onChange={handleChange}
              className="checkbox checkbox-secondary"
            />
          </label>
          <label className="label">
            <span className="label-text-alt">
              Позволяет использовать атрибут для создания вариантов товаров
            </span>
          </label>
        </div>

        {formData.is_variant_compatible && (
          <div className="form-control">
            <label className="label cursor-pointer">
              <span className="label-text flex items-center gap-2">
                📦 {t('attributes.affectsStock')}
              </span>
              <input
                type="checkbox"
                name="affects_stock"
                checked={formData.affects_stock}
                onChange={handleChange}
                className="checkbox checkbox-warning"
              />
            </label>
            <label className="label">
              <span className="label-text-alt">
                Влияет ли данный атрибут на учет остатков товара (например,
                размер или цвет могут влиять на остатки)
              </span>
            </label>
          </div>
        )}
      </div>

      {/* Секция настроек вариативного атрибута */}
      {formData.is_variant_compatible && (
        <>
          <div className="divider flex items-center gap-2">
            <span className="text-lg">🔄</span>
            <span>Настройки вариативного атрибута</span>
          </div>

          <div className="alert alert-info">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              className="stroke-current shrink-0 w-6 h-6"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth="2"
                d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
              ></path>
            </svg>
            <div>
              <h3 className="font-bold">Вариативный атрибут</h3>
              <div className="text-xs">
                Эти настройки определяют, как атрибут будет использоваться при
                создании вариантов товаров
              </div>
            </div>
          </div>

          <div className="form-control">
            <label className="label">
              <span className="label-text">Тип вариативного атрибута</span>
            </label>
            <select
              name="variant_type"
              value={formData.variant_type || 'multiselect'}
              onChange={handleChange}
              className="select select-bordered"
            >
              <option value="text">Текст</option>
              <option value="number">Число</option>
              <option value="select">Выбор (один)</option>
              <option value="multiselect">Множественный выбор</option>
              <option value="boolean">Да/Нет</option>
              <option value="date">Дата</option>
              <option value="range">Диапазон</option>
            </select>
            <label className="label">
              <span className="label-text-alt">
                Для большинства случаев рекомендуется &ldquo;Множественный
                выбор&rdquo;
              </span>
            </label>
          </div>

          <div className="form-control">
            <label className="label">
              <span className="label-text">Порядок сортировки</span>
            </label>
            <input
              type="number"
              name="variant_sort_order"
              value={formData.variant_sort_order || 0}
              onChange={handleChange}
              className="input input-bordered"
              min="0"
              placeholder="0"
            />
            <label className="label">
              <span className="label-text-alt">
                Чем меньше число, тем выше в списке будет атрибут
              </span>
            </label>
          </div>

          <div className="space-y-2">
            <div className="form-control">
              <label className="label cursor-pointer">
                <span className="label-text">Обязательный для вариантов</span>
                <input
                  type="checkbox"
                  name="variant_is_required"
                  checked={formData.variant_is_required}
                  onChange={handleChange}
                  className="checkbox checkbox-primary"
                />
              </label>
              <label className="label">
                <span className="label-text-alt">
                  Если включено, пользователь обязан выбрать значение для
                  создания вариантов
                </span>
              </label>
            </div>

            <div className="form-control">
              <label className="label cursor-pointer">
                <span className="label-text flex items-center gap-2">
                  📦 <span>Влияет на остатки товара</span>
                </span>
                <input
                  type="checkbox"
                  name="variant_affects_stock"
                  checked={formData.variant_affects_stock}
                  onChange={handleChange}
                  className="checkbox checkbox-warning"
                />
              </label>
              <label className="label">
                <span className="label-text-alt">
                  Если включено, каждый вариант будет иметь отдельный учёт
                  остатков
                </span>
              </label>
            </div>
          </div>
        </>
      )}

      <div className="flex gap-2 items-center">
        <button
          type="button"
          onClick={handleTranslate}
          className="btn btn-secondary"
          disabled={isTranslating || !formData.display_name}
        >
          {isTranslating ? (
            <>
              <span className="loading loading-spinner loading-sm"></span>
              {t('attributes.translating')}
            </>
          ) : (
            <>🌍 {tCommon('common.translate')}</>
          )}
        </button>
        {translations.display_name &&
          Object.keys(translations.display_name).length > 0 && (
            <div className="flex items-center gap-2 text-sm text-success">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                className="h-4 w-4"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
                />
              </svg>
              <span>Переводы загружены</span>
            </div>
          )}
      </div>

      {/* Секция отображения переводов */}
      {translations.display_name &&
        Object.keys(translations.display_name).length > 0 && (
          <div className="mt-6 animate-in slide-in-from-bottom duration-300">
            <div className="divider">{t('attributes.translations')}</div>

            {/* Переводы названия атрибута */}
            <div className="card bg-base-100 shadow-sm border border-base-300 transition-all duration-300 hover:shadow-md">
              <div className="card-body">
                <div className="flex items-center justify-between mb-3 flex-wrap gap-2">
                  <h4 className="card-title text-base">
                    {t('attributes.displayNameTranslations')}
                  </h4>
                  <div className="flex items-center gap-2">
                    <button
                      type="button"
                      onClick={() =>
                        setIsEditingTranslations(!isEditingTranslations)
                      }
                      className="btn btn-ghost btn-xs"
                    >
                      {isEditingTranslations ? '✓ Готово' : '✏️ Редактировать'}
                    </button>
                    {isEditingTranslations && (
                      <button
                        type="button"
                        onClick={() => {
                          if (confirm('Очистить все переводы?')) {
                            setTranslations({});
                            setIsEditingTranslations(false);
                          }
                        }}
                        className="btn btn-ghost btn-xs text-error"
                      >
                        🗑️ Очистить
                      </button>
                    )}
                  </div>
                </div>
                <div className="space-y-3">
                  {Object.entries(translations.display_name).map(
                    ([lang, text]) => (
                      <div key={lang} className="flex items-center gap-3">
                        <div className="badge badge-primary badge-outline font-semibold">
                          {lang.toUpperCase()}
                        </div>
                        <div className="flex-1">
                          <input
                            type="text"
                            value={text}
                            onChange={(e) =>
                              handleTranslationChange(
                                'display_name',
                                lang,
                                e.target.value
                              )
                            }
                            readOnly={!isEditingTranslations}
                            className={`input input-bordered input-sm w-full ${
                              isEditingTranslations ? '' : 'bg-base-200'
                            }`}
                          />
                        </div>
                      </div>
                    )
                  )}
                </div>
              </div>
            </div>

            {/* Переводы опций */}
            {translations.options &&
              Object.keys(translations.options).length > 0 && (
                <div className="card bg-base-100 shadow-sm border border-base-300 mt-4 transition-all duration-300 hover:shadow-md">
                  <div className="card-body">
                    <h4 className="card-title text-base">
                      {t('attributes.optionTranslations')}
                    </h4>
                    <div className="space-y-4">
                      {Object.entries(translations.options).map(
                        ([optionValue, langTranslations]) => (
                          <div
                            key={optionValue}
                            className="border border-base-300 rounded-lg p-3 hover:shadow-sm transition-shadow"
                          >
                            <div className="font-medium text-sm mb-2 flex items-center gap-2">
                              <span className="text-base-content/70">
                                Опция:
                              </span>
                              <span className="font-mono bg-base-200 px-2 py-0.5 rounded">
                                {optionValue}
                              </span>
                            </div>
                            <div className="space-y-2 pl-2">
                              {Object.entries(langTranslations).map(
                                ([lang, translation]) => (
                                  <div
                                    key={lang}
                                    className="flex items-center gap-3"
                                  >
                                    <div className="badge badge-sm badge-ghost font-semibold min-w-[3rem] justify-center">
                                      {lang.toUpperCase()}
                                    </div>
                                    {isEditingTranslations ? (
                                      <input
                                        type="text"
                                        value={translation}
                                        onChange={(e) =>
                                          handleTranslationChange(
                                            'options',
                                            lang,
                                            e.target.value,
                                            optionValue
                                          )
                                        }
                                        className="input input-bordered input-xs flex-1"
                                      />
                                    ) : (
                                      <div className="text-sm flex-1 px-2">
                                        {translation}
                                      </div>
                                    )}
                                  </div>
                                )
                              )}
                            </div>
                          </div>
                        )
                      )}
                    </div>
                  </div>
                </div>
              )}
          </div>
        )}

      <div className="flex gap-2 pt-4">
        <button type="submit" className="btn btn-primary">
          {tCommon('common.save')}
        </button>
        <button type="button" onClick={onCancel} className="btn btn-ghost">
          {tCommon('common.cancel')}
        </button>
      </div>
    </form>
  );
}
