'use client';

import { useState } from 'react';
import { useTranslations } from 'next-intl';
import SearchAnalytics from './SearchAnalytics';
import SearchWeights from './SearchWeights';
import BehaviorAnalytics from './BehaviorAnalytics';
import SynonymManager from './SynonymManager';
import WeightOptimization from './WeightOptimization';

type TabType =
  | 'analytics'
  | 'behavior'
  | 'weights'
  | 'synonyms'
  | 'optimization';

export default function SearchDashboard() {
  const t = useTranslations('admin.search');
  const [activeTab, setActiveTab] = useState<TabType>('analytics');

  const tabs: { key: TabType; label: string; icon: string }[] = [
    {
      key: 'analytics',
      label: t('tabs.analytics'),
      icon: 'M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z',
    },
    {
      key: 'behavior',
      label: t('tabs.behavior'),
      icon: 'M16 8v8m-4-5v5m-4-2v2m-2 4h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z',
    },
    {
      key: 'weights',
      label: t('tabs.weights'),
      icon: 'M12 6V4m0 2a2 2 0 100 4m0-4a2 2 0 110 4m-6 8a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4m6 6v10m6-2a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4',
    },
    {
      key: 'synonyms',
      label: t('tabs.synonyms'),
      icon: 'M7 7h.01M7 3h5c1.1 0 2 .9 2 2v1M7 7l2-2M7 7l2 2m5-4v1m-2 0v1m-2 0v1m4-1l-2-2m2 2l2 2m-2-2v2m0-4v2',
    },
    {
      key: 'optimization',
      label: t('tabs.optimization'),
      icon: 'M13 10V3L4 14h7v7l9-11h-7z',
    },
  ];

  return (
    <div className="min-h-screen bg-base-200">
      <div className="container mx-auto px-4 py-8">
        {/* Заголовок */}
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-base-content mb-2">
            {t('title')}
          </h1>
          <p className="text-base-content/70">{t('description')}</p>
        </div>

        {/* Навигация по табам */}
        <div className="tabs tabs-boxed bg-base-100 mb-6 p-1">
          {tabs.map((tab) => (
            <button
              key={tab.key}
              className={`tab tab-lg gap-2 ${
                activeTab === tab.key ? 'tab-active' : ''
              }`}
              onClick={() => setActiveTab(tab.key)}
            >
              <svg
                className="w-5 h-5"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d={tab.icon}
                />
              </svg>
              {tab.label}
            </button>
          ))}
        </div>

        {/* Контент табов */}
        <div className="bg-base-100 rounded-xl shadow-lg p-6">
          {activeTab === 'analytics' && <SearchAnalytics />}
          {activeTab === 'behavior' && <BehaviorAnalytics />}
          {activeTab === 'weights' && <SearchWeights />}
          {activeTab === 'synonyms' && <SynonymManager />}
          {activeTab === 'optimization' && <WeightOptimization />}
        </div>
      </div>
    </div>
  );
}
