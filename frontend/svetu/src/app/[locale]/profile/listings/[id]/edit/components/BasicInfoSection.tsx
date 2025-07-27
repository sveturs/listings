'use client';

import { useTranslations } from 'next-intl';
import { useEffect, useState } from 'react';

interface BasicInfoSectionProps {
  data: {
    title: string;
    description: string;
    price: number;
    condition: string;
  };
  errors?: Record<string, string>;
  onChange: (field: string, value: string | number) => void;
}

export function BasicInfoSection({
  data,
  errors = {},
  onChange,
}: BasicInfoSectionProps) {
  const t = useTranslations('profile.listings.editListing');
  const [titleSEO, setTitleSEO] = useState<'poor' | 'good' | 'excellent'>(
    'good'
  );

  // SEO анализ заголовка
  useEffect(() => {
    const length = data.title.length;
    if (length < 20 || length > 70) {
      setTitleSEO('poor');
    } else if (length >= 40 && length <= 60) {
      setTitleSEO('excellent');
    } else {
      setTitleSEO('good');
    }
  }, [data.title]);

  const _getSEOColor = () => {
    switch (titleSEO) {
      case 'poor':
        return 'text-error';
      case 'good':
        return 'text-warning';
      case 'excellent':
        return 'text-success';
    }
  };

  const getSEOText = () => {
    switch (titleSEO) {
      case 'poor':
        return t('seo.poor');
      case 'good':
        return t('seo.good');
      case 'excellent':
        return t('seo.excellent');
    }
  };

  return (
    <div className="space-y-6">
      <div>
        <h3 className="text-lg font-semibold mb-4 flex items-center gap-2">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            strokeWidth={1.5}
            stroke="currentColor"
            className="w-5 h-5 text-primary"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              d="M19.5 14.25v-2.625a3.375 3.375 0 00-3.375-3.375h-1.5A1.125 1.125 0 0113.5 7.125v-1.5a3.375 3.375 0 00-3.375-3.375H8.25m0 12.75h7.5m-7.5 3H12M10.5 2.25H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 00-9-9z"
            />
          </svg>
          {t('sections.basicInfo')}
        </h3>

        {/* Title */}
        <div className="form-control">
          <label className="label">
            <span className="label-text font-semibold text-base">
              {t('fields.title')} *
            </span>
            <div
              className={`badge badge-sm ${titleSEO === 'excellent' ? 'badge-success' : titleSEO === 'good' ? 'badge-warning' : 'badge-error'} gap-1`}
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
                strokeWidth={2}
                stroke="currentColor"
                className="w-3 h-3"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  d="M3 13.125C3 12.504 3.504 12 4.125 12h2.25c.621 0 1.125.504 1.125 1.125v6.75C7.5 20.496 6.996 21 6.375 21h-2.25A1.125 1.125 0 013 19.875v-6.75zM9.75 8.625c0-.621.504-1.125 1.125-1.125h2.25c.621 0 1.125.504 1.125 1.125v11.25c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 01-1.125-1.125V8.625zM16.5 4.125c0-.621.504-1.125 1.125-1.125h2.25C20.496 3 21 3.504 21 4.125v15.75c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 01-1.125-1.125V4.125z"
                />
              </svg>
              SEO: {getSEOText()}
            </div>
          </label>
          <input
            type="text"
            className={`input input-bordered w-full ${errors.title ? 'input-error' : ''} focus:input-primary transition-all`}
            value={data.title}
            onChange={(e) => onChange('title', e.target.value)}
            required
            maxLength={100}
            placeholder={t('fields.titlePlaceholder')}
          />
          <label className="label">
            <span className="label-text-alt text-error">{errors.title}</span>
            <span
              className={`label-text-alt ${data.title.length > 90 ? 'text-warning' : ''}`}
            >
              {data.title.length}/100
            </span>
          </label>
        </div>
      </div>

      {/* Description */}
      <div className="form-control">
        <label className="label">
          <span className="label-text font-semibold text-base">
            {t('fields.description')} *
          </span>
          <div className="dropdown dropdown-end">
            <label tabIndex={0} className="btn btn-primary btn-xs gap-1">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
                strokeWidth={2}
                stroke="currentColor"
                className="w-3 h-3"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  d="M9.813 15.904L9 18.75l-.813-2.846a4.5 4.5 0 00-3.09-3.09L2.25 12l2.846-.813a4.5 4.5 0 003.09-3.09L9 5.25l.813 2.846a4.5 4.5 0 003.09 3.09L15.75 12l-2.846.813a4.5 4.5 0 00-3.09 3.09zM18.259 8.715L18 9.75l-.259-1.035a3.375 3.375 0 00-2.455-2.456L14.25 6l1.036-.259a3.375 3.375 0 002.455-2.456L18 2.25l.259 1.035a3.375 3.375 0 002.456 2.456L21.75 6l-1.035.259a3.375 3.375 0 00-2.456 2.456zM16.894 20.567L16.5 21.75l-.394-1.183a2.25 2.25 0 00-1.423-1.423L13.5 18.75l1.183-.394a2.25 2.25 0 001.423-1.423l.394-1.183.394 1.183a2.25 2.25 0 001.423 1.423l1.183.394-1.183.394a2.25 2.25 0 00-1.423 1.423z"
                />
              </svg>
              AI
            </label>
            <ul
              tabIndex={0}
              className="dropdown-content menu p-2 shadow-lg bg-base-100 rounded-box w-52 z-10 border border-base-300"
            >
              <li>
                <a onClick={() => {}} className="gap-2">
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    fill="none"
                    viewBox="0 0 24 24"
                    strokeWidth={1.5}
                    stroke="currentColor"
                    className="w-4 h-4"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      d="M16.862 4.487l1.687-1.688a1.875 1.875 0 112.652 2.652L10.582 16.07a4.5 4.5 0 01-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 011.13-1.897l8.932-8.931zm0 0L19.5 7.125M18 14v4.75A2.25 2.25 0 0115.75 21H5.25A2.25 2.25 0 013 18.75V8.25A2.25 2.25 0 015.25 6H10"
                    />
                  </svg>
                  {t('ai.improveText')}
                </a>
              </li>
              <li>
                <a onClick={() => {}} className="gap-2">
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    fill="none"
                    viewBox="0 0 24 24"
                    strokeWidth={1.5}
                    stroke="currentColor"
                    className="w-4 h-4"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      d="M12 4.5v15m7.5-7.5h-15"
                    />
                  </svg>
                  {t('ai.addDetails')}
                </a>
              </li>
              <li>
                <a onClick={() => {}} className="gap-2">
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    fill="none"
                    viewBox="0 0 24 24"
                    strokeWidth={1.5}
                    stroke="currentColor"
                    className="w-4 h-4"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      d="M3.75 3v11.25A2.25 2.25 0 006 16.5h2.25M3.75 3h-1.5m1.5 0h16.5m0 0h1.5m-1.5 0v11.25A2.25 2.25 0 0118 16.5h-2.25m-7.5 0h7.5m-7.5 0l-1 3m8.5-3l1 3m0 0l.5 1.5m-.5-1.5h-9.5m0 0l-.5 1.5M9 11.25v1.5M12 9v3.75m3-6v6"
                    />
                  </svg>
                  {t('ai.optimizeSEO')}
                </a>
              </li>
            </ul>
          </div>
        </label>
        <textarea
          className={`textarea textarea-bordered h-32 ${errors.description ? 'textarea-error' : ''} focus:textarea-primary transition-all`}
          value={data.description}
          onChange={(e) => onChange('description', e.target.value)}
          required
          maxLength={1000}
          placeholder={t('fields.descriptionPlaceholder')}
        />
        <label className="label">
          <span className="label-text-alt text-error">
            {errors.description}
          </span>
          <span
            className={`label-text-alt ${data.description.length > 900 ? 'text-warning' : ''}`}
          >
            {data.description.length}/1000
          </span>
        </label>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        {/* Price */}
        <div className="form-control">
          <label className="label">
            <span className="label-text font-semibold text-base flex items-center gap-2">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
                strokeWidth={1.5}
                stroke="currentColor"
                className="w-4 h-4"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  d="M12 6v12m-3-2.818l.879.659c1.171.879 3.07.879 4.242 0 1.172-.879 1.172-2.303 0-3.182C13.536 12.219 12.768 12 12 12c-.725 0-1.45-.22-2.003-.659-1.106-.879-1.106-2.303 0-3.182s2.9-.879 4.006 0l.415.33M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                />
              </svg>
              {t('fields.price')} *
            </span>
          </label>
          <div className="join w-full">
            <input
              type="number"
              className={`input input-bordered join-item flex-1 ${errors.price ? 'input-error' : ''} focus:input-primary transition-all`}
              value={data.price}
              onChange={(e) => onChange('price', Number(e.target.value))}
              required
              min="0"
              step="0.01"
              placeholder="0.00"
            />
            <select className="select select-bordered join-item focus:select-primary">
              <option>RSD</option>
              <option>EUR</option>
              <option>USD</option>
            </select>
          </div>
          {errors.price && (
            <label className="label">
              <span className="label-text-alt text-error">{errors.price}</span>
            </label>
          )}
          <div className="text-xs text-base-content/60 mt-1">
            {t('fields.priceHint')}
          </div>
        </div>

        {/* Condition */}
        <div className="form-control">
          <label className="label">
            <span className="label-text font-semibold text-base flex items-center gap-2">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
                strokeWidth={1.5}
                stroke="currentColor"
                className="w-4 h-4"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  d="M9 12.75L11.25 15 15 9.75m-3-7.036A11.959 11.959 0 013.598 6 11.99 11.99 0 003 9.749c0 5.592 3.824 10.29 9 11.623 5.176-1.332 9-6.03 9-11.622 0-1.31-.21-2.571-.598-3.751h-.152c-3.196 0-6.1-1.248-8.25-3.285z"
                />
              </svg>
              {t('fields.condition')} *
            </span>
          </label>
          <select
            className={`select select-bordered w-full ${errors.condition ? 'select-error' : ''} focus:select-primary transition-all`}
            value={data.condition}
            onChange={(e) => onChange('condition', e.target.value)}
            required
          >
            <option value="new">{t('condition.new')}</option>
            <option value="used">{t('condition.used')}</option>
            <option value="refurbished">{t('condition.refurbished')}</option>
          </select>
          {errors.condition && (
            <label className="label">
              <span className="label-text-alt text-error">
                {errors.condition}
              </span>
            </label>
          )}
        </div>
      </div>

      {/* Helpful tips */}
      <div className="alert shadow-lg">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          fill="none"
          viewBox="0 0 24 24"
          className="stroke-info shrink-0 w-6 h-6"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth="2"
            d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
          ></path>
        </svg>
        <div>
          <h3 className="font-bold">{t('tips.title')}</h3>
          <div className="text-xs">{t('tips.goodListing')}</div>
        </div>
      </div>
    </div>
  );
}
