'use client';

import { useState, useEffect, useMemo, useCallback, memo } from 'react';
import { useTranslations } from 'next-intl';
import { UnifiedAttributeField } from './UnifiedAttributeField';
import IntuitiveAttributeField from './IntuitiveAttributeField';
import { unifiedAttributeService } from '@/services/unifiedAttributeService';
import type { components } from '@/types/generated/api';

type UnifiedAttribute =
  components['schemas']['backend_internal_domain_models.UnifiedAttribute'];
type UnifiedAttributeValue =
  components['schemas']['backend_internal_domain_models.UnifiedAttributeValue'];

interface SmartAttributeFiltersProps {
  categoryId?: number;
  onFiltersChange: (filters: Record<number, UnifiedAttributeValue>) => void;
  initialFilters?: Record<number, UnifiedAttributeValue>;
  className?: string;
}

interface FilterGroup {
  id: string;
  name: string;
  icon: string;
  attributes: UnifiedAttribute[];
  priority: number;
  isPopular?: boolean;
}

interface SmartSuggestion {
  attributeId: number;
  value: string;
  confidence: number;
  reason: string;
}

function SmartAttributeFiltersComponent({
  categoryId,
  onFiltersChange,
  initialFilters = {},
  className = '',
}: SmartAttributeFiltersProps) {
  const t = useTranslations('common');
  const tFilters = useTranslations('filters');

  const [attributes, setAttributes] = useState<UnifiedAttribute[]>([]);
  const [filters, setFilters] =
    useState<Record<number, UnifiedAttributeValue>>(initialFilters);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [searchQuery, setSearchQuery] = useState('');
  const [expandedGroups, setExpandedGroups] = useState<Set<string>>(
    new Set(['popular', 'basic'])
  );
  const [smartSuggestions, setSmartSuggestions] = useState<SmartSuggestion[]>(
    []
  );
  const [showAdvanced, setShowAdvanced] = useState(false);

  // Загрузка атрибутов для категории
  useEffect(() => {
    const loadAttributes = async () => {
      if (!categoryId) {
        setLoading(false);
        return;
      }

      setLoading(true);
      setError(null);

      try {
        const response =
          await unifiedAttributeService.getCategoryAttributes(categoryId);

        if (response.success && response.data) {
          const activeAttributes = response.data.filter(
            (attr) => attr.is_active !== false && attr.is_filterable !== false
          );
          setAttributes(activeAttributes);
          await generateSmartSuggestions(activeAttributes);
        } else {
          throw new Error(response.error || 'Failed to load attributes');
        }
      } catch (err) {
        console.error('Error loading attributes:', err);
        setError(t('load_error'));
      } finally {
        setLoading(false);
      }
    };

    loadAttributes();
  }, [categoryId, t]);

  // Генерация умных предложений на основе популярности и релевантности
  const generateSmartSuggestions = async (attrs: UnifiedAttribute[]) => {
    const suggestions: SmartSuggestion[] = [];

    // Логика определения популярных фильтров на основе названий атрибутов
    const popularPatterns = [
      {
        pattern: /price|cost|стоимость|цена/i,
        reason: 'Most users filter by price',
        confidence: 0.95,
      },
      {
        pattern: /year|год/i,
        reason: 'Year is commonly filtered',
        confidence: 0.9,
      },
      {
        pattern: /brand|make|марка|бренд/i,
        reason: 'Brand filtering is very popular',
        confidence: 0.88,
      },
      {
        pattern: /condition|состояние/i,
        reason: 'Condition is important for buyers',
        confidence: 0.85,
      },
      {
        pattern: /type|тип|вид/i,
        reason: 'Type helps narrow search',
        confidence: 0.8,
      },
      {
        pattern: /location|город|место/i,
        reason: 'Location is crucial for search',
        confidence: 0.9,
      },
    ];

    attrs.forEach((attr) => {
      if (!attr.id || !attr.name) return;

      const name = attr.name.toLowerCase();
      const matchedPattern = popularPatterns.find((p) => p.pattern.test(name));

      if (
        matchedPattern &&
        attr.attribute_type === 'select' &&
        attr.options?.length
      ) {
        // Для select полей предлагаем популярные опции
        if (Array.isArray(attr.options) && attr.options.length > 0) {
          const topOption = attr.options[0]; // Можно улучшить логику выбора
          suggestions.push({
            attributeId: attr.id,
            value: String(topOption),
            confidence: matchedPattern.confidence,
            reason: matchedPattern.reason,
          });
        }
      } else if (matchedPattern && attr.attribute_type === 'number') {
        // Для числовых полей не предлагаем конкретные значения
        suggestions.push({
          attributeId: attr.id,
          value: '',
          confidence: matchedPattern.confidence,
          reason: matchedPattern.reason,
        });
      }
    });

    // Сортируем по confidence
    suggestions.sort((a, b) => b.confidence - a.confidence);
    setSmartSuggestions(suggestions.slice(0, 5)); // Топ 5 предложений
  };

  // Фильтрация атрибутов по поисковому запросу
  const filteredAttributes = useMemo(() => {
    if (!searchQuery.trim()) return attributes;

    const query = searchQuery.toLowerCase();
    return attributes.filter(
      (attr) =>
        attr.name?.toLowerCase().includes(query) ||
        attr.display_name?.toLowerCase().includes(query)
    );
  }, [attributes, searchQuery]);

  // Умная группировка атрибутов
  const groupedAttributes = useMemo((): FilterGroup[] => {
    const groups = new Map<string, FilterGroup>();

    // Получаем ID популярных атрибутов из suggestions
    const popularAttributeIds = new Set(
      smartSuggestions.map((s) => s.attributeId)
    );

    const predefinedGroups = {
      popular: {
        name: tFilters('groups.popular'),
        icon: '⭐',
        priority: 1,
        isPopular: true,
      },
      price: { name: tFilters('groups.price'), icon: '💰', priority: 2 },
      basic: { name: tFilters('groups.basic'), icon: '🏷️', priority: 3 },
      technical: {
        name: tFilters('groups.technical'),
        icon: '⚙️',
        priority: 4,
      },
      condition: {
        name: tFilters('groups.condition'),
        icon: '✨',
        priority: 5,
      },
      location: { name: tFilters('groups.location'), icon: '📍', priority: 6 },
      other: { name: tFilters('groups.other'), icon: '📋', priority: 7 },
    };

    filteredAttributes.forEach((attr) => {
      if (!attr.name) return;

      let groupId = 'other';
      const name = attr.name.toLowerCase();

      // Сначала проверяем популярные атрибуты
      if (popularAttributeIds.has(attr.id!)) {
        groupId = 'popular';
      } else if (
        ['price', 'cost', 'стоимость', 'цена'].some((key) => name.includes(key))
      ) {
        groupId = 'price';
      } else if (
        ['brand', 'model', 'type', 'category', 'name', 'марка', 'модель'].some(
          (key) => name.includes(key)
        )
      ) {
        groupId = 'basic';
      } else if (
        [
          'year',
          'engine',
          'fuel',
          'transmission',
          'power',
          'год',
          'двигатель',
        ].some((key) => name.includes(key))
      ) {
        groupId = 'technical';
      } else if (
        ['condition', 'warranty', 'used', 'new', 'состояние', 'гарантия'].some(
          (key) => name.includes(key)
        )
      ) {
        groupId = 'condition';
      } else if (
        ['location', 'city', 'region', 'address', 'город', 'регион'].some(
          (key) => name.includes(key)
        )
      ) {
        groupId = 'location';
      }

      if (!groups.has(groupId)) {
        const groupInfo =
          predefinedGroups[groupId as keyof typeof predefinedGroups] ||
          predefinedGroups.other;
        groups.set(groupId, {
          id: groupId,
          name: groupInfo.name,
          icon: groupInfo.icon,
          priority: groupInfo.priority,
          isPopular: 'isPopular' in groupInfo ? groupInfo.isPopular : false,
          attributes: [],
        });
      }

      groups.get(groupId)!.attributes.push(attr);
    });

    // Сортировка групп по приоритету и атрибутов внутри групп
    return Array.from(groups.values())
      .sort((a, b) => a.priority - b.priority)
      .map((group) => ({
        ...group,
        attributes: group.attributes.sort((a, b) => {
          // В популярной группе сортируем по confidence из suggestions
          if (group.isPopular) {
            const aConfidence =
              smartSuggestions.find((s) => s.attributeId === a.id)
                ?.confidence || 0;
            const bConfidence =
              smartSuggestions.find((s) => s.attributeId === b.id)
                ?.confidence || 0;
            return bConfidence - aConfidence;
          }
          return (a.sort_order || 0) - (b.sort_order || 0);
        }),
      }))
      .filter((group) => group.attributes.length > 0); // Убираем пустые группы
  }, [filteredAttributes, smartSuggestions, tFilters]);

  // Обработчик изменения фильтра
  const handleFilterChange = useCallback(
    (attributeId: number, value: UnifiedAttributeValue) => {
      const newFilters = { ...filters };

      // Если значение пустое, удаляем фильтр
      const isEmpty =
        !value.text_value &&
        value.numeric_value === undefined &&
        value.boolean_value === undefined &&
        !value.date_value &&
        (!value.json_value ||
          (Array.isArray(value.json_value) && value.json_value.length === 0));

      if (isEmpty) {
        delete newFilters[attributeId];
      } else {
        newFilters[attributeId] = value;
      }

      setFilters(newFilters);
      onFiltersChange(newFilters);
    },
    [filters, onFiltersChange]
  );

  // Применить умное предложение
  const applySuggestion = (suggestion: SmartSuggestion) => {
    const attr = attributes.find((a) => a.id === suggestion.attributeId);
    if (!attr || !attr.id) return;

    const value: UnifiedAttributeValue = {
      attribute_id: attr.id,
      text_value: suggestion.value,
    };

    handleFilterChange(attr.id, value);
  };

  // Очистить все фильтры
  const clearAllFilters = () => {
    setFilters({});
    onFiltersChange({});
  };

  // Переключение группы
  const toggleGroup = (groupId: string) => {
    setExpandedGroups((prev) => {
      const next = new Set(prev);
      if (next.has(groupId)) {
        next.delete(groupId);
      } else {
        next.add(groupId);
      }
      return next;
    });
  };

  const activeFiltersCount = Object.keys(filters).length;

  if (loading) {
    return (
      <div className="flex items-center justify-center py-8">
        <div className="loading loading-spinner loading-lg"></div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="alert alert-error">
        <span>{error}</span>
      </div>
    );
  }

  return (
    <div className={`space-y-4 ${className}`}>
      {/* Заголовок с поиском */}
      <div className="flex items-center justify-between">
        <h3 className="text-lg font-semibold flex items-center gap-2">
          🔍 {tFilters('title')}
          {activeFiltersCount > 0 && (
            <div className="badge badge-primary">{activeFiltersCount}</div>
          )}
        </h3>

        <div className="flex items-center gap-2">
          {activeFiltersCount > 0 && (
            <button
              onClick={clearAllFilters}
              className="btn btn-sm btn-outline btn-error"
            >
              {tFilters('clear_all')}
            </button>
          )}
          <button
            onClick={() => setShowAdvanced(!showAdvanced)}
            className="btn btn-sm btn-outline"
          >
            {showAdvanced
              ? tFilters('hide_advanced')
              : tFilters('show_advanced')}
          </button>
        </div>
      </div>

      {/* Поиск по атрибутам */}
      <div className="form-control">
        <input
          type="text"
          placeholder={tFilters('search_placeholder')}
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          className="input input-bordered"
        />
      </div>

      {/* Умные предложения */}
      {smartSuggestions.length > 0 && activeFiltersCount === 0 && (
        <div className="alert alert-info">
          <div className="w-full">
            <h4 className="font-bold mb-2">
              💡 {tFilters('smart_suggestions.title')}
            </h4>
            <div className="space-y-2">
              {smartSuggestions.slice(0, 3).map((suggestion) => {
                const attr = attributes.find(
                  (a) => a.id === suggestion.attributeId
                );
                if (!attr) return null;

                return (
                  <div
                    key={suggestion.attributeId}
                    className="flex items-center justify-between"
                  >
                    <span className="text-sm">
                      {attr.display_name || attr.name} - {suggestion.reason}
                    </span>
                    <button
                      onClick={() => applySuggestion(suggestion)}
                      className="btn btn-xs btn-primary"
                    >
                      {tFilters('smart_suggestions.apply')}
                    </button>
                  </div>
                );
              })}
            </div>
          </div>
        </div>
      )}

      {/* Группы фильтров */}
      <div className="space-y-3">
        {groupedAttributes.map((group) => {
          const isExpanded = expandedGroups.has(group.id);
          const groupHasActiveFilters = group.attributes.some(
            (attr) => attr.id && filters[attr.id]
          );

          // Показываем только популярные и базовые группы по умолчанию
          if (
            !showAdvanced &&
            !['popular', 'basic', 'price'].includes(group.id)
          ) {
            return null;
          }

          return (
            <div key={group.id} className="card bg-base-100 shadow-sm border">
              <div
                className="card-body py-3 cursor-pointer"
                onClick={() => toggleGroup(group.id)}
              >
                <div className="flex items-center justify-between">
                  <h4 className="text-md font-medium flex items-center gap-2">
                    <span>{group.icon}</span>
                    {group.name}
                    <div className="badge badge-neutral badge-sm">
                      {group.attributes.length}
                    </div>
                    {groupHasActiveFilters && (
                      <div className="badge badge-primary badge-sm">✓</div>
                    )}
                  </h4>
                  <svg
                    className={`w-4 h-4 transition-transform ${isExpanded ? 'rotate-180' : ''}`}
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M19 9l-7 7-7-7"
                    />
                  </svg>
                </div>
              </div>

              {isExpanded && (
                <div className="card-body pt-0">
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    {group.attributes.map((attr) => {
                      if (!attr.id) return null;

                      return (
                        <IntuitiveAttributeField
                          key={attr.id}
                          attribute={attr}
                          value={filters[attr.id]}
                          onChange={(value) =>
                            handleFilterChange(attr.id!, value)
                          }
                          className="form-control-sm"
                          enableAutocomplete={true}
                        />
                      );
                    })}
                  </div>
                </div>
              )}
            </div>
          );
        })}
      </div>

      {/* Счетчик результатов (placeholder для будущей интеграции) */}
      {activeFiltersCount > 0 && (
        <div className="alert alert-success">
          <span>
            🎯 {tFilters('filters_applied', { count: activeFiltersCount })}
          </span>
        </div>
      )}
    </div>
  );
}

// Мемоизированный экспорт для оптимизации производительности
export const SmartAttributeFilters = memo(
  SmartAttributeFiltersComponent,
  (prevProps, nextProps) => {
    // Пользовательское сравнение для оптимизации перерендеринга
    return (
      prevProps.categoryId === nextProps.categoryId &&
      prevProps.className === nextProps.className &&
      JSON.stringify(prevProps.initialFilters) ===
        JSON.stringify(nextProps.initialFilters) &&
      prevProps.onFiltersChange === nextProps.onFiltersChange
    );
  }
);

export default SmartAttributeFilters;
