'use client';

import { useState, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import { toast } from '@/utils/toast';
import { VariantAttribute } from '../page';

interface VariantAttributeFormProps {
  attribute?: VariantAttribute | null;
  onSave: (data: Partial<VariantAttribute>) => void;
  onCancel: () => void;
}

export default function VariantAttributeForm({
  attribute,
  onSave,
  onCancel,
}: VariantAttributeFormProps) {
  const t = useTranslations('admin.variantAttributes');
  const tCommon = useTranslations('admin.common');

  const [formData, setFormData] = useState<Partial<VariantAttribute>>({
    name: '',
    display_name: '',
    type: 'text',
    is_required: false,
    sort_order: 0,
    affects_stock: false,
  });

  useEffect(() => {
    if (attribute) {
      setFormData({
        name: attribute.name || '',
        display_name: attribute.display_name || '',
        type: attribute.type || 'text',
        is_required: attribute.is_required || false,
        sort_order: attribute.sort_order || 0,
        affects_stock: attribute.affects_stock || false,
      });
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
            ? value ? Number(value) : 0
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

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    if (!formData.name || !formData.display_name || !formData.type) {
      toast.error(t('validationError'));
      return;
    }

    onSave(formData);
  };

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
          placeholder={t('displayNamePlaceholder')}
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
          placeholder={t('systemNamePlaceholder')}
        />
        <label className="label">
          <span className="label-text-alt">{t('systemNameHint')}</span>
        </label>
      </div>

      <div className="form-control">
        <label className="label">
          <span className="label-text">{t('type')} *</span>
        </label>
        <select
          name="type"
          value={formData.type}
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
        </select>
      </div>

      <div className="form-control">
        <label className="label">
          <span className="label-text">{t('sortOrder')}</span>
        </label>
        <input
          type="number"
          name="sort_order"
          value={formData.sort_order}
          onChange={handleChange}
          className="input input-bordered"
          min="0"
          placeholder="0"
        />
        <label className="label">
          <span className="label-text-alt">{t('sortOrderHint')}</span>
        </label>
      </div>

      <div className="divider">{t('settings')}</div>

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
          <label className="label">
            <span className="label-text-alt">{t('isRequiredHint')}</span>
          </label>
        </div>

        <div className="form-control">
          <label className="label cursor-pointer">
            <span className="label-text flex items-center gap-2">
              📦 {t('affectsStock')}
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
            <span className="label-text-alt">{t('affectsStockHint')}</span>
          </label>
        </div>
      </div>

      <div className="flex gap-2 pt-4">
        <button type="submit" className="btn btn-primary">
          {tCommon('save')}
        </button>
        <button
          type="button"
          onClick={onCancel}
          className="btn btn-ghost"
        >
          {tCommon('cancel')}
        </button>
      </div>
    </form>
  );
}