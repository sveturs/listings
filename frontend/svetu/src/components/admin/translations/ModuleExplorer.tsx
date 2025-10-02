'use client';

import { useState, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import {
  FolderIcon,
  CheckCircleIcon,
  ExclamationTriangleIcon,
  MagnifyingGlassIcon,
} from '@heroicons/react/24/outline';
import { apiClient } from '@/services/api-client';

interface ModuleData {
  name: string;
  keys: number;
  complete: number;
  incomplete: number;
  placeholders: number;
  missing: number;
  languages: Record<
    string,
    {
      total: number;
      complete: number;
      incomplete: number;
      placeholders: number;
      missing: number;
    }
  >;
}

interface ModuleExplorerProps {
  onSelectModule: (module: string) => void;
  selectedModule: string | null;
}

export default function ModuleExplorer({
  onSelectModule,
  selectedModule,
}: ModuleExplorerProps) {
  const t = useTranslations('admin');
  const [modules, setModules] = useState<ModuleData[]>([]);
  const [loading, setLoading] = useState(false);
  const [searchTerm, setSearchTerm] = useState('');
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetchModules();
  }, []);

  const fetchModules = async () => {
    setLoading(true);
    setError(null);
    try {
      const response = await apiClient.get(
        '/admin/translations/frontend/modules'
      );

      if (!response.data) {
        throw new Error('Failed to fetch modules');
      }

      setModules(response.data.data || []);
    } catch (err) {
      setError((err as Error).message);
    } finally {
      setLoading(false);
    }
  };

  const filteredModules = modules.filter((module) =>
    module.name.toLowerCase().includes(searchTerm.toLowerCase())
  );

  const getModuleStatus = (module: ModuleData) => {
    if (module.placeholders > 0 || module.missing > 0) {
      return 'warning';
    }
    if (module.complete === module.keys * 3) {
      // 3 languages
      return 'success';
    }
    return 'info';
  };

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'success':
        return <CheckCircleIcon className="h-5 w-5 text-success" />;
      case 'warning':
        return <ExclamationTriangleIcon className="h-5 w-5 text-warning" />;
      default:
        return <FolderIcon className="h-5 w-5 text-info" />;
    }
  };

  const getCompletionPercentage = (module: ModuleData) => {
    const total = module.keys * 3; // 3 languages
    if (total === 0) return 100;
    return Math.round((module.complete / total) * 100);
  };

  return (
    <div className="bg-base-100 rounded-lg p-4">
      <h3 className="text-lg font-semibold mb-4">
        {t('translations.modules')}
      </h3>

      {/* Search */}
      <div className="form-control mb-4">
        <div className="input-group">
          <span>
            <MagnifyingGlassIcon className="h-5 w-5" />
          </span>
          <input
            type="text"
            placeholder={t('translations.searchModules')}
            className="input input-bordered w-full"
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
          />
        </div>
      </div>

      {/* Module List */}
      {loading ? (
        <div className="flex justify-center py-8">
          <span className="loading loading-spinner"></span>
        </div>
      ) : error ? (
        <div className="alert alert-error">
          <span>{error}</span>
        </div>
      ) : (
        <div className="space-y-2 max-h-[600px] overflow-y-auto">
          {filteredModules.map((module) => {
            const status = getModuleStatus(module);
            const completion = getCompletionPercentage(module);
            const isSelected = selectedModule === module.name;

            return (
              <button
                key={module.name}
                className={`w-full text-left p-3 rounded-lg transition-all ${
                  isSelected
                    ? 'bg-primary text-primary-content'
                    : 'bg-base-200 hover:bg-base-300'
                }`}
                onClick={() => onSelectModule(module.name)}
              >
                <div className="flex items-start justify-between">
                  <div className="flex items-start space-x-2">
                    {getStatusIcon(status)}
                    <div>
                      <div className="font-medium">{module.name}</div>
                      <div
                        className={`text-sm ${isSelected ? 'text-primary-content/70' : 'text-base-content/60'}`}
                      >
                        {module.keys} {t('translations.keys')}
                      </div>
                    </div>
                  </div>
                  <div className="text-right">
                    <div
                      className={`text-sm font-medium ${isSelected ? 'text-primary-content' : ''}`}
                    >
                      {completion}%
                    </div>
                    {module.placeholders > 0 && (
                      <div
                        className={`text-xs ${isSelected ? 'text-primary-content/70' : 'text-warning'}`}
                      >
                        {module.placeholders} {t('translations.placeholders')}
                      </div>
                    )}
                  </div>
                </div>

                {/* Progress Bar */}
                <div className="mt-2">
                  <progress
                    className={`progress ${
                      status === 'success'
                        ? 'progress-success'
                        : status === 'warning'
                          ? 'progress-warning'
                          : 'progress-info'
                    } h-1`}
                    value={completion}
                    max="100"
                  />
                </div>
              </button>
            );
          })}
        </div>
      )}
    </div>
  );
}
