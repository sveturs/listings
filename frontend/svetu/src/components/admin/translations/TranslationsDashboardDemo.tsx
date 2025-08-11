'use client';

import { useState } from 'react';
import { useTranslations } from 'next-intl';
import {
  ChartBarIcon,
  ExclamationTriangleIcon,
  CheckCircleIcon,
  ArrowPathIcon,
} from '@heroicons/react/24/outline';
import AITranslationsDemo from './AITranslationsDemo';

// Demo статистика
const DEMO_STATS = {
  total_translations: 14286,
  complete_translations: 13996,
  missing_translations: 1,
  placeholder_count: 322,
  language_stats: {
    en: {
      total: 504,
      complete: 504,
      machine_translated: 458,
      verified: 60,
      coverage: 100,
    },
    ru: {
      total: 543,
      complete: 543,
      machine_translated: 422,
      verified: 135,
      coverage: 100,
    },
    sr: {
      total: 508,
      complete: 508,
      machine_translated: 443,
      verified: 79,
      coverage: 100,
    },
  },
};

export default function TranslationsDashboardDemo() {
  const t = useTranslations('admin');
  const [activeTab, setActiveTab] = useState<'ai' | 'sync' | 'stats'>('ai');
  const [statistics] = useState(DEMO_STATS);

  return (
    <div className="container mx-auto px-4 py-8">
      {/* Header */}
      <div className="mb-8">
        <h1 className="text-3xl font-bold mb-2">{t('translations.title')}</h1>
        <p className="text-base-content/60">{t('translations.description')}</p>

        <div className="alert alert-warning mt-4">
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
          <span>Демо-режим: Некоторые функции недоступны без авторизации</span>
        </div>
      </div>

      {/* Quick Stats */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-8">
        <div className="stat bg-base-100 rounded-lg">
          <div className="stat-figure text-primary">
            <ChartBarIcon className="h-8 w-8" />
          </div>
          <div className="stat-title">
            {t('translations.totalTranslations')}
          </div>
          <div className="stat-value text-primary">
            {statistics.total_translations.toLocaleString()}
          </div>
          <div className="stat-desc">{t('translations.acrossAllModules')}</div>
        </div>

        <div className="stat bg-base-100 rounded-lg">
          <div className="stat-figure text-success">
            <CheckCircleIcon className="h-8 w-8" />
          </div>
          <div className="stat-title">{t('translations.complete')}</div>
          <div className="stat-value text-success">
            {Math.round(
              (statistics.complete_translations /
                statistics.total_translations) *
                100
            )}
            %
          </div>
          <div className="stat-desc">
            {statistics.complete_translations.toLocaleString()}{' '}
            {t('translations.translations')}
          </div>
        </div>

        <div className="stat bg-base-100 rounded-lg">
          <div className="stat-figure text-warning">
            <ExclamationTriangleIcon className="h-8 w-8" />
          </div>
          <div className="stat-title">Плейсхолдеры</div>
          <div className="stat-value text-warning">
            {statistics.placeholder_count}
          </div>
          <div className="stat-desc">{t('translations.needsTranslation')}</div>
        </div>

        <div className="stat bg-base-100 rounded-lg">
          <div className="stat-figure text-error">
            <ExclamationTriangleIcon className="h-8 w-8" />
          </div>
          <div className="stat-title">{t('translations.missing')}</div>
          <div className="stat-value text-error">
            {statistics.missing_translations}
          </div>
          <div className="stat-desc">{t('translations.notFound')}</div>
        </div>
      </div>

      {/* Tabs */}
      <div className="tabs tabs-boxed mb-6">
        <button
          className={`tab ${activeTab === 'ai' ? 'tab-active' : ''}`}
          onClick={() => setActiveTab('ai')}
        >
          🤖 AI Переводы
        </button>
        <button
          className={`tab ${activeTab === 'sync' ? 'tab-active' : ''}`}
          onClick={() => setActiveTab('sync')}
        >
          <ArrowPathIcon className="h-4 w-4 mr-2" />
          {t('translations.synchronization')}
        </button>
        <button
          className={`tab ${activeTab === 'stats' ? 'tab-active' : ''}`}
          onClick={() => setActiveTab('stats')}
        >
          <ChartBarIcon className="h-4 w-4 mr-2" />
          {t('translations.statistics')}
        </button>
      </div>

      {/* Tab Content */}
      {activeTab === 'ai' && <AITranslationsDemo />}

      {activeTab === 'sync' && (
        <div className="grid gap-6">
          <div className="card bg-base-100 shadow-sm">
            <div className="card-body">
              <h4 className="font-semibold mb-3">Синхронизация переводов</h4>

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
                    ></path>
                  </svg>
                  <span>
                    Система автоматически синхронизирует переводы между Frontend
                    JSON файлами и базой данных
                  </span>
                </div>

                <div className="flex gap-4">
                  <button className="btn btn-primary" disabled>
                    <ArrowPathIcon className="h-4 w-4 mr-2" />
                    Frontend → База данных
                  </button>
                  <button className="btn btn-secondary" disabled>
                    <ArrowPathIcon className="h-4 w-4 mr-2" />
                    База данных → Frontend
                  </button>
                </div>

                <div className="text-sm text-base-content/60">
                  Функции синхронизации требуют авторизации администратора
                </div>
              </div>
            </div>
          </div>
        </div>
      )}

      {activeTab === 'stats' && (
        <div className="grid gap-6">
          <div className="card bg-base-100 shadow-sm">
            <div className="card-body">
              <h4 className="font-semibold mb-4">Статистика по языкам</h4>

              <div className="grid md:grid-cols-3 gap-4">
                {Object.entries(statistics.language_stats).map(
                  ([lang, stats]: [string, any]) => (
                    <div key={lang} className="card bg-base-200">
                      <div className="card-body p-4">
                        <h5 className="font-medium">{lang.toUpperCase()}</h5>
                        <div className="space-y-2 text-sm">
                          <div className="flex justify-between">
                            <span>Всего:</span>
                            <span className="font-medium">{stats.total}</span>
                          </div>
                          <div className="flex justify-between">
                            <span>Завершено:</span>
                            <span className="text-success">
                              {stats.complete}
                            </span>
                          </div>
                          <div className="flex justify-between">
                            <span>Машинный перевод:</span>
                            <span className="text-warning">
                              {stats.machine_translated}
                            </span>
                          </div>
                          <div className="flex justify-between">
                            <span>Проверено:</span>
                            <span className="text-info">{stats.verified}</span>
                          </div>
                          <div className="flex justify-between">
                            <span>Покрытие:</span>
                            <span className="font-bold text-primary">
                              {stats.coverage}%
                            </span>
                          </div>
                        </div>
                      </div>
                    </div>
                  )
                )}
              </div>
            </div>
          </div>

          <div className="card bg-base-100 shadow-sm">
            <div className="card-body">
              <h4 className="font-semibold mb-4">
                Модули с недостающими переводами
              </h4>

              <div className="space-y-2">
                <div className="flex justify-between items-center p-3 bg-warning/10 rounded-lg">
                  <span className="font-medium">storefronts</span>
                  <span className="badge badge-warning">148 плейсхолдеров</span>
                </div>
                <div className="flex justify-between items-center p-3 bg-warning/10 rounded-lg">
                  <span className="font-medium">reviews</span>
                  <span className="badge badge-warning">79 плейсхолдеров</span>
                </div>
                <div className="flex justify-between items-center p-3 bg-warning/10 rounded-lg">
                  <span className="font-medium">orders</span>
                  <span className="badge badge-warning">15 плейсхолдеров</span>
                </div>
                <div className="flex justify-between items-center p-3 bg-warning/10 rounded-lg">
                  <span className="font-medium">search</span>
                  <span className="badge badge-warning">15 плейсхолдеров</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
