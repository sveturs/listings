'use client';

import { useState, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import { toast } from 'react-hot-toast';

interface SearchConfig {
  fuzzySearchEnabled: boolean;
  fuzzySearchThreshold: number;
  maxResults: number;
  defaultSort: 'relevance' | 'date' | 'price';
  enableSynonyms: boolean;
  enableTransliteration: boolean;
  cacheEnabled: boolean;
  cacheTimeout: number;
  analyticsEnabled: boolean;
  debugMode: boolean;
}

export default function SearchConfig() {
  const t = useTranslations();
  const [config, setConfig] = useState<SearchConfig>({
    fuzzySearchEnabled: true,
    fuzzySearchThreshold: 0.3,
    maxResults: 100,
    defaultSort: 'relevance',
    enableSynonyms: true,
    enableTransliteration: true,
    cacheEnabled: true,
    cacheTimeout: 300,
    analyticsEnabled: true,
    debugMode: false,
  });
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [originalConfig, setOriginalConfig] = useState<SearchConfig | null>(
    null
  );

  useEffect(() => {
    fetchConfig();
  }, []);

  const fetchConfig = async () => {
    try {
      setLoading(true);
      // Mock data for now since API may not be implemented
      const mockConfig: SearchConfig = {
        fuzzySearchEnabled: true,
        fuzzySearchThreshold: 0.3,
        maxResults: 100,
        defaultSort: 'relevance',
        enableSynonyms: true,
        enableTransliteration: true,
        cacheEnabled: true,
        cacheTimeout: 300,
        analyticsEnabled: true,
        debugMode: false,
      };
      setConfig(mockConfig);
      setOriginalConfig(mockConfig);
    } catch (error) {
      console.error('Error fetching search config:', error);
      toast.error(t('admin.search.config.fetchError'));
    } finally {
      setLoading(false);
    }
  };

  const handleSave = async () => {
    try {
      setSaving(true);
      // Mock save - in real implementation, this would call an API
      await new Promise((resolve) => setTimeout(resolve, 1000));
      setOriginalConfig(config);
      toast.success(t('admin.search.config.saved'));
    } catch (error) {
      console.error('Error saving search config:', error);
      toast.error(t('admin.search.config.saveError'));
    } finally {
      setSaving(false);
    }
  };

  const handleReset = () => {
    if (originalConfig) {
      setConfig(originalConfig);
      toast(t('admin.search.config.reset'));
    }
  };

  const hasChanges =
    originalConfig && JSON.stringify(config) !== JSON.stringify(originalConfig);

  if (loading) {
    return (
      <div className="flex justify-center items-center h-64">
        <div className="loading loading-spinner loading-lg"></div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="card bg-base-100 shadow-xl">
        <div className="card-body">
          <h2 className="card-title">{t('admin.search.config.general')}</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="form-control">
              <label className="label">
                <span className="label-text">
                  {t('admin.search.config.maxResults')}
                </span>
              </label>
              <input
                type="number"
                className="input input-bordered"
                value={config.maxResults}
                onChange={(e) =>
                  setConfig({
                    ...config,
                    maxResults: parseInt(e.target.value) || 100,
                  })
                }
                min="10"
                max="1000"
              />
            </div>
            <div className="form-control">
              <label className="label">
                <span className="label-text">
                  {t('admin.search.config.defaultSort')}
                </span>
              </label>
              <select
                className="select select-bordered"
                value={config.defaultSort}
                onChange={(e) =>
                  setConfig({
                    ...config,
                    defaultSort: e.target.value as
                      | 'relevance'
                      | 'date'
                      | 'price',
                  })
                }
              >
                <option value="relevance">
                  {t('admin.search.config.relevance')}
                </option>
                <option value="date">{t('admin.search.config.date')}</option>
                <option value="price">{t('admin.search.config.price')}</option>
              </select>
            </div>
          </div>
        </div>
      </div>

      <div className="card bg-base-100 shadow-xl">
        <div className="card-body">
          <h2 className="card-title">{t('admin.search.config.fuzzySearch')}</h2>
          <div className="form-control">
            <label className="label cursor-pointer">
              <span className="label-text">
                {t('admin.search.config.enableFuzzySearch')}
              </span>
              <input
                type="checkbox"
                className="toggle toggle-primary"
                checked={config.fuzzySearchEnabled}
                onChange={(e) =>
                  setConfig({ ...config, fuzzySearchEnabled: e.target.checked })
                }
              />
            </label>
          </div>
          {config.fuzzySearchEnabled && (
            <div className="form-control">
              <label className="label">
                <span className="label-text">
                  {t('admin.search.config.fuzzyThreshold')}
                </span>
                <span className="label-text-alt">
                  {config.fuzzySearchThreshold}
                </span>
              </label>
              <input
                type="range"
                className="range range-primary"
                min="0"
                max="1"
                step="0.1"
                value={config.fuzzySearchThreshold}
                onChange={(e) =>
                  setConfig({
                    ...config,
                    fuzzySearchThreshold: parseFloat(e.target.value),
                  })
                }
              />
              <div className="w-full flex justify-between text-xs px-2">
                <span>0.0</span>
                <span>0.5</span>
                <span>1.0</span>
              </div>
            </div>
          )}
        </div>
      </div>

      <div className="card bg-base-100 shadow-xl">
        <div className="card-body">
          <h2 className="card-title">{t('admin.search.config.features')}</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="form-control">
              <label className="label cursor-pointer">
                <span className="label-text">
                  {t('admin.search.config.enableSynonyms')}
                </span>
                <input
                  type="checkbox"
                  className="toggle toggle-primary"
                  checked={config.enableSynonyms}
                  onChange={(e) =>
                    setConfig({ ...config, enableSynonyms: e.target.checked })
                  }
                />
              </label>
            </div>
            <div className="form-control">
              <label className="label cursor-pointer">
                <span className="label-text">
                  {t('admin.search.config.enableTransliteration')}
                </span>
                <input
                  type="checkbox"
                  className="toggle toggle-primary"
                  checked={config.enableTransliteration}
                  onChange={(e) =>
                    setConfig({
                      ...config,
                      enableTransliteration: e.target.checked,
                    })
                  }
                />
              </label>
            </div>
            <div className="form-control">
              <label className="label cursor-pointer">
                <span className="label-text">
                  {t('admin.search.config.enableAnalytics')}
                </span>
                <input
                  type="checkbox"
                  className="toggle toggle-primary"
                  checked={config.analyticsEnabled}
                  onChange={(e) =>
                    setConfig({ ...config, analyticsEnabled: e.target.checked })
                  }
                />
              </label>
            </div>
            <div className="form-control">
              <label className="label cursor-pointer">
                <span className="label-text">
                  {t('admin.search.config.debugMode')}
                </span>
                <input
                  type="checkbox"
                  className="toggle toggle-warning"
                  checked={config.debugMode}
                  onChange={(e) =>
                    setConfig({ ...config, debugMode: e.target.checked })
                  }
                />
              </label>
            </div>
          </div>
        </div>
      </div>

      <div className="card bg-base-100 shadow-xl">
        <div className="card-body">
          <h2 className="card-title">{t('admin.search.config.caching')}</h2>
          <div className="form-control">
            <label className="label cursor-pointer">
              <span className="label-text">
                {t('admin.search.config.enableCache')}
              </span>
              <input
                type="checkbox"
                className="toggle toggle-primary"
                checked={config.cacheEnabled}
                onChange={(e) =>
                  setConfig({ ...config, cacheEnabled: e.target.checked })
                }
              />
            </label>
          </div>
          {config.cacheEnabled && (
            <div className="form-control">
              <label className="label">
                <span className="label-text">
                  {t('admin.search.config.cacheTimeout')}
                </span>
                <span className="label-text-alt">{config.cacheTimeout}s</span>
              </label>
              <input
                type="number"
                className="input input-bordered"
                value={config.cacheTimeout}
                onChange={(e) =>
                  setConfig({
                    ...config,
                    cacheTimeout: parseInt(e.target.value) || 300,
                  })
                }
                min="60"
                max="3600"
              />
            </div>
          )}
        </div>
      </div>

      <div className="card bg-base-100 shadow-xl">
        <div className="card-body">
          <div className="card-actions justify-end">
            <button
              className="btn btn-outline"
              onClick={handleReset}
              disabled={!hasChanges || saving}
            >
              {t('admin.search.config.reset')}
            </button>
            <button
              className={`btn btn-primary ${saving ? 'loading' : ''}`}
              onClick={handleSave}
              disabled={!hasChanges || saving}
            >
              {saving ? '' : t('admin.search.config.save')}
            </button>
          </div>
        </div>
      </div>

      {config.debugMode && (
        <div className="alert alert-warning">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            className="stroke-current shrink-0 h-6 w-6"
            fill="none"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth="2"
              d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.98-.833-2.75 0L3.098 16.5c-.77.833.192 2.5 1.732 2.5z"
            />
          </svg>
          <span>{t('admin.search.config.debugWarning')}</span>
        </div>
      )}
    </div>
  );
}
