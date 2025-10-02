'use client';

import { useState, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import { toast } from '@/utils/toast';
import type { VariantAttributeFull } from '@/types/variant-attributes';
import { apiClient } from '@/services/api-client';

interface VariantAttributeFormProps {
  attribute?: VariantAttributeFull;
  onSuccess?: () => void;
  onCancel?: () => void;
}

export default function VariantAttributeForm({
  attribute,
  onSuccess,
  onCancel,
}: VariantAttributeFormProps) {
  const _t = useTranslations('admin');
  const [loading, setLoading] = useState(false);
  const [formData, setFormData] = useState<Partial<VariantAttributeFull>>({
    name: '',
    display_name: '',
    attribute_type: 'text',
    is_required: false,
    sort_order: 0,
    affects_stock: false,
    ...attribute,
  });

  useEffect(() => {
    if (attribute) {
      setFormData(attribute);
    }
  }, [attribute]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);

    try {
      const isEdit = !!formData.id;

      const url = isEdit
        ? `/admin/variant-attributes/${formData.id}`
        : `/admin/variant-attributes`;

      const response = isEdit
        ? await apiClient.put(url, formData)
        : await apiClient.post(url, formData);

      if (response.data) {
        toast.success(
          isEdit ? 'Вариативный атрибут обновлен' : 'Вариативный атрибут создан'
        );
        onSuccess?.();
      } else {
        toast.error('Ошибка сохранения атрибута');
      }
    } catch (error) {
      console.error('Error saving variant attribute:', error);
      toast.error('Ошибка сохранения атрибута');
    } finally {
      setLoading(false);
    }
  };

  const attributeTypes = [
    { value: 'text', label: 'Текст' },
    { value: 'select', label: 'Выбор из списка' },
    { value: 'multiselect', label: 'Множественный выбор' },
    { value: 'number', label: 'Число' },
    { value: 'boolean', label: 'Да/Нет' },
    { value: 'color', label: 'Цвет' },
  ];

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div className="form-control">
        <label className="label">
          <span className="label-text">Системное имя*</span>
        </label>
        <input
          type="text"
          value={formData.name}
          onChange={(e) => setFormData({ ...formData, name: e.target.value })}
          className="input input-bordered"
          required
          placeholder="например: color, size, memory"
          disabled={!!formData.id} // Не позволяем менять системное имя
        />
        <label className="label">
          <span className="label-text-alt">
            Используется в коде, только латинские буквы и подчеркивания
          </span>
        </label>
      </div>

      <div className="form-control">
        <label className="label">
          <span className="label-text">Отображаемое имя*</span>
        </label>
        <input
          type="text"
          value={formData.display_name}
          onChange={(e) =>
            setFormData({ ...formData, display_name: e.target.value })
          }
          className="input input-bordered"
          required
          placeholder="например: Цвет, Размер, Объем памяти"
        />
        <label className="label">
          <span className="label-text-alt">
            Название, которое увидят пользователи
          </span>
        </label>
      </div>

      <div className="form-control">
        <label className="label">
          <span className="label-text">Тип атрибута*</span>
        </label>
        <select
          value={formData.attribute_type}
          onChange={(e) =>
            setFormData({ ...formData, attribute_type: e.target.value })
          }
          className="select select-bordered"
          required
        >
          {attributeTypes.map((type) => (
            <option key={type.value} value={type.value}>
              {type.label}
            </option>
          ))}
        </select>
      </div>

      <div className="form-control">
        <label className="label">
          <span className="label-text">Порядок сортировки</span>
        </label>
        <input
          type="number"
          value={formData.sort_order}
          onChange={(e) =>
            setFormData({
              ...formData,
              sort_order: parseInt(e.target.value) || 0,
            })
          }
          className="input input-bordered"
          placeholder="0"
        />
        <label className="label">
          <span className="label-text-alt">
            Порядок отображения в списках (меньше = выше)
          </span>
        </label>
      </div>

      <div className="divider">Настройки</div>

      <div className="form-control">
        <label className="label cursor-pointer">
          <span className="label-text">Обязательный атрибут</span>
          <input
            type="checkbox"
            checked={formData.is_required}
            onChange={(e) =>
              setFormData({ ...formData, is_required: e.target.checked })
            }
            className="checkbox checkbox-primary"
          />
        </label>
        <label className="label">
          <span className="label-text-alt">
            Обязателен для заполнения при создании варианта
          </span>
        </label>
      </div>

      <div className="form-control">
        <label className="label cursor-pointer">
          <span className="label-text">Влияет на учет остатков</span>
          <input
            type="checkbox"
            checked={formData.affects_stock}
            onChange={(e) =>
              setFormData({ ...formData, affects_stock: e.target.checked })
            }
            className="checkbox checkbox-warning"
          />
        </label>
        <label className="label">
          <span className="label-text-alt">
            Каждая комбинация значений будет иметь отдельный учет остатков
          </span>
        </label>
      </div>

      <div className="modal-action">
        <button
          type="button"
          className="btn btn-ghost"
          onClick={onCancel}
          disabled={loading}
        >
          Отмена
        </button>
        <button
          type="submit"
          className={`btn btn-primary ${loading ? 'loading' : ''}`}
          disabled={loading}
        >
          {loading ? 'Сохранение...' : formData.id ? 'Обновить' : 'Создать'}
        </button>
      </div>
    </form>
  );
}
