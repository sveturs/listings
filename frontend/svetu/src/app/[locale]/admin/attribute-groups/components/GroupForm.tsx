'use client';

import { useState, useEffect } from 'react';
import { AttributeGroup, adminApi } from '@/services/admin';
import { useTranslations } from 'next-intl';
import { toast } from '@/utils/toast';

interface GroupFormProps {
  group?: AttributeGroup | null;
  onSave: (data: Partial<AttributeGroup>) => void;
  onCancel: () => void;
}

export default function GroupForm({ group, onSave, onCancel }: GroupFormProps) {
  const t = useTranslations('admin');
  const tCommon = useTranslations('admin');

  const [formData, setFormData] = useState<Partial<AttributeGroup>>({
    name: '',
    display_name: '',
    description: '',
    icon: '',
    sort_order: 0,
    is_active: true,
  });

  const [isTranslating, setIsTranslating] = useState(false);

  useEffect(() => {
    if (group) {
      setFormData({
        name: group.name || '',
        display_name: group.display_name || '',
        description: group.description || '',
        icon: group.icon || '',
        sort_order: group.sort_order || 0,
        is_active: group.is_active !== false,
      });
    }
  }, [group]);

  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>
  ) => {
    const { name, value, type } = e.target;

    setFormData((prev) => ({
      ...prev,
      [name]:
        type === 'checkbox'
          ? (e.target as HTMLInputElement).checked
          : type === 'number'
            ? Number(value)
            : value,
    }));

    // Auto-generate system name from display name
    if (name === 'display_name' && !group) {
      const systemName = value
        .toLowerCase()
        .replace(/[^a-z0-9_]/g, '_')
        .replace(/_+/g, '_')
        .replace(/^_|_$/g, '');
      setFormData((prev) => ({ ...prev, name: systemName }));
    }
  };

  const handleTranslate = async () => {
    if (!formData.display_name) {
      toast.error('Введите название группы для перевода');
      return;
    }

    setIsTranslating(true);
    try {
      const _translations = await adminApi.translate(formData.display_name);

      // Also translate description if present
      if (formData.description) {
        const _descTranslations = await adminApi.translate(
          formData.description
        );
        // Here we would save description translations
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

    if (!formData.name || !formData.display_name) {
      toast.error('Заполните обязательные поля');
      return;
    }

    onSave(formData);
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div className="form-control">
        <label className="label">
          <span className="label-text">{t('attributeGroups.groupName')} *</span>
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
          <span className="label-text">{tCommon('common.systemName')} *</span>
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
          <span className="label-text">{t('attributeGroups.description')}</span>
        </label>
        <textarea
          name="description"
          value={formData.description}
          onChange={handleChange}
          className="textarea textarea-bordered"
          rows={3}
        />
      </div>

      <div className="form-control">
        <label className="label">
          <span className="label-text">{t('attributeGroups.icon')}</span>
        </label>
        <input
          type="text"
          name="icon"
          value={formData.icon}
          onChange={handleChange}
          className="input input-bordered"
          placeholder="🏷️"
        />
      </div>

      <div className="form-control">
        <label className="label">
          <span className="label-text">{tCommon('common.sortOrder')}</span>
        </label>
        <input
          type="number"
          name="sort_order"
          value={formData.sort_order}
          onChange={handleChange}
          className="input input-bordered"
          min="0"
        />
      </div>

      <div className="form-control">
        <label className="label cursor-pointer">
          <span className="label-text">{tCommon('common.isActive')}</span>
          <input
            type="checkbox"
            name="is_active"
            checked={formData.is_active}
            onChange={handleChange}
            className="checkbox checkbox-primary"
          />
        </label>
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
              {tCommon('common.translating')}
            </>
          ) : (
            <>🌍 {tCommon('common.translate')}</>
          )}
        </button>
      </div>

      <div className="flex gap-2 pt-4">
        <button type="submit" className="btn btn-primary flex-1">
          {tCommon('common.save')}
        </button>
        <button
          type="button"
          onClick={onCancel}
          className="btn btn-ghost flex-1"
        >
          {tCommon('common.cancel')}
        </button>
      </div>
    </form>
  );
}
