'use client';

import { useState, useEffect, useRef, useMemo } from 'react';
import { useTranslations } from 'next-intl';
import { useAttributeAutocomplete } from '@/hooks/useAttributeAutocomplete';
import type { components } from '@/types/generated/api';

type UnifiedAttribute =
  components['schemas']['backend_internal_domain_models.UnifiedAttribute'];
type UnifiedAttributeValue =
  components['schemas']['backend_internal_domain_models.UnifiedAttributeValue'];

interface AutocompleteAttributeFieldProps {
  attribute: UnifiedAttribute;
  value?: UnifiedAttributeValue;
  onChange: (value: UnifiedAttributeValue) => void;
  className?: string;
}

interface SuggestionItem {
  value: string;
  label: string;
  type: 'popular' | 'recent' | 'suggestion' | 'exact';
  confidence: number;
}

export function AutocompleteAttributeField({
  attribute,
  value,
  onChange,
  className = '',
}: AutocompleteAttributeFieldProps) {
  const t = useTranslations('common');
  const tFilters = useTranslations('filters');

  const [inputValue, setInputValue] = useState(value?.text_value || '');
  const [showSuggestions, setShowSuggestions] = useState(false);
  const [selectedIndex, setSelectedIndex] = useState(-1);
  const [suggestions, setSuggestions] = useState<SuggestionItem[]>([]);

  const inputRef = useRef<HTMLInputElement>(null);
  const suggestionsRef = useRef<HTMLDivElement>(null);

  // Используем хук для управления автокомплитом
  const { popularValues, recentValues, saveValue, getFilteredSuggestions } =
    useAttributeAutocomplete({
      attributeId: attribute.id!,
      attributeName: attribute.name || 'unknown',
    });

  // Получение умных предложений на основе типа атрибута и популярных значений
  const generateSmartSuggestions = useMemo(() => {
    const attributeName = (attribute.name || '').toLowerCase();
    const smartPatterns: Record<string, string[]> = {
      // Цены
      price: ['50000', '100000', '150000', '200000', '300000', '500000'],
      cost: ['1000', '5000', '10000', '15000', '20000', '30000'],
      цена: ['50000', '100000', '150000', '200000', '300000', '500000'],
      стоимость: ['1000', '5000', '10000', '15000', '20000', '30000'],

      // Годы
      year: ['2024', '2023', '2022', '2021', '2020', '2019', '2018'],
      год: ['2024', '2023', '2022', '2021', '2020', '2019', '2018'],

      // Бренды (общие)
      brand: ['Apple', 'Samsung', 'BMW', 'Mercedes', 'Audi', 'Toyota', 'Honda'],
      марка: ['BMW', 'Mercedes', 'Audi', 'Toyota', 'Honda', 'Volkswagen'],
      бренд: ['Apple', 'Samsung', 'Sony', 'LG', 'Xiaomi', 'Huawei'],

      // Состояние
      condition: ['Новое', 'Отличное', 'Хорошее', 'Удовлетворительное'],
      состояние: ['Новое', 'Отличное', 'Хорошее', 'Удовлетворительное'],

      // Локации (примеры для Сербии)
      location: ['Београд', 'Нови Сад', 'Ниш', 'Крагујевац', 'Суботица'],
      город: ['Белград', 'Новосибирск', 'Екатеринбург', 'Нижний Новгород'],

      // Типы
      type: ['Стандартный', 'Премиум', 'Эконом', 'Люкс'],
      тип: ['Стандартный', 'Премиум', 'Эконом', 'Люкс'],
    };

    // Найти подходящий паттерн
    for (const [pattern, values] of Object.entries(smartPatterns)) {
      if (attributeName.includes(pattern)) {
        return values;
      }
    }

    // Если специфичные паттерны не найдены, используем опции атрибута
    if (attribute.options && Array.isArray(attribute.options)) {
      return attribute.options.map(String);
    }

    return [];
  }, [attribute]);

  // Создание списка предложений с использованием хука
  const createSuggestions = useMemo(() => {
    const hookSuggestions = getFilteredSuggestions(inputValue);
    const suggestions: SuggestionItem[] = [];

    // Преобразуем предложения из хука в формат компонента
    hookSuggestions.forEach(({ value, type }) => {
      let confidence = 0.5;
      let suggestionType: SuggestionItem['type'] = 'suggestion';

      if (
        inputValue.trim() &&
        value.toLowerCase() === inputValue.toLowerCase()
      ) {
        suggestionType = 'exact';
        confidence = 1.0;
      } else if (type === 'popular') {
        suggestionType = 'popular';
        confidence = 0.9;
      } else if (type === 'recent') {
        suggestionType = 'recent';
        confidence = 0.7;
      }

      suggestions.push({
        value,
        label: value,
        type: suggestionType,
        confidence,
      });
    });

    // Добавляем умные предложения если не хватает предложений
    if (suggestions.length < 6 && !inputValue.trim()) {
      generateSmartSuggestions.forEach((val) => {
        if (!suggestions.find((s) => s.value === val)) {
          suggestions.push({
            value: val,
            label: val,
            type: 'suggestion',
            confidence: 0.5,
          });
        }
      });
    }

    return suggestions.slice(0, 8);
  }, [inputValue, getFilteredSuggestions, generateSmartSuggestions]);

  // Обновление предложений при изменении поиска
  useEffect(() => {
    setSuggestions(createSuggestions);
    setSelectedIndex(-1);
  }, [createSuggestions]);

  // Обработка изменения ввода
  const handleInputChange = (newValue: string) => {
    setInputValue(newValue);

    // Создаем объект значения для родительского компонента
    const attributeValue: UnifiedAttributeValue = {
      attribute_id: attribute.id!,
      text_value: newValue.trim(),
    };

    onChange(attributeValue);
  };

  // Выбор предложения
  const selectSuggestion = (suggestion: SuggestionItem) => {
    setInputValue(suggestion.value);
    setShowSuggestions(false);
    setSelectedIndex(-1);

    const attributeValue: UnifiedAttributeValue = {
      attribute_id: attribute.id!,
      text_value: suggestion.value,
    };

    onChange(attributeValue);

    // Сохраняем значение через хук
    saveValue(suggestion.value);
  };

  // Обработка клавиатуры
  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (!showSuggestions || suggestions.length === 0) return;

    switch (e.key) {
      case 'ArrowDown':
        e.preventDefault();
        setSelectedIndex((prev) =>
          prev < suggestions.length - 1 ? prev + 1 : prev
        );
        break;
      case 'ArrowUp':
        e.preventDefault();
        setSelectedIndex((prev) => (prev > 0 ? prev - 1 : -1));
        break;
      case 'Enter':
        e.preventDefault();
        if (selectedIndex >= 0) {
          selectSuggestion(suggestions[selectedIndex]);
        } else {
          setShowSuggestions(false);
        }
        break;
      case 'Escape':
        setShowSuggestions(false);
        setSelectedIndex(-1);
        break;
    }
  };

  // Получение иконки для типа предложения
  const getSuggestionIcon = (type: SuggestionItem['type']) => {
    switch (type) {
      case 'exact':
        return '🎯';
      case 'popular':
        return '⭐';
      case 'recent':
        return '🕒';
      case 'suggestion':
        return '💡';
      default:
        return '';
    }
  };

  // Получение описания типа предложения
  const getSuggestionTypeLabel = (type: SuggestionItem['type']) => {
    switch (type) {
      case 'exact':
        return t('autocomplete.exact_match');
      case 'popular':
        return tFilters('smart_suggestions.most_used');
      case 'recent':
        return t('autocomplete.recently_used');
      case 'suggestion':
        return tFilters('smart_suggestions.recommended');
      default:
        return '';
    }
  };

  return (
    <div className={`form-control relative ${className}`}>
      <label className="label">
        <span className="label-text font-medium">
          {attribute.display_name || attribute.name}
          {attribute.is_required && <span className="text-error"> *</span>}
        </span>
      </label>

      <div className="relative">
        <input
          ref={inputRef}
          type="text"
          value={inputValue}
          onChange={(e) => handleInputChange(e.target.value)}
          onFocus={() => setShowSuggestions(true)}
          onBlur={(e) => {
            // Задержка для обработки клика по предложению
            setTimeout(() => {
              if (!suggestionsRef.current?.contains(e.relatedTarget as Node)) {
                setShowSuggestions(false);
              }
            }, 150);
          }}
          onKeyDown={handleKeyDown}
          placeholder={attribute.description || t('autocomplete.enter_value')}
          className={`input input-bordered w-full pr-10 ${className.includes('has-error') ? 'input-error' : ''}`}
        />

        {/* Иконка поиска */}
        <div className="absolute inset-y-0 right-0 pr-3 flex items-center pointer-events-none">
          <svg
            className="h-4 w-4 text-gray-400"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
            />
          </svg>
        </div>

        {/* Список предложений */}
        {showSuggestions && suggestions.length > 0 && (
          <div
            ref={suggestionsRef}
            className="absolute z-50 w-full mt-1 bg-base-100 border border-base-300 rounded-lg shadow-lg max-h-60 overflow-y-auto"
          >
            {suggestions.map((suggestion, index) => (
              <div
                key={`${suggestion.value}-${suggestion.type}`}
                className={`px-4 py-2 cursor-pointer flex items-center justify-between hover:bg-base-200 ${
                  index === selectedIndex
                    ? 'bg-primary text-primary-content'
                    : ''
                }`}
                onClick={() => selectSuggestion(suggestion)}
              >
                <div className="flex items-center gap-2 flex-1">
                  <span className="text-lg">
                    {getSuggestionIcon(suggestion.type)}
                  </span>
                  <span className="font-medium">{suggestion.label}</span>
                </div>
                <div className="text-xs opacity-70">
                  {getSuggestionTypeLabel(suggestion.type)}
                </div>
              </div>
            ))}

            {/* Подсказка по управлению */}
            <div className="px-4 py-2 border-t border-base-300 bg-base-50">
              <div className="text-xs text-base-content opacity-60">
                ↑↓ {t('autocomplete.navigate')}, Enter {t('select')}, Esc{' '}
                {t('close')}
              </div>
            </div>
          </div>
        )}
      </div>

      {/* Описание атрибута */}
      {attribute.description && (
        <label className="label">
          <span className="label-text-alt opacity-70">
            {attribute.description}
          </span>
        </label>
      )}
    </div>
  );
}
