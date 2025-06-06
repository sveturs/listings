'use client';

import { useState, useEffect } from 'react';
import { Attribute, adminApi } from '@/services/admin';
import { useTranslations } from 'next-intl';
import { toast } from '@/utils/toast';
import IconPicker from '@/components/IconPicker';

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
  const t = useTranslations('admin.attributes');
  const tCommon = useTranslations('admin.common');

  const [formData, setFormData] = useState<Partial<Attribute>>({
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
    unit: '',
    min_value: undefined,
    max_value: undefined,
    min_length: undefined,
    max_length: undefined,
    pattern: '',
    default_value: '',
  });

  const [options, setOptions] = useState<SelectOption[]>([]);
  const [isTranslating, setIsTranslating] = useState(false);

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
        unit: attribute.unit || '',
        min_value: attribute.min_value,
        max_value: attribute.max_value,
        min_length: attribute.min_length,
        max_length: attribute.max_length,
        pattern: attribute.pattern || '',
        default_value: attribute.default_value || '',
      });

      // Parse options if select type
      if (attribute.attribute_type === 'select' && attribute.options) {
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

  const handleTranslate = async () => {
    if (!formData.display_name) {
      toast.error('Введите отображаемое имя для перевода');
      return;
    }

    setIsTranslating(true);
    try {
      const _translations = await adminApi.translate(formData.display_name);

      // Translate options if select type
      if (formData.attribute_type === 'select' && options.length > 0) {
        for (const option of options) {
          if (option.value) {
            const _optionTranslations = await adminApi.translate(option.value);
            // Here we would save option translations
          }
        }
      }

      toast.success('Переводы получены');
    } catch (error) {
      toast.error('Ошибка при переводе');
      console.error('Translation error:', error);
    } finally {
      setIsTranslating(false);
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

    // Add options for select type
    if (formData.attribute_type === 'select') {
      dataToSave.options = options
        .filter((opt) => opt.value)
        .map((opt) => opt.value);
    }

    onSave(dataToSave);
  };

  const showValidationFields = ['text', 'number', 'range'].includes(
    formData.attribute_type || ''
  );
  const showNumberFields = ['number', 'range'].includes(
    formData.attribute_type || ''
  );
  const showTextFields = formData.attribute_type === 'text';
  const showSelectOptions = formData.attribute_type === 'select';

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div className="form-control">
        <label className="label">
          <span className="label-text">{t('displayName')} *</span>
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
          <span className="label-text">{t('systemName')} *</span>
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
          <span className="label-text">{t('type')} *</span>
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
          <span className="label-text">{t('icon')}</span>
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
            <span className="label-text">{t('options')}</span>
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
              {t('addOption')}
            </button>
          </div>
        </div>
      )}

      {/* Unit field for numbers */}
      {showNumberFields && (
        <div className="form-control">
          <label className="label">
            <span className="label-text">{t('unit')}</span>
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

      {/* Validation Fields */}
      {showValidationFields && <div className="divider">{t('validation')}</div>}

      {showNumberFields && (
        <>
          <div className="grid grid-cols-2 gap-2">
            <div className="form-control">
              <label className="label">
                <span className="label-text">{t('minValue')}</span>
              </label>
              <input
                type="number"
                name="min_value"
                value={formData.min_value || ''}
                onChange={handleChange}
                className="input input-bordered"
              />
            </div>
            <div className="form-control">
              <label className="label">
                <span className="label-text">{t('maxValue')}</span>
              </label>
              <input
                type="number"
                name="max_value"
                value={formData.max_value || ''}
                onChange={handleChange}
                className="input input-bordered"
              />
            </div>
          </div>
        </>
      )}

      {showTextFields && (
        <>
          <div className="grid grid-cols-2 gap-2">
            <div className="form-control">
              <label className="label">
                <span className="label-text">{t('minLength')}</span>
              </label>
              <input
                type="number"
                name="min_length"
                value={formData.min_length || ''}
                onChange={handleChange}
                className="input input-bordered"
              />
            </div>
            <div className="form-control">
              <label className="label">
                <span className="label-text">{t('maxLength')}</span>
              </label>
              <input
                type="number"
                name="max_length"
                value={formData.max_length || ''}
                onChange={handleChange}
                className="input input-bordered"
              />
            </div>
          </div>
          <div className="form-control">
            <label className="label">
              <span className="label-text">{t('pattern')}</span>
            </label>
            <input
              type="text"
              name="pattern"
              value={formData.pattern}
              onChange={handleChange}
              className="input input-bordered"
              placeholder="^[A-Z][0-9]+$"
            />
          </div>
        </>
      )}

      <div className="form-control">
        <label className="label">
          <span className="label-text">{t('defaultValue')}</span>
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
            <span className="label-text">{t('isRequired')}</span>
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
            <span className="label-text">{t('isSearchable')}</span>
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
            <span className="label-text">{t('isFilterable')}</span>
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
            <span className="label-text">{t('showInCard')}</span>
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
            <span className="label-text">{t('showInList')}</span>
            <input
              type="checkbox"
              name="show_in_list"
              checked={formData.show_in_list}
              onChange={handleChange}
              className="checkbox checkbox-primary"
            />
          </label>
        </div>
      </div>

      <div className="flex gap-2">
        <button
          type="button"
          onClick={handleTranslate}
          className="btn btn-secondary"
          disabled={isTranslating || !formData.display_name}
        >
          {isTranslating ? (
            <>
              <span className="loading loading-spinner loading-sm"></span>
              {t('translating')}
            </>
          ) : (
            <>🌍 {tCommon('translate')}</>
          )}
        </button>
      </div>

      <div className="flex gap-2 pt-4">
        <button type="submit" className="btn btn-primary flex-1">
          {tCommon('save')}
        </button>
        <button
          type="button"
          onClick={onCancel}
          className="btn btn-ghost flex-1"
        >
          {tCommon('cancel')}
        </button>
      </div>
    </form>
  );
}
