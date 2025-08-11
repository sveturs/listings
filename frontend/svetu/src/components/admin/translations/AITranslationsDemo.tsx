'use client';

import { useState } from 'react';
import { useTranslations } from 'next-intl';

const DEMO_PROVIDERS = [
  {
    id: 'openai',
    name: 'OpenAI GPT-4',
    type: 'openai',
    model: 'gpt-4-turbo-preview',
    enabled: true,
    configured: true,
  },
  {
    id: 'anthropic',
    name: 'Anthropic Claude',
    type: 'anthropic',
    model: 'claude-3-opus',
    enabled: false,
    configured: false,
  },
  {
    id: 'deepl',
    name: 'DeepL API',
    type: 'deepl',
    enabled: false,
    configured: false,
  },
  {
    id: 'google',
    name: 'Google Translate',
    type: 'google',
    enabled: false,
    configured: true,
  },
];

const DEMO_MODULES = ['common', 'marketplace', 'auth', 'admin', 'orders'];

export default function AITranslationsDemo() {
  const _t = useTranslations('admin');
  const [batchMode, setBatchMode] = useState(false);
  const [activeProvider, setActiveProvider] = useState('openai');
  const [singleText, setSingleText] = useState('');
  const [singleKey, setSingleKey] = useState('');
  const [selectedModules, setSelectedModules] = useState<string[]>(['common']);
  const [demoResults, setDemoResults] = useState<any[]>([]);
  const [isTranslating, setIsTranslating] = useState(false);

  const handleDemoTranslate = () => {
    setIsTranslating(true);
    setTimeout(() => {
      setDemoResults([
        {
          key: singleKey || 'demo.example',
          module: 'common',
          translations: {
            ru: '[RU] ' + (singleText || 'Пример текста'),
            sr: '[SR] ' + (singleText || 'Primer teksta'),
            en: '[EN] ' + (singleText || 'Example text'),
          },
          provider: activeProvider,
          confidence: 0.95,
        },
      ]);
      setIsTranslating(false);
    }, 1500);
  };

  const handleBatchTranslate = () => {
    setIsTranslating(true);
    setTimeout(() => {
      const results = selectedModules.flatMap((module) => [
        {
          key: `${module}.title`,
          module,
          translations: {
            ru: `[RU] Заголовок модуля ${module}`,
            sr: `[SR] Naslov modula ${module}`,
            en: `[EN] Module title ${module}`,
          },
          provider: activeProvider,
          confidence: 0.92,
        },
        {
          key: `${module}.description`,
          module,
          translations: {
            ru: `[RU] Описание модуля ${module}`,
            sr: `[SR] Opis modula ${module}`,
            en: `[EN] Module description ${module}`,
          },
          provider: activeProvider,
          confidence: 0.88,
        },
      ]);
      setDemoResults(results);
      setIsTranslating(false);
    }, 2000);
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h3 className="text-xl font-semibold">AI-powered переводы (Демо)</h3>
          <p className="text-sm text-base-content/70 mt-1">
            Демонстрация автоматического перевода с помощью AI
          </p>
        </div>
      </div>

      {/* Provider Selection */}
      <div className="card bg-base-100 shadow-sm">
        <div className="card-body">
          <h4 className="font-semibold mb-3">AI Провайдеры</h4>

          <div className="grid md:grid-cols-2 lg:grid-cols-4 gap-3">
            {DEMO_PROVIDERS.map((provider) => (
              <div
                key={provider.id}
                className={`card border cursor-pointer transition-all ${
                  provider.id === activeProvider
                    ? 'border-primary bg-primary/5 shadow-md'
                    : 'border-base-300 hover:border-primary/50'
                }`}
                onClick={() =>
                  provider.enabled && setActiveProvider(provider.id)
                }
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
                    {provider.id === activeProvider && (
                      <span className="text-primary">✓</span>
                    )}
                  </div>

                  {provider.configured && (
                    <div className="text-xs text-success mt-2">✓ Настроен</div>
                  )}
                  {provider.model && (
                    <div className="text-xs text-base-content/60 mt-1">
                      {provider.model}
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

            <div className="space-y-4">
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

              <div className="flex gap-3">
                <div className="badge badge-outline">Исходный: EN</div>
                <div className="badge badge-primary">→ RU</div>
                <div className="badge badge-primary">→ SR</div>
              </div>
            </div>

            <div className="mt-4">
              <button
                onClick={handleDemoTranslate}
                disabled={isTranslating}
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

            <div className="space-y-4">
              <div className="form-control">
                <label className="label">
                  <span className="label-text">Выберите модули</span>
                </label>
                <div className="flex flex-wrap gap-2">
                  {DEMO_MODULES.map((module) => (
                    <label key={module} className="label cursor-pointer gap-2">
                      <input
                        type="checkbox"
                        checked={selectedModules.includes(module)}
                        onChange={(e) => {
                          if (e.target.checked) {
                            setSelectedModules([...selectedModules, module]);
                          } else {
                            setSelectedModules(
                              selectedModules.filter((m) => m !== module)
                            );
                          }
                        }}
                        className="checkbox checkbox-primary checkbox-sm"
                      />
                      <span className="label-text">{module}</span>
                    </label>
                  ))}
                </div>
              </div>

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
                  Будут переведены все недостающие тексты в выбранных модулях
                </span>
              </div>
            </div>

            <div className="mt-4">
              <button
                onClick={handleBatchTranslate}
                disabled={isTranslating || selectedModules.length === 0}
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
      {demoResults.length > 0 && (
        <div className="card bg-base-100 shadow-sm">
          <div className="card-body">
            <div className="flex justify-between items-center mb-4">
              <h4 className="font-semibold">
                Результаты перевода ({demoResults.length})
              </h4>
              <button className="btn btn-success btn-sm">
                ✓ Применить все переводы
              </button>
            </div>

            <div className="space-y-3 max-h-96 overflow-y-auto">
              {demoResults.map((result, idx) => (
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
                            <span className="text-sm">{text as string}</span>
                          </div>
                        )
                      )}
                    </div>
                  </div>
                </div>
              ))}
            </div>

            <div className="alert alert-success mt-4">
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
                  d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
                />
              </svg>
              <span>
                Демо-перевод выполнен успешно! В реальной системе переводы будут
                сохранены.
              </span>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
