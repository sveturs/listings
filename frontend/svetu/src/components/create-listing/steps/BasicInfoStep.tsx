'use client';

import { useState } from 'react';
import { useTranslations } from 'next-intl';
import { useCreateListing } from '@/contexts/CreateListingContext';

interface BasicInfoStepProps {
  onNext: () => void;
  onBack: () => void;
}

export default function BasicInfoStep({ onNext, onBack }: BasicInfoStepProps) {
  const t = useTranslations('create_listing');
  const tCreate_listing.basic_info = useTranslations('create_listing.basic_info');
  const tCommon = useTranslations('common');
  const tCreate_listing.regional_tips = useTranslations('create_listing.regional_tips');
  const { state, setBasicInfo, setLocalization } = useCreateListing();
  const [formData, setFormData] = useState({
    title: state.title || '',
    description: state.description || '',
    price: state.price || 0,
    currency: state.currency || 'RSD',
    condition: state.condition || 'used',
  });

  const [scriptMode, setScriptMode] = useState<'cyrillic' | 'latin' | 'mixed'>(
    state.localization.script || 'cyrillic'
  );

  const currencies = [
    { code: 'RSD', symbol: 'РСД', name: 'Српски динар', popular: true },
    { code: 'EUR', symbol: '€', name: 'Евро', popular: true },
    { code: 'HRK', symbol: 'kn', name: 'Хрватска куна', popular: false },
    { code: 'MKD', symbol: 'ден', name: 'Македонски денар', popular: false },
  ];

  const conditions = [
    {
      id: 'new',
      label: 'condition.new',
      icon: '✨',
      description: 'condition.new_desc',
    },
    {
      id: 'used',
      label: 'condition.used',
      icon: '👍',
      description: 'condition.used_desc',
      popular: true,
    },
    {
      id: 'refurbished',
      label: 'condition.refurbished',
      icon: '🔧',
      description: 'condition.refurbished_desc',
    },
  ];

  const canProceed =
    formData.title.trim() && formData.description.trim() && formData.price > 0;

  const convertCyrillic = (text: string) => {
    const cyrillicMap: { [key: string]: string } = {
      a: 'а',
      b: 'б',
      v: 'в',
      g: 'г',
      d: 'д',
      đ: 'ђ',
      e: 'е',
      ž: 'ж',
      z: 'з',
      i: 'и',
      j: 'ј',
      k: 'к',
      l: 'л',
      m: 'м',
      n: 'н',
      nj: 'њ',
      o: 'о',
      p: 'п',
      r: 'р',
      s: 'с',
      t: 'т',
      ć: 'ћ',
      u: 'у',
      f: 'ф',
      h: 'х',
      c: 'ц',
      č: 'ч',
      dž: 'џ',
      š: 'ш',
    };

    return text.replace(/[a-zA-ZđžćčšnjdzDŽŽĆČŠNJDZ]/g, (match) => {
      return cyrillicMap[match.toLowerCase()] || match;
    });
  };

  const handleScriptChange = (newScript: 'cyrillic' | 'latin' | 'mixed') => {
    setScriptMode(newScript);
    setLocalization({ script: newScript });
    if (newScript === 'cyrillic' && scriptMode === 'latin') {
      setFormData((prev) => ({
        ...prev,
        title: convertCyrillic(prev.title),
        description: convertCyrillic(prev.description),
      }));
    }
  };

  return (
    <div className="max-w-2xl mx-auto">
      <div className="card bg-base-100 shadow-lg">
        <div className="card-body">
          <h2 className="card-title text-2xl mb-4 flex items-center">
            📝 {tCreate_listing.basic_info('title')}
          </h2>
          <p className="text-base-content/70 mb-6">
            {tCreate_listing.basic_info('description')}
          </p>

          <div className="space-y-6">
            {/* Переключатель скрипта */}
            <div className="form-control">
              <label className="label">
                <span className="label-text font-medium">
                  🔤 {t('script_mode')}
                </span>
              </label>
              <div className="flex gap-2">
                <button
                  type="button"
                  onClick={() => handleScriptChange('cyrillic')}
                  className={`btn btn-sm ${scriptMode === 'cyrillic' ? 'btn-primary' : 'btn-outline'}`}
                >
                  Ћирилица
                </button>
                <button
                  type="button"
                  onClick={() => handleScriptChange('latin')}
                  className={`btn btn-sm ${scriptMode === 'latin' ? 'btn-primary' : 'btn-outline'}`}
                >
                  Latinica
                </button>
                <button
                  type="button"
                  onClick={() => handleScriptChange('mixed')}
                  className={`btn btn-sm ${scriptMode === 'mixed' ? 'btn-primary' : 'btn-outline'}`}
                >
                  Мешано
                </button>
              </div>
            </div>

            {/* Название объявления */}
            <div className="form-control">
              <label className="label">
                <span className="label-text font-medium">
                  📋 {t('title')}
                </span>
                <span className="label-text-alt text-error">*</span>
              </label>
              <input
                type="text"
                placeholder={
                  scriptMode === 'cyrillic'
                    ? 'нпр. Фрижидер Бош, добро стање'
                    : 'npr. Frižider Bosch, dobro stanje'
                }
                className="input input-bordered"
                value={formData.title}
                onChange={(e) => {
                  const newTitle = e.target.value;
                  setFormData((prev) => ({ ...prev, title: newTitle }));
                  setBasicInfo({ title: newTitle });
                }}
                maxLength={80}
              />
              <label className="label">
                <span className="label-text-alt text-base-content/60">
                  {formData.title.length}/80 {t('characters')}
                </span>
              </label>
            </div>

            {/* Описание */}
            <div className="form-control">
              <label className="label">
                <span className="label-text font-medium">
                  📄 {t('description')}
                </span>
                <span className="label-text-alt text-error">*</span>
              </label>
              <textarea
                placeholder={
                  scriptMode === 'cyrillic'
                    ? 'Опишите детаљно ваш производ. Наведите стање, узраст, разлог продаје...'
                    : 'Opišite detaljno vaš proizvod. Navedite stanje, uzrast, razlog prodaje...'
                }
                className="textarea textarea-bordered h-32"
                value={formData.description}
                onChange={(e) => {
                  const newDescription = e.target.value;
                  setFormData((prev) => ({
                    ...prev,
                    description: newDescription,
                  }));
                  setBasicInfo({ description: newDescription });
                }}
                maxLength={1000}
              />
              <label className="label">
                <span className="label-text-alt text-base-content/60">
                  {formData.description.length}/1000{' '}
                  {t('characters')}
                </span>
              </label>
            </div>

            {/* Цена и валута */}
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <div className="form-control">
                <label className="label">
                  <span className="label-text font-medium">
                    💰 {t('price')}
                  </span>
                  <span className="label-text-alt text-error">*</span>
                </label>
                <input
                  type="number"
                  placeholder="0"
                  className="input input-bordered"
                  value={formData.price || ''}
                  onChange={(e) => {
                    const newPrice = parseInt(e.target.value) || 0;
                    setFormData((prev) => ({
                      ...prev,
                      price: newPrice,
                    }));
                    setBasicInfo({ price: newPrice });
                  }}
                  min="0"
                  step="50"
                />
              </div>

              <div className="form-control">
                <label className="label">
                  <span className="label-text font-medium">
                    💱 {t('currency')}
                  </span>
                </label>
                <select
                  className="select select-bordered"
                  value={formData.currency}
                  onChange={(e) => {
                    const newCurrency = e.target.value as any;
                    setFormData((prev) => ({
                      ...prev,
                      currency: newCurrency,
                    }));
                    setBasicInfo({ currency: newCurrency });
                  }}
                >
                  {currencies
                    .filter((c) => c.popular)
                    .map((currency) => (
                      <option key={currency.code} value={currency.code}>
                        {currency.symbol} - {currency.name}
                      </option>
                    ))}
                  <option disabled>────────────────</option>
                  {currencies
                    .filter((c) => !c.popular)
                    .map((currency) => (
                      <option key={currency.code} value={currency.code}>
                        {currency.symbol} - {currency.name}
                      </option>
                    ))}
                </select>
              </div>
            </div>

            {/* Состояние товара */}
            <div className="form-control">
              <label className="label">
                <span className="label-text font-medium">
                  🏷️ {t('condition')}
                </span>
                <span className="label-text-alt text-error">*</span>
              </label>

              <div className="grid gap-3">
                {conditions.map((condition) => (
                  <label key={condition.id} className="cursor-pointer">
                    <input
                      type="radio"
                      name="condition"
                      value={condition.id}
                      checked={formData.condition === condition.id}
                      onChange={(e) => {
                        const newCondition = e.target.value as any;
                        setFormData((prev) => ({
                          ...prev,
                          condition: newCondition,
                        }));
                        setBasicInfo({ condition: newCondition });
                      }}
                      className="sr-only"
                    />
                    <div
                      className={`
                      card border-2 transition-all duration-200
                      ${
                        formData.condition === condition.id
                          ? 'border-primary bg-primary/5'
                          : 'border-base-300 hover:border-primary/50'
                      }
                    `}
                    >
                      <div className="card-body p-4">
                        <div className="flex items-start gap-3">
                          <span className="text-2xl">{condition.icon}</span>
                          <div className="flex-1">
                            <div className="flex items-center gap-2">
                              <h3 className="font-medium">
                                {t(condition.label)}
                              </h3>
                              {condition.popular && (
                                <span className="badge badge-primary badge-sm">
                                  {tCommon('popular')}
                                </span>
                              )}
                            </div>
                            <p className="text-sm text-base-content/60 mt-1">
                              {t(condition.description)}
                            </p>
                          </div>
                          {formData.condition === condition.id && (
                            <svg
                              className="w-6 h-6 text-primary"
                              fill="currentColor"
                              viewBox="0 0 20 20"
                            >
                              <path
                                fillRule="evenodd"
                                d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
                                clipRule="evenodd"
                              />
                            </svg>
                          )}
                        </div>
                      </div>
                    </div>
                  </label>
                ))}
              </div>
            </div>

            {/* Подсказки для региональных пользователей */}
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
                ></path>
              </svg>
              <div className="text-sm">
                <p className="font-medium">
                  💡 {tCreate_listing.regional_tips('title')}
                </p>
                <ul className="text-xs mt-2 space-y-1">
                  <li>• {tCreate_listing.regional_tips('pricing')}</li>
                  <li>• {tCreate_listing.regional_tips('description')}</li>
                  <li>• {tCreate_listing.regional_tips('honesty')}</li>
                </ul>
              </div>
            </div>
          </div>

          {/* Кнопки навигации */}
          <div className="card-actions justify-between mt-6">
            <button className="btn btn-outline" onClick={onBack}>
              ← {tCommon('back')}
            </button>
            <button
              className={`btn btn-primary ${!canProceed ? 'btn-disabled' : ''}`}
              onClick={onNext}
              disabled={!canProceed}
            >
              {tCommon('continue')} →
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
