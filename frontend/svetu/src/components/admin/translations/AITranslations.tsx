'use client';

import { useState, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import { toast } from '@/utils/toast';
import { tokenManager } from '@/utils/tokenManager';

interface AIProvider {
  id: string;
  name: string;
  type: 'openai' | 'anthropic' | 'google' | 'deepl';
  apiKey: string;
  endpoint?: string;
  model?: string;
  enabled: boolean;
  maxTokens?: number;
  temperature?: number;
}

interface TranslationItem {
  key: string;
  module: string;
  sourceLanguage: string;
  sourceText: string;
  targetLanguages: string[];
  context?: string;
}

interface TranslationResult {
  key: string;
  module: string;
  translations: Record<string, string>;
  provider: string;
  confidence?: number;
  alternativeTranslations?: Record<string, string[]>;
}

interface AITranslationsProps {
  onTranslationComplete?: () => void;
}

const SUPPORTED_LANGUAGES = [
  { code: 'en', name: 'English' },
  { code: 'ru', name: 'Русский' },
  { code: 'sr', name: 'Српски' },
];

const AI_PROVIDERS: AIProvider[] = [
  {
    id: 'openai',
    name: 'OpenAI GPT-4',
    type: 'openai',
    apiKey: '',
    model: 'gpt-4-turbo-preview',
    enabled: false,
    maxTokens: 2000,
    temperature: 0.3,
  },
  {
    id: 'anthropic',
    name: 'Anthropic Claude',
    type: 'anthropic',
    apiKey: '',
    model: 'claude-3-opus-20240229',
    enabled: false,
    maxTokens: 2000,
    temperature: 0.3,
  },
  {
    id: 'deepl',
    name: 'DeepL API',
    type: 'deepl',
    apiKey: '',
    endpoint: 'https://api.deepl.com/v2/translate',
    enabled: false,
  },
  {
    id: 'google',
    name: 'Google Translate',
    type: 'google',
    apiKey: '',
    endpoint: 'https://translation.googleapis.com/language/translate/v2',
    enabled: false,
  },
];

export default function AITranslations({
  onTranslationComplete,
}: AITranslationsProps) {
  const t = useTranslations('admin');
  const [providers, setProviders] = useState<AIProvider[]>(AI_PROVIDERS);
  const [activeProvider, setActiveProvider] = useState<string>('');
  const [_translationItems, _setTranslationItems] = useState<TranslationItem[]>(
    []
  );
  const [isTranslating, setIsTranslating] = useState(false);
  const [results, setResults] = useState<TranslationResult[]>([]);
  const [_showProviderSettings, _setShowProviderSettings] = useState(false);
  const [editingProvider, setEditingProvider] = useState<AIProvider | null>(
    null
  );

  // Batch translation settings
  const [batchMode, setBatchMode] = useState(false);
  const [selectedModule, _setSelectedModule] = useState('');
  const [modules, setModules] = useState<string[]>([]);
  const [sourceLanguage, setSourceLanguage] = useState('en');
  const [targetLanguages, setTargetLanguages] = useState<string[]>([
    'ru',
    'sr',
  ]);
  const [missingOnly, setMissingOnly] = useState(true);

  // Single translation
  const [singleText, setSingleText] = useState('');
  const [singleKey, setSingleKey] = useState('');
  const [singleModule, setSingleModule] = useState('common');
  const [singleContext, setSingleContext] = useState('');

  useEffect(() => {
    fetchProviders();
    fetchModules();
  }, []);

  const fetchProviders = async () => {
    try {
      const response = await fetch('/api/v1/admin/translations/ai/providers', {
        headers: {
          Authorization: `Bearer ${tokenManager.getAccessToken()}`,
        },
      });

      if (response.ok) {
        const data = await response.json();
        if (data.data) {
          setProviders(data.data);
          const active = data.data.find((p: AIProvider) => p.enabled);
          if (active) {
            setActiveProvider(active.id);
          }
        }
      }
    } catch (error) {
      console.error('Error fetching providers:', error);
    }
  };

  const fetchModules = async () => {
    try {
      const response = await fetch(
        '/api/v1/admin/translations/frontend/modules',
        {
          headers: {
            Authorization: `Bearer ${tokenManager.getAccessToken()}`,
          },
        }
      );

      if (response.ok) {
        const data = await response.json();
        if (data.data) {
          setModules(data.data.map((m: any) => m.name));
        }
      }
    } catch (error) {
      console.error('Error fetching modules:', error);
    }
  };

  const handleProviderUpdate = async (provider: AIProvider) => {
    try {
      const response = await fetch(
        `/api/v1/admin/translations/ai/providers/${provider.id}`,
        {
          method: 'PUT',
          headers: {
            'Content-Type': 'application/json',
            Authorization: `Bearer ${tokenManager.getAccessToken()}`,
          },
          body: JSON.stringify(provider),
        }
      );

      if (response.ok) {
        toast.success('Провайдер обновлен');
        setProviders((prev) =>
          prev.map((p) => (p.id === provider.id ? provider : p))
        );
        if (provider.enabled) {
          setActiveProvider(provider.id);
          // Disable other providers
          setProviders((prev) =>
            prev.map((p) =>
              p.id !== provider.id ? { ...p, enabled: false } : p
            )
          );
        }
        setEditingProvider(null);
      } else {
        toast.error('Ошибка обновления провайдера');
      }
    } catch (error) {
      console.error('Error updating provider:', error);
      toast.error('Ошибка при обновлении провайдера');
    }
  };

  const handleSingleTranslate = async () => {
    if (!singleText || !singleKey) {
      toast.error('Введите текст и ключ для перевода');
      return;
    }

    if (!activeProvider) {
      toast.error('Выберите AI провайдера');
      return;
    }

    setIsTranslating(true);
    try {
      const response = await fetch('/api/v1/admin/translations/ai/translate', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${tokenManager.getAccessToken()}`,
        },
        body: JSON.stringify({
          provider: activeProvider,
          text: singleText,
          key: singleKey,
          module: singleModule,
          source_language: sourceLanguage,
          target_languages: targetLanguages,
          context: singleContext || undefined,
        }),
      });

      if (response.ok) {
        const data = await response.json();
        setResults([data.data]);
        toast.success('Перевод выполнен успешно');
      } else {
        toast.error('Ошибка при переводе');
      }
    } catch (error) {
      console.error('Error translating:', error);
      toast.error('Ошибка при выполнении перевода');
    } finally {
      setIsTranslating(false);
    }
  };

  const handleBatchTranslate = async () => {
    if (!selectedModule && !batchMode) {
      toast.error('Выберите модуль для перевода');
      return;
    }

    if (!activeProvider) {
      toast.error('Выберите AI провайдера');
      return;
    }

    setIsTranslating(true);
    try {
      const response = await fetch('/api/v1/admin/translations/ai/batch', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${tokenManager.getAccessToken()}`,
        },
        body: JSON.stringify({
          provider: activeProvider,
          modules: batchMode ? modules : [selectedModule],
          source_language: sourceLanguage,
          target_languages: targetLanguages,
          missing_only: missingOnly,
        }),
      });

      if (response.ok) {
        const data = await response.json();
        setResults(data.data.results || []);
        toast.success(`Переведено ${data.data.translated_count} текстов`);
        onTranslationComplete?.();
      } else {
        toast.error('Ошибка при массовом переводе');
      }
    } catch (error) {
      console.error('Error batch translating:', error);
      toast.error('Ошибка при выполнении массового перевода');
    } finally {
      setIsTranslating(false);
    }
  };

  const handleApplyTranslations = async () => {
    if (results.length === 0) {
      toast.error('Нет переводов для применения');
      return;
    }

    try {
      const updates = results.flatMap((result) =>
        Object.entries(result.translations).map(([lang, text]) => ({
          key: result.key,
          module: result.module,
          language: lang,
          value: text,
        }))
      );

      const response = await fetch('/api/v1/admin/translations/ai/apply', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${tokenManager.getAccessToken()}`,
        },
        body: JSON.stringify({ translations: updates }),
      });

      if (response.ok) {
        toast.success('Переводы применены успешно');
        setResults([]);
        onTranslationComplete?.();
      } else {
        toast.error('Ошибка при применении переводов');
      }
    } catch (error) {
      console.error('Error applying translations:', error);
      toast.error('Ошибка при применении переводов');
    }
  };

  const renderProviderSettings = () => (
    <div className="modal modal-open">
      <div className="modal-box max-w-2xl">
        <h3 className="font-bold text-lg mb-4">
          Настройки AI провайдера: {editingProvider?.name}
        </h3>

        {editingProvider && (
          <div className="space-y-4">
            <div className="form-control">
              <label className="label">
                <span className="label-text">API ключ</span>
              </label>
              <input
                type="password"
                value={editingProvider.apiKey}
                onChange={(e) =>
                  setEditingProvider({
                    ...editingProvider,
                    apiKey: e.target.value,
                  })
                }
                className="input input-bordered"
                placeholder="Введите API ключ..."
              />
            </div>

            {editingProvider.endpoint && (
              <div className="form-control">
                <label className="label">
                  <span className="label-text">Endpoint URL</span>
                </label>
                <input
                  type="text"
                  value={editingProvider.endpoint}
                  onChange={(e) =>
                    setEditingProvider({
                      ...editingProvider,
                      endpoint: e.target.value,
                    })
                  }
                  className="input input-bordered"
                />
              </div>
            )}

            {editingProvider.model && (
              <div className="form-control">
                <label className="label">
                  <span className="label-text">Модель</span>
                </label>
                <input
                  type="text"
                  value={editingProvider.model}
                  onChange={(e) =>
                    setEditingProvider({
                      ...editingProvider,
                      model: e.target.value,
                    })
                  }
                  className="input input-bordered"
                />
              </div>
            )}

            {editingProvider.maxTokens !== undefined && (
              <div className="form-control">
                <label className="label">
                  <span className="label-text">Max Tokens</span>
                </label>
                <input
                  type="number"
                  value={editingProvider.maxTokens}
                  onChange={(e) =>
                    setEditingProvider({
                      ...editingProvider,
                      maxTokens: parseInt(e.target.value),
                    })
                  }
                  className="input input-bordered"
                />
              </div>
            )}

            {editingProvider.temperature !== undefined && (
              <div className="form-control">
                <label className="label">
                  <span className="label-text">Temperature (0-1)</span>
                </label>
                <input
                  type="number"
                  step="0.1"
                  min="0"
                  max="1"
                  value={editingProvider.temperature}
                  onChange={(e) =>
                    setEditingProvider({
                      ...editingProvider,
                      temperature: parseFloat(e.target.value),
                    })
                  }
                  className="input input-bordered"
                />
              </div>
            )}

            <div className="form-control">
              <label className="label cursor-pointer">
                <span className="label-text">Активировать провайдера</span>
                <input
                  type="checkbox"
                  checked={editingProvider.enabled}
                  onChange={(e) =>
                    setEditingProvider({
                      ...editingProvider,
                      enabled: e.target.checked,
                    })
                  }
                  className="checkbox checkbox-primary"
                />
              </label>
            </div>
          </div>
        )}

        <div className="modal-action">
          <button
            onClick={() => setEditingProvider(null)}
            className="btn btn-ghost"
          >
            Отмена
          </button>
          <button
            onClick={() =>
              editingProvider && handleProviderUpdate(editingProvider)
            }
            className="btn btn-primary"
          >
            Сохранить
          </button>
        </div>
      </div>
    </div>
  );

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h3 className="text-xl font-semibold">AI-powered переводы</h3>
          <p className="text-sm text-base-content/70 mt-1">
            Автоматический перевод текстов с помощью AI
          </p>
        </div>
      </div>

      {/* Provider Selection */}
      <div className="card bg-base-100 shadow-sm">
        <div className="card-body">
          <h4 className="font-semibold mb-3">AI Провайдеры</h4>

          <div className="grid md:grid-cols-2 lg:grid-cols-4 gap-3">
            {providers.map((provider) => (
              <div
                key={provider.id}
                className={`card border ${
                  provider.enabled
                    ? 'border-primary bg-primary/5'
                    : 'border-base-300'
                }`}
              >
                <div className="card-body p-3">
                  <div className="flex justify-between items-start">
                    <div>
                      <h5 className="font-medium text-sm">{provider.name}</h5>
                      <div className="mt-1">
                        {provider.enabled ? (
                          <span className="badge badge-success badge-xs">
                            Активен
                          </span>
                        ) : (
                          <span className="badge badge-ghost badge-xs">
                            Неактивен
                          </span>
                        )}
                      </div>
                    </div>
                    <button
                      onClick={() => setEditingProvider(provider)}
                      className="btn btn-ghost btn-xs"
                    >
                      ⚙️
                    </button>
                  </div>

                  {provider.apiKey && (
                    <div className="text-xs text-success mt-2">
                      ✓ API ключ настроен
                    </div>
                  )}
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>

      {/* Translation Mode Tabs */}
      <div className="tabs tabs-boxed">
        <a
          className={`tab ${!batchMode ? 'tab-active' : ''}`}
          onClick={() => setBatchMode(false)}
        >
          Одиночный перевод
        </a>
        <a
          className={`tab ${batchMode ? 'tab-active' : ''}`}
          onClick={() => setBatchMode(true)}
        >
          Массовый перевод
        </a>
      </div>

      {/* Single Translation */}
      {!batchMode && (
        <div className="card bg-base-100 shadow-sm">
          <div className="card-body">
            <h4 className="font-semibold mb-3">Одиночный перевод</h4>

            <div className="grid md:grid-cols-2 gap-4">
              <div className="form-control">
                <label className="label">
                  <span className="label-text">Ключ перевода</span>
                </label>
                <input
                  type="text"
                  value={singleKey}
                  onChange={(e) => setSingleKey(e.target.value)}
                  className="input input-bordered"
                  placeholder="например: common.welcomeMessage"
                />
              </div>

              <div className="form-control">
                <label className="label">
                  <span className="label-text">Модуль</span>
                </label>
                <select
                  value={singleModule}
                  onChange={(e) => setSingleModule(e.target.value)}
                  className="select select-bordered"
                >
                  {modules.map((module) => (
                    <option key={module} value={module}>
                      {module}
                    </option>
                  ))}
                </select>
              </div>

              <div className="form-control md:col-span-2">
                <label className="label">
                  <span className="label-text">Текст для перевода</span>
                </label>
                <textarea
                  value={singleText}
                  onChange={(e) => setSingleText(e.target.value)}
                  className="textarea textarea-bordered"
                  rows={3}
                  placeholder="Введите текст на исходном языке..."
                />
              </div>

              <div className="form-control md:col-span-2">
                <label className="label">
                  <span className="label-text">Контекст (опционально)</span>
                </label>
                <input
                  type="text"
                  value={singleContext}
                  onChange={(e) => setSingleContext(e.target.value)}
                  className="input input-bordered"
                  placeholder="например: Заголовок на главной странице"
                />
              </div>

              <div className="form-control">
                <label className="label">
                  <span className="label-text">Исходный язык</span>
                </label>
                <select
                  value={sourceLanguage}
                  onChange={(e) => setSourceLanguage(e.target.value)}
                  className="select select-bordered"
                >
                  {SUPPORTED_LANGUAGES.map((lang) => (
                    <option key={lang.code} value={lang.code}>
                      {lang.name}
                    </option>
                  ))}
                </select>
              </div>

              <div className="form-control">
                <label className="label">
                  <span className="label-text">Целевые языки</span>
                </label>
                <div className="flex gap-3">
                  {SUPPORTED_LANGUAGES.filter(
                    (l) => l.code !== sourceLanguage
                  ).map((lang) => (
                    <label
                      key={lang.code}
                      className="label cursor-pointer gap-2"
                    >
                      <input
                        type="checkbox"
                        checked={targetLanguages.includes(lang.code)}
                        onChange={(e) => {
                          if (e.target.checked) {
                            setTargetLanguages([...targetLanguages, lang.code]);
                          } else {
                            setTargetLanguages(
                              targetLanguages.filter((l) => l !== lang.code)
                            );
                          }
                        }}
                        className="checkbox checkbox-primary checkbox-sm"
                      />
                      <span className="label-text">{lang.name}</span>
                    </label>
                  ))}
                </div>
              </div>
            </div>

            <div className="mt-4">
              <button
                onClick={handleSingleTranslate}
                disabled={isTranslating || !activeProvider}
                className="btn btn-primary"
              >
                {isTranslating ? (
                  <>
                    <span className="loading loading-spinner loading-sm"></span>
                    Перевод...
                  </>
                ) : (
                  '🤖 Перевести'
                )}
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Batch Translation */}
      {batchMode && (
        <div className="card bg-base-100 shadow-sm">
          <div className="card-body">
            <h4 className="font-semibold mb-3">Массовый перевод</h4>

            <div className="grid md:grid-cols-2 gap-4">
              <div className="form-control">
                <label className="label">
                  <span className="label-text">Выберите модули</span>
                </label>
                <select
                  multiple
                  value={modules}
                  onChange={(e) => {
                    const selected = Array.from(
                      e.target.selectedOptions,
                      (option) => option.value
                    );
                    setModules(selected);
                  }}
                  className="select select-bordered h-32"
                  size={5}
                >
                  {modules.map((module) => (
                    <option key={module} value={module}>
                      {module}
                    </option>
                  ))}
                </select>
                <label className="label">
                  <span className="label-text-alt">
                    Удерживайте Ctrl для множественного выбора
                  </span>
                </label>
              </div>

              <div className="space-y-4">
                <div className="form-control">
                  <label className="label">
                    <span className="label-text">Исходный язык</span>
                  </label>
                  <select
                    value={sourceLanguage}
                    onChange={(e) => setSourceLanguage(e.target.value)}
                    className="select select-bordered"
                  >
                    {SUPPORTED_LANGUAGES.map((lang) => (
                      <option key={lang.code} value={lang.code}>
                        {lang.name}
                      </option>
                    ))}
                  </select>
                </div>

                <div className="form-control">
                  <label className="label cursor-pointer">
                    <span className="label-text">
                      Только отсутствующие переводы
                    </span>
                    <input
                      type="checkbox"
                      checked={missingOnly}
                      onChange={(e) => setMissingOnly(e.target.checked)}
                      className="checkbox checkbox-primary"
                    />
                  </label>
                </div>
              </div>

              <div className="form-control md:col-span-2">
                <label className="label">
                  <span className="label-text">Целевые языки</span>
                </label>
                <div className="flex gap-4">
                  {SUPPORTED_LANGUAGES.filter(
                    (l) => l.code !== sourceLanguage
                  ).map((lang) => (
                    <label
                      key={lang.code}
                      className="label cursor-pointer gap-2"
                    >
                      <input
                        type="checkbox"
                        checked={targetLanguages.includes(lang.code)}
                        onChange={(e) => {
                          if (e.target.checked) {
                            setTargetLanguages([...targetLanguages, lang.code]);
                          } else {
                            setTargetLanguages(
                              targetLanguages.filter((l) => l !== lang.code)
                            );
                          }
                        }}
                        className="checkbox checkbox-primary"
                      />
                      <span className="label-text">{lang.name}</span>
                    </label>
                  ))}
                </div>
              </div>
            </div>

            <div className="mt-4">
              <button
                onClick={handleBatchTranslate}
                disabled={isTranslating || !activeProvider}
                className="btn btn-primary"
              >
                {isTranslating ? (
                  <>
                    <span className="loading loading-spinner loading-sm"></span>
                    Выполняется перевод...
                  </>
                ) : (
                  '🤖 Запустить массовый перевод'
                )}
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Translation Results */}
      {results.length > 0 && (
        <div className="card bg-base-100 shadow-sm">
          <div className="card-body">
            <div className="flex justify-between items-center mb-4">
              <h4 className="font-semibold">
                Результаты перевода ({results.length})
              </h4>
              <button
                onClick={handleApplyTranslations}
                className="btn btn-success btn-sm"
              >
                ✓ Применить все переводы
              </button>
            </div>

            <div className="space-y-3 max-h-96 overflow-y-auto">
              {results.map((result, idx) => (
                <div key={idx} className="card bg-base-200">
                  <div className="card-body p-4">
                    <div className="flex justify-between items-start">
                      <div>
                        <div className="font-mono text-sm font-semibold">
                          {result.key}
                        </div>
                        <div className="text-xs text-base-content/60">
                          Модуль: {result.module} | Провайдер: {result.provider}
                        </div>
                      </div>
                      {result.confidence && (
                        <div className="badge badge-ghost badge-sm">
                          Уверенность: {Math.round(result.confidence * 100)}%
                        </div>
                      )}
                    </div>

                    <div className="grid gap-2 mt-3">
                      {Object.entries(result.translations).map(
                        ([lang, text]) => (
                          <div key={lang} className="flex gap-2">
                            <span className="badge badge-primary badge-sm">
                              {lang.toUpperCase()}
                            </span>
                            <span className="text-sm">{text}</span>
                          </div>
                        )
                      )}
                    </div>

                    {result.alternativeTranslations && (
                      <details className="mt-2">
                        <summary className="text-xs cursor-pointer text-base-content/60">
                          Альтернативные варианты
                        </summary>
                        <div className="mt-2 space-y-1">
                          {Object.entries(result.alternativeTranslations).map(
                            ([lang, alternatives]) => (
                              <div key={lang} className="text-xs">
                                <span className="font-semibold">
                                  {lang.toUpperCase()}:
                                </span>
                                {alternatives.map((alt, i) => (
                                  <span key={i} className="ml-2">
                                    {i + 1}. {alt}
                                  </span>
                                ))}
                              </div>
                            )
                          )}
                        </div>
                      </details>
                    )}
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>
      )}

      {/* Provider Settings Modal */}
      {editingProvider && renderProviderSettings()}
    </div>
  );
}
