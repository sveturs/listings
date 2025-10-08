'use client';

import { useTranslations } from 'next-intl';
import { useState, useEffect, useCallback, useMemo } from 'react';
import { apiClient } from '@/services/api-client';
import debounce from 'lodash/debounce';

interface SEOData {
  keywords: string;
  slug: string;
}

interface SEOSectionProps {
  data: SEOData;
  basicData: {
    title: string;
    description: string;
  };
  errors?: Record<string, string>;
  onChange: (field: keyof SEOData, value: string) => void;
  listingId?: number;
}

interface SEOAnalysis {
  titleScore: number;
  descriptionScore: number;
  overallScore: number;
  recommendations: string[];
}

export function SEOSection({
  data,
  basicData,
  errors = {},
  onChange,
  listingId,
}: SEOSectionProps) {
  const t = useTranslations('profile');
  const [analysis, setAnalysis] = useState<SEOAnalysis>({
    titleScore: 0,
    descriptionScore: 0,
    overallScore: 0,
    recommendations: [],
  });
  const [slugStatus, setSlugStatus] = useState<{
    checking: boolean;
    available?: boolean;
    suggestion?: string;
  }>({
    checking: false,
  });

  // Check slug availability
  const checkSlugAvailability = useCallback(
    async (slug: string) => {
      if (!slug) {
        setSlugStatus({ checking: false });
        return;
      }

      setSlugStatus({ checking: true });

      try {
        const response = await apiClient.post<{
          data: {
            available: boolean;
            suggestion?: string;
          };
        }>('/api/v1/c2c/listings/check-slug', {
          slug,
          exclude_id: listingId || 0,
        });

        if (response.data && response.data.data) {
          setSlugStatus({
            checking: false,
            available: response.data.data.available,
            suggestion: response.data.data.suggestion,
          });

          // If slug is not available and we have a suggestion, use it
          if (!response.data.data.available && response.data.data.suggestion) {
            onChange('slug', response.data.data.suggestion);
          }
        }
      } catch (error) {
        console.error('Error checking slug:', error);
        setSlugStatus({ checking: false });
      }
    },
    [listingId, onChange]
  );

  // Debounced slug check
  const debouncedSlugCheck = useMemo(
    () =>
      debounce((slug: string) => {
        checkSlugAvailability(slug);
      }, 500),
    [checkSlugAvailability]
  );

  // Auto-generate slug based on title
  useEffect(() => {
    if (!data.slug && basicData.title) {
      const slug = basicData.title
        .toLowerCase()
        .replace(/[^a-z0-9\s-]/g, '')
        .replace(/\s+/g, '-')
        .substring(0, 50);
      onChange('slug', slug);
    }
  }, [basicData.title, data.slug, onChange]);

  // Check slug when it changes
  useEffect(() => {
    if (data.slug) {
      debouncedSlugCheck(data.slug);
    }
  }, [data.slug, debouncedSlugCheck]);

  // Reset slug status when slug changes from parent
  useEffect(() => {
    setSlugStatus({ checking: false });
  }, [data.slug]);

  // SEO Analysis
  useEffect(() => {
    const analyzeTitle = () => {
      const title = basicData.title;
      let score = 0;
      const recommendations: string[] = [];

      if (title.length >= 30 && title.length <= 60) {
        score += 40;
      } else if (title.length < 30) {
        recommendations.push(t('seo.analysis.titleTooShort'));
      } else {
        recommendations.push(t('seo.analysis.titleTooLong'));
      }

      if (title.length > 0) score += 20;
      if (/[А-Яа-я]/.test(title) || /[A-Za-z]/.test(title)) score += 20;
      if (!/^\s|\s$/.test(title)) score += 10;
      if (!/\s{2,}/.test(title)) score += 10;

      return { score, recommendations };
    };

    const analyzeDescription = () => {
      const desc = basicData.description;
      let score = 0;
      const recommendations: string[] = [];

      if (desc.length >= 120 && desc.length <= 155) {
        score += 40;
      } else if (desc.length < 120) {
        recommendations.push(t('seo.analysis.descTooShort'));
      } else {
        recommendations.push(t('seo.analysis.descTooLong'));
      }

      if (desc.length > 0) score += 20;
      if (/[А-Яа-я]/.test(desc) || /[A-Za-z]/.test(desc)) score += 20;
      if (desc.includes(basicData.title.split(' ')[0])) score += 10;
      if (desc.endsWith('.') || desc.endsWith('!') || desc.endsWith('?'))
        score += 10;

      return { score, recommendations };
    };

    const titleAnalysis = analyzeTitle();
    const descAnalysis = analyzeDescription();
    const overallScore = Math.round(
      (titleAnalysis.score + descAnalysis.score) / 2
    );

    setAnalysis({
      titleScore: titleAnalysis.score,
      descriptionScore: descAnalysis.score,
      overallScore,
      recommendations: [
        ...titleAnalysis.recommendations,
        ...descAnalysis.recommendations,
      ],
    });
  }, [basicData.title, basicData.description, t]);

  const getScoreColor = (score: number) => {
    if (score >= 80) return 'text-success';
    if (score >= 60) return 'text-warning';
    return 'text-error';
  };

  const getScoreBadgeClass = (score: number) => {
    if (score >= 80) return 'badge-success';
    if (score >= 60) return 'badge-warning';
    return 'badge-error';
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
              d="M21 21l-5.197-5.197m0 0A7.5 7.5 0 105.196 5.196a7.5 7.5 0 0010.607 10.607z"
            />
          </svg>
          {t('seo.title')}
        </h3>

        {/* SEO Score */}
        <div className="card bg-base-100 border border-base-300 mb-6">
          <div className="card-body">
            <h4 className="card-title text-base flex items-center gap-2">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
                strokeWidth={1.5}
                stroke="currentColor"
                className="w-5 h-5"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  d="M3 13.125C3 12.504 3.504 12 4.125 12h2.25c.621 0 1.125.504 1.125 1.125v6.75C7.5 20.496 6.996 21 6.375 21h-2.25A1.125 1.125 0 013 19.875v-6.75zM9.75 8.625c0-.621.504-1.125 1.125-1.125h2.25c.621 0 1.125.504 1.125 1.125v11.25c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 01-1.125-1.125V8.625zM16.5 4.125c0-.621.504-1.125 1.125-1.125h2.25C20.496 3 21 3.504 21 4.125v15.75c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 01-1.125-1.125V4.125z"
                />
              </svg>
              {t('seo.analysis.title')}
              <div
                className={`badge ${getScoreBadgeClass(analysis.overallScore)}`}
              >
                {analysis.overallScore}/100
              </div>
            </h4>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mt-4">
              <div className="stat bg-base-200 rounded-lg">
                <div className="stat-title">{t('seo.analysis.titleScore')}</div>
                <div
                  className={`stat-value text-2xl ${getScoreColor(analysis.titleScore)}`}
                >
                  {analysis.titleScore}/100
                </div>
              </div>
              <div className="stat bg-base-200 rounded-lg">
                <div className="stat-title">{t('seo.analysis.descScore')}</div>
                <div
                  className={`stat-value text-2xl ${getScoreColor(analysis.descriptionScore)}`}
                >
                  {analysis.descriptionScore}/100
                </div>
              </div>
            </div>

            {analysis.recommendations.length > 0 && (
              <div className="mt-4">
                <h5 className="font-medium mb-2">
                  {t('seo.analysis.recommendations')}
                </h5>
                <div className="space-y-1">
                  {analysis.recommendations.map((rec, index) => (
                    <div
                      key={index}
                      className="flex items-center gap-2 text-sm"
                    >
                      <svg
                        xmlns="http://www.w3.org/2000/svg"
                        fill="none"
                        viewBox="0 0 24 24"
                        strokeWidth={1.5}
                        stroke="currentColor"
                        className="w-4 h-4 text-warning"
                      >
                        <path
                          strokeLinecap="round"
                          strokeLinejoin="round"
                          d="M12 9v3.75m0-10.036A11.959 11.959 0 013.598 6 11.99 11.99 0 003 9.75c0 5.592 3.824 10.29 9 11.622C17.176 20.04 21 15.342 21 9.75c0-1.31-.21-2.57-.598-3.75h-.152c-3.196 0-6.1-1.249-8.25-3.286zm0 13.036h.008v.008H12v-.008z"
                        />
                      </svg>
                      {rec}
                    </div>
                  ))}
                </div>
              </div>
            )}
          </div>
        </div>

        {/* URL Slug */}
        <div className="form-control mb-4">
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
                  d="M13.19 8.688a4.5 4.5 0 011.242 7.244l-4.5 4.5a4.5 4.5 0 01-6.364-6.364l1.757-1.757m13.35-.622l1.757-1.757a4.5 4.5 0 00-6.364-6.364l-4.5 4.5a4.5 4.5 0 001.242 7.244"
                />
              </svg>
              {t('seo.slug')}
            </span>
          </label>
          <div className="join w-full">
            <span className="join-item bg-base-200 flex items-center px-3 text-sm text-base-content/60">
              svetuplatform.com/listings/
            </span>
            <div className="relative join-item flex-1">
              <input
                type="text"
                className={`input input-bordered w-full pr-10 ${
                  errors.slug ? 'input-error' : ''
                } ${
                  !slugStatus.checking && slugStatus.available === true
                    ? 'input-success'
                    : ''
                } ${
                  !slugStatus.checking && slugStatus.available === false
                    ? 'input-warning'
                    : ''
                } focus:input-primary transition-all`}
                value={data.slug}
                onChange={(e) => {
                  const slug = e.target.value
                    .toLowerCase()
                    .replace(/[^a-z0-9\s-]/g, '')
                    .replace(/\s+/g, '-');
                  onChange('slug', slug);
                }}
                placeholder={t('seo.slugPlaceholder')}
                maxLength={50}
              />
              {/* Status indicator */}
              <div className="absolute inset-y-0 right-0 flex items-center pr-3">
                {slugStatus.checking && (
                  <span className="loading loading-spinner loading-sm text-primary"></span>
                )}
                {!slugStatus.checking && slugStatus.available === true && (
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    fill="none"
                    viewBox="0 0 24 24"
                    strokeWidth={2}
                    stroke="currentColor"
                    className="w-5 h-5 text-success"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      d="M9 12.75L11.25 15 15 9.75M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                    />
                  </svg>
                )}
                {!slugStatus.checking && slugStatus.available === false && (
                  <div
                    className="tooltip tooltip-left"
                    data-tip={t('seo.slugAutoChanged')}
                  >
                    <svg
                      xmlns="http://www.w3.org/2000/svg"
                      fill="none"
                      viewBox="0 0 24 24"
                      strokeWidth={2}
                      stroke="currentColor"
                      className="w-5 h-5 text-warning"
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        d="M12 9v3.75m9-.75a9 9 0 11-18 0 9 9 0 0118 0zm-9 3.75h.008v.008H12v-.008z"
                      />
                    </svg>
                  </div>
                )}
              </div>
            </div>
          </div>
          {errors.slug && (
            <label className="label">
              <span className="label-text-alt text-error">{errors.slug}</span>
            </label>
          )}
          {!errors.slug && slugStatus.suggestion && (
            <label className="label">
              <span className="label-text-alt text-warning">
                {t('seo.slugChanged', { suggestion: slugStatus.suggestion })}
              </span>
            </label>
          )}
        </div>

        {/* Keywords */}
        <div className="form-control mb-4">
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
                  d="M15.75 5.25a3 3 0 013 3m3 0a6 6 0 01-7.029 5.912c-.563-.097-1.159.026-1.563.43L10.5 17.25H8.25v2.25H6v2.25H2.25v-2.818c0-.597.237-1.17.659-1.591l6.499-6.499c.404-.404.527-1 .43-1.563A6 6 0 1121.75 8.25z"
                />
              </svg>
              {t('seo.keywords')}
            </span>
          </label>
          <input
            type="text"
            className={`input input-bordered w-full ${errors.keywords ? 'input-error' : ''} focus:input-primary transition-all`}
            value={data.keywords}
            onChange={(e) => onChange('keywords', e.target.value)}
            placeholder={t('seo.keywordsPlaceholder')}
          />
          <label className="label">
            <span className="label-text-alt text-error">{errors.keywords}</span>
            <span className="label-text-alt">{t('seo.keywordsHint')}</span>
          </label>
        </div>

        {/* Preview */}
        <div className="card bg-base-100 border border-base-300">
          <div className="card-body">
            <h4 className="card-title text-base flex items-center gap-2">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
                strokeWidth={1.5}
                stroke="currentColor"
                className="w-5 h-5"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  d="M2.036 12.322a1.012 1.012 0 010-.639C3.423 7.51 7.36 4.5 12 4.5c4.638 0 8.573 3.007 9.963 7.178.07.207.07.431 0 .639C20.577 16.49 16.64 19.5 12 19.5c-4.638 0-8.573-3.007-9.963-7.178z"
                />
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
                />
              </svg>
              {t('seo.preview.title')}
            </h4>

            <div className="mockup-browser border border-base-300 bg-base-100">
              <div className="mockup-browser-toolbar">
                <div className="input border border-base-300">
                  svetuplatform.com/listings/{data.slug || 'your-listing'}
                </div>
              </div>
              <div className="flex flex-col px-4 py-4 border-t border-base-300">
                <div className="text-blue-600 text-lg hover:underline cursor-pointer">
                  {basicData.title || t('seo.preview.defaultTitle')}
                </div>
                <div className="text-green-700 text-sm">
                  svetuplatform.com › listings › {data.slug || 'your-listing'}
                </div>
                <div className="text-gray-600 text-sm mt-1">
                  {basicData.description?.substring(0, 155) ||
                    t('seo.preview.defaultDescription')}
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* SEO Tips */}
        <div className="alert shadow-lg mt-6">
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
            <h3 className="font-bold">{t('seo.tips.title')}</h3>
            <div className="text-xs">{t('seo.tips.content')}</div>
          </div>
        </div>
      </div>
    </div>
  );
}
