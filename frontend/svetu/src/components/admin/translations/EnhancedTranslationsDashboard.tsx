'use client';

import { useState } from 'react';
import { useTranslations } from 'next-intl';
import {
  ChartBarIcon,
  ClockIcon,
  DocumentTextIcon,
  ArrowPathIcon,
  PlayIcon,
  DocumentArrowDownIcon,
} from '@heroicons/react/24/outline';

// Import our new components
import BulkTranslationManager from './BulkTranslationManager';
import VersionHistoryViewer from './VersionHistoryViewer';
import AuditLogViewer from './AuditLogViewer';
import ExportImportManager from './ExportImportManager';

// Import existing demo component for backward compatibility
import AITranslationsDemo from './AITranslationsDemo';
import AICostsMonitor from './AICostsMonitor';

export default function EnhancedTranslationsDashboard() {
  const _t = useTranslations('admin');
  const [activeTab, setActiveTab] = useState<
    'overview' | 'bulk' | 'audit' | 'export' | 'ai' | 'sync' | 'stats' | 'costs'
  >('overview');
  const [_showVersionHistory, _setShowVersionHistory] = useState(false);
  const [_versionHistoryParams, _setVersionHistoryParams] = useState<{
    entityType: string;
    entityId: number;
  } | null>(null);

  // Demo statistics for overview
  const [statistics] = useState({
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
  });

  const _openVersionHistory = (_entityType: string, _entityId: number) => {
    // setVersionHistoryParams({ entityType, entityId });
    // setShowVersionHistory(true);
  };

  const _closeVersionHistory = () => {
    // setShowVersionHistory(false);
    // setVersionHistoryParams(null);
  };

  return (
    <div className="container mx-auto px-4 py-8">
      {/* Header */}
      <div className="mb-8">
        <h1 className="text-3xl font-bold mb-2">Система переводов</h1>
        <p className="text-base-content/60">
          Расширенное управление переводами с версионированием, аудитом и
          массовыми операциями
        </p>
      </div>

      {/* Navigation Tabs */}
      <div className="tabs tabs-boxed mb-6">
        <button
          className={`tab ${activeTab === 'overview' ? 'tab-active' : ''}`}
          onClick={() => setActiveTab('overview')}
        >
          <ChartBarIcon className="h-4 w-4 mr-2" />
          Обзор
        </button>
        <button
          className={`tab ${activeTab === 'bulk' ? 'tab-active' : ''}`}
          onClick={() => setActiveTab('bulk')}
        >
          <PlayIcon className="h-4 w-4 mr-2" />
          Массовый перевод
        </button>
        <button
          className={`tab ${activeTab === 'audit' ? 'tab-active' : ''}`}
          onClick={() => setActiveTab('audit')}
        >
          <DocumentTextIcon className="h-4 w-4 mr-2" />
          Аудит
        </button>
        <button
          className={`tab ${activeTab === 'export' ? 'tab-active' : ''}`}
          onClick={() => setActiveTab('export')}
        >
          <DocumentArrowDownIcon className="h-4 w-4 mr-2" />
          Импорт/Экспорт
        </button>
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
          Синхронизация
        </button>
        <button
          className={`tab ${activeTab === 'stats' ? 'tab-active' : ''}`}
          onClick={() => setActiveTab('stats')}
        >
          <ChartBarIcon className="h-4 w-4 mr-2" />
          Статистика
        </button>
        <button
          className={`tab ${activeTab === 'costs' ? 'tab-active' : ''}`}
          onClick={() => setActiveTab('costs')}
        >
          💰 Расходы AI
        </button>
      </div>

      {/* Tab Content */}
      <div className="min-h-96">
        {/* Overview Tab */}
        {activeTab === 'overview' && (
          <div className="space-y-8">
            {/* Quick Stats */}
            <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
              <div className="stat bg-base-100 rounded-lg">
                <div className="stat-figure text-primary">
                  <ChartBarIcon className="h-8 w-8" />
                </div>
                <div className="stat-title">Всего переводов</div>
                <div className="stat-value text-primary">
                  {statistics.total_translations.toLocaleString()}
                </div>
                <div className="stat-desc">Во всех модулях</div>
              </div>

              <div className="stat bg-base-100 rounded-lg">
                <div className="stat-figure text-success">
                  <ChartBarIcon className="h-8 w-8" />
                </div>
                <div className="stat-title">Завершено</div>
                <div className="stat-value text-success">
                  {Math.round(
                    (statistics.complete_translations /
                      statistics.total_translations) *
                      100
                  )}
                  %
                </div>
                <div className="stat-desc">
                  {statistics.complete_translations.toLocaleString()} переводов
                </div>
              </div>

              <div className="stat bg-base-100 rounded-lg">
                <div className="stat-figure text-warning">
                  <ChartBarIcon className="h-8 w-8" />
                </div>
                <div className="stat-title">Плейсхолдеры</div>
                <div className="stat-value text-warning">
                  {statistics.placeholder_count}
                </div>
                <div className="stat-desc">Требуют перевода</div>
              </div>

              <div className="stat bg-base-100 rounded-lg">
                <div className="stat-figure text-error">
                  <ChartBarIcon className="h-8 w-8" />
                </div>
                <div className="stat-title">Отсутствуют</div>
                <div className="stat-value text-error">
                  {statistics.missing_translations}
                </div>
                <div className="stat-desc">Ключи не найдены</div>
              </div>
            </div>

            {/* Quick Actions */}
            <div className="grid md:grid-cols-3 gap-6">
              <div className="card bg-base-100 shadow-sm">
                <div className="card-body">
                  <h3 className="card-title">Массовые операции</h3>
                  <p className="text-base-content/60 mb-4">
                    Переводите множество элементов одновременно
                  </p>
                  <button
                    className="btn btn-primary"
                    onClick={() => setActiveTab('bulk')}
                  >
                    Открыть менеджер
                  </button>
                </div>
              </div>

              <div className="card bg-base-100 shadow-sm">
                <div className="card-body">
                  <h3 className="card-title">Журнал изменений</h3>
                  <p className="text-base-content/60 mb-4">
                    Просматривайте историю всех операций
                  </p>
                  <button
                    className="btn btn-secondary"
                    onClick={() => setActiveTab('audit')}
                  >
                    Открыть аудит
                  </button>
                </div>
              </div>

              <div className="card bg-base-100 shadow-sm">
                <div className="card-body">
                  <h3 className="card-title">Импорт/Экспорт</h3>
                  <p className="text-base-content/60 mb-4">
                    Работайте с переводами через файлы
                  </p>
                  <button
                    className="btn btn-accent"
                    onClick={() => setActiveTab('export')}
                  >
                    Управление файлами
                  </button>
                </div>
              </div>
            </div>

            {/* Language Statistics */}
            <div className="card bg-base-100 shadow-sm">
              <div className="card-body">
                <h3 className="card-title mb-4">Покрытие по языкам</h3>
                <div className="grid md:grid-cols-3 gap-4">
                  {Object.entries(statistics.language_stats).map(
                    ([lang, stats]: [string, any]) => (
                      <div key={lang} className="card bg-base-200">
                        <div className="card-body p-4">
                          <h4 className="font-medium text-lg">
                            {lang.toUpperCase()}
                            {lang === 'sr' && ' 🇷🇸'}
                            {lang === 'en' && ' 🇺🇸'}
                            {lang === 'ru' && ' 🇷🇺'}
                          </h4>
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
                              <span>Машинный:</span>
                              <span className="text-warning">
                                {stats.machine_translated}
                              </span>
                            </div>
                            <div className="flex justify-between">
                              <span>Проверено:</span>
                              <span className="text-info">
                                {stats.verified}
                              </span>
                            </div>
                            <div className="progress progress-primary w-full">
                              <div
                                className="progress-bar"
                                style={{ width: `${stats.coverage}%` }}
                              ></div>
                            </div>
                            <div className="text-center font-bold text-primary">
                              {stats.coverage}%
                            </div>
                          </div>
                        </div>
                      </div>
                    )
                  )}
                </div>
              </div>
            </div>

            {/* Recent Activity - Demo */}
            <div className="card bg-base-100 shadow-sm">
              <div className="card-body">
                <h3 className="card-title mb-4 flex items-center gap-2">
                  <ClockIcon className="h-5 w-5" />
                  Последняя активность
                </h3>
                <div className="space-y-3">
                  <div className="flex items-center justify-between p-3 bg-base-200 rounded-lg">
                    <div className="flex items-center gap-3">
                      <div className="badge badge-success">create</div>
                      <span>Создан перевод для category #15</span>
                    </div>
                    <span className="text-sm text-base-content/60">
                      2 минуты назад
                    </span>
                  </div>
                  <div className="flex items-center justify-between p-3 bg-base-200 rounded-lg">
                    <div className="flex items-center gap-3">
                      <div className="badge badge-info">translate</div>
                      <span>Массовый перевод 25 атрибутов</span>
                    </div>
                    <span className="text-sm text-base-content/60">
                      5 минут назад
                    </span>
                  </div>
                  <div className="flex items-center justify-between p-3 bg-base-200 rounded-lg">
                    <div className="flex items-center gap-3">
                      <div className="badge badge-warning">update</div>
                      <span>Обновлен перевод для listing #123</span>
                    </div>
                    <span className="text-sm text-base-content/60">
                      10 минут назад
                    </span>
                  </div>
                </div>

                <div className="mt-4">
                  <button
                    className="btn btn-outline btn-sm"
                    onClick={() => setActiveTab('audit')}
                  >
                    Посмотреть полный журнал
                  </button>
                </div>
              </div>
            </div>
          </div>
        )}

        {/* Bulk Translation Tab */}
        {activeTab === 'bulk' && <BulkTranslationManager />}

        {/* Audit Tab */}
        {activeTab === 'audit' && <AuditLogViewer />}

        {/* Export/Import Tab */}
        {activeTab === 'export' && <ExportImportManager />}

        {/* AI Translations Tab (Legacy) */}
        {activeTab === 'ai' && <AITranslationsDemo />}

        {/* Synchronization Tab (Legacy) */}
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
                      Система автоматически синхронизирует переводы между
                      Frontend JSON файлами и базой данных
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

        {/* Costs Tab */}
        {activeTab === 'costs' && <AICostsMonitor />}

        {/* Statistics Tab (Legacy) */}
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
                              <span className="text-info">
                                {stats.verified}
                              </span>
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
                    <span className="badge badge-warning">
                      148 плейсхолдеров
                    </span>
                  </div>
                  <div className="flex justify-between items-center p-3 bg-warning/10 rounded-lg">
                    <span className="font-medium">reviews</span>
                    <span className="badge badge-warning">
                      79 плейсхолдеров
                    </span>
                  </div>
                  <div className="flex justify-between items-center p-3 bg-warning/10 rounded-lg">
                    <span className="font-medium">orders</span>
                    <span className="badge badge-warning">
                      15 плейсхолдеров
                    </span>
                  </div>
                  <div className="flex justify-between items-center p-3 bg-warning/10 rounded-lg">
                    <span className="font-medium">search</span>
                    <span className="badge badge-warning">
                      15 плейсхолдеров
                    </span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        )}
      </div>

      {/* Version History Modal */}
      {showVersionHistory && versionHistoryParams && (
        <VersionHistoryViewer
          entityType={versionHistoryParams.entityType}
          entityId={versionHistoryParams.entityId}
          onClose={closeVersionHistory}
        />
      )}
    </div>
  );
}
