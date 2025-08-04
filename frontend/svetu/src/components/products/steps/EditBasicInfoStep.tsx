'use client';

import { useState, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import { useEditProduct } from '@/contexts/EditProductContext';

interface EditBasicInfoStepProps {
  onNext: () => void;
  onBack: () => void;
}

export default function EditBasicInfoStep({
  onNext,
  onBack,
}: EditBasicInfoStepProps) {
  const t = useTranslations('storefronts');
  const tCommon = useTranslations('common');
  const { state, setProductData, setError, clearError } = useEditProduct();

  const [formData, setFormData] = useState({
    name: state.productData.name || '',
    description: state.productData.description || '',
    price: state.productData.price || 0,
    stock_quantity: state.productData.stock_quantity || 0,
    sku: state.productData.sku || '',
    barcode: state.productData.barcode || '',
    is_active: state.productData.is_active ?? true,
  });

  useEffect(() => {
    setProductData(formData);
  }, [formData, setProductData]);

  const handleChange = (field: string, value: any) => {
    setFormData((prev) => ({ ...prev, [field]: value }));
    clearError(field);
  };

  const handleNext = () => {
    // Валидация
    const errors: Record<string, string> = {};

    if (!formData.name || formData.name.length < 3) {
      errors.name = t('nameRequired');
    }

    if (!formData.description || formData.description.length < 10) {
      errors.description = t('descriptionRequired');
    }

    if (!formData.price || formData.price <= 0) {
      errors.price = t('priceRequired');
    }

    // Показываем ошибки
    Object.entries(errors).forEach(([field, message]) => {
      setError(field, message);
    });

    if (Object.keys(errors).length === 0) {
      onNext();
    }
  };

  return (
    <div className="space-y-6">
      {/* Заголовок */}
      <div className="text-center">
        <div className="w-16 h-16 bg-primary/20 rounded-full flex items-center justify-center mx-auto mb-4">
          <span className="text-2xl">📝</span>
        </div>
        <h3 className="text-2xl font-bold text-base-content mb-2">
          {t('products.steps.basic')}
        </h3>
        <p className="text-base-content/70">
          {t('basicInformationDescription')}
        </p>
      </div>

      {/* Форма */}
      <div className="space-y-6">
        {/* Название */}
        <div>
          <label className="label">
            <span className="label-text font-semibold">
              {t('productName')} *
            </span>
          </label>
          <input
            type="text"
            className={`input input-bordered w-full ${state.errors.name ? 'input-error' : ''}`}
            placeholder={t('productNamePlaceholder')}
            value={formData.name}
            onChange={(e) => handleChange('name', e.target.value)}
          />
          {state.errors.name && (
            <div className="label">
              <span className="label-text-alt text-error">
                {state.errors.name}
              </span>
            </div>
          )}
        </div>

        {/* Описание */}
        <div>
          <label className="label">
            <span className="label-text font-semibold">
              {t('description')} *
            </span>
          </label>
          <textarea
            className={`textarea textarea-bordered w-full h-32 ${state.errors.description ? 'textarea-error' : ''}`}
            placeholder={t('descriptionPlaceholder')}
            value={formData.description}
            onChange={(e) => handleChange('description', e.target.value)}
          />
          {state.errors.description && (
            <div className="label">
              <span className="label-text-alt text-error">
                {state.errors.description}
              </span>
            </div>
          )}
        </div>

        {/* Цена и валюта */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label className="label">
              <span className="label-text font-semibold">
                {t('price')} *
              </span>
            </label>
            <input
              type="number"
              className={`input input-bordered w-full ${state.errors.price ? 'input-error' : ''}`}
              placeholder="0"
              min="0"
              step="0.01"
              value={formData.price}
              onChange={(e) =>
                handleChange('price', parseFloat(e.target.value) || 0)
              }
            />
            {state.errors.price && (
              <div className="label">
                <span className="label-text-alt text-error">
                  {state.errors.price}
                </span>
              </div>
            )}
          </div>

          <div>
            <label className="label">
              <span className="label-text font-semibold">
                {t('currency')}
              </span>
            </label>
            <select className="select select-bordered w-full" disabled>
              <option value="RSD">RSD</option>
            </select>
          </div>
        </div>

        {/* Количество на складе */}
        <div>
          <label className="label">
            <span className="label-text font-semibold">
              {t('stockQuantity')}
            </span>
          </label>
          <input
            type="number"
            className="input input-bordered w-full"
            placeholder="0"
            min="0"
            value={formData.stock_quantity}
            onChange={(e) =>
              handleChange('stock_quantity', parseInt(e.target.value) || 0)
            }
          />
          <div className="label">
            <span className="label-text-alt text-base-content/60">
              {t('stockQuantityHelp')}
            </span>
          </div>
        </div>

        {/* SKU и Barcode */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label className="label">
              <span className="label-text font-semibold">
                {t('sku')}
              </span>
            </label>
            <input
              type="text"
              className="input input-bordered w-full"
              placeholder={t('skuPlaceholder')}
              value={formData.sku}
              onChange={(e) => handleChange('sku', e.target.value)}
            />
          </div>

          <div>
            <label className="label">
              <span className="label-text font-semibold">
                {t('barcode')}
              </span>
            </label>
            <input
              type="text"
              className="input input-bordered w-full"
              placeholder={t('barcodePlaceholder')}
              value={formData.barcode}
              onChange={(e) => handleChange('barcode', e.target.value)}
            />
          </div>
        </div>

        {/* Активен */}
        <div className="form-control">
          <label className="cursor-pointer label">
            <span className="label-text font-semibold">
              {t('activeProduct')}
            </span>
            <input
              type="checkbox"
              className="toggle toggle-primary"
              checked={formData.is_active}
              onChange={(e) => handleChange('is_active', e.target.checked)}
            />
          </label>
          <div className="label">
            <span className="label-text-alt text-base-content/60">
              {t('activeProductHelp')}
            </span>
          </div>
        </div>
      </div>

      {/* Кнопки навигации */}
      <div className="flex justify-between">
        <button onClick={onBack} className="btn btn-outline btn-lg">
          {tCommon('back')}
        </button>
        <button onClick={handleNext} className="btn btn-primary btn-lg">
          {tCommon('continue')}
        </button>
      </div>
    </div>
  );
}
