'use client';

import React, { useState } from 'react';
import { useTranslations, useLocale } from 'next-intl';
import { useCreateAIProduct } from '@/contexts/CreateAIProductContext';

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
  });

  const handleSave = () => {
    setAIData({
      title: editedData.title,
      description: editedData.description,
      price: editedData.price,
      stockQuantity: editedData.stockQuantity,
    });
    setView('publish');
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

      {/* Category Info */}
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
          <div className="font-semibold">
            {t('aiDetectedCategory') || 'AI Detected Category'}
          </div>
          <div className="text-sm">{state.aiData.category}</div>
        </div>
      </div>

      {/* Actions */}
      <div className="flex justify-between gap-3">
        <button onClick={() => setView('process')} className="btn btn-outline">
          {t('back') || 'Back'}
        </button>
        <button onClick={handleSave} className="btn btn-primary px-8">
          {t('continueToPublish') || 'Continue to Publish'}
        </button>
      </div>
    </div>
  );
}
