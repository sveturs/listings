'use client';

import { useState, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import { toast } from '@/utils/toast';
import Link from 'next/link';
import { tokenManager } from '@/utils/tokenManager';
import VariantAttributeForm from './components/VariantAttributeForm';

import configManager from '@/config';

interface VariantAttribute {
  id: number;
  name: string;
  display_name: string;
  type: string;
  is_required: boolean;
  sort_order: number;
  affects_stock: boolean;
  created_at: string;
  updated_at: string;
}

interface VariantMapping {
  id: number;
  variant_attribute_id: number;
  category_id: number;
  sort_order: number;
  is_required: boolean;
  attribute?: VariantAttribute;
  category?: {
    id: number;
    name: string;
  };
}

export default function VariantAttributesClient() {
  const _t = useTranslations('admin');
  const [attributes, setAttributes] = useState<VariantAttribute[]>([]);
  const [mappings, setMappings] = useState<VariantMapping[]>([]);
  const [loading, setLoading] = useState(true);
  const [selectedCategory, setSelectedCategory] = useState<number | null>(null);
  const [categories, setCategories] = useState<any[]>([]);
  const [showFormModal, setShowFormModal] = useState(false);
  const [selectedAttribute, setSelectedAttribute] =
    useState<VariantAttribute | null>(null);

  useEffect(() => {
    fetchVariantAttributes();
    fetchCategories();
  }, []);

  useEffect(() => {
    if (selectedCategory) {
      fetchCategoryMappings(selectedCategory);
    }
  }, [selectedCategory]);

  const fetchVariantAttributes = async () => {
    try {
      const token = tokenManager.getAccessToken();
      const apiUrl = configManager.getApiUrl();
      const response = await fetch(
        `${apiUrl}/api/v1/admin/variant-attributes`,
        {
          headers: {
            Authorization: token ? `Bearer ${token}` : '',
          },
        }
      );
      if (response.ok) {
        const data = await response.json();
        setAttributes(data.data || []);
      }
    } catch (error) {
      console.error('Error fetching variant attributes:', error);
      toast.error('Ошибка загрузки вариативных атрибутов');
    } finally {
      setLoading(false);
    }
  };

  const fetchCategories = async () => {
    try {
      const token = tokenManager.getAccessToken();
      const apiUrl = configManager.getApiUrl();
      const response = await fetch(
        `${apiUrl}/api/v1/marketplace/category-tree`,
        {
          headers: {
            Authorization: token ? `Bearer ${token}` : '',
          },
        }
      );
      if (response.ok) {
        const data = await response.json();
        setCategories(data.data || []);
      }
    } catch (error) {
      console.error('Error fetching categories:', error);
    }
  };

  const fetchCategoryMappings = async (categoryId: number) => {
    try {
      const token = tokenManager.getAccessToken();
      const apiUrl = configManager.getApiUrl();
      const response = await fetch(
        `${apiUrl}/api/v1/admin/variant-attributes/mappings?category_id=${categoryId}`,
        {
          headers: {
            Authorization: token ? `Bearer ${token}` : '',
          },
        }
      );
      if (response.ok) {
        const data = await response.json();
        setMappings(data.data || []);
      }
    } catch (error) {
      console.error('Error fetching mappings:', error);
    }
  };

  const handleDelete = async (attr: VariantAttribute) => {
    if (
      !confirm(
        `Удалить атрибут "${attr.display_name}"? Это действие нельзя отменить.`
      )
    ) {
      return;
    }

    try {
      const token = tokenManager.getAccessToken();
      const apiUrl = configManager.getApiUrl();
      const response = await fetch(
        `${apiUrl}/api/v1/admin/variant-attributes/${attr.id}`,
        {
          method: 'DELETE',
          headers: {
            Authorization: token ? `Bearer ${token}` : '',
          },
        }
      );

      if (response.ok) {
        toast.success('Атрибут удален');
        fetchVariantAttributes();
      } else {
        const error = await response.json();
        toast.error(error.error || 'Ошибка удаления атрибута');
      }
    } catch (error) {
      console.error('Error deleting variant attribute:', error);
      toast.error('Ошибка удаления атрибута');
    }
  };

  const handleToggleMapping = async (
    attribute: VariantAttribute,
    categoryId: number,
    isEnabled: boolean
  ) => {
    try {
      if (isEnabled) {
        const token = tokenManager.getAccessToken();
        const apiUrl = configManager.getApiUrl();
        const response = await fetch(
          `${apiUrl}/api/v1/admin/variant-attributes/mappings`,
          {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
              Authorization: token ? `Bearer ${token}` : '',
            },
            body: JSON.stringify({
              variant_attribute_id: attribute.id,
              category_id: categoryId,
              sort_order: 0,
              is_required: false,
            }),
          }
        );

        if (response.ok) {
          toast.success('Атрибут добавлен к категории');
          fetchCategoryMappings(categoryId);
        }
      } else {
        const mapping = mappings.find(
          (m) =>
            m.variant_attribute_id === attribute.id &&
            m.category_id === categoryId
        );
        if (mapping) {
          const apiUrl = configManager.getApiUrl();
          const response = await fetch(
            `${apiUrl}/api/v1/admin/variant-attributes/mappings/${mapping.id}`,
            {
              method: 'DELETE',
            }
          );

          if (response.ok) {
            toast.success('Атрибут удален из категории');
            fetchCategoryMappings(categoryId);
          }
        }
      }
    } catch (error) {
      console.error('Error toggling mapping:', error);
      toast.error('Ошибка при изменении связи');
    }
  };

  const renderCategoryTree = (items: any[], level = 0) => {
    return items.map((category) => (
      <div key={category.id}>
        <div
          className={`p-2 hover:bg-base-200 rounded cursor-pointer flex items-center gap-2`}
          style={{ paddingLeft: `${level * 20 + 8}px` }}
          onClick={() => setSelectedCategory(category.id)}
        >
          {category.icon && <span>{category.icon}</span>}
          <span className={selectedCategory === category.id ? 'font-bold' : ''}>
            {category.name}
          </span>
          {category.children?.length > 0 && (
            <span className="badge badge-sm">{category.children.length}</span>
          )}
        </div>
        {category.children && renderCategoryTree(category.children, level + 1)}
      </div>
    ));
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <span className="loading loading-spinner loading-lg"></span>
      </div>
    );
  }

  return (
    <>
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Список вариативных атрибутов */}
        <div className="lg:col-span-1">
          <div className="card bg-base-100 shadow-xl">
            <div className="card-body">
              <div className="flex justify-between items-start mb-4">
                <div>
                  <h2 className="card-title">Вариативные атрибуты</h2>
                  <p className="text-sm text-base-content/70">
                    Атрибуты для создания вариантов
                  </p>
                </div>
                <button
                  className="btn btn-primary btn-sm"
                  onClick={() => {
                    setSelectedAttribute(null);
                    setShowFormModal(true);
                  }}
                >
                  + Добавить
                </button>
              </div>

              <div className="space-y-2 max-h-[600px] overflow-y-auto">
                {attributes.map((attr) => (
                  <div
                    key={attr.id}
                    className="p-3 border rounded-lg hover:bg-base-200 transition-colors"
                  >
                    <div className="flex justify-between items-start">
                      <div className="flex-1">
                        <div className="font-medium">{attr.display_name}</div>
                        <div className="text-sm text-base-content/70">
                          {attr.name} • {attr.type}
                        </div>
                        <div className="flex gap-2 mt-2">
                          {attr.affects_stock && (
                            <span className="badge badge-sm badge-warning">
                              📦 Влияет на остатки
                            </span>
                          )}
                          {attr.is_required && (
                            <span className="badge badge-sm badge-info">
                              ✅ Обязательный
                            </span>
                          )}
                        </div>
                      </div>
                      <div className="flex gap-1">
                        <button
                          className="btn btn-ghost btn-xs"
                          onClick={() => {
                            setSelectedAttribute(attr);
                            setShowFormModal(true);
                          }}
                        >
                          ✏️
                        </button>
                        <button
                          className="btn btn-ghost btn-xs text-error"
                          onClick={() => handleDelete(attr)}
                        >
                          🗑️
                        </button>
                      </div>
                    </div>
                  </div>
                ))}
              </div>

              <div className="mt-4">
                <Link
                  href="/admin/attributes"
                  className="btn btn-primary btn-sm w-full"
                >
                  Управление атрибутами
                </Link>
              </div>
            </div>
          </div>
        </div>

        {/* Категории */}
        <div className="lg:col-span-1">
          <div className="card bg-base-100 shadow-xl">
            <div className="card-body">
              <h2 className="card-title">Категории</h2>
              <p className="text-sm text-base-content/70 mb-4">
                Выберите категорию для настройки вариантов
              </p>

              <div className="max-h-[600px] overflow-y-auto">
                {renderCategoryTree(categories)}
              </div>
            </div>
          </div>
        </div>

        {/* Настройки для выбранной категории */}
        <div className="lg:col-span-1">
          <div className="card bg-base-100 shadow-xl">
            <div className="card-body">
              {selectedCategory ? (
                <>
                  <h2 className="card-title">Вариативные атрибуты категории</h2>
                  <p className="text-sm text-base-content/70 mb-4">
                    Настройте какие атрибуты могут использоваться как варианты
                  </p>

                  <div className="space-y-2">
                    {attributes.map((attr) => {
                      const mapping = mappings.find(
                        (m) => m.variant_attribute_id === attr.id
                      );
                      const isEnabled = !!mapping;

                      return (
                        <div key={attr.id} className="form-control">
                          <label className="label cursor-pointer">
                            <div className="flex-1">
                              <span className="label-text font-medium">
                                {attr.display_name}
                              </span>
                              <div className="text-xs text-base-content/60">
                                {attr.name}
                                {attr.affects_stock && ' • 📦 Остатки'}
                                {attr.is_required && ' • ✅ Обязательный'}
                              </div>
                            </div>
                            <input
                              type="checkbox"
                              checked={isEnabled}
                              onChange={(e) =>
                                handleToggleMapping(
                                  attr,
                                  selectedCategory,
                                  e.target.checked
                                )
                              }
                              className="checkbox checkbox-primary"
                            />
                          </label>

                          {mapping && (
                            <div className="ml-4 mt-2 p-2 bg-base-200 rounded">
                              <div className="flex items-center gap-2">
                                <label className="label cursor-pointer p-0">
                                  <span className="label-text text-xs">
                                    Обязательный
                                  </span>
                                  <input
                                    type="checkbox"
                                    checked={mapping.is_required}
                                    onChange={async (e) => {
                                      try {
                                        const token =
                                          tokenManager.getAccessToken();
                                        const apiUrl =
                                          configManager.getApiUrl();
                                        const response = await fetch(
                                          `${apiUrl}/api/v1/admin/variant-attributes/mappings/${mapping.id}`,
                                          {
                                            method: 'PATCH',
                                            headers: {
                                              'Content-Type':
                                                'application/json',
                                              Authorization: token
                                                ? `Bearer ${token}`
                                                : '',
                                            },
                                            body: JSON.stringify({
                                              is_required: e.target.checked,
                                            }),
                                          }
                                        );

                                        if (response.ok) {
                                          fetchCategoryMappings(
                                            selectedCategory
                                          );
                                        }
                                      } catch (error) {
                                        console.error(
                                          'Error updating mapping:',
                                          error
                                        );
                                      }
                                    }}
                                    className="checkbox checkbox-xs ml-2"
                                  />
                                </label>
                              </div>
                            </div>
                          )}
                        </div>
                      );
                    })}
                  </div>

                  {mappings.length > 0 && (
                    <div className="alert alert-info mt-4">
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
                      <div className="text-sm">
                        Активно {mappings.length} вариативных атрибутов для этой
                        категории
                      </div>
                    </div>
                  )}
                </>
              ) : (
                <div className="text-center py-8 text-base-content/60">
                  Выберите категорию для настройки вариативных атрибутов
                </div>
              )}
            </div>
          </div>
        </div>
      </div>

      {/* Modal для формы создания/редактирования */}
      {showFormModal && (
        <div className="modal modal-open">
          <div className="modal-box max-w-2xl">
            <h3 className="font-bold text-lg mb-4">
              {selectedAttribute
                ? `Редактировать: ${selectedAttribute.display_name}`
                : 'Новый вариативный атрибут'}
            </h3>

            <VariantAttributeForm
              attribute={selectedAttribute || undefined}
              onSuccess={() => {
                setShowFormModal(false);
                fetchVariantAttributes();
                setSelectedAttribute(null);
              }}
              onCancel={() => {
                setShowFormModal(false);
                setSelectedAttribute(null);
              }}
            />
          </div>
          <div
            className="modal-backdrop"
            onClick={() => setShowFormModal(false)}
          ></div>
        </div>
      )}
    </>
  );
}
