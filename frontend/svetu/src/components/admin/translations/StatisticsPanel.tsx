'use client';

import { useTranslations } from 'next-intl';
import { ArrowPathIcon } from '@heroicons/react/24/outline';

interface StatisticsPanelProps {
  statistics: any;
  onRefresh: () => void;
}

export default function StatisticsPanel({
  statistics,
  onRefresh,
}: StatisticsPanelProps) {
  const t = useTranslations('admin');

  if (!statistics) {
    return (
      <div className="bg-base-100 rounded-lg p-8">
        <div className="text-center">
          <p className="text-base-content/60">
            {t('translations.noStatistics')}
          </p>
          <button className="btn btn-primary mt-4" onClick={onRefresh}>
            <ArrowPathIcon className="h-5 w-5 mr-2" />
            {t('translations.loadStatistics')}
          </button>
        </div>
      </div>
    );
  }

  const languages = Object.keys(statistics.language_stats || {});

  return (
    <div className="space-y-6">
      {/* Language Statistics */}
      <div className="bg-base-100 rounded-lg p-6">
        <h3 className="text-lg font-semibold mb-4">
          {t('translations.languageStatistics')}
        </h3>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          {languages.map((lang) => {
            const stats = statistics.language_stats[lang];
            return (
              <div key={lang} className="card bg-base-200">
                <div className="card-body">
                  <h4 className="card-title">{lang.toUpperCase()}</h4>
                  <div className="space-y-2">
                    <div className="flex justify-between">
                      <span>{t('translations.total')}:</span>
                      <span className="font-semibold">{stats.total}</span>
                    </div>
                    <div className="flex justify-between">
                      <span>{t('translations.complete')}:</span>
                      <span className="font-semibold text-success">
                        {stats.complete}
                      </span>
                    </div>
                    <div className="flex justify-between">
                      <span>{t('translations.machineTranslated')}:</span>
                      <span className="font-semibold text-warning">
                        {stats.machine_translated}
                      </span>
                    </div>
                    <div className="flex justify-between">
                      <span>{t('translations.verified')}:</span>
                      <span className="font-semibold text-info">
                        {stats.verified}
                      </span>
                    </div>
                  </div>
                  <div className="mt-4">
                    <div className="flex justify-between text-sm mb-1">
                      <span>{t('translations.coverage')}</span>
                      <span>{Math.round(stats.coverage)}%</span>
                    </div>
                    <progress
                      className="progress progress-primary"
                      value={stats.coverage}
                      max="100"
                    />
                  </div>
                </div>
              </div>
            );
          })}
        </div>
      </div>

      {/* Module Statistics */}
      <div className="bg-base-100 rounded-lg p-6">
        <h3 className="text-lg font-semibold mb-4">
          {t('translations.moduleStatistics')}
        </h3>
        <div className="overflow-x-auto">
          <table className="table">
            <thead>
              <tr>
                <th>{t('translations.module')}</th>
                <th>{t('translations.keys')}</th>
                <th>{t('translations.complete')}</th>
                <th>{t('translations.incomplete')}</th>
                <th>{t('translations.placeholders')}</th>
                <th>{t('translations.missing')}</th>
                <th>{t('translations.completion')}</th>
              </tr>
            </thead>
            <tbody>
              {statistics.module_stats?.map((module: any) => {
                const total = module.keys * 3; // 3 languages
                const completion =
                  total > 0 ? Math.round((module.complete / total) * 100) : 0;
                return (
                  <tr key={module.name}>
                    <td className="font-medium">{module.name}</td>
                    <td>{module.keys}</td>
                    <td className="text-success">{module.complete}</td>
                    <td className="text-warning">{module.incomplete}</td>
                    <td className="text-error">{module.placeholders}</td>
                    <td className="text-base-content/60">{module.missing}</td>
                    <td>
                      <div className="flex items-center gap-2">
                        <progress
                          className={`progress ${
                            completion === 100
                              ? 'progress-success'
                              : completion > 80
                                ? 'progress-warning'
                                : 'progress-error'
                          } w-20`}
                          value={completion}
                          max="100"
                        />
                        <span className="text-sm">{completion}%</span>
                      </div>
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>
      </div>

      {/* Recent Changes */}
      {statistics.recent_changes && statistics.recent_changes.length > 0 && (
        <div className="bg-base-100 rounded-lg p-6">
          <h3 className="text-lg font-semibold mb-4">
            {t('translations.recentChanges')}
          </h3>
          <div className="space-y-2">
            {statistics.recent_changes.map((change: any, index: number) => (
              <div
                key={index}
                className="flex items-center justify-between py-2 border-b"
              >
                <div>
                  <span className="font-medium">{change.action}</span>
                  {change.entity_type && (
                    <span className="ml-2 text-sm text-base-content/60">
                      ({change.entity_type})
                    </span>
                  )}
                </div>
                <span className="text-sm text-base-content/60">
                  {new Date(change.created_at).toLocaleString()}
                </span>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Refresh Button */}
      <div className="flex justify-end">
        <button className="btn btn-primary" onClick={onRefresh}>
          <ArrowPathIcon className="h-5 w-5 mr-2" />
          {t('translations.refreshStatistics')}
        </button>
      </div>
    </div>
  );
}
