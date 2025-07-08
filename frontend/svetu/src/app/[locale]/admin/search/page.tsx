'use client';

import { useTranslations } from 'next-intl';
import { useState } from 'react';
import SearchWeights from './components/SearchWeights';
import SearchAnalytics from './components/SearchAnalytics';
import SynonymManager from './components/SynonymManager';
import SearchDashboard from './components/SearchDashboard';
import TransliterationConfig from './components/TransliterationConfig';
import SearchConfig from './components/SearchConfig';

export default function SearchAdminPage() {
  const t = useTranslations();
  const [activeTab, setActiveTab] = useState<
    | 'dashboard'
    | 'weights'
    | 'analytics'
    | 'synonyms'
    | 'transliteration'
    | 'config'
  >('dashboard');

  return (
    <div className="container mx-auto p-4">
      <h1 className="text-3xl font-bold mb-6">{t('admin.search.title')}</h1>

      <div className="tabs tabs-boxed mb-6 flex-wrap">
        <button
          className={`tab ${activeTab === 'dashboard' ? 'tab-active' : ''}`}
          onClick={() => setActiveTab('dashboard')}
        >
          {t('admin.search.dashboard.tab')}
        </button>
        <button
          className={`tab ${activeTab === 'analytics' ? 'tab-active' : ''}`}
          onClick={() => setActiveTab('analytics')}
        >
          {t('admin.search.analytics.tab')}
        </button>
        <button
          className={`tab ${activeTab === 'weights' ? 'tab-active' : ''}`}
          onClick={() => setActiveTab('weights')}
        >
          {t('admin.search.weights.tab')}
        </button>
        <button
          className={`tab ${activeTab === 'synonyms' ? 'tab-active' : ''}`}
          onClick={() => setActiveTab('synonyms')}
        >
          {t('admin.search.synonyms.tab')}
        </button>
        <button
          className={`tab ${activeTab === 'transliteration' ? 'tab-active' : ''}`}
          onClick={() => setActiveTab('transliteration')}
        >
          {t('admin.search.transliteration.tab')}
        </button>
        <button
          className={`tab ${activeTab === 'config' ? 'tab-active' : ''}`}
          onClick={() => setActiveTab('config')}
        >
          {t('admin.search.config.tab')}
        </button>
      </div>

      <div className="bg-base-200 rounded-lg p-6">
        {activeTab === 'dashboard' && <SearchDashboard />}
        {activeTab === 'analytics' && <SearchAnalytics />}
        {activeTab === 'weights' && <SearchWeights />}
        {activeTab === 'synonyms' && <SynonymManager />}
        {activeTab === 'transliteration' && <TransliterationConfig />}
        {activeTab === 'config' && <SearchConfig />}
      </div>
    </div>
  );
}
