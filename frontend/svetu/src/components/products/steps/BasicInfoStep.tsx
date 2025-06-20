'use client';

import { useState, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import { useCreateProduct } from '@/contexts/CreateProductContext';

interface BasicInfoStepProps {
  onNext: () => void;
  onBack: () => void;
}

export default function BasicInfoStep({ onNext, onBack }: BasicInfoStepProps) {
  const t = useTranslations();
  const { state, setProductData, setError, clearError } = useCreateProduct();
  const [formData, setFormData] = useState({
    name: state.productData.name || '',
    description: state.productData.description || '',
    price: state.productData.price || 0,
    currency: state.productData.currency || 'RSD',
    stock_quantity: state.productData.stock_quantity || 0,
    sku: state.productData.sku || '',
    barcode: state.productData.barcode || '',
    is_active:
      state.productData.is_active !== undefined
        ? state.productData.is_active
        : true,
  });

  useEffect(() => {
    // Синхронизируем с глобальным состоянием
    setProductData(formData);
  }, [formData, setProductData]);

  const handleChange = (
    e: React.ChangeEvent<
      HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement
    >
  ) => {
    const { name, value, type } = e.target;

    if (type === 'checkbox') {
      const checked = (e.target as HTMLInputElement).checked;
      setFormData((prev) => ({ ...prev, [name]: checked }));
      clearError(name);
    } else if (name === 'price' || name === 'stock_quantity') {
      const numValue = parseFloat(value) || 0;
      setFormData((prev) => ({ ...prev, [name]: numValue }));
      clearError(name);
    } else {
      setFormData((prev) => ({ ...prev, [name]: value }));
      clearError(name);
    }
  };

  const validateForm = (): boolean => {
    let isValid = true;

    if (!formData.name || formData.name.length < 3) {
      setError('name', t('storefronts.products.nameRequired'));
      isValid = false;
    }

    if (!formData.description || formData.description.length < 10) {
      setError('description', t('storefronts.products.descriptionRequired'));
      isValid = false;
    }

    if (formData.price <= 0) {
      setError('price', t('storefronts.products.priceRequired'));
      isValid = false;
    }

    return isValid;
  };

  const handleNext = () => {
    if (validateForm()) {
      onNext();
    }
  };

  return (
    <div className="max-w-4xl mx-auto">
      <div className="text-center mb-8">
        <h2 className="text-3xl font-bold text-base-content mb-4">
          {t('storefronts.products.basicInformation')}
        </h2>
        <p className="text-lg text-base-content/70">
          {t('storefronts.products.basicInformationDescription')}
        </p>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
        {/* Левая колонка - основная информация */}
        <div className="space-y-6">
          <div className="card bg-base-100 shadow-xl">
            <div className="card-body">
              <h3 className="card-title text-xl mb-4 flex items-center gap-2">
                <span className="text-2xl">📝</span>
                {t('storefronts.products.productDetails')}
              </h3>

              {/* Название товара */}
              <div className="form-control">
                <label className="label">
                  <span className="label-text font-semibold">
                    {t('storefronts.products.productName')} *
                  </span>
                </label>
                <input
                  type="text"
                  name="name"
                  value={formData.name}
                  onChange={handleChange}
                  className={`input input-bordered input-lg ${state.errors.name ? 'input-error' : ''}`}
                  placeholder={t('storefronts.products.productNamePlaceholder')}
                />
                {state.errors.name && (
                  <label className="label">
                    <span className="label-text-alt text-error">
                      {state.errors.name}
                    </span>
                  </label>
                )}
              </div>

              {/* Описание */}
              <div className="form-control">
                <label className="label">
                  <span className="label-text font-semibold">
                    {t('storefronts.products.description')} *
                  </span>
                  <span className="label-text-alt text-base-content/60">
                    {t('storefronts.products.descriptionHelp')}
                  </span>
                </label>
                <textarea
                  name="description"
                  value={formData.description}
                  onChange={handleChange}
                  className={`textarea textarea-bordered h-32 ${state.errors.description ? 'textarea-error' : ''}`}
                  placeholder={t('storefronts.products.descriptionPlaceholder')}
                />
                {state.errors.description && (
                  <label className="label">
                    <span className="label-text-alt text-error">
                      {state.errors.description}
                    </span>
                  </label>
                )}
              </div>
            </div>
          </div>
        </div>

        {/* Правая колонка - ценообразование и инвентарь */}
        <div className="space-y-6">
          {/* Цена */}
          <div className="card bg-base-100 shadow-xl">
            <div className="card-body">
              <h3 className="card-title text-xl mb-4 flex items-center gap-2">
                <span className="text-2xl">💰</span>
                {t('storefronts.products.pricing')}
              </h3>

              <div className="grid grid-cols-2 gap-4">
                <div className="form-control">
                  <label className="label">
                    <span className="label-text font-semibold">
                      {t('storefronts.products.price')} *
                    </span>
                  </label>
                  <div className="input-group">
                    <input
                      type="number"
                      name="price"
                      value={formData.price}
                      onChange={handleChange}
                      className={`input input-bordered w-full ${state.errors.price ? 'input-error' : ''}`}
                      min="0"
                      step="0.01"
                    />
                    <span className="bg-base-200 text-base-content px-4 flex items-center">
                      {formData.currency}
                    </span>
                  </div>
                  {state.errors.price && (
                    <label className="label">
                      <span className="label-text-alt text-error">
                        {state.errors.price}
                      </span>
                    </label>
                  )}
                </div>

                <div className="form-control">
                  <label className="label">
                    <span className="label-text font-semibold">
                      {t('storefronts.products.currency')}
                    </span>
                  </label>
                  <select
                    name="currency"
                    value={formData.currency}
                    onChange={handleChange}
                    className="select select-bordered"
                  >
                    <option value="RSD">RSD - Serbian Dinar</option>
                    <option value="EUR">EUR - Euro</option>
                    <option value="USD">USD - US Dollar</option>
                  </select>
                </div>
              </div>
            </div>
          </div>

          {/* Инвентарь */}
          <div className="card bg-base-100 shadow-xl">
            <div className="card-body">
              <h3 className="card-title text-xl mb-4 flex items-center gap-2">
                <span className="text-2xl">📦</span>
                {t('storefronts.products.inventory')}
              </h3>

              <div className="grid grid-cols-2 gap-4">
                <div className="form-control">
                  <label className="label">
                    <span className="label-text font-semibold">
                      {t('storefronts.products.stockQuantity')}
                    </span>
                    <span className="label-text-alt text-base-content/60">
                      {t('storefronts.products.stockQuantityHelp')}
                    </span>
                  </label>
                  <input
                    type="number"
                    name="stock_quantity"
                    value={formData.stock_quantity}
                    onChange={handleChange}
                    className="input input-bordered"
                    min="0"
                  />
                </div>

                <div className="form-control">
                  <label className="label">
                    <span className="label-text font-semibold">
                      {t('storefronts.products.sku')}
                    </span>
                  </label>
                  <input
                    type="text"
                    name="sku"
                    value={formData.sku}
                    onChange={handleChange}
                    className="input input-bordered"
                    placeholder={t('storefronts.products.skuPlaceholder')}
                  />
                </div>
              </div>

              <div className="form-control">
                <label className="label">
                  <span className="label-text font-semibold">
                    {t('storefronts.products.barcode')}
                  </span>
                </label>
                <input
                  type="text"
                  name="barcode"
                  value={formData.barcode}
                  onChange={handleChange}
                  className="input input-bordered"
                  placeholder={t('storefronts.products.barcodePlaceholder')}
                />
              </div>

              <div className="form-control">
                <label className="label cursor-pointer">
                  <span className="label-text font-semibold">
                    {t('storefronts.products.activeProduct')}
                  </span>
                  <input
                    type="checkbox"
                    name="is_active"
                    checked={formData.is_active}
                    onChange={handleChange}
                    className="toggle toggle-primary"
                  />
                </label>
                <label className="label">
                  <span className="label-text-alt text-base-content/60">
                    {t('storefronts.products.activeProductHelp')}
                  </span>
                </label>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Навигация */}
      <div className="flex justify-between items-center mt-8">
        <button onClick={onBack} className="btn btn-outline btn-lg px-8">
          <svg
            className="w-5 h-5 mr-2"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M15 19l-7-7 7-7"
            />
          </svg>
          {t('common.back')}
        </button>

        <button onClick={handleNext} className="btn btn-primary btn-lg px-8">
          {t('common.next')}
          <svg
            className="w-5 h-5 ml-2"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M9 5l7 7-7 7"
            />
          </svg>
        </button>
      </div>
    </div>
  );
}
