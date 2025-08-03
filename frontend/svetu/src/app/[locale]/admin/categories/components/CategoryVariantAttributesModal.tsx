'use client';

import { useState, useEffect, useCallback } from 'react';
import { useTranslations } from 'next-intl';
import { toast } from '@/utils/toast';
import { adminApi, VariantAttribute } from '@/services/admin';

// Type for the category variant attribute association
interface CategoryVariantAttribute {
  variant_attribute_id: number;
  variant_attribute_name: string;
  is_enabled: boolean;
  is_required: boolean;
  sort_order: number;
}

interface CategoryVariantAttributesModalProps {
  categoryId: number;
  categoryName: string;
  isOpen: boolean;
  onClose: () => void;
}

export default function CategoryVariantAttributesModal({
  categoryId,
  categoryName,
  isOpen,
  onClose,
}: CategoryVariantAttributesModalProps) {
  const t = useTranslations('admin');
  const [loading, setLoading] = useState(false);
  const [saving, setSaving] = useState(false);
  const [allVariantAttributes, setAllVariantAttributes] = useState<
    VariantAttribute[]
  >([]);
  const [selectedAttributes, setSelectedAttributes] = useState<{
    [key: string]: {
      enabled: boolean;
      sortOrder: number;
      isRequired: boolean;
    };
  }>({});

  const loadData = useCallback(async () => {
    setLoading(true);
    try {
      // Загружаем все вариативные атрибуты
      const variantAttrs = await adminApi.variantAttributes.getAll();
      setAllVariantAttributes(variantAttrs.data);

      // Загружаем текущие настройки для категории
      const categoryVariantAttrs =
        await adminApi.categories.getVariantAttributes(categoryId);

      // Преобразуем в удобный формат
      const selected: typeof selectedAttributes = {};
      categoryVariantAttrs.forEach((attr: CategoryVariantAttribute) => {
        selected[attr.variant_attribute_name] = {
          enabled: true,
          sortOrder: attr.sort_order,
          isRequired: attr.is_required,
        };
      });
      setSelectedAttributes(selected);
    } catch (error) {
      console.error('Failed to load data:', error);
      toast.error(t('common.error'));
    } finally {
      setLoading(false);
    }
  }, [categoryId, t]);

  useEffect(() => {
    if (isOpen) {
      loadData();
    }
  }, [isOpen, categoryId, loadData]);

  const handleToggleAttribute = (attrName: string) => {
    setSelectedAttributes((prev) => {
      const current = prev[attrName];
      if (current?.enabled) {
        // Отключаем
        return {
          ...prev,
          [attrName]: {
            ...current,
            enabled: false,
          },
        };
      } else {
        // Включаем с дефолтными значениями
        return {
          ...prev,
          [attrName]: {
            enabled: true,
            sortOrder: Object.keys(prev).filter((k) => prev[k].enabled).length,
            isRequired: false,
          },
        };
      }
    });
  };

  const handleToggleRequired = (attrName: string) => {
    setSelectedAttributes((prev) => ({
      ...prev,
      [attrName]: {
        ...prev[attrName],
        isRequired: !prev[attrName].isRequired,
      },
    }));
  };

  const handleSortOrderChange = (attrName: string, value: number) => {
    setSelectedAttributes((prev) => ({
      ...prev,
      [attrName]: {
        ...prev[attrName],
        sortOrder: value,
      },
    }));
  };

  const handleSave = async () => {
    setSaving(true);
    try {
      // Преобразуем в формат для API
      const variantAttributes = Object.entries(selectedAttributes)
        .filter(([_, config]) => config.enabled)
        .map(([name, config]) => ({
          variant_attribute_name: name,
          sort_order: config.sortOrder,
          is_required: config.isRequired,
        }));

      await adminApi.categories.updateVariantAttributes(categoryId, {
        variant_attributes: variantAttributes,
      });

      toast.success(t('common.saveSuccess'));
      onClose();
    } catch (error) {
      console.error('Failed to save:', error);
      toast.error(t('common.error'));
    } finally {
      setSaving(false);
    }
  };

  if (!isOpen) return null;

  return (
    <div className="modal modal-open">
      <div className="modal-box w-11/12 max-w-4xl">
        <div className="flex justify-between items-center mb-4">
          <h2 className="text-xl font-bold">
            {t('categories.variantAttributes')}: {categoryName}
          </h2>
          <button
            className="btn btn-sm btn-circle btn-ghost"
            onClick={onClose}
            disabled={saving}
          >
            ✕
          </button>
        </div>

        {loading ? (
          <div className="flex justify-center items-center h-64">
            <span className="loading loading-spinner loading-lg"></span>
          </div>
        ) : (
          <div className="space-y-4">
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
                />
              </svg>
              <div>
                <div className="text-sm">
                  Выберите вариативные атрибуты, которые будут доступны для
                  товаров в этой категории. Порядок отображения можно настроить
                  с помощью числового поля.
                </div>
              </div>
            </div>

            <div className="overflow-x-auto">
              <table className="table table-zebra">
                <thead>
                  <tr>
                    <th>{t('common.active')}</th>
                    <th>{t('attributes.name')}</th>
                    <th>{t('attributes.displayName')}</th>
                    <th>{t('attributes.type')}</th>
                    <th>{t('variantAttributes.affectsStock')}</th>
                    <th>{t('attributes.required')}</th>
                    <th>{t('attributes.sortOrder')}</th>
                  </tr>
                </thead>
                <tbody>
                  {allVariantAttributes.map((attr) => {
                    const config = selectedAttributes[attr.name] || {
                      enabled: false,
                      sortOrder: 0,
                      isRequired: false,
                    };
                    return (
                      <tr
                        key={attr.id}
                        className={config.enabled ? '' : 'opacity-50'}
                      >
                        <td>
                          <label>
                            <input
                              type="checkbox"
                              className="checkbox checkbox-primary"
                              checked={config.enabled}
                              onChange={() => handleToggleAttribute(attr.name)}
                            />
                          </label>
                        </td>
                        <td>
                          <code className="text-sm">{attr.name}</code>
                        </td>
                        <td>{attr.display_name}</td>
                        <td>
                          <span className="badge badge-outline">
                            {t(`variantAttributes.types.${attr.type}`)}
                          </span>
                        </td>
                        <td>
                          {attr.affects_stock ? (
                            <span className="badge badge-warning badge-sm">
                              📦 Да
                            </span>
                          ) : (
                            <span className="badge badge-ghost badge-sm">
                              Нет
                            </span>
                          )}
                        </td>
                        <td>
                          <label>
                            <input
                              type="checkbox"
                              className="checkbox checkbox-sm"
                              checked={config.isRequired}
                              onChange={() => handleToggleRequired(attr.name)}
                              disabled={!config.enabled}
                            />
                          </label>
                        </td>
                        <td>
                          <input
                            type="number"
                            className="input input-bordered input-sm w-20"
                            value={config.sortOrder}
                            onChange={(e) =>
                              handleSortOrderChange(
                                attr.name,
                                parseInt(e.target.value) || 0
                              )
                            }
                            disabled={!config.enabled}
                            min="0"
                          />
                        </td>
                      </tr>
                    );
                  })}
                </tbody>
              </table>
            </div>

            <div className="flex justify-end gap-2 pt-4">
              <button
                className="btn btn-ghost"
                onClick={onClose}
                disabled={saving}
              >
                {t('common.cancel')}
              </button>
              <button
                className="btn btn-primary"
                onClick={handleSave}
                disabled={saving}
              >
                {saving && (
                  <span className="loading loading-spinner loading-sm mr-2"></span>
                )}
                {t('common.save')}
              </button>
            </div>
          </div>
        )}
      </div>
      <div className="modal-backdrop" onClick={onClose}></div>
    </div>
  );
}
