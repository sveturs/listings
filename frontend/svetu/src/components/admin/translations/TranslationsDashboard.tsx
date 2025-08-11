'use client';

import { useState, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import { toast } from 'react-hot-toast';
import { tokenManager } from '@/utils/tokenManager';
import {
  ChartBarIcon,
  LanguageIcon,
  ExclamationTriangleIcon,
  CheckCircleIcon,
  ArrowPathIcon,
  FunnelIcon,
} from '@heroicons/react/24/outline';
import ModuleExplorer from './ModuleExplorer';
import TranslationEditor from './TranslationEditor';
import StatisticsPanel from './StatisticsPanel';
import SyncManager from './SyncManager';
import ConflictResolver from './ConflictResolver';
import AITranslations from './AITranslations';

interface TranslationStats {
  total_translations: number;
  complete_translations: number;
  missing_translations: number;
  placeholder_count: number;
  language_stats: Record<string, any>;
  module_stats: any[];
  recent_changes: any[];
}

export default function TranslationsDashboard() {
  const t = useTranslations('admin');
  const [activeTab, setActiveTab] = useState<
    'modules' | 'database' | 'sync' | 'conflicts' | 'ai' | 'stats'
  >('modules');
  const [selectedModule, setSelectedModule] = useState<string | null>(null);
  const [statistics, setStatistics] = useState<TranslationStats | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetchStatistics();
  }, []);

  const fetchStatistics = async () => {
    setLoading(true);
    setError(null);
    try {
      const token = tokenManager.getAccessToken();
      if (!token) {
        throw new Error('No authentication token available');
      }

      const response = await fetch(
        '/api/v1/admin/translations/stats/overview',
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );

      if (!response.ok) {
        throw new Error('Failed to fetch statistics');
      }

      const data = await response.json();
      setStatistics(data.data);
    } catch (err) {
      setError((err as Error).message);
      toast.error(t('translations.fetchError'));
    } finally {
      setLoading(false);
    }
  };

  const renderTabContent = () => {
    switch (activeTab) {
      case 'modules':
        return (
          <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
            <div className="lg:col-span-1">
              <ModuleExplorer
                onSelectModule={setSelectedModule}
                selectedModule={selectedModule}
              />
            </div>
            <div className="lg:col-span-2">
              {selectedModule ? (
                <TranslationEditor module={selectedModule} />
              ) : (
                <div className="bg-base-200 rounded-lg p-12 text-center">
                  <LanguageIcon className="h-16 w-16 mx-auto text-base-content/30 mb-4" />
                  <p className="text-base-content/60">
                    {t('translations.selectModulePrompt')}
                  </p>
                </div>
              )}
            </div>
          </div>
        );
      case 'database':
        return (
          <div className="bg-base-100 rounded-lg p-6">
            <h3 className="text-lg font-semibold mb-4">
              {t('translations.databaseTranslations')}
            </h3>
            <p className="text-base-content/60">
              {t('translations.comingSoon')}
            </p>
          </div>
        );
      case 'sync':
        return <SyncManager />;
      case 'conflicts':
        return <ConflictResolver onConflictResolved={fetchStatistics} />;
      case 'ai':
        return <AITranslations onTranslationComplete={fetchStatistics} />;
      case 'stats':
        return (
          <StatisticsPanel
            statistics={statistics}
            onRefresh={fetchStatistics}
          />
        );
      default:
        return null;
    }
  };

  return (
    <div className="container mx-auto px-4 py-8">
      {/* Header */}
      <div className="mb-8">
        <h1 className="text-3xl font-bold mb-2">{t('translations.title')}</h1>
        <p className="text-base-content/60">{t('translations.description')}</p>
      </div>

      {/* Quick Stats */}
      {statistics && (
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
            <div className="stat-desc">
              {t('translations.acrossAllModules')}
            </div>
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
            <div className="stat-title">{t('translations.placeholders')}</div>
            <div className="stat-value text-warning">
              {statistics.placeholder_count}
            </div>
            <div className="stat-desc">
              {t('translations.needsTranslation')}
            </div>
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
      )}

      {/* Tabs */}
      <div className="tabs tabs-boxed mb-6">
        <button
          className={`tab ${activeTab === 'modules' ? 'tab-active' : ''}`}
          onClick={() => setActiveTab('modules')}
        >
          <LanguageIcon className="h-4 w-4 mr-2" />
          {t('translations.frontendModules')}
        </button>
        <button
          className={`tab ${activeTab === 'database' ? 'tab-active' : ''}`}
          onClick={() => setActiveTab('database')}
        >
          <FunnelIcon className="h-4 w-4 mr-2" />
          {t('translations.databaseTranslations')}
        </button>
        <button
          className={`tab ${activeTab === 'sync' ? 'tab-active' : ''}`}
          onClick={() => setActiveTab('sync')}
        >
          <ArrowPathIcon className="h-4 w-4 mr-2" />
          {t('translations.synchronization')}
        </button>
        <button
          className={`tab ${activeTab === 'conflicts' ? 'tab-active' : ''}`}
          onClick={() => setActiveTab('conflicts')}
        >
          <ExclamationTriangleIcon className="h-4 w-4 mr-2" />
          Конфликты
        </button>
        <button
          className={`tab ${activeTab === 'ai' ? 'tab-active' : ''}`}
          onClick={() => setActiveTab('ai')}
        >
          🤖 AI Переводы
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
      {loading && !statistics ? (
        <div className="flex justify-center items-center py-12">
          <span className="loading loading-spinner loading-lg"></span>
        </div>
      ) : error ? (
        <div className="alert alert-error">
          <ExclamationTriangleIcon className="h-6 w-6" />
          <span>{error}</span>
        </div>
      ) : (
        renderTabContent()
      )}
    </div>
  );
}
