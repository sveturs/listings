'use client';

import React, { useState } from 'react';
import { useTranslations, useLocale, NextIntlClientProvider } from 'next-intl';
import { useCreateAIProduct } from '@/contexts/CreateAIProductContext';
import { CategoryTreeSelector } from '@/components/common/CategoryTreeSelector';

interface EnhanceViewProps {
  storefrontId: number | null;
  storefrontSlug: string;
}

export default function EnhanceView({
  storefrontId: _storefrontId,
  storefrontSlug: _storefrontSlug,
}: EnhanceViewProps) {
  const t = useTranslations('storefronts');
  const locale = useLocale();
  const { state, setAIData, setView } = useCreateAIProduct();

  // Получаем контент на текущей локали
  const getLocalizedContent = () => {
    const translations = state.aiData.translations || {};
    // Если есть перевод на текущую локаль, используем его
    if (translations[locale]) {
      return {
        title: translations[locale].title,
        description: translations[locale].description,
      };
    }
    // Иначе используем оригинал (английский)
    return {
      title: state.aiData.title,
      description: state.aiData.description,
    };
  };

  const localizedContent = getLocalizedContent();

  const [editedData, setEditedData] = useState({
    title: localizedContent.title,
    description: localizedContent.description,
    price: state.aiData.price,
    stockQuantity: state.aiData.stockQuantity,
    categoryId: state.aiData.categoryId,
    category: state.aiData.category,
  });

  const handleSave = () => {
    setAIData({
      title: editedData.title,
      description: editedData.description,
      price: editedData.price,
      stockQuantity: editedData.stockQuantity,
      categoryId: editedData.categoryId,
      category: editedData.category,
    });
    setView('variants');
  };

  const handleCategoryChange = async (categoryId: number | number[]) => {
    const id = Array.isArray(categoryId) ? categoryId[0] : categoryId;
    // Загружаем информацию о категории для получения имени
    try {
      const response = await fetch(
        `/api/v1/marketplace/categories?page=1&limit=1000`
      );
      if (response.ok) {
        const data = await response.json();
        const category = data.data?.find((cat: any) => cat.id === id);
        if (category) {
          setEditedData((prev) => ({
            ...prev,
            categoryId: id,
            category: category.name,
          }));
        }
      }
    } catch (error) {
      console.error('Failed to load category name:', error);
      setEditedData((prev) => ({ ...prev, categoryId: id }));
    }
  };

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-2xl font-bold mb-2">
          {t('enhanceProduct') || 'Enhance Your Product'}
        </h2>
        <p className="text-base-content/70">
          {t('enhanceDescription') || 'Review and edit AI-generated content'}
        </p>
      </div>

      {/* Title */}
      <div className="form-control">
        <label className="label">
          <span className="label-text font-semibold">
            {t('productName') || 'Title'}
          </span>
        </label>
        <input
          type="text"
          value={editedData.title}
          onChange={(e) =>
            setEditedData((prev) => ({ ...prev, title: e.target.value }))
          }
          className="input input-bordered w-full"
          placeholder={t('productName') || 'Product title'}
        />
      </div>

      {/* Description */}
      <div className="form-control">
        <label className="label">
          <span className="label-text font-semibold">
            {t('description') || 'Description'}
          </span>
        </label>
        <textarea
          value={editedData.description}
          onChange={(e) =>
            setEditedData((prev) => ({ ...prev, description: e.target.value }))
          }
          className="textarea textarea-bordered h-32"
          placeholder={t('productDescription') || 'Product description'}
        />
      </div>

      {/* Price & Stock */}
      <div className="grid grid-cols-2 gap-4">
        <div className="form-control">
          <label className="label">
            <span className="label-text font-semibold">
              {t('price') || 'Price'}
            </span>
          </label>
          <input
            type="number"
            value={editedData.price}
            onChange={(e) =>
              setEditedData((prev) => ({
                ...prev,
                price: Number(e.target.value),
              }))
            }
            className="input input-bordered"
            min="0"
          />
        </div>
        <div className="form-control">
          <label className="label">
            <span className="label-text font-semibold">
              {t('stockQuantity') || 'Stock'}
            </span>
          </label>
          <input
            type="number"
            value={editedData.stockQuantity}
            onChange={(e) =>
              setEditedData((prev) => ({
                ...prev,
                stockQuantity: Number(e.target.value),
              }))
            }
            className="input input-bordered"
            min="0"
          />
        </div>
      </div>

      {/* Category Selection */}
      <div className="form-control">
        <label className="label">
          <span className="label-text font-semibold">
            {t('category') || 'Category'}
          </span>
          <span className="label-text-alt text-info">
            {t('aiDetected') || 'AI Detected'}: {state.aiData.category}
          </span>
        </label>
        <NextIntlClientProvider
          locale={locale}
          messages={{
            marketplace: {
              selectCategory: t('selectCategory'),
              searchCategories: t('searchCategories'),
              categoriesSelected: t.raw('categoriesSelected'),
              apply: t('apply'),
              cancel: t('cancel'),
              categoriesLoadError: t('categoriesLoadError'),
              noCategoriesFound: t('noCategoriesFound'),
            },
          }}
        >
          <CategoryTreeSelector
            value={editedData.categoryId}
            onChange={handleCategoryChange}
            placeholder={t('selectCategory') || 'Select category'}
            showPath={true}
            allowParentSelection={false}
          />
        </NextIntlClientProvider>
      </div>

      {/* Actions */}
      <div className="flex justify-between gap-3">
        <button onClick={() => setView('process')} className="btn btn-outline">
          {t('back') || 'Back'}
        </button>
        <button onClick={handleSave} className="btn btn-primary px-8">
          {t('continueToVariants') || 'Continue to Variants'}
        </button>
      </div>
    </div>
  );
}
