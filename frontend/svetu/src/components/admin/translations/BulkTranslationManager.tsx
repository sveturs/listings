'use client';

import { useState, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import {
  translationAdminApi,
  BulkTranslateRequest,
  BulkTranslateResult,
} from '@/services/translationAdminApi';
import { adminApi } from '@/services/admin';
import {
  PlayIcon,
  ChartBarIcon,
  DocumentArrowDownIcon,
  DocumentArrowUpIcon,
  CogIcon,
  UserIcon,
  CalendarIcon,
  TagIcon,
  CurrencyEuroIcon,
  ArrowPathIcon,
  XMarkIcon,
} from '@heroicons/react/24/outline';

interface EntityOption {
  id: number;
  name: string;
  type: 'category' | 'attribute' | 'listing';
  description?: string;
  icon?: string;
  parent_name?: string;
  price?: number;
  status?: string;
  user_name?: string;
  created_at?: string;
}

export default function BulkTranslationManager() {
  const _t = useTranslations('admin.translations');

  // State management
  const [isLoading, setIsLoading] = useState(false);
  const [entities, setEntities] = useState<EntityOption[]>([]);
  const [selectedEntities, setSelectedEntities] = useState<number[]>([]);
  const [entityType, setEntityType] = useState<
    'category' | 'attribute' | 'listing'
  >('category');
  const [sourceLanguage, setSourceLanguage] = useState('sr');
  const [targetLanguages, setTargetLanguages] = useState<string[]>([
    'en',
    'ru',
  ]);
  const [autoApprove, setAutoApprove] = useState(false);
  const [overwriteExisting, setOverwriteExisting] = useState(false);
  const [providerId, setProviderId] = useState<number | undefined>(undefined);
  const [providers, setProviders] = useState<any[]>([]);
  const [listingFilters, setListingFilters] = useState({
    categoryId: '',
    userId: '',
    onlyActive: true,
  });

  // Results state
  const [result, setResult] = useState<BulkTranslateResult | null>(null);
  const [progress, setProgress] = useState<number>(0);
  const [expandedSections, setExpandedSections] = useState({
    successful: false,
    failed: true,
    skipped: false,
  });

  const LANGUAGES = [
    { code: 'sr', name: 'Srpski', flag: '🇷🇸' },
    { code: 'en', name: 'English', flag: '🇺🇸' },
    { code: 'ru', name: 'Русский', flag: '🇷🇺' },
  ];

  // Load data on mount
  useEffect(() => {
    loadEntities();
    loadProviders();
  }, [entityType, listingFilters]);

  const loadEntities = async () => {
    try {
      setIsLoading(true);
      let data: EntityOption[] = [];

      if (entityType === 'category') {
        const categories = await adminApi.categories.getAll();
        data = categories.map((cat) => ({
          id: cat.id,
          name: cat.name,
          type: 'category' as const,
          icon: cat.icon || '📁',
          parent_name: cat.parent_name,
          description:
            cat.description || `${cat.listing_count || 0} объявлений`,
        }));
      } else if (entityType === 'attribute') {
        const response = await adminApi.attributes.getAll(1, 1000);
        data = response.data.map((attr) => ({
          id: attr.id,
          name: attr.display_name || attr.name,
          type: 'attribute' as const,
          icon: attr.icon || '⚙️',
          description: `Тип: ${attr.attribute_type}${attr.is_required ? ' • Обязательный' : ''}${attr.is_filterable ? ' • Фильтр' : ''}`,
        }));
      } else if (entityType === 'listing') {
        // Загружаем объявления для перевода
        try {
          // Строим параметры запроса
          const params = new URLSearchParams();
          params.append('limit', '100');
          if (listingFilters.categoryId) {
            params.append('category_id', listingFilters.categoryId);
          }
          if (listingFilters.userId) {
            params.append('user_id', listingFilters.userId);
          }
          if (listingFilters.onlyActive) {
            params.append('status', 'active');
          }

          const response = await fetch(
            `/api/v1/marketplace/listings?${params.toString()}`,
            {
              method: 'GET',
              headers: {
                'Content-Type': 'application/json',
              },
              credentials: 'include',
            }
          );

          if (response.ok) {
            const result = await response.json();
            console.log('Listings API response:', result);

            // API возвращает вложенную структуру: { data: { success: true, data: [...] } }
            let listings = [];

            // Проверяем вложенную структуру
            if (
              result.data &&
              result.data.data &&
              Array.isArray(result.data.data)
            ) {
              listings = result.data.data;
            } else if (result.data && Array.isArray(result.data)) {
              listings = result.data;
            } else if (result.listings && Array.isArray(result.listings)) {
              listings = result.listings;
            } else if (Array.isArray(result)) {
              listings = result;
            } else {
              console.error('Unexpected listings response structure:', result);
            }

            data = listings.map((listing: any) => ({
              id: listing.id,
              name: listing.title || `Объявление #${listing.id}`,
              type: 'listing' as const,
              price: listing.price,
              status: listing.status,
              user_name: listing.user?.name || listing.user_name,
              description: listing.description
                ? listing.description.length > 100
                  ? listing.description.substring(0, 100) + '...'
                  : listing.description
                : listing.category_name
                  ? `Категория: ${listing.category_name}`
                  : '',
              created_at: listing.created_at,
              icon:
                listing.images && listing.images.length > 0
                  ? listing.images[0].public_url || listing.images[0].url
                  : listing.image_url || null,
              parent_name: listing.category_name,
            }));
          }
        } catch (err) {
          console.error('Failed to load listings:', err);
        }
      }

      setEntities(data);
      // Сбрасываем выбранные элементы при смене типа
      setSelectedEntities([]);
    } catch (error) {
      console.error('Failed to load entities:', error);
    } finally {
      setIsLoading(false);
    }
  };

  const loadProviders = async () => {
    try {
      const response = await translationAdminApi.getProviders();
      if (response.success && response.data) {
        setProviders(response.data);

        // Select first active provider by default
        const activeProvider = response.data.find((p) => p.is_active);
        if (activeProvider) {
          setProviderId(activeProvider.id);
        }
      } else {
        setProviders([]);
        setProviderId(undefined);
      }
    } catch (error) {
      console.error('Failed to load providers:', error);
      // Providers API не реализован, используем значения по умолчанию
      setProviders([]);
      setProviderId(undefined);
    }
  };

  const handleEntitySelect = (entityId: number) => {
    setSelectedEntities((prev) =>
      prev.includes(entityId)
        ? prev.filter((id) => id !== entityId)
        : [...prev, entityId]
    );
  };

  const handleSelectAll = () => {
    if (selectedEntities.length === entities.length) {
      setSelectedEntities([]);
    } else {
      setSelectedEntities(entities.map((e) => e.id));
    }
  };

  const handleTargetLanguageToggle = (langCode: string) => {
    setTargetLanguages((prev) =>
      prev.includes(langCode)
        ? prev.filter((code) => code !== langCode)
        : [...prev, langCode]
    );
  };

  const retryFailedTranslations = async () => {
    if (
      !result?.details?.failed_items ||
      result.details.failed_items.length === 0
    ) {
      return;
    }

    // Собираем ID неудачных элементов
    const failedIds = result.details.failed_items.map((item) => item.entity_id);

    // Запускаем перевод только для неудачных элементов
    setSelectedEntities(failedIds);
    await startBulkTranslation();
  };

  const toggleSection = (section: 'successful' | 'failed' | 'skipped') => {
    setExpandedSections((prev) => ({
      ...prev,
      [section]: !prev[section],
    }));
  };

  const startBulkTranslation = async () => {
    if (selectedEntities.length === 0 || targetLanguages.length === 0) {
      alert('Выберите сущности и целевые языки');
      return;
    }

    try {
      setIsLoading(true);
      setProgress(0);

      const request: BulkTranslateRequest = {
        entity_type: entityType,
        entity_ids: selectedEntities,
        source_language: sourceLanguage,
        target_languages: targetLanguages,
        provider_id: providerId,
        auto_approve: autoApprove,
        overwrite_existing: overwriteExisting,
      };

      // Simulate progress (в реальности можно использовать WebSocket или polling)
      const progressInterval = setInterval(() => {
        setProgress((prev) => Math.min(prev + 10, 90));
      }, 500);

      const translationResult =
        await translationAdminApi.bulkTranslate(request);

      clearInterval(progressInterval);
      setProgress(100);

      if (translationResult.success && translationResult.data) {
        setResult(translationResult.data);
        
        // Обновляем список сущностей после успешного перевода
        // чтобы отразить новые переводы в UI
        await loadEntities();
        
        // Сбрасываем выбранные элементы
        setSelectedEntities([]);
      } else {
        // Показываем ошибку если запрос не успешен
        setResult({
          total_processed: selectedEntities.length,
          successful: 0,
          failed: selectedEntities.length,
          skipped: 0,
          errors: [
            translationResult.error || 'Неизвестная ошибка при переводе',
          ],
        });
      }

      setTimeout(() => setProgress(0), 2000);
    } catch (error) {
      console.error('Bulk translation failed:', error);
      alert('Ошибка при выполнении массового перевода');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold">Массовый перевод</h2>
          <p className="text-base-content/60 mt-1">
            Переводите множество сущностей одновременно с помощью AI
          </p>
        </div>
        <div className="flex gap-2">
          <button className="btn btn-outline btn-sm gap-2">
            <DocumentArrowDownIcon className="h-4 w-4" />
            Экспорт
          </button>
          <button className="btn btn-outline btn-sm gap-2">
            <DocumentArrowUpIcon className="h-4 w-4" />
            Импорт
          </button>
        </div>
      </div>

      {/* Configuration */}
      <div className="grid md:grid-cols-2 gap-6">
        {/* Left Column - Settings */}
        <div className="space-y-4">
          {/* Entity Type Selection */}
          <div className="card bg-base-100 shadow-sm">
            <div className="card-body">
              <h3 className="card-title text-lg">Тип сущности</h3>
              <div className="form-control">
                <label className="label cursor-pointer">
                  <span className="label-text">Категории</span>
                  <input
                    type="radio"
                    name="entity-type"
                    className="radio"
                    checked={entityType === 'category'}
                    onChange={() => setEntityType('category')}
                  />
                </label>
                <label className="label cursor-pointer">
                  <span className="label-text">Атрибуты</span>
                  <input
                    type="radio"
                    name="entity-type"
                    className="radio"
                    checked={entityType === 'attribute'}
                    onChange={() => setEntityType('attribute')}
                  />
                </label>
                <label className="label cursor-pointer">
                  <span className="label-text">Объявления</span>
                  <input
                    type="radio"
                    name="entity-type"
                    className="radio"
                    checked={entityType === 'listing'}
                    onChange={() => setEntityType('listing')}
                  />
                </label>
              </div>
            </div>
          </div>

          {/* Language Settings */}
          <div className="card bg-base-100 shadow-sm">
            <div className="card-body">
              <h3 className="card-title text-lg">Языки</h3>

              {/* Source Language */}
              <div className="form-control">
                <label className="label">
                  <span className="label-text">Исходный язык</span>
                </label>
                <select
                  className="select select-bordered"
                  value={sourceLanguage}
                  onChange={(e) => setSourceLanguage(e.target.value)}
                >
                  {LANGUAGES.map((lang) => (
                    <option key={lang.code} value={lang.code}>
                      {lang.flag} {lang.name}
                    </option>
                  ))}
                </select>
              </div>

              {/* Target Languages */}
              <div className="form-control">
                <label className="label">
                  <span className="label-text">Целевые языки</span>
                </label>
                <div className="space-y-2">
                  {LANGUAGES.filter((lang) => lang.code !== sourceLanguage).map(
                    (lang) => (
                      <label key={lang.code} className="label cursor-pointer">
                        <span className="label-text">
                          {lang.flag} {lang.name}
                        </span>
                        <input
                          type="checkbox"
                          className="checkbox"
                          checked={targetLanguages.includes(lang.code)}
                          onChange={() => handleTargetLanguageToggle(lang.code)}
                        />
                      </label>
                    )
                  )}
                </div>
              </div>
            </div>
          </div>

          {/* Provider Selection */}
          {providers.length > 0 && (
            <div className="card bg-base-100 shadow-sm">
              <div className="card-body">
                <h3 className="card-title text-lg">Провайдер перевода</h3>
                <select
                  className="select select-bordered"
                  value={providerId || ''}
                  onChange={(e) =>
                    setProviderId(
                      e.target.value ? Number(e.target.value) : undefined
                    )
                  }
                >
                  <option value="">Автоматический выбор</option>
                  {providers
                    .filter((p) => p.is_active)
                    .map((provider) => (
                      <option key={provider.id} value={provider.id}>
                        {provider.name} ({provider.type})
                      </option>
                    ))}
                </select>
              </div>
            </div>
          )}

          {/* Listing Filters - только для объявлений */}
          {entityType === 'listing' && (
            <div className="card bg-base-100 shadow-sm">
              <div className="card-body">
                <h3 className="card-title text-lg">Фильтры объявлений</h3>
                <div className="form-control">
                  <label className="label">
                    <span className="label-text">
                      ID категории (опционально)
                    </span>
                  </label>
                  <input
                    type="number"
                    placeholder="Например: 1"
                    className="input input-bordered input-sm"
                    value={listingFilters.categoryId}
                    onChange={(e) =>
                      setListingFilters((prev) => ({
                        ...prev,
                        categoryId: e.target.value,
                      }))
                    }
                  />
                </div>
                <div className="form-control">
                  <label className="label">
                    <span className="label-text">
                      ID пользователя (опционально)
                    </span>
                  </label>
                  <input
                    type="number"
                    placeholder="Например: 2"
                    className="input input-bordered input-sm"
                    value={listingFilters.userId}
                    onChange={(e) =>
                      setListingFilters((prev) => ({
                        ...prev,
                        userId: e.target.value,
                      }))
                    }
                  />
                </div>
                <div className="form-control">
                  <label className="label cursor-pointer">
                    <span className="label-text">Только активные</span>
                    <input
                      type="checkbox"
                      className="checkbox"
                      checked={listingFilters.onlyActive}
                      onChange={(e) =>
                        setListingFilters((prev) => ({
                          ...prev,
                          onlyActive: e.target.checked,
                        }))
                      }
                    />
                  </label>
                </div>
              </div>
            </div>
          )}

          {/* Options */}
          <div className="card bg-base-100 shadow-sm">
            <div className="card-body">
              <h3 className="card-title text-lg">Опции</h3>
              <div className="form-control">
                <label className="label cursor-pointer">
                  <span className="label-text">Автоматическое одобрение</span>
                  <input
                    type="checkbox"
                    className="toggle"
                    checked={autoApprove}
                    onChange={(e) => setAutoApprove(e.target.checked)}
                  />
                </label>
                <label className="label cursor-pointer">
                  <span className="label-text">Перезаписать существующие</span>
                  <input
                    type="checkbox"
                    className="toggle"
                    checked={overwriteExisting}
                    onChange={(e) => setOverwriteExisting(e.target.checked)}
                  />
                </label>
              </div>
            </div>
          </div>
        </div>

        {/* Right Column - Entity Selection */}
        <div className="space-y-4">
          <div className="card bg-base-100 shadow-sm">
            <div className="card-body">
              <div className="flex items-center justify-between mb-4">
                <h3 className="card-title text-lg">
                  Выбор{' '}
                  {entityType === 'category'
                    ? 'категорий'
                    : entityType === 'attribute'
                      ? 'атрибутов'
                      : 'объявлений'}
                </h3>
                <button
                  className="btn btn-ghost btn-sm"
                  onClick={handleSelectAll}
                >
                  {selectedEntities.length === entities.length
                    ? 'Снять все'
                    : 'Выбрать все'}
                </button>
              </div>

              <div className="max-h-[500px] overflow-y-auto space-y-2 pr-2">
                {isLoading ? (
                  <div className="flex justify-center py-8">
                    <span className="loading loading-spinner loading-md"></span>
                  </div>
                ) : entities.length === 0 ? (
                  <div className="text-center py-8 text-base-content/60">
                    <div className="text-6xl mb-4 opacity-20">
                      {entityType === 'category'
                        ? '📁'
                        : entityType === 'attribute'
                          ? '⚙️'
                          : '📋'}
                    </div>
                    Нет доступных{' '}
                    {entityType === 'category'
                      ? 'категорий'
                      : entityType === 'attribute'
                        ? 'атрибутов'
                        : 'объявлений'}
                  </div>
                ) : (
                  <div className="grid gap-2">
                    {entities.map((entity) => (
                      <div
                        key={entity.id}
                        className={`
                          relative border rounded-lg p-3 transition-all cursor-pointer
                          ${
                            selectedEntities.includes(entity.id)
                              ? 'border-primary bg-primary/5 shadow-md ring-2 ring-primary/20'
                              : 'border-base-300 hover:border-primary/50 hover:bg-base-200/50 hover:shadow-sm'
                          }
                        `}
                        onClick={() => handleEntitySelect(entity.id)}
                      >
                        {/* Selection indicator */}
                        {selectedEntities.includes(entity.id) && (
                          <div className="absolute -top-2 -right-2 w-6 h-6 bg-primary rounded-full flex items-center justify-center">
                            <svg
                              className="w-4 h-4 text-white"
                              fill="currentColor"
                              viewBox="0 0 20 20"
                            >
                              <path
                                fillRule="evenodd"
                                d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
                                clipRule="evenodd"
                              />
                            </svg>
                          </div>
                        )}

                        <div className="flex items-start gap-3">
                          {/* Icon/Image */}
                          <div className="flex-shrink-0">
                            {entity.type === 'listing' && entity.icon ? (
                              <div className="relative overflow-hidden rounded-lg">
                                <img
                                  src={entity.icon}
                                  alt=""
                                  className="w-16 h-16 object-cover transition-transform hover:scale-110"
                                  onError={(e) => {
                                    const target = e.target as HTMLImageElement;
                                    target.style.display = 'none';
                                    const parent = target.parentElement;
                                    if (parent) {
                                      parent.innerHTML =
                                        '<div class="w-16 h-16 bg-base-200 flex items-center justify-center text-2xl">📋</div>';
                                    }
                                  }}
                                />
                              </div>
                            ) : (
                              <div
                                className={`w-16 h-16 rounded-lg flex items-center justify-center text-2xl transition-all ${
                                  selectedEntities.includes(entity.id)
                                    ? 'bg-primary/10'
                                    : 'bg-base-200'
                                }`}
                              >
                                {entity.icon ||
                                  (entityType === 'category'
                                    ? '📁'
                                    : entityType === 'attribute'
                                      ? '⚙️'
                                      : '📋')}
                              </div>
                            )}
                          </div>

                          {/* Content */}
                          <div className="flex-1 min-w-0">
                            <div className="flex items-start justify-between gap-2">
                              <div className="flex-1">
                                <h4 className="font-medium text-sm leading-tight">
                                  {entity.name}
                                  {entity.parent_name && (
                                    <span className="text-xs text-base-content/60 ml-2">
                                      в {entity.parent_name}
                                    </span>
                                  )}
                                </h4>
                                
                                {/* Translation status badges */}
                                {entity.translations && (
                                  <div className="flex gap-1 mt-1">
                                    {entity.translations.en && (
                                      <span className="badge badge-success badge-xs">EN</span>
                                    )}
                                    {entity.translations.ru && (
                                      <span className="badge badge-success badge-xs">RU</span>
                                    )}
                                    {entity.translations.sr && (
                                      <span className="badge badge-success badge-xs">SR</span>
                                    )}
                                  </div>
                                )}

                                {/* Description */}
                                {entity.description && (
                                  <p className="text-xs text-base-content/60 mt-1 line-clamp-2">
                                    {entity.description}
                                  </p>
                                )}

                                {/* Additional info */}
                                <div className="flex items-center gap-3 mt-2 text-xs text-base-content/50">
                                  {entity.price !== undefined && (
                                    <span className="flex items-center gap-1 font-medium text-success">
                                      <CurrencyEuroIcon className="h-3 w-3" />
                                      {entity.price.toLocaleString()}
                                    </span>
                                  )}
                                  {entity.status && (
                                    <span
                                      className={`badge badge-xs ${
                                        entity.status === 'active'
                                          ? 'badge-success'
                                          : entity.status === 'pending'
                                            ? 'badge-warning'
                                            : entity.status === 'draft'
                                              ? 'badge-info'
                                              : 'badge-ghost'
                                      }`}
                                    >
                                      {entity.status === 'active'
                                        ? '✓ Активно'
                                        : entity.status === 'pending'
                                          ? '⏳ Ожидает'
                                          : entity.status === 'draft'
                                            ? '📝 Черновик'
                                            : entity.status}
                                    </span>
                                  )}
                                  {entity.user_name && (
                                    <span className="flex items-center gap-1">
                                      <UserIcon className="h-3 w-3" />
                                      {entity.user_name}
                                    </span>
                                  )}
                                  {entity.created_at && (
                                    <span className="flex items-center gap-1">
                                      <CalendarIcon className="h-3 w-3" />
                                      {new Date(
                                        entity.created_at
                                      ).toLocaleDateString('ru-RU')}
                                    </span>
                                  )}
                                  <span className="flex items-center gap-1 text-base-content/40">
                                    <TagIcon className="h-3 w-3" />#{entity.id}
                                  </span>
                                </div>
                              </div>

                              {/* Checkbox */}
                              <div className="flex-shrink-0">
                                <input
                                  type="checkbox"
                                  className="checkbox checkbox-primary checkbox-sm"
                                  checked={selectedEntities.includes(entity.id)}
                                  onChange={(e) => {
                                    e.stopPropagation();
                                    handleEntitySelect(entity.id);
                                  }}
                                />
                              </div>
                            </div>
                          </div>
                        </div>
                      </div>
                    ))}
                  </div>
                )}
              </div>

              <div className="text-sm text-base-content/60 mt-4">
                Выбрано: {selectedEntities.length} из {entities.length}
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Progress */}
      {progress > 0 && (
        <div className="card bg-base-100 shadow-sm">
          <div className="card-body">
            <h3 className="card-title text-lg">Прогресс выполнения</h3>
            <div className="flex items-center gap-4">
              <progress
                className="progress progress-primary w-full"
                value={progress}
                max="100"
              ></progress>
              <span className="text-sm font-medium">{progress}%</span>
            </div>
          </div>
        </div>
      )}

      {/* Results */}
      {result && (
        <div className="space-y-4">
          {/* Statistics */}
          <div className="card bg-base-100 shadow-sm">
            <div className="card-body">
              <div className="flex items-center justify-between mb-4">
                <h3 className="card-title text-lg">Результаты перевода</h3>
                <div className="flex items-center gap-4">
                  {result.processing_time && (
                    <span className="text-sm text-base-content/60">
                      Время обработки:{' '}
                      {(result.processing_time / 1000).toFixed(2)} сек
                    </span>
                  )}
                  <button
                    className="btn btn-ghost btn-sm btn-circle"
                    onClick={() => {
                      setResult(null);
                      setProgress(0);
                    }}
                    title="Закрыть результаты"
                  >
                    <XMarkIcon className="h-5 w-5" />
                  </button>
                </div>
              </div>

              <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                <div className="stat">
                  <div className="stat-figure text-primary">
                    <ChartBarIcon className="h-8 w-8 opacity-60" />
                  </div>
                  <div className="stat-title">Обработано</div>
                  <div className="stat-value text-primary">
                    {result.total_processed}
                  </div>
                  <div className="stat-desc">всего элементов</div>
                </div>

                <div className="stat">
                  <div className="stat-figure text-success">
                    <svg
                      className="h-8 w-8 opacity-60"
                      fill="none"
                      viewBox="0 0 24 24"
                      stroke="currentColor"
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={2}
                        d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
                      />
                    </svg>
                  </div>
                  <div className="stat-title">Успешно</div>
                  <div className="stat-value text-success">
                    {result.successful}
                  </div>
                  <div className="stat-desc">
                    {result.successful > 0 &&
                      `${((result.successful / result.total_processed) * 100).toFixed(0)}%`}
                  </div>
                </div>

                <div className="stat">
                  <div className="stat-figure text-error">
                    <svg
                      className="h-8 w-8 opacity-60"
                      fill="none"
                      viewBox="0 0 24 24"
                      stroke="currentColor"
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={2}
                        d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                      />
                    </svg>
                  </div>
                  <div className="stat-title">Ошибки</div>
                  <div className="stat-value text-error">{result.failed}</div>
                  <div className="stat-desc">
                    {result.failed > 0 &&
                      `${((result.failed / result.total_processed) * 100).toFixed(0)}%`}
                  </div>
                </div>

                <div className="stat">
                  <div className="stat-figure text-warning">
                    <svg
                      className="h-8 w-8 opacity-60"
                      fill="none"
                      viewBox="0 0 24 24"
                      stroke="currentColor"
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={2}
                        d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                      />
                    </svg>
                  </div>
                  <div className="stat-title">Пропущено</div>
                  <div className="stat-value text-warning">
                    {result.skipped}
                  </div>
                  <div className="stat-desc">
                    {result.skipped > 0 &&
                      `${((result.skipped / result.total_processed) * 100).toFixed(0)}%`}
                  </div>
                </div>
              </div>

              {result.provider_used && (
                <div className="mt-4 flex items-center gap-2 text-sm text-base-content/60">
                  <CogIcon className="h-4 w-4" />
                  <span>
                    Использован провайдер:{' '}
                    <strong>{result.provider_used}</strong>
                  </span>
                </div>
              )}
            </div>
          </div>

          {/* Successful Items */}
          {result.details?.successful_items &&
            result.details.successful_items.length > 0 && (
              <div className="card bg-success/5 border border-success/20">
                <div className="card-body">
                  <div
                    className="flex items-center justify-between cursor-pointer"
                    onClick={() => toggleSection('successful')}
                  >
                    <h4 className="font-semibold text-success flex items-center gap-2">
                      <svg
                        className="h-5 w-5"
                        fill="none"
                        viewBox="0 0 24 24"
                        stroke="currentColor"
                      >
                        <path
                          strokeLinecap="round"
                          strokeLinejoin="round"
                          strokeWidth={2}
                          d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
                        />
                      </svg>
                      Успешно переведено (
                      {result.details.successful_items.length})
                    </h4>
                    <svg
                      className={`h-5 w-5 transition-transform ${expandedSections.successful ? 'rotate-180' : ''}`}
                      fill="none"
                      viewBox="0 0 24 24"
                      stroke="currentColor"
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={2}
                        d="M19 9l-7 7-7-7"
                      />
                    </svg>
                  </div>
                  {expandedSections.successful && (
                    <div className="max-h-60 overflow-y-auto space-y-2 mt-3">
                      {result.details.successful_items.map((item, index) => (
                        <div
                          key={index}
                          className="p-3 bg-base-100 rounded border border-success/10"
                        >
                          <div className="flex items-center justify-between mb-2">
                            <span className="text-sm font-medium">
                              <strong>#{item.entity_id}</strong>{' '}
                              {item.entity_name}
                            </span>
                            <div className="flex gap-1">
                              {item.languages.map((lang) => (
                                <span
                                  key={lang}
                                  className="badge badge-success badge-sm"
                                >
                                  {lang.toUpperCase()}
                                </span>
                              ))}
                            </div>
                          </div>
                          {/* Показываем переведенные значения */}
                          {item.translations && (
                            <div className="space-y-1 mt-2">
                              {Object.entries(item.translations).map(([lang, text]) => (
                                <div key={lang} className="flex gap-2 text-xs">
                                  <span className="font-semibold text-base-content/70 w-8">
                                    {lang.toUpperCase()}:
                                  </span>
                                  <span className="text-base-content/90 flex-1">
                                    {text as string}
                                  </span>
                                </div>
                              ))}
                            </div>
                          )}
                        </div>
                      ))}
                    </div>
                  )}
                </div>
              </div>
            )}

          {/* Failed Items */}
          {result.details?.failed_items &&
            result.details.failed_items.length > 0 && (
              <div className="card bg-error/5 border border-error/20">
                <div className="card-body">
                  <div className="flex items-center justify-between mb-3">
                    <div
                      className="flex items-center gap-2 cursor-pointer flex-1"
                      onClick={() => toggleSection('failed')}
                    >
                      <h4 className="font-semibold text-error flex items-center gap-2">
                        <svg
                          className="h-5 w-5"
                          fill="none"
                          viewBox="0 0 24 24"
                          stroke="currentColor"
                        >
                          <path
                            strokeLinecap="round"
                            strokeLinejoin="round"
                            strokeWidth={2}
                            d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                          />
                        </svg>
                        Ошибки перевода ({result.details.failed_items.length})
                      </h4>
                      <svg
                        className={`h-5 w-5 transition-transform ${expandedSections.failed ? 'rotate-180' : ''}`}
                        fill="none"
                        viewBox="0 0 24 24"
                        stroke="currentColor"
                      >
                        <path
                          strokeLinecap="round"
                          strokeLinejoin="round"
                          strokeWidth={2}
                          d="M19 9l-7 7-7-7"
                        />
                      </svg>
                    </div>
                    <button
                      className="btn btn-error btn-sm gap-2 ml-4"
                      onClick={retryFailedTranslations}
                      disabled={isLoading}
                    >
                      <ArrowPathIcon className="h-4 w-4" />
                      Повторить
                    </button>
                  </div>
                  {expandedSections.failed && (
                    <div className="max-h-60 overflow-y-auto space-y-2">
                      {result.details.failed_items.map((item, index) => (
                        <div
                          key={index}
                          className="p-3 bg-base-100 rounded border border-error/10"
                        >
                          <div className="flex items-start justify-between mb-1">
                            <span className="text-sm font-medium">
                              <strong>#{item.entity_id}</strong>{' '}
                              {item.entity_name}
                            </span>
                            {item.language && (
                              <span className="badge badge-error badge-sm">
                                {item.language.toUpperCase()}
                              </span>
                            )}
                          </div>
                          <p className="text-xs text-error mt-1">
                            ❌ {item.error}
                          </p>
                        </div>
                      ))}
                    </div>
                  )}
                </div>
              </div>
            )}

          {/* Skipped Items */}
          {result.details?.skipped_items &&
            result.details.skipped_items.length > 0 && (
              <div className="card bg-warning/5 border border-warning/20">
                <div className="card-body">
                  <div
                    className="flex items-center justify-between cursor-pointer"
                    onClick={() => toggleSection('skipped')}
                  >
                    <h4 className="font-semibold text-warning flex items-center gap-2">
                      <svg
                        className="h-5 w-5"
                        fill="none"
                        viewBox="0 0 24 24"
                        stroke="currentColor"
                      >
                        <path
                          strokeLinecap="round"
                          strokeLinejoin="round"
                          strokeWidth={2}
                          d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                        />
                      </svg>
                      Пропущенные элементы (
                      {result.details.skipped_items.length})
                    </h4>
                    <svg
                      className={`h-5 w-5 transition-transform ${expandedSections.skipped ? 'rotate-180' : ''}`}
                      fill="none"
                      viewBox="0 0 24 24"
                      stroke="currentColor"
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={2}
                        d="M19 9l-7 7-7-7"
                      />
                    </svg>
                  </div>
                  {expandedSections.skipped && (
                    <div className="max-h-60 overflow-y-auto space-y-2 mt-3">
                      {result.details.skipped_items.map((item, index) => (
                        <div
                          key={index}
                          className="p-3 bg-base-100 rounded border border-warning/10"
                        >
                          <div className="flex items-start justify-between mb-1">
                            <span className="text-sm font-medium">
                              <strong>#{item.entity_id}</strong>{' '}
                              {item.entity_name}
                            </span>
                            {item.existing_languages &&
                              item.existing_languages.length > 0 && (
                                <div className="flex gap-1">
                                  {item.existing_languages.map((lang) => (
                                    <span
                                      key={lang}
                                      className="badge badge-warning badge-sm"
                                    >
                                      {lang.toUpperCase()}
                                    </span>
                                  ))}
                                </div>
                              )}
                          </div>
                          <p className="text-xs text-warning-content/70 mt-1">
                            ⚠️ Причина: {item.reason}
                          </p>
                        </div>
                      ))}
                    </div>
                  )}
                </div>
              </div>
            )}

          {/* Legacy errors display (if no details provided) */}
          {!result.details && result.errors && result.errors.length > 0 && (
            <div className="card bg-error/5 border border-error/20">
              <div className="card-body">
                <h4 className="font-semibold text-error mb-2">Общие ошибки:</h4>
                <div className="bg-error/10 p-3 rounded-lg">
                  <ul className="list-disc list-inside text-sm space-y-1">
                    {result.errors.map((error, index) => (
                      <li key={index} className="text-error">
                        {error}
                      </li>
                    ))}
                  </ul>
                </div>
              </div>
            </div>
          )}
        </div>
      )}

      {/* Action Button */}
      <div className="card bg-base-100 shadow-sm">
        <div className="card-body">
          <button
            className="btn btn-primary btn-lg w-full gap-2"
            onClick={startBulkTranslation}
            disabled={
              isLoading ||
              selectedEntities.length === 0 ||
              targetLanguages.length === 0
            }
          >
            {isLoading ? (
              <>
                <span className="loading loading-spinner"></span>
                Выполняется перевод...
              </>
            ) : (
              <>
                <PlayIcon className="h-5 w-5" />
                Начать массовый перевод
              </>
            )}
          </button>

          <div className="text-center text-sm text-base-content/60 mt-2">
            {selectedEntities.length > 0 &&
              targetLanguages.length > 0 &&
              `Будет создано ~${selectedEntities.length * targetLanguages.length} переводов`}
          </div>
        </div>
      </div>
    </div>
  );
}
