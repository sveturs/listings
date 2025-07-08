'use client';

import { useTranslations } from 'next-intl';
import { useState } from 'react';
import SearchWeights from '../admin/search/components/SearchWeights';
import SearchAnalytics from '../admin/search/components/SearchAnalytics';
import SynonymManager from '../admin/search/components/SynonymManager';
import SearchDashboard from '../admin/search/components/SearchDashboard';
import TransliterationConfig from '../admin/search/components/TransliterationConfig';
import SearchConfig from '../admin/search/components/SearchConfig';

export default function AdminTestPage() {
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
      <h1 className="text-3xl font-bold mb-6 text-center">🔧 Admin Panel Test (No Auth Guard)</h1>

      <div className="tabs tabs-boxed mb-6 flex-wrap">
        <button
          className={`tab ${activeTab === 'dashboard' ? 'tab-active' : ''}`}
          onClick={() => setActiveTab('dashboard')}
        >
          Dashboard
        </button>
        <button
          className={`tab ${activeTab === 'analytics' ? 'tab-active' : ''}`}
          onClick={() => setActiveTab('analytics')}
        >
          Analytics
        </button>
        <button
          className={`tab ${activeTab === 'weights' ? 'tab-active' : ''}`}
          onClick={() => setActiveTab('weights')}
        >
          Search Weights
        </button>
        <button
          className={`tab ${activeTab === 'synonyms' ? 'tab-active' : ''}`}
          onClick={() => setActiveTab('synonyms')}
        >
          Synonyms
        </button>
        <button
          className={`tab ${activeTab === 'transliteration' ? 'tab-active' : ''}`}
          onClick={() => setActiveTab('transliteration')}
        >
          Transliteration
        </button>
        <button
          className={`tab ${activeTab === 'config' ? 'tab-active' : ''}`}
          onClick={() => setActiveTab('config')}
        >
          Config
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